package database

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migration represents a single database migration
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// Migrate runs all pending database migrations
func (db *DB) Migrate() error {
	// Create migrations table if it doesn't exist
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all available migrations
	migrations, err := db.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedVersions, err := db.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Filter out already applied migrations
	pendingMigrations := filterPendingMigrations(migrations, appliedVersions)

	if len(pendingMigrations) == 0 {
		return nil // No pending migrations
	}

	// Apply pending migrations
	for _, migration := range pendingMigrations {
		if err := db.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.conn.Exec(query)
	return err
}

// loadMigrations loads all migration files from the embedded filesystem
func (db *DB) loadMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Parse version from filename (e.g., "001_initial_schema.sql")
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		// Read migration content
		content, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		// Extract name from filename (remove version and extension)
		name := strings.TrimSuffix(parts[1], ".sql")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			SQL:     string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a set of applied migration versions
func (db *DB) getAppliedMigrations() (map[int]bool, error) {
	query := "SELECT version FROM migrations ORDER BY version"
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// filterPendingMigrations returns migrations that haven't been applied yet
func filterPendingMigrations(migrations []Migration, applied map[int]bool) []Migration {
	var pending []Migration
	for _, migration := range migrations {
		if !applied[migration.Version] {
			pending = append(pending, migration)
		}
	}
	return pending
}

// applyMigration applies a single migration within a transaction
func (db *DB) applyMigration(migration Migration) error {
	return db.ExecTx(func(tx *sql.Tx) error {
		// Execute the migration SQL
		if _, err := tx.Exec(migration.SQL); err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}

		// Record the migration as applied
		query := "INSERT INTO migrations (version, name) VALUES (?, ?)"
		if _, err := tx.Exec(query, migration.Version, migration.Name); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})
}

// GetMigrationStatus returns the current migration status
func (db *DB) GetMigrationStatus() ([]MigrationStatus, error) {
	// Get all available migrations
	migrations, err := db.loadMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations with timestamps
	appliedMigrations, err := db.getAppliedMigrationsWithTimestamp()
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	var status []MigrationStatus
	for _, migration := range migrations {
		migrationStatus := MigrationStatus{
			Version: migration.Version,
			Name:    migration.Name,
			Applied: false,
		}

		if appliedInfo, exists := appliedMigrations[migration.Version]; exists {
			migrationStatus.Applied = true
			migrationStatus.AppliedAt = &appliedInfo.AppliedAt
		}

		status = append(status, migrationStatus)
	}

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   int    `json:"version"`
	Name      string `json:"name"`
	Applied   bool   `json:"applied"`
	AppliedAt *string `json:"applied_at,omitempty"`
}

// AppliedMigrationInfo contains information about an applied migration
type AppliedMigrationInfo struct {
	Version   int
	Name      string
	AppliedAt string
}

// getAppliedMigrationsWithTimestamp returns applied migrations with their timestamps
func (db *DB) getAppliedMigrationsWithTimestamp() (map[int]AppliedMigrationInfo, error) {
	query := "SELECT version, name, applied_at FROM migrations ORDER BY version"
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]AppliedMigrationInfo)
	for rows.Next() {
		var info AppliedMigrationInfo
		if err := rows.Scan(&info.Version, &info.Name, &info.AppliedAt); err != nil {
			return nil, err
		}
		applied[info.Version] = info
	}

	return applied, rows.Err()
}

// ResetDatabase drops all tables and re-runs migrations (USE WITH CAUTION)
func (db *DB) ResetDatabase() error {
	return db.ExecTx(func(tx *sql.Tx) error {
		// Get all table names
		tables, err := db.GetTableInfo()
		if err != nil {
			return fmt.Errorf("failed to get table info: %w", err)
		}

		// Drop all tables
		for _, table := range tables {
			if _, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}

		return nil
	})
}

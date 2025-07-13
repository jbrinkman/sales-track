package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB represents the database connection and configuration
type DB struct {
	conn     *sql.DB
	filePath string
}

// Config represents database configuration options
type Config struct {
	FilePath    string // Path to SQLite database file
	InMemory    bool   // Use in-memory database for testing
	AutoMigrate bool   // Automatically run migrations on startup
}

// New creates a new database connection with the given configuration
func New(config Config) (*DB, error) {
	var dsn string
	var filePath string

	if config.InMemory {
		// Use in-memory database for testing
		dsn = ":memory:"
		filePath = ":memory:"
	} else {
		// Use file-based database
		if config.FilePath == "" {
			// Default to sales_track.db in current directory
			config.FilePath = "sales_track.db"
		}

		// Ensure directory exists
		dir := filepath.Dir(config.FilePath)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}
		}

		dsn = config.FilePath
		filePath = config.FilePath
	}

	// Open database connection
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure SQLite settings
	if err := configureSQLite(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to configure SQLite: %w", err)
	}

	db := &DB{
		conn:     conn,
		filePath: filePath,
	}

	// Run migrations if requested
	if config.AutoMigrate {
		if err := db.Migrate(); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return db, nil
}

// configureSQLite sets up SQLite-specific configuration
func configureSQLite(conn *sql.DB) error {
	// Enable foreign key constraints
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Set journal mode to WAL for better concurrency
	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return fmt.Errorf("failed to set journal mode: %w", err)
	}

	// Set synchronous mode to NORMAL for better performance
	if _, err := conn.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		return fmt.Errorf("failed to set synchronous mode: %w", err)
	}

	// Set cache size (negative value means KB, positive means pages)
	if _, err := conn.Exec("PRAGMA cache_size = -64000"); err != nil { // 64MB cache
		return fmt.Errorf("failed to set cache size: %w", err)
	}

	// Set temp store to memory for better performance
	if _, err := conn.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		return fmt.Errorf("failed to set temp store: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Conn returns the underlying sql.DB connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// FilePath returns the database file path
func (db *DB) FilePath() string {
	return db.filePath
}

// Ping tests the database connection
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// Stats returns database connection statistics
func (db *DB) Stats() sql.DBStats {
	return db.conn.Stats()
}

// BeginTx starts a new transaction with the given options
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}

// ExecTx executes a function within a transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (db *DB) ExecTx(fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// IsHealthy checks if the database connection is healthy
func (db *DB) IsHealthy() bool {
	if db.conn == nil {
		return false
	}

	// Test with a simple query
	var result int
	err := db.conn.QueryRow("SELECT 1").Scan(&result)
	return err == nil && result == 1
}

// GetVersion returns the SQLite version
func (db *DB) GetVersion() (string, error) {
	var version string
	err := db.conn.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to get SQLite version: %w", err)
	}
	return version, nil
}

// GetTableInfo returns information about database tables
func (db *DB) GetTableInfo() ([]string, error) {
	rows, err := db.conn.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query table info: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table rows: %w", err)
	}

	return tables, nil
}

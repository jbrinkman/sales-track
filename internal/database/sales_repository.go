package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"sales-track/internal/models"
)

// SalesRepository handles database operations for sales records
type SalesRepository struct {
	db *DB
}

// NewSalesRepository creates a new sales repository
func NewSalesRepository(db *DB) *SalesRepository {
	return &SalesRepository{db: db}
}

// Create inserts a new sales record into the database
func (r *SalesRepository) Create(record models.CreateSalesRecordRequest) (*models.SalesRecord, error) {
	// Parse the date string
	date, err := time.Parse("2006-01-02", record.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	query := `
		INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.conn.Exec(query,
		record.Store,
		record.Vendor,
		date,
		record.Description,
		record.SalePrice,
		record.Commission,
		record.Remaining,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert sales record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	// Fetch and return the created record
	return r.GetByID(id)
}

// GetByID retrieves a sales record by its ID
func (r *SalesRepository) GetByID(id int64) (*models.SalesRecord, error) {
	query := `
		SELECT id, store, vendor, date, description, sale_price, commission, remaining, created_at, updated_at
		FROM sales_records
		WHERE id = ?
	`

	var record models.SalesRecord
	err := r.db.conn.QueryRow(query, id).Scan(
		&record.ID,
		&record.Store,
		&record.Vendor,
		&record.Date,
		&record.Description,
		&record.SalePrice,
		&record.Commission,
		&record.Remaining,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sales record with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get sales record: %w", err)
	}

	return &record, nil
}

// Update updates an existing sales record
func (r *SalesRepository) Update(id int64, updates models.UpdateSalesRecordRequest) (*models.SalesRecord, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}

	if updates.Store != nil {
		setParts = append(setParts, "store = ?")
		args = append(args, *updates.Store)
	}
	if updates.Vendor != nil {
		setParts = append(setParts, "vendor = ?")
		args = append(args, *updates.Vendor)
	}
	if updates.Date != nil {
		date, err := time.Parse("2006-01-02", *updates.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		setParts = append(setParts, "date = ?")
		args = append(args, date)
	}
	if updates.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *updates.Description)
	}
	if updates.SalePrice != nil {
		setParts = append(setParts, "sale_price = ?")
		args = append(args, *updates.SalePrice)
	}
	if updates.Commission != nil {
		setParts = append(setParts, "commission = ?")
		args = append(args, *updates.Commission)
	}
	if updates.Remaining != nil {
		setParts = append(setParts, "remaining = ?")
		args = append(args, *updates.Remaining)
	}

	if len(setParts) == 0 {
		return r.GetByID(id) // No updates, return existing record
	}

	// Add updated_at timestamp
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id) // Add ID for WHERE clause

	query := fmt.Sprintf("UPDATE sales_records SET %s WHERE id = ?", strings.Join(setParts, ", "))

	_, err := r.db.conn.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update sales record: %w", err)
	}

	// Return updated record
	return r.GetByID(id)
}

// Delete removes a sales record from the database
func (r *SalesRepository) Delete(id int64) error {
	query := "DELETE FROM sales_records WHERE id = ?"
	result, err := r.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sales record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sales record with ID %d not found", id)
	}

	return nil
}

// List retrieves sales records with optional filtering and pagination
func (r *SalesRepository) List(filter models.SalesRecordFilter) (*models.SalesRecordList, error) {
	// Build WHERE clause
	whereParts := []string{}
	args := []interface{}{}

	if filter.Store != nil {
		whereParts = append(whereParts, "store = ?")
		args = append(args, *filter.Store)
	}
	if filter.Vendor != nil {
		whereParts = append(whereParts, "vendor = ?")
		args = append(args, *filter.Vendor)
	}
	if filter.DateFrom != nil {
		whereParts = append(whereParts, "date >= ?")
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil {
		whereParts = append(whereParts, "date <= ?")
		args = append(args, *filter.DateTo)
	}
	if filter.MinPrice != nil {
		whereParts = append(whereParts, "sale_price >= ?")
		args = append(args, *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		whereParts = append(whereParts, "sale_price <= ?")
		args = append(args, *filter.MaxPrice)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Build ORDER BY clause
	orderBy := "ORDER BY date DESC" // Default sort
	if filter.SortBy != nil && filter.SortOrder != nil {
		validSortFields := map[string]bool{
			"date":       true,
			"store":      true,
			"vendor":     true,
			"sale_price": true,
			"created_at": true,
		}
		validSortOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		if validSortFields[*filter.SortBy] && validSortOrders[*filter.SortOrder] {
			orderBy = fmt.Sprintf("ORDER BY %s %s", *filter.SortBy, strings.ToUpper(*filter.SortOrder))
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sales_records %s", whereClause)
	var total int64
	err := r.db.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Build LIMIT and OFFSET
	limit := 50 // Default limit
	if filter.Limit != nil && *filter.Limit > 0 {
		limit = *filter.Limit
	}

	offset := 0
	if filter.Offset != nil && *filter.Offset > 0 {
		offset = *filter.Offset
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT id, store, vendor, date, description, sale_price, commission, remaining, created_at, updated_at
		FROM sales_records
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	queryArgs := append(args, limit, offset)
	rows, err := r.db.conn.Query(query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sales records: %w", err)
	}
	defer rows.Close()

	var records []models.SalesRecord
	for rows.Next() {
		var record models.SalesRecord
		err := rows.Scan(
			&record.ID,
			&record.Store,
			&record.Vendor,
			&record.Date,
			&record.Description,
			&record.SalePrice,
			&record.Commission,
			&record.Remaining,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales record: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sales records: %w", err)
	}

	// Calculate pagination info
	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division

	return &models.SalesRecordList{
		Records:    records,
		Total:      total,
		Page:       page,
		PageSize:   limit,
		TotalPages: totalPages,
	}, nil
}

// CreateBatch inserts multiple sales records in a single transaction
func (r *SalesRepository) CreateBatch(records []models.CreateSalesRecordRequest) ([]models.SalesRecord, error) {
	var createdRecords []models.SalesRecord

	err := r.db.ExecTx(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(`
			INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, record := range records {
			// Parse the date string
			date, err := time.Parse("2006-01-02", record.Date)
			if err != nil {
				return fmt.Errorf("invalid date format for record: %w", err)
			}

			result, err := stmt.Exec(
				record.Store,
				record.Vendor,
				date,
				record.Description,
				record.SalePrice,
				record.Commission,
				record.Remaining,
			)
			if err != nil {
				return fmt.Errorf("failed to insert sales record: %w", err)
			}

			id, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get last insert ID: %w", err)
			}

			// Fetch the created record (within the transaction)
			var createdRecord models.SalesRecord
			err = tx.QueryRow(`
				SELECT id, store, vendor, date, description, sale_price, commission, remaining, created_at, updated_at
				FROM sales_records WHERE id = ?
			`, id).Scan(
				&createdRecord.ID,
				&createdRecord.Store,
				&createdRecord.Vendor,
				&createdRecord.Date,
				&createdRecord.Description,
				&createdRecord.SalePrice,
				&createdRecord.Commission,
				&createdRecord.Remaining,
				&createdRecord.CreatedAt,
				&createdRecord.UpdatedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to fetch created record: %w", err)
			}

			createdRecords = append(createdRecords, createdRecord)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdRecords, nil
}

// GetStats returns basic statistics about the sales records
func (r *SalesRepository) GetStats() (*models.DatabaseStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_records,
			COALESCE(MIN(date), '') as earliest_date,
			COALESCE(MAX(date), '') as latest_date,
			COALESCE(SUM(sale_price), 0) as total_sales,
			COALESCE(AVG(sale_price), 0) as avg_sale_price,
			COUNT(DISTINCT store) as unique_stores,
			COUNT(DISTINCT vendor) as unique_vendors,
			COALESCE(MAX(updated_at), '') as last_updated
		FROM sales_records
	`

	var stats models.DatabaseStats
	var earliestDateStr, latestDateStr, lastUpdatedStr string
	
	err := r.db.conn.QueryRow(query).Scan(
		&stats.TotalRecords,
		&earliestDateStr,
		&latestDateStr,
		&stats.TotalSales,
		&stats.AvgSalePrice,
		&stats.UniqueStores,
		&stats.UniqueVendors,
		&lastUpdatedStr,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get database stats: %w", err)
	}

	// Parse date strings
	if earliestDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", earliestDateStr); err == nil {
			stats.EarliestDate = parsed
		}
	}
	if latestDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", latestDateStr); err == nil {
			stats.LatestDate = parsed
		}
	}
	if lastUpdatedStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", lastUpdatedStr); err == nil {
			stats.LastUpdated = parsed
		}
	}

	return &stats, nil
}

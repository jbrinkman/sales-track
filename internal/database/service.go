package database

import (
	"database/sql"
	"fmt"

	"sales-track/internal/models"
)

// Service provides high-level database operations
// It combines multiple repositories and provides a unified API
type Service struct {
	db                *DB
	salesRepo         *SalesRepository
	reportingRepo     *ReportingRepository
}

// NewService creates a new database service
func NewService(config Config) (*Service, error) {
	db, err := New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return &Service{
		db:                db,
		salesRepo:         NewSalesRepository(db),
		reportingRepo:     NewReportingRepository(db),
	}, nil
}

// Close closes the database connection
func (s *Service) Close() error {
	return s.db.Close()
}

// Health checks if the database service is healthy
func (s *Service) Health() error {
	if !s.db.IsHealthy() {
		return fmt.Errorf("database connection is not healthy")
	}
	return nil
}

// GetDB returns the underlying database connection (for advanced usage)
func (s *Service) GetDB() *DB {
	return s.db
}

// ===== SALES RECORD OPERATIONS =====

// CreateSalesRecord creates a new sales record
func (s *Service) CreateSalesRecord(record models.CreateSalesRecordRequest) (*models.SalesRecord, error) {
	return s.salesRepo.Create(record)
}

// GetSalesRecord retrieves a sales record by ID
func (s *Service) GetSalesRecord(id int64) (*models.SalesRecord, error) {
	return s.salesRepo.GetByID(id)
}

// UpdateSalesRecord updates an existing sales record
func (s *Service) UpdateSalesRecord(id int64, updates models.UpdateSalesRecordRequest) (*models.SalesRecord, error) {
	return s.salesRepo.Update(id, updates)
}

// DeleteSalesRecord removes a sales record
func (s *Service) DeleteSalesRecord(id int64) error {
	return s.salesRepo.Delete(id)
}

// ListSalesRecords retrieves sales records with filtering and pagination
func (s *Service) ListSalesRecords(filter models.SalesRecordFilter) (*models.SalesRecordList, error) {
	return s.salesRepo.List(filter)
}

// CreateSalesRecordsBatch creates multiple sales records in a single transaction
func (s *Service) CreateSalesRecordsBatch(records []models.CreateSalesRecordRequest) ([]models.SalesRecord, error) {
	return s.salesRepo.CreateBatch(records)
}

// GetDatabaseStats returns overall database statistics
func (s *Service) GetDatabaseStats() (*models.DatabaseStats, error) {
	return s.salesRepo.GetStats()
}

// ===== REPORTING OPERATIONS =====

// GetYearlySummary returns yearly sales summary
func (s *Service) GetYearlySummary() ([]models.YearlySummary, error) {
	return s.reportingRepo.GetYearlySummary()
}

// GetMonthlySummary returns monthly sales summary, optionally filtered by year
func (s *Service) GetMonthlySummary(year *string) ([]models.MonthlySummary, error) {
	return s.reportingRepo.GetMonthlySummary(year)
}

// GetDailySummary returns daily sales summary, optionally filtered by year and month
func (s *Service) GetDailySummary(year *string, month *string) ([]models.DailySummary, error) {
	return s.reportingRepo.GetDailySummary(year, month)
}

// GetStorePerformance returns store performance analytics
func (s *Service) GetStorePerformance() ([]models.StorePerformance, error) {
	return s.reportingRepo.GetStorePerformance()
}

// GetVendorPerformance returns vendor performance analytics
func (s *Service) GetVendorPerformance() ([]models.VendorPerformance, error) {
	return s.reportingRepo.GetVendorPerformance()
}

// GetPivotTableData returns hierarchical data for pivot table display
func (s *Service) GetPivotTableData(year *string) (*PivotTableData, error) {
	return s.reportingRepo.GetPivotTableData(year)
}

// GetDrillDownData returns detailed records for a specific time period
func (s *Service) GetDrillDownData(year string, month *string, day *string) ([]models.SalesRecord, error) {
	return s.reportingRepo.GetDrillDownData(year, month, day)
}

// GetCustomSummary returns custom aggregated data
func (s *Service) GetCustomSummary(groupBy string, year *string, store *string, vendor *string) ([]models.SalesSummary, error) {
	return s.reportingRepo.GetCustomSummary(groupBy, year, store, vendor)
}

// ===== MIGRATION OPERATIONS =====

// RunMigrations executes all pending database migrations
func (s *Service) RunMigrations() error {
	return s.db.Migrate()
}

// GetMigrationStatus returns the current migration status
func (s *Service) GetMigrationStatus() ([]MigrationStatus, error) {
	return s.db.GetMigrationStatus()
}

// ResetDatabase drops all tables and re-runs migrations (USE WITH CAUTION)
func (s *Service) ResetDatabase() error {
	if err := s.db.ResetDatabase(); err != nil {
		return err
	}
	return s.db.Migrate()
}

// ===== UTILITY OPERATIONS =====

// GetVersion returns the SQLite version
func (s *Service) GetVersion() (string, error) {
	return s.db.GetVersion()
}

// GetTableInfo returns information about database tables
func (s *Service) GetTableInfo() ([]string, error) {
	return s.db.GetTableInfo()
}

// ExecTx executes a function within a transaction
func (s *Service) ExecTx(fn func(*Service) error) error {
	return s.db.ExecTx(func(tx *sql.Tx) error {
		// Create a temporary service with the transaction
		// Note: This is a simplified approach. In a production system,
		// you might want to create transaction-aware repositories
		return fn(s)
	})
}

// ===== CONVENIENCE METHODS =====

// ImportSalesData is a convenience method for importing sales data
// It validates the data and creates records in batches for better performance
func (s *Service) ImportSalesData(records []models.CreateSalesRecordRequest) (*ImportResult, error) {
	if len(records) == 0 {
		return &ImportResult{
			TotalRecords:    0,
			SuccessfulRecords: 0,
			FailedRecords:   0,
			Errors:          []string{},
		}, nil
	}

	// Validate records first
	var validRecords []models.CreateSalesRecordRequest
	var errors []string

	for i, record := range records {
		if err := validateSalesRecord(record); err != nil {
			errors = append(errors, fmt.Sprintf("Record %d: %v", i+1, err))
			continue
		}
		validRecords = append(validRecords, record)
	}

	// Import valid records
	var createdRecords []models.SalesRecord
	if len(validRecords) > 0 {
		var err error
		createdRecords, err = s.salesRepo.CreateBatch(validRecords)
		if err != nil {
			return nil, fmt.Errorf("failed to import sales data: %w", err)
		}
	}

	return &ImportResult{
		TotalRecords:      len(records),
		SuccessfulRecords: len(createdRecords),
		FailedRecords:     len(records) - len(createdRecords),
		Errors:            errors,
		CreatedRecords:    createdRecords,
	}, nil
}

// ImportResult represents the result of a data import operation
type ImportResult struct {
	TotalRecords      int                   `json:"total_records"`
	SuccessfulRecords int                   `json:"successful_records"`
	FailedRecords     int                   `json:"failed_records"`
	Errors            []string              `json:"errors,omitempty"`
	CreatedRecords    []models.SalesRecord `json:"created_records,omitempty"`
}

// validateSalesRecord performs basic validation on a sales record
func validateSalesRecord(record models.CreateSalesRecordRequest) error {
	if record.Store == "" {
		return fmt.Errorf("store is required")
	}
	if record.Vendor == "" {
		return fmt.Errorf("vendor is required")
	}
	if record.Date == "" {
		return fmt.Errorf("date is required")
	}
	if record.Description == "" {
		return fmt.Errorf("description is required")
	}
	if record.SalePrice < 0 {
		return fmt.Errorf("sale price cannot be negative")
	}
	if record.Commission < 0 {
		return fmt.Errorf("commission cannot be negative")
	}
	if record.Remaining < 0 {
		return fmt.Errorf("remaining cannot be negative")
	}
	return nil
}

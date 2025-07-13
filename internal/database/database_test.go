package database

import (
	"database/sql"
	"testing"

	"sales-track/internal/models"
)

// TestDatabaseConnection tests basic database connection and configuration
func TestDatabaseConnection(t *testing.T) {
	// Test in-memory database
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test connection health
	if !db.IsHealthy() {
		t.Error("Database should be healthy")
	}

	// Test ping
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}

	// Test version
	version, err := db.GetVersion()
	if err != nil {
		t.Errorf("Failed to get version: %v", err)
	}
	if version == "" {
		t.Error("Version should not be empty")
	}

	// Test table info
	tables, err := db.GetTableInfo()
	if err != nil {
		t.Errorf("Failed to get table info: %v", err)
	}

	// Should have sales_records table after migration
	found := false
	for _, table := range tables {
		if table == "sales_records" {
			found = true
			break
		}
	}
	if !found {
		t.Error("sales_records table should exist after migration")
	}
}

// TestSalesRepository tests CRUD operations for sales records
func TestSalesRepository(t *testing.T) {
	// Setup test database
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewSalesRepository(db)

	// Test Create
	createReq := models.CreateSalesRecordRequest{
		Store:       "Test Store",
		Vendor:      "Test Vendor",
		Date:        "2024-01-15",
		Description: "Test Product",
		SalePrice:   100.50,
		Commission:  10.05,
		Remaining:   90.45,
	}

	created, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create sales record: %v", err)
	}

	if created.ID == 0 {
		t.Error("Created record should have an ID")
	}
	if created.Store != createReq.Store {
		t.Errorf("Expected store %s, got %s", createReq.Store, created.Store)
	}
	if created.SalePrice != createReq.SalePrice {
		t.Errorf("Expected sale price %.2f, got %.2f", createReq.SalePrice, created.SalePrice)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get sales record: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}

	// Test Update
	newStore := "Updated Store"
	updateReq := models.UpdateSalesRecordRequest{
		Store: &newStore,
	}

	updated, err := repo.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update sales record: %v", err)
	}

	if updated.Store != newStore {
		t.Errorf("Expected updated store %s, got %s", newStore, updated.Store)
	}

	// Test List
	filter := models.SalesRecordFilter{}
	list, err := repo.List(filter)
	if err != nil {
		t.Fatalf("Failed to list sales records: %v", err)
	}

	if list.Total != 1 {
		t.Errorf("Expected 1 record, got %d", list.Total)
	}
	if len(list.Records) != 1 {
		t.Errorf("Expected 1 record in list, got %d", len(list.Records))
	}

	// Test Delete
	err = repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete sales record: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted record")
	}
}

// TestSalesRepositoryBatch tests batch operations
func TestSalesRepositoryBatch(t *testing.T) {
	// Setup test database
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewSalesRepository(db)

	// Test batch create
	records := []models.CreateSalesRecordRequest{
		{
			Store:       "Store A",
			Vendor:      "Vendor 1",
			Date:        "2024-01-15",
			Description: "Product A",
			SalePrice:   100.00,
			Commission:  10.00,
			Remaining:   90.00,
		},
		{
			Store:       "Store B",
			Vendor:      "Vendor 2",
			Date:        "2024-01-16",
			Description: "Product B",
			SalePrice:   200.00,
			Commission:  20.00,
			Remaining:   180.00,
		},
	}

	created, err := repo.CreateBatch(records)
	if err != nil {
		t.Fatalf("Failed to create batch records: %v", err)
	}

	if len(created) != 2 {
		t.Errorf("Expected 2 created records, got %d", len(created))
	}

	// Test stats
	stats, err := repo.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalRecords != 2 {
		t.Errorf("Expected 2 total records, got %d", stats.TotalRecords)
	}
	if stats.TotalSales != 300.00 {
		t.Errorf("Expected total sales 300.00, got %.2f", stats.TotalSales)
	}
}

// TestReportingRepository tests reporting and analytics operations
func TestReportingRepository(t *testing.T) {
	// Setup test database with sample data
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	salesRepo := NewSalesRepository(db)
	reportingRepo := NewReportingRepository(db)

	// Insert test data
	testRecords := []models.CreateSalesRecordRequest{
		{
			Store:       "Store A",
			Vendor:      "Vendor 1",
			Date:        "2024-01-15",
			Description: "Product A",
			SalePrice:   100.00,
			Commission:  10.00,
			Remaining:   90.00,
		},
		{
			Store:       "Store A",
			Vendor:      "Vendor 2",
			Date:        "2024-02-15",
			Description: "Product B",
			SalePrice:   200.00,
			Commission:  20.00,
			Remaining:   180.00,
		},
		{
			Store:       "Store B",
			Vendor:      "Vendor 1",
			Date:        "2024-01-20",
			Description: "Product C",
			SalePrice:   150.00,
			Commission:  15.00,
			Remaining:   135.00,
		},
	}

	_, err = salesRepo.CreateBatch(testRecords)
	if err != nil {
		t.Fatalf("Failed to create test records: %v", err)
	}

	// Test yearly summary
	yearly, err := reportingRepo.GetYearlySummary()
	if err != nil {
		t.Fatalf("Failed to get yearly summary: %v", err)
	}

	if len(yearly) != 1 {
		t.Errorf("Expected 1 year, got %d", len(yearly))
	}
	if yearly[0].Year != "2024" {
		t.Errorf("Expected year 2024, got %s", yearly[0].Year)
	}
	if yearly[0].ItemsSold != 3 {
		t.Errorf("Expected 3 items sold, got %d", yearly[0].ItemsSold)
	}

	// Test monthly summary
	monthly, err := reportingRepo.GetMonthlySummary(nil)
	if err != nil {
		t.Fatalf("Failed to get monthly summary: %v", err)
	}

	if len(monthly) != 2 {
		t.Errorf("Expected 2 months, got %d", len(monthly))
	}

	// Test store performance
	storePerf, err := reportingRepo.GetStorePerformance()
	if err != nil {
		t.Fatalf("Failed to get store performance: %v", err)
	}

	if len(storePerf) != 2 {
		t.Errorf("Expected 2 stores, got %d", len(storePerf))
	}

	// Test vendor performance
	vendorPerf, err := reportingRepo.GetVendorPerformance()
	if err != nil {
		t.Fatalf("Failed to get vendor performance: %v", err)
	}

	if len(vendorPerf) != 2 {
		t.Errorf("Expected 2 vendors, got %d", len(vendorPerf))
	}

	// Test pivot table data
	pivotData, err := reportingRepo.GetPivotTableData(nil)
	if err != nil {
		t.Fatalf("Failed to get pivot table data: %v", err)
	}

	if len(pivotData.YearlyData) != 1 {
		t.Errorf("Expected 1 year in pivot data, got %d", len(pivotData.YearlyData))
	}
	if len(pivotData.MonthlyData) != 2 {
		t.Errorf("Expected 2 months in pivot data, got %d", len(pivotData.MonthlyData))
	}
}

// TestDatabaseService tests the high-level service layer
func TestDatabaseService(t *testing.T) {
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	service, err := NewService(config)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Close()

	// Test health check
	if err := service.Health(); err != nil {
		t.Errorf("Service should be healthy: %v", err)
	}

	// Test import functionality
	records := []models.CreateSalesRecordRequest{
		{
			Store:       "Test Store",
			Vendor:      "Test Vendor",
			Date:        "2024-01-15",
			Description: "Test Product",
			SalePrice:   100.00,
			Commission:  10.00,
			Remaining:   90.00,
		},
	}

	result, err := service.ImportSalesData(records)
	if err != nil {
		t.Fatalf("Failed to import sales data: %v", err)
	}

	if result.TotalRecords != 1 {
		t.Errorf("Expected 1 total record, got %d", result.TotalRecords)
	}
	if result.SuccessfulRecords != 1 {
		t.Errorf("Expected 1 successful record, got %d", result.SuccessfulRecords)
	}
	if result.FailedRecords != 0 {
		t.Errorf("Expected 0 failed records, got %d", result.FailedRecords)
	}

	// Test database stats
	stats, err := service.GetDatabaseStats()
	if err != nil {
		t.Fatalf("Failed to get database stats: %v", err)
	}

	if stats.TotalRecords != 1 {
		t.Errorf("Expected 1 total record, got %d", stats.TotalRecords)
	}
}

// TestMigrations tests the migration system
func TestMigrations(t *testing.T) {
	config := Config{
		InMemory:    true,
		AutoMigrate: false, // Don't auto-migrate
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test manual migration
	err = db.Migrate()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Test migration status
	status, err := db.GetMigrationStatus()
	if err != nil {
		t.Fatalf("Failed to get migration status: %v", err)
	}

	if len(status) == 0 {
		t.Error("Expected at least one migration")
	}

	// Check that the first migration was applied
	if !status[0].Applied {
		t.Error("First migration should be applied")
	}
}

// TestTransactions tests transaction handling
func TestTransactions(t *testing.T) {
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewSalesRepository(db)

	// Test successful transaction
	err = db.ExecTx(func(tx *sql.Tx) error {
		// Create record using the repository (which uses the main connection)
		// In a real transaction scenario, you'd pass the tx to repository methods
		// For this test, we'll just verify the transaction mechanism works
		return nil
	})

	if err != nil {
		t.Fatalf("Transaction should succeed: %v", err)
	}

	// Test creating a record normally (outside transaction)
	_, err = repo.Create(models.CreateSalesRecordRequest{
		Store:       "TX Test Store",
		Vendor:      "TX Test Vendor",
		Date:        "2024-01-15",
		Description: "TX Test Product",
		SalePrice:   100.00,
		Commission:  10.00,
		Remaining:   90.00,
	})

	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Verify record was created
	filter := models.SalesRecordFilter{}
	list, err := repo.List(filter)
	if err != nil {
		t.Fatalf("Failed to list records: %v", err)
	}

	if list.Total != 1 {
		t.Errorf("Expected 1 record after creation, got %d", list.Total)
	}
}

// BenchmarkSalesRecordCreate benchmarks sales record creation
func BenchmarkSalesRecordCreate(b *testing.B) {
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewSalesRepository(db)

	record := models.CreateSalesRecordRequest{
		Store:       "Benchmark Store",
		Vendor:      "Benchmark Vendor",
		Date:        "2024-01-15",
		Description: "Benchmark Product",
		SalePrice:   100.00,
		Commission:  10.00,
		Remaining:   90.00,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.Create(record)
		if err != nil {
			b.Fatalf("Failed to create record: %v", err)
		}
	}
}

// BenchmarkSalesRecordQuery benchmarks sales record queries
func BenchmarkSalesRecordQuery(b *testing.B) {
	config := Config{
		InMemory:    true,
		AutoMigrate: true,
	}

	db, err := New(config)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewSalesRepository(db)

	// Create test data
	for i := 0; i < 1000; i++ {
		_, err := repo.Create(models.CreateSalesRecordRequest{
			Store:       "Store A",
			Vendor:      "Vendor 1",
			Date:        "2024-01-15",
			Description: "Product",
			SalePrice:   100.00,
			Commission:  10.00,
			Remaining:   90.00,
		})
		if err != nil {
			b.Fatalf("Failed to create test record: %v", err)
		}
	}

	filter := models.SalesRecordFilter{
		Limit: intPtr(50),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.List(filter)
		if err != nil {
			b.Fatalf("Failed to query records: %v", err)
		}
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

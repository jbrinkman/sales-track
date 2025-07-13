package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"sales-track/internal/database"
)

// Test HTML data for testing
const testHTMLData = `
<table>
	<tr>
		<th>Store</th>
		<th>Vendor</th>
		<th>Date</th>
		<th>Description</th>
		<th>Sale Price</th>
		<th>Commission</th>
		<th>Remaining</th>
	</tr>
	<tr>
		<td>Test Store</td>
		<td>Test Vendor</td>
		<td>2024-01-15</td>
		<td>Test Product</td>
		<td>$100.00</td>
		<td>$10.00</td>
		<td>$90.00</td>
	</tr>
	<tr>
		<td>Another Store</td>
		<td>Another Vendor</td>
		<td>2024-01-16</td>
		<td>Another Product</td>
		<td>$200.00</td>
		<td>$20.00</td>
		<td>$180.00</td>
	</tr>
</table>
`

const testConsignableData = `
<tr>
	<td>Consignable Store</td>
	<td>Consignable Vendor</td>
	<td>2024-01-17</td>
	<td>Consignable Product</td>
	<td>$150.00</td>
	<td>$15.00</td>
	<td>$135.00</td>
</tr>
`

func setupTestApp(t *testing.T) *App {
	// Create temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	app := NewApp()
	
	// Initialize with test database and auto-migration
	config := database.Config{
		FilePath:    dbPath,
		InMemory:    false,
		AutoMigrate: true, // Enable auto-migration for tests
	}
	
	dbService, err := database.NewService(config)
	if err != nil {
		t.Fatalf("Failed to create test database service: %v", err)
	}
	
	app.dbService = dbService
	app.ctx = context.Background()
	
	return app
}

func TestApp_ImportHTMLData(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	result, err := app.ImportHTMLData(testHTMLData)
	if err != nil {
		t.Fatalf("ImportHTMLData failed: %v", err)
	}

	// Debug output
	t.Logf("Import result: Success=%v, TotalRows=%d, ParsedRows=%d, ImportedRows=%d", 
		result.Success, result.TotalRows, result.ParsedRows, result.ImportedRows)
	t.Logf("Error message: %s", result.ErrorMessage)
	t.Logf("Parse errors: %d", len(result.ParseErrors))
	t.Logf("Import errors: %d", len(result.ImportErrors))
	
	if len(result.ParseErrors) > 0 {
		for i, parseErr := range result.ParseErrors {
			t.Logf("Parse error %d: Row %d, Column %s, Message: %s", i, parseErr.Row, parseErr.Column, parseErr.Message)
		}
	}
	
	if len(result.ImportErrors) > 0 {
		for i, importErr := range result.ImportErrors {
			t.Logf("Import error %d: %s", i, importErr.Error)
		}
	}

	// Verify results
	if !result.Success {
		t.Errorf("Expected success=true, got %v", result.Success)
	}

	if result.TotalRows != 2 {
		t.Errorf("Expected TotalRows=2, got %d", result.TotalRows)
	}

	if result.ParsedRows != 2 {
		t.Errorf("Expected ParsedRows=2, got %d", result.ParsedRows)
	}

	if result.ImportedRows != 2 {
		t.Errorf("Expected ImportedRows=2, got %d", result.ImportedRows)
	}

	if len(result.ImportedRecords) == 0 {
		t.Errorf("Expected imported records, got none")
		return
	}

	if len(result.ImportedRecords) != 2 {
		t.Errorf("Expected 2 imported records, got %d", len(result.ImportedRecords))
	}

	// Verify first record
	record := result.ImportedRecords[0]
	if record.Store != "Test Store" {
		t.Errorf("Expected Store='Test Store', got '%s'", record.Store)
	}

	if record.SalePrice != 100.00 {
		t.Errorf("Expected SalePrice=100.00, got %f", record.SalePrice)
	}
}

func TestApp_ImportHTMLDataBatch(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	result, err := app.ImportHTMLDataBatch(testHTMLData)
	if err != nil {
		t.Fatalf("ImportHTMLDataBatch failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Errorf("Expected success=true, got %v", result.Success)
	}

	if result.ImportedRows != 2 {
		t.Errorf("Expected ImportedRows=2, got %d", result.ImportedRows)
	}
}

func TestApp_ImportHTMLDataWithOptions_Consignable(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	options := ImportOptions{
		UseConsignableFormat: true,
		UseBatchImport:       true,
	}

	result, err := app.ImportHTMLDataWithOptions(testConsignableData, options)
	if err != nil {
		t.Fatalf("ImportHTMLDataWithOptions failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Errorf("Expected success=true, got %v", result.Success)
	}

	if result.ImportedRows != 1 {
		t.Errorf("Expected ImportedRows=1, got %d", result.ImportedRows)
	}

	// Verify the record was imported correctly
	record := result.ImportedRecords[0]
	if record.Store != "Consignable Store" {
		t.Errorf("Expected Store='Consignable Store', got '%s'", record.Store)
	}
}

func TestApp_ValidateHTMLData(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	result, err := app.ValidateHTMLData(testHTMLData)
	if err != nil {
		t.Fatalf("ValidateHTMLData failed: %v", err)
	}

	// Verify validation results
	if !result.Valid {
		t.Errorf("Expected Valid=true, got %v", result.Valid)
	}

	if result.TotalRows != 2 {
		t.Errorf("Expected TotalRows=2, got %d", result.TotalRows)
	}

	if result.ValidRows != 2 {
		t.Errorf("Expected ValidRows=2, got %d", result.ValidRows)
	}

	if result.InvalidRows != 0 {
		t.Errorf("Expected InvalidRows=0, got %d", result.InvalidRows)
	}
}

func TestApp_ValidateHTMLData_Invalid(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	invalidHTML := `
	<table>
		<tr>
			<th>Store</th>
			<th>Vendor</th>
			<th>Date</th>
			<th>Description</th>
			<th>Sale Price</th>
		</tr>
		<tr>
			<td></td>
			<td>Test Vendor</td>
			<td>invalid-date</td>
			<td></td>
			<td>not-a-price</td>
		</tr>
	</table>
	`

	result, err := app.ValidateHTMLData(invalidHTML)
	if err != nil {
		t.Fatalf("ValidateHTMLData failed: %v", err)
	}

	// Should have validation errors
	if result.Valid {
		t.Errorf("Expected Valid=false for invalid data, got %v", result.Valid)
	}

	if result.InvalidRows == 0 {
		t.Errorf("Expected InvalidRows>0 for invalid data, got %d", result.InvalidRows)
	}

	if len(result.Errors) == 0 {
		t.Errorf("Expected validation errors for invalid data, got none")
	}
}

func TestApp_GetDatabaseHealth(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	health, err := app.GetDatabaseHealth()
	if err != nil {
		t.Fatalf("GetDatabaseHealth failed: %v", err)
	}

	if !health.Connected {
		t.Errorf("Expected Connected=true, got %v", health.Connected)
	}

	if health.Error != "" {
		t.Errorf("Expected no error, got '%s'", health.Error)
	}
}

func TestApp_GetDatabaseHealth_NotInitialized(t *testing.T) {
	app := NewApp()

	health, err := app.GetDatabaseHealth()
	if err != nil {
		t.Fatalf("GetDatabaseHealth failed: %v", err)
	}

	if health.Connected {
		t.Errorf("Expected Connected=false for uninitialized service, got %v", health.Connected)
	}

	if health.Error == "" {
		t.Errorf("Expected error message for uninitialized service, got empty string")
	}
}

func TestApp_GetImportStatistics(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	// Import some test data first
	_, err := app.ImportHTMLData(testHTMLData)
	if err != nil {
		t.Fatalf("Failed to import test data: %v", err)
	}

	// Get statistics
	stats, err := app.GetImportStatistics()
	if err != nil {
		t.Fatalf("GetImportStatistics failed: %v", err)
	}

	if stats.TotalRecords != 2 {
		t.Errorf("Expected TotalRecords=2, got %d", stats.TotalRecords)
	}

	if stats.TotalSales != 300.00 { // 100 + 200
		t.Errorf("Expected TotalSales=300.00, got %f", stats.TotalSales)
	}

	if stats.AveragePrice != 150.00 { // (100 + 200) / 2
		t.Errorf("Expected AveragePrice=150.00, got %f", stats.AveragePrice)
	}
}

func TestApp_GetRecentImports(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	// Import some test data first
	_, err := app.ImportHTMLData(testHTMLData)
	if err != nil {
		t.Fatalf("Failed to import test data: %v", err)
	}

	// Get recent imports
	records, err := app.GetRecentImports(10)
	if err != nil {
		t.Fatalf("GetRecentImports failed: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 recent records, got %d", len(records))
	}

	// Records should be sorted by created_at desc, so newest first
	// Both records should have recent timestamps
	for _, record := range records {
		if time.Since(record.CreatedAt) > time.Minute {
			t.Errorf("Expected recent record, but CreatedAt is %v", record.CreatedAt)
		}
	}
}

func TestApp_ImportHTMLData_DatabaseError(t *testing.T) {
	app := NewApp() // No database service initialized

	result, err := app.ImportHTMLData(testHTMLData)
	if err == nil {
		t.Fatalf("Expected error for uninitialized database service, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result for error case, got %+v", result)
	}
}

func TestApp_ImportHTMLData_ParseError(t *testing.T) {
	app := setupTestApp(t)
	defer app.dbService.Close()

	invalidHTML := "not valid html"

	result, err := app.ImportHTMLData(invalidHTML)
	if err != nil {
		t.Fatalf("ImportHTMLData should not return error for parse failures: %v", err)
	}

	if result.Success {
		t.Errorf("Expected Success=false for invalid HTML, got %v", result.Success)
	}

	if result.ErrorMessage == "" {
		t.Errorf("Expected error message for invalid HTML, got empty string")
	}
}

// Benchmark tests
func BenchmarkApp_ImportHTMLData(b *testing.B) {
	// Create temporary database
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "bench.db")
	
	app := NewApp()
	config := database.Config{
		FilePath:    dbPath,
		InMemory:    false,
		AutoMigrate: true,
	}
	
	dbService, err := database.NewService(config)
	if err != nil {
		b.Fatalf("Failed to create benchmark database service: %v", err)
	}
	defer dbService.Close()
	
	app.dbService = dbService
	app.ctx = context.Background()

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := app.ImportHTMLData(testHTMLData)
		if err != nil {
			b.Fatalf("ImportHTMLData failed: %v", err)
		}
	}
}

func BenchmarkApp_ImportHTMLDataBatch(b *testing.B) {
	// Create temporary database
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "bench_batch.db")
	
	app := NewApp()
	config := database.Config{
		FilePath:    dbPath,
		InMemory:    false,
		AutoMigrate: true,
	}
	
	dbService, err := database.NewService(config)
	if err != nil {
		b.Fatalf("Failed to create benchmark database service: %v", err)
	}
	defer dbService.Close()
	
	app.dbService = dbService
	app.ctx = context.Background()

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := app.ImportHTMLDataBatch(testHTMLData)
		if err != nil {
			b.Fatalf("ImportHTMLDataBatch failed: %v", err)
		}
	}
}

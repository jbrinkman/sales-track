# Database Layer

This package provides a comprehensive database layer for the Sales Track application, built on SQLite with Go.

## Overview

The database layer is organized into several components:

- **Models**: Go structs representing database entities and API requests/responses
- **Connection Management**: Database connection, configuration, and health monitoring
- **Migration System**: Automated database schema migrations with version control
- **Repositories**: Data access layer with CRUD operations and queries
- **Service Layer**: High-level API combining multiple repositories
- **Testing**: Comprehensive unit tests and benchmarks

## Architecture

```
┌─────────────────┐
│   Service       │  ← High-level API, combines repositories
├─────────────────┤
│  Repositories   │  ← Data access layer (Sales, Reporting)
├─────────────────┤
│   Connection    │  ← Database connection and configuration
├─────────────────┤
│   Migrations    │  ← Schema version control and updates
├─────────────────┤
│     Models      │  ← Go structs and data types
└─────────────────┘
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    "sales-track/internal/database"
    "sales-track/internal/models"
)

func main() {
    // Create database service
    config := database.Config{
        FilePath:    "sales_track.db",
        AutoMigrate: true,
    }
    
    service, err := database.NewService(config)
    if err != nil {
        log.Fatal(err)
    }
    defer service.Close()
    
    // Create a sales record
    record := models.CreateSalesRecordRequest{
        Store:       "Downtown Store",
        Vendor:      "Electronics Plus",
        Date:        "2024-01-15",
        Description: "Samsung TV",
        SalePrice:   899.99,
        Commission:  89.99,
        Remaining:   810.00,
    }
    
    created, err := service.CreateSalesRecord(record)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Created record with ID: %d", created.ID)
    
    // Get yearly summary
    yearly, err := service.GetYearlySummary()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, year := range yearly {
        log.Printf("Year %s: %d items, $%.2f total sales", 
            year.Year, year.ItemsSold, year.TotalSales)
    }
}
```

### Configuration Options

```go
config := database.Config{
    FilePath:    "path/to/database.db",  // Database file path
    InMemory:    false,                  // Use in-memory DB (for testing)
    AutoMigrate: true,                   // Run migrations on startup
}
```

## Components

### 1. Models (`models/sales_record.go`)

Defines all data structures used throughout the application:

- **`SalesRecord`**: Main entity representing a sales transaction
- **`CreateSalesRecordRequest`**: Data for creating new records
- **`UpdateSalesRecordRequest`**: Data for updating existing records
- **`SalesRecordFilter`**: Filtering and pagination options
- **Summary Types**: `YearlySummary`, `MonthlySummary`, `DailySummary`
- **Analytics Types**: `StorePerformance`, `VendorPerformance`

### 2. Connection Management (`connection.go`)

Handles database connections and configuration:

```go
// Create connection
db, err := database.New(config)

// Health check
if db.IsHealthy() {
    log.Println("Database is healthy")
}

// Get SQLite version
version, err := db.GetVersion()

// Execute in transaction
err = db.ExecTx(func(tx *sql.Tx) error {
    // Your transactional code here
    return nil
})
```

### 3. Migration System (`migrations.go`)

Automated schema migrations with version control:

```go
// Run all pending migrations
err := db.Migrate()

// Get migration status
status, err := db.GetMigrationStatus()
for _, migration := range status {
    log.Printf("Migration %d (%s): Applied=%v", 
        migration.Version, migration.Name, migration.Applied)
}

// Reset database (USE WITH CAUTION)
err := db.ResetDatabase()
```

### 4. Sales Repository (`sales_repository.go`)

CRUD operations for sales records:

```go
repo := database.NewSalesRepository(db)

// Create single record
record, err := repo.Create(createRequest)

// Create multiple records in batch
records, err := repo.CreateBatch(createRequests)

// Get by ID
record, err := repo.GetByID(123)

// Update record
updated, err := repo.Update(123, updateRequest)

// Delete record
err := repo.Delete(123)

// List with filtering and pagination
filter := models.SalesRecordFilter{
    Store:     stringPtr("Downtown Store"),
    DateFrom:  timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
    Limit:     intPtr(50),
    SortBy:    stringPtr("date"),
    SortOrder: stringPtr("desc"),
}
list, err := repo.List(filter)

// Get database statistics
stats, err := repo.GetStats()
```

### 5. Reporting Repository (`reporting_repository.go`)

Analytics and reporting queries:

```go
repo := database.NewReportingRepository(db)

// Yearly summary (pivot table top level)
yearly, err := repo.GetYearlySummary()

// Monthly summary with optional year filter
monthly, err := repo.GetMonthlySummary(stringPtr("2024"))

// Daily summary with year/month filters
daily, err := repo.GetDailySummary(stringPtr("2024"), stringPtr("01"))

// Store performance analytics
stores, err := repo.GetStorePerformance()

// Vendor performance analytics
vendors, err := repo.GetVendorPerformance()

// Pivot table data (hierarchical)
pivotData, err := repo.GetPivotTableData(stringPtr("2024"))

// Drill-down to specific time period
records, err := repo.GetDrillDownData("2024", stringPtr("01"), stringPtr("15"))

// Custom aggregations
summary, err := repo.GetCustomSummary("month", stringPtr("2024"), nil, nil)
```

### 6. Service Layer (`service.go`)

High-level API combining all repositories:

```go
service, err := database.NewService(config)

// All sales repository methods available
record, err := service.CreateSalesRecord(createRequest)
list, err := service.ListSalesRecords(filter)

// All reporting repository methods available
yearly, err := service.GetYearlySummary()
pivotData, err := service.GetPivotTableData(nil)

// Convenience methods
result, err := service.ImportSalesData(records)
stats, err := service.GetDatabaseStats()

// Migration operations
err = service.RunMigrations()
status, err := service.GetMigrationStatus()

// Health and utility
err = service.Health()
version, err := service.GetVersion()
```

## Data Import

The service layer provides a convenient `ImportSalesData` method for bulk imports:

```go
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
    // ... more records
}

result, err := service.ImportSalesData(records)
if err != nil {
    log.Fatal(err)
}

log.Printf("Imported %d/%d records successfully", 
    result.SuccessfulRecords, result.TotalRecords)

if len(result.Errors) > 0 {
    log.Printf("Errors: %v", result.Errors)
}
```

## Reporting and Analytics

### Pivot Table Data

The core Excel replacement functionality:

```go
// Get hierarchical data for pivot table
pivotData, err := service.GetPivotTableData(nil) // All years
// or
pivotData, err := service.GetPivotTableData(stringPtr("2024")) // Specific year

// Access different levels
for _, year := range pivotData.YearlyData {
    log.Printf("Year %s: %d items, $%.2f remaining", 
        year.Year, year.ItemsSold, year.TotalRemaining)
}

for _, month := range pivotData.MonthlyData {
    log.Printf("Month %s: %d items, $%.2f remaining", 
        month.YearMonth, month.ItemsSold, month.TotalRemaining)
}
```

### Drill-Down Functionality

```go
// Drill down to specific time periods
yearRecords, err := service.GetDrillDownData("2024", nil, nil)
monthRecords, err := service.GetDrillDownData("2024", stringPtr("01"), nil)
dayRecords, err := service.GetDrillDownData("2024", stringPtr("01"), stringPtr("15"))
```

### Performance Analytics

```go
// Store performance
stores, err := service.GetStorePerformance()
for _, store := range stores {
    log.Printf("Store %s: $%.2f total sales, %.2f avg sale", 
        store.Store, store.TotalSales, store.AvgSalePrice)
}

// Vendor performance
vendors, err := service.GetVendorPerformance()
for _, vendor := range vendors {
    log.Printf("Vendor %s: %d items, %d stores", 
        vendor.Vendor, vendor.TotalItems, vendor.UniqueStores)
}
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./internal/database

# Run tests with coverage
go test -cover ./internal/database

# Run benchmarks
go test -bench=. ./internal/database

# Run specific test
go test -run TestSalesRepository ./internal/database
```

### Test Coverage

The test suite covers:

- ✅ Database connection and configuration
- ✅ Migration system functionality
- ✅ CRUD operations for sales records
- ✅ Batch operations and transactions
- ✅ Filtering and pagination
- ✅ Reporting and analytics queries
- ✅ Service layer integration
- ✅ Error handling and edge cases
- ✅ Performance benchmarks

## Performance Considerations

### Indexing Strategy

The database uses strategic indexing for optimal performance:

- **Primary Index**: Date (DESC) - Most common query pattern
- **Secondary Indexes**: Store, Vendor - Filtering operations
- **Composite Indexes**: Store+Date, Vendor+Date - Complex queries
- **Specialized Indexes**: Remaining balance, Created timestamp

### Query Optimization

- **Views**: Pre-built aggregations for common reporting patterns
- **Batch Operations**: Efficient bulk inserts with transactions
- **Connection Pooling**: Optimized SQLite configuration
- **Prepared Statements**: Reusable queries for better performance

### Memory Usage

- **64MB Cache**: Configured for optimal performance
- **WAL Mode**: Better concurrency and crash recovery
- **In-Memory Temp**: Faster temporary operations

## Error Handling

The database layer provides comprehensive error handling:

```go
// Specific error types
record, err := service.GetSalesRecord(999)
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        // Handle not found case
    } else {
        // Handle other database errors
    }
}

// Validation errors during import
result, err := service.ImportSalesData(records)
if err != nil {
    log.Fatal(err) // Critical error
}

// Individual record validation errors
for _, errMsg := range result.Errors {
    log.Printf("Validation error: %s", errMsg)
}
```

## Integration with Wails

The database layer is designed to integrate seamlessly with Wails:

```go
// In your Wails app context
type App struct {
    ctx context.Context
    db  *database.Service
}

func NewApp() *App {
    config := database.Config{
        FilePath:    "sales_track.db",
        AutoMigrate: true,
    }
    
    db, err := database.NewService(config)
    if err != nil {
        log.Fatal(err)
    }
    
    return &App{db: db}
}

// Wails context methods
func (a *App) CreateSalesRecord(record models.CreateSalesRecordRequest) (*models.SalesRecord, error) {
    return a.db.CreateSalesRecord(record)
}

func (a *App) GetYearlySummary() ([]models.YearlySummary, error) {
    return a.db.GetYearlySummary()
}
```

## Future Enhancements

Planned improvements:

- **Connection Pooling**: Multiple connections for high concurrency
- **Query Builder**: Fluent API for complex queries
- **Caching Layer**: Redis integration for frequently accessed data
- **Audit Logging**: Track all data changes
- **Data Validation**: Enhanced validation with custom rules
- **Backup/Restore**: Automated backup and restore functionality

## Troubleshooting

### Common Issues

1. **Migration Failures**: Check file permissions and disk space
2. **Connection Issues**: Verify database file path and permissions
3. **Performance Issues**: Check index usage with `EXPLAIN QUERY PLAN`
4. **Lock Issues**: Ensure proper transaction handling

### Debug Mode

Enable detailed logging:

```go
// Add debug logging to connection
db.conn.SetMaxOpenConns(1) // Force single connection for debugging
```

### Health Checks

```go
// Regular health monitoring
if err := service.Health(); err != nil {
    log.Printf("Database health check failed: %v", err)
    // Implement reconnection logic
}
```

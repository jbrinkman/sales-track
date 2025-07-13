package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"sales-track/internal/database"
	"sales-track/internal/models"
	"sales-track/internal/parser"
)

// App struct
type App struct {
	ctx       context.Context
	dbService *database.Service
	parser    *parser.HTMLTableParser
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		parser: parser.NewHTMLTableParser(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Initialize database service
	dbPath := filepath.Join(".", "sales_track.db")
	config := database.Config{
		FilePath:    dbPath,
		InMemory:    false,
		AutoMigrate: true, // Enable auto-migration
	}
	
	dbService, err := database.NewService(config)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return
	}
	
	a.dbService = dbService
	log.Println("Database service initialized successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ImportHTMLData imports HTML table data into the database
func (a *App) ImportHTMLData(htmlData string) (*ImportResult, error) {
	if a.dbService == nil {
		return nil, fmt.Errorf("database service not initialized")
	}

	// Parse HTML data
	parseResult, err := a.parser.ParseHTML(htmlData)
	if err != nil {
		return &ImportResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to parse HTML data: %v", err),
		}, nil
	}

	// Convert parsed records to database format and import
	var importedRecords []models.SalesRecord
	var importErrors []ImportError

	for _, record := range parseResult.Records {
		// Import individual record
		savedRecord, err := a.dbService.CreateSalesRecord(record)
		if err != nil {
			importErrors = append(importErrors, ImportError{
				Record: record,
				Error:  err.Error(),
			})
			continue
		}
		importedRecords = append(importedRecords, *savedRecord)
	}

	// Prepare result
	result := &ImportResult{
		Success:           len(importedRecords) > 0,
		TotalRows:         parseResult.TotalRows,
		ParsedRows:        parseResult.SuccessCount,
		ImportedRows:      len(importedRecords),
		ParseErrors:       parseResult.Errors,
		ImportErrors:      importErrors,
		ProcessingTime:    parseResult.Statistics.ProcessingTime,
		ImportedRecords:   importedRecords,
		ColumnMapping:     parseResult.ColumnMapping,
		DataTypesDetected: parseResult.Statistics.DataTypesDetected,
	}

	if len(importErrors) > 0 {
		result.ErrorMessage = fmt.Sprintf("Imported %d of %d records. %d records failed to import.", 
			len(importedRecords), parseResult.SuccessCount, len(importErrors))
	}

	return result, nil
}

// ImportHTMLDataBatch imports HTML data using batch operations for better performance
func (a *App) ImportHTMLDataBatch(htmlData string) (*ImportResult, error) {
	if a.dbService == nil {
		return nil, fmt.Errorf("database service not initialized")
	}

	// Parse HTML data
	parseResult, err := a.parser.ParseHTML(htmlData)
	if err != nil {
		return &ImportResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to parse HTML data: %v", err),
		}, nil
	}

	// Use batch import for better performance
	importedRecords, err := a.dbService.CreateSalesRecordsBatch(parseResult.Records)
	if err != nil {
		return &ImportResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to import records: %v", err),
			TotalRows:    parseResult.TotalRows,
			ParsedRows:   parseResult.SuccessCount,
			ParseErrors:  parseResult.Errors,
		}, nil
	}

	// Prepare result
	result := &ImportResult{
		Success:           true,
		TotalRows:         parseResult.TotalRows,
		ParsedRows:        parseResult.SuccessCount,
		ImportedRows:      len(importedRecords),
		ProcessingTime:    parseResult.Statistics.ProcessingTime,
		ImportedRecords:   importedRecords,
		ColumnMapping:     parseResult.ColumnMapping,
		DataTypesDetected: parseResult.Statistics.DataTypesDetected,
	}

	return result, nil
}

// ImportHTMLDataWithOptions imports HTML data with parsing options
func (a *App) ImportHTMLDataWithOptions(htmlData string, options ImportOptions) (*ImportResult, error) {
	if a.dbService == nil {
		return nil, fmt.Errorf("database service not initialized")
	}

	// Configure parser based on options
	if options.UseConsignableFormat {
		a.parser.SetConsignableMapping()
	} else if len(options.CustomColumnMapping) > 0 {
		a.parser.SetPositionalMapping(options.CustomColumnMapping)
	}

	// Set strict mode if requested
	a.parser.StrictMode = options.StrictMode

	// Use batch import if available
	if options.UseBatchImport {
		return a.ImportHTMLDataBatch(htmlData)
	}

	return a.ImportHTMLData(htmlData)
}

// GetImportStatistics returns statistics about imported data
func (a *App) GetImportStatistics() (*ImportStatistics, error) {
	if a.dbService == nil {
		return nil, fmt.Errorf("database service not initialized")
	}

	// Get database statistics
	stats, err := a.dbService.GetDatabaseStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get database statistics: %v", err)
	}

	// Calculate recent records (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	limit := 1 // We just want the count
	recentFilter := models.SalesRecordFilter{
		DateFrom: &thirtyDaysAgo,
		Limit:    &limit,
	}
	
	recentList, err := a.dbService.ListSalesRecords(recentFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent records: %v", err)
	}

	return &ImportStatistics{
		TotalRecords:  int(stats.TotalRecords),
		RecentRecords: int(recentList.Total),
		TotalSales:    stats.TotalSales,
		AveragePrice:  stats.AvgSalePrice,
	}, nil
}

// ValidateHTMLData validates HTML data without importing
func (a *App) ValidateHTMLData(htmlData string) (*ValidationResult, error) {
	// Parse HTML data without importing
	parseResult, err := a.parser.ParseHTML(htmlData)
	if err != nil {
		return &ValidationResult{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("Failed to parse HTML data: %v", err),
		}, nil
	}

	return &ValidationResult{
		Valid:             parseResult.SuccessCount > 0,
		TotalRows:         parseResult.TotalRows,
		ValidRows:         parseResult.SuccessCount,
		InvalidRows:       parseResult.ErrorCount,
		Errors:            parseResult.Errors,
		Warnings:          parseResult.Warnings,
		ColumnMapping:     parseResult.ColumnMapping,
		DataTypesDetected: parseResult.Statistics.DataTypesDetected,
		ProcessingTime:    parseResult.Statistics.ProcessingTime,
	}, nil
}

// GetDatabaseHealth returns database connection health status
func (a *App) GetDatabaseHealth() (*DatabaseHealth, error) {
	if a.dbService == nil {
		return &DatabaseHealth{
			Connected: false,
			Error:     "Database service not initialized",
		}, nil
	}

	// Check database health
	err := a.dbService.Health()
	if err != nil {
		return &DatabaseHealth{
			Connected: false,
			Error:     err.Error(),
		}, nil
	}

	return &DatabaseHealth{
		Connected: true,
	}, nil
}

// GetRecentImports returns recently imported sales records
func (a *App) GetRecentImports(limit int) ([]models.SalesRecord, error) {
	if a.dbService == nil {
		return nil, fmt.Errorf("database service not initialized")
	}

	sortBy := "created_at"
	sortOrder := "desc"
	filter := models.SalesRecordFilter{
		Limit:     &limit,
		SortBy:    &sortBy,
		SortOrder: &sortOrder,
	}

	result, err := a.dbService.ListSalesRecords(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent imports: %v", err)
	}

	return result.Records, nil
}

package main

import (
	"time"

	"sales-track/internal/models"
	"sales-track/internal/parser"
)

// ImportResult represents the result of an HTML data import operation
type ImportResult struct {
	Success           bool                      `json:"success"`
	TotalRows         int                       `json:"total_rows"`
	ParsedRows        int                       `json:"parsed_rows"`
	ImportedRows      int                       `json:"imported_rows"`
	ErrorMessage      string                    `json:"error_message,omitempty"`
	ParseErrors       []parser.ParseError       `json:"parse_errors,omitempty"`
	ImportErrors      []ImportError             `json:"import_errors,omitempty"`
	ProcessingTime    time.Duration             `json:"processing_time"`
	ImportedRecords   []models.SalesRecord      `json:"imported_records,omitempty"`
	ColumnMapping     map[string]int            `json:"column_mapping"`
	DataTypesDetected map[string]string         `json:"data_types_detected"`
}

// ImportError represents an error that occurred during database import
type ImportError struct {
	Record models.CreateSalesRecordRequest `json:"record"`
	Error  string                          `json:"error"`
}

// ImportOptions provides configuration options for HTML data import
type ImportOptions struct {
	UseConsignableFormat bool     `json:"use_consignable_format"`
	CustomColumnMapping  []string `json:"custom_column_mapping,omitempty"`
	StrictMode           bool     `json:"strict_mode"`
	UseBatchImport       bool     `json:"use_batch_import"`
}

// ValidationResult represents the result of HTML data validation
type ValidationResult struct {
	Valid             bool                      `json:"valid"`
	TotalRows         int                       `json:"total_rows"`
	ValidRows         int                       `json:"valid_rows"`
	InvalidRows       int                       `json:"invalid_rows"`
	ErrorMessage      string                    `json:"error_message,omitempty"`
	Errors            []parser.ParseError       `json:"errors,omitempty"`
	Warnings          []parser.ParseWarning     `json:"warnings,omitempty"`
	ColumnMapping     map[string]int            `json:"column_mapping"`
	DataTypesDetected map[string]string         `json:"data_types_detected"`
	ProcessingTime    time.Duration             `json:"processing_time"`
}

// ImportStatistics provides statistics about imported data
type ImportStatistics struct {
	TotalRecords  int     `json:"total_records"`
	RecentRecords int     `json:"recent_records"`
	TotalSales    float64 `json:"total_sales"`
	AveragePrice  float64 `json:"average_price"`
}

// DatabaseHealth represents the health status of the database connection
type DatabaseHealth struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

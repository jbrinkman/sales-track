package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// SalesRecord represents a single sales transaction record
// This struct maps directly to the sales_records table in the database
type SalesRecord struct {
	ID          int64     `json:"id" db:"id"`
	Store       string    `json:"store" db:"store"`
	Vendor      string    `json:"vendor" db:"vendor"`
	Date        time.Time `json:"date" db:"date"`
	Description string    `json:"description" db:"description"`
	SalePrice   float64   `json:"sale_price" db:"sale_price"`
	Commission  float64   `json:"commission" db:"commission"`
	Remaining   float64   `json:"remaining" db:"remaining"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NullTime handles nullable time fields from SQLite
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface for database/sql
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return nil
	case string:
		// Try to parse SQLite datetime format
		if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			nt.Time, nt.Valid = t, true
			return nil
		}
		// Try to parse SQLite date format
		if t, err := time.Parse("2006-01-02", v); err == nil {
			nt.Time, nt.Valid = t, true
			return nil
		}
		return fmt.Errorf("cannot parse time string: %s", v)
	default:
		return fmt.Errorf("cannot scan %T into NullTime", value)
	}
}

// Value implements the driver Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// CreateSalesRecordRequest represents the data needed to create a new sales record
// Used for API requests and data import operations
type CreateSalesRecordRequest struct {
	Store       string  `json:"store" validate:"required,min=1,max=100"`
	Vendor      string  `json:"vendor" validate:"required,min=1,max=100"`
	Date        string  `json:"date" validate:"required"` // Date as string for parsing
	Description string  `json:"description" validate:"required,min=1"`
	SalePrice   float64 `json:"sale_price" validate:"required,min=0"`
	Commission  float64 `json:"commission" validate:"min=0"`
	Remaining   float64 `json:"remaining" validate:"min=0"`
}

// UpdateSalesRecordRequest represents the data that can be updated for a sales record
type UpdateSalesRecordRequest struct {
	Store       *string  `json:"store,omitempty" validate:"omitempty,min=1,max=100"`
	Vendor      *string  `json:"vendor,omitempty" validate:"omitempty,min=1,max=100"`
	Date        *string  `json:"date,omitempty"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=1"`
	SalePrice   *float64 `json:"sale_price,omitempty" validate:"omitempty,min=0"`
	Commission  *float64 `json:"commission,omitempty" validate:"omitempty,min=0"`
	Remaining   *float64 `json:"remaining,omitempty" validate:"omitempty,min=0"`
}

// SalesRecordFilter represents filtering options for querying sales records
type SalesRecordFilter struct {
	Store     *string    `json:"store,omitempty"`
	Vendor    *string    `json:"vendor,omitempty"`
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
	MinPrice  *float64   `json:"min_price,omitempty"`
	MaxPrice  *float64   `json:"max_price,omitempty"`
	Limit     *int       `json:"limit,omitempty"`
	Offset    *int       `json:"offset,omitempty"`
	SortBy    *string    `json:"sort_by,omitempty"`    // date, store, vendor, sale_price
	SortOrder *string    `json:"sort_order,omitempty"` // asc, desc
}

// SalesRecordList represents a paginated list of sales records
type SalesRecordList struct {
	Records    []SalesRecord `json:"records"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// SalesSummary represents aggregated sales data
type SalesSummary struct {
	Period        string  `json:"period"`         // Year, Month, or Date
	ItemsSold     int64   `json:"items_sold"`     // Count of records
	TotalSales    float64 `json:"total_sales"`    // Sum of sale_price
	TotalCommission float64 `json:"total_commission"` // Sum of commission
	TotalRemaining  float64 `json:"total_remaining"`  // Sum of remaining
	UniqueStores    int64   `json:"unique_stores"`    // Count of distinct stores
	UniqueVendors   int64   `json:"unique_vendors"`   // Count of distinct vendors
}

// YearlySummary represents yearly aggregated data
type YearlySummary struct {
	Year            string  `json:"year"`
	ItemsSold       int64   `json:"items_sold"`
	TotalSales      float64 `json:"total_sales"`
	TotalCommission float64 `json:"total_commission"`
	TotalRemaining  float64 `json:"total_remaining"`
	UniqueStores    int64   `json:"unique_stores"`
	UniqueVendors   int64   `json:"unique_vendors"`
}

// MonthlySummary represents monthly aggregated data
type MonthlySummary struct {
	Year            string  `json:"year"`
	Month           string  `json:"month"`
	YearMonth       string  `json:"year_month"`
	ItemsSold       int64   `json:"items_sold"`
	TotalSales      float64 `json:"total_sales"`
	TotalCommission float64 `json:"total_commission"`
	TotalRemaining  float64 `json:"total_remaining"`
	UniqueStores    int64   `json:"unique_stores"`
	UniqueVendors   int64   `json:"unique_vendors"`
}

// DailySummary represents daily aggregated data
type DailySummary struct {
	Date            time.Time `json:"date"`
	Year            string    `json:"year"`
	Month           string    `json:"month"`
	Day             string    `json:"day"`
	YearMonth       string    `json:"year_month"`
	ItemsSold       int64     `json:"items_sold"`
	TotalSales      float64   `json:"total_sales"`
	TotalCommission float64   `json:"total_commission"`
	TotalRemaining  float64   `json:"total_remaining"`
	UniqueStores    int64     `json:"unique_stores"`
	UniqueVendors   int64     `json:"unique_vendors"`
}

// StorePerformance represents store-based analytics
type StorePerformance struct {
	Store           string    `json:"store"`
	TotalItems      int64     `json:"total_items"`
	TotalSales      float64   `json:"total_sales"`
	TotalCommission float64   `json:"total_commission"`
	TotalRemaining  float64   `json:"total_remaining"`
	AvgSalePrice    float64   `json:"avg_sale_price"`
	FirstSaleDate   time.Time `json:"first_sale_date"`
	LastSaleDate    time.Time `json:"last_sale_date"`
	UniqueVendors   int64     `json:"unique_vendors"`
}

// VendorPerformance represents vendor-based analytics
type VendorPerformance struct {
	Vendor          string    `json:"vendor"`
	TotalItems      int64     `json:"total_items"`
	TotalSales      float64   `json:"total_sales"`
	TotalCommission float64   `json:"total_commission"`
	TotalRemaining  float64   `json:"total_remaining"`
	AvgSalePrice    float64   `json:"avg_sale_price"`
	FirstSaleDate   time.Time `json:"first_sale_date"`
	LastSaleDate    time.Time `json:"last_sale_date"`
	UniqueStores    int64     `json:"unique_stores"`
}

// DatabaseStats represents overall database statistics
type DatabaseStats struct {
	TotalRecords    int64     `json:"total_records"`
	EarliestDate    time.Time `json:"earliest_date"`
	LatestDate      time.Time `json:"latest_date"`
	TotalSales      float64   `json:"total_sales"`
	AvgSalePrice    float64   `json:"avg_sale_price"`
	UniqueStores    int64     `json:"unique_stores"`
	UniqueVendors   int64     `json:"unique_vendors"`
	LastUpdated     time.Time `json:"last_updated"`
}

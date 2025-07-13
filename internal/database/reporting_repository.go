package database

import (
	"fmt"
	"time"

	"sales-track/internal/models"
)

// ReportingRepository handles database operations for reporting and analytics
type ReportingRepository struct {
	db *DB
}

// NewReportingRepository creates a new reporting repository
func NewReportingRepository(db *DB) *ReportingRepository {
	return &ReportingRepository{db: db}
}

// GetYearlySummary returns yearly sales summary data
func (r *ReportingRepository) GetYearlySummary() ([]models.YearlySummary, error) {
	query := `
		SELECT 
			year,
			items_sold,
			total_sales,
			total_commission,
			total_remaining,
			unique_stores,
			unique_vendors
		FROM v_yearly_sales_summary
		ORDER BY year DESC
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query yearly summary: %w", err)
	}
	defer rows.Close()

	var summaries []models.YearlySummary
	for rows.Next() {
		var summary models.YearlySummary
		err := rows.Scan(
			&summary.Year,
			&summary.ItemsSold,
			&summary.TotalSales,
			&summary.TotalCommission,
			&summary.TotalRemaining,
			&summary.UniqueStores,
			&summary.UniqueVendors,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan yearly summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating yearly summaries: %w", err)
	}

	return summaries, nil
}

// GetMonthlySummary returns monthly sales summary data, optionally filtered by year
func (r *ReportingRepository) GetMonthlySummary(year *string) ([]models.MonthlySummary, error) {
	query := `
		SELECT 
			year,
			month,
			year_month,
			items_sold,
			total_sales,
			total_commission,
			total_remaining,
			unique_stores,
			unique_vendors
		FROM v_monthly_sales_summary
	`

	args := []interface{}{}
	if year != nil {
		query += " WHERE year = ?"
		args = append(args, *year)
	}

	query += " ORDER BY year DESC, month DESC"

	rows, err := r.db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly summary: %w", err)
	}
	defer rows.Close()

	var summaries []models.MonthlySummary
	for rows.Next() {
		var summary models.MonthlySummary
		err := rows.Scan(
			&summary.Year,
			&summary.Month,
			&summary.YearMonth,
			&summary.ItemsSold,
			&summary.TotalSales,
			&summary.TotalCommission,
			&summary.TotalRemaining,
			&summary.UniqueStores,
			&summary.UniqueVendors,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monthly summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating monthly summaries: %w", err)
	}

	return summaries, nil
}

// GetDailySummary returns daily sales summary data, optionally filtered by year and month
func (r *ReportingRepository) GetDailySummary(year *string, month *string) ([]models.DailySummary, error) {
	query := `
		SELECT 
			date,
			year,
			month,
			day,
			year_month,
			items_sold,
			total_sales,
			total_commission,
			total_remaining,
			unique_stores,
			unique_vendors
		FROM v_daily_sales_summary
	`

	args := []interface{}{}
	whereParts := []string{}

	if year != nil {
		whereParts = append(whereParts, "year = ?")
		args = append(args, *year)
	}
	if month != nil {
		whereParts = append(whereParts, "month = ?")
		args = append(args, *month)
	}

	if len(whereParts) > 0 {
		query += " WHERE " + whereParts[0]
		if len(whereParts) > 1 {
			query += " AND " + whereParts[1]
		}
	}

	query += " ORDER BY date DESC"

	rows, err := r.db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily summary: %w", err)
	}
	defer rows.Close()

	var summaries []models.DailySummary
	for rows.Next() {
		var summary models.DailySummary
		err := rows.Scan(
			&summary.Date,
			&summary.Year,
			&summary.Month,
			&summary.Day,
			&summary.YearMonth,
			&summary.ItemsSold,
			&summary.TotalSales,
			&summary.TotalCommission,
			&summary.TotalRemaining,
			&summary.UniqueStores,
			&summary.UniqueVendors,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating daily summaries: %w", err)
	}

	return summaries, nil
}

// GetStorePerformance returns store performance analytics
func (r *ReportingRepository) GetStorePerformance() ([]models.StorePerformance, error) {
	query := `
		SELECT 
			store,
			total_items,
			total_sales,
			total_commission,
			total_remaining,
			avg_sale_price,
			first_sale_date,
			last_sale_date,
			unique_vendors
		FROM v_store_performance
		ORDER BY total_sales DESC
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query store performance: %w", err)
	}
	defer rows.Close()

	var performances []models.StorePerformance
	for rows.Next() {
		var performance models.StorePerformance
		var firstSaleDateStr, lastSaleDateStr string
		
		err := rows.Scan(
			&performance.Store,
			&performance.TotalItems,
			&performance.TotalSales,
			&performance.TotalCommission,
			&performance.TotalRemaining,
			&performance.AvgSalePrice,
			&firstSaleDateStr,
			&lastSaleDateStr,
			&performance.UniqueVendors,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan store performance: %w", err)
		}
		
		// Parse date strings
		if parsed, err := time.Parse("2006-01-02", firstSaleDateStr); err == nil {
			performance.FirstSaleDate = parsed
		}
		if parsed, err := time.Parse("2006-01-02", lastSaleDateStr); err == nil {
			performance.LastSaleDate = parsed
		}
		
		performances = append(performances, performance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating store performances: %w", err)
	}

	return performances, nil
}

// GetVendorPerformance returns vendor performance analytics
func (r *ReportingRepository) GetVendorPerformance() ([]models.VendorPerformance, error) {
	query := `
		SELECT 
			vendor,
			total_items,
			total_sales,
			total_commission,
			total_remaining,
			avg_sale_price,
			first_sale_date,
			last_sale_date,
			unique_stores
		FROM v_vendor_performance
		ORDER BY total_sales DESC
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vendor performance: %w", err)
	}
	defer rows.Close()

	var performances []models.VendorPerformance
	for rows.Next() {
		var performance models.VendorPerformance
		var firstSaleDateStr, lastSaleDateStr string
		
		err := rows.Scan(
			&performance.Vendor,
			&performance.TotalItems,
			&performance.TotalSales,
			&performance.TotalCommission,
			&performance.TotalRemaining,
			&performance.AvgSalePrice,
			&firstSaleDateStr,
			&lastSaleDateStr,
			&performance.UniqueStores,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vendor performance: %w", err)
		}
		
		// Parse date strings
		if parsed, err := time.Parse("2006-01-02", firstSaleDateStr); err == nil {
			performance.FirstSaleDate = parsed
		}
		if parsed, err := time.Parse("2006-01-02", lastSaleDateStr); err == nil {
			performance.LastSaleDate = parsed
		}
		
		performances = append(performances, performance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating vendor performances: %w", err)
	}

	return performances, nil
}

// GetPivotTableData returns hierarchical data for pivot table display
// This is the core function for the Excel replacement workflow
func (r *ReportingRepository) GetPivotTableData(year *string) (*PivotTableData, error) {
	// Get yearly data
	yearlyData, err := r.GetYearlySummary()
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly data: %w", err)
	}

	// Filter yearly data if year is specified
	if year != nil {
		filteredYearly := []models.YearlySummary{}
		for _, yearly := range yearlyData {
			if yearly.Year == *year {
				filteredYearly = append(filteredYearly, yearly)
				break
			}
		}
		yearlyData = filteredYearly
	}

	// Get monthly data
	monthlyData, err := r.GetMonthlySummary(year)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly data: %w", err)
	}

	// Get daily data
	dailyData, err := r.GetDailySummary(year, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily data: %w", err)
	}

	return &PivotTableData{
		YearlyData:  yearlyData,
		MonthlyData: monthlyData,
		DailyData:   dailyData,
	}, nil
}

// PivotTableData represents hierarchical data for pivot table display
type PivotTableData struct {
	YearlyData  []models.YearlySummary  `json:"yearly_data"`
	MonthlyData []models.MonthlySummary `json:"monthly_data"`
	DailyData   []models.DailySummary   `json:"daily_data"`
}

// GetDrillDownData returns detailed records for a specific time period
func (r *ReportingRepository) GetDrillDownData(year string, month *string, day *string) ([]models.SalesRecord, error) {
	query := `
		SELECT id, store, vendor, date, description, sale_price, commission, remaining, created_at, updated_at
		FROM sales_records
		WHERE strftime('%Y', date) = ?
	`
	args := []interface{}{year}

	if month != nil {
		query += " AND strftime('%m', date) = ?"
		args = append(args, *month)
	}

	if day != nil {
		query += " AND strftime('%d', date) = ?"
		args = append(args, *day)
	}

	query += " ORDER BY date DESC, id DESC"

	rows, err := r.db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query drill-down data: %w", err)
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
		return nil, fmt.Errorf("error iterating drill-down records: %w", err)
	}

	return records, nil
}

// GetCustomSummary returns custom aggregated data based on grouping criteria
func (r *ReportingRepository) GetCustomSummary(groupBy string, year *string, store *string, vendor *string) ([]models.SalesSummary, error) {
	// Validate groupBy parameter
	validGroupBy := map[string]string{
		"year":   "strftime('%Y', date)",
		"month":  "strftime('%Y-%m', date)",
		"day":    "date",
		"store":  "store",
		"vendor": "vendor",
	}

	groupByClause, valid := validGroupBy[groupBy]
	if !valid {
		return nil, fmt.Errorf("invalid groupBy parameter: %s", groupBy)
	}

	query := fmt.Sprintf(`
		SELECT 
			%s as period,
			COUNT(*) as items_sold,
			SUM(sale_price) as total_sales,
			SUM(commission) as total_commission,
			SUM(remaining) as total_remaining,
			COUNT(DISTINCT store) as unique_stores,
			COUNT(DISTINCT vendor) as unique_vendors
		FROM sales_records
	`, groupByClause)

	args := []interface{}{}
	whereParts := []string{}

	if year != nil {
		whereParts = append(whereParts, "strftime('%Y', date) = ?")
		args = append(args, *year)
	}
	if store != nil {
		whereParts = append(whereParts, "store = ?")
		args = append(args, *store)
	}
	if vendor != nil {
		whereParts = append(whereParts, "vendor = ?")
		args = append(args, *vendor)
	}

	if len(whereParts) > 0 {
		query += " WHERE " + whereParts[0]
		for i := 1; i < len(whereParts); i++ {
			query += " AND " + whereParts[i]
		}
	}

	query += fmt.Sprintf(" GROUP BY %s ORDER BY period DESC", groupByClause)

	rows, err := r.db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom summary: %w", err)
	}
	defer rows.Close()

	var summaries []models.SalesSummary
	for rows.Next() {
		var summary models.SalesSummary
		err := rows.Scan(
			&summary.Period,
			&summary.ItemsSold,
			&summary.TotalSales,
			&summary.TotalCommission,
			&summary.TotalRemaining,
			&summary.UniqueStores,
			&summary.UniqueVendors,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan custom summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating custom summaries: %w", err)
	}

	return summaries, nil
}

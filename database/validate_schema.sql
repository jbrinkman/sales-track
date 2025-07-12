-- Database Schema Validation Script
-- This script validates the database schema and tests all constraints, indexes, and views
-- Run this script after creating the database to ensure everything works correctly

-- ============================================================================
-- SCHEMA VALIDATION
-- ============================================================================

-- Check if main table exists
SELECT name FROM sqlite_master WHERE type='table' AND name='sales_records';

-- Verify table structure
PRAGMA table_info(sales_records);

-- Check all indexes exist
SELECT name, sql FROM sqlite_master WHERE type='index' AND tbl_name='sales_records';

-- Check all views exist
SELECT name FROM sqlite_master WHERE type='view';

-- Check triggers exist
SELECT name, sql FROM sqlite_master WHERE type='trigger';

-- ============================================================================
-- CONSTRAINT TESTING
-- ============================================================================

-- Test positive value constraints (these should fail)
-- INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) 
-- VALUES ('Test Store', 'Test Vendor', '2024-01-01', 'Test Item', -10.00, 5.00, 0.00);

-- INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) 
-- VALUES ('Test Store', 'Test Vendor', '2024-01-01', 'Test Item', 10.00, -5.00, 0.00);

-- Test empty string constraints (these should fail)
-- INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) 
-- VALUES ('', 'Test Vendor', '2024-01-01', 'Test Item', 10.00, 1.00, 9.00);

-- INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) 
-- VALUES ('Test Store', '', '2024-01-01', 'Test Item', 10.00, 1.00, 9.00);

-- Test valid insertion (this should succeed)
INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) 
VALUES ('Test Store', 'Test Vendor', '2024-01-01', 'Test Item', 10.00, 1.00, 9.00);

-- Verify the test record was inserted
SELECT * FROM sales_records WHERE store = 'Test Store';

-- Test the update trigger
UPDATE sales_records SET sale_price = 15.00 WHERE store = 'Test Store';

-- Verify updated_at was automatically updated
SELECT id, sale_price, created_at, updated_at FROM sales_records WHERE store = 'Test Store';

-- Clean up test record
DELETE FROM sales_records WHERE store = 'Test Store';

-- ============================================================================
-- INDEX PERFORMANCE TESTING
-- ============================================================================

-- Test date index performance
EXPLAIN QUERY PLAN SELECT * FROM sales_records WHERE date >= '2024-01-01' ORDER BY date DESC;

-- Test store index performance
EXPLAIN QUERY PLAN SELECT * FROM sales_records WHERE store = 'Downtown Store';

-- Test vendor index performance
EXPLAIN QUERY PLAN SELECT * FROM sales_records WHERE vendor = 'Electronics Plus';

-- Test composite index performance
EXPLAIN QUERY PLAN SELECT * FROM sales_records WHERE store = 'Downtown Store' AND date >= '2024-01-01' ORDER BY date DESC;

-- ============================================================================
-- VIEW TESTING
-- ============================================================================

-- Test yearly summary view
SELECT * FROM v_yearly_sales_summary;

-- Test monthly summary view (limit to recent months)
SELECT * FROM v_monthly_sales_summary LIMIT 10;

-- Test daily summary view (limit to recent days)
SELECT * FROM v_daily_sales_summary LIMIT 10;

-- Test store performance view
SELECT * FROM v_store_performance;

-- Test vendor performance view
SELECT * FROM v_vendor_performance;

-- ============================================================================
-- REPORTING QUERY TESTING
-- ============================================================================

-- Test pivot table style queries
SELECT 
    strftime('%Y', date) as year,
    strftime('%m', date) as month,
    COUNT(*) as items_sold,
    SUM(remaining) as sum_remaining
FROM sales_records 
GROUP BY strftime('%Y', date), strftime('%m', date)
ORDER BY year DESC, month DESC;

-- Test drill-down queries
SELECT 
    date,
    store,
    vendor,
    description,
    sale_price,
    commission,
    remaining
FROM sales_records 
WHERE strftime('%Y-%m', date) = '2024-06'
ORDER BY date DESC;

-- Test filtering by store
SELECT 
    date,
    vendor,
    description,
    sale_price,
    remaining
FROM sales_records 
WHERE store = 'Downtown Store'
ORDER BY date DESC;

-- Test filtering by vendor
SELECT 
    date,
    store,
    description,
    sale_price,
    remaining
FROM sales_records 
WHERE vendor = 'Electronics Plus'
ORDER BY date DESC;

-- Test date range queries
SELECT 
    date,
    store,
    vendor,
    description,
    sale_price,
    remaining
FROM sales_records 
WHERE date BETWEEN '2024-05-01' AND '2024-05-31'
ORDER BY date DESC;

-- ============================================================================
-- PERFORMANCE STATISTICS
-- ============================================================================

-- Database size information
SELECT 
    COUNT(*) as total_records,
    MIN(date) as earliest_date,
    MAX(date) as latest_date,
    SUM(sale_price) as total_sales,
    AVG(sale_price) as avg_sale_price,
    COUNT(DISTINCT store) as unique_stores,
    COUNT(DISTINCT vendor) as unique_vendors
FROM sales_records;

-- Index usage statistics (SQLite specific)
-- PRAGMA index_list(sales_records);
-- PRAGMA index_info(idx_sales_records_date);

-- ============================================================================
-- DATA INTEGRITY CHECKS
-- ============================================================================

-- Check for any null values in required fields
SELECT 'Null store values' as issue, COUNT(*) as count FROM sales_records WHERE store IS NULL OR store = '';
SELECT 'Null vendor values' as issue, COUNT(*) as count FROM sales_records WHERE vendor IS NULL OR vendor = '';
SELECT 'Null date values' as issue, COUNT(*) as count FROM sales_records WHERE date IS NULL;
SELECT 'Null description values' as issue, COUNT(*) as count FROM sales_records WHERE description IS NULL OR description = '';

-- Check for negative values
SELECT 'Negative sale_price' as issue, COUNT(*) as count FROM sales_records WHERE sale_price < 0;
SELECT 'Negative commission' as issue, COUNT(*) as count FROM sales_records WHERE commission < 0;
SELECT 'Negative remaining' as issue, COUNT(*) as count FROM sales_records WHERE remaining < 0;

-- Check for unrealistic values
SELECT 'Very high sale_price' as issue, COUNT(*) as count FROM sales_records WHERE sale_price > 10000;
SELECT 'Commission > sale_price' as issue, COUNT(*) as count FROM sales_records WHERE commission > sale_price;

-- ============================================================================
-- SUMMARY REPORT
-- ============================================================================

SELECT 'Database validation completed successfully' as status;

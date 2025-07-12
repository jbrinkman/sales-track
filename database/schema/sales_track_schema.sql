-- Sales Track Database Schema
-- SQLite database schema for sales data storage and reporting
-- Version: 1.0
-- Created: 2025-07-12

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- ============================================================================
-- SALES RECORDS TABLE
-- ============================================================================
-- Main table for storing individual sales transactions
-- Supports the core workflow of importing HTML table data and generating reports

CREATE TABLE IF NOT EXISTS sales_records (
    -- Primary key
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Core sales data fields (from HTML import)
    store VARCHAR(100) NOT NULL,                    -- Store identifier/name
    vendor VARCHAR(100) NOT NULL,                   -- Vendor/supplier name
    date DATE NOT NULL,                             -- Sale date (YYYY-MM-DD format)
    description TEXT NOT NULL,                      -- Item/product description
    sale_price DECIMAL(10,2) NOT NULL,             -- Sale amount (dollars.cents)
    commission DECIMAL(10,2) NOT NULL DEFAULT 0.00, -- Commission amount (dollars.cents)
    remaining DECIMAL(10,2) NOT NULL DEFAULT 0.00,  -- Remaining balance (dollars.cents)
    
    -- Metadata fields
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Last update timestamp
    
    -- Data validation constraints
    CONSTRAINT chk_sale_price_positive CHECK (sale_price >= 0),
    CONSTRAINT chk_commission_positive CHECK (commission >= 0),
    CONSTRAINT chk_remaining_positive CHECK (remaining >= 0),
    CONSTRAINT chk_date_format CHECK (date IS NOT NULL AND date != ''),
    CONSTRAINT chk_store_not_empty CHECK (LENGTH(TRIM(store)) > 0),
    CONSTRAINT chk_vendor_not_empty CHECK (LENGTH(TRIM(vendor)) > 0),
    CONSTRAINT chk_description_not_empty CHECK (LENGTH(TRIM(description)) > 0)
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================
-- Indexes optimized for reporting queries and data filtering

-- Primary reporting index: date-based queries (most common)
CREATE INDEX IF NOT EXISTS idx_sales_records_date 
ON sales_records(date DESC);

-- Store-based filtering and grouping
CREATE INDEX IF NOT EXISTS idx_sales_records_store 
ON sales_records(store);

-- Vendor-based filtering and grouping  
CREATE INDEX IF NOT EXISTS idx_sales_records_vendor 
ON sales_records(vendor);

-- Composite index for date range queries by store
CREATE INDEX IF NOT EXISTS idx_sales_records_store_date 
ON sales_records(store, date DESC);

-- Composite index for date range queries by vendor
CREATE INDEX IF NOT EXISTS idx_sales_records_vendor_date 
ON sales_records(vendor, date DESC);

-- Composite index for pivot table queries (year/month grouping)
CREATE INDEX IF NOT EXISTS idx_sales_records_date_store_vendor 
ON sales_records(date DESC, store, vendor);

-- Index for financial calculations and remaining balance queries
CREATE INDEX IF NOT EXISTS idx_sales_records_remaining 
ON sales_records(remaining DESC) WHERE remaining > 0;

-- Metadata index for data management
CREATE INDEX IF NOT EXISTS idx_sales_records_created_at 
ON sales_records(created_at DESC);

-- ============================================================================
-- TRIGGERS FOR AUTOMATIC TIMESTAMP UPDATES
-- ============================================================================
-- Automatically update the updated_at timestamp when records are modified

CREATE TRIGGER IF NOT EXISTS trg_sales_records_updated_at
    AFTER UPDATE ON sales_records
    FOR EACH ROW
BEGIN
    UPDATE sales_records 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE id = NEW.id;
END;

-- ============================================================================
-- VIEWS FOR COMMON REPORTING QUERIES
-- ============================================================================

-- View for yearly sales summary (pivot table foundation)
CREATE VIEW IF NOT EXISTS v_yearly_sales_summary AS
SELECT 
    strftime('%Y', date) as year,
    COUNT(*) as items_sold,
    SUM(sale_price) as total_sales,
    SUM(commission) as total_commission,
    SUM(remaining) as total_remaining,
    COUNT(DISTINCT store) as unique_stores,
    COUNT(DISTINCT vendor) as unique_vendors
FROM sales_records
GROUP BY strftime('%Y', date)
ORDER BY year DESC;

-- View for monthly sales summary (drill-down level 1)
CREATE VIEW IF NOT EXISTS v_monthly_sales_summary AS
SELECT 
    strftime('%Y', date) as year,
    strftime('%m', date) as month,
    strftime('%Y-%m', date) as year_month,
    COUNT(*) as items_sold,
    SUM(sale_price) as total_sales,
    SUM(commission) as total_commission,
    SUM(remaining) as total_remaining,
    COUNT(DISTINCT store) as unique_stores,
    COUNT(DISTINCT vendor) as unique_vendors
FROM sales_records
GROUP BY strftime('%Y-%m', date)
ORDER BY year DESC, month DESC;

-- View for daily sales summary (drill-down level 2)
CREATE VIEW IF NOT EXISTS v_daily_sales_summary AS
SELECT 
    date,
    strftime('%Y', date) as year,
    strftime('%m', date) as month,
    strftime('%d', date) as day,
    strftime('%Y-%m', date) as year_month,
    COUNT(*) as items_sold,
    SUM(sale_price) as total_sales,
    SUM(commission) as total_commission,
    SUM(remaining) as total_remaining,
    COUNT(DISTINCT store) as unique_stores,
    COUNT(DISTINCT vendor) as unique_vendors
FROM sales_records
GROUP BY date
ORDER BY date DESC;

-- View for store performance analysis
CREATE VIEW IF NOT EXISTS v_store_performance AS
SELECT 
    store,
    COUNT(*) as total_items,
    SUM(sale_price) as total_sales,
    SUM(commission) as total_commission,
    SUM(remaining) as total_remaining,
    AVG(sale_price) as avg_sale_price,
    MIN(date) as first_sale_date,
    MAX(date) as last_sale_date,
    COUNT(DISTINCT vendor) as unique_vendors
FROM sales_records
GROUP BY store
ORDER BY total_sales DESC;

-- View for vendor performance analysis
CREATE VIEW IF NOT EXISTS v_vendor_performance AS
SELECT 
    vendor,
    COUNT(*) as total_items,
    SUM(sale_price) as total_sales,
    SUM(commission) as total_commission,
    SUM(remaining) as total_remaining,
    AVG(sale_price) as avg_sale_price,
    MIN(date) as first_sale_date,
    MAX(date) as last_sale_date,
    COUNT(DISTINCT store) as unique_stores
FROM sales_records
GROUP BY vendor
ORDER BY total_sales DESC;

-- ============================================================================
-- SAMPLE DATA FOR TESTING (OPTIONAL)
-- ============================================================================
-- Uncomment the following section to insert sample data for development/testing

/*
INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) VALUES
('Store A', 'Vendor 1', '2024-01-15', 'Product Alpha', 150.00, 15.00, 135.00),
('Store A', 'Vendor 2', '2024-01-16', 'Product Beta', 200.00, 20.00, 180.00),
('Store B', 'Vendor 1', '2024-01-17', 'Product Gamma', 75.50, 7.55, 67.95),
('Store A', 'Vendor 3', '2024-02-01', 'Product Delta', 300.00, 30.00, 270.00),
('Store C', 'Vendor 2', '2024-02-15', 'Product Epsilon', 125.75, 12.58, 113.17),
('Store B', 'Vendor 1', '2024-03-01', 'Product Zeta', 89.99, 9.00, 80.99),
('Store A', 'Vendor 1', '2024-03-15', 'Product Eta', 175.25, 17.53, 157.72);
*/

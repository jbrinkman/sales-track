-- Migration: 001_initial_schema.sql
-- Description: Initial database schema for Sales Track application
-- Created: 2025-07-12
-- Version: 1.0

-- This migration creates the foundational database structure for the Sales Track application
-- It includes the main sales_records table, indexes, triggers, and views for reporting

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- ============================================================================
-- CREATE SALES RECORDS TABLE
-- ============================================================================

CREATE TABLE sales_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    store VARCHAR(100) NOT NULL,
    vendor VARCHAR(100) NOT NULL,
    date DATE NOT NULL,
    description TEXT NOT NULL,
    sale_price DECIMAL(10,2) NOT NULL,
    commission DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    remaining DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_sale_price_positive CHECK (sale_price >= 0),
    CONSTRAINT chk_commission_positive CHECK (commission >= 0),
    CONSTRAINT chk_remaining_positive CHECK (remaining >= 0),
    CONSTRAINT chk_date_format CHECK (date IS NOT NULL AND date != ''),
    CONSTRAINT chk_store_not_empty CHECK (LENGTH(TRIM(store)) > 0),
    CONSTRAINT chk_vendor_not_empty CHECK (LENGTH(TRIM(vendor)) > 0),
    CONSTRAINT chk_description_not_empty CHECK (LENGTH(TRIM(description)) > 0)
);

-- ============================================================================
-- CREATE INDEXES
-- ============================================================================

CREATE INDEX idx_sales_records_date ON sales_records(date DESC);
CREATE INDEX idx_sales_records_store ON sales_records(store);
CREATE INDEX idx_sales_records_vendor ON sales_records(vendor);
CREATE INDEX idx_sales_records_store_date ON sales_records(store, date DESC);
CREATE INDEX idx_sales_records_vendor_date ON sales_records(vendor, date DESC);
CREATE INDEX idx_sales_records_date_store_vendor ON sales_records(date DESC, store, vendor);
CREATE INDEX idx_sales_records_remaining ON sales_records(remaining DESC) WHERE remaining > 0;
CREATE INDEX idx_sales_records_created_at ON sales_records(created_at DESC);

-- ============================================================================
-- CREATE TRIGGERS
-- ============================================================================

CREATE TRIGGER trg_sales_records_updated_at
    AFTER UPDATE ON sales_records
    FOR EACH ROW
BEGIN
    UPDATE sales_records 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE id = NEW.id;
END;

-- ============================================================================
-- CREATE VIEWS
-- ============================================================================

CREATE VIEW v_yearly_sales_summary AS
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

CREATE VIEW v_monthly_sales_summary AS
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

CREATE VIEW v_daily_sales_summary AS
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

CREATE VIEW v_store_performance AS
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

CREATE VIEW v_vendor_performance AS
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

# Sales Track Database

This directory contains the database schema, migrations, and related files for the Sales Track application.

## Directory Structure

```
database/
├── README.md                    # This file
├── DATABASE_DESIGN.md          # Comprehensive database design documentation
├── schema/
│   └── sales_track_schema.sql  # Complete database schema with comments
├── migrations/
│   └── 001_initial_schema.sql  # Initial database migration
├── sample_data.sql             # Sample data for development/testing
└── validate_schema.sql         # Schema validation and testing script
```

## Quick Start

### 1. Create Database
```bash
# Create a new SQLite database
sqlite3 sales_track.db < database/migrations/001_initial_schema.sql
```

### 2. Add Sample Data (Optional)
```bash
# Add sample data for development/testing
sqlite3 sales_track.db < database/sample_data.sql
```

### 3. Validate Schema
```bash
# Run validation tests
sqlite3 sales_track.db < database/validate_schema.sql
```

## Files Description

### Schema Files
- **`schema/sales_track_schema.sql`** - Complete database schema with detailed comments, constraints, indexes, views, and triggers
- **`migrations/001_initial_schema.sql`** - Migration file for initial database setup

### Data Files
- **`sample_data.sql`** - Realistic sample data spanning 6 months with multiple stores and vendors
- **`validate_schema.sql`** - Comprehensive validation script to test schema integrity

### Documentation
- **`DATABASE_DESIGN.md`** - Detailed documentation covering design decisions, performance considerations, and scalability planning

## Database Schema Overview

### Main Table: `sales_records`
- **Primary Key**: `id` (auto-increment)
- **Core Fields**: store, vendor, date, description, sale_price, commission, remaining
- **Metadata**: created_at, updated_at (with automatic trigger)
- **Constraints**: Positive values, non-empty strings, valid dates

### Indexes (Performance Optimized)
- Date-based queries (primary use case)
- Store and vendor filtering
- Composite indexes for complex queries
- Partial index for outstanding balances

### Views (Reporting)
- `v_yearly_sales_summary` - Annual aggregations
- `v_monthly_sales_summary` - Monthly drill-down
- `v_daily_sales_summary` - Daily drill-down
- `v_store_performance` - Store analytics
- `v_vendor_performance` - Vendor analytics

### Triggers
- Automatic `updated_at` timestamp maintenance

## Usage Examples

### Basic Queries
```sql
-- Get recent sales
SELECT * FROM sales_records ORDER BY date DESC LIMIT 10;

-- Sales by store
SELECT store, COUNT(*), SUM(sale_price) FROM sales_records GROUP BY store;

-- Monthly summary
SELECT * FROM v_monthly_sales_summary WHERE year = '2024';
```

### Reporting Queries
```sql
-- Pivot table style (Year > Month > Date)
SELECT year, month, items_sold, total_remaining 
FROM v_monthly_sales_summary 
ORDER BY year DESC, month DESC;

-- Drill-down to specific month
SELECT * FROM sales_records 
WHERE strftime('%Y-%m', date) = '2024-06' 
ORDER BY date DESC;
```

## Performance Notes

- **Optimized for**: Date range queries, store/vendor filtering, pivot table reporting
- **Expected Performance**: Sub-millisecond for typical queries with proper indexing
- **Scalability**: Designed to handle 100K+ records efficiently
- **Memory Usage**: Minimal overhead, suitable for desktop application

## Migration Strategy

- **Sequential Numbering**: Migration files numbered 001, 002, etc.
- **Forward Only**: Each migration builds on the previous
- **Data Preservation**: Migrations designed to preserve existing data
- **Rollback**: Include rollback instructions in migration comments

## Development Workflow

1. **Schema Changes**: Create new migration file
2. **Testing**: Use sample data and validation script
3. **Documentation**: Update DATABASE_DESIGN.md
4. **Integration**: Test with Go application layer

## Backup and Recovery

### Backup
```bash
# Full database backup
cp sales_track.db sales_track_backup_$(date +%Y%m%d).db

# SQL dump backup
sqlite3 sales_track.db .dump > sales_track_backup.sql
```

### Recovery
```bash
# Restore from file backup
cp sales_track_backup_20240712.db sales_track.db

# Restore from SQL dump
sqlite3 sales_track_new.db < sales_track_backup.sql
```

## Troubleshooting

### Common Issues
1. **Constraint Violations**: Check data validation in import process
2. **Performance Issues**: Verify indexes are being used (EXPLAIN QUERY PLAN)
3. **Lock Issues**: Ensure proper connection management in application

### Validation
Run the validation script to check schema integrity:
```bash
sqlite3 sales_track.db < database/validate_schema.sql
```

## Integration with Go Application

The database layer will be accessed through:
- **Database Package**: `internal/database` or similar
- **Models**: Go structs matching table schema
- **Queries**: Prepared statements for performance
- **Migrations**: Automated migration runner

See the Go application code for specific implementation details.

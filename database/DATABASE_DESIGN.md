# Sales Track Database Design

This document describes the database schema design for the Sales Track application, including design decisions, performance considerations, and future scalability planning.

## Overview

The Sales Track database is designed to efficiently store and query sales transaction data imported from HTML tables. The schema prioritizes:

- **Performance**: Optimized for pivot table-style reporting queries
- **Data Integrity**: Comprehensive constraints and validation
- **Scalability**: Indexed for large datasets and future growth
- **Simplicity**: Single main table with supporting views

## Database Technology

**SQLite** was chosen for the following reasons:
- **Local Storage**: No server setup required for desktop application
- **Performance**: Excellent for read-heavy workloads with proper indexing
- **Reliability**: ACID compliance and crash-safe transactions
- **Portability**: Single file database, easy backup and distribution
- **Go Integration**: Excellent SQLite drivers available for Go

## Schema Design

### Core Table: `sales_records`

The main table stores individual sales transactions with the following structure:

#### Primary Fields (from HTML import)
- `id` - Auto-incrementing primary key
- `store` - Store identifier/name (VARCHAR 100)
- `vendor` - Vendor/supplier name (VARCHAR 100) 
- `date` - Sale date in YYYY-MM-DD format (DATE)
- `description` - Item/product description (TEXT)
- `sale_price` - Sale amount in dollars.cents (DECIMAL 10,2)
- `commission` - Commission amount in dollars.cents (DECIMAL 10,2)
- `remaining` - Remaining balance in dollars.cents (DECIMAL 10,2)

#### Metadata Fields
- `created_at` - Record creation timestamp (DATETIME)
- `updated_at` - Last modification timestamp (DATETIME)

### Data Types Rationale

| Field | Type | Rationale |
|-------|------|-----------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | SQLite optimization, efficient joins |
| `store` | VARCHAR(100) | Reasonable length for store names, indexed |
| `vendor` | VARCHAR(100) | Reasonable length for vendor names, indexed |
| `date` | DATE | Proper date type for range queries and grouping |
| `description` | TEXT | Variable length for product descriptions |
| `sale_price` | DECIMAL(10,2) | Precise currency handling, up to $99,999,999.99 |
| `commission` | DECIMAL(10,2) | Precise currency handling, matches sale_price |
| `remaining` | DECIMAL(10,2) | Precise currency handling, matches sale_price |
| `created_at` | DATETIME | Full timestamp for audit trail |
| `updated_at` | DATETIME | Full timestamp for change tracking |

## Constraints and Validation

### Data Integrity Constraints
- **Positive Values**: All monetary fields must be >= 0
- **Non-Empty Fields**: Store, vendor, and description cannot be empty
- **Date Validation**: Date field cannot be null or empty
- **String Validation**: Text fields are trimmed and checked for content

### Business Logic Constraints
- Sale price, commission, and remaining must be non-negative
- Store and vendor names must contain actual content (not just whitespace)
- Dates must be provided for all transactions

## Indexing Strategy

The indexing strategy is optimized for the application's primary use cases:

### Primary Indexes
1. **Date Index** (`idx_sales_records_date`)
   - Most common query pattern for reporting
   - Descending order for recent-first display

2. **Store Index** (`idx_sales_records_store`)
   - Store-based filtering and grouping

3. **Vendor Index** (`idx_sales_records_vendor`)
   - Vendor-based filtering and grouping

### Composite Indexes
4. **Store + Date** (`idx_sales_records_store_date`)
   - Optimizes store-specific date range queries

5. **Vendor + Date** (`idx_sales_records_vendor_date`)
   - Optimizes vendor-specific date range queries

6. **Date + Store + Vendor** (`idx_sales_records_date_store_vendor`)
   - Optimizes complex pivot table queries

### Specialized Indexes
7. **Remaining Balance** (`idx_sales_records_remaining`)
   - Partial index for outstanding balances only
   - Optimizes financial reporting queries

8. **Created At** (`idx_sales_records_created_at`)
   - Data management and audit queries

## Views for Reporting

Pre-built views optimize common reporting patterns:

### Hierarchical Summary Views
- `v_yearly_sales_summary` - Top-level pivot table data
- `v_monthly_sales_summary` - Month-level drill-down
- `v_daily_sales_summary` - Day-level drill-down

### Performance Analysis Views
- `v_store_performance` - Store-based analytics
- `v_vendor_performance` - Vendor-based analytics

### View Benefits
- **Performance**: Pre-aggregated calculations
- **Consistency**: Standardized business logic
- **Simplicity**: Complex queries abstracted into simple SELECT statements

## Triggers

### Automatic Timestamp Updates
- `trg_sales_records_updated_at` - Automatically updates `updated_at` on record modification
- Ensures accurate change tracking without application-level management

## Performance Considerations

### Query Optimization
- **Date Range Queries**: Optimized with descending date index
- **Grouping Operations**: Composite indexes support GROUP BY operations
- **Filtering**: Individual field indexes support WHERE clauses
- **Sorting**: Index order matches common sort requirements

### Expected Performance
- **Small Dataset** (< 10K records): Sub-millisecond queries
- **Medium Dataset** (10K - 100K records): Single-digit millisecond queries
- **Large Dataset** (100K+ records): Still performant with proper indexing

### Memory Usage
- **Index Overhead**: ~20-30% of table size for all indexes
- **View Storage**: Views are virtual, no additional storage
- **SQLite Efficiency**: Minimal memory footprint for desktop application

## Scalability Planning

### Current Capacity
- **Record Limit**: Practically unlimited (SQLite supports 281TB databases)
- **Performance**: Optimized for up to 1M+ records with current indexing
- **Storage**: Efficient storage with minimal overhead

### Future Enhancements
1. **Additional Tables**: Easy to add related tables (categories, customers, etc.)
2. **Partitioning**: Date-based partitioning if needed for very large datasets
3. **Archive Strategy**: Older data can be moved to separate archive tables
4. **Replication**: SQLite supports backup and replication strategies

### Migration Strategy
- **Version Control**: Migration files numbered sequentially
- **Rollback Support**: Each migration can include rollback instructions
- **Data Preservation**: Migrations designed to preserve existing data

## Data Import Considerations

### HTML Table Mapping
The schema directly maps to the expected HTML table structure:
- Column order matches typical HTML table layout
- Data types handle common formatting variations
- Constraints validate imported data quality

### Error Handling
- **Invalid Data**: Constraints prevent bad data insertion
- **Duplicate Detection**: Can be implemented at application level
- **Data Cleaning**: Text fields trimmed and validated

## Backup and Recovery

### Backup Strategy
- **File-based**: Simple file copy of SQLite database
- **Export Options**: SQL dump, CSV export via views
- **Incremental**: Based on `created_at` and `updated_at` timestamps

### Recovery Options
- **Point-in-time**: Using backup files and transaction logs
- **Data Validation**: Constraints ensure data integrity after recovery
- **Migration Replay**: Can rebuild database from migration files

## Security Considerations

### Data Protection
- **Local Storage**: Data remains on user's machine
- **No Network Exposure**: SQLite doesn't expose network ports
- **File Permissions**: Standard OS file permissions apply

### Input Validation
- **SQL Injection**: Parameterized queries prevent injection attacks
- **Data Validation**: Constraints prevent malformed data
- **Type Safety**: Strong typing prevents data corruption

## Testing Strategy

### Unit Testing
- **Schema Validation**: Test all constraints and triggers
- **Index Performance**: Verify query performance with test data
- **View Accuracy**: Validate aggregation calculations

### Integration Testing
- **Data Import**: Test with various HTML table formats
- **Reporting**: Verify pivot table calculations
- **Performance**: Load testing with large datasets

### Sample Data
- **Development**: Sample data included in schema file
- **Testing**: Comprehensive test datasets for various scenarios
- **Performance**: Large datasets for performance validation

## Conclusion

This database design provides a solid foundation for the Sales Track application with:

- **Efficient Storage**: Normalized structure with minimal redundancy
- **Fast Queries**: Comprehensive indexing for all common access patterns
- **Data Integrity**: Robust constraints and validation rules
- **Scalability**: Designed to handle growth in data volume and complexity
- **Maintainability**: Clear structure with good documentation

The schema supports the core application requirements while providing flexibility for future enhancements and optimizations.

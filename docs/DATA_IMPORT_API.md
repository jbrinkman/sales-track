# Data Import API Documentation

The Data Import API provides Wails context methods for importing HTML table data into the Sales Track application. This API bridges the HTML parser and database layer to provide a complete data import solution.

## Overview

The Data Import API consists of several methods exposed through the Wails application context:

- **ImportHTMLData** - Import HTML table data with individual record processing
- **ImportHTMLDataBatch** - Import HTML table data using batch operations for better performance
- **ImportHTMLDataWithOptions** - Import with configurable parsing options
- **ValidateHTMLData** - Validate HTML data without importing
- **GetImportStatistics** - Get statistics about imported data
- **GetDatabaseHealth** - Check database connection health
- **GetRecentImports** - Get recently imported records

## API Methods

### ImportHTMLData

Imports HTML table data into the database with individual record processing.

**Signature:**
```go
func (a *App) ImportHTMLData(htmlData string) (*ImportResult, error)
```

**Parameters:**
- `htmlData` (string) - HTML table data to import

**Returns:**
- `ImportResult` - Detailed import results
- `error` - Error if the operation fails

**Example Usage:**
```javascript
// Frontend JavaScript
const htmlData = `
<table>
  <tr>
    <th>Store</th>
    <th>Vendor</th>
    <th>Date</th>
    <th>Description</th>
    <th>Sale Price</th>
    <th>Commission</th>
    <th>Remaining</th>
  </tr>
  <tr>
    <td>Downtown Store</td>
    <td>Electronics Plus</td>
    <td>2024-01-15</td>
    <td>Samsung TV</td>
    <td>$899.99</td>
    <td>$89.99</td>
    <td>$810.00</td>
  </tr>
</table>
`;

const result = await ImportHTMLData(htmlData);
console.log(`Imported ${result.imported_rows} of ${result.total_rows} records`);
```

### ImportHTMLDataBatch

Imports HTML table data using batch operations for better performance with large datasets.

**Signature:**
```go
func (a *App) ImportHTMLDataBatch(htmlData string) (*ImportResult, error)
```

**Performance:**
- ~175μs per operation (individual records)
- ~162μs per operation (batch processing)
- ~7% performance improvement for batch operations

### ImportHTMLDataWithOptions

Imports HTML data with configurable parsing options.

**Signature:**
```go
func (a *App) ImportHTMLDataWithOptions(htmlData string, options ImportOptions) (*ImportResult, error)
```

**ImportOptions:**
```go
type ImportOptions struct {
    UseConsignableFormat bool     `json:"use_consignable_format"`
    CustomColumnMapping  []string `json:"custom_column_mapping,omitempty"`
    StrictMode           bool     `json:"strict_mode"`
    UseBatchImport       bool     `json:"use_batch_import"`
}
```

**Example - Consignable Format:**
```javascript
const options = {
    use_consignable_format: true,
    use_batch_import: true
};

const result = await ImportHTMLDataWithOptions(htmlData, options);
```

**Example - Custom Column Mapping:**
```javascript
const options = {
    custom_column_mapping: ["vendor", "store", "date", "description", "sale_price"],
    strict_mode: true
};

const result = await ImportHTMLDataWithOptions(htmlData, options);
```

### ValidateHTMLData

Validates HTML data without importing to check for parsing errors.

**Signature:**
```go
func (a *App) ValidateHTMLData(htmlData string) (*ValidationResult, error)
```

**Returns ValidationResult:**
```go
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
```

**Example Usage:**
```javascript
const validation = await ValidateHTMLData(htmlData);
if (!validation.valid) {
    console.log(`Validation failed: ${validation.error_message}`);
    validation.errors.forEach(error => {
        console.log(`Row ${error.row}, Column ${error.column}: ${error.message}`);
    });
}
```

### GetImportStatistics

Returns statistics about imported data.

**Signature:**
```go
func (a *App) GetImportStatistics() (*ImportStatistics, error)
```

**Returns ImportStatistics:**
```go
type ImportStatistics struct {
    TotalRecords  int     `json:"total_records"`
    RecentRecords int     `json:"recent_records"`  // Last 30 days
    TotalSales    float64 `json:"total_sales"`
    AveragePrice  float64 `json:"average_price"`
}
```

### GetDatabaseHealth

Checks database connection health.

**Signature:**
```go
func (a *App) GetDatabaseHealth() (*DatabaseHealth, error)
```

**Returns DatabaseHealth:**
```go
type DatabaseHealth struct {
    Connected bool   `json:"connected"`
    Error     string `json:"error,omitempty"`
}
```

### GetRecentImports

Returns recently imported sales records.

**Signature:**
```go
func (a *App) GetRecentImports(limit int) ([]models.SalesRecord, error)
```

**Parameters:**
- `limit` (int) - Maximum number of records to return

**Returns:**
- Array of `SalesRecord` objects sorted by creation date (newest first)

## Data Types

### ImportResult

Complete result of an import operation:

```go
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
```

### ImportError

Error that occurred during database import:

```go
type ImportError struct {
    Record models.CreateSalesRecordRequest `json:"record"`
    Error  string                          `json:"error"`
}
```

## Supported HTML Formats

### Standard HTML Tables
```html
<table>
  <tr>
    <th>Store</th>
    <th>Vendor</th>
    <th>Date</th>
    <th>Description</th>
    <th>Sale Price</th>
    <th>Commission</th>
    <th>Remaining</th>
  </tr>
  <tr>
    <td>Store Name</td>
    <td>Vendor Name</td>
    <td>2024-01-15</td>
    <td>Product Description</td>
    <td>$100.00</td>
    <td>$10.00</td>
    <td>$90.00</td>
  </tr>
</table>
```

### Headerless Rows (Consignable Format)
```html
<tr>
  <td>Store Name</td>
  <td>Vendor Name</td>
  <td>2024-01-15</td>
  <td>Product Description</td>
  <td>$100.00</td>
  <td>$10.00</td>
  <td>$90.00</td>
</tr>
```

### Tab-Delimited Data
```
Store	Vendor	Date	Description	Sale Price	Commission	Remaining
Store Name	Vendor Name	2024-01-15	Product Description	100.00	10.00	90.00
```

## Column Recognition

The parser intelligently recognizes various column name variations:

- **Store**: store, shop, location, outlet, branch, store name, shop name
- **Vendor**: vendor, supplier, brand, manufacturer, company, vendor name, supplier name
- **Date**: date, sale date, transaction date, order date, purchase date, sold date
- **Description**: description, item, product, item description, product description, details, name, product name
- **Sale Price**: sale price, price, amount, total, sale amount, selling price, cost, value
- **Commission**: commission, fee, commission amount, commission fee, comm, commission %, commission rate
- **Remaining**: remaining, balance, remaining balance, outstanding, due, remaining amount, balance due

## Error Handling

The API provides comprehensive error handling:

### Parse Errors
- Invalid HTML structure
- Missing required columns
- Invalid data formats (dates, currency)
- Empty required fields

### Import Errors
- Database connection issues
- Constraint violations
- Transaction failures

### Example Error Response
```json
{
  "success": false,
  "total_rows": 5,
  "parsed_rows": 3,
  "imported_rows": 2,
  "error_message": "Imported 2 of 3 records. 1 records failed to import.",
  "parse_errors": [
    {
      "row": 4,
      "column": "date",
      "message": "Invalid date format: unable to parse date: invalid-date",
      "value": "invalid-date"
    }
  ],
  "import_errors": [
    {
      "record": {...},
      "error": "duplicate key value violates unique constraint"
    }
  ]
}
```

## Performance Characteristics

### Benchmarks (Apple M3 Max)
- **Individual Import**: ~175μs per operation, 27KB memory, 435 allocations
- **Batch Import**: ~162μs per operation, 27KB memory, 407 allocations
- **Validation Only**: ~98μs per operation (HTML parsing only)

### Scalability
- **Small datasets** (1-10 records): Sub-millisecond processing
- **Medium datasets** (10-100 records): 1-10ms processing
- **Large datasets** (100-1000 records): 10-100ms processing

## Best Practices

### 1. Use Batch Import for Large Datasets
```javascript
const options = {
    use_batch_import: true
};
const result = await ImportHTMLDataWithOptions(htmlData, options);
```

### 2. Validate Before Importing
```javascript
const validation = await ValidateHTMLData(htmlData);
if (validation.valid) {
    const result = await ImportHTMLData(htmlData);
}
```

### 3. Handle Partial Imports
```javascript
const result = await ImportHTMLData(htmlData);
if (result.imported_rows < result.parsed_rows) {
    console.log(`Warning: Only ${result.imported_rows} of ${result.parsed_rows} records imported`);
    result.import_errors.forEach(error => {
        console.log(`Import error: ${error.error}`);
    });
}
```

### 4. Monitor Database Health
```javascript
const health = await GetDatabaseHealth();
if (!health.connected) {
    console.error(`Database error: ${health.error}`);
}
```

### 5. Use Consignable Format for Headerless Data
```javascript
const options = {
    use_consignable_format: true
};
const result = await ImportHTMLDataWithOptions(headerlessHTML, options);
```

## Integration Examples

### Vue.js Component Integration
```vue
<template>
  <div class="import-section">
    <textarea v-model="htmlData" placeholder="Paste HTML table data here"></textarea>
    <button @click="importData" :disabled="importing">
      {{ importing ? 'Importing...' : 'Import Data' }}
    </button>
    <div v-if="result" class="result">
      <p>Imported {{ result.imported_rows }} of {{ result.total_rows }} records</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      htmlData: '',
      importing: false,
      result: null
    }
  },
  methods: {
    async importData() {
      this.importing = true;
      try {
        this.result = await ImportHTMLData(this.htmlData);
      } catch (error) {
        console.error('Import failed:', error);
      } finally {
        this.importing = false;
      }
    }
  }
}
</script>
```

## Error Recovery

The API is designed for graceful error recovery:

1. **Parse errors** don't stop processing - valid rows are still imported
2. **Import errors** are isolated - one failed record doesn't affect others
3. **Detailed error reporting** helps users fix data issues
4. **Partial success** is supported - some records can be imported even if others fail

This comprehensive API provides a robust foundation for the Sales Track application's data import functionality.

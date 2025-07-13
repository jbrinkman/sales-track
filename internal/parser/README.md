# HTML Table Parser

This package provides comprehensive HTML table parsing functionality for the Sales Track application, designed to extract sales data from various HTML table formats commonly found on websites.

## Overview

The HTML Table Parser is a core component that enables users to copy HTML table data from websites and import it directly into the Sales Track application. It handles various HTML formats, data validation, and provides detailed feedback about the parsing process.

## Features

### 🔍 **Flexible HTML Parsing**
- **Multiple Table Formats**: Handles standard HTML tables, tables with CSS classes, and nested structures
- **Automatic Table Detection**: Finds and selects the best table when multiple tables are present
- **Delimited Data Support**: Converts tab-separated and pipe-separated data to HTML tables
- **Robust HTML Processing**: Handles malformed HTML and various encoding issues

### 📊 **Intelligent Column Mapping**
- **Flexible Column Names**: Recognizes various column name variations (e.g., "Store", "Shop", "Location")
- **Smart Matching**: Uses fuzzy matching to map columns even with different naming conventions
- **Required Field Validation**: Ensures all essential columns are present
- **Optional Field Handling**: Gracefully handles missing optional columns

### 💰 **Advanced Data Type Parsing**
- **Currency Parsing**: Handles various currency formats ($, €, £, ¥) with commas and parentheses
- **Date Parsing**: Supports multiple date formats (ISO, US, European, natural language)
- **Number Validation**: Validates numeric data with proper error handling
- **Text Normalization**: Cleans and normalizes text data

### 🛡️ **Comprehensive Error Handling**
- **Detailed Error Messages**: Provides specific error information for each parsing issue
- **Row-Level Validation**: Reports errors at the individual row and column level
- **Warning System**: Distinguishes between critical errors and minor warnings
- **Parsing Statistics**: Provides detailed statistics about the parsing operation

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "sales-track/internal/parser"
)

func main() {
    // Create a new parser
    htmlParser := parser.NewHTMLTableParser()
    
    // HTML table data (copied from a website)
    htmlData := `
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
    `
    
    // Parse the HTML
    result, err := htmlParser.ParseHTML(htmlData)
    if err != nil {
        fmt.Printf("Parsing failed: %v\n", err)
        return
    }
    
    // Display results
    fmt.Printf("Successfully parsed %d records\n", result.SuccessCount)
    fmt.Printf("Errors: %d, Warnings: %d\n", result.ErrorCount, len(result.Warnings))
    
    // Access parsed records
    for i, record := range result.Records {
        fmt.Printf("Record %d: %s - %s - $%.2f\n", 
            i+1, record.Store, record.Description, record.SalePrice)
    }
}
```

### Advanced Configuration

```go
// Create parser with strict mode (requires exact column matches)
parser := &parser.HTMLTableParser{
    StrictMode: true,
}

result, err := parser.ParseHTML(htmlData)
```

## Supported Data Formats

### HTML Table Formats

#### Standard HTML Table
```html
<table>
    <tr>
        <th>Store</th>
        <th>Vendor</th>
        <th>Date</th>
        <th>Description</th>
        <th>Sale Price</th>
    </tr>
    <tr>
        <td>Store A</td>
        <td>Vendor 1</td>
        <td>2024-01-15</td>
        <td>Product X</td>
        <td>$100.00</td>
    </tr>
</table>
```

#### Table with CSS Classes and Styling
```html
<table class="data-table" border="1" style="width:100%">
    <thead>
        <tr style="background-color: #f0f0f0;">
            <th class="header-cell">Store Location</th>
            <th class="header-cell">Vendor Name</th>
            <!-- ... -->
        </tr>
    </thead>
    <tbody>
        <tr class="data-row">
            <td class="data-cell">Downtown Branch</td>
            <td class="data-cell">Tech Solutions Inc.</td>
            <!-- ... -->
        </tr>
    </tbody>
</table>
```

### Delimited Data Formats

#### Tab-Separated Values
```
Store	Vendor	Date	Description	Sale Price
Store A	Vendor 1	2024-01-15	Product X	100.00
Store B	Vendor 2	2024-01-16	Product Y	200.00
```

#### Pipe-Separated Values
```
Store|Vendor|Date|Description|Sale Price
Store A|Vendor 1|2024-01-15|Product X|100.00
Store B|Vendor 2|2024-01-16|Product Y|200.00
```

## Column Mapping

The parser recognizes various column name variations:

### Store Column
- `store`, `shop`, `location`, `outlet`, `branch`
- `store name`, `shop name`

### Vendor Column
- `vendor`, `supplier`, `brand`, `manufacturer`, `company`
- `vendor name`, `supplier name`

### Date Column
- `date`, `sale date`, `transaction date`, `order date`
- `purchase date`, `sold date`

### Description Column
- `description`, `item`, `product`, `item description`
- `product description`, `details`, `name`, `product name`

### Sale Price Column
- `sale price`, `price`, `amount`, `total`, `sale amount`
- `selling price`, `cost`, `value`

### Commission Column
- `commission`, `fee`, `commission amount`, `commission fee`
- `comm`, `commission %`, `commission rate`

### Remaining Column
- `remaining`, `balance`, `remaining balance`, `outstanding`
- `due`, `remaining amount`, `balance due`

## Data Type Support

### Currency Formats
- **Dollar**: `$100.00`, `$1,234.56`
- **Euro**: `€123.45`
- **Pound**: `£99.99`
- **Yen**: `¥1000`
- **Plain Numbers**: `100.00`, `1,234.56`
- **Negative Values**: `(50.00)`, `-50.00`

### Date Formats
- **ISO Format**: `2024-01-15`
- **US Format**: `01/15/2024`, `1/15/2024`
- **European Format**: `15/01/2024`
- **Natural Language**: `Jan 15, 2024`, `January 15, 2024`
- **Alternative**: `15 Jan 2024`, `2024/01/15`

## Error Handling

### Error Types

#### Parse Errors (Critical)
- Missing required fields
- Invalid data formats
- Unparseable values

#### Warnings (Non-Critical)
- Optional field parsing issues
- Data type assumptions
- Format inconsistencies

### Error Information
```go
type ParseError struct {
    Row     int    `json:"row"`          // Row number (1-based)
    Column  string `json:"column"`       // Column name
    Message string `json:"message"`      // Error description
    Value   string `json:"value"`        // Original value that caused error
}
```

### Example Error Handling
```go
result, err := parser.ParseHTML(htmlData)
if err != nil {
    // Critical parsing failure
    fmt.Printf("Parsing failed: %v\n", err)
    return
}

// Check for row-level errors
if result.ErrorCount > 0 {
    fmt.Printf("Found %d rows with errors:\n", result.ErrorCount)
    for _, parseErr := range result.Errors {
        fmt.Printf("Row %d, Column '%s': %s (Value: '%s')\n", 
            parseErr.Row, parseErr.Column, parseErr.Message, parseErr.Value)
    }
}

// Check for warnings
if len(result.Warnings) > 0 {
    fmt.Printf("Found %d warnings:\n", len(result.Warnings))
    for _, warning := range result.Warnings {
        fmt.Printf("Row %d, Column '%s': %s\n", 
            warning.Row, warning.Column, warning.Message)
    }
}
```

## Parse Results

### ParseResult Structure
```go
type ParseResult struct {
    Records       []models.CreateSalesRecordRequest // Successfully parsed records
    TotalRows     int                               // Total data rows processed
    SuccessCount  int                               // Successfully parsed rows
    ErrorCount    int                               // Rows with critical errors
    Errors        []ParseError                      // Detailed error information
    Warnings      []ParseWarning                    // Non-critical warnings
    ColumnMapping map[string]int                    // Column name to index mapping
    Statistics    ParseStatistics                   // Parsing statistics
}
```

### Statistics Information
```go
type ParseStatistics struct {
    TablesFound       int                    // Number of HTML tables found
    HeadersDetected   []string               // Column headers detected
    DataTypesDetected map[string]string      // Detected data types per column
    ValueRanges       map[string]ValueRange  // Value ranges for numeric columns
    ProcessingTime    time.Duration          // Time taken to parse
}
```

## Integration with Database Layer

The parser outputs `models.CreateSalesRecordRequest` structures that are directly compatible with the database layer:

```go
// Parse HTML data
result, err := parser.ParseHTML(htmlData)
if err != nil {
    return err
}

// Import successful records into database
if result.SuccessCount > 0 {
    importResult, err := dbService.ImportSalesData(result.Records)
    if err != nil {
        return fmt.Errorf("database import failed: %w", err)
    }
    
    fmt.Printf("Imported %d records successfully\n", importResult.SuccessfulRecords)
}
```

## Performance Characteristics

### Benchmarks
- **Small Tables** (1-10 rows): ~100μs
- **Medium Tables** (10-100 rows): ~1ms
- **Large Tables** (100-1000 rows): ~10ms

### Memory Usage
- **Efficient Processing**: Streaming HTML parsing
- **Minimal Allocation**: Reuses data structures where possible
- **Garbage Collection Friendly**: Minimal temporary allocations

## Testing

### Comprehensive Test Suite
The parser includes extensive unit tests covering:

- **Basic HTML table parsing**
- **Various column name variations**
- **Multiple data formats (HTML, tab-delimited, pipe-delimited)**
- **Error handling and validation**
- **Currency and date parsing**
- **Real-world HTML examples**
- **Performance benchmarks**

### Running Tests
```bash
# Run all parser tests
go test ./internal/parser -v

# Run with coverage
go test ./internal/parser -cover

# Run benchmarks
go test ./internal/parser -bench=.
```

## Common Use Cases

### Website Data Import
1. User copies HTML table from a sales website
2. Pastes data into Sales Track application
3. Parser automatically detects format and extracts data
4. Validates and imports records into database

### Spreadsheet Migration
1. Export data from Excel/Google Sheets as HTML
2. Import HTML data using the parser
3. Automatic column mapping and data validation
4. Seamless migration to Sales Track database

### Data Validation
1. Parse data to identify formatting issues
2. Review errors and warnings before import
3. Clean data based on parser feedback
4. Re-import validated data

## Error Recovery

### Partial Import Strategy
- Parse all rows, collecting errors for problematic rows
- Import successful rows immediately
- Provide detailed feedback for failed rows
- Allow user to fix and re-import failed rows

### Data Cleaning Suggestions
- Suggest corrections for common formatting issues
- Provide examples of expected formats
- Highlight specific problematic values
- Offer automatic cleaning options where safe

## Future Enhancements

### Planned Features
- **Excel File Support**: Direct .xlsx file parsing
- **CSV Import**: Native CSV file support
- **Data Preview**: Preview parsed data before import
- **Custom Column Mapping**: User-defined column mappings
- **Batch Processing**: Handle multiple tables/files
- **Data Transformation**: Custom data transformation rules

### Integration Improvements
- **Wails Integration**: Direct integration with frontend
- **Real-time Validation**: Live parsing feedback
- **Progress Tracking**: Progress bars for large imports
- **Undo/Redo**: Rollback failed imports

## Troubleshooting

### Common Issues

1. **"No HTML tables found"**
   - Ensure data contains `<table>` tags
   - Try tab-delimited format if copying from spreadsheet

2. **"Missing required columns"**
   - Check column headers match expected names
   - Verify all required columns are present

3. **Date parsing errors**
   - Use consistent date format
   - Supported formats: YYYY-MM-DD, MM/DD/YYYY, etc.

4. **Currency parsing errors**
   - Remove extra formatting characters
   - Use standard currency symbols ($, €, £)

### Debug Mode
Enable detailed logging for troubleshooting:
```go
parser := parser.NewHTMLTableParser()
parser.StrictMode = true // Enable strict validation
```

This comprehensive HTML table parser provides the foundation for seamless data import in the Sales Track application, handling the complexity of various HTML formats while providing detailed feedback and robust error handling.

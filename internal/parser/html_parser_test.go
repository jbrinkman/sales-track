package parser

import (
	"strings"
	"testing"
	"time"
)

// Test data constants for better maintainability and readability
const (
	basicTableHTML = `
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
		<tr>
			<td>Mall Location</td>
			<td>Home & Garden</td>
			<td>01/16/2024</td>
			<td>Patio Set</td>
			<td>1299.00</td>
			<td>129.90</td>
			<td>1169.10</td>
		</tr>
	</table>
	`

	variousColumnNamesHTML = `
	<table>
		<tr>
			<th>Shop Name</th>
			<th>Supplier</th>
			<th>Sale Date</th>
			<th>Product</th>
			<th>Amount</th>
			<th>Fee</th>
			<th>Balance</th>
		</tr>
		<tr>
			<td>Test Store</td>
			<td>Test Vendor</td>
			<td>2024-02-01</td>
			<td>Test Product</td>
			<td>$100.00</td>
			<td>$10.00</td>
			<td>$90.00</td>
		</tr>
	</table>
	`

	headerlessRowsHTML = `
	<tr>
		<td>Downtown Store</td>
		<td>Electronics Plus</td>
		<td>2024-01-15</td>
		<td>Samsung TV</td>
		<td>$899.99</td>
		<td>$89.99</td>
		<td>$810.00</td>
	</tr>
	<tr>
		<td>Mall Location</td>
		<td>Home & Garden</td>
		<td>01/16/2024</td>
		<td>Patio Set</td>
		<td>1299.00</td>
		<td>129.90</td>
		<td>1169.10</td>
	</tr>
	`

	errorHandlingHTML = `
	<table>
		<tr>
			<th>Store</th>
			<th>Vendor</th>
			<th>Date</th>
			<th>Description</th>
			<th>Sale Price</th>
		</tr>
		<tr>
			<td></td>
			<td>Test Vendor</td>
			<td>invalid-date</td>
			<td></td>
			<td>not-a-number</td>
		</tr>
		<tr>
			<td>Valid Store</td>
			<td>Valid Vendor</td>
			<td>2024-01-15</td>
			<td>Valid Product</td>
			<td>100.00</td>
		</tr>
	</table>
	`

	realWorldConsignableHTML = `
	<tr class="odd">
		<td>Downtown Branch</td>
		<td>Tech Solutions Inc.</td>
		<td>March 15, 2024</td>
		<td>Laptop Computer - Dell XPS 13</td>
		<td>$1,299.99</td>
		<td>$129.99</td>
		<td>$1,170.00</td>
	</tr>
	<tr class="even">
		<td>Mall Outlet</td>
		<td>Home Essentials LLC</td>
		<td>03/16/2024</td>
		<td>Kitchen Appliance Set</td>
		<td>$899.50</td>
		<td>$89.95</td>
		<td>$809.55</td>
	</tr>
	<tr class="odd">
		<td>Westside Store</td>
		<td>Fashion Forward Co.</td>
		<td>2024-03-17</td>
		<td>Designer Handbag Collection</td>
		<td>$2,450.00</td>
		<td>$245.00</td>
		<td>$2,205.00</td>
	</tr>
	`

	tabDelimitedData = `Store	Vendor	Date	Description	Sale Price	Commission	Remaining
Downtown Store	Electronics Plus	2024-01-15	Samsung TV	899.99	89.99	810.00
Mall Location	Home & Garden	2024-01-16	Patio Set	1299.00	129.90	1169.10`
)

// TestParseHTML_BasicTable tests parsing a basic HTML table
func TestParseHTML_BasicTable(t *testing.T) {
	parser := NewHTMLTableParser()
	
	result, err := parser.ParseHTML(basicTableHTML)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	// Check basic statistics
	if result.TotalRows != 2 {
		t.Errorf("Expected 2 total rows, got %d", result.TotalRows)
	}
	
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful records, got %d", result.SuccessCount)
	}
	
	if result.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", result.ErrorCount)
	}
	
	if len(result.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result.Records))
	}
	
	// Check first record
	record1 := result.Records[0]
	if record1.Store != "Downtown Store" {
		t.Errorf("Expected store 'Downtown Store', got '%s'", record1.Store)
	}
	if record1.Vendor != "Electronics Plus" {
		t.Errorf("Expected vendor 'Electronics Plus', got '%s'", record1.Vendor)
	}
	if record1.Date != "2024-01-15" {
		t.Errorf("Expected date '2024-01-15', got '%s'", record1.Date)
	}
	if record1.Description != "Samsung TV" {
		t.Errorf("Expected description 'Samsung TV', got '%s'", record1.Description)
	}
	if record1.SalePrice != 899.99 {
		t.Errorf("Expected sale price 899.99, got %f", record1.SalePrice)
	}
	if record1.Commission != 89.99 {
		t.Errorf("Expected commission 89.99, got %f", record1.Commission)
	}
	if record1.Remaining != 810.00 {
		t.Errorf("Expected remaining 810.00, got %f", record1.Remaining)
	}
	
	// Check second record with different date format
	record2 := result.Records[1]
	if record2.Date != "2024-01-16" {
		t.Errorf("Expected date '2024-01-16', got '%s'", record2.Date)
	}
	if record2.SalePrice != 1299.00 {
		t.Errorf("Expected sale price 1299.00, got %f", record2.SalePrice)
	}
}

// TestParseHTML_VariousColumnNames tests parsing with different column name variations
func TestParseHTML_VariousColumnNames(t *testing.T) {
	parser := NewHTMLTableParser()
	
	result, err := parser.ParseHTML(variousColumnNamesHTML)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful record, got %d", result.SuccessCount)
	}
	
	record := result.Records[0]
	if record.Store != "Test Store" {
		t.Errorf("Expected store 'Test Store', got '%s'", record.Store)
	}
	if record.Vendor != "Test Vendor" {
		t.Errorf("Expected vendor 'Test Vendor', got '%s'", record.Vendor)
	}
}

// TestParseHTML_TabDelimited tests parsing tab-delimited data
func TestParseHTML_TabDelimited(t *testing.T) {
	parser := NewHTMLTableParser()
	
	result, err := parser.ParseHTML(tabDelimitedData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful records, got %d", result.SuccessCount)
	}
	
	if len(result.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result.Records))
	}
}

// TestParseHTML_ErrorHandling tests error handling for invalid data
func TestParseHTML_ErrorHandling(t *testing.T) {
	parser := NewHTMLTableParser()
	
	result, err := parser.ParseHTML(errorHandlingHTML)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.TotalRows != 2 {
		t.Errorf("Expected 2 total rows, got %d", result.TotalRows)
	}
	
	if result.ErrorCount != 1 {
		t.Errorf("Expected 1 error, got %d", result.ErrorCount)
	}
	
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 success, got %d", result.SuccessCount)
	}
	
	// Check that we have errors for the first row
	if len(result.Errors) == 0 {
		t.Error("Expected parsing errors, got none")
	}
	
	// Check that the valid record was parsed correctly
	if len(result.Records) != 1 {
		t.Errorf("Expected 1 valid record, got %d", len(result.Records))
	}
	
	if len(result.Records) > 0 {
		record := result.Records[0]
		if record.Store != "Valid Store" {
			t.Errorf("Expected store 'Valid Store', got '%s'", record.Store)
		}
	}
}

// TestParseHTML_MissingColumns tests handling of missing required columns
func TestParseHTML_MissingColumns(t *testing.T) {
	parser := NewHTMLTableParser()
	
	htmlData := `
	<table>
		<tr>
			<th>Store</th>
			<th>Date</th>
		</tr>
		<tr>
			<td>Test Store</td>
			<td>2024-01-15</td>
		</tr>
	</table>
	`
	
	_, err := parser.ParseHTML(htmlData)
	if err == nil {
		t.Error("Expected error for missing required columns, got none")
	}
	
	if !strings.Contains(err.Error(), "missing required columns") {
		t.Errorf("Expected 'missing required columns' error, got: %v", err)
	}
}

// TestParseHTML_NoTables tests handling when no tables are found
func TestParseHTML_NoTables(t *testing.T) {
	parser := NewHTMLTableParser()
	
	htmlData := `<div>This is not a table</div>`
	
	_, err := parser.ParseHTML(htmlData)
	if err == nil {
		t.Error("Expected error for no tables found, got none")
	}
	
	if !strings.Contains(err.Error(), "no HTML tables found") {
		t.Errorf("Expected 'no HTML tables found' error, got: %v", err)
	}
}

// TestParseHTML_MultipleTables tests handling multiple tables (should pick the largest)
func TestParseHTML_MultipleTables(t *testing.T) {
	parser := NewHTMLTableParser()
	
	htmlData := `
	<table>
		<tr><th>Small</th></tr>
		<tr><td>Table</td></tr>
	</table>
	<table>
		<tr>
			<th>Store</th>
			<th>Vendor</th>
			<th>Date</th>
			<th>Description</th>
			<th>Sale Price</th>
		</tr>
		<tr>
			<td>Store 1</td>
			<td>Vendor 1</td>
			<td>2024-01-15</td>
			<td>Product 1</td>
			<td>100.00</td>
		</tr>
		<tr>
			<td>Store 2</td>
			<td>Vendor 2</td>
			<td>2024-01-16</td>
			<td>Product 2</td>
			<td>200.00</td>
		</tr>
	</table>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.Statistics.TablesFound != 2 {
		t.Errorf("Expected 2 tables found, got %d", result.Statistics.TablesFound)
	}
	
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful records (from larger table), got %d", result.SuccessCount)
	}
}

// TestParseCurrency tests currency parsing with various formats
func TestParseCurrency(t *testing.T) {
	parser := NewHTMLTableParser()
	
	testCases := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"$100.00", 100.00, false},
		{"100.50", 100.50, false},
		{"1,234.56", 1234.56, false},
		{"$1,234.56", 1234.56, false},
		{"(50.00)", -50.00, false},
		{"", 0.00, false},
		{"not-a-number", 0.00, true},
		{"€123.45", 123.45, false},
		{"£99.99", 99.99, false},
	}
	
	for _, tc := range testCases {
		result, err := parser.parseCurrency(tc.input)
		
		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', got none", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("For input '%s', expected %f, got %f", tc.input, tc.expected, result)
			}
		}
	}
}

// TestParseDate tests date parsing with various formats
func TestParseDate(t *testing.T) {
	parser := NewHTMLTableParser()
	
	testCases := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"2024-01-15", "2024-01-15", false},
		{"01/15/2024", "2024-01-15", false},
		{"1/15/2024", "2024-01-15", false},
		{"Jan 15, 2024", "2024-01-15", false},
		{"January 15, 2024", "2024-01-15", false},
		{"15 Jan 2024", "2024-01-15", false},
		{"invalid-date", "", true},
		{"", "", true},
	}
	
	for _, tc := range testCases {
		result, err := parser.parseDate(tc.input)
		
		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', got none", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("For input '%s', expected '%s', got '%s'", tc.input, tc.expected, result)
			}
		}
	}
}

// TestDetectDataType tests data type detection
func TestDetectDataType(t *testing.T) {
	parser := NewHTMLTableParser()
	
	testCases := []struct {
		values   []string
		expected string
	}{
		{[]string{"2024-01-15", "2024-01-16", "2024-01-17"}, "date"},
		{[]string{"$100.00", "$200.50", "$300.75"}, "currency"},
		{[]string{"100", "200", "300"}, "number"},
		{[]string{"Store A", "Store B", "Store C"}, "text"},
		{[]string{}, "unknown"},
	}
	
	for _, tc := range testCases {
		result := parser.detectDataType(tc.values)
		if result != tc.expected {
			t.Errorf("For values %v, expected type '%s', got '%s'", tc.values, tc.expected, result)
		}
	}
}

// TestLooksLikeDate tests date pattern recognition
func TestLooksLikeDate(t *testing.T) {
	parser := NewHTMLTableParser()
	
	testCases := []struct {
		input    string
		expected bool
	}{
		{"2024-01-15", true},
		{"01/15/2024", true},
		{"Jan 15, 2024", true},
		{"January 15, 2024", true},
		{"not a date", false},
		{"123.45", false},
		{"", false},
	}
	
	for _, tc := range testCases {
		result := parser.looksLikeDate(tc.input)
		if result != tc.expected {
			t.Errorf("For input '%s', expected %t, got %t", tc.input, tc.expected, result)
		}
	}
}

// TestLooksLikeCurrency tests currency pattern recognition
func TestLooksLikeCurrency(t *testing.T) {
	parser := NewHTMLTableParser()
	
	testCases := []struct {
		input    string
		expected bool
	}{
		{"$100.00", true},
		{"123.45", true},
		{"(50.00)", true},
		{"not currency", false},
		{"", false},
	}
	
	for _, tc := range testCases {
		result := parser.looksLikeCurrency(tc.input)
		if result != tc.expected {
			t.Errorf("For input '%s', expected %t, got %t", tc.input, tc.expected, result)
		}
	}
}

// TestParseHTML_RealWorldExample tests with a more realistic HTML table
func TestParseHTML_RealWorldExample(t *testing.T) {
	parser := NewHTMLTableParser()
	
	// Simulate HTML that might be copied from a website
	htmlData := `
	<table class="data-table" border="1">
		<thead>
			<tr style="background-color: #f0f0f0;">
				<th>Store Location</th>
				<th>Vendor Name</th>
				<th>Transaction Date</th>
				<th>Item Description</th>
				<th>Sale Amount</th>
				<th>Commission Fee</th>
				<th>Remaining Balance</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>Downtown Branch</td>
				<td>Tech Solutions Inc.</td>
				<td>March 15, 2024</td>
				<td>Laptop Computer - Dell XPS 13</td>
				<td>$1,299.99</td>
				<td>$129.99</td>
				<td>$1,170.00</td>
			</tr>
			<tr>
				<td>Mall Outlet</td>
				<td>Home Essentials LLC</td>
				<td>03/16/2024</td>
				<td>Kitchen Appliance Set</td>
				<td>$899.50</td>
				<td>$89.95</td>
				<td>$809.55</td>
			</tr>
			<tr>
				<td>Westside Store</td>
				<td>Fashion Forward Co.</td>
				<td>2024-03-17</td>
				<td>Designer Handbag Collection</td>
				<td>$2,450.00</td>
				<td>$245.00</td>
				<td>$2,205.00</td>
			</tr>
		</tbody>
	</table>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.SuccessCount != 3 {
		t.Errorf("Expected 3 successful records, got %d", result.SuccessCount)
	}
	
	if result.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d. Errors: %v", result.ErrorCount, result.Errors)
	}
	
	// Check that different date formats were parsed correctly
	expectedDates := []string{"2024-03-15", "2024-03-16", "2024-03-17"}
	for i, record := range result.Records {
		if record.Date != expectedDates[i] {
			t.Errorf("Record %d: expected date '%s', got '%s'", i, expectedDates[i], record.Date)
		}
	}
	
	// Check that currency values with commas were parsed correctly
	if result.Records[0].SalePrice != 1299.99 {
		t.Errorf("Expected first record sale price 1299.99, got %f", result.Records[0].SalePrice)
	}
	
	if result.Records[2].SalePrice != 2450.00 {
		t.Errorf("Expected third record sale price 2450.00, got %f", result.Records[2].SalePrice)
	}
}

// BenchmarkParseHTML benchmarks the HTML parsing performance
func BenchmarkParseHTML(b *testing.B) {
	parser := NewHTMLTableParser()
	
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
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseHTML(htmlData)
		if err != nil {
			b.Fatalf("ParseHTML failed: %v", err)
		}
	}
}

// TestParseHTML_HeaderlessRows tests parsing table rows without headers using positional mapping
func TestParseHTML_HeaderlessRows(t *testing.T) {
	parser := NewHTMLTableParser()
	parser.SetConsignableMapping() // Use standard Consignable format
	
	result, err := parser.ParseHTML(headerlessRowsHTML)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	// Check basic statistics
	if result.TotalRows != 2 {
		t.Errorf("Expected 2 total rows, got %d", result.TotalRows)
	}
	
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful records, got %d", result.SuccessCount)
	}
	
	if result.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d. Errors: %v", result.ErrorCount, result.Errors)
	}
	
	if len(result.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result.Records))
	}
	
	// Check first record
	record1 := result.Records[0]
	if record1.Store != "Downtown Store" {
		t.Errorf("Expected store 'Downtown Store', got '%s'", record1.Store)
	}
	if record1.Vendor != "Electronics Plus" {
		t.Errorf("Expected vendor 'Electronics Plus', got '%s'", record1.Vendor)
	}
	if record1.Date != "2024-01-15" {
		t.Errorf("Expected date '2024-01-15', got '%s'", record1.Date)
	}
	if record1.Description != "Samsung TV" {
		t.Errorf("Expected description 'Samsung TV', got '%s'", record1.Description)
	}
	if record1.SalePrice != 899.99 {
		t.Errorf("Expected sale price 899.99, got %f", record1.SalePrice)
	}
	if record1.Commission != 89.99 {
		t.Errorf("Expected commission 89.99, got %f", record1.Commission)
	}
	if record1.Remaining != 810.00 {
		t.Errorf("Expected remaining 810.00, got %f", record1.Remaining)
	}
	
	// Check second record with different date format
	record2 := result.Records[1]
	if record2.Date != "2024-01-16" {
		t.Errorf("Expected date '2024-01-16', got '%s'", record2.Date)
	}
	if record2.SalePrice != 1299.00 {
		t.Errorf("Expected sale price 1299.00, got %f", record2.SalePrice)
	}
}

// TestParseHTML_CustomPositionalMapping tests custom positional column mapping
func TestParseHTML_CustomPositionalMapping(t *testing.T) {
	parser := NewHTMLTableParser()
	
	// Custom column order: Vendor, Store, Date, Description, Sale Price
	parser.SetPositionalMapping([]string{
		"vendor",
		"store", 
		"date",
		"description",
		"sale_price",
	})
	
	htmlData := `
	<tr>
		<td>Tech Vendor</td>
		<td>Test Store</td>
		<td>2024-02-01</td>
		<td>Test Product</td>
		<td>$100.00</td>
	</tr>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful record, got %d", result.SuccessCount)
	}
	
	record := result.Records[0]
	if record.Vendor != "Tech Vendor" {
		t.Errorf("Expected vendor 'Tech Vendor', got '%s'", record.Vendor)
	}
	if record.Store != "Test Store" {
		t.Errorf("Expected store 'Test Store', got '%s'", record.Store)
	}
}

// TestParseHTML_HeaderlessWithTbody tests parsing tbody sections without table wrapper
func TestParseHTML_HeaderlessWithTbody(t *testing.T) {
	parser := NewHTMLTableParser()
	parser.SetConsignableMapping()
	
	htmlData := `
	<tbody>
		<tr>
			<td>Store A</td>
			<td>Vendor 1</td>
			<td>2024-01-15</td>
			<td>Product A</td>
			<td>$150.00</td>
			<td>$15.00</td>
			<td>$135.00</td>
		</tr>
		<tr>
			<td>Store B</td>
			<td>Vendor 2</td>
			<td>2024-01-16</td>
			<td>Product B</td>
			<td>$250.00</td>
			<td>$25.00</td>
			<td>$225.00</td>
		</tr>
	</tbody>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful records, got %d", result.SuccessCount)
	}
	
	if len(result.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result.Records))
	}
}

// TestParseHTML_HeaderlessInsufficientColumns tests error handling for insufficient columns
func TestParseHTML_HeaderlessInsufficientColumns(t *testing.T) {
	parser := NewHTMLTableParser()
	parser.SetConsignableMapping() // Expects 7 columns
	
	// Only 3 columns provided
	htmlData := `
	<tr>
		<td>Store A</td>
		<td>Vendor 1</td>
		<td>2024-01-15</td>
	</tr>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	// Should have 1 total row but 0 successful records due to missing required fields
	if result.TotalRows != 1 {
		t.Errorf("Expected 1 total row, got %d", result.TotalRows)
	}
	
	if result.SuccessCount != 0 {
		t.Errorf("Expected 0 successful records due to missing fields, got %d", result.SuccessCount)
	}
	
	if result.ErrorCount != 1 {
		t.Errorf("Expected 1 error record, got %d", result.ErrorCount)
	}
	
	// Should have errors for missing required fields (description, sale_price)
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing required fields")
	}
	
	// Check that we have errors for the missing required fields
	foundDescriptionError := false
	foundSalePriceError := false
	
	for _, parseErr := range result.Errors {
		if parseErr.Column == "description" && strings.Contains(parseErr.Message, "required but empty") {
			foundDescriptionError = true
		}
		if parseErr.Column == "sale_price" && strings.Contains(parseErr.Message, "required but empty") {
			foundSalePriceError = true
		}
	}
	
	if !foundDescriptionError {
		t.Error("Expected error for missing description field")
	}
	if !foundSalePriceError {
		t.Error("Expected error for missing sale_price field")
	}
}

// TestParseHTML_HeaderlessWithErrors tests error handling in headerless mode
func TestParseHTML_HeaderlessWithErrors(t *testing.T) {
	parser := NewHTMLTableParser()
	parser.SetConsignableMapping()
	
	htmlData := `
	<tr>
		<td></td>
		<td>Valid Vendor</td>
		<td>invalid-date</td>
		<td>Valid Product</td>
		<td>not-a-price</td>
		<td>10.00</td>
		<td>90.00</td>
	</tr>
	<tr>
		<td>Valid Store</td>
		<td>Valid Vendor</td>
		<td>2024-01-15</td>
		<td>Valid Product</td>
		<td>$199.99</td>
		<td>$19.99</td>
		<td>$180.00</td>
	</tr>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.TotalRows != 2 {
		t.Errorf("Expected 2 total rows, got %d", result.TotalRows)
	}
	
	if result.ErrorCount != 1 {
		t.Errorf("Expected 1 error, got %d", result.ErrorCount)
	}
	
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 success, got %d", result.SuccessCount)
	}
	
	// Check that we have errors for the first row
	if len(result.Errors) == 0 {
		t.Error("Expected parsing errors, got none")
	}
	
	// Check that the valid record was parsed correctly
	if len(result.Records) != 1 {
		t.Errorf("Expected 1 valid record, got %d", len(result.Records))
	}
	
	if len(result.Records) > 0 {
		record := result.Records[0]
		if record.Store != "Valid Store" {
			t.Errorf("Expected store 'Valid Store', got '%s'", record.Store)
		}
		if record.SalePrice != 199.99 {
			t.Errorf("Expected sale price 199.99, got %f", record.SalePrice)
		}
	}
}

// TestParseHTML_RealWorldConsignableExample tests with realistic Consignable HTML
func TestParseHTML_RealWorldConsignableExample(t *testing.T) {
	parser := NewHTMLTableParser()
	parser.SetConsignableMapping()
	
	result, err := parser.ParseHTML(realWorldConsignableHTML)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.SuccessCount != 3 {
		t.Errorf("Expected 3 successful records, got %d", result.SuccessCount)
	}
	
	if result.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d. Errors: %v", result.ErrorCount, result.Errors)
	}
	
	// Check that different date formats were parsed correctly
	expectedDates := []string{"2024-03-15", "2024-03-16", "2024-03-17"}
	for i, record := range result.Records {
		if record.Date != expectedDates[i] {
			t.Errorf("Record %d: expected date '%s', got '%s'", i, expectedDates[i], record.Date)
		}
	}
	
	// Check that currency values with commas were parsed correctly
	if result.Records[0].SalePrice != 1299.99 {
		t.Errorf("Expected first record sale price 1299.99, got %f", result.Records[0].SalePrice)
	}
	
	if result.Records[2].SalePrice != 2450.00 {
		t.Errorf("Expected third record sale price 2450.00, got %f", result.Records[2].SalePrice)
	}
}

// TestLooksLikeTableRows tests the table row detection logic
func TestLooksLikeTableRows(t *testing.T) {
	parser := NewHTMLTableParser()
	
	testCases := []struct {
		input    string
		expected bool
	}{
		{"<tr><td>data</td></tr>", true},
		{"<table><tr><td>data</td></tr></table>", false},
		{"<TR><TD>data</TD></TR>", true}, // Case insensitive
		{"<div>not table rows</div>", false},
		{"<tbody><tr><td>data</td></tr></tbody>", true},
		{"", false},
	}
	
	for _, tc := range testCases {
		result := parser.looksLikeTableRows(tc.input)
		if result != tc.expected {
			t.Errorf("For input '%s', expected %t, got %t", tc.input, tc.expected, result)
		}
	}
}
// TestParseHTML_ProcessingTime tests that processing time is recorded
func TestParseHTML_ProcessingTime(t *testing.T) {
	parser := NewHTMLTableParser()
	
	htmlData := `
	<table>
		<tr><th>Store</th><th>Vendor</th><th>Date</th><th>Description</th><th>Sale Price</th></tr>
		<tr><td>Test Store</td><td>Test Vendor</td><td>2024-01-15</td><td>Test Product</td><td>100.00</td></tr>
	</table>
	`
	
	result, err := parser.ParseHTML(htmlData)
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	
	if result.Statistics.ProcessingTime <= 0 {
		t.Error("Expected processing time to be recorded and greater than 0")
	}
	
	if result.Statistics.ProcessingTime > time.Second {
		t.Error("Processing time seems unusually high for a simple table")
	}
}

package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"sales-track/internal/models"
)

// Package-level compiled regex patterns for performance optimization
// These are compiled once at package initialization instead of on every function call
var (
	datePatterns = []*regexp.Regexp{
		regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}`),
		regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{4}`),
		regexp.MustCompile(`^\d{1,2}-\d{1,2}-\d{4}`),
		regexp.MustCompile(`^[A-Za-z]{3,9}\s+\d{1,2},?\s+\d{4}`),
	}
	
	currencyPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^\$\d+\.?\d*`),
		regexp.MustCompile(`^\d+\.\d{2}$`),
		regexp.MustCompile(`^\(\d+\.?\d*\)$`), // Negative in parentheses
	}
)

// HTMLTableParser handles parsing HTML table data into sales records
type HTMLTableParser struct {
	// Configuration options
	StrictMode bool // If true, requires exact column matches
	
	// Positional mapping for headerless tables
	UsePositionalMapping bool     // Enable positional column mapping
	PositionalColumns    []string // Column names in order for positional mapping
}

// NewHTMLTableParser creates a new HTML table parser
func NewHTMLTableParser() *HTMLTableParser {
	return &HTMLTableParser{
		StrictMode: false,
	}
}

// SetPositionalMapping configures the parser to use positional column mapping
// for headerless tables. Columns should be in the order they appear in the HTML.
func (p *HTMLTableParser) SetPositionalMapping(columns []string) {
	p.UsePositionalMapping = true
	p.PositionalColumns = columns
}

// SetConsignableMapping configures the parser for the standard Consignable format:
// Store, Vendor, Date, Description, Sale Price, Commission, Remaining
func (p *HTMLTableParser) SetConsignableMapping() {
	p.SetPositionalMapping([]string{
		"store",
		"vendor", 
		"date",
		"description",
		"sale_price",
		"commission",
		"remaining",
	})
}

// ParseResult contains the results of parsing HTML table data
type ParseResult struct {
	Records       []models.CreateSalesRecordRequest `json:"records"`
	TotalRows     int                               `json:"total_rows"`
	SuccessCount  int                               `json:"success_count"`
	ErrorCount    int                               `json:"error_count"`
	Errors        []ParseError                      `json:"errors,omitempty"`
	Warnings      []ParseWarning                    `json:"warnings,omitempty"`
	ColumnMapping map[string]int                    `json:"column_mapping"`
	Statistics    ParseStatistics                   `json:"statistics"`
}

// ParseError represents an error that occurred during parsing
type ParseError struct {
	Row     int    `json:"row"`
	Column  string `json:"column,omitempty"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ParseWarning represents a warning that occurred during parsing
type ParseWarning struct {
	Row     int    `json:"row"`
	Column  string `json:"column,omitempty"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ParseStatistics contains statistics about the parsing operation
type ParseStatistics struct {
	TablesFound       int                    `json:"tables_found"`
	HeadersDetected   []string               `json:"headers_detected"`
	DataTypesDetected map[string]string      `json:"data_types_detected"`
	ValueRanges       map[string]ValueRange  `json:"value_ranges,omitempty"`
	ProcessingTime    time.Duration          `json:"processing_time"`
}

// ValueRange represents the range of values found in a column
type ValueRange struct {
	Min   interface{} `json:"min,omitempty"`
	Max   interface{} `json:"max,omitempty"`
	Count int         `json:"count"`
}

// Required columns for sales record validation
var requiredColumns = []string{"store", "vendor", "date", "description", "sale_price"}

// validateRequiredColumns consolidates validation logic for required columns
func (p *HTMLTableParser) validateRequiredColumns(mapping map[string]int, context string) error {
	missingColumns := []string{}
	
	for _, col := range requiredColumns {
		if _, exists := mapping[col]; !exists {
			missingColumns = append(missingColumns, col)
		}
	}
	
	if len(missingColumns) > 0 {
		return fmt.Errorf("%s missing required columns: %v", context, missingColumns)
	}
	
	return nil
}
var ColumnMapping = map[string][]string{
	"store": {
		"store", "shop", "location", "outlet", "branch", "store name", "shop name",
	},
	"vendor": {
		"vendor", "supplier", "brand", "manufacturer", "company", "vendor name", "supplier name",
	},
	"date": {
		"date", "sale date", "transaction date", "order date", "purchase date", "sold date",
	},
	"description": {
		"description", "item", "product", "item description", "product description", "details", "name", "product name",
	},
	"sale_price": {
		"sale price", "price", "amount", "total", "sale amount", "selling price", "cost", "value",
	},
	"commission": {
		"commission", "fee", "commission amount", "commission fee", "comm", "commission %", "commission rate",
	},
	"remaining": {
		"remaining", "balance", "remaining balance", "outstanding", "due", "remaining amount", "balance due",
	},
}

// ParseHTML parses HTML table data and extracts sales records
func (p *HTMLTableParser) ParseHTML(htmlData string) (*ParseResult, error) {
	startTime := time.Now()
	
	result := &ParseResult{
		Records:       []models.CreateSalesRecordRequest{},
		ColumnMapping: make(map[string]int),
		Statistics: ParseStatistics{
			DataTypesDetected: make(map[string]string),
			ValueRanges:       make(map[string]ValueRange),
		},
	}

	// Clean and prepare HTML data
	cleanHTML := p.cleanHTML(htmlData)
	
	// Parse HTML
	doc, err := html.Parse(strings.NewReader(cleanHTML))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find all tables
	tables := p.findTables(doc)
	result.Statistics.TablesFound = len(tables)
	
	if len(tables) == 0 {
		return nil, fmt.Errorf("no HTML tables found in the provided data")
	}

	// Process the first table (or the largest table if multiple)
	table := p.selectBestTable(tables)
	
	// Extract table data
	tableData, err := p.extractTableData(table)
	if err != nil {
		return nil, fmt.Errorf("failed to extract table data: %w", err)
	}

	if len(tableData) == 0 {
		return nil, fmt.Errorf("no data rows found in table")
	}

	result.TotalRows = len(tableData) - 1 // Subtract header row

	// Detect headers and create column mapping
	headers := tableData[0]
	result.Statistics.HeadersDetected = headers
	
	columnMapping, err := p.createColumnMapping(headers)
	if err != nil {
		return nil, fmt.Errorf("failed to map columns: %w", err)
	}
	result.ColumnMapping = columnMapping

	// Parse data rows
	for i, row := range tableData[1:] {
		rowNum := i + 2 // +2 because we skip header and want 1-based indexing
		
		record, parseErrors, warnings := p.parseRow(row, columnMapping, rowNum)
		
		if len(parseErrors) > 0 {
			result.Errors = append(result.Errors, parseErrors...)
			result.ErrorCount++
		} else {
			result.Records = append(result.Records, record)
			result.SuccessCount++
		}
		
		if len(warnings) > 0 {
			result.Warnings = append(result.Warnings, warnings...)
		}
	}

	// Calculate statistics
	p.calculateStatistics(result, tableData)
	result.Statistics.ProcessingTime = time.Since(startTime)

	return result, nil
}

// cleanHTML cleans and normalizes HTML data
func (p *HTMLTableParser) cleanHTML(htmlData string) string {
	// Remove common problematic characters and normalize whitespace
	cleaned := strings.TrimSpace(htmlData)
	
	// Check if this looks like table rows without a table wrapper
	if p.looksLikeTableRows(cleaned) {
		return p.wrapTableRows(cleaned)
	}
	
	// If it doesn't look like HTML, wrap it in a basic table structure
	if !strings.Contains(strings.ToLower(cleaned), "<table") {
		// Try to detect if it's tab-separated or other delimited data
		if strings.Contains(cleaned, "\t") || strings.Contains(cleaned, "|") {
			return p.convertDelimitedToHTML(cleaned)
		}
	}
	
	// Ensure we have a complete HTML document structure
	if !strings.Contains(strings.ToLower(cleaned), "<html") {
		cleaned = fmt.Sprintf("<html><body>%s</body></html>", cleaned)
	}
	
	return cleaned
}

// looksLikeTableRows checks if the HTML looks like table rows without table wrapper
func (p *HTMLTableParser) looksLikeTableRows(htmlData string) bool {
	lower := strings.ToLower(htmlData)
	
	// Check if it contains <tr> tags but no <table> tag
	hasTR := strings.Contains(lower, "<tr")
	hasTable := strings.Contains(lower, "<table")
	
	return hasTR && !hasTable
}

// wrapTableRows wraps table rows in a proper table structure
func (p *HTMLTableParser) wrapTableRows(rowsHTML string) string {
	var htmlBuilder strings.Builder
	
	htmlBuilder.WriteString("<html><body><table>")
	
	// If using positional mapping, add synthetic headers
	if p.UsePositionalMapping && len(p.PositionalColumns) > 0 {
		htmlBuilder.WriteString("<thead><tr>")
		for _, col := range p.PositionalColumns {
			// Convert internal column names to display names
			displayName := p.getDisplayColumnName(col)
			htmlBuilder.WriteString(fmt.Sprintf("<th>%s</th>", displayName))
		}
		htmlBuilder.WriteString("</tr></thead>")
	}
	
	htmlBuilder.WriteString("<tbody>")
	htmlBuilder.WriteString(rowsHTML)
	htmlBuilder.WriteString("</tbody></table></body></html>")
	
	return htmlBuilder.String()
}

// getDisplayColumnName converts internal column names to display names
func (p *HTMLTableParser) getDisplayColumnName(internalName string) string {
	displayNames := map[string]string{
		"store":       "Store",
		"vendor":      "Vendor", 
		"date":        "Date",
		"description": "Description",
		"sale_price":  "Sale Price",
		"commission":  "Commission",
		"remaining":   "Remaining",
	}
	
	if display, exists := displayNames[internalName]; exists {
		return display
	}
	return internalName
}

// convertDelimitedToHTML converts tab-separated or pipe-separated data to HTML table
func (p *HTMLTableParser) convertDelimitedToHTML(data string) string {
	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return data
	}

	var delimiter string
	if strings.Contains(lines[0], "\t") {
		delimiter = "\t"
	} else if strings.Contains(lines[0], "|") {
		delimiter = "|"
	} else {
		return data // Can't detect delimiter, return as-is
	}

	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<html><body><table>")
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		cells := strings.Split(line, delimiter)
		if i == 0 {
			htmlBuilder.WriteString("<thead><tr>")
			for _, cell := range cells {
				htmlBuilder.WriteString(fmt.Sprintf("<th>%s</th>", strings.TrimSpace(cell)))
			}
			htmlBuilder.WriteString("</tr></thead><tbody>")
		} else {
			htmlBuilder.WriteString("<tr>")
			for _, cell := range cells {
				htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td>", strings.TrimSpace(cell)))
			}
			htmlBuilder.WriteString("</tr>")
		}
	}
	
	htmlBuilder.WriteString("</tbody></table></body></html>")
	return htmlBuilder.String()
}

// findTables finds all table elements in the HTML document
func (p *HTMLTableParser) findTables(n *html.Node) []*html.Node {
	var tables []*html.Node
	
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {
			tables = append(tables, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	
	traverse(n)
	return tables
}

// selectBestTable selects the best table to parse (largest by row count)
func (p *HTMLTableParser) selectBestTable(tables []*html.Node) *html.Node {
	if len(tables) == 1 {
		return tables[0]
	}
	
	bestTable := tables[0]
	maxRows := p.countTableRows(bestTable)
	
	for _, table := range tables[1:] {
		rowCount := p.countTableRows(table)
		if rowCount > maxRows {
			bestTable = table
			maxRows = rowCount
		}
	}
	
	return bestTable
}

// countTableRows counts the number of rows in a table
func (p *HTMLTableParser) countTableRows(table *html.Node) int {
	count := 0
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			count++
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	traverse(table)
	return count
}

// extractTableData extracts all cell data from a table
func (p *HTMLTableParser) extractTableData(table *html.Node) ([][]string, error) {
	var rows [][]string
	
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			row := p.extractRowData(node)
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	
	traverse(table)
	return rows, nil
}

// extractRowData extracts cell data from a table row
func (p *HTMLTableParser) extractRowData(row *html.Node) []string {
	var cells []string
	
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && (node.Data == "td" || node.Data == "th") {
			cellText := p.extractTextContent(node)
			cells = append(cells, strings.TrimSpace(cellText))
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	
	traverse(row)
	return cells
}

// extractTextContent extracts text content from an HTML node
func (p *HTMLTableParser) extractTextContent(node *html.Node) string {
	var text strings.Builder
	
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	
	traverse(node)
	return text.String()
}

// createColumnMapping creates a mapping from expected columns to actual column indices
func (p *HTMLTableParser) createColumnMapping(headers []string) (map[string]int, error) {
	mapping := make(map[string]int)
	
	// If using positional mapping, create mapping based on position
	if p.UsePositionalMapping && len(p.PositionalColumns) > 0 {
		// Check if we have enough columns
		if len(headers) < len(p.PositionalColumns) {
			return nil, fmt.Errorf("positional mapping expects %d columns, but only %d headers found", 
				len(p.PositionalColumns), len(headers))
		}
		
		for i, col := range p.PositionalColumns {
			if i < len(headers) {
				mapping[col] = i
			}
		}
		
		// Use consolidated validation
		if err := p.validateRequiredColumns(mapping, "positional mapping"); err != nil {
			return nil, fmt.Errorf("%w. Expected %d columns, got %d headers", 
				err, len(p.PositionalColumns), len(headers))
		}
		
		return mapping, nil
	}
	
	// Original header-based mapping logic
	// Normalize headers for comparison
	normalizedHeaders := make([]string, len(headers))
	for i, header := range headers {
		normalizedHeaders[i] = strings.ToLower(strings.TrimSpace(header))
	}
	
	// Try to match each expected column
	for expectedCol, variations := range ColumnMapping {
		found := false
		for _, variation := range variations {
			for i, header := range normalizedHeaders {
				if strings.Contains(header, strings.ToLower(variation)) || 
				   strings.Contains(strings.ToLower(variation), header) {
					mapping[expectedCol] = i
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		
		if !found && p.StrictMode {
			return nil, fmt.Errorf("required column '%s' not found in headers: %v", expectedCol, headers)
		}
	}
	
	// Use consolidated validation
	if err := p.validateRequiredColumns(mapping, "header-based mapping"); err != nil {
		return nil, fmt.Errorf("%w. Available headers: %v", err, headers)
	}
	
	return mapping, nil
}

// parseRow parses a single data row into a sales record
func (p *HTMLTableParser) parseRow(row []string, columnMapping map[string]int, rowNum int) (models.CreateSalesRecordRequest, []ParseError, []ParseWarning) {
	var record models.CreateSalesRecordRequest
	var errors []ParseError
	var warnings []ParseWarning
	
	// Helper function to get cell value safely
	getCell := func(column string) string {
		if idx, exists := columnMapping[column]; exists && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}
	
	// Parse Store
	record.Store = getCell("store")
	if record.Store == "" {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  "store",
			Message: "Store field is required but empty",
		})
	}
	
	// Parse Vendor
	record.Vendor = getCell("vendor")
	if record.Vendor == "" {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  "vendor",
			Message: "Vendor field is required but empty",
		})
	}
	
	// Parse Date
	dateStr := getCell("date")
	if dateStr == "" {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  "date",
			Message: "Date field is required but empty",
		})
	} else {
		parsedDate, err := p.parseDate(dateStr)
		if err != nil {
			errors = append(errors, ParseError{
				Row:     rowNum,
				Column:  "date",
				Message: fmt.Sprintf("Invalid date format: %v", err),
				Value:   dateStr,
			})
		} else {
			record.Date = parsedDate
		}
	}
	
	// Parse Description
	record.Description = getCell("description")
	if record.Description == "" {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  "description",
			Message: "Description field is required but empty",
		})
	}
	
	// Parse Sale Price
	salePriceStr := getCell("sale_price")
	if salePriceStr == "" {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  "sale_price",
			Message: "Sale price field is required but empty",
		})
	} else {
		price, err := p.parseCurrency(salePriceStr)
		if err != nil {
			errors = append(errors, ParseError{
				Row:     rowNum,
				Column:  "sale_price",
				Message: fmt.Sprintf("Invalid sale price format: %v", err),
				Value:   salePriceStr,
			})
		} else {
			record.SalePrice = price
		}
	}
	
	// Parse Commission (optional)
	commissionStr := getCell("commission")
	if commissionStr != "" {
		commission, err := p.parseCurrency(commissionStr)
		if err != nil {
			warnings = append(warnings, ParseWarning{
				Row:     rowNum,
				Column:  "commission",
				Message: fmt.Sprintf("Invalid commission format, using 0.00: %v", err),
				Value:   commissionStr,
			})
			record.Commission = 0.00
		} else {
			record.Commission = commission
		}
	}
	
	// Parse Remaining (optional)
	remainingStr := getCell("remaining")
	if remainingStr != "" {
		remaining, err := p.parseCurrency(remainingStr)
		if err != nil {
			warnings = append(warnings, ParseWarning{
				Row:     rowNum,
				Column:  "remaining",
				Message: fmt.Sprintf("Invalid remaining format, using 0.00: %v", err),
				Value:   remainingStr,
			})
			record.Remaining = 0.00
		} else {
			record.Remaining = remaining
		}
	}
	
	return record, errors, warnings
}

// parseDate parses various date formats
func (p *HTMLTableParser) parseDate(dateStr string) (string, error) {
	// Common date formats to try
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"1/2/2006",
		"02/01/2006",
		"2/1/2006",
		"2006/01/02",
		"2006/1/2",
		"Jan 2, 2006",
		"January 2, 2006",
		"2 Jan 2006",
		"2 January 2006",
		"2006-01-02 15:04:05",
		"01/02/2006 15:04:05",
	}
	
	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed.Format("2006-01-02"), nil
		}
	}
	
	return "", fmt.Errorf("unable to parse date: %s", dateStr)
}

// parseCurrency parses currency values, handling various formats
func (p *HTMLTableParser) parseCurrency(currencyStr string) (float64, error) {
	// Remove common currency symbols and formatting
	cleaned := strings.TrimSpace(currencyStr)
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")
	cleaned = strings.ReplaceAll(cleaned, "¥", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	
	// Handle parentheses for negative numbers
	if strings.HasPrefix(cleaned, "(") && strings.HasSuffix(cleaned, ")") {
		cleaned = "-" + strings.Trim(cleaned, "()")
	}
	
	if cleaned == "" {
		return 0.0, nil
	}
	
	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0, fmt.Errorf("invalid currency format: %s", currencyStr)
	}
	
	return value, nil
}

// calculateStatistics calculates parsing statistics
func (p *HTMLTableParser) calculateStatistics(result *ParseResult, tableData [][]string) {
	if len(tableData) < 2 {
		return
	}
	
	headers := tableData[0]
	
	// Analyze data types for each column
	for i, header := range headers {
		if i >= len(tableData[1]) {
			continue
		}
		
		// Sample a few values to determine data type
		sampleValues := []string{}
		for j := 1; j < len(tableData) && j < 6; j++ { // Sample first 5 data rows
			if i < len(tableData[j]) {
				sampleValues = append(sampleValues, tableData[j][i])
			}
		}
		
		dataType := p.detectDataType(sampleValues)
		result.Statistics.DataTypesDetected[header] = dataType
	}
}

// detectDataType attempts to detect the data type of a column based on sample values
func (p *HTMLTableParser) detectDataType(values []string) string {
	if len(values) == 0 {
		return "unknown"
	}
	
	dateCount := 0
	currencyCount := 0
	numberCount := 0
	
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		
		// Check if it looks like a date
		if p.looksLikeDate(value) {
			dateCount++
		}
		
		// Check if it looks like currency
		if p.looksLikeCurrency(value) {
			currencyCount++
		}
		
		// Check if it's a number
		if _, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", ""), 64); err == nil {
			numberCount++
		}
	}
	
	total := len(values)
	if dateCount > total/2 {
		return "date"
	}
	if currencyCount > total/2 {
		return "currency"
	}
	if numberCount > total/2 {
		return "number"
	}
	
	return "text"
}

// looksLikeDate checks if a string looks like a date using pre-compiled patterns
func (p *HTMLTableParser) looksLikeDate(value string) bool {
	// Use pre-compiled regex patterns for better performance
	for _, pattern := range datePatterns {
		if pattern.MatchString(value) {
			return true
		}
	}
	
	return false
}

// looksLikeCurrency checks if a string looks like a currency value using pre-compiled patterns
func (p *HTMLTableParser) looksLikeCurrency(value string) bool {
	// Use pre-compiled regex patterns for better performance
	for _, pattern := range currencyPatterns {
		if pattern.MatchString(value) {
			return true
		}
	}
	
	return false
}

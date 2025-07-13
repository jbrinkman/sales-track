package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"sales-track/internal/parser"
)

// RunConsignableDemo demonstrates the Consignable HTML parser workflow
func RunConsignableDemo() {
	fmt.Println("=== Consignable HTML Parser Demo ===\n")
	
	// Create parser configured for Consignable format
	htmlParser := parser.NewHTMLTableParser()
	htmlParser.SetConsignableMapping() // Store, Vendor, Date, Description, Sale Price, Commission, Remaining
	
	// Get the directory of this example file
	exampleDir, _ := os.Getwd()
	exampleDir = filepath.Join(exampleDir, "internal", "parser", "examples")
	
	// Demo 1: Parse Consignable row fragment
	fmt.Println("1. Parsing Consignable HTML rows (headerless)...")
	consignableFile := filepath.Join(exampleDir, "consignable_rows.html")
	if htmlData, err := ioutil.ReadFile(consignableFile); err == nil {
		result, err := htmlParser.ParseHTML(string(htmlData))
		if err != nil {
			fmt.Printf("Consignable parsing failed: %v\n", err)
		} else {
			printConsignableResults("Consignable Rows", result)
		}
	} else {
		fmt.Printf("Could not read Consignable file: %v\n", err)
	}
	
	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	
	// Demo 2: Parse inline Consignable-style data
	fmt.Println("2. Parsing inline Consignable data...")
	inlineHTML := `
	<tr>
		<td>My Store</td>
		<td>My Vendor</td>
		<td>2024-07-13</td>
		<td>Sample Product</td>
		<td>$199.99</td>
		<td>$19.99</td>
		<td>$180.00</td>
	</tr>
	<tr>
		<td>Another Store</td>
		<td>Another Vendor</td>
		<td>July 14, 2024</td>
		<td>Another Product</td>
		<td>$299.50</td>
		<td>$29.95</td>
		<td>$269.55</td>
	</tr>
	`
	
	result, err := htmlParser.ParseHTML(inlineHTML)
	if err != nil {
		fmt.Printf("Inline parsing failed: %v\n", err)
	} else {
		printConsignableResults("Inline Consignable Data", result)
	}
	
	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	
	// Demo 3: Show what happens with wrong column count
	fmt.Println("3. Handling insufficient columns...")
	insufficientHTML := `
	<tr>
		<td>Store Only</td>
		<td>Vendor Only</td>
		<td>2024-07-13</td>
	</tr>
	`
	
	_, err = htmlParser.ParseHTML(insufficientHTML)
	if err != nil {
		fmt.Printf("Expected error for insufficient columns: %v\n", err)
	}
	
	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	
	// Demo 4: Custom column mapping
	fmt.Println("4. Custom column mapping example...")
	customParser := parser.NewHTMLTableParser()
	customParser.SetPositionalMapping([]string{
		"vendor",      // Column 0
		"store",       // Column 1
		"date",        // Column 2
		"description", // Column 3
		"sale_price",  // Column 4
	})
	
	customHTML := `
	<tr>
		<td>Custom Vendor</td>
		<td>Custom Store</td>
		<td>2024-07-13</td>
		<td>Custom Product</td>
		<td>$99.99</td>
	</tr>
	`
	
	result, err = customParser.ParseHTML(customHTML)
	if err != nil {
		fmt.Printf("Custom parsing failed: %v\n", err)
	} else {
		printConsignableResults("Custom Column Order", result)
	}
}

func printConsignableResults(title string, result *parser.ParseResult) {
	fmt.Printf("=== %s Results ===\n", title)
	fmt.Printf("Total Rows: %d\n", result.TotalRows)
	fmt.Printf("Successful: %d\n", result.SuccessCount)
	fmt.Printf("Errors: %d\n", result.ErrorCount)
	fmt.Printf("Processing Time: %v\n", result.Statistics.ProcessingTime)
	
	// Show successful records in a table format
	if result.SuccessCount > 0 {
		fmt.Printf("\nParsed Records:\n")
		fmt.Printf("%-15s %-20s %-12s %-25s %10s %10s %10s\n", 
			"Store", "Vendor", "Date", "Description", "Sale", "Comm", "Remaining")
		fmt.Printf("%s\n", strings.Repeat("-", 110))
		
		for _, record := range result.Records {
			// Truncate long descriptions for display
			desc := record.Description
			if len(desc) > 25 {
				desc = desc[:22] + "..."
			}
			
			fmt.Printf("%-15s %-20s %-12s %-25s %10.2f %10.2f %10.2f\n",
				truncate(record.Store, 15),
				truncate(record.Vendor, 20),
				record.Date,
				desc,
				record.SalePrice,
				record.Commission,
				record.Remaining)
		}
	}
	
	// Show errors if any
	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, parseErr := range result.Errors {
			fmt.Printf("  Row %d, Column '%s': %s",
				parseErr.Row, parseErr.Column, parseErr.Message)
			if parseErr.Value != "" {
				fmt.Printf(" (Value: '%s')", parseErr.Value)
			}
			fmt.Println()
		}
	}
	
	// Show warnings if any
	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("  Row %d, Column '%s': %s\n",
				warning.Row, warning.Column, warning.Message)
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

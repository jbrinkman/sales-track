package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sales-track/internal/parser"
)

func main() {
	// Get the directory of this example file
	exampleDir := filepath.Dir(os.Args[0])
	
	fmt.Println("=== HTML Table Parser Demo ===\n")
	
	// Create parser
	htmlParser := parser.NewHTMLTableParser()
	
	// Demo 1: Parse HTML file
	fmt.Println("1. Parsing HTML table file...")
	htmlFile := filepath.Join(exampleDir, "sample_table.html")
	if htmlData, err := ioutil.ReadFile(htmlFile); err == nil {
		result, err := htmlParser.ParseHTML(string(htmlData))
		if err != nil {
			log.Printf("HTML parsing failed: %v", err)
		} else {
			printResults("HTML Table", result)
		}
	} else {
		fmt.Printf("Could not read HTML file: %v\n", err)
	}
	
	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
	
	// Demo 2: Parse tab-delimited data
	fmt.Println("2. Parsing tab-delimited data...")
	tabFile := filepath.Join(exampleDir, "sample_tab_delimited.txt")
	if tabData, err := ioutil.ReadFile(tabFile); err == nil {
		result, err := htmlParser.ParseHTML(string(tabData))
		if err != nil {
			log.Printf("Tab-delimited parsing failed: %v", err)
		} else {
			printResults("Tab-Delimited Data", result)
		}
	} else {
		fmt.Printf("Could not read tab file: %v\n", err)
	}
	
	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
	
	// Demo 3: Parse inline HTML with errors
	fmt.Println("3. Parsing HTML with validation errors...")
	problemHTML := `
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
			<td>Valid Vendor</td>
			<td>invalid-date</td>
			<td>Valid Product</td>
			<td>not-a-price</td>
		</tr>
		<tr>
			<td>Valid Store</td>
			<td>Valid Vendor</td>
			<td>2024-01-15</td>
			<td>Valid Product</td>
			<td>$199.99</td>
		</tr>
	</table>
	`
	
	result, err := htmlParser.ParseHTML(problemHTML)
	if err != nil {
		log.Printf("Problem HTML parsing failed: %v", err)
	} else {
		printResults("HTML with Errors", result)
	}
	
	fmt.Println("\n" + strings.Repeat("=", 70) + "\n")
	
	// Demo 4: Run Consignable-specific demo
	RunConsignableDemo()
}

func printResults(title string, result *parser.ParseResult) {
	fmt.Printf("=== %s Results ===\n", title)
	fmt.Printf("Total Rows: %d\n", result.TotalRows)
	fmt.Printf("Successful: %d\n", result.SuccessCount)
	fmt.Printf("Errors: %d\n", result.ErrorCount)
	fmt.Printf("Warnings: %d\n", len(result.Warnings))
	fmt.Printf("Processing Time: %v\n", result.Statistics.ProcessingTime)
	fmt.Printf("Tables Found: %d\n", result.Statistics.TablesFound)
	
	if len(result.Statistics.HeadersDetected) > 0 {
		fmt.Printf("Headers: %v\n", result.Statistics.HeadersDetected)
	}
	
	if len(result.Statistics.DataTypesDetected) > 0 {
		fmt.Printf("Data Types: %v\n", result.Statistics.DataTypesDetected)
	}
	
	// Show successful records
	if result.SuccessCount > 0 {
		fmt.Printf("\nSuccessful Records:\n")
		for i, record := range result.Records {
			fmt.Printf("  %d. %s | %s | %s | %s | $%.2f\n",
				i+1, record.Store, record.Vendor, record.Date, 
				record.Description, record.SalePrice)
		}
	}
	
	// Show errors
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
	
	// Show warnings
	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("  Row %d, Column '%s': %s\n",
				warning.Row, warning.Column, warning.Message)
		}
	}
}

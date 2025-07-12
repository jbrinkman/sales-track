# Sales Track

A modern desktop application built with Go and Wails v2.10 to replace Excel-based sales tracking workflows. Features HTML table data import, SQLite storage, and pivot table-style reporting with drill-down capabilities.

## Overview

Sales Track modernizes the process of tracking sales data for physical storefronts by replacing manual Excel spreadsheet workflows with an intuitive desktop application.

### Current Workflow (Excel-based)
- Copy HTML table data from online website
- Paste into Excel spreadsheet with columns: Store, Vendor, Date, Description, Sale Price, Commission, Remaining
- Use pivot table to group by Year > Month > Date
- Display: Labels, Items Sold, Sum of Remaining

### New Workflow (Sales Track)
- Paste HTML table data directly into the application
- Automatic parsing and validation
- SQLite database storage
- Interactive pivot table-style reports with drill-down
- Year-based filtering and export capabilities

## Features

- **HTML Data Import**: Paste HTML table data directly from websites
- **Automatic Parsing**: Extract Store, Vendor, Date, Description, Sale Price, Commission, Remaining
- **SQLite Storage**: Reliable local database storage
- **Pivot Reports**: Group data by Year > Month > Date with drill-down navigation
- **Year Filtering**: Focus on specific years for targeted analysis
- **Data Export**: Export filtered data to CSV/Excel formats
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Technical Stack

- **Backend**: Go with Wails v2.10 framework
- **Database**: SQLite
- **Frontend**: Vue + TypeScript + UnoCSS
- **Platform**: Cross-platform desktop application

## Installation

### Prerequisites

- Go 1.19 or later
- Node.js 16 or later
- Wails v2.10

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/sales-track.git
   cd sales-track
   ```

2. Install dependencies:
   ```bash
   # Install Go dependencies
   go mod tidy
   
   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. Run in development mode:
   ```bash
   wails dev
   ```

4. Build for production:
   ```bash
   wails build
   ```

## Usage

1. **Import Data**: 
   - Copy HTML table data from your sales website
   - Paste into the import section of Sales Track
   - Click "Import" to parse and store the data

2. **View Reports**:
   - Navigate to the Reports section
   - Use the year filter to focus on specific periods
   - Click on year/month/date entries to drill down
   - View individual line items at the lowest level

3. **Export Data**:
   - Filter data as needed
   - Use the export function to save to CSV or Excel

## Data Schema

The application stores sales data with the following fields:
- **Store**: Store identifier
- **Vendor**: Vendor name
- **Date**: Sale date
- **Description**: Item description
- **Sale Price**: Sale amount
- **Commission**: Commission amount
- **Remaining**: Remaining balance

## Development

### Project Structure

```
sales-track/
├── app/                    # Go backend application
├── frontend/              # Vue frontend application
├── build/                 # Build configurations
├── database/             # Database migrations and schema
├── docs/                 # Documentation
└── tests/                # Test files
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or feature requests, please open an issue on GitHub.

# gocroissant

[![Version](https://img.shields.io/badge/version-v0.3.0-blue)](https://github.com/beyondcivic/gocroissant/releases/tag/v0.3.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/gocroissant.svg)](https://pkg.go.dev/github.com/beyondcivic/gocroissant)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go implementation for working with the [ML Commons Croissant](https://github.com/mlcommons/croissant) metadata format - a standardized way to describe machine learning datasets using JSON-LD.

## Overview

Croissant is an open metadata standard designed to improve dataset documentation, searchability, and usage in machine learning workflows. This library simplifies the creation of Croissant-compatible metadata from CSV data sources by:

- **Automatically inferring schema types** from dataset content
- **Generating complete, valid JSON-LD metadata** compliant with the Croissant specification
- **Providing comprehensive validation tools** to ensure compatibility
- **Supporting metadata comparison** for schema compatibility checking
- **Supporting the full Croissant specification** with extensible architecture

This project provides both a command-line interface and a Go library for working with Croissant metadata format.

## Key Features

- ✅ **Metadata Generation**: Convert CSV files to Croissant JSON-LD metadata
- ✅ **Schema Validation**: Comprehensive validation with detailed error reporting
- ✅ **Metadata Comparison**: Compare two metadata files for compatibility
- ✅ **Type Inference**: Automatic detection of data types (Boolean, Integer, Float, Date, URL, Text)
- ✅ **File Analysis**: Detailed CSV file inspection and statistics
- ✅ **JSON-LD Support**: Full JSON-LD processing and validation
- ✅ **CLI & Library**: Both command-line tool and Go library interfaces
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows

## Getting Started

### Prerequisites

- Go 1.24 or later
- Nix 2.25.4 or later (optional but recommended)
- PowerShell v7.5.1 or later (for building)

### Installation

#### Option 1: Install from Source

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/gocroissant.git
cd gocroissant
```

2. Build the application:

```bash
go build -o gocroissant .
```

#### Option 2: Using Nix (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/gocroissant.git
cd gocroissant
```

2. Prepare the environment using Nix flakes:

```bash
nix develop
```

3. Build the application:

```bash
./build.ps1
```

#### Option 3: Go Install

```bash
go install github.com/beyondcivic/gocroissant@latest
```

## Quick Start

### Command Line Interface

The `gocroissant` tool provides several commands for working with Croissant metadata:

```bash
# Generate metadata from CSV
gocroissant generate data.csv -o metadata.jsonld

# Validate existing metadata
gocroissant validate metadata.jsonld

# Compare two metadata files for compatibility
gocroissant match reference.jsonld candidate.jsonld

# Analyze CSV file structure
gocroissant info data.csv

# Show version information
gocroissant version
```

### Go Library Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
)

func main() {
	// Generate metadata from CSV
	outputPath, err := croissant.GenerateMetadata("data.csv", "dataset.jsonld")
	if err != nil {
		log.Fatalf("Error generating metadata: %v", err)
	}
	fmt.Printf("Metadata saved to: %s\n", outputPath)

	// Validate metadata
	issues, err := croissant.ValidateFile("dataset.jsonld")
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	if issues.HasErrors() {
		fmt.Println("Validation failed:")
		fmt.Println(issues.Report())
	} else {
		fmt.Println("Validation passed!")
	}

	// Compare two metadata files
	ref, err := croissant.LoadMetadataFromFile("reference.jsonld")
	if err != nil {
		log.Fatalf("Error loading reference: %v", err)
	}

	cand, err := croissant.LoadMetadataFromFile("candidate.jsonld")
	if err != nil {
		log.Fatalf("Error loading candidate: %v", err)
	}

	result := croissant.MatchMetadata(*ref, *cand)
	if result.IsMatch {
		fmt.Println("Files are compatible!")
	} else {
		fmt.Printf("Compatibility issues found: %d missing, %d mismatched\n",
			len(result.MissingFields), len(result.TypeMismatches))
	}
}
```

## Detailed Command Reference

### `generate` - Generate Metadata from CSV

Convert a CSV file to Croissant metadata format with automatic type inference.

```bash
gocroissant generate [CSV_FILE] [OPTIONS]
```

**Options:**

- `-o, --output`: Output file path (default: `[filename]_metadata.jsonld`)
- `-v, --validate`: Validate generated metadata and show issues
- `--strict`: Enable strict validation mode
- `--check-files`: Check if referenced files exist during validation

**Examples:**

```bash
# Basic generation
gocroissant generate data.csv

# With custom output path
gocroissant generate data.csv -o my-dataset.jsonld

# Generate and validate
gocroissant generate data.csv -o metadata.jsonld --validate

# Strict validation with file checking
gocroissant generate data.csv --validate --strict --check-files
```

### `validate` - Validate Existing Metadata

Validate a Croissant metadata file for compliance with the specification.

```bash
gocroissant validate [METADATA_FILE] [OPTIONS]
```

**Options:**

- `--strict`: Enable strict validation mode
- `--check-files`: Check if referenced files exist
- `--check-urls`: Validate URLs by making HTTP requests

**Examples:**

```bash
# Basic validation
gocroissant validate metadata.jsonld

# Strict validation with file and URL checking
gocroissant validate metadata.jsonld --strict --check-files --check-urls
```

### `match` - Compare Metadata Compatibility

Compare two Croissant metadata files to check schema compatibility. The candidate file can have additional fields, but must contain all fields from the reference with compatible data types.

```bash
gocroissant match [REFERENCE_FILE] [CANDIDATE_FILE] [OPTIONS]
```

**Options:**

- `-v, --verbose`: Show detailed information including extra fields in candidate

**Examples:**

```bash
# Basic compatibility check
gocroissant match reference.jsonld candidate.jsonld

# Verbose output showing all details
gocroissant match reference.jsonld candidate.jsonld --verbose
```

**Exit Codes:**

- `0`: Files are compatible
- `1`: Files are incompatible or error occurred

### `info` - Analyze CSV File

Display detailed information about a CSV file's structure, columns, and inferred data types.

```bash
gocroissant info [CSV_FILE] [OPTIONS]
```

**Options:**

- `--sample-size`: Number of rows to sample for type inference (default: 10)

**Examples:**

```bash
# Basic file analysis
gocroissant info data.csv

# Analyze with larger sample size
gocroissant info data.csv --sample-size 100
```

### `version` - Show Version Information

Display version, build information, and system details.

```bash
gocroissant version
```

## Data Type Inference

The tool automatically detects and maps data types to schema.org types:

| Detected Pattern           | Schema.org Type | Description     |
| -------------------------- | --------------- | --------------- |
| `true`, `false`, `1`, `0`  | `sc:Boolean`    | Boolean values  |
| `123`, `-456`              | `sc:Integer`    | Whole numbers   |
| `3.14`, `2.5e10`           | `sc:Float`      | Decimal numbers |
| `2023-01-01`, `01/15/2023` | `sc:Date`       | Date values     |
| `https://example.com`      | `sc:URL`        | Web URLs        |
| Everything else            | `sc:Text`       | Text content    |

## Configuration

The application supports configuration through environment variables:

- `CROISSANT_OUTPUT_PATH`: Default output path for generated metadata

If no output path is provided, the default format is `[filename]_metadata.jsonld`.

## Validation Features

The validation system provides comprehensive checking:

### Error Types

- **Mandatory fields**: Missing required properties
- **Type validation**: Incorrect `@type` values
- **Reference validation**: Invalid field or file references
- **Data type validation**: Invalid `dataType` specifications
- **JSON-LD structure**: Malformed JSON-LD documents

### Validation Modes

- **Standard mode**: Basic compliance checking
- **Strict mode**: Enhanced validation with additional warnings
- **File checking**: Verify referenced files exist
- **URL validation**: Check URL accessibility (optional)

## Schema Compatibility

The `match` command performs intelligent compatibility checking:

### Compatible Scenarios

- ✅ Candidate has additional fields (extra columns allowed)
- ✅ Numeric type compatibility (`sc:Number` ↔ `sc:Float`, `sc:Integer`)
- ✅ Exact field name and type matches

### Incompatible Scenarios

- ❌ Missing required fields from reference
- ❌ Type mismatches (e.g., `sc:Text` vs `sc:Integer`)
- ❌ Invalid JSON-LD structure

## Examples

### Example 1: Basic Metadata Generation

```bash
# Generate metadata for a simple CSV
$ gocroissant generate sales_data.csv -o sales_metadata.jsonld

Generating Croissant metadata for 'sales_data.csv'...
✓ Croissant metadata generated successfully and saved to: sales_metadata.jsonld
```

### Example 2: Generation with Validation

```bash
$ gocroissant generate sales_data.csv -o sales_metadata.jsonld --validate

Generating Croissant metadata for 'sales_data.csv'...
✓ Validation passed with no issues.
✓ Croissant metadata generated successfully and saved to: sales_metadata.jsonld
```

### Example 3: Metadata Validation

```bash
# Successful validation
$ gocroissant validate sales_metadata.jsonld
Validating Croissant metadata file 'sales_metadata.jsonld'...
✓ Validation passed with no issues.

# Validation with issues
$ gocroissant validate ./samples_jsonld/missing_fields.jsonld

Found the following 3 error(s) during the validation:
  -  [Metadata(mydataset) > FileObject(a-csv-table)] Property "https://schema.org/contentUrl" is mandatory, but does not exist.
  -  [Metadata(mydataset) > RecordSet(a-record-set) > Field(first-field)] The field does not specify a valid http://mlcommons.org/croissant/dataType, neither does any of its predecessor. Got:
  -  [Metadata(mydataset)] The current JSON-LD doesn't extend https://schema.org/Dataset.

Found the following 1 warning(s) during the validation:
  -  [Metadata(mydataset)] Property "http://purl.org/dc/terms/conformsTo" is recommended, but does not exist.
exit status 1
```

### Example 4: Metadata Comparison

```bash
# Compatible schemas
$ gocroissant match reference.jsonld candidate.jsonld

Comparing metadata files...
Reference: reference.jsonld
Candidate: candidate.jsonld

✓ Compatibility check PASSED
The candidate metadata is compatible with the reference.

Matched fields (6):
  ✓ transaction_id
  ✓ timestamp
  ✓ location
  ✓ water_flow_rate
  ✓ precipitation
  ✓ turbidity

Summary:
  Matched: 6
  Missing: 0
  Type mismatches: 0
  Extra fields: 0

# Incompatible schemas
$ gocroissant match reference.jsonld incompatible.jsonld

Comparing metadata files...
Reference: reference.jsonld
Candidate: incompatible.jsonld

✗ Compatibility check FAILED
The candidate metadata is NOT compatible with the reference.

Missing fields (3):
  ✗ transaction_id (required by reference but not found in candidate)
  ✗ timestamp (required by reference but not found in candidate)
  ✗ location (required by reference but not found in candidate)

Type mismatches (1):
  ✗ water_flow_rate: reference expects 'sc:Float', candidate has 'sc:Text'

Candidate has 2 additional field(s) (use --verbose to see details)

Summary:
  Matched: 2
  Missing: 3
  Type mismatches: 1
  Extra fields: 2
```

### Example 5: CSV File Analysis

```bash
$ gocroissant info sample_data.csv --sample-size 20

CSV File Information: sample_data.csv
=====================================
File Size: 15247 bytes
Total Rows: 101 (including header)
Data Rows: 100
Columns: 5
Sample Size: 20 rows

Column Information:
-------------------
1. id (sc:Integer)
2. name (sc:Text)
3. price (sc:Float)
4. available (sc:Boolean)
5. created_date (sc:Date)
```

## API Reference

### Core Functions

#### `GenerateMetadata(csvPath, outputPath string) (string, error)`

Generates Croissant metadata from a CSV file.

#### `GenerateMetadataWithValidation(csvPath, outputPath string) (*Metadata, error)`

Generates metadata and returns the parsed Metadata struct for further processing.

#### `ValidateFile(filePath string) (*Issues, error)`

Validates a Croissant metadata file and returns validation issues.

#### `ValidateJSON(data []byte) (*Issues, error)`

Validates Croissant metadata from JSON bytes.

#### `MatchMetadata(reference, candidate Metadata) *MatchResult`

Compares two metadata objects for schema compatibility.

#### `LoadMetadataFromFile(filePath string) (*Metadata, error)`

Loads and parses a Croissant metadata file.

### Data Structures

#### `Metadata`

Represents the complete Croissant metadata structure.

#### `MatchResult`

Contains the results of a metadata comparison:

- `IsMatch bool`: Whether the schemas are compatible
- `MissingFields []string`: Fields required by reference but missing in candidate
- `TypeMismatches []FieldMismatch`: Fields with incompatible data types
- `ExtraFields []string`: Additional fields in candidate
- `MatchedFields []string`: Successfully matched fields

#### `Issues`

Contains validation results with errors and warnings.

## Architecture

The library is organized into several key components:

### Core Package (`pkg/croissant`)

- **Generation**: CSV parsing and metadata creation
- **Validation**: JSON-LD and Croissant specification validation
- **Matching**: Schema compatibility checking
- **JSON-LD Processing**: JSON-LD document handling and validation
- **Type Inference**: Automatic data type detection from CSV content

### Command Line Interface (`cmd/gocroissant`)

- **Cobra-based CLI** with subcommands for each major function
- **Comprehensive help system** with detailed usage examples
- **Flexible output options** and validation modes

## Development

### Adding New Data Types

To add support for new data types, modify the `InferDataType` function in `pkg/croissant/croissant.go`:

```go
func InferDataType(value string) string {
	// Existing data type detection...

	// Add your new data type detection here
	if myCustomTypeDetector(value) {
		return "sc:MyCustomType"
	}

	// Default to Text
	return "sc:Text"
}
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `go test ./...`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to the branch: `git push origin feature/new-feature`
7. Submit a pull request

### Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Build Environment

### Using Nix (Recommended)

Use Nix flakes to set up the build environment:

```bash
nix develop
```

### Manual Build

Check the build arguments in `build.ps1`:

```bash
# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Then run:

```bash
./build.ps1
```

Or build manually:

```bash
go build -o gocroissant .
```

## Related Projects

- [ML Commons Croissant](https://github.com/mlcommons/croissant) - Official Croissant specification
- [Croissant Editor](https://github.com/mlcommons/croissant/tree/main/editor) - Web-based metadata editor
- [Python Croissant](https://github.com/mlcommons/croissant/tree/main/python) - Python implementation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

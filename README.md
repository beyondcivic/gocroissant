# gocroissant
Basic library that generates a JSON-LD file compatible with ML Commons Croissant metadata format from a CSV file.

This project provides a Go library and command-line tool for converting CSV files to Croissant metadata format.

## Library Structure

The project is structured as follows:

```
croissant-metadata/
├── croissant/
│   └── croissant.go    # Library code
├── main.go             # CLI application
├── go.mod              # Go module file
└── go.sum              # Go dependencies checksum
```

## Getting Started

### Prerequisites

- Go 1.24 or later

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/beyondcivic/gocroissant.git
   cd gocroissant
   ```

2. Build the application:
   ```
   go build -o gocroissant
   ```

### Usage

#### Command Line Interface

```bash
# Generate metadata with default output path
./gocroissant data.csv

# Specify output path
./gocroissant data.csv -o output.json
```

#### Using the Library in Your Go Code

```go
package main

import (
	"fmt"
	"log"

	"github.com/beyondcivic/gocroissant/pkg/croissant"
)

func main() {
	outputPath, err := croissant.GenerateMetadata("data.csv", "output.json")
	if err != nil {
		log.Fatalf("Error generating metadata: %v", err)
	}
	fmt.Printf("Metadata saved to: %s\n", outputPath)
}
```

## Features

- Automatically infers field data types from CSV content
- Calculates SHA-256 hash for file verification
- Generates Croissant metadata in JSON-LD format
- Configurable output path

## Configuration

The application supports configuration through environment variables with the prefix `CROISSANT_`.

## Development

### Adding New Data Types

To add support for new data types, modify the `InferDataType` function in `croissant/croissant.go`:

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

## License

TODO.

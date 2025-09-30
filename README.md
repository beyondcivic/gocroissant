# gocroissant

[![Version](https://img.shields.io/badge/version-v0.2.7-blue)](https://github.com/beyondcivic/gocroissant/releases/tag/v0.2.7)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-TBD-red)](LICENSE)

A Go implementation for working with the [ML Commons Croissant](https://github.com/mlcommons/croissant) metadata format - a standardized way to describe machine learning datasets using JSON-LD.

## Overview

Croissant is an open metadata standard designed to improve dataset documentation, searchability, and usage in machine learning workflows. This library simplifies the creation of Croissant-compatible metadata from CSV data sources by:

- Automatically inferring schema types from dataset content
- Generating complete, valid JSON-LD metadata
- Providing validation tools to ensure compatibility
- Supporting the full Croissant specification

This project provides both a command-line interface and a Go library for converting CSV files to Croissant metadata format.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Nix 2.25.4 or later
- Powershell v7.5.1 or later

### Installation

1. Clone the repository:

```bash
   git clone https://github.com/beyondcivic/gocroissant.git
   cd gocroissant
```

2. Prepare the environment using NIX flakes (optional but recommended):

```bash
nix develop
```

3. Build the application (NOTE: requires powershell v7.5.1):

```bash
./build.ps1
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
	outputPath, err := croissant.GenerateMetadata("data.csv", "dataset.jsonld")
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

NOTE: Currently, only `CROISSANT_OUTPUT_PATH` is supported to specify the output file path for generated metadata.

NOTE: If no output path is provided explicitly, the default output path `metadata.jsonld` will be used.

## Usage Examples

### Generate metadata without validation

```bash
$ gocroissant generate data.csv -o metadata.jsonld
```

### Generate metadata with validation

```bash
$ gocroissant generate data.csv -o metadata.jsonld -v

Validation passed with no issues.
Croissant metadata generated and saved to: metadata.json
```

### Generate metadata with validation but without saving to a file

```bash
$ gocroissant generate data.csv -v
Validation passed with no issues.
```

### Validate an existing metadata file

```bash
$ gocroissant validate metadata.json
Validation passed with no issues.
```

### Example with issues

```
$ gocroissant validate ./samples_jsonld/missing_fields.jsonld

Found the following 3 error(s) during the validation:
  -  [Metadata(mydataset) > FileObject(a-csv-table)] Property "https://schema.org/contentUrl" is mandatory, but does not exist.
  -  [Metadata(mydataset) > RecordSet(a-record-set) > Field(first-field)] The field does not specify a valid http://mlcommons.org/croissant/dataType, neither does any of its predecessor. Got:
  -  [Metadata(mydataset)] The current JSON-LD doesn't extend https://schema.org/Dataset.

Found the following 1 warning(s) during the validation:
  -  [Metadata(mydataset)] Property "http://purl.org/dc/terms/conformsTo" is recommended, but does not exist.
exit status 1
```

```bash
$ gocroissant validate ./samples_jsonld/invalid_references.jsonld

Found the following 1 error(s) during the validation:
  -  [Metadata(mydataset) > FileObject(a-csv-table)] "a-csv-table" should have an attribute "@type": "http://mlcommons.org/croissant/FileObject" or "@type": "http://mlcommons.org/croissant/FileSet". Got sc:WRONG_TYPE instead.
exit status 1
```

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

## Build environment

_NOTE: Using Nix is optional, you can also rely on your local Go installation._

Use NIX flakes to setup the build environment.

```bash
nix develop
```

Check the build arguments in build.ps1, e.g.:

```bash
# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Then run the following command to build the project:

```bash
./build.ps1
```

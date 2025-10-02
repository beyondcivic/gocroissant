# API Documentation

This document provides comprehensive API documentation for the gocroissant library.

## Package Overview

The `github.com/beyondcivic/gocroissant/pkg/croissant` package provides functionality for working with the ML Commons Croissant metadata format.

## Core Functions

### Metadata Generation

#### `GenerateMetadata(csvPath, outputPath string) (string, error)`

Generates Croissant metadata from a CSV file and saves it to the specified output path.

**Parameters:**

- `csvPath`: Path to the input CSV file
- `outputPath`: Path where the generated metadata will be saved

**Returns:**

- `string`: The actual output path used (may differ from input if empty)
- `error`: Any error that occurred during generation

**Example:**

```go
outputPath, err := croissant.GenerateMetadata("data.csv", "metadata.jsonld")
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Metadata saved to: %s\n", outputPath)
```

#### `GenerateMetadataWithValidation(csvPath, outputPath string) (*Metadata, error)`

Generates metadata and returns the parsed Metadata struct for further processing.

**Parameters:**

- `csvPath`: Path to the input CSV file
- `outputPath`: Path where the generated metadata will be saved (can be empty)

**Returns:**

- `*Metadata`: The generated metadata structure
- `error`: Any error that occurred during generation

**Example:**

```go
metadata, err := croissant.GenerateMetadataWithValidation("data.csv", "metadata.jsonld")
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Generated dataset: %s with %d record sets\n",
    metadata.Name, len(metadata.RecordSets))
```

### Validation

#### `ValidateFile(filePath string) (*Issues, error)`

Validates a Croissant metadata file and returns detailed validation results.

**Parameters:**

- `filePath`: Path to the metadata file to validate

**Returns:**

- `*Issues`: Validation results containing errors and warnings
- `error`: Any error that occurred during validation

**Example:**

```go
issues, err := croissant.ValidateFile("metadata.jsonld")
if err != nil {
    log.Fatalf("Validation error: %v", err)
}

if issues.HasErrors() {
    fmt.Println("Validation failed:")
    fmt.Println(issues.Report())
} else {
    fmt.Println("Validation passed!")
}
```

#### `ValidateJSON(data []byte) (*Issues, error)`

Validates Croissant metadata from JSON bytes.

**Parameters:**

- `data`: JSON-LD content as bytes

**Returns:**

- `*Issues`: Validation results
- `error`: Any error that occurred

#### `ValidateJSONWithOptions(data []byte, options ValidationOptions) (*Issues, error)`

Validates metadata with custom validation options.

**Parameters:**

- `data`: JSON-LD content as bytes
- `options`: Validation configuration options

**Returns:**

- `*Issues`: Validation results
- `error`: Any error that occurred

**Example:**

```go
options := croissant.ValidationOptions{
    StrictMode:      true,
    CheckDataTypes:  true,
    ValidateURLs:    false,
    CheckFileExists: true,
}

issues, err := croissant.ValidateJSONWithOptions(data, options)
```

### Schema Comparison

#### `MatchMetadata(reference, candidate Metadata) *MatchResult`

Compares two metadata objects for schema compatibility.

**Parameters:**

- `reference`: The reference metadata to compare against
- `candidate`: The candidate metadata to check for compatibility

**Returns:**

- `*MatchResult`: Detailed comparison results

**Example:**

```go
ref, _ := croissant.LoadMetadataFromFile("reference.jsonld")
cand, _ := croissant.LoadMetadataFromFile("candidate.jsonld")

result := croissant.MatchMetadata(*ref, *cand)
if result.IsMatch {
    fmt.Printf("Compatible! %d fields matched\n", len(result.MatchedFields))
} else {
    fmt.Printf("Incompatible: %d missing, %d mismatched\n",
        len(result.MissingFields), len(result.TypeMismatches))
}
```

#### `LoadMetadataFromFile(filePath string) (*Metadata, error)`

Loads and parses a Croissant metadata file from disk.

**Parameters:**

- `filePath`: Path to the metadata file

**Returns:**

- `*Metadata`: The parsed metadata structure
- `error`: Any error that occurred

## Data Structures

### `Metadata`

Represents the complete Croissant metadata structure with all required and optional fields.

**Key Fields:**

- `Name string`: Dataset name
- `Description string`: Dataset description
- `RecordSets []RecordSet`: Array of record sets
- `Distributions []Distribution`: Array of file distributions
- `ConformsTo string`: Croissant specification version

### `MatchResult`

Contains the results of comparing two metadata files.

**Fields:**

- `IsMatch bool`: Whether schemas are compatible
- `MissingFields []string`: Fields required by reference but missing in candidate
- `TypeMismatches []FieldMismatch`: Fields with incompatible data types
- `ExtraFields []string`: Additional fields in candidate
- `MatchedFields []string`: Successfully matched fields

### `FieldMismatch`

Represents a field with incompatible data types between reference and candidate.

**Fields:**

- `FieldName string`: Name of the mismatched field
- `ReferenceType string`: Expected data type from reference
- `CandidateType string`: Actual data type in candidate

### `Issues`

Contains validation results with errors and warnings.

**Methods:**

- `HasErrors() bool`: Returns true if validation errors exist
- `HasWarnings() bool`: Returns true if validation warnings exist
- `ErrorCount() int`: Returns the number of errors
- `WarningCount() int`: Returns the number of warnings
- `Report() string`: Returns a formatted report of all issues

### `ValidationOptions`

Configuration options for metadata validation.

**Fields:**

- `StrictMode bool`: Enable additional warnings and strict checking
- `CheckDataTypes bool`: Validate data type specifications
- `ValidateURLs bool`: Check URL accessibility (requires network calls)
- `CheckFileExists bool`: Verify referenced files exist on disk

## Type Inference

The library automatically infers schema.org data types:

| Input Pattern              | Detected Type | Schema.org Type |
| -------------------------- | ------------- | --------------- |
| `true`, `false`, `1`, `0`  | Boolean       | `sc:Boolean`    |
| `123`, `-456`              | Integer       | `sc:Integer`    |
| `3.14`, `2.5e10`           | Float         | `sc:Float`      |
| `2023-01-01`, `01/15/2023` | Date          | `sc:Date`       |
| `https://example.com`      | URL           | `sc:URL`        |
| Everything else            | Text          | `sc:Text`       |

## Error Handling

All functions return Go errors following standard conventions. Common error types include:

- **File I/O errors**: File not found, permission denied
- **JSON parsing errors**: Invalid JSON syntax
- **JSON-LD validation errors**: Invalid JSON-LD structure
- **Croissant validation errors**: Non-compliance with specification
- **CSV parsing errors**: Invalid CSV structure or encoding

## Best Practices

1. **Always check errors**: All functions return errors that should be handled
2. **Use validation options**: Configure validation based on your requirements
3. **Handle missing files gracefully**: Use `CheckFileExists` option when appropriate
4. **Cache metadata objects**: Avoid repeated file parsing when possible
5. **Use type-safe comparisons**: Leverage the `MatchMetadata` function for schema compatibility

## Version Compatibility

This library is compatible with:

- Croissant specification version 1.0
- JSON-LD 1.1
- Go 1.24+

## Related Documentation

- [Command Line Reference](cmd/gocroissant.md)
- [Examples](../examples/)
- [ML Commons Croissant Specification](https://github.com/mlcommons/croissant)

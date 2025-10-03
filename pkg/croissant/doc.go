/*
Package croissant provides comprehensive functionality for working with the ML Commons Croissant
metadata format - a standardized way to describe machine learning datasets using JSON-LD.

# Overview

The Croissant metadata format is an open standard designed to improve dataset documentation,
searchability, and usage in machine learning workflows. This package simplifies working with
Croissant metadata by providing:

  - Automatic metadata generation from CSV files with intelligent type inference
  - Comprehensive validation tools with detailed error reporting
  - Schema compatibility checking for dataset evolution
  - Full JSON-LD processing and validation support
  - Extensible architecture supporting the complete Croissant specification

# Quick Start

Generate metadata from a CSV file:

	outputPath, err := croissant.GenerateMetadata("data.csv", "metadata.jsonld")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Metadata generated: %s\n", outputPath)

Validate existing metadata:

	issues, err := croissant.ValidateFile("metadata.jsonld")
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	if issues.HasErrors() {
		fmt.Println("Validation failed:")
		fmt.Println(issues.Report())
	} else {
		fmt.Println("✓ Validation passed!")
	}

Compare metadata files for compatibility:

	ref, _ := croissant.LoadMetadataFromFile("reference.jsonld")
	cand, _ := croissant.LoadMetadataFromFile("candidate.jsonld")

	result := croissant.MatchMetadata(*ref, *cand)
	if result.IsMatch {
		fmt.Printf("✓ Compatible schemas with %d matched fields\n", len(result.MatchedFields))
	} else {
		fmt.Printf("✗ Incompatible: %d missing, %d type mismatches\n",
			len(result.MissingFields), len(result.TypeMismatches))
	}

# Features

## Metadata Generation

The package automatically generates Croissant-compliant metadata from CSV files:

  - Intelligent data type inference (Boolean, Integer, Float, Date, URL, Text)
  - SHA-256 hash calculation for file integrity verification
  - Configurable output paths and validation options
  - Support for environment variable configuration

## Validation System

Comprehensive validation with multiple modes:

  - JSON-LD structure validation using the json-gold library
  - Croissant specification compliance checking
  - Configurable validation modes (standard, strict)
  - Optional file existence and URL accessibility verification
  - Detailed error reporting with contextual information

## Schema Compatibility

Advanced schema comparison for dataset evolution:

  - Field-by-field compatibility analysis
  - Intelligent type compatibility (numeric type flexibility)
  - Support for schema evolution (additional fields allowed)
  - Detailed reporting of matches, mismatches, and missing fields

# Data Type Inference

The package automatically maps CSV content to appropriate schema.org types:

	Input Pattern              → Detected Type → Schema.org Type
	===========================================================
	true, false, 1, 0         → Boolean       → sc:Boolean
	123, -456                 → Integer       → sc:Integer
	3.14, 2.5e10             → Float         → sc:Float
	2023-01-01, 01/15/2023   → Date          → sc:Date
	https://example.com       → URL           → sc:URL
	Everything else           → Text          → sc:Text

# Validation Options

Customize validation behavior using ValidationOptions:

	options := croissant.ValidationOptions{
		StrictMode:      true,  // Enable additional warnings
		CheckDataTypes:  true,  // Validate data type specifications
		ValidateURLs:    false, // Skip network calls for URL validation
		CheckFileExists: true,  // Verify referenced files exist
	}

	issues, err := croissant.ValidateJSONWithOptions(data, options)

# Schema Compatibility Rules

When comparing metadata files, the following rules apply:

  - All fields in the reference must exist in the candidate
  - Field data types must be compatible (exact or compatible numeric types)
  - Additional fields in the candidate are allowed (backward compatibility)
  - Compatible numeric types: sc:Number accepts sc:Float and sc:Integer

# Error Handling

All functions follow Go error handling conventions. Common error scenarios:

  - File I/O errors (file not found, permission denied)
  - JSON parsing errors (invalid JSON syntax)
  - JSON-LD validation errors (invalid JSON-LD structure)
  - Croissant validation errors (specification non-compliance)
  - CSV parsing errors (invalid structure or encoding)

# Performance Considerations

  - Metadata objects can be cached to avoid repeated file parsing
  - Large CSV files are processed incrementally for memory efficiency
  - URL validation is optional to avoid network latency
  - File existence checks can be disabled for performance

# Examples

See the examples directory for comprehensive usage examples:

  - Basic metadata generation and validation
  - Advanced validation with custom options
  - Schema compatibility checking
  - Error handling patterns

# Related Tools

This package includes a command-line tool (gocroissant) that provides:

  - generate: Convert CSV files to Croissant metadata
  - validate: Validate existing metadata files
  - match: Compare metadata files for compatibility
  - info: Analyze CSV file structure
  - version: Display version information

# Specification Compliance

This implementation supports:

  - Croissant specification version 1.0
  - JSON-LD 1.1 processing
  - Schema.org vocabulary
  - Full Croissant metadata structure

# License

MIT License - see LICENSE file for details.

# Related Projects

  - ML Commons Croissant: https://github.com/mlcommons/croissant
  - Croissant Editor: Web-based metadata editor
  - Python Croissant: Python implementation
*/
package croissant

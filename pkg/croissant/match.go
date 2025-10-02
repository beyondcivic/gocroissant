package croissant

import (
	"os"
	"strings"
)

// MatchResult represents the result of comparing two Croissant metadata files for schema compatibility.
// It provides detailed information about field matches, mismatches, and additional fields.
//
// The comparison checks whether a candidate metadata file is compatible with a reference metadata file.
// Compatibility means:
//   - All fields from the reference must exist in the candidate
//   - Field data types must be compatible (exact match or compatible numeric types)
//   - The candidate may have additional fields (this doesn't affect compatibility)
//
// Example usage:
//
//	ref, _ := croissant.LoadMetadataFromFile("reference.jsonld")
//	cand, _ := croissant.LoadMetadataFromFile("candidate.jsonld")
//	result := croissant.MatchMetadata(*ref, *cand)
//
//	if result.IsMatch {
//		fmt.Printf("Compatible! %d fields matched\n", len(result.MatchedFields))
//	} else {
//		fmt.Printf("Issues: %d missing, %d type mismatches\n",
//			len(result.MissingFields), len(result.TypeMismatches))
//	}
type MatchResult struct {
	// IsMatch indicates whether the candidate is compatible with the reference.
	// True if all reference fields exist in candidate with compatible types.
	IsMatch bool

	// MissingFields lists field names that exist in reference but not in candidate.
	// These represent compatibility violations.
	MissingFields []string

	// TypeMismatches lists fields that exist in both files but have incompatible data types.
	// These represent compatibility violations.
	TypeMismatches []FieldMismatch

	// ExtraFields lists field names that exist in candidate but not in reference.
	// These do not affect compatibility but may be useful for information.
	ExtraFields []string

	// MatchedFields lists field names that exist in both files with compatible types.
	// These represent successful matches.
	MatchedFields []string
}

// FieldMismatch represents a field that exists in both metadata files but has incompatible data types.
// This indicates a schema compatibility issue that prevents the candidate from being used
// as a drop-in replacement for the reference.
type FieldMismatch struct {
	// FieldName is the name of the field that has a type mismatch.
	FieldName string

	// ReferenceType is the data type expected by the reference metadata.
	ReferenceType string

	// CandidateType is the data type found in the candidate metadata.
	CandidateType string
}

// MatchMetadata compares two Croissant metadata objects to check if the candidate
// is compatible with the reference. The candidate can have additional fields,
// but all reference fields must exist in the candidate with matching data types.
//
// Compatibility Rules:
//   - All fields in the reference must exist in the candidate
//   - Field data types must be compatible (see type compatibility below)
//   - Additional fields in the candidate are allowed and don't affect compatibility
//
// Type Compatibility:
//   - Exact type matches (sc:Text = sc:Text)
//   - Numeric type compatibility (sc:Number accepts sc:Float, sc:Integer)
//   - Schema.org prefix normalization (sc:Text = https://schema.org/Text)
//
// The function returns a MatchResult containing detailed information about:
//   - Whether the schemas are compatible (IsMatch)
//   - Successfully matched fields (MatchedFields)
//   - Missing required fields (MissingFields)
//   - Type mismatches (TypeMismatches)
//   - Additional fields in candidate (ExtraFields)
//
// Example:
//
//	reference, err := croissant.LoadMetadataFromFile("reference.jsonld")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	candidate, err := croissant.LoadMetadataFromFile("candidate.jsonld")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	result := croissant.MatchMetadata(*reference, *candidate)
//
//	if result.IsMatch {
//		fmt.Printf("✓ Compatible schemas with %d matched fields\n", len(result.MatchedFields))
//		if len(result.ExtraFields) > 0 {
//			fmt.Printf("  Candidate has %d additional fields\n", len(result.ExtraFields))
//		}
//	} else {
//		fmt.Printf("✗ Incompatible schemas:\n")
//		if len(result.MissingFields) > 0 {
//			fmt.Printf("  Missing %d required fields\n", len(result.MissingFields))
//		}
//		if len(result.TypeMismatches) > 0 {
//			fmt.Printf("  %d type mismatches found\n", len(result.TypeMismatches))
//		}
//	}
func MatchMetadata(reference Metadata, candidate Metadata) *MatchResult {
	result := &MatchResult{
		IsMatch:        true,
		MissingFields:  []string{},
		TypeMismatches: []FieldMismatch{},
		ExtraFields:    []string{},
		MatchedFields:  []string{},
	}

	// Create maps for easy lookup
	refFields := make(map[string]Field)
	candFields := make(map[string]Field)

	// Collect all fields from all record sets in reference
	for _, rs := range reference.RecordSets {
		for _, field := range rs.Fields {
			refFields[field.Name] = field
		}
	}

	// Collect all fields from all record sets in candidate
	for _, rs := range candidate.RecordSets {
		for _, field := range rs.Fields {
			candFields[field.Name] = field
		}
	}

	// Check each reference field
	for refFieldName, refField := range refFields {
		if candField, exists := candFields[refFieldName]; exists {
			// Field exists, check data type compatibility
			if !areDataTypesCompatible(refField.DataType, candField.DataType) {
				result.TypeMismatches = append(result.TypeMismatches, FieldMismatch{
					FieldName:     refFieldName,
					ReferenceType: getDataTypeString(refField.DataType),
					CandidateType: getDataTypeString(candField.DataType),
				})
				result.IsMatch = false
			} else {
				result.MatchedFields = append(result.MatchedFields, refFieldName)
			}
		} else {
			// Field is missing in candidate
			result.MissingFields = append(result.MissingFields, refFieldName)
			result.IsMatch = false
		}
	}

	// Find extra fields in candidate (this doesn't affect compatibility)
	for candFieldName := range candFields {
		if _, exists := refFields[candFieldName]; !exists {
			result.ExtraFields = append(result.ExtraFields, candFieldName)
		}
	}

	return result
}

// areDataTypesCompatible checks if two DataType objects are compatible
func areDataTypesCompatible(ref DataType, cand DataType) bool {
	refTypes := ref.GetTypes()
	candTypes := cand.GetTypes()

	if len(refTypes) == 0 || len(candTypes) == 0 {
		return false
	}

	// Normalize and check type compatibility
	refPrimaryType := normalizeDataType(refTypes[0])
	candPrimaryType := normalizeDataType(candTypes[0])

	// Direct match
	if refPrimaryType == candPrimaryType {
		return true
	}

	// Check for compatible numeric types
	if isCompatibleNumericType(refPrimaryType, candPrimaryType) {
		return true
	}

	return false
}

// isCompatibleNumericType checks if two numeric types are compatible
func isCompatibleNumericType(refType, candType string) bool {
	// Define type compatibility mappings
	numericTypes := map[string][]string{
		"number":  {"number", "integer", "float"},
		"integer": {"number", "integer"},
		"float":   {"number", "float"},
	}

	// Check if candidate type is compatible with reference type
	if compatibleTypes, exists := numericTypes[refType]; exists {
		for _, compatibleType := range compatibleTypes {
			if candType == compatibleType {
				return true
			}
		}
	}

	return false
}

// normalizeDataType normalizes data type strings for comparison
func normalizeDataType(dataType string) string {
	// Remove common prefixes and normalize to lowercase
	dataType = strings.ToLower(dataType)
	dataType = strings.TrimPrefix(dataType, "sc:")
	dataType = strings.TrimPrefix(dataType, "https://schema.org/")
	dataType = strings.TrimPrefix(dataType, "http://schema.org/")
	return dataType
}

// getDataTypeString returns a string representation of a DataType
func getDataTypeString(dt DataType) string {
	types := dt.GetTypes()
	if len(types) == 0 {
		return "unknown"
	}
	if len(types) == 1 {
		return types[0]
	}
	return "[" + strings.Join(types, ", ") + "]"
}

// LoadMetadataFromFile loads and parses a Croissant metadata file from disk.
// It validates the JSON-LD structure and parses it into a Metadata object.
//
// The function performs the following steps:
//  1. Reads the file from the specified path
//  2. Validates that the content is valid JSON-LD using the json-gold library
//  3. Parses the JSON-LD into a Croissant Metadata structure
//  4. Returns the parsed metadata or an error if any step fails
//
// Supported file formats:
//   - JSON-LD files (.jsonld, .json)
//   - Must conform to Croissant metadata specification
//   - Must be valid JSON-LD documents
//
// Example:
//
//	metadata, err := croissant.LoadMetadataFromFile("dataset.jsonld")
//	if err != nil {
//		log.Fatalf("Failed to load metadata: %v", err)
//	}
//
//	fmt.Printf("Loaded dataset: %s\n", metadata.Name)
//	fmt.Printf("Record sets: %d\n", len(metadata.RecordSets))
//	fmt.Printf("Distributions: %d\n", len(metadata.Distributions))
//
// Common errors:
//   - File not found or permission denied
//   - Invalid JSON syntax
//   - Invalid JSON-LD structure
//   - Non-compliant Croissant metadata format
func LoadMetadataFromFile(filePath string) (*Metadata, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, CroissantError{Message: "failed to read file", Value: err}
	}

	processor := NewJSONLDProcessor()

	// Validate JSON-LD structure
	if err := processor.ValidateJSONLD(data); err != nil {
		return nil, CroissantError{Message: "invalid JSON-LD document", Value: err}
	}

	// Parse the metadata
	metadata, err := processor.ParseCroissantMetadata(data)
	if err != nil {
		return nil, CroissantError{Message: "failed to parse Croissant metadata", Value: err}
	}

	return metadata, nil
}

// datatypes.go
package croissant

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// InferDataType infers the schema.org data type from a value
func InferDataType(value string) string {
	// Trim whitespace
	value = strings.TrimSpace(value)
	if value == "" {
		return "sc:Text"
	}

	// Try to parse as boolean
	lowerVal := strings.ToLower(value)
	if lowerVal == "true" || lowerVal == "false" {
		return "sc:Boolean"
	}

	// Try to parse as integer
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "sc:Integer"
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "sc:Number"
	}

	// Try to parse as date (various formats)
	dateFormats := []string{
		"2006-01-02",
		"01/02/2006",
		"2006/01/02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
	}
	for _, format := range dateFormats {
		if _, err := time.Parse(format, value); err == nil {
			return "sc:DateTime"
		}
	}

	// Try to parse as URL
	if _, err := url.ParseRequestURI(value); err == nil && (strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")) {
		return "sc:URL"
	}

	// Try to detect email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegex.MatchString(value) {
		return "sc:Text" // Email is still text but we could add a custom type if needed
	}

	// Default to Text
	return "sc:Text"
}

// IsValidDataType checks if a dataType is valid according to Croissant specification
func IsValidDataType(dataType string) bool {
	validTypes := map[string]bool{
		// Schema.org types
		"sc:Text":           true,
		"sc:Boolean":        true,
		"sc:Integer":        true,
		"sc:Number":         true,
		"sc:DateTime":       true,
		"sc:URL":            true,
		"sc:ImageObject":    true,
		"sc:VideoObject":    true,
		"sc:Enumeration":    true,
		"sc:GeoShape":       true,
		"sc:GeoCoordinates": true,

		// Croissant-specific types
		"cr:Label":            true,
		"cr:Split":            true,
		"cr:BoundingBox":      true,
		"cr:SegmentationMask": true,

		// Croissant Split types
		"cr:TrainingSplit":   true,
		"cr:ValidationSplit": true,
		"cr:TestSplit":       true,
	}

	// Also accept Wikidata entities (wd:Q...)
	if strings.HasPrefix(dataType, "wd:Q") {
		return true
	}

	return validTypes[dataType]
}

// InferSemanticDataType attempts to infer semantic data types for ML datasets
func InferSemanticDataType(fieldName, value string, context map[string]interface{}) []string {
	fieldNameLower := strings.ToLower(fieldName)

	// Check for split-related fields
	if strings.Contains(fieldNameLower, "split") {
		splitValues := []string{"train", "training", "val", "validation", "test", "testing"}
		for _, splitVal := range splitValues {
			if strings.Contains(strings.ToLower(value), splitVal) {
				return []string{"cr:Split", "sc:Text"}
			}
		}
	}

	// Check for label-related fields
	labelFields := []string{"label", "class", "category", "target", "annotation"}
	for _, labelField := range labelFields {
		if strings.Contains(fieldNameLower, labelField) {
			baseType := InferDataType(value)
			return []string{baseType, "cr:Label"}
		}
	}

	// Check for bounding box patterns (arrays of 4 numbers)
	if strings.Contains(fieldNameLower, "bbox") || strings.Contains(fieldNameLower, "box") {
		// This would need more sophisticated parsing for actual bounding box detection
		return []string{"cr:BoundingBox"}
	}

	// Check for enumeration patterns
	if context != nil {
		if enumValues, exists := context["enumValues"]; exists {
			if enumSlice, ok := enumValues.([]string); ok {
				for _, enumVal := range enumSlice {
					if value == enumVal {
						return []string{"sc:Enumeration", "sc:Text"}
					}
				}
			}
		}
	}

	// Default to basic type inference
	return []string{InferDataType(value)}
}

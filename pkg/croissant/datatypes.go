// datatypes.go
package croissant

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Supported value data types:
// Schema.org data types
const VT_scText string = "sc:Text"
const VT_scBool string = "sc:Boolean"
const VT_scInt string = "sc:Integer"
const VT_scNum string = "sc:Number"
const VT_scDateT string = "sc:DateTime"
const VT_scURL string = "sc:URL"
const VT_scImage string = "sc:ImageObject"
const VT_scVideo string = "sc:VideoObject"
const VT_scEnum string = "sc:Enumeration"
const VT_scGeoShape string = "sc:GeoShape"
const VT_scGeoCoord string = "sc:GeoCoordinates"

// Croissant-specific types
const VT_crLabel string = "cr:Label"
const VT_crSplit string = "cr:Split"
const VT_crBBox string = "cr:BoundingBox"
const VT_crSegMask string = "cr:SegmentationMask"

// Croissant Split types
const VT_crSplitTrain string = "cr:TrainingSplit"
const VT_crSplitVal string = "cr:ValidationSplit"
const VT_crSplitTest string = "cr:TestSplit"

// Wikidata entities (wd:Q...)
const VT_wdPrefix string = "wd:Q"

// InferDataType infers the schema.org data type from a value.
// Returns in
func InferDataType(value string) string {
	// Trim whitespace
	value = strings.TrimSpace(value)
	if value == "" {
		return VT_scText
	}

	// Try to parse as boolean
	lowerVal := strings.ToLower(value)
	if lowerVal == "true" || lowerVal == "false" {
		return VT_scBool
	}

	// Try to parse as integer
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return VT_scInt
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return VT_scNum
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
			return VT_scDateT
		}
	}

	// Try to parse as URL
	if _, err := url.ParseRequestURI(value); err == nil && (strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")) {
		return VT_scURL
	}

	// Try to detect email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegex.MatchString(value) {
		return VT_scText // Email is still text but we could add a custom type if needed
	}

	// Default to Text
	return VT_scText
}

// IsValidDataType checks if a dataType is valid according to Croissant specification
func IsValidDataType(dataType string) bool {
	validTypes := map[string]bool{
		// Schema.org types
		VT_scText:     true,
		VT_scBool:     true,
		VT_scInt:      true,
		VT_scNum:      true,
		VT_scDateT:    true,
		VT_scURL:      true,
		VT_scImage:    true,
		VT_scVideo:    true,
		VT_scEnum:     true,
		VT_scGeoShape: true,
		VT_scGeoCoord: true,

		// Croissant-specific types
		VT_crLabel:   true,
		VT_crSplit:   true,
		VT_crBBox:    true,
		VT_crSegMask: true,

		// Croissant Split types
		VT_crSplitTrain: true,
		VT_crSplitVal:   true,
		VT_crSplitTest:  true,
	}

	// Also accept Wikidata entities (wd:Q...)
	if strings.HasPrefix(dataType, VT_wdPrefix) {
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
				return []string{VT_crSplit, VT_scText}
			}
		}
	}

	// Check for label-related fields
	labelFields := []string{"label", "class", "category", "target", "annotation"}
	for _, labelField := range labelFields {
		if strings.Contains(fieldNameLower, labelField) {
			baseType := InferDataType(value)
			return []string{baseType, VT_crLabel}
		}
	}

	// Check for bounding box patterns (arrays of 4 numbers)
	if strings.Contains(fieldNameLower, "bbox") || strings.Contains(fieldNameLower, "box") {
		// This would need more sophisticated parsing for actual bounding box detection
		return []string{VT_crBBox}
	}

	// Check for enumeration patterns
	if context != nil {
		if enumValues, exists := context["enumValues"]; exists {
			if enumSlice, ok := enumValues.([]string); ok {
				for _, enumVal := range enumSlice {
					if value == enumVal {
						return []string{VT_scEnum, VT_scText}
					}
				}
			}
		}
	}

	// Default to basic type inference
	return []string{InferDataType(value)}
}

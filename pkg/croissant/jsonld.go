// jsonld.go
package croissant

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/piprate/json-gold/ld"
)

// JSONLDProcessor handles JSON-LD processing using json-gold library
type JSONLDProcessor struct {
	processor *ld.JsonLdProcessor
	options   *ld.JsonLdOptions
}

// NewJSONLDProcessor creates a new JSON-LD processor
func NewJSONLDProcessor() *JSONLDProcessor {
	return &JSONLDProcessor{
		processor: ld.NewJsonLdProcessor(),
		options:   ld.NewJsonLdOptions(""),
	}
}

// ParseJSONLD parses and expands JSON-LD document to a normalized form
func (j *JSONLDProcessor) ParseJSONLD(data []byte) (map[string]interface{}, error) {
	// First, parse as regular JSON to get a map
	var jsonDoc interface{}
	if err := json.Unmarshal(data, &jsonDoc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Expand the JSON-LD document using json-gold
	expanded, err := j.processor.Expand(jsonDoc, j.options)
	if err != nil {
		return nil, fmt.Errorf("failed to expand JSON-LD: %w", err)
	}

	// Convert expanded document back to map[string]interface{}
	if len(expanded) == 0 {
		return make(map[string]interface{}), nil
	}

	// For Croissant documents, we typically expect a single document
	if expandedMap, ok := expanded[0].(map[string]interface{}); ok {
		return expandedMap, nil
	}

	return nil, fmt.Errorf("unexpected expanded JSON-LD structure")
}

// CompactJSONLD compacts an expanded JSON-LD document with the given context
func (j *JSONLDProcessor) CompactJSONLD(expanded interface{}, context map[string]interface{}) (map[string]interface{}, error) {
	compacted, err := j.processor.Compact(expanded, context, j.options)
	if err != nil {
		return nil, fmt.Errorf("failed to compact JSON-LD: %w", err)
	}

	// The compacted result should already be a map[string]interface{}
	return compacted, nil
}

// ValidateJSONLD validates that the document is valid JSON-LD
func (j *JSONLDProcessor) ValidateJSONLD(data []byte) error {
	var jsonDoc interface{}
	if err := json.Unmarshal(data, &jsonDoc); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Try to expand the document - this will fail if it's not valid JSON-LD
	_, err := j.processor.Expand(jsonDoc, j.options)
	if err != nil {
		return fmt.Errorf("invalid JSON-LD: %w", err)
	}

	return nil
}

// ParseCroissantMetadata parses Croissant JSON-LD metadata and converts it to our Metadata struct
func (j *JSONLDProcessor) ParseCroissantMetadata(data []byte) (*Metadata, error) {
	// First validate that it's valid JSON-LD
	if err := j.ValidateJSONLD(data); err != nil {
		return nil, err
	}

	// Parse as regular JSON into our Metadata struct
	// The json-gold validation ensures it's properly structured JSON-LD
	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse Croissant metadata: %w", err)
	}

	return &metadata, nil
}

// GetExpandedProperty retrieves a property from expanded JSON-LD using its full IRI
func GetExpandedProperty(expanded map[string]interface{}, property string) interface{} {
	// Try direct property access first
	if val, exists := expanded[property]; exists {
		return val
	}

	// Try common prefixes
	prefixes := []string{
		"https://schema.org/",
		"http://mlcommons.org/croissant/",
		"http://purl.org/dc/terms/",
	}

	for _, prefix := range prefixes {
		fullProperty := prefix + strings.TrimPrefix(property, prefix)
		if val, exists := expanded[fullProperty]; exists {
			return val
		}
	}

	return nil
}

// GetPropertyValue extracts a simple string value from a JSON-LD property
func GetPropertyValue(property interface{}) string {
	if property == nil {
		return ""
	}

	// Handle array of values (common in JSON-LD)
	if arr, ok := property.([]interface{}); ok && len(arr) > 0 {
		if valueObj, ok := arr[0].(map[string]interface{}); ok {
			if value, ok := valueObj["@value"].(string); ok {
				return value
			}
		}
		// Fallback to string conversion
		if str, ok := arr[0].(string); ok {
			return str
		}
	}

	// Handle single value object
	if valueObj, ok := property.(map[string]interface{}); ok {
		if value, ok := valueObj["@value"].(string); ok {
			return value
		}
	}

	// Handle simple string
	if str, ok := property.(string); ok {
		return str
	}

	return ""
}

// GetPropertyArray extracts an array of values from a JSON-LD property
func GetPropertyArray(property interface{}) []interface{} {
	if property == nil {
		return nil
	}

	if arr, ok := property.([]interface{}); ok {
		return arr
	}

	// If it's a single value, return it as an array
	return []interface{}{property}
}

// ExtractCroissantProperties extracts common Croissant properties from expanded JSON-LD
func ExtractCroissantProperties(expanded map[string]interface{}) map[string]interface{} {
	properties := make(map[string]interface{})

	// Common Schema.org properties
	properties["name"] = GetPropertyValue(GetExpandedProperty(expanded, "name"))
	properties["description"] = GetPropertyValue(GetExpandedProperty(expanded, "description"))
	properties["url"] = GetPropertyValue(GetExpandedProperty(expanded, "url"))
	properties["version"] = GetPropertyValue(GetExpandedProperty(expanded, "version"))
	properties["datePublished"] = GetPropertyValue(GetExpandedProperty(expanded, "datePublished"))

	// Croissant-specific properties
	properties["conformsTo"] = GetPropertyValue(GetExpandedProperty(expanded, "conformsTo"))
	properties["distribution"] = GetPropertyArray(GetExpandedProperty(expanded, "distribution"))
	properties["recordSet"] = GetPropertyArray(GetExpandedProperty(expanded, "recordSet"))

	// Type information
	properties["@type"] = GetPropertyValue(GetExpandedProperty(expanded, "@type"))

	return properties
}

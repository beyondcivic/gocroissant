// structs.go
package croissant

import (
	"encoding/json"
	"fmt"
)

// Field represents a field in the Croissant metadata
type Field struct {
	ID          string      `json:"@id"`
	Type        string      `json:"@type"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	DataType    DataType    `json:"dataType"`
	Source      FieldSource `json:"source,omitempty"`
	Repeated    bool        `json:"repeated,omitempty"`
	Examples    interface{} `json:"examples,omitempty"`
	SubField    []Field     `json:"subField,omitempty"`
	References  *FieldRef   `json:"references,omitempty"`
}

// FieldSource represents the source information for a field
type FieldSource struct {
	Extract    Extract    `json:"extract,omitempty"`
	FileObject FileObject `json:"fileObject,omitempty"`
	FileSet    FileObject `json:"fileSet,omitempty"`
	Transform  *Transform `json:"transform,omitempty"`
	Format     string     `json:"format,omitempty"`
}

// Extract represents the extraction information for a field source
type Extract struct {
	Column       string `json:"column,omitempty"`
	JSONPath     string `json:"jsonPath,omitempty"`
	Regex        string `json:"regex,omitempty"`
	Separator    string `json:"separator,omitempty"`
	FileProperty string `json:"fileProperty,omitempty"`
}

// FileObject represents a file object reference
type FileObject struct {
	ID string `json:"@id"`
}

// KeyRef represents a key reference in a composite key
type KeyRef struct {
	ID string `json:"@id"`
}

// FieldRef represents a reference to another field
type FieldRef struct {
	ID    string  `json:"@id,omitempty"`
	Field *KeyRef `json:"field,omitempty"`
}

// DataType represents a data type that can be either a single string or an array of strings
type DataType struct {
	// Single dataType case: just a string value
	SingleType *string `json:"-"`
	// Array dataType case: array of string values
	ArrayType []string `json:"-"`
}

// RecordSetKey represents a record set key that can be either a single key or composite key
type RecordSetKey struct {
	// Single key case: just an ID reference
	SingleKey *KeyRef `json:"-"`
	// Composite key case: array of ID references
	CompositeKey []KeyRef `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for RecordSetKey
func (k RecordSetKey) MarshalJSON() ([]byte, error) {
	if k.SingleKey != nil {
		return json.Marshal(k.SingleKey)
	}
	if k.CompositeKey != nil {
		return json.Marshal(k.CompositeKey)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for RecordSetKey
func (k *RecordSetKey) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single key first
	var singleKey KeyRef
	if err := json.Unmarshal(data, &singleKey); err == nil && singleKey.ID != "" {
		k.SingleKey = &singleKey
		return nil
	}

	// Try to unmarshal as an array of keys
	var compositeKey []KeyRef
	if err := json.Unmarshal(data, &compositeKey); err == nil && len(compositeKey) > 0 {
		k.CompositeKey = compositeKey
		return nil
	}

	// Return error if neither format worked
	return fmt.Errorf("key must be either a single key object or an array of key objects")
}

// IsComposite returns true if this is a composite key
func (k RecordSetKey) IsComposite() bool {
	return len(k.CompositeKey) > 0
}

// GetKeyIDs returns all key IDs (single or composite)
func (k RecordSetKey) GetKeyIDs() []string {
	if k.SingleKey != nil {
		return []string{k.SingleKey.ID}
	}
	if k.CompositeKey != nil {
		ids := make([]string, len(k.CompositeKey))
		for i, key := range k.CompositeKey {
			ids[i] = key.ID
		}
		return ids
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for DataType
func (d DataType) MarshalJSON() ([]byte, error) {
	if d.SingleType != nil {
		return json.Marshal(*d.SingleType)
	}
	if d.ArrayType != nil {
		return json.Marshal(d.ArrayType)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for DataType
func (d *DataType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single string first
	var singleType string
	if err := json.Unmarshal(data, &singleType); err == nil && singleType != "" {
		d.SingleType = &singleType
		return nil
	}

	// Try to unmarshal as an array of strings
	var arrayType []string
	if err := json.Unmarshal(data, &arrayType); err == nil && len(arrayType) > 0 {
		d.ArrayType = arrayType
		return nil
	}

	// Return error if neither format worked
	return fmt.Errorf("dataType must be either a string or an array of strings")
}

// IsArray returns true if this is an array of data types
func (d DataType) IsArray() bool {
	return len(d.ArrayType) > 0
}

// GetTypes returns all data types (single or array)
func (d DataType) GetTypes() []string {
	if d.SingleType != nil {
		return []string{*d.SingleType}
	}
	if d.ArrayType != nil {
		return d.ArrayType
	}
	return nil
}

// GetFirstType returns the first data type (useful for backward compatibility)
func (d DataType) GetFirstType() string {
	if d.SingleType != nil {
		return *d.SingleType
	}
	if len(d.ArrayType) > 0 {
		return d.ArrayType[0]
	}
	return ""
}

// Distribution represents a file in the Croissant metadata
type Distribution struct {
	ID             string      `json:"@id"`
	Type           string      `json:"@type"`
	Name           string      `json:"name"`
	Description    string      `json:"description,omitempty"`
	ContentSize    string      `json:"contentSize,omitempty"`
	ContentURL     string      `json:"contentUrl,omitempty"`
	EncodingFormat string      `json:"encodingFormat"`
	SHA256         string      `json:"sha256,omitempty"`
	MD5            string      `json:"md5,omitempty"`
	ContainedIn    *FileObject `json:"containedIn,omitempty"`
	Includes       string      `json:"includes,omitempty"`
}

// RecordSet represents a record set in the Croissant metadata
type RecordSet struct {
	ID          string                   `json:"@id"`
	Type        string                   `json:"@type"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	DataType    *DataType                `json:"dataType,omitempty"`
	Fields      []Field                  `json:"field"`
	Key         *RecordSetKey            `json:"key,omitempty"`
	Data        []map[string]interface{} `json:"data,omitempty"`
}

// Context represents the complete JSON-LD context for Croissant 1.0
type Context struct {
	Language      string          `json:"@language"`
	Vocab         string          `json:"@vocab"`
	CiteAs        string          `json:"citeAs"`
	Column        string          `json:"column"`
	ConformsTo    string          `json:"conformsTo"`
	CR            string          `json:"cr"`
	DCT           string          `json:"dct"`
	RAI           string          `json:"rai,omitempty"`
	WD            string          `json:"wd,omitempty"`
	Data          DataContext     `json:"data"`
	DataType      DataTypeContext `json:"dataType"`
	Examples      DataContext     `json:"examples"`
	Extract       string          `json:"extract"`
	Field         string          `json:"field"`
	FileObject    string          `json:"fileObject"`
	FileProperty  string          `json:"fileProperty"`
	FileSet       string          `json:"fileSet"`
	Format        string          `json:"format"`
	Includes      string          `json:"includes"`
	IsLiveDataset string          `json:"isLiveDataset"`
	JSONPath      string          `json:"jsonPath"`
	Key           string          `json:"key"`
	MD5           string          `json:"md5"`
	ParentField   string          `json:"parentField"`
	Path          string          `json:"path"`
	RecordSet     string          `json:"recordSet"`
	References    string          `json:"references"`
	Regex         string          `json:"regex"`
	Repeated      string          `json:"repeated"`
	Replace       string          `json:"replace"`
	SC            string          `json:"sc"`
	Separator     string          `json:"separator"`
	Source        string          `json:"source"`
	SubField      string          `json:"subField"`
	Transform     string          `json:"transform"`
}

// DataContext represents the data field in the context
type DataContext struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

// DataTypeContext represents the dataType field in the context
type DataTypeContext struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

// Metadata represents the complete Croissant metadata
type Metadata struct {
	Context       Context        `json:"@context"`
	Type          string         `json:"@type"`
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	ConformsTo    string         `json:"conformsTo"`
	DatePublished string         `json:"datePublished,omitempty"`
	Version       string         `json:"version,omitempty"`
	URL           string         `json:"url,omitempty"`
	License       string         `json:"license,omitempty"`
	CiteAs        string         `json:"citeAs,omitempty"`
	Creator       interface{}    `json:"creator,omitempty"`
	Publisher     interface{}    `json:"publisher,omitempty"`
	Keywords      []string       `json:"keywords,omitempty"`
	Distributions []Distribution `json:"distribution"`
	RecordSets    []RecordSet    `json:"recordSet"`
	IsLiveDataset bool           `json:"isLiveDataset,omitempty"`
}

// Transform represents a data transformation
type Transform struct {
	Type      string `json:"@type"`
	Regex     string `json:"regex,omitempty"`
	Replace   string `json:"replace,omitempty"`
	Format    string `json:"format,omitempty"`
	JSONPath  string `json:"jsonPath,omitempty"`
	Separator string `json:"separator,omitempty"`
}

// Source represents a more complete source definition
type Source struct {
	Extract    Extract     `json:"extract,omitempty"`
	FileObject FileObject  `json:"fileObject,omitempty"`
	Field      string      `json:"field,omitempty"`
	Transform  []Transform `json:"transform,omitempty"`
}

// NewSingleKey creates a RecordSetKey with a single key reference
func NewSingleKey(keyID string) *RecordSetKey {
	return &RecordSetKey{
		SingleKey: &KeyRef{ID: keyID},
	}
}

// NewCompositeKey creates a RecordSetKey with multiple key references
func NewCompositeKey(keyIDs ...string) *RecordSetKey {
	keys := make([]KeyRef, len(keyIDs))
	for i, id := range keyIDs {
		keys[i] = KeyRef{ID: id}
	}
	return &RecordSetKey{
		CompositeKey: keys,
	}
}

// NewSingleDataType creates a DataType with a single type
func NewNullableSingleDataType(dataType string) *DataType {
	return &DataType{
		SingleType: &dataType,
	}
}

// NewSingleDataType creates a DataType with a single type
func NewSingleDataType(dataType string) DataType {
	return DataType{
		SingleType: &dataType,
	}
}

// NewArrayDataType creates a DataType with multiple types
func NewArrayDataType(dataTypes ...string) DataType {
	return DataType{
		ArrayType: dataTypes,
	}
}

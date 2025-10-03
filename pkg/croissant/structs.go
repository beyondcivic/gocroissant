// structs.go
package croissant

import (
	"encoding/json"
)

// Field represents a field in the Croissant metadata.
type Field struct {
	ID          string        `json:"@id"`
	Type        string        `json:"@type"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	DataType    DataType      `json:"dataType"`
	Source      FieldSource   `json:"source,omitempty"`
	Repeated    bool          `json:"repeated,omitempty"`
	Examples    interface{}   `json:"examples,omitempty"`
	SubField    []Field       `json:"subField,omitempty"`
	ParentField FieldRefSlice `json:"parentField,omitempty"`
	References  FieldRefSlice `json:"references,omitempty"`
}

// FieldSource represents the source information for a field.
type FieldSource struct {
	Extract    *Extract    `json:"extract,omitempty"`
	FileObject *FileObject `json:"fileObject,omitempty"`
	FileSet    *FileObject `json:"fileSet,omitempty"`
	Transform  *Transform  `json:"transform,omitempty"`
	Format     string      `json:"format,omitempty"`
}

// Extract represents the extraction information for a field source.
type Extract struct {
	// Extraction method
	FileProperty string `json:"fileProperty,omitempty"`
	// Name of the column (field) that contains values.
	Column string `json:"column,omitempty"`
	// A JSONPATH expression that extracts values.
	JSONPath  string `json:"jsonPath,omitempty"`
	Regex     string `json:"regex,omitempty"`
	Separator string `json:"separator,omitempty"`
}

// FileObject represents a file object reference.
type FileObject struct {
	ID string `json:"@id"`
}

// KeyRef represents a key reference in a composite key.
type KeyRef struct {
	ID string `json:"@id"`
}

// FieldRef represents a reference to another field.
type FieldRef struct {
	ID    string  `json:"@id,omitempty"`
	Field *KeyRef `json:"field,omitempty"`
}

// Parses ONE or MANY FieldRefs.
type FieldRefSlice []FieldRef

// In some test files, references are nested under a "field" property.
// In cases of reformatting, the property will be omitted.
//
// Accepts:
//   - "references": { "@id": "..." }
//   - "references": { [{"@id": "..."}, {"@id": "..."}...] }
//   - "references": { field: {"@id": "..."} }
func (ref *FieldRefSlice) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as a NestedFieldRef first
	type NestedFieldRef struct {
		Field *FieldRef `json:"field,omitempty"`
	}
	var singleNested NestedFieldRef
	if err := json.Unmarshal(data, &singleNested); err == nil {
		*ref = []FieldRef{*singleNested.Field}

		return nil
	}

	var single FieldRef
	if err := json.Unmarshal(data, &single); err == nil {
		*ref = []FieldRef{single}

		return nil
	}
	// Try unmarshaling as a []FieldRef
	var multi []FieldRef
	if err := json.Unmarshal(data, &multi); err == nil {
		*ref = multi

		return nil
	}
	// Otherwise, error
	return CroissantError{
		Message: "FieldRefSlice: cannot unmarshal",
		Value:   string(data),
	}
}

func (ref FieldRefSlice) MarshalJSON() ([]byte, error) {
	switch len(ref) {
	case 0:
		return []byte("{}"), nil
	case 1:
		return json.Marshal(ref[0])
	default:
		return json.Marshal([]FieldRef(ref))
	}
}

// DataType represents a data type that can be either a single string or an array of strings.
// It is represented internally as a list.
type DataType []string

// RecordSetKey represents a record set key that can be either a single key or composite key
type RecordSetKey []KeyRef

// MarshalJSON implements custom JSON marshaling for RecordSetKey
func (key RecordSetKey) MarshalJSON() ([]byte, error) {
	switch len(key) {
	case 0:
		return []byte("{}"), nil
	case 1:
		return json.Marshal(key[0])
	default:
		return json.Marshal(key)
	}
}

// UnmarshalJSON implements custom JSON unmarshaling for RecordSetKey
func (key *RecordSetKey) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single key first
	var singleKey KeyRef
	if err := json.Unmarshal(data, &singleKey); err == nil && singleKey.ID != "" {
		*key = []KeyRef{singleKey}
		return nil
	}

	// Try to unmarshal as an array of keys
	var compositeKey []KeyRef
	if err := json.Unmarshal(data, &compositeKey); err == nil {
		*key = compositeKey
		return nil
	}

	// Return error if neither format worked
	return CroissantError{Message: "key must be either a single key object or an array of key objects"}
}

// IsComposite returns true if this is a composite key.
func (k RecordSetKey) IsComposite() bool {
	return len(k) > 1
}

// GetKeyIDs returns all key IDs (single or composite).
func (k RecordSetKey) GetKeyIDs() []string {
	if k == nil {
		return nil
	}
	ids := make([]string, len(k))
	for i, key := range k {
		ids[i] = key.ID
	}
	return ids
}

// MarshalJSON implements custom JSON marshaling for DataType.
func (d DataType) MarshalJSON() ([]byte, error) {
	switch len(d) {
	case 0:
		return []byte("{}"), nil
	case 1:
		return json.Marshal(d[0])
	default:
		return json.Marshal(d)
	}
}

// UnmarshalJSON implements custom JSON unmarshaling for DataType.
func (d *DataType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single string first
	var singleType string
	if err := json.Unmarshal(data, &singleType); err == nil && singleType != "" {
		*d = []string{singleType}

		return nil
	}

	// Try to unmarshal as an array of strings
	var arrayType []string
	if err := json.Unmarshal(data, &arrayType); err == nil && len(arrayType) > 0 {
		*d = arrayType

		return nil
	}

	// Return error if neither format worked
	return CroissantError{Message: "dataType must be either a string or an array of strings"}
}

// IsArray returns true if this is an array of data types.
func (d DataType) IsArray() bool {
	return len(d) > 1
}

// GetTypes returns all data types (single or array).
func (d DataType) GetTypes() []string {
	return d
}

// GetFirstType returns the first data type (useful for backward compatibility).
func (d DataType) GetFirstType() string {
	if len(d) > 0 {
		return d[0]
	}

	return ""
}

// Distribution represents a file in the Croissant metadata.
type Distribution struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
	// The name of the file.
	Name string `json:"name"`
	// Description of the file.
	Description string `json:"description,omitempty"`
	// File size in kb, mb, gb etc...
	// Defaults to bytes if unit not specified.
	ContentSize string `json:"contentSize,omitempty"`
	// URL to the actual bytes of the file object.
	ContentURL string `json:"contentUrl,omitempty"`
	// Format of the file, given as a MIME type.
	EncodingFormat string `json:"encodingFormat"`
	// SHA256 checksum of the file contents.
	SHA256 string `json:"sha256,omitempty"`
	// MD5 checksum of the file contents.
	MD5 string `json:"md5,omitempty"`
	// A FileObject or FileSet this resource is contained in.
	ContainedIn *FileObjectRef `json:"containedIn,omitempty"`
	// A glob pattern of the files to include (FileSet).
	Includes string `json:"includes,omitempty"`
	// A glob pattern of the files to exclude (FileSet).
	Excludes string `json:"excludes,omitempty"`
}

// RecordSet represents a record set in the Croissant metadata.
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

// Context represents the complete JSON-LD context for Croissant 1.0.
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

// DataContext represents the data field in the context.
type DataContext struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

// DataTypeContext represents the dataType field in the context.
type DataTypeContext struct {
	ID   string `json:"@id"`
	Type string `json:"@type"`
}

// Metadata represents the complete Croissant metadata for a dataset.
type Metadata struct {
	Context Context `json:"@context"`
	// Dataset Type.  Must by `schema.org/Dataset`
	Type string `json:"@type"`
	// Name of the dataset.
	Name string `json:"name"`
	// Description of the dataset.
	Description string `json:"description,omitempty"`
	// Versioned schema the croissant metadata conforms to.
	ConformsTo string `json:"conformsTo"`
	// Date the dataset was published.
	DatePublished string `json:"datePublished,omitempty"`
	// Version of the dataset.
	// Either an single int, or a MAJOR.MINOR.PATCH sematic version string.
	// [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
	Version string `json:"version,omitempty"`
	// Url of the dataset, usually a webpage.
	URL string `json:"url,omitempty"`
	// Licenses of the dataset.
	// Spec suggests using references from https://spdx.org/licenses/.
	License string `json:"license,omitempty"`
	// A citation to the dataset itself, or a citation for a publication that describes the dataset.
	// Ideally, citations should be expressed using the bibtex format.
	// Note that this is different from schema.org/citation, which is used to make a citation to another publication from this dataset.
	CiteAs string `json:"citeAs,omitempty"`
	// Creator(s) of the dataset.
	Creator interface{} `json:"creator,omitempty"`
	// Publisher(s) of the dataset.
	Publisher interface{} `json:"publisher,omitempty"`
	// A set of keywords associated with the dataset, either as free text, or a DefinedTerm with a formal definition.
	Keywords []string `json:"keywords,omitempty"`
	// Set of FileObject and FileSet definitions that describe the raw files of the dataset.
	Distributions []Distribution `json:"distribution"`
	// Set of RecordSet definitions that describe the content of the dataset.
	RecordSets []RecordSet `json:"recordSet"`
	// If true, dataset is non-static and may change over time.
	// Distribution resources may not contain a checksum if they are expected to change.
	IsLiveDataset bool `json:"isLiveDataset,omitempty"`
}

// Transform represents a data transformation.
type Transform struct {
	Type      string `json:"@type"`
	Regex     string `json:"regex,omitempty"`
	Replace   string `json:"replace,omitempty"`
	Format    string `json:"format,omitempty"`
	JSONPath  string `json:"jsonPath,omitempty"`
	Separator string `json:"separator,omitempty"`
}

// Source represents a more complete source definition.
type Source struct {
	Extract    *Extract    `json:"extract,omitempty"`
	FileObject *FileObject `json:"fileObject,omitempty"`
	Field      string      `json:"field,omitempty"`
	Transform  []Transform `json:"transform,omitempty"`
}

// NewRecordSetKey creates a RecordSetKey with a single key reference
func NewRecordSetKey(keyID string) *RecordSetKey {
	return &RecordSetKey{
		KeyRef{ID: keyID},
	}
}

// NewCompositeKey creates a RecordSetKey with multiple key references.
func NewCompositeKey(keyIDs ...string) *RecordSetKey {
	keys := make(RecordSetKey, len(keyIDs))
	for i, id := range keyIDs {
		keys[i] = KeyRef{ID: id}
	}
	return &keys
}

// NewSingleDataType creates a DataType with a single type.
func NewNullableSingleDataType(dataType string) *DataType {
	return &DataType{dataType}
}

// NewSingleDataType creates a DataType with a single type.
func NewSingleDataType(dataType string) DataType {
	return DataType{dataType}
}

// NewArrayDataType creates a DataType with multiple types.
func NewArrayDataType(dataTypes ...string) DataType {
	return dataTypes
}

// ValidateSource validates the source configuration.
func (fs FieldSource) ValidateSource() bool {
	// If no source is configured, it's invalid unless it's a parent field with subfields
	hasFileObject := fs.FileObject.ID != "" || fs.FileSet.ID != ""
	hasExtract := fs.Extract.Column != "" || fs.Extract.JSONPath != "" || fs.Extract.FileProperty != "" || fs.Extract.Regex != ""

	// A valid source needs either a file object reference with extraction info, or other valid configurations
	return hasFileObject && (hasExtract || fs.Format != "")
}

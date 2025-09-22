// structs.go
package croissant

// Field represents a field in the Croissant metadata
type Field struct {
	ID          string      `json:"@id"`
	Type        string      `json:"@type"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	DataType    string      `json:"dataType"`
	Source      FieldSource `json:"source"`
	Repeated    bool        `json:"repeated,omitempty"`
	Examples    interface{} `json:"examples,omitempty"`
}

// FieldSource represents the source information for a field
type FieldSource struct {
	Extract    Extract    `json:"extract"`
	FileObject FileObject `json:"fileObject"`
}

// Extract represents the extraction information for a field source
type Extract struct {
	Column    string `json:"column,omitempty"`
	JSONPath  string `json:"jsonPath,omitempty"`
	Regex     string `json:"regex,omitempty"`
	Separator string `json:"separator,omitempty"`
}

// FileObject represents a file object reference
type FileObject struct {
	ID string `json:"@id"`
}

// Distribution represents a file in the Croissant metadata
type Distribution struct {
	ID             string `json:"@id"`
	Type           string `json:"@type"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	ContentSize    string `json:"contentSize,omitempty"`
	ContentURL     string `json:"contentUrl"`
	EncodingFormat string `json:"encodingFormat"`
	SHA256         string `json:"sha256,omitempty"`
	MD5            string `json:"md5,omitempty"`
}

// RecordSet represents a record set in the Croissant metadata
type RecordSet struct {
	ID          string  `json:"@id"`
	Type        string  `json:"@type"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Fields      []Field `json:"field"`
	Key         string  `json:"key,omitempty"`
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

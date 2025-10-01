# croissant

```go
import "github.com/beyondcivic/gocroissant/pkg/croissant"
```

Package croissant provides functionality for working with the ML Commons Croissant metadata format \- a standardized way to describe machine learning datasets using JSON\-LD.

This package simplifies the creation of Croissant\-compatible metadata from CSV data sources by automatically inferring schema types from dataset content, generating complete valid JSON\-LD metadata, providing validation tools to ensure compatibility, and supporting the full Croissant specification.

### Basic Usage

Generate metadata from a CSV file:

```
outputPath, err := croissant.GenerateMetadata("data.csv", "dataset.jsonld")
if err != nil {
	log.Fatalf("Error generating metadata: %v", err)
}
fmt.Printf("Metadata saved to: %s\n", outputPath)
```

### Data Type Inference

The package automatically infers schema.org data types from CSV content:

- Boolean values \(true/false\)
- Integer numbers
- Floating\-point numbers
- Dates in various formats
- URLs
- Default to Text for other content

### Validation

Validate existing Croissant metadata:

```
issues, err := croissant.ValidateMetadata("metadata.jsonld")
if err != nil {
	log.Fatalf("Validation error: %v", err)
}
if len(issues) == 0 {
	fmt.Println("Validation passed")
}
```

issues.go

jsonld.go

metadata\_node.go

node.go

structs.go

utils.go

validation.go

## Index

- [func CalculateSHA256\(filePath string\) \(string, error\)](<#CalculateSHA256>)
- [func CountCSVRows\(csvPath string\) \(int, error\)](<#CountCSVRows>)
- [func DetectCSVDelimiter\(csvPath string\) \(rune, error\)](<#DetectCSVDelimiter>)
- [func ExtractCroissantProperties\(expanded map\[string\]interface\{\}\) map\[string\]interface\{\}](<#ExtractCroissantProperties>)
- [func GenerateMetadata\(csvPath string, outputPath string\) \(string, error\)](<#GenerateMetadata>)
- [func GetCSVColumnTypes\(csvPath string, sampleSize int\) \(\[\]string, \[\]string, error\)](<#GetCSVColumnTypes>)
- [func GetCSVColumns\(csvPath string\) \(\[\]string, \[\]string, error\)](<#GetCSVColumns>)
- [func GetCSVColumnsAndSampleRows\(csvPath string, maxRows int\) \(\[\]string, \[\]\[\]string, error\)](<#GetCSVColumnsAndSampleRows>)
- [func GetExpandedProperty\(expanded map\[string\]interface\{\}, property string\) interface\{\}](<#GetExpandedProperty>)
- [func GetFileStats\(filePath string\) \(map\[string\]interface\{\}, error\)](<#GetFileStats>)
- [func GetPropertyArray\(property interface\{\}\) \[\]interface\{\}](<#GetPropertyArray>)
- [func GetPropertyValue\(property interface\{\}\) string](<#GetPropertyValue>)
- [func InferDataType\(value string\) string](<#InferDataType>)
- [func InferSemanticDataType\(fieldName, value string, context map\[string\]interface\{\}\) \[\]string](<#InferSemanticDataType>)
- [func IsCSVFile\(filePath string\) bool](<#IsCSVFile>)
- [func IsValidDataType\(dataType string\) bool](<#IsValidDataType>)
- [func ParseCSVWithOptions\(csvPath string, delimiter rune, hasHeader bool\) \(\[\]string, \[\]\[\]string, error\)](<#ParseCSVWithOptions>)
- [func SanitizeFileName\(fileName string\) string](<#SanitizeFileName>)
- [func ValidateCSVStructure\(csvPath string\) error](<#ValidateCSVStructure>)
- [func ValidateCrossReferences\(node \*MetadataNode, issues \*Issues\)](<#ValidateCrossReferences>)
- [func ValidateDistributionNode\(dist \*DistributionNode, issues \*Issues, options ValidationOptions\)](<#ValidateDistributionNode>)
- [func ValidateFieldNode\(field \*FieldNode, issues \*Issues, options ValidationOptions\)](<#ValidateFieldNode>)
- [func ValidateMetadataNode\(node \*MetadataNode, issues \*Issues, options ValidationOptions\)](<#ValidateMetadataNode>)
- [func ValidateOutputPath\(outputPath string\) error](<#ValidateOutputPath>)
- [func ValidateRecordSetNode\(rs \*RecordSetNode, issues \*Issues, options ValidationOptions\)](<#ValidateRecordSetNode>)
- [type BaseNode](<#BaseNode>)
  - [func \(n \*BaseNode\) GetID\(\) string](<#BaseNode.GetID>)
  - [func \(n \*BaseNode\) GetName\(\) string](<#BaseNode.GetName>)
  - [func \(n \*BaseNode\) GetParent\(\) Node](<#BaseNode.GetParent>)
  - [func \(n \*BaseNode\) SetParent\(parent Node\)](<#BaseNode.SetParent>)
- [type Context](<#Context>)
  - [func CreateDefaultContext\(\) Context](<#CreateDefaultContext>)
- [type CroissantError](<#CroissantError>)
  - [func \(e CroissantError\) Error\(\) string](<#CroissantError.Error>)
- [type DataContext](<#DataContext>)
- [type DataType](<#DataType>)
  - [func NewArrayDataType\(dataTypes ...string\) DataType](<#NewArrayDataType>)
  - [func NewNullableSingleDataType\(dataType string\) \*DataType](<#NewNullableSingleDataType>)
  - [func NewSingleDataType\(dataType string\) DataType](<#NewSingleDataType>)
  - [func \(d DataType\) GetFirstType\(\) string](<#DataType.GetFirstType>)
  - [func \(d DataType\) GetTypes\(\) \[\]string](<#DataType.GetTypes>)
  - [func \(d DataType\) IsArray\(\) bool](<#DataType.IsArray>)
  - [func \(d DataType\) MarshalJSON\(\) \(\[\]byte, error\)](<#DataType.MarshalJSON>)
  - [func \(d \*DataType\) UnmarshalJSON\(data \[\]byte\) error](<#DataType.UnmarshalJSON>)
- [type DataTypeContext](<#DataTypeContext>)
- [type Distribution](<#Distribution>)
- [type DistributionNode](<#DistributionNode>)
  - [func \(d \*DistributionNode\) Validate\(issues \*Issues\)](<#DistributionNode.Validate>)
- [type Extract](<#Extract>)
- [type ExtractNode](<#ExtractNode>)
- [type Field](<#Field>)
- [type FieldNode](<#FieldNode>)
  - [func \(f \*FieldNode\) Validate\(issues \*Issues\)](<#FieldNode.Validate>)
- [type FieldRef](<#FieldRef>)
- [type FieldRefSlice](<#FieldRefSlice>)
  - [func \(ref FieldRefSlice\) MarshalJSON\(\) \(\[\]byte, error\)](<#FieldRefSlice.MarshalJSON>)
  - [func \(ref \*FieldRefSlice\) UnmarshalJSON\(data \[\]byte\) error](<#FieldRefSlice.UnmarshalJSON>)
- [type FieldSource](<#FieldSource>)
  - [func \(fs FieldSource\) ValidateSource\(\) bool](<#FieldSource.ValidateSource>)
- [type FileObject](<#FileObject>)
- [type FileObjectRef](<#FileObjectRef>)
- [type Issue](<#Issue>)
- [type IssueType](<#IssueType>)
- [type Issues](<#Issues>)
  - [func NewIssues\(\) \*Issues](<#NewIssues>)
  - [func ValidateFile\(filePath string\) \(\*Issues, error\)](<#ValidateFile>)
  - [func ValidateJSON\(data \[\]byte\) \(\*Issues, error\)](<#ValidateJSON>)
  - [func ValidateJSONWithOptions\(data \[\]byte, options ValidationOptions\) \(\*Issues, error\)](<#ValidateJSONWithOptions>)
  - [func ValidateMetadata\(metadata Metadata\) \*Issues](<#ValidateMetadata>)
  - [func ValidateMetadataWithOptions\(metadata Metadata, options ValidationOptions\) \*Issues](<#ValidateMetadataWithOptions>)
  - [func \(i \*Issues\) AddError\(message string, node ...Node\)](<#Issues.AddError>)
  - [func \(i \*Issues\) AddWarning\(message string, node ...Node\)](<#Issues.AddWarning>)
  - [func \(i \*Issues\) ErrorCount\(\) int](<#Issues.ErrorCount>)
  - [func \(i \*Issues\) HasErrors\(\) bool](<#Issues.HasErrors>)
  - [func \(i \*Issues\) HasWarnings\(\) bool](<#Issues.HasWarnings>)
  - [func \(i \*Issues\) Report\(\) string](<#Issues.Report>)
  - [func \(i \*Issues\) WarningCount\(\) int](<#Issues.WarningCount>)
- [type JSONLDProcessor](<#JSONLDProcessor>)
  - [func NewJSONLDProcessor\(\) \*JSONLDProcessor](<#NewJSONLDProcessor>)
  - [func \(j \*JSONLDProcessor\) CompactJSONLD\(expanded interface\{\}, context map\[string\]interface\{\}\) \(map\[string\]interface\{\}, error\)](<#JSONLDProcessor.CompactJSONLD>)
  - [func \(j \*JSONLDProcessor\) ParseCroissantMetadata\(data \[\]byte\) \(\*Metadata, error\)](<#JSONLDProcessor.ParseCroissantMetadata>)
  - [func \(j \*JSONLDProcessor\) ParseJSONLD\(data \[\]byte\) \(map\[string\]interface\{\}, error\)](<#JSONLDProcessor.ParseJSONLD>)
  - [func \(j \*JSONLDProcessor\) ValidateJSONLD\(data \[\]byte\) error](<#JSONLDProcessor.ValidateJSONLD>)
- [type KeyRef](<#KeyRef>)
- [type Metadata](<#Metadata>)
- [type MetadataNode](<#MetadataNode>)
  - [func FromMetadata\(metadata Metadata\) \*MetadataNode](<#FromMetadata>)
  - [func NewMetadataNode\(\) \*MetadataNode](<#NewMetadataNode>)
  - [func \(m \*MetadataNode\) Validate\(issues \*Issues\)](<#MetadataNode.Validate>)
- [type MetadataWithValidation](<#MetadataWithValidation>)
  - [func GenerateMetadataWithValidation\(csvPath string, outputPath string\) \(\*MetadataWithValidation, error\)](<#GenerateMetadataWithValidation>)
  - [func NewMetadataWithValidation\(metadata Metadata\) \*MetadataWithValidation](<#NewMetadataWithValidation>)
  - [func \(m \*MetadataWithValidation\) GetIssues\(\) \*Issues](<#MetadataWithValidation.GetIssues>)
  - [func \(m \*MetadataWithValidation\) HasErrors\(\) bool](<#MetadataWithValidation.HasErrors>)
  - [func \(m \*MetadataWithValidation\) HasWarnings\(\) bool](<#MetadataWithValidation.HasWarnings>)
  - [func \(m \*MetadataWithValidation\) Report\(\) string](<#MetadataWithValidation.Report>)
  - [func \(m \*MetadataWithValidation\) Validate\(\)](<#MetadataWithValidation.Validate>)
  - [func \(m \*MetadataWithValidation\) ValidateWithOptions\(options ValidationOptions\)](<#MetadataWithValidation.ValidateWithOptions>)
- [type Node](<#Node>)
- [type RecordSet](<#RecordSet>)
  - [func CreateEnumerationRecordSet\(id, name string, values \[\]string, urls \[\]string\) RecordSet](<#CreateEnumerationRecordSet>)
  - [func CreateSplitRecordSet\(\) RecordSet](<#CreateSplitRecordSet>)
- [type RecordSetKey](<#RecordSetKey>)
  - [func NewCompositeKey\(keyIDs ...string\) \*RecordSetKey](<#NewCompositeKey>)
  - [func NewSingleKey\(keyID string\) \*RecordSetKey](<#NewSingleKey>)
  - [func \(k RecordSetKey\) GetKeyIDs\(\) \[\]string](<#RecordSetKey.GetKeyIDs>)
  - [func \(k RecordSetKey\) IsComposite\(\) bool](<#RecordSetKey.IsComposite>)
  - [func \(k RecordSetKey\) MarshalJSON\(\) \(\[\]byte, error\)](<#RecordSetKey.MarshalJSON>)
  - [func \(k \*RecordSetKey\) UnmarshalJSON\(data \[\]byte\) error](<#RecordSetKey.UnmarshalJSON>)
- [type RecordSetNode](<#RecordSetNode>)
  - [func \(r \*RecordSetNode\) Validate\(issues \*Issues\)](<#RecordSetNode.Validate>)
- [type Source](<#Source>)
- [type SourceNode](<#SourceNode>)
  - [func \(s \*SourceNode\) ValidateSource\(\) bool](<#SourceNode.ValidateSource>)
- [type Transform](<#Transform>)
- [type ValidationOptions](<#ValidationOptions>)
  - [func DefaultValidationOptions\(\) ValidationOptions](<#DefaultValidationOptions>)


<a name="CalculateSHA256"></a>
## func [CalculateSHA256](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L17>)

```go
func CalculateSHA256(filePath string) (string, error)
```

CalculateSHA256 calculates the SHA\-256 hash of a file

<a name="CountCSVRows"></a>
## func [CountCSVRows](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L260>)

```go
func CountCSVRows(csvPath string) (int, error)
```

CountCSVRows counts the total number of rows in a CSV file \(including header\)

<a name="DetectCSVDelimiter"></a>
## func [DetectCSVDelimiter](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L148>)

```go
func DetectCSVDelimiter(csvPath string) (rune, error)
```

DetectCSVDelimiter attempts to detect the CSV delimiter

<a name="ExtractCroissantProperties"></a>
## func [ExtractCroissantProperties](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L169>)

```go
func ExtractCroissantProperties(expanded map[string]interface{}) map[string]interface{}
```

ExtractCroissantProperties extracts common Croissant properties from expanded JSON\-LD

<a name="GenerateMetadata"></a>
## func [GenerateMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L303>)

```go
func GenerateMetadata(csvPath string, outputPath string) (string, error)
```

GenerateMetadata generates Croissant metadata from a CSV file \(simple API\)

<a name="GetCSVColumnTypes"></a>
## func [GetCSVColumnTypes](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L345>)

```go
func GetCSVColumnTypes(csvPath string, sampleSize int) ([]string, []string, error)
```

GetCSVColumnTypes analyzes a CSV file and returns inferred data types for each column

<a name="GetCSVColumns"></a>
## func [GetCSVColumns](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L33>)

```go
func GetCSVColumns(csvPath string) ([]string, []string, error)
```

GetCSVColumns reads the column names and first row from a CSV file

<a name="GetCSVColumnsAndSampleRows"></a>
## func [GetCSVColumnsAndSampleRows](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L69>)

```go
func GetCSVColumnsAndSampleRows(csvPath string, maxRows int) ([]string, [][]string, error)
```

GetCSVColumnsAndSampleRows reads column names and multiple sample rows for better type inference

<a name="GetExpandedProperty"></a>
## func [GetExpandedProperty](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L97>)

```go
func GetExpandedProperty(expanded map[string]interface{}, property string) interface{}
```

GetExpandedProperty retrieves a property from expanded JSON\-LD using its full IRI

<a name="GetFileStats"></a>
## func [GetFileStats](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L240>)

```go
func GetFileStats(filePath string) (map[string]interface{}, error)
```

GetFileStats returns basic statistics about a file

<a name="GetPropertyArray"></a>
## func [GetPropertyArray](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L155>)

```go
func GetPropertyArray(property interface{}) []interface{}
```

GetPropertyArray extracts an array of values from a JSON\-LD property

<a name="GetPropertyValue"></a>
## func [GetPropertyValue](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L121>)

```go
func GetPropertyValue(property interface{}) string
```

GetPropertyValue extracts a simple string value from a JSON\-LD property

<a name="InferDataType"></a>
## func [InferDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L55>)

```go
func InferDataType(value string) string
```

InferDataType infers the schema.org data type from a value

<a name="InferSemanticDataType"></a>
## func [InferSemanticDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L145>)

```go
func InferSemanticDataType(fieldName, value string, context map[string]interface{}) []string
```

InferSemanticDataType attempts to infer semantic data types for ML datasets

<a name="IsCSVFile"></a>
## func [IsCSVFile](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L402>)

```go
func IsCSVFile(filePath string) bool
```

IsCSVFile checks if a file appears to be a CSV file based on extension

<a name="IsValidDataType"></a>
## func [IsValidDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L109>)

```go
func IsValidDataType(dataType string) bool
```

IsValidDataType checks if a dataType is valid according to Croissant specification

<a name="ParseCSVWithOptions"></a>
## func [ParseCSVWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L186>)

```go
func ParseCSVWithOptions(csvPath string, delimiter rune, hasHeader bool) ([]string, [][]string, error)
```

ParseCSVWithOptions parses a CSV file with custom options

<a name="SanitizeFileName"></a>
## func [SanitizeFileName](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L408>)

```go
func SanitizeFileName(fileName string) string
```

SanitizeFileName removes or replaces invalid characters in filenames

<a name="ValidateCSVStructure"></a>
## func [ValidateCSVStructure](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L292>)

```go
func ValidateCSVStructure(csvPath string) error
```

ValidateCSVStructure performs basic validation on CSV file structure

<a name="ValidateCrossReferences"></a>
## func [ValidateCrossReferences](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L323>)

```go
func ValidateCrossReferences(node *MetadataNode, issues *Issues)
```

ValidateCrossReferences validates that all references are valid

<a name="ValidateDistributionNode"></a>
## func [ValidateDistributionNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L149>)

```go
func ValidateDistributionNode(dist *DistributionNode, issues *Issues, options ValidationOptions)
```

ValidateDistributionNode validates a distribution node

<a name="ValidateFieldNode"></a>
## func [ValidateFieldNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L251>)

```go
func ValidateFieldNode(field *FieldNode, issues *Issues, options ValidationOptions)
```

ValidateFieldNode validates a field node

<a name="ValidateMetadataNode"></a>
## func [ValidateMetadataNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L95>)

```go
func ValidateMetadataNode(node *MetadataNode, issues *Issues, options ValidationOptions)
```

ValidateMetadataNode performs comprehensive validation of a metadata node

<a name="ValidateOutputPath"></a>
## func [ValidateOutputPath](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/utils.go#L118>)

```go
func ValidateOutputPath(outputPath string) error
```

ValidateOutputPath validates if the given path is a valid file path

<a name="ValidateRecordSetNode"></a>
## func [ValidateRecordSetNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L191>)

```go
func ValidateRecordSetNode(rs *RecordSetNode, issues *Issues, options ValidationOptions)
```

ValidateRecordSetNode validates a record set node

<a name="BaseNode"></a>
## type [BaseNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L14-L18>)

BaseNode implements common functionality for all nodes

```go
type BaseNode struct {
    ID   string `json:"@id,omitempty"`
    Name string `json:"name,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="BaseNode.GetID"></a>
### func \(\*BaseNode\) [GetID](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L24>)

```go
func (n *BaseNode) GetID() string
```



<a name="BaseNode.GetName"></a>
### func \(\*BaseNode\) [GetName](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L20>)

```go
func (n *BaseNode) GetName() string
```



<a name="BaseNode.GetParent"></a>
### func \(\*BaseNode\) [GetParent](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L28>)

```go
func (n *BaseNode) GetParent() Node
```



<a name="BaseNode.SetParent"></a>
### func \(\*BaseNode\) [SetParent](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L32>)

```go
func (n *BaseNode) SetParent(parent Node)
```



<a name="Context"></a>
## type [Context](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L257-L293>)

Context represents the complete JSON\-LD context for Croissant 1.0

```go
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
```

<a name="CreateDefaultContext"></a>
### func [CreateDefaultContext](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L254>)

```go
func CreateDefaultContext() Context
```

CreateDefaultContext creates the ML Commons Croissant 1.0 compliant context

<a name="CroissantError"></a>
## type [CroissantError](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/error.go#L5-L10>)



```go
type CroissantError struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value any
}
```

<a name="CroissantError.Error"></a>
### func \(CroissantError\) [Error](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/error.go#L12>)

```go
func (e CroissantError) Error() string
```



<a name="DataContext"></a>
## type [DataContext](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L296-L299>)

DataContext represents the data field in the context

```go
type DataContext struct {
    ID   string `json:"@id"`
    Type string `json:"@type"`
}
```

<a name="DataType"></a>
## type [DataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L105-L110>)

DataType represents a data type that can be either a single string or an array of strings

```go
type DataType struct {
    // Single dataType case: just a string value
    SingleType *string `json:"-"`
    // Array dataType case: array of string values
    ArrayType []string `json:"-"`
}
```

<a name="NewArrayDataType"></a>
### func [NewArrayDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L378>)

```go
func NewArrayDataType(dataTypes ...string) DataType
```

NewArrayDataType creates a DataType with multiple types

<a name="NewNullableSingleDataType"></a>
### func [NewNullableSingleDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L364>)

```go
func NewNullableSingleDataType(dataType string) *DataType
```

NewSingleDataType creates a DataType with a single type

<a name="NewSingleDataType"></a>
### func [NewSingleDataType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L371>)

```go
func NewSingleDataType(dataType string) DataType
```

NewSingleDataType creates a DataType with a single type

<a name="DataType.GetFirstType"></a>
### func \(DataType\) [GetFirstType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L219>)

```go
func (d DataType) GetFirstType() string
```

GetFirstType returns the first data type \(useful for backward compatibility\)

<a name="DataType.GetTypes"></a>
### func \(DataType\) [GetTypes](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L208>)

```go
func (d DataType) GetTypes() []string
```

GetTypes returns all data types \(single or array\)

<a name="DataType.IsArray"></a>
### func \(DataType\) [IsArray](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L203>)

```go
func (d DataType) IsArray() bool
```

IsArray returns true if this is an array of data types

<a name="DataType.MarshalJSON"></a>
### func \(DataType\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L172>)

```go
func (d DataType) MarshalJSON() ([]byte, error)
```

MarshalJSON implements custom JSON marshaling for DataType

<a name="DataType.UnmarshalJSON"></a>
### func \(\*DataType\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L183>)

```go
func (d *DataType) UnmarshalJSON(data []byte) error
```

UnmarshalJSON implements custom JSON unmarshaling for DataType

<a name="DataTypeContext"></a>
## type [DataTypeContext](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L302-L305>)

DataTypeContext represents the dataType field in the context

```go
type DataTypeContext struct {
    ID   string `json:"@id"`
    Type string `json:"@type"`
}
```

<a name="Distribution"></a>
## type [Distribution](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L230-L242>)

Distribution represents a file in the Croissant metadata

```go
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
```

<a name="DistributionNode"></a>
## type [DistributionNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L173-L181>)

DistributionNode represents a file distribution

```go
type DistributionNode struct {
    BaseNode
    Type           string `json:"@type"`
    ContentSize    string `json:"contentSize,omitempty"`
    ContentURL     string `json:"contentUrl,omitempty"`
    EncodingFormat string `json:"encodingFormat,omitempty"`
    SHA256         string `json:"sha256,omitempty"`
    MD5            string `json:"md5,omitempty"`
}
```

<a name="DistributionNode.Validate"></a>
### func \(\*DistributionNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L184>)

```go
func (d *DistributionNode) Validate(issues *Issues)
```

Validate validates the distribution node

<a name="Extract"></a>
## type [Extract](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L44-L50>)

Extract represents the extraction information for a field source

```go
type Extract struct {
    Column       string `json:"column,omitempty"`
    JSONPath     string `json:"jsonPath,omitempty"`
    Regex        string `json:"regex,omitempty"`
    Separator    string `json:"separator,omitempty"`
    FileProperty string `json:"fileProperty,omitempty"`
}
```

<a name="ExtractNode"></a>
## type [ExtractNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L417-L422>)

ExtractNode represents extraction details

```go
type ExtractNode struct {
    Regex        string `json:"regex,omitempty"`
    Column       string `json:"column,omitempty"`
    JSONPath     string `json:"jsonPath,omitempty"`
    FileProperty string `json:"fileProperty,omitempty"`
}
```

<a name="Field"></a>
## type [Field](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L9-L32>)

Field represents a field in the Croissant metadata

```go
type Field struct {
    ID   string `json:"@id"`
    Type string `json:"@type"`
    // Name of the field.
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    // Data type of the field identified by the class URI.
    // Usually either an atomic type (e.g, sc:Integer) or a semantic type (e.g., sc:GeoLocation).
    DataType DataType `json:"dataType"`
    // The data source of the field.
    // Represented as a reference to a FileObject or FileSet's contents.
    Source FieldSource `json:"source,omitempty"`
    // If true, field is a list of `dataType` values.
    Repeated bool `json:"repeated,omitempty"`
    // Examples of field values.
    Examples interface{} `json:"examples,omitempty"`
    // Additional fields defined within this one.
    SubField []Field `json:"subField,omitempty"`
    // A special case of SubField.
    // References one or more Fields in the same RecordSet.
    ParentField []FieldRefSlice `json:"parentField,omitempty"`
    // References one or more Fields that are part of a separate RecordSet.
    References []FieldRefSlice `json:"references,omitempty"`
}
```

<a name="FieldNode"></a>
## type [FieldNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L342-L352>)

FieldNode represents a field

```go
type FieldNode struct {
    BaseNode
    Type        string          `json:"@type"`
    Description string          `json:"description,omitempty"`
    DataType    DataType        `json:"dataType,omitempty"`
    Source      SourceNode      `json:"source,omitempty"`
    Repeated    bool            `json:"repeated,omitempty"`
    Examples    interface{}     `json:"examples,omitempty"`
    SubField    []*FieldNode    `json:"subField,omitempty"`
    References  []FieldRefSlice `json:"references,omitempty"`
}
```

<a name="FieldNode.Validate"></a>
### func \(\*FieldNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L355>)

```go
func (f *FieldNode) Validate(issues *Issues)
```

Validate validates the field node

<a name="FieldRef"></a>
## type [FieldRef](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L63-L66>)

FieldRef represents a reference to another field

```go
type FieldRef struct {
    ID    string  `json:"@id,omitempty"`
    Field *KeyRef `json:"field,omitempty"`
}
```

<a name="FieldRefSlice"></a>
## type [FieldRefSlice](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L69>)

Parses ONE or MANY FieldRefs.

```go
type FieldRefSlice []FieldRef
```

<a name="FieldRefSlice.MarshalJSON"></a>
### func \(FieldRefSlice\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L93>)

```go
func (ref FieldRefSlice) MarshalJSON() ([]byte, error)
```



<a name="FieldRefSlice.UnmarshalJSON"></a>
### func \(\*FieldRefSlice\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L71>)

```go
func (ref *FieldRefSlice) UnmarshalJSON(data []byte) error
```



<a name="FieldSource"></a>
## type [FieldSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L35-L41>)

FieldSource represents the source information for a field

```go
type FieldSource struct {
    Extract    Extract    `json:"extract,omitempty"`
    FileObject FileObject `json:"fileObject,omitempty"`
    FileSet    FileObject `json:"fileSet,omitempty"`
    Transform  *Transform `json:"transform,omitempty"`
    Format     string     `json:"format,omitempty"`
}
```

<a name="FieldSource.ValidateSource"></a>
### func \(FieldSource\) [ValidateSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L385>)

```go
func (fs FieldSource) ValidateSource() bool
```

ValidateSource validates the source configuration

<a name="FileObject"></a>
## type [FileObject](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L53-L55>)

FileObject represents a file object reference

```go
type FileObject struct {
    ID string `json:"@id"`
}
```

<a name="FileObjectRef"></a>
## type [FileObjectRef](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L425-L427>)

FileObjectRef represents a reference to a file object

```go
type FileObjectRef struct {
    ID string `json:"@id"`
}
```

<a name="Issue"></a>
## type [Issue](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L19-L23>)

Issue represents a single validation issue

```go
type Issue struct {
    Type    IssueType
    Message string
    Context string // For context like "Metadata(mydataset) > FileObject(a-csv-table)"
}
```

<a name="IssueType"></a>
## type [IssueType](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L11>)

IssueType represents the type of issue \(error or warning\)

```go
type IssueType int
```

<a name="ErrorIssue"></a>

```go
const (
    ErrorIssue IssueType = iota
    WarningIssue
)
```

<a name="Issues"></a>
## type [Issues](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L26-L29>)

Issues represents a collection of validation issues

```go
type Issues struct {
    // contains filtered or unexported fields
}
```

<a name="NewIssues"></a>
### func [NewIssues](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L32>)

```go
func NewIssues() *Issues
```

NewIssues creates a new Issues instance

<a name="ValidateFile"></a>
### func [ValidateFile](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L31>)

```go
func ValidateFile(filePath string) (*Issues, error)
```

ValidateFile validates a Croissant metadata file and returns issues

<a name="ValidateJSON"></a>
### func [ValidateJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L41>)

```go
func ValidateJSON(data []byte) (*Issues, error)
```

ValidateJSON validates Croissant metadata in JSON\-LD format and returns issues

<a name="ValidateJSONWithOptions"></a>
### func [ValidateJSONWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L60>)

```go
func ValidateJSONWithOptions(data []byte, options ValidationOptions) (*Issues, error)
```

ValidateJSONWithOptions validates Croissant metadata in JSON\-LD format with options and returns issues

<a name="ValidateMetadata"></a>
### func [ValidateMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L79>)

```go
func ValidateMetadata(metadata Metadata) *Issues
```

ValidateMetadata validates a Metadata struct and returns issues

<a name="ValidateMetadataWithOptions"></a>
### func [ValidateMetadataWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L84>)

```go
func ValidateMetadataWithOptions(metadata Metadata, options ValidationOptions) *Issues
```

ValidateMetadataWithOptions validates a Metadata struct with specific options

<a name="Issues.AddError"></a>
### func \(\*Issues\) [AddError](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L40>)

```go
func (i *Issues) AddError(message string, node ...Node)
```

AddError adds a new error to the issues collection

<a name="Issues.AddWarning"></a>
### func \(\*Issues\) [AddWarning](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L52>)

```go
func (i *Issues) AddWarning(message string, node ...Node)
```

AddWarning adds a new warning to the issues collection

<a name="Issues.ErrorCount"></a>
### func \(\*Issues\) [ErrorCount](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L74>)

```go
func (i *Issues) ErrorCount() int
```

ErrorCount returns the number of errors

<a name="Issues.HasErrors"></a>
### func \(\*Issues\) [HasErrors](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L64>)

```go
func (i *Issues) HasErrors() bool
```

HasErrors returns true if there are any errors

<a name="Issues.HasWarnings"></a>
### func \(\*Issues\) [HasWarnings](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L69>)

```go
func (i *Issues) HasWarnings() bool
```

HasWarnings returns true if there are any warnings

<a name="Issues.Report"></a>
### func \(\*Issues\) [Report](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L84>)

```go
func (i *Issues) Report() string
```

Report generates a human\-readable report of all issues

<a name="Issues.WarningCount"></a>
### func \(\*Issues\) [WarningCount](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/issues.go#L79>)

```go
func (i *Issues) WarningCount() int
```

WarningCount returns the number of warnings

<a name="JSONLDProcessor"></a>
## type [JSONLDProcessor](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L12-L15>)

JSONLDProcessor handles JSON\-LD processing using json\-gold library

```go
type JSONLDProcessor struct {
    // contains filtered or unexported fields
}
```

<a name="NewJSONLDProcessor"></a>
### func [NewJSONLDProcessor](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L18>)

```go
func NewJSONLDProcessor() *JSONLDProcessor
```

NewJSONLDProcessor creates a new JSON\-LD processor

<a name="JSONLDProcessor.CompactJSONLD"></a>
### func \(\*JSONLDProcessor\) [CompactJSONLD](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L53>)

```go
func (j *JSONLDProcessor) CompactJSONLD(expanded interface{}, context map[string]interface{}) (map[string]interface{}, error)
```

CompactJSONLD compacts an expanded JSON\-LD document with the given context

<a name="JSONLDProcessor.ParseCroissantMetadata"></a>
### func \(\*JSONLDProcessor\) [ParseCroissantMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L80>)

```go
func (j *JSONLDProcessor) ParseCroissantMetadata(data []byte) (*Metadata, error)
```

ParseCroissantMetadata parses Croissant JSON\-LD metadata and converts it to our Metadata struct

<a name="JSONLDProcessor.ParseJSONLD"></a>
### func \(\*JSONLDProcessor\) [ParseJSONLD](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L26>)

```go
func (j *JSONLDProcessor) ParseJSONLD(data []byte) (map[string]interface{}, error)
```

ParseJSONLD parses and expands JSON\-LD document to a normalized form

<a name="JSONLDProcessor.ValidateJSONLD"></a>
### func \(\*JSONLDProcessor\) [ValidateJSONLD](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/jsonld.go#L64>)

```go
func (j *JSONLDProcessor) ValidateJSONLD(data []byte) error
```

ValidateJSONLD validates that the document is valid JSON\-LD

<a name="KeyRef"></a>
## type [KeyRef](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L58-L60>)

KeyRef represents a key reference in a composite key

```go
type KeyRef struct {
    ID string `json:"@id"`
}
```

<a name="Metadata"></a>
## type [Metadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L308-L325>)

Metadata represents the complete Croissant metadata

```go
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
```

<a name="MetadataNode"></a>
## type [MetadataNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L7-L18>)

MetadataNode represents a Croissant metadata document

```go
type MetadataNode struct {
    BaseNode
    Context       Context             `json:"@context"`
    Type          string              `json:"@type"`
    Description   string              `json:"description,omitempty"`
    ConformsTo    string              `json:"conformsTo,omitempty"`
    DatePublished string              `json:"datePublished,omitempty"`
    Version       string              `json:"version,omitempty"`
    Distributions []*DistributionNode `json:"distribution"`
    RecordSets    []*RecordSetNode    `json:"recordSet"`
    Issues        *Issues             `json:"-"` // Not serialized to JSON
}
```

<a name="FromMetadata"></a>
### func [FromMetadata](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L65>)

```go
func FromMetadata(metadata Metadata) *MetadataNode
```

FromMetadata converts a Metadata struct to a MetadataNode

<a name="NewMetadataNode"></a>
### func [NewMetadataNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L21>)

```go
func NewMetadataNode() *MetadataNode
```

NewMetadataNode creates a new MetadataNode

<a name="MetadataNode.Validate"></a>
### func \(\*MetadataNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L35>)

```go
func (m *MetadataNode) Validate(issues *Issues)
```

Validate validates the metadata node

<a name="MetadataWithValidation"></a>
## type [MetadataWithValidation](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L369-L373>)

AddValidationToMetadata adds validation functionality to the Metadata struct

```go
type MetadataWithValidation struct {
    Metadata
    // contains filtered or unexported fields
}
```

<a name="GenerateMetadataWithValidation"></a>
### func [GenerateMetadataWithValidation](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L318>)

```go
func GenerateMetadataWithValidation(csvPath string, outputPath string) (*MetadataWithValidation, error)
```

GenerateMetadataWithValidation generates Croissant metadata with validation from a CSV file

<a name="NewMetadataWithValidation"></a>
### func [NewMetadataWithValidation](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L376>)

```go
func NewMetadataWithValidation(metadata Metadata) *MetadataWithValidation
```

NewMetadataWithValidation creates a new MetadataWithValidation instance

<a name="MetadataWithValidation.GetIssues"></a>
### func \(\*MetadataWithValidation\) [GetIssues](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L419>)

```go
func (m *MetadataWithValidation) GetIssues() *Issues
```

GetIssues returns the validation issues

<a name="MetadataWithValidation.HasErrors"></a>
### func \(\*MetadataWithValidation\) [HasErrors](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L403>)

```go
func (m *MetadataWithValidation) HasErrors() bool
```

HasErrors returns true if there are validation errors

<a name="MetadataWithValidation.HasWarnings"></a>
### func \(\*MetadataWithValidation\) [HasWarnings](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L411>)

```go
func (m *MetadataWithValidation) HasWarnings() bool
```

HasWarnings returns true if there are validation warnings

<a name="MetadataWithValidation.Report"></a>
### func \(\*MetadataWithValidation\) [Report](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L395>)

```go
func (m *MetadataWithValidation) Report() string
```

Report returns a string report of validation issues

<a name="MetadataWithValidation.Validate"></a>
### func \(\*MetadataWithValidation\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L384>)

```go
func (m *MetadataWithValidation) Validate()
```

Validate runs validation on the metadata

<a name="MetadataWithValidation.ValidateWithOptions"></a>
### func \(\*MetadataWithValidation\) [ValidateWithOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L389>)

```go
func (m *MetadataWithValidation) ValidateWithOptions(options ValidationOptions)
```

ValidateWithOptions runs validation with specific options

<a name="Node"></a>
## type [Node](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/node.go#L5-L11>)

Node represents a node in the Croissant metadata structure

```go
type Node interface {
    GetName() string
    GetID() string
    GetParent() Node
    SetParent(Node)
    Validate(*Issues)
}
```

<a name="RecordSet"></a>
## type [RecordSet](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L245-L254>)

RecordSet represents a record set in the Croissant metadata

```go
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
```

<a name="CreateEnumerationRecordSet"></a>
### func [CreateEnumerationRecordSet](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L191>)

```go
func CreateEnumerationRecordSet(id, name string, values []string, urls []string) RecordSet
```

CreateEnumerationRecordSet creates a RecordSet for categorical/enumeration data

<a name="CreateSplitRecordSet"></a>
### func [CreateSplitRecordSet](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/croissant.go#L239>)

```go
func CreateSplitRecordSet() RecordSet
```

CreateSplitRecordSet creates a standard ML split RecordSet

<a name="RecordSetKey"></a>
## type [RecordSetKey](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L113-L118>)

RecordSetKey represents a record set key that can be either a single key or composite key

```go
type RecordSetKey struct {
    // Single key case: just an ID reference
    SingleKey *KeyRef `json:"-"`
    // Composite key case: array of ID references
    CompositeKey []KeyRef `json:"-"`
}
```

<a name="NewCompositeKey"></a>
### func [NewCompositeKey](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L353>)

```go
func NewCompositeKey(keyIDs ...string) *RecordSetKey
```

NewCompositeKey creates a RecordSetKey with multiple key references

<a name="NewSingleKey"></a>
### func [NewSingleKey](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L346>)

```go
func NewSingleKey(keyID string) *RecordSetKey
```

NewSingleKey creates a RecordSetKey with a single key reference

<a name="RecordSetKey.GetKeyIDs"></a>
### func \(RecordSetKey\) [GetKeyIDs](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L157>)

```go
func (k RecordSetKey) GetKeyIDs() []string
```

GetKeyIDs returns all key IDs \(single or composite\)

<a name="RecordSetKey.IsComposite"></a>
### func \(RecordSetKey\) [IsComposite](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L152>)

```go
func (k RecordSetKey) IsComposite() bool
```

IsComposite returns true if this is a composite key

<a name="RecordSetKey.MarshalJSON"></a>
### func \(RecordSetKey\) [MarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L121>)

```go
func (k RecordSetKey) MarshalJSON() ([]byte, error)
```

MarshalJSON implements custom JSON marshaling for RecordSetKey

<a name="RecordSetKey.UnmarshalJSON"></a>
### func \(\*RecordSetKey\) [UnmarshalJSON](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L132>)

```go
func (k *RecordSetKey) UnmarshalJSON(data []byte) error
```

UnmarshalJSON implements custom JSON unmarshaling for RecordSetKey

<a name="RecordSetNode"></a>
## type [RecordSetNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L207-L215>)

RecordSetNode represents a record set

```go
type RecordSetNode struct {
    BaseNode
    Type        string                   `json:"@type"`
    Description string                   `json:"description,omitempty"`
    DataType    DataType                 `json:"dataType,omitempty"`
    Fields      []*FieldNode             `json:"field"`
    Key         *RecordSetKey            `json:"key,omitempty"`
    Data        []map[string]interface{} `json:"data,omitempty"`
}
```

<a name="RecordSetNode.Validate"></a>
### func \(\*RecordSetNode\) [Validate](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L218>)

```go
func (r *RecordSetNode) Validate(issues *Issues)
```

Validate validates the record set node

<a name="Source"></a>
## type [Source](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L338-L343>)

Source represents a more complete source definition

```go
type Source struct {
    Extract    Extract     `json:"extract,omitempty"`
    FileObject FileObject  `json:"fileObject,omitempty"`
    Field      string      `json:"field,omitempty"`
    Transform  []Transform `json:"transform,omitempty"`
}
```

<a name="SourceNode"></a>
## type [SourceNode](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L399-L405>)

SourceNode represents a source

```go
type SourceNode struct {
    Extract    ExtractNode   `json:"extract,omitempty"`
    FileObject FileObjectRef `json:"fileObject,omitempty"`
    FileSet    FileObjectRef `json:"fileSet,omitempty"`
    Transform  *Transform    `json:"transform,omitempty"`
    Format     string        `json:"format,omitempty"`
}
```

<a name="SourceNode.ValidateSource"></a>
### func \(\*SourceNode\) [ValidateSource](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/metadata_node.go#L408>)

```go
func (s *SourceNode) ValidateSource() bool
```

ValidateSource validates the source node

<a name="Transform"></a>
## type [Transform](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/structs.go#L328-L335>)

Transform represents a data transformation

```go
type Transform struct {
    Type      string `json:"@type"`
    Regex     string `json:"regex,omitempty"`
    Replace   string `json:"replace,omitempty"`
    Format    string `json:"format,omitempty"`
    JSONPath  string `json:"jsonPath,omitempty"`
    Separator string `json:"separator,omitempty"`
}
```

<a name="ValidationOptions"></a>
## type [ValidationOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L13-L18>)

ValidationOptions represents options for validation

```go
type ValidationOptions struct {
    StrictMode      bool
    CheckDataTypes  bool
    ValidateURLs    bool
    CheckFileExists bool
}
```

<a name="DefaultValidationOptions"></a>
### func [DefaultValidationOptions](<https://github.com:beyondcivic/gocroissant/blob/main/pkg/croissant/validation.go#L21>)

```go
func DefaultValidationOptions() ValidationOptions
```

DefaultValidationOptions returns default validation options
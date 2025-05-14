// issues.go
package croissant

import (
	"fmt"
	"sort"
	"strings"
)

// IssueType represents the type of issue (error or warning)
type IssueType int

const (
	ErrorIssue IssueType = iota
	WarningIssue
)

// Issue represents a single validation issue
type Issue struct {
	Type    IssueType
	Message string
	Context string // For context like "Metadata(mydataset) > FileObject(a-csv-table)"
}

// Issues represents a collection of validation issues
type Issues struct {
	errors   map[string]struct{}
	warnings map[string]struct{}
}

// NewIssues creates a new Issues instance
func NewIssues() *Issues {
	return &Issues{
		errors:   make(map[string]struct{}),
		warnings: make(map[string]struct{}),
	}
}

// AddError adds a new error to the issues collection
func (i *Issues) AddError(message string, node ...Node) {
	var context string
	if len(node) > 0 {
		context = getIssueContext(node[0])
		if context != "" {
			message = fmt.Sprintf("[%s] %s", context, message)
		}
	}
	i.errors[message] = struct{}{}
}

// AddWarning adds a new warning to the issues collection
func (i *Issues) AddWarning(message string, node ...Node) {
	var context string
	if len(node) > 0 {
		context = getIssueContext(node[0])
		if context != "" {
			message = fmt.Sprintf("[%s] %s", context, message)
		}
	}
	i.warnings[message] = struct{}{}
}

// HasErrors returns true if there are any errors
func (i *Issues) HasErrors() bool {
	return len(i.errors) > 0
}

// HasWarnings returns true if there are any warnings
func (i *Issues) HasWarnings() bool {
	return len(i.warnings) > 0
}

// ErrorCount returns the number of errors
func (i *Issues) ErrorCount() int {
	return len(i.errors)
}

// WarningCount returns the number of warnings
func (i *Issues) WarningCount() int {
	return len(i.warnings)
}

// Report generates a human-readable report of all issues
func (i *Issues) Report() string {
	var result strings.Builder

	// Sort before printing because maps are not ordered
	if len(i.errors) > 0 {
		errors := make([]string, 0, len(i.errors))
		for err := range i.errors {
			errors = append(errors, err)
		}
		sort.Strings(errors)

		result.WriteString(fmt.Sprintf("Found the following %d error(s) during the validation:\n", len(errors)))
		for _, err := range errors {
			result.WriteString(fmt.Sprintf("  -  %s\n", err))
		}
	}

	if len(i.warnings) > 0 {
		warnings := make([]string, 0, len(i.warnings))
		for warn := range i.warnings {
			warnings = append(warnings, warn)
		}
		sort.Strings(warnings)

		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("Found the following %d warning(s) during the validation:\n", len(warnings)))
		for _, warn := range warnings {
			result.WriteString(fmt.Sprintf("  -  %s\n", warn))
		}
	}

	return strings.TrimSpace(result.String())
}

// getIssueContext generates a context string for an issue based on the node
func getIssueContext(node Node) string {
	if node == nil {
		return ""
	}

	// Build up context by traversing parent hierarchy
	var parts []string
	current := node

	for current != nil {
		nodeType := getNodeType(current)
		nodeName := current.GetName()
		if nodeName == "" {
			parts = append([]string{nodeType + "()"}, parts...)
		} else {
			parts = append([]string{nodeType + "(" + nodeName + ")"}, parts...)
		}
		current = current.GetParent()
	}

	return strings.Join(parts, " > ")
}

// getNodeType returns the type name of a node
func getNodeType(node Node) string {
	switch node.(type) {
	case *MetadataNode:
		return "Metadata"
	case *DistributionNode:
		return "FileObject"
	case *RecordSetNode:
		return "RecordSet"
	case *FieldNode:
		return "Field"
	default:
		return "Node"
	}
}

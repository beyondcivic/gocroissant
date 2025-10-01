// File: pkg/croissant/croissant_test.go
package croissant

import (
	"testing"
)

// TestInferDataType tests the InferDataType function for various types.
func TestInferDataType(t *testing.T) {
	cases := []struct {
		val      string
		expected string
	}{
		{"true", "sc:Boolean"},
		{"false", "sc:Boolean"},
		{"123", "sc:Integer"},
		{"-456", "sc:Integer"},
		{"3.14", "sc:Number"},
		{"2.5e10", "sc:Number"},
		{"2023-01-01", "sc:DateTime"},
		{"01/15/2023", "sc:DateTime"},
		{"https://example.com", "sc:URL"},
		{"http://foo.org", "sc:URL"},
		{"test@example.com", "sc:Text"},
		{"foo bar", "sc:Text"},
		{"", "sc:Text"},
	}
	for _, c := range cases {
		got := InferDataType(c.val)
		if got != c.expected {
			t.Errorf("InferDataType(%q) = %q, want %q", c.val, got, c.expected)
		}
	}
}

// TestIsValidDataType tests the IsValidDataType function.
func TestIsValidDataType(t *testing.T) {
	valid := []string{
		"sc:Text", "sc:Boolean", "sc:Integer", "sc:Number", "sc:DateTime", "sc:URL",
		"sc:ImageObject", "sc:VideoObject", "sc:Enumeration", "sc:GeoShape", "sc:GeoCoordinates",
		"cr:Label", "cr:Split", "cr:BoundingBox", "cr:SegmentationMask",
		"cr:TrainingSplit", "cr:ValidationSplit", "cr:TestSplit",
		"wd:Q123", // Wikidata
	}
	for _, typ := range valid {
		if !IsValidDataType(typ) {
			t.Errorf("IsValidDataType(%q) = false, want true", typ)
		}
	}
	invalid := []string{"sc:Foo", "bar", "wd:X123"}
	for _, typ := range invalid {
		if IsValidDataType(typ) {
			t.Errorf("IsValidDataType(%q) = true, want false", typ)
		}
	}
}

// TestInferSemanticDataType tests basic semantic inference.
func TestInferSemanticDataType(t *testing.T) {
	labelFields := []string{"label", "Label", "category", "target", "annotation"}
	for _, field := range labelFields {
		types := InferSemanticDataType(field, "cat", nil)
		if len(types) < 2 || types[1] != "cr:Label" {
			t.Errorf("InferSemanticDataType(%q, \"cat\") should include cr:Label, got %#v", field, types)
		}
	}
	if types := InferSemanticDataType("split", "train", nil); len(types) < 2 || types[0] != "cr:Split" {
		t.Errorf("InferSemanticDataType(\"split\", \"train\") should include cr:Split, got %#v", types)
	}
	if types := InferSemanticDataType("bbox", "10,20,30,40", nil); types[0] != "cr:BoundingBox" {
		t.Errorf("InferSemanticDataType(\"bbox\", ...) should be cr:BoundingBox, got %#v", types)
	}
	// enum context
	ctx := map[string]interface{}{"enumValues": []string{"a", "b"}}
	if types := InferSemanticDataType("foo", "a", ctx); types[0] != "sc:Enumeration" {
		t.Errorf("InferSemanticDataType(..., \"a\", ctx) should include sc:Enumeration, got %#v", types)
	}
}

// TestCleanFieldName
func TestCleanFieldName(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"foo bar", "foo_bar"},
		{"foo-bar", "foo_bar"},
		{" foo ", "foo"},
		{"foo__bar", "foo_bar"},
		{"123abc", "field_123abc"},
		{"a!@#b", "a_b"},
	}
	for _, c := range cases {
		got := cleanFieldName(c.in)
		if got != c.want {
			t.Errorf("cleanFieldName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

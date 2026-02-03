package resume

import (
	"strings"
	"testing"
)

const testYAML = `
cv:
  name: John Doe
  email: john@example.com
  sections:
    experience:
      - company: TechCo
        position: Engineer
        highlights:
          - Built APIs
          - Led team
      - company: StartupXYZ
        position: Developer
    skills:
      technical:
        - Python
        - Go
`

func TestParsePath(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedLen    int
		expectedKeys   []string
		expectedIndex0 int
		expectError    bool
	}{
		{
			name:           "simple key",
			path:           "cv.name",
			expectedLen:    2,
			expectedKeys:   []string{"cv", "name"},
			expectedIndex0: -1,
			expectError:    false,
		},
		{
			name:           "nested keys",
			path:           "cv.sections.experience",
			expectedLen:    3,
			expectedKeys:   []string{"cv", "sections", "experience"},
			expectedIndex0: -1,
			expectError:    false,
		},
		{
			name:           "array index",
			path:           "cv.sections.experience[0]",
			expectedLen:    3,
			expectedKeys:   []string{"cv", "sections", "experience"},
			expectedIndex0: 0,
			expectError:    false,
		},
		{
			name:           "array of objects",
			path:           "cv.sections.experience[0].company",
			expectedLen:    4,
			expectedKeys:   []string{"cv", "sections", "experience", "company"},
			expectedIndex0: 0,
			expectError:    false,
		},
		{
			name:        "invalid path",
			path:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, err := parsePath(tt.path)

			if tt.expectError {
				if err == nil {
					t.Error("parsePath() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("parsePath() error = %v, want nil", err)
				return
			}

			if len(segments) != tt.expectedLen {
				t.Errorf("parsePath() returned %d segments, want %d", len(segments), tt.expectedLen)
			}

			for i, expectedKey := range tt.expectedKeys {
				if i >= len(segments) {
					break
				}
				if segments[i].key != expectedKey {
					t.Errorf("segment[%d].key = %s, want %s", i, segments[i].key, expectedKey)
				}
			}

			if tt.expectedIndex0 >= 0 && len(segments) > 2 {
				if segments[2].index != tt.expectedIndex0 {
					t.Errorf("segment[2].index = %d, want %d", segments[2].index, tt.expectedIndex0)
				}
			}
		})
	}
}

func TestParsePath_InvalidSyntax(t *testing.T) {
	_, err := parsePath("")
	if err == nil {
		t.Error("parsePath() expected error for empty path, got nil")
	}
}

func TestGetNodeByPath_SimpleKey(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	node, err := doc.GetNodeByPath("cv.name")
	if err != nil {
		t.Fatalf("GetNodeByPath() error = %v, want nil", err)
	}

	if node.Value != "John Doe" {
		t.Errorf("node.Value = %s, want %s", node.Value, "John Doe")
	}
}

func TestGetNodeByPath_NestedKeys(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	node, err := doc.GetNodeByPath("cv.sections.experience")
	if err != nil {
		t.Fatalf("GetNodeByPath() error = %v, want nil", err)
	}

	if node.Kind != 2 { // yaml.SequenceNode
		t.Errorf("node.Kind = %d, want 2 (SequenceNode)", node.Kind)
	}
}

func TestGetNodeByPath_ArrayIndex(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	node, err := doc.GetNodeByPath("cv.sections.experience[0]")
	if err != nil {
		t.Fatalf("GetNodeByPath() error = %v, want nil", err)
	}

	if node.Kind != 4 { // yaml.MappingNode
		t.Errorf("node.Kind = %d, want 4 (MappingNode)", node.Kind)
	}
}

func TestGetNodeByPath_ArrayOfObjects(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	node, err := doc.GetNodeByPath("cv.sections.experience[0].company")
	if err != nil {
		t.Fatalf("GetNodeByPath() error = %v, want nil", err)
	}

	if node.Value != "TechCo" {
		t.Errorf("node.Value = %s, want %s", node.Value, "TechCo")
	}
}

func TestGetNodeByPath_InvalidPath(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	_, err = doc.GetNodeByPath("cv.invalid.path")
	if err == nil {
		t.Error("GetNodeByPath() expected error for invalid path, got nil")
	}
}

func TestGetNodeByPath_OutOfBounds(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	_, err = doc.GetNodeByPath("cv.sections.experience[99]")
	if err == nil {
		t.Error("GetNodeByPath() expected error for out of bounds index, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("GetNodeByPath() error should mention out of bounds, got: %v", err)
	}
}

func TestSetNodeValue_Scalar(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	err = doc.SetNodeValue("cv.name", "Jane Doe")
	if err != nil {
		t.Fatalf("SetNodeValue() error = %v, want nil", err)
	}

	value, err := doc.GetScalarValue("cv.name")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if value != "Jane Doe" {
		t.Errorf("GetScalarValue() = %s, want %s", value, "Jane Doe")
	}
}

func TestSetNodeValue_ComplexValue(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	complexValue := map[string]interface{}{
		"company":  "NewCo",
		"position": "Senior Engineer",
	}

	err = doc.SetNodeValue("cv.sections.experience[0]", complexValue)
	if err != nil {
		t.Fatalf("SetNodeValue() error = %v, want nil", err)
	}

	company, err := doc.GetScalarValue("cv.sections.experience[0].company")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if company != "NewCo" {
		t.Errorf("company = %s, want %s", company, "NewCo")
	}
}

func TestSetNodeValue_ArrayElement(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	err = doc.SetNodeValue("cv.sections.experience[0].highlights[0]", "Updated highlight")
	if err != nil {
		t.Fatalf("SetNodeValue() error = %v, want nil", err)
	}

	value, err := doc.GetScalarValue("cv.sections.experience[0].highlights[0]")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if value != "Updated highlight" {
		t.Errorf("value = %s, want %s", value, "Updated highlight")
	}
}

func TestUpdateNodeValue_ScalarNode(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	node, err := doc.GetNodeByPath("cv.name")
	if err != nil {
		t.Fatalf("GetNodeByPath() error = %v", err)
	}

	err = doc.updateNodeValue(node, "Test Name")
	if err != nil {
		t.Fatalf("updateNodeValue() error = %v, want nil", err)
	}

	if node.Value != "Test Name" {
		t.Errorf("node.Value = %s, want %s", node.Value, "Test Name")
	}
}

func TestUpdateNodeValue_ComplexType(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	node, err := doc.GetNodeByPath("cv.sections.experience[0]")
	if err != nil {
		t.Fatalf("GetNodeByPath() error = %v", err)
	}

	newValue := map[string]interface{}{
		"company":  "TestCo",
		"position": "Tester",
	}

	err = doc.updateNodeValue(node, newValue)
	if err != nil {
		t.Fatalf("updateNodeValue() error = %v, want nil", err)
	}

	// Verify the update
	company, err := doc.GetScalarValue("cv.sections.experience[0].company")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if company != "TestCo" {
		t.Errorf("company = %s, want %s", company, "TestCo")
	}
}

func TestUpdateNodeValue_DocumentNodeUnwrapping(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	// This tests the fix for DocumentNode unwrapping (line 216-220 in yaml_manipulator.go)
	arrayValue := []string{"Item1", "Item2"}
	err = doc.SetNodeValue("cv.sections.skills.technical", arrayValue)
	if err != nil {
		t.Fatalf("SetNodeValue() error = %v, want nil", err)
	}

	length, err := doc.GetArrayLength("cv.sections.skills.technical")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}
	if length != 2 {
		t.Errorf("array length = %d, want 2", length)
	}

	value0, _ := doc.GetScalarValue("cv.sections.skills.technical[0]")
	if value0 != "Item1" {
		t.Errorf("value[0] = %s, want Item1", value0)
	}
}

func TestReorderArray(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	// Get original first company
	originalFirst, _ := doc.GetScalarValue("cv.sections.experience[0].company")

	// Reorder: swap first two elements
	err = doc.ReorderArray("cv.sections.experience", []int{1, 0})
	if err != nil {
		t.Fatalf("ReorderArray() error = %v, want nil", err)
	}

	// Verify reorder
	newFirst, err := doc.GetScalarValue("cv.sections.experience[0].company")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}

	if newFirst == originalFirst {
		t.Error("ReorderArray() did not change the order")
	}
}

func TestReorderArray_InvalidOrder(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	// Wrong length
	err = doc.ReorderArray("cv.sections.experience", []int{0})
	if err == nil {
		t.Error("ReorderArray() expected error for wrong length, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "does not match array length") {
		t.Errorf("ReorderArray() error should mention array length, got: %v", err)
	}
}

func TestReorderArray_InvalidIndex(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	// Invalid index
	err = doc.ReorderArray("cv.sections.experience", []int{0, 99})
	if err == nil {
		t.Error("ReorderArray() expected error for invalid index, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid index") {
		t.Errorf("ReorderArray() error should mention invalid index, got: %v", err)
	}
}

func TestReorderArray_NotAnArray(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	err = doc.ReorderArray("cv.name", []int{0})
	if err == nil {
		t.Error("ReorderArray() expected error for non-array, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "does not point to an array") {
		t.Errorf("ReorderArray() error should mention array requirement, got: %v", err)
	}
}

func TestAddArrayElement(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	originalLength, _ := doc.GetArrayLength("cv.sections.experience[0].highlights")

	err = doc.AddArrayElement("cv.sections.experience[0].highlights", "New highlight", 0)
	if err != nil {
		t.Fatalf("AddArrayElement() error = %v, want nil", err)
	}

	newLength, err := doc.GetArrayLength("cv.sections.experience[0].highlights")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}

	if newLength != originalLength+1 {
		t.Errorf("array length = %d, want %d", newLength, originalLength+1)
	}

	// Verify it was added at position 0
	value, err := doc.GetScalarValue("cv.sections.experience[0].highlights[0]")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if value != "New highlight" {
		t.Errorf("value = %s, want %s", value, "New highlight")
	}
}

func TestAddArrayElement_WithPosition(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	err = doc.AddArrayElement("cv.sections.experience[0].highlights", "Middle highlight", 1)
	if err != nil {
		t.Fatalf("AddArrayElement() error = %v, want nil", err)
	}

	value, err := doc.GetScalarValue("cv.sections.experience[0].highlights[1]")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if value != "Middle highlight" {
		t.Errorf("value = %s, want %s", value, "Middle highlight")
	}
}

func TestAddArrayElement_ToNonArray(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	err = doc.AddArrayElement("cv.name", "test", 0)
	if err == nil {
		t.Error("AddArrayElement() expected error for non-array, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "does not point to an array") {
		t.Errorf("AddArrayElement() error should mention array requirement, got: %v", err)
	}
}

func TestRemoveArrayElement(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	originalLength, _ := doc.GetArrayLength("cv.sections.experience[0].highlights")

	err = doc.RemoveArrayElement("cv.sections.experience[0].highlights", 0)
	if err != nil {
		t.Fatalf("RemoveArrayElement() error = %v, want nil", err)
	}

	newLength, err := doc.GetArrayLength("cv.sections.experience[0].highlights")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}

	if newLength != originalLength-1 {
		t.Errorf("array length = %d, want %d", newLength, originalLength-1)
	}
}

func TestRemoveArrayElement_OutOfBounds(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	err = doc.RemoveArrayElement("cv.sections.experience[0].highlights", 99)
	if err == nil {
		t.Error("RemoveArrayElement() expected error for out of bounds, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("RemoveArrayElement() error should mention out of bounds, got: %v", err)
	}
}

func TestGetArrayLength(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	length, err := doc.GetArrayLength("cv.sections.experience")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v, want nil", err)
	}

	if length != 2 {
		t.Errorf("GetArrayLength() = %d, want 2", length)
	}
}

func TestGetArrayLength_NotAnArray(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	_, err = doc.GetArrayLength("cv.name")
	if err == nil {
		t.Error("GetArrayLength() expected error for non-array, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "does not point to an array") {
		t.Errorf("GetArrayLength() error should mention array requirement, got: %v", err)
	}
}

func TestGetScalarValue(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	value, err := doc.GetScalarValue("cv.name")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v, want nil", err)
	}

	if value != "John Doe" {
		t.Errorf("GetScalarValue() = %s, want %s", value, "John Doe")
	}
}

func TestGetScalarValue_NotScalar(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	_, err = doc.GetScalarValue("cv.sections")
	if err == nil {
		t.Error("GetScalarValue() expected error for non-scalar, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "does not point to a scalar") {
		t.Errorf("GetScalarValue() error should mention scalar requirement, got: %v", err)
	}
}

func TestToYAML_RoundTrip(t *testing.T) {
	doc, err := ParseYAMLDocument(testYAML)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() error = %v", err)
	}

	// Modify the document
	err = doc.SetNodeValue("cv.name", "Jane Smith")
	if err != nil {
		t.Fatalf("SetNodeValue() error = %v", err)
	}

	// Convert back to YAML
	yamlStr, err := doc.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML() error = %v, want nil", err)
	}

	// Parse again
	doc2, err := ParseYAMLDocument(yamlStr)
	if err != nil {
		t.Fatalf("ParseYAMLDocument() on round-trip error = %v", err)
	}

	// Verify the change persisted
	value, err := doc2.GetScalarValue("cv.name")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if value != "Jane Smith" {
		t.Errorf("after round-trip, cv.name = %s, want Jane Smith", value)
	}
}

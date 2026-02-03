package changes

import (
	"fmt"
	"testing"

	"tailor/internal/resume"
)

// createTestChange creates a test change with default values
func createTestChange(id, changeType, path string) Change {
	return Change{
		ID:         id,
		Type:       ChangeType(changeType),
		Path:       path,
		Operation:  Operation{NewValue: "test"},
		Reason:     "test reason",
		Confidence: "high",
		Priority:   1,
		Status:     StatusPending,
	}
}

// createTestChangeSet creates a test changeset with the specified number of changes
func createTestChangeSet(numChanges int) *ChangeSet {
	changes := make([]Change, numChanges)
	for i := 0; i < numChanges; i++ {
		changes[i] = createTestChange(
			fmt.Sprintf("change_%03d", i+1),
			"modify",
			fmt.Sprintf("cv.field[%d]", i),
		)
	}
	return &ChangeSet{
		Changes: changes,
		Summary: ChangeSummary{
			TotalChanges:     numChanges,
			SectionsAffected: []string{"test"},
			KeyOptimizations: []string{"test"},
		},
	}
}

// createTestYAMLDoc creates a YAML test document from the provided content
func createTestYAMLDoc(t *testing.T, yamlContent string) *resume.YAMLDocument {
	doc, err := resume.ParseYAMLDocument(yamlContent)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}
	return doc
}

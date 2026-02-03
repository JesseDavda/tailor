package changes

import (
	"os"
	"strings"
	"testing"
)

func TestParseChangeSet_ValidJSON(t *testing.T) {
	jsonData, err := os.ReadFile("testdata/valid_changeset.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	cs, err := ParseChangeSet(string(jsonData))
	if err != nil {
		t.Fatalf("ParseChangeSet() error = %v, want nil", err)
	}

	if cs == nil {
		t.Fatal("ParseChangeSet() returned nil changeset")
	}

	if len(cs.Changes) != 1 {
		t.Errorf("ParseChangeSet() got %d changes, want 1", len(cs.Changes))
	}

	if cs.Summary.TotalChanges != 1 {
		t.Errorf("Summary.TotalChanges = %d, want 1", cs.Summary.TotalChanges)
	}

	// Check first change details
	change := cs.Changes[0]
	if change.ID != "change_001" {
		t.Errorf("Change.ID = %s, want change_001", change.ID)
	}
	if change.Type != ChangeTypeModify {
		t.Errorf("Change.Type = %s, want %s", change.Type, ChangeTypeModify)
	}
	if change.Status != StatusPending {
		t.Errorf("Change.Status = %s, want %s", change.Status, StatusPending)
	}
}

func TestParseChangeSet_ValidJSON_WithCoverLetter(t *testing.T) {
	jsonData, err := os.ReadFile("testdata/changeset_with_cover_letter.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	cs, err := ParseChangeSet(string(jsonData))
	if err != nil {
		t.Fatalf("ParseChangeSet() error = %v, want nil", err)
	}

	if !cs.HasCoverLetter() {
		t.Error("ParseChangeSet() did not parse cover letter")
	}

	if cs.CoverLetter == nil {
		t.Fatal("CoverLetter is nil")
	}

	if len(cs.CoverLetter.BulletPoints) != 3 {
		t.Errorf("CoverLetter.BulletPoints has %d items, want 3", len(cs.CoverLetter.BulletPoints))
	}

	if cs.CoverLetter.FullLetter == "" {
		t.Error("CoverLetter.FullLetter is empty")
	}
}

func TestParseChangeSet_InvalidJSON(t *testing.T) {
	invalidJSON := `{"changes": [invalid json`

	_, err := ParseChangeSet(invalidJSON)
	if err == nil {
		t.Error("ParseChangeSet() expected error for invalid JSON, got nil")
	}
}

func TestParseChangeSet_MarkdownCodeBlocks(t *testing.T) {
	jsonWithMarkdown := "```json\n" + `{
		"changes": [{
			"id": "test",
			"type": "modify",
			"path": "cv.name",
			"operation": {"new_value": "Test"},
			"reason": "test",
			"confidence": "high",
			"priority": 1
		}],
		"summary": {
			"total_changes": 1,
			"sections_affected": ["test"],
			"key_optimizations": ["test"]
		}
	}` + "\n```"

	cs, err := ParseChangeSet(jsonWithMarkdown)
	if err != nil {
		t.Fatalf("ParseChangeSet() error = %v, want nil", err)
	}

	if cs == nil {
		t.Fatal("ParseChangeSet() returned nil changeset")
	}

	if len(cs.Changes) != 1 {
		t.Errorf("ParseChangeSet() got %d changes, want 1", len(cs.Changes))
	}
}

func TestParseChangeSet_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "missing changes array",
			json: `{"summary": {"total_changes": 0, "sections_affected": [], "key_optimizations": []}}`,
		},
		{
			name: "missing summary",
			json: `{"changes": []}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseChangeSet(tt.json)
			if err == nil {
				t.Error("ParseChangeSet() expected error for missing required fields, got nil")
			}
		})
	}
}

func TestCleanMarkdownCodeBlocks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "json code block",
			input:    "```json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "plain code block",
			input:    "```\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "no code block",
			input:    "{\"key\": \"value\"}",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "with extra whitespace",
			input:    "```json\n  {\"key\": \"value\"}  \n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "multiple lines in block",
			input:    "```json\n{\n  \"key\": \"value\"\n}\n```",
			expected: "{\n  \"key\": \"value\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanMarkdownCodeBlocks(tt.input)
			if got != tt.expected {
				t.Errorf("cleanMarkdownCodeBlocks() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestValidateChangeSet_NoChanges(t *testing.T) {
	cs := &ChangeSet{
		Changes: []Change{},
		Summary: ChangeSummary{
			TotalChanges:     0,
			SectionsAffected: []string{},
			KeyOptimizations: []string{},
		},
	}

	err := validateChangeSet(cs)
	if err == nil {
		t.Error("validateChangeSet() expected error for empty changes, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "no changes") {
		t.Errorf("validateChangeSet() error message should mention 'no changes', got: %v", err)
	}
}

func TestValidateChangeSet_InvalidChangeType(t *testing.T) {
	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:   "test",
				Type: "invalid_type",
				Path: "cv.name",
			},
		},
		Summary: ChangeSummary{
			TotalChanges: 1,
		},
	}

	err := validateChangeSet(cs)
	if err == nil {
		t.Error("validateChangeSet() expected error for invalid change type, got nil")
	}
}

func TestValidateChange_AllTypes(t *testing.T) {
	tests := []struct {
		name      string
		change    Change
		wantError bool
	}{
		{
			name: "valid modify",
			change: Change{
				ID:         "test",
				Type:       ChangeTypeModify,
				Path:       "cv.name",
				Operation:  Operation{NewValue: "Test"},
				Confidence: "high",
				Priority:   1,
			},
			wantError: false,
		},
		{
			name: "valid reorder",
			change: Change{
				ID:   "test",
				Type: ChangeTypeReorder,
				Path: "cv.sections.experience",
				Operation: Operation{
					NewOrder: []int{1, 0, 2},
				},
				Confidence: "high",
				Priority:   2,
			},
			wantError: false,
		},
		{
			name: "valid add",
			change: Change{
				ID:         "test",
				Type:       ChangeTypeAdd,
				Path:       "cv.skills",
				Operation:  Operation{NewValue: "Go"},
				Confidence: "medium",
				Priority:   3,
			},
			wantError: false,
		},
		{
			name: "valid remove",
			change: Change{
				ID:   "test",
				Type: ChangeTypeRemove,
				Path: "cv.skills[0]",
				Operation: Operation{
					Index: intPtr(0),
				},
				Confidence: "low",
				Priority:   3,
			},
			wantError: false,
		},
		{
			name: "missing ID",
			change: Change{
				ID:   "",
				Type: ChangeTypeModify,
				Path: "cv.name",
			},
			wantError: true,
		},
		{
			name: "missing path",
			change: Change{
				ID:   "test",
				Type: ChangeTypeModify,
				Path: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChange(&tt.change)
			if (err != nil) != tt.wantError {
				t.Errorf("validateChange() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateChange_InvalidConfidence(t *testing.T) {
	change := Change{
		ID:         "test",
		Type:       ChangeTypeModify,
		Path:       "cv.name",
		Confidence: "invalid",
		Priority:   1,
	}

	err := validateChange(&change)
	if err == nil {
		t.Error("validateChange() expected error for invalid confidence, got nil")
	}
}

func TestValidateChange_InvalidPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority int
	}{
		{"priority 0", 0},
		{"priority 4", 4},
		{"priority -1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := Change{
				ID:         "test",
				Type:       ChangeTypeModify,
				Path:       "cv.name",
				Confidence: "high",
				Priority:   tt.priority,
			}

			err := validateChange(&change)
			if err == nil {
				t.Error("validateChange() expected error for invalid priority, got nil")
			}
		})
	}
}

func TestValidateCoverLetter_Valid(t *testing.T) {
	cl := &CoverLetterData{
		BulletPoints: []string{
			"Point 1 with sufficient length",
			"Point 2 with sufficient length",
			"Point 3 with sufficient length",
		},
		FullLetter: "This is a full letter with enough content to be valid and meaningful. It contains more than 100 characters which is the minimum requirement for a valid cover letter.",
	}

	err := validateCoverLetter(cl)
	if err != nil {
		t.Errorf("validateCoverLetter() error = %v, want nil", err)
	}
}

func TestValidateCoverLetter_EmptyBulletPoints(t *testing.T) {
	cl := &CoverLetterData{
		BulletPoints: []string{},
		FullLetter:   "This is a full letter",
	}

	err := validateCoverLetter(cl)
	if err == nil {
		t.Error("validateCoverLetter() expected error for empty bullet points, got nil")
	}
}

func TestValidateCoverLetter_ShortLetter(t *testing.T) {
	cl := &CoverLetterData{
		BulletPoints: []string{"Point 1", "Point 2"},
		FullLetter:   "Too short",
	}

	err := validateCoverLetter(cl)
	if err == nil {
		t.Error("validateCoverLetter() expected error for short letter, got nil")
	}
}

func TestValidateCoverLetter_Nil(t *testing.T) {
	// validateCoverLetter returns nil for nil input (cover letter is optional)
	err := validateCoverLetter(nil)
	if err != nil {
		t.Errorf("validateCoverLetter() error = %v, want nil for optional cover letter", err)
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}

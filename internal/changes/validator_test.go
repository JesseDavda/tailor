package changes

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
        highlights:
          - Created features
    skills:
      technical:
        - Python
        - Go
`

func TestValidator_ValidateChange_ModifyValid(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeModify,
		Path: "cv.name",
		Operation: Operation{
			NewValue: "Jane Doe",
		},
	}

	err := validator.ValidateChange(&change)
	if err != nil {
		t.Errorf("ValidateChange() error = %v, want nil", err)
	}
}

func TestValidator_ValidateChange_ModifyInvalidPath(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeModify,
		Path: "cv.invalid.path",
		Operation: Operation{
			NewValue: "test",
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for invalid path, got nil")
	}
}

func TestValidator_ValidateChange_ModifyTypeMismatch(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	// Trying to modify an array directly (should specify element)
	change := Change{
		ID:   "test",
		Type: ChangeTypeModify,
		Path: "cv.sections.experience",
		Operation: Operation{
			NewValue: "not an array",
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for array modification, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "specify an element index") {
		t.Errorf("ValidateChange() error should mention element index, got: %v", err)
	}
}

func TestValidator_ValidateChange_ReorderValid(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeReorder,
		Path: "cv.sections.experience",
		Operation: Operation{
			NewOrder: []int{1, 0}, // Swap first two elements
		},
	}

	err := validator.ValidateChange(&change)
	if err != nil {
		t.Errorf("ValidateChange() error = %v, want nil", err)
	}
}

func TestValidator_ValidateChange_ReorderInvalidLength(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeReorder,
		Path: "cv.sections.experience",
		Operation: Operation{
			NewOrder: []int{0}, // Wrong length (array has 2 elements)
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for invalid reorder length, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "does not match array length") {
		t.Errorf("ValidateChange() error should mention array length mismatch, got: %v", err)
	}
}

func TestValidator_ValidateChange_ReorderDuplicateIndex(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeReorder,
		Path: "cv.sections.experience",
		Operation: Operation{
			NewOrder: []int{0, 0}, // Duplicate index
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for duplicate index, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "duplicate index") {
		t.Errorf("ValidateChange() error should mention duplicate index, got: %v", err)
	}
}

func TestValidator_ValidateChange_ReorderOutOfBounds(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeReorder,
		Path: "cv.sections.experience",
		Operation: Operation{
			NewOrder: []int{0, 5}, // Index 5 is out of bounds
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for out of bounds index, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid index") {
		t.Errorf("ValidateChange() error should mention invalid index, got: %v", err)
	}
}

func TestValidator_ValidateChange_AddValid(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeAdd,
		Path: "cv.sections.experience[0].highlights",
		Operation: Operation{
			NewValue: "New highlight",
		},
	}

	err := validator.ValidateChange(&change)
	if err != nil {
		t.Errorf("ValidateChange() error = %v, want nil", err)
	}
}

func TestValidator_ValidateChange_AddWithPosition(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	position := 1
	change := Change{
		ID:   "test",
		Type: ChangeTypeAdd,
		Path: "cv.sections.experience[0].highlights",
		Operation: Operation{
			NewValue: "New highlight",
			Position: &position,
		},
	}

	err := validator.ValidateChange(&change)
	if err != nil {
		t.Errorf("ValidateChange() error = %v, want nil", err)
	}
}

func TestValidator_ValidateChange_AddInvalidPosition(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	position := 999
	change := Change{
		ID:   "test",
		Type: ChangeTypeAdd,
		Path: "cv.sections.experience[0].highlights",
		Operation: Operation{
			NewValue: "New highlight",
			Position: &position,
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for invalid position, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid position") {
		t.Errorf("ValidateChange() error should mention invalid position, got: %v", err)
	}
}

func TestValidator_ValidateChange_AddToNonArray(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeAdd,
		Path: "cv.name", // name is a string, not an array
		Operation: Operation{
			NewValue: "test",
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for adding to non-array, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "must point to array") {
		t.Errorf("ValidateChange() error should mention array requirement, got: %v", err)
	}
}

func TestValidator_ValidateChange_RemoveValid(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	index := 0
	change := Change{
		ID:   "test",
		Type: ChangeTypeRemove,
		Path: "cv.sections.experience[0].highlights",
		Operation: Operation{
			Index: &index,
		},
	}

	err := validator.ValidateChange(&change)
	if err != nil {
		t.Errorf("ValidateChange() error = %v, want nil", err)
	}
}

func TestValidator_ValidateChange_RemoveOutOfBounds(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	index := 999
	change := Change{
		ID:   "test",
		Type: ChangeTypeRemove,
		Path: "cv.sections.experience[0].highlights",
		Operation: Operation{
			Index: &index,
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for out of bounds index, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid index") {
		t.Errorf("ValidateChange() error should mention invalid index, got: %v", err)
	}
}

func TestValidator_ValidateChange_RemoveMissingIndex(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	change := Change{
		ID:   "test",
		Type: ChangeTypeRemove,
		Path: "cv.sections.experience[0].highlights",
		Operation: Operation{
			Index: nil,
		},
	}

	err := validator.ValidateChange(&change)
	if err == nil {
		t.Error("ValidateChange() expected error for missing index, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "requires index") {
		t.Errorf("ValidateChange() error should mention required index, got: %v", err)
	}
}

func TestValidator_ValidateChanges_MultipleChanges(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:     "1",
				Type:   ChangeTypeModify,
				Path:   "cv.name",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "Jane Doe",
				},
			},
			{
				ID:     "2",
				Type:   ChangeTypeModify,
				Path:   "cv.invalid.path",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "test",
				},
			},
			{
				ID:     "3",
				Type:   ChangeTypeModify,
				Path:   "cv.email",
				Status: StatusPending, // Should be skipped
				Operation: Operation{
					NewValue: "test@example.com",
				},
			},
		},
	}

	errors := validator.ValidateChanges(cs)

	// Should have 1 error (change 2 has invalid path, change 3 is pending so skipped)
	if len(errors) != 1 {
		t.Errorf("ValidateChanges() returned %d errors, want 1", len(errors))
	}

	if len(errors) > 0 && !strings.Contains(errors[0].Error(), "change 2") {
		t.Errorf("ValidateChanges() error should mention change 2, got: %v", errors[0])
	}
}

func TestValidator_DetectConflicts(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	validator := NewValidator(doc)

	tests := []struct {
		name           string
		changes        []Change
		expectConflict bool
	}{
		{
			name: "no conflicts",
			changes: []Change{
				{
					ID:     "1",
					Type:   ChangeTypeModify,
					Path:   "cv.name",
					Status: StatusApproved,
					Operation: Operation{NewValue: "Jane"},
				},
				{
					ID:     "2",
					Type:   ChangeTypeModify,
					Path:   "cv.email",
					Status: StatusApproved,
					Operation: Operation{NewValue: "jane@example.com"},
				},
			},
			expectConflict: false,
		},
		{
			name: "same path conflict",
			changes: []Change{
				{
					ID:     "1",
					Type:   ChangeTypeModify,
					Path:   "cv.name",
					Status: StatusApproved,
					Operation: Operation{NewValue: "Jane"},
				},
				{
					ID:     "2",
					Type:   ChangeTypeModify,
					Path:   "cv.name",
					Status: StatusApproved,
					Operation: Operation{NewValue: "John"},
				},
			},
			expectConflict: true,
		},
		{
			name: "pending changes ignored",
			changes: []Change{
				{
					ID:     "1",
					Type:   ChangeTypeModify,
					Path:   "cv.name",
					Status: StatusApproved,
					Operation: Operation{NewValue: "Jane"},
				},
				{
					ID:     "2",
					Type:   ChangeTypeModify,
					Path:   "cv.name",
					Status: StatusPending,
					Operation: Operation{NewValue: "John"},
				},
			},
			expectConflict: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChangeSet{Changes: tt.changes}
			conflicts := validator.DetectConflicts(cs)

			hasConflict := len(conflicts) > 0
			if hasConflict != tt.expectConflict {
				t.Errorf("DetectConflicts() hasConflict = %v, want %v (conflicts: %v)",
					hasConflict, tt.expectConflict, conflicts)
			}
		})
	}
}

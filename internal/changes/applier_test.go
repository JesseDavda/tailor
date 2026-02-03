package changes

import (
	"strings"
	"testing"

	"tailor/internal/terminal"
)

func TestApplier_ApplyChanges_Modify(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

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
				Priority: 1,
			},
		},
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 1", result.SuccessCount)
	}
	if result.FailureCount != 0 {
		t.Errorf("ApplyChanges() FailureCount = %d, want 0", result.FailureCount)
	}

	// Verify the change was applied
	value, err := doc.GetScalarValue("cv.name")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if value != "Jane Doe" {
		t.Errorf("GetScalarValue() = %q, want %q", value, "Jane Doe")
	}
}

func TestApplier_ApplyChanges_Reorder(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:     "1",
				Type:   ChangeTypeReorder,
				Path:   "cv.sections.experience",
				Status: StatusApproved,
				Operation: Operation{
					NewOrder: []int{1, 0}, // Swap first two elements
				},
				Priority: 1,
			},
		},
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 1", result.SuccessCount)
	}

	// Verify order changed - first element should now be StartupXYZ
	company, err := doc.GetScalarValue("cv.sections.experience[0].company")
	if err != nil {
		t.Fatalf("GetScalarValue() error = %v", err)
	}
	if company != "StartupXYZ" {
		t.Errorf("After reorder, first company = %q, want %q", company, "StartupXYZ")
	}
}

func TestApplier_ApplyChanges_Add(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:     "1",
				Type:   ChangeTypeAdd,
				Path:   "cv.sections.experience[0].highlights",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "New highlight",
				},
				Priority: 1,
			},
		},
	}

	// Get original length
	originalLength, err := doc.GetArrayLength("cv.sections.experience[0].highlights")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 1", result.SuccessCount)
	}

	// Verify element was added
	newLength, err := doc.GetArrayLength("cv.sections.experience[0].highlights")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}
	if newLength != originalLength+1 {
		t.Errorf("After add, array length = %d, want %d", newLength, originalLength+1)
	}
}

func TestApplier_ApplyChanges_Remove(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

	index := 0
	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:     "1",
				Type:   ChangeTypeRemove,
				Path:   "cv.sections.experience[0].highlights",
				Status: StatusApproved,
				Operation: Operation{
					Index: &index,
				},
				Priority: 1,
			},
		},
	}

	// Get original length
	originalLength, err := doc.GetArrayLength("cv.sections.experience[0].highlights")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 1", result.SuccessCount)
	}

	// Verify element was removed
	newLength, err := doc.GetArrayLength("cv.sections.experience[0].highlights")
	if err != nil {
		t.Fatalf("GetArrayLength() error = %v", err)
	}
	if newLength != originalLength-1 {
		t.Errorf("After remove, array length = %d, want %d", newLength, originalLength-1)
	}
}

func TestApplier_ApplyChanges_MultipleChanges(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

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
				Priority: 1,
			},
			{
				ID:     "2",
				Type:   ChangeTypeModify,
				Path:   "cv.email",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "jane@example.com",
				},
				Priority: 2,
			},
			{
				ID:     "3",
				Type:   ChangeTypeModify,
				Path:   "cv.sections.experience[0].position",
				Status: StatusPending, // Should be skipped
				Operation: Operation{
					NewValue: "Senior Engineer",
				},
				Priority: 1,
			},
		},
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	// Should apply 2 out of 3 changes (one is pending)
	if result.SuccessCount != 2 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 2", result.SuccessCount)
	}

	// Verify first change
	name, _ := doc.GetScalarValue("cv.name")
	if name != "Jane Doe" {
		t.Errorf("cv.name = %q, want %q", name, "Jane Doe")
	}

	// Verify second change
	email, _ := doc.GetScalarValue("cv.email")
	if email != "jane@example.com" {
		t.Errorf("cv.email = %q, want %q", email, "jane@example.com")
	}

	// Verify third change was NOT applied
	position, _ := doc.GetScalarValue("cv.sections.experience[0].position")
	if position == "Senior Engineer" {
		t.Error("Pending change should not have been applied")
	}
}

func TestApplier_ApplyChanges_PartialFailure(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

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
				Priority: 1,
			},
			{
				ID:     "2",
				Type:   ChangeTypeModify,
				Path:   "cv.invalid.path",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "test",
				},
				Priority: 2,
			},
		},
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 1", result.SuccessCount)
	}
	if result.FailureCount != 1 {
		t.Errorf("ApplyChanges() FailureCount = %d, want 1", result.FailureCount)
	}
	if len(result.Errors) != 1 {
		t.Errorf("ApplyChanges() Errors length = %d, want 1", len(result.Errors))
	}
}

func TestApplier_ApplyChanges_NoApprovedChanges(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:     "1",
				Type:   ChangeTypeModify,
				Path:   "cv.name",
				Status: StatusPending,
				Operation: Operation{
					NewValue: "Jane Doe",
				},
			},
		},
	}

	result, err := applier.ApplyChanges(doc, cs)
	if err != nil {
		t.Fatalf("ApplyChanges() error = %v, want nil", err)
	}

	if result.SuccessCount != 0 {
		t.Errorf("ApplyChanges() SuccessCount = %d, want 0", result.SuccessCount)
	}
	if result.FailureCount != 0 {
		t.Errorf("ApplyChanges() FailureCount = %d, want 0", result.FailureCount)
	}
}

func TestApplier_SortChangesForApplication(t *testing.T) {
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

	changes := []Change{
		{
			ID:       "1",
			Type:     ChangeTypeModify,
			Path:     "cv.name",
			Priority: 2,
		},
		{
			ID:       "2",
			Type:     ChangeTypeReorder,
			Path:     "cv.sections.experience",
			Priority: 3,
		},
		{
			ID:       "3",
			Type:     ChangeTypeModify,
			Path:     "cv.sections.experience[0].company",
			Priority: 1,
		},
	}

	sorted := applier.sortChangesForApplication(changes)

	// Reorder should come first
	if sorted[0].Type != ChangeTypeReorder {
		t.Errorf("First change should be reorder, got %s", sorted[0].Type)
	}

	// Then by priority (1 before 2)
	if sorted[1].ID != "3" {
		t.Errorf("Second change should be ID 3 (priority 1), got %s", sorted[1].ID)
	}
	if sorted[2].ID != "1" {
		t.Errorf("Third change should be ID 1 (priority 2), got %s", sorted[2].ID)
	}
}

func TestApplier_ValidateBeforeApply_Success(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

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
		},
	}

	err := applier.ValidateBeforeApply(doc, cs)
	if err != nil {
		t.Errorf("ValidateBeforeApply() error = %v, want nil", err)
	}
}

func TestApplier_ValidateBeforeApply_ValidationError(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

	cs := &ChangeSet{
		Changes: []Change{
			{
				ID:     "1",
				Type:   ChangeTypeModify,
				Path:   "cv.invalid.path",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "test",
				},
			},
		},
	}

	err := applier.ValidateBeforeApply(doc, cs)
	if err == nil {
		t.Error("ValidateBeforeApply() expected error for invalid path, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("ValidateBeforeApply() error should mention validation failure, got: %v", err)
	}
}

func TestApplier_ValidateBeforeApply_WithConflicts(t *testing.T) {
	doc := createTestYAMLDoc(t, testYAML)
	logger := terminal.NewColorLogger()
	applier := NewApplier(logger)

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
				Path:   "cv.name",
				Status: StatusApproved,
				Operation: Operation{
					NewValue: "John Smith",
				},
			},
		},
	}

	// Should not return error for conflicts, just warn
	err := applier.ValidateBeforeApply(doc, cs)
	if err != nil {
		t.Errorf("ValidateBeforeApply() error = %v, conflicts should only warn", err)
	}
}

func TestApplyResult_ErrorTracking(t *testing.T) {
	result := &ApplyResult{
		SuccessCount: 2,
		FailureCount: 1,
		Errors:       []error{},
	}

	if result.SuccessCount != 2 {
		t.Errorf("SuccessCount = %d, want 2", result.SuccessCount)
	}
	if result.FailureCount != 1 {
		t.Errorf("FailureCount = %d, want 1", result.FailureCount)
	}
	if len(result.Errors) != 0 {
		t.Errorf("Errors length = %d, want 0", len(result.Errors))
	}
}

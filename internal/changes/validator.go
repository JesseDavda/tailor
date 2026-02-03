package changes

import (
	"fmt"
	"reflect"
	"tailor/internal/resume"
)

// Validator validates changes before they are applied
type Validator struct {
	doc *resume.YAMLDocument
}

// NewValidator creates a new validator for a YAML document
func NewValidator(doc *resume.YAMLDocument) *Validator {
	return &Validator{doc: doc}
}

// ValidateChange validates a single change against the document
func (v *Validator) ValidateChange(change *Change) error {
	// Check if path exists
	_, err := v.doc.GetNodeByPath(change.Path)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	// Type-specific validation
	switch change.Type {
	case ChangeTypeModify:
		return v.validateModify(change)
	case ChangeTypeReorder:
		return v.validateReorder(change)
	case ChangeTypeAdd:
		return v.validateAdd(change)
	case ChangeTypeRemove:
		return v.validateRemove(change)
	default:
		return fmt.Errorf("unknown change type: %s", change.Type)
	}
}

// validateModify validates a modify operation
func (v *Validator) validateModify(change *Change) error {
	if change.Operation.NewValue == nil {
		return fmt.Errorf("modify operation requires new_value")
	}

	// Verify the path points to a modifiable node
	node, err := v.doc.GetNodeByPath(change.Path)
	if err != nil {
		return err
	}

	// Allow array modification if the new value is also an array (complete replacement)
	if node.Kind == 2 { // yaml.SequenceNode
		// Check if new_value is an array/slice type
		if !isArrayValue(change.Operation.NewValue) {
			return fmt.Errorf("cannot modify array directly, specify an element index")
		}
	}

	return nil
}

// isArrayValue checks if a value is an array or slice type
func isArrayValue(value any) bool {
	if value == nil {
		return false
	}

	// Use reflection to check if the value is a slice or array
	v := reflect.ValueOf(value)
	kind := v.Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

// validateReorder validates a reorder operation
func (v *Validator) validateReorder(change *Change) error {
	if len(change.Operation.NewOrder) == 0 {
		return fmt.Errorf("reorder operation requires new_order")
	}

	// Get the array length
	length, err := v.doc.GetArrayLength(change.Path)
	if err != nil {
		return fmt.Errorf("reorder validation failed: %w", err)
	}

	// Validate newOrder indices
	if len(change.Operation.NewOrder) != length {
		return fmt.Errorf("newOrder length (%d) does not match array length (%d)",
			len(change.Operation.NewOrder), length)
	}

	// Check all indices are valid and unique
	seen := make(map[int]bool)
	for _, idx := range change.Operation.NewOrder {
		if idx < 0 || idx >= length {
			return fmt.Errorf("invalid index %d in newOrder (array length: %d)", idx, length)
		}
		if seen[idx] {
			return fmt.Errorf("duplicate index %d in newOrder", idx)
		}
		seen[idx] = true
	}

	return nil
}

// validateAdd validates an add operation
func (v *Validator) validateAdd(change *Change) error {
	if change.Operation.NewValue == nil {
		return fmt.Errorf("add operation requires new_value")
	}

	// Verify path points to an array
	length, err := v.doc.GetArrayLength(change.Path)
	if err != nil {
		return fmt.Errorf("add validation failed, path must point to array: %w", err)
	}

	// Validate position if specified
	if change.Operation.Position != nil {
		pos := *change.Operation.Position
		if pos < 0 || pos > length {
			return fmt.Errorf("invalid position %d for array of length %d", pos, length)
		}
	}

	return nil
}

// validateRemove validates a remove operation
func (v *Validator) validateRemove(change *Change) error {
	// Verify path points to an array
	length, err := v.doc.GetArrayLength(change.Path)
	if err != nil {
		return fmt.Errorf("remove validation failed, path must point to array: %w", err)
	}

	// Validate index if specified
	if change.Operation.Index != nil {
		idx := *change.Operation.Index
		if idx < 0 || idx >= length {
			return fmt.Errorf("invalid index %d for array of length %d", idx, length)
		}
	} else {
		return fmt.Errorf("remove operation requires index")
	}

	return nil
}

// ValidateChanges validates all approved changes in a changeset
func (v *Validator) ValidateChanges(changeSet *ChangeSet) []error {
	var errors []error

	for i, change := range changeSet.Changes {
		if !change.IsApproved() {
			continue
		}

		if err := v.ValidateChange(&change); err != nil {
			errors = append(errors, fmt.Errorf("change %d (%s): %w", i+1, change.ID, err))
		}
	}

	return errors
}

// DetectConflicts detects potential conflicts between changes
func (v *Validator) DetectConflicts(changeSet *ChangeSet) []string {
	var conflicts []string
	approvedChanges := changeSet.GetApprovedChanges()

	// Check for changes affecting the same path
	pathCounts := make(map[string]int)
	for _, change := range approvedChanges {
		pathCounts[change.Path]++
	}

	for path, count := range pathCounts {
		if count > 1 {
			conflicts = append(conflicts, fmt.Sprintf("Multiple changes target the same path: %s", path))
		}
	}

	return conflicts
}

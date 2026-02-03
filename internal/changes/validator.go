package changes

import (
	"fmt"
	"reflect"
	"tailor/internal/resume"
)

type Validator struct {
	doc *resume.YAMLDocument
}

func NewValidator(doc *resume.YAMLDocument) *Validator {
	return &Validator{doc: doc}
}

func (v *Validator) ValidateChange(change *Change) error {
	_, err := v.doc.GetNodeByPath(change.Path)
	if err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

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

func (v *Validator) validateModify(change *Change) error {
	if change.Operation.NewValue == nil {
		return fmt.Errorf("modify operation requires new_value")
	}

	node, err := v.doc.GetNodeByPath(change.Path)
	if err != nil {
		return err
	}

	if node.Kind == 2 { // yaml.SequenceNode
		if !isArrayValue(change.Operation.NewValue) {
			return fmt.Errorf("cannot modify array directly, specify an element index")
		}
	}

	return nil
}

func isArrayValue(value any) bool {
	if value == nil {
		return false
	}

	v := reflect.ValueOf(value)
	kind := v.Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func (v *Validator) validateReorder(change *Change) error {
	if len(change.Operation.NewOrder) == 0 {
		return fmt.Errorf("reorder operation requires new_order")
	}

	length, err := v.doc.GetArrayLength(change.Path)
	if err != nil {
		return fmt.Errorf("reorder validation failed: %w", err)
	}

	if len(change.Operation.NewOrder) != length {
		return fmt.Errorf("newOrder length (%d) does not match array length (%d)",
			len(change.Operation.NewOrder), length)
	}

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

func (v *Validator) validateAdd(change *Change) error {
	if change.Operation.NewValue == nil {
		return fmt.Errorf("add operation requires new_value")
	}

	length, err := v.doc.GetArrayLength(change.Path)
	if err != nil {
		return fmt.Errorf("add validation failed, path must point to array: %w", err)
	}

	if change.Operation.Position != nil {
		pos := *change.Operation.Position
		if pos < 0 || pos > length {
			return fmt.Errorf("invalid position %d for array of length %d", pos, length)
		}
	}

	return nil
}

func (v *Validator) validateRemove(change *Change) error {
	length, err := v.doc.GetArrayLength(change.Path)
	if err != nil {
		return fmt.Errorf("remove validation failed, path must point to array: %w", err)
	}

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

func (v *Validator) DetectConflicts(changeSet *ChangeSet) []string {
	var conflicts []string
	approvedChanges := changeSet.GetApprovedChanges()

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

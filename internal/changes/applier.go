package changes

import (
	"fmt"
	"sort"
	"strings"
	"tailor/internal/resume"
	"tailor/internal/terminal"
)

// Applier applies approved changes to a YAML document
type Applier struct {
	logger *terminal.ColorLogger
}

// NewApplier creates a new applier
func NewApplier(logger *terminal.ColorLogger) *Applier {
	return &Applier{
		logger: logger,
	}
}

// ApplyResult contains the result of applying changes
type ApplyResult struct {
	SuccessCount int
	FailureCount int
	Errors       []error
}

// ApplyChanges applies all approved changes to the document
func (a *Applier) ApplyChanges(doc *resume.YAMLDocument, changeSet *ChangeSet) (*ApplyResult, error) {
	result := &ApplyResult{
		Errors: []error{},
	}

	// Get only approved changes
	approvedChanges := changeSet.GetApprovedChanges()
	if len(approvedChanges) == 0 {
		a.logger.Info("No approved changes to apply")
		return result, nil
	}

	a.logger.Info(fmt.Sprintf("Applying %d approved changes", len(approvedChanges)))

	// Sort changes for safe application
	sortedChanges := a.sortChangesForApplication(approvedChanges)

	// Apply each change
	for i, change := range sortedChanges {
		a.logger.Info(fmt.Sprintf("[%d/%d] Applying %s to %s",
			i+1, len(sortedChanges), change.Type, change.Path))

		if err := a.applyChange(doc, &change); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors,
				fmt.Errorf("failed to apply change %s: %w", change.ID, err))
			a.logger.Error(err.Error())
		} else {
			result.SuccessCount++
			a.logger.Success("Applied")
		}
	}

	a.logger.Info(fmt.Sprintf("Applied %d/%d changes successfully",
		result.SuccessCount, result.SuccessCount+result.FailureCount))

	return result, nil
}

// sortChangesForApplication sorts changes to minimize conflicts
func (a *Applier) sortChangesForApplication(changes []Change) []Change {
	sorted := make([]Change, len(changes))
	copy(sorted, changes)

	sort.SliceStable(sorted, func(i, j int) bool {
		// Reorders first (they affect indices)
		if sorted[i].Type == ChangeTypeReorder && sorted[j].Type != ChangeTypeReorder {
			return true
		}
		if sorted[i].Type != ChangeTypeReorder && sorted[j].Type == ChangeTypeReorder {
			return false
		}

		// Then by priority (lower number = higher priority)
		if sorted[i].Priority != sorted[j].Priority {
			return sorted[i].Priority < sorted[j].Priority
		}

		// Then by path depth (deeper first to avoid index shifts)
		depthI := strings.Count(sorted[i].Path, ".")
		depthJ := strings.Count(sorted[j].Path, ".")
		if depthI != depthJ {
			return depthI > depthJ
		}

		// Finally, keep original order
		return false
	})

	return sorted
}

// applyChange applies a single change to the document
func (a *Applier) applyChange(doc *resume.YAMLDocument, change *Change) error {
	switch change.Type {
	case ChangeTypeModify:
		return doc.SetNodeValue(change.Path, change.Operation.NewValue)

	case ChangeTypeReorder:
		return doc.ReorderArray(change.Path, change.Operation.NewOrder)

	case ChangeTypeAdd:
		position := 0
		if change.Operation.Position != nil {
			position = *change.Operation.Position
		}
		return doc.AddArrayElement(change.Path, change.Operation.NewValue, position)

	case ChangeTypeRemove:
		if change.Operation.Index == nil {
			return fmt.Errorf("remove operation requires index")
		}
		return doc.RemoveArrayElement(change.Path, *change.Operation.Index)

	default:
		return fmt.Errorf("unknown change type: %s", change.Type)
	}
}

// ValidateBeforeApply validates all changes before applying any
func (a *Applier) ValidateBeforeApply(doc *resume.YAMLDocument, changeSet *ChangeSet) error {
	validator := NewValidator(doc)

	// Validate each approved change
	errors := validator.ValidateChanges(changeSet)
	if len(errors) > 0 {
		var errMsgs []string
		for _, err := range errors {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Errorf("validation failed:\n%s", strings.Join(errMsgs, "\n"))
	}

	// Check for conflicts
	conflicts := validator.DetectConflicts(changeSet)
	if len(conflicts) > 0 {
		a.logger.Warn("Detected conflicts:")
		for _, conflict := range conflicts {
			a.logger.Warn(fmt.Sprintf("  - %s", conflict))
		}
	}

	return nil
}

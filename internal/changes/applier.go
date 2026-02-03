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

func (a *Applier) ApplyChanges(doc *resume.YAMLDocument, changeSet *ChangeSet) (*ApplyResult, error) {
	result := &ApplyResult{
		Errors: []error{},
	}

	approvedChanges := changeSet.GetApprovedChanges()
	if len(approvedChanges) == 0 {
		a.logger.Info("No approved changes to apply")
		return result, nil
	}

	a.logger.Info(fmt.Sprintf("Applying %d approved changes", len(approvedChanges)))

	sortedChanges := a.sortChangesForApplication(approvedChanges)

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

func (a *Applier) sortChangesForApplication(changes []Change) []Change {
	sorted := make([]Change, len(changes))
	copy(sorted, changes)

	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Type == ChangeTypeReorder && sorted[j].Type != ChangeTypeReorder {
			return true
		}
		if sorted[i].Type != ChangeTypeReorder && sorted[j].Type == ChangeTypeReorder {
			return false
		}

		if sorted[i].Priority != sorted[j].Priority {
			return sorted[i].Priority < sorted[j].Priority
		}

		depthI := strings.Count(sorted[i].Path, ".")
		depthJ := strings.Count(sorted[j].Path, ".")
		if depthI != depthJ {
			return depthI > depthJ
		}

		return false
	})

	return sorted
}

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

func (a *Applier) ValidateBeforeApply(doc *resume.YAMLDocument, changeSet *ChangeSet) error {
	validator := NewValidator(doc)

	errors := validator.ValidateChanges(changeSet)
	if len(errors) > 0 {
		var errMsgs []string
		for _, err := range errors {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Errorf("validation failed:\n%s", strings.Join(errMsgs, "\n"))
	}

	conflicts := validator.DetectConflicts(changeSet)
	if len(conflicts) > 0 {
		a.logger.Warn("Detected conflicts:")
		for _, conflict := range conflicts {
			a.logger.Warn(fmt.Sprintf("  - %s", conflict))
		}
	}

	return nil
}

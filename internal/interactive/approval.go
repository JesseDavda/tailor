package interactive

import (
	"fmt"
	"strings"
	"tailor/internal/changes"
	"tailor/internal/terminal"
)

type ApprovalMode string

const (
	ModeOneByOne ApprovalMode = "one-by-one"
	ModeBatch    ApprovalMode = "batch"
	ModeAutoHigh ApprovalMode = "auto-high"
)

type Approver struct {
	mode    ApprovalMode
	session *terminal.InteractiveSession
}

// NewApprover creates a new approver with the specified mode
func NewApprover(mode string, session *terminal.InteractiveSession) *Approver {
	approvalMode := ApprovalMode(mode)

	// Validate mode
	switch approvalMode {
	case ModeOneByOne, ModeBatch, ModeAutoHigh:
		// Valid
	default:
		approvalMode = ModeOneByOne // Default to one-by-one
	}

	return &Approver{
		mode:    approvalMode,
		session: session,
	}
}

// ReviewChanges interactively reviews all changes in a changeset
func (a *Approver) ReviewChanges(changeSet *changes.ChangeSet) error {
	if changeSet == nil || len(changeSet.Changes) == 0 {
		return fmt.Errorf("no changes to review")
	}

	// Show summary
	a.printSummary(changeSet)

	// Review based on mode
	switch a.mode {
	case ModeOneByOne:
		return a.reviewOneByOne(changeSet)
	case ModeBatch:
		return a.reviewBatch(changeSet)
	case ModeAutoHigh:
		return a.reviewAutoHigh(changeSet)
	default:
		return a.reviewOneByOne(changeSet)
	}
}

// printSummary displays a summary of the changeset
func (a *Approver) printSummary(changeSet *changes.ChangeSet) {
	fmt.Println(bold("\n=== Change Summary ==="))
	fmt.Printf("Total changes: %d\n", changeSet.Summary.TotalChanges)
	fmt.Printf("Sections affected: %s\n", strings.Join(changeSet.Summary.SectionsAffected, ", "))
	if len(changeSet.Summary.KeyOptimizations) > 0 {
		fmt.Printf("Key optimizations: %s\n", strings.Join(changeSet.Summary.KeyOptimizations, ", "))
	}
}

// reviewOneByOne reviews each change individually
func (a *Approver) reviewOneByOne(changeSet *changes.ChangeSet) error {
	totalChanges := len(changeSet.Changes)

	for i := range changeSet.Changes {
		change := &changeSet.Changes[i]

		// Display the change
		a.displayChange(i+1, totalChanges, change)

		// Prompt for decision
		decision, err := a.promptDecision()
		if err != nil {
			return err
		}

		// Apply decision
		switch decision {
		case "Approve":
			change.Approve()
		case "Reject":
			change.Reject()
		case "Quit":
			fmt.Println("\nReview interrupted. Changes reviewed so far will be saved.")
			return nil
		}
	}

	// Show final summary
	a.printFinalSummary(changeSet)

	return nil
}

// reviewBatch reviews changes grouped by section
func (a *Approver) reviewBatch(changeSet *changes.ChangeSet) error {
	// Group changes by section (extracted from path)
	sections := make(map[string][]*changes.Change)
	for i := range changeSet.Changes {
		change := &changeSet.Changes[i]
		section := extractSection(change.Path)
		sections[section] = append(sections[section], change)
	}

	// Review each section
	for section, sectionChanges := range sections {
		fmt.Println(bold(fmt.Sprintf("\n=== Section: %s (%d changes) ===", section, len(sectionChanges))))

		// Show all changes in section
		for i, change := range sectionChanges {
			a.displayChange(i+1, len(sectionChanges), change)
			fmt.Println()
		}

		// Prompt for batch decision
		decision, err := a.promptBatchDecision(section)
		if err != nil {
			return err
		}

		// Apply decision to all changes in section
		for _, change := range sectionChanges {
			switch decision {
			case "Approve All":
				change.Approve()
			case "Reject All":
				change.Reject()
			case "Review Individually":
				// Review this change individually
				subDecision, err := a.promptDecision()
				if err != nil {
					return err
				}
				switch subDecision {
				case "Approve":
					change.Approve()
				case "Reject":
					change.Reject()
				case "Skip":
					change.Skip()
				case "Quit":
					return nil
				}
			case "Skip All":
				change.Skip()
			case "Quit":
				return nil
			}
		}
	}

	a.printFinalSummary(changeSet)
	return nil
}

// reviewAutoHigh auto-approves high priority+confidence, reviews others
func (a *Approver) reviewAutoHigh(changeSet *changes.ChangeSet) error {
	autoApproved := 0
	totalChanges := len(changeSet.Changes)

	for i := range changeSet.Changes {
		change := &changeSet.Changes[i]

		// Auto-approve high confidence + priority 1
		if strings.ToLower(change.Confidence) == "high" && change.Priority == 1 {
			change.Approve()
			autoApproved++
			fmt.Printf("\n[Auto-approved] %s\n", change.Path)
			fmt.Printf("Reason: %s\n", change.Reason)
			continue
		}

		// Review others manually
		a.displayChange(i+1-autoApproved, totalChanges-autoApproved, change)

		decision, err := a.promptDecision()
		if err != nil {
			return err
		}

		switch decision {
		case "Approve":
			change.Approve()
		case "Reject":
			change.Reject()
		case "Skip":
			change.Skip()
		case "Quit":
			fmt.Println("\nReview interrupted.")
			a.printFinalSummary(changeSet)
			return nil
		}
	}

	fmt.Printf("\n%s\n", green(fmt.Sprintf("Auto-approved %d high-priority changes", autoApproved)))
	a.printFinalSummary(changeSet)
	return nil
}

// displayChange shows a single change with formatting
func (a *Approver) displayChange(num, total int, change *changes.Change) {
	// Header
	fmt.Print(FormatChangeHeader(num, total, string(change.Type), change.Path))

	// Priority and confidence
	fmt.Println(FormatPriority(change.Priority, change.Confidence))

	// Reason
	fmt.Print(FormatReason(change.Reason))

	// Show the change based on type
	switch change.Type {
	case changes.ChangeTypeModify:
		fmt.Println(FormatDiff(change.Operation.OldValue, change.Operation.NewValue))

	case changes.ChangeTypeAdd:
		fmt.Println(FormatAddDiff(change.Operation.NewValue))

	case changes.ChangeTypeRemove:
		fmt.Println(FormatRemoveDiff(change.Operation.OldValue))

	case changes.ChangeTypeReorder:
		// For reorder, we'd need to fetch the actual items
		// For now, just show the indices
		fmt.Printf("%s\n", yellow(fmt.Sprintf("Reorder from %v to %v",
			change.Operation.OldOrder, change.Operation.NewOrder)))
	}
}

// promptDecision prompts the user for a decision on a single change
func (a *Approver) promptDecision() (string, error) {
	items := []string{"Approve", "Reject", "Quit"}
	_, result, err := a.session.RunPrompt("Decision", items)
	return result, err
}

// promptBatchDecision prompts for a batch decision on a section
func (a *Approver) promptBatchDecision(section string) (string, error) {
	items := []string{"Approve All", "Reject All", "Review Individually", "Skip All", "Quit"}
	_, result, err := a.session.RunPrompt(fmt.Sprintf("Decide for section '%s'", section), items)
	return result, err
}

// printFinalSummary shows the final approval statistics
func (a *Approver) printFinalSummary(changeSet *changes.ChangeSet) {
	counts := changeSet.CountByStatus()
	fmt.Print(FormatSummary(
		counts[changes.StatusApproved],
		counts[changes.StatusRejected],
		counts[changes.StatusSkipped],
	))
}

// extractSection extracts the section name from a path
// Example: "cv.sections.experience[0].highlights[2]" -> "experience"
func extractSection(path string) string {
	parts := strings.Split(path, ".")
	for i, part := range parts {
		if part == "sections" && i+1 < len(parts) {
			section := parts[i+1]
			// Remove array index if present
			if idx := strings.Index(section, "["); idx > 0 {
				section = section[:idx]
			}
			return section
		}
	}
	return "unknown"
}

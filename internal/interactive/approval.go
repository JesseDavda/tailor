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

func NewApprover(mode string, session *terminal.InteractiveSession) *Approver {
	approvalMode := ApprovalMode(mode)

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

func (a *Approver) ReviewChanges(changeSet *changes.ChangeSet) error {
	if changeSet == nil || len(changeSet.Changes) == 0 {
		return fmt.Errorf("no changes to review")
	}

	a.printSummary(changeSet)

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

func (a *Approver) printSummary(changeSet *changes.ChangeSet) {
	fmt.Println(bold("\n=== Change Summary ==="))
	fmt.Printf("Total changes: %d\n", changeSet.Summary.TotalChanges)
	fmt.Printf("Sections affected: %s\n", strings.Join(changeSet.Summary.SectionsAffected, ", "))
	if len(changeSet.Summary.KeyOptimizations) > 0 {
		fmt.Printf("Key optimizations: %s\n", strings.Join(changeSet.Summary.KeyOptimizations, ", "))
	}
}

func (a *Approver) reviewOneByOne(changeSet *changes.ChangeSet) error {
	totalChanges := len(changeSet.Changes)

	for i := range changeSet.Changes {
		change := &changeSet.Changes[i]

		a.displayChange(i+1, totalChanges, change)

		decision, err := a.promptDecision()
		if err != nil {
			return err
		}

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

	a.printFinalSummary(changeSet)

	return nil
}

func (a *Approver) reviewBatch(changeSet *changes.ChangeSet) error {
	sections := make(map[string][]*changes.Change)
	for i := range changeSet.Changes {
		change := &changeSet.Changes[i]
		section := extractSection(change.Path)
		sections[section] = append(sections[section], change)
	}

	for section, sectionChanges := range sections {
		fmt.Println(bold(fmt.Sprintf("\n=== Section: %s (%d changes) ===", section, len(sectionChanges))))

		for i, change := range sectionChanges {
			a.displayChange(i+1, len(sectionChanges), change)
			fmt.Println()
		}

		decision, err := a.promptBatchDecision(section)
		if err != nil {
			return err
		}

		for _, change := range sectionChanges {
			switch decision {
			case "Approve All":
				change.Approve()
			case "Reject All":
				change.Reject()
			case "Review Individually":
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

func (a *Approver) reviewAutoHigh(changeSet *changes.ChangeSet) error {
	autoApproved := 0
	totalChanges := len(changeSet.Changes)

	for i := range changeSet.Changes {
		change := &changeSet.Changes[i]

		if strings.ToLower(change.Confidence) == "high" && change.Priority == 1 {
			change.Approve()
			autoApproved++
			fmt.Printf("\n[Auto-approved] %s\n", change.Path)
			fmt.Printf("Reason: %s\n", change.Reason)
			continue
		}

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

func (a *Approver) displayChange(num, total int, change *changes.Change) {
	fmt.Print(FormatChangeHeader(num, total, string(change.Type), change.Path))

	fmt.Println(FormatPriority(change.Priority, change.Confidence))

	fmt.Print(FormatReason(change.Reason))

	switch change.Type {
	case changes.ChangeTypeModify:
		fmt.Println(FormatDiff(change.Operation.OldValue, change.Operation.NewValue))

	case changes.ChangeTypeAdd:
		fmt.Println(FormatAddDiff(change.Operation.NewValue))

	case changes.ChangeTypeRemove:
		fmt.Println(FormatRemoveDiff(change.Operation.OldValue))

	case changes.ChangeTypeReorder:
		fmt.Printf("%s\n", yellow(fmt.Sprintf("Reorder from %v to %v",
			change.Operation.OldOrder, change.Operation.NewOrder)))
	}
}

func (a *Approver) promptDecision() (string, error) {
	items := []string{"Approve", "Reject", "Quit"}
	_, result, err := a.session.RunPrompt("Decision", items)
	return result, err
}

func (a *Approver) promptBatchDecision(section string) (string, error) {
	items := []string{"Approve All", "Reject All", "Review Individually", "Skip All", "Quit"}
	_, result, err := a.session.RunPrompt(fmt.Sprintf("Decide for section '%s'", section), items)
	return result, err
}

func (a *Approver) printFinalSummary(changeSet *changes.ChangeSet) {
	counts := changeSet.CountByStatus()
	fmt.Print(FormatSummary(
		counts[changes.StatusApproved],
		counts[changes.StatusRejected],
		counts[changes.StatusSkipped],
	))
}

func extractSection(path string) string {
	parts := strings.Split(path, ".")
	for i, part := range parts {
		if part == "sections" && i+1 < len(parts) {
			section := parts[i+1]
			if idx := strings.Index(section, "["); idx > 0 {
				section = section[:idx]
			}
			return section
		}
	}
	return "unknown"
}

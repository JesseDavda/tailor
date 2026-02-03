package interactive

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

func FormatDiff(oldValue, newValue any) string {
	oldStr := fmt.Sprintf("%v", oldValue)
	newStr := fmt.Sprintf("%v", newValue)

	var result strings.Builder

	if oldStr != "" {
		result.WriteString(red("- " + oldStr))
		result.WriteString("\n")
	}

	if newStr != "" {
		result.WriteString(green("+ " + newStr))
	}

	return result.String()
}

func FormatReorderDiff(oldOrder, newOrder []int, items []string) string {
	var result strings.Builder

	result.WriteString(yellow("Old order:\n"))
	for i, idx := range oldOrder {
		if idx < len(items) {
			result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, items[idx]))
		}
	}

	result.WriteString("\n")
	result.WriteString(green("New order:\n"))
	for i, idx := range newOrder {
		if idx < len(items) {
			result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, items[idx]))
		}
	}

	return result.String()
}

func FormatChangeHeader(changeNum, totalChanges int, changeType, path string) string {
	var result strings.Builder

	result.WriteString(bold(fmt.Sprintf("\n--- Change %d of %d ---\n", changeNum, totalChanges)))
	result.WriteString(fmt.Sprintf("Type: %s\n", cyan(changeType)))
	result.WriteString(fmt.Sprintf("Location: %s\n", path))

	return result.String()
}

func FormatPriority(priority int, confidence string) string {
	var priorityStr string
	switch priority {
	case 1:
		priorityStr = "P1 (Critical)"
	case 2:
		priorityStr = "P2 (Important)"
	case 3:
		priorityStr = "P3 (Minor)"
	default:
		priorityStr = fmt.Sprintf("P%d", priority)
	}

	confidenceColor := green
	switch strings.ToLower(confidence) {
	case "high":
		confidenceColor = green
	case "medium":
		confidenceColor = yellow
	case "low":
		confidenceColor = red
	}

	return fmt.Sprintf("Priority: %s | Confidence: %s",
		yellow(priorityStr),
		confidenceColor(strings.Title(confidence)))
}

func FormatReason(reason string) string {
	return fmt.Sprintf("\nReason: %s\n", reason)
}

func FormatSummary(approved, rejected, skipped int) string {
	var result strings.Builder

	result.WriteString(bold("\n=== Review Summary ===\n"))
	result.WriteString(fmt.Sprintf("%s: %d\n", green("Approved"), approved))
	result.WriteString(fmt.Sprintf("%s: %d\n", red("Rejected"), rejected))
	result.WriteString(fmt.Sprintf("%s: %d\n", yellow("Skipped"), skipped))

	return result.String()
}

func FormatAddDiff(newValue any) string {
	return green(fmt.Sprintf("+ ADD: %v", newValue))
}

func FormatRemoveDiff(oldValue any) string {
	return red(fmt.Sprintf("- REMOVE: %v", oldValue))
}

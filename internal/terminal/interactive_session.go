package terminal

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// InteractiveSession handles collapsed prompt display
type InteractiveSession struct {
	manager    *Manager
	selections []string
	cyan       *color.Color
}

// newInteractiveSession creates a new interactive session
func newInteractiveSession(manager *Manager) *InteractiveSession {
	return &InteractiveSession{
		manager:    manager,
		selections: []string{},
		cyan:       color.New(color.FgCyan),
	}
}

// RunPrompt shows a promptui.Select, captures selection, and displays collapsed summary
func (s *InteractiveSession) RunPrompt(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	idx, result, err := prompt.Run()
	if err != nil {
		return idx, result, fmt.Errorf("prompt failed: %w", err)
	}

	// Capture selection for summary
	s.selections = append(s.selections, fmt.Sprintf("%s: %s", label, result))

	// Clear the prompt lines using ANSI escape codes
	// promptui typically uses 2 header lines + number of items
	linesToClear := len(items) + 2
	for i := 0; i < linesToClear; i++ {
		fmt.Print("\033[A\033[K") // Move up one line, clear it
	}

	// Print collapsed summary
	fmt.Printf("✓ %s\n", s.cyan.Sprintf("%s: %s", label, result))

	return idx, result, nil
}

// ShowSummary displays all captured selections
func (s *InteractiveSession) ShowSummary() {
	if len(s.selections) > 0 {
		fmt.Println("\nReview Summary:")
		for _, selection := range s.selections {
			fmt.Printf("  • %s\n", s.cyan.Sprint(selection))
		}
	}
}

// Resume restarts the spinner with an optional message
func (s *InteractiveSession) Resume(msg string) {
	if msg != "" {
		s.manager.UpdateSpinner(msg)
	}
}

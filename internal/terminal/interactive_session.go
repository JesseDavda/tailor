package terminal

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

type InteractiveSession struct {
	manager    *Manager
	selections []string
	cyan       *color.Color
}

func newInteractiveSession(manager *Manager) *InteractiveSession {
	return &InteractiveSession{
		manager:    manager,
		selections: []string{},
		cyan:       color.New(color.FgCyan),
	}
}

func (s *InteractiveSession) RunPrompt(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	idx, result, err := prompt.Run()
	if err != nil {
		return idx, result, fmt.Errorf("prompt failed: %w", err)
	}

	s.selections = append(s.selections, fmt.Sprintf("%s: %s", label, result))

	linesToClear := len(items) + 2
	for i := 0; i < linesToClear; i++ {
		fmt.Print("\033[A\033[K") // Move up one line, clear it
	}

	fmt.Printf("✓ %s\n", s.cyan.Sprintf("%s: %s", label, result))

	return idx, result, nil
}

func (s *InteractiveSession) ShowSummary() {
	if len(s.selections) > 0 {
		fmt.Println("\nReview Summary:")
		for _, selection := range s.selections {
			fmt.Printf("  • %s\n", s.cyan.Sprint(selection))
		}
	}
}

func (s *InteractiveSession) Resume(msg string) {
	if msg != "" {
		s.manager.UpdateSpinner(msg)
	}
}

package executor

import (
	"context"
	"tailor/internal/terminal"
	"tailor/pkg/config"
)

// Executor orchestrates the resume tailoring workflow
type Executor struct {
	config   *config.Config
	terminal *terminal.Manager
	ctx      context.Context
}

// New creates a new Executor with the provided configuration
func New(cfg *config.Config) *Executor {
	return &Executor{
		config:   cfg,
		terminal: terminal.NewManager("Starting Tailor"),
		ctx:      context.Background(),
	}
}

// Run executes the complete resume tailoring workflow
func (e *Executor) Run() error {
	e.terminal.Start()
	defer e.terminal.Shutdown()

	// Load inputs
	masterYAML, err := e.loadResume()
	if err != nil {
		return err
	}

	jobDesc, err := e.loadJobDescription()
	if err != nil {
		return err
	}

	// Get changes from Claude
	changeSet, err := e.analyzeResume(masterYAML, jobDesc)
	if err != nil {
		return err
	}

	// Interactive review
	if err := e.reviewChanges(changeSet); err != nil {
		return err
	}

	// Check if any changes were approved
	approvedCount := len(changeSet.GetApprovedChanges())
	if approvedCount == 0 {
		e.terminal.Info("No changes approved. Exiting.")
		return nil
	}

	// Apply changes
	tailoredYAML, err := e.applyChanges(masterYAML, changeSet)
	if err != nil {
		return err
	}

	// Write output
	return e.writeOutput(tailoredYAML)
}

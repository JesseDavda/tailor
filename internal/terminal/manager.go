package terminal

import (
	"context"

	pin "github.com/yarlson/pin"
)

// Manager coordinates all terminal output (spinner, logger, interactive prompts)
type Manager struct {
	pin           *pin.Pin
	logger        *ColorLogger
	ctx           context.Context
	cancel        context.CancelFunc
	spinnerActive bool
}

// NewManager creates a new terminal manager with the initial spinner message
func NewManager(initialMsg string) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		pin:           pin.New(initialMsg),
		logger:        NewColorLogger(),
		ctx:           ctx,
		cancel:        cancel,
		spinnerActive: false,
	}
}

// Start starts the spinner
func (m *Manager) Start() {
	if !m.spinnerActive {
		m.cancel = m.pin.Start(m.ctx)
		m.spinnerActive = true
	}
}

func (m *Manager) StopSpinner() {
	if m.spinnerActive {
		m.pin.Stop()
		m.spinnerActive = false
	}
}

// UpdateSpinner updates the spinner message
func (m *Manager) UpdateSpinner(msg string) {
	if !m.spinnerActive {
		m.Start()
	}

	m.pin.UpdateMessage(msg)
}

// Fail stops the spinner with a failure message
func (m *Manager) Fail(msg string) {
	m.pin.Fail(msg)
	m.spinnerActive = false
}

// Success logs a success message
func (m *Manager) Success(msg string) {
	m.logger.Success(msg)
}

// Error logs an error message
func (m *Manager) Error(msg string) {
	m.logger.Error(msg)
}

// Warn logs a warning message
func (m *Manager) Warn(msg string) {
	m.logger.Warn(msg)
}

// Info logs an info message
func (m *Manager) Info(msg string) {
	m.logger.Info(msg)
}

// Section logs a section header
func (m *Manager) Section(title string) {
	m.logger.Section(title)
}

func (m *Manager) PauseForInteractive() *InteractiveSession {
	m.StopSpinner()
	return newInteractiveSession(m)
}

// GetLogger returns the underlying color logger
func (m *Manager) GetLogger() *ColorLogger {
	return m.logger
}

// Shutdown stops the spinner gracefully
func (m *Manager) Shutdown() {
	if m.spinnerActive {
		m.cancel()
		m.spinnerActive = false
	}
}

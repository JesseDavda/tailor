package terminal

import (
	"context"

	pin "github.com/yarlson/pin"
)

type Manager struct {
	pin           *pin.Pin
	logger        *ColorLogger
	ctx           context.Context
	cancel        context.CancelFunc
	spinnerActive bool
}

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

func (m *Manager) UpdateSpinner(msg string) {
	if !m.spinnerActive {
		m.Start()
	}

	m.pin.UpdateMessage(msg)
}

func (m *Manager) Fail(msg string) {
	m.pin.Fail(msg)
	m.spinnerActive = false
}

func (m *Manager) Success(msg string) {
	m.logger.Success(msg)
}

func (m *Manager) Error(msg string) {
	m.logger.Error(msg)
}

func (m *Manager) Warn(msg string) {
	m.logger.Warn(msg)
}

func (m *Manager) Info(msg string) {
	m.logger.Info(msg)
}

func (m *Manager) Section(title string) {
	m.logger.Section(title)
}

func (m *Manager) PauseForInteractive() *InteractiveSession {
	m.StopSpinner()
	return newInteractiveSession(m)
}

func (m *Manager) GetLogger() *ColorLogger {
	return m.logger
}

func (m *Manager) Shutdown() {
	if m.spinnerActive {
		m.cancel()
		m.spinnerActive = false
	}
}

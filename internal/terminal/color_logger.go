package terminal

import (
	"github.com/fatih/color"
)

// ColorLogger provides colored logging output
type ColorLogger struct {
	success *color.Color
	error   *color.Color
	warning *color.Color
	info    *color.Color
	section *color.Color
}

// NewColorLogger creates a new colored logger
func NewColorLogger() *ColorLogger {
	return &ColorLogger{
		success: color.New(color.FgGreen, color.Bold),
		error:   color.New(color.FgRed, color.Bold),
		warning: color.New(color.FgYellow),
		info:    color.New(color.FgCyan),
		section: color.New(color.FgHiCyan, color.Bold, color.Underline),
	}
}

// Success logs a success message with green color
func (l *ColorLogger) Success(msg string) {
	l.success.Printf("✓ %s\n", msg)
}

// Error logs an error message with red color
func (l *ColorLogger) Error(msg string) {
	l.error.Printf("✗ %s\n", msg)
}

// Warn logs a warning message with yellow color
func (l *ColorLogger) Warn(msg string) {
	l.warning.Printf("⚠ %s\n", msg)
}

// Info logs an info message with cyan color
func (l *ColorLogger) Info(msg string) {
	l.info.Printf("• %s\n", msg)
}

// Section logs a section header with bold cyan color
func (l *ColorLogger) Section(title string) {
	l.section.Printf("\n%s\n", title)
}

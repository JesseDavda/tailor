package logger

import (
	"fmt"
	"strings"
)

// Logger provides simple, always-verbose logging for the application
type Logger struct{}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{}
}

// Info prints an informational message with a checkmark
func (l *Logger) Info(msg string) {
	fmt.Println("✓", msg)
}

// Step prints a step message without a newline (for ongoing operations)
func (l *Logger) Step(msg string) {
	fmt.Print("▸ ", msg, "... ")
}

// Success prints a success message with a checkmark
func (l *Logger) Success(msg string) {
	fmt.Println("✓", msg)
}

// Warn prints a warning message
func (l *Logger) Warn(msg string) {
	fmt.Println("⚠", msg)
}

// Error prints an error message
func (l *Logger) Error(msg string) {
	fmt.Println("✗", msg)
}

// Section prints a section header with underline
func (l *Logger) Section(title string) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))
}

package terminal

import (
	"github.com/fatih/color"
)

type ColorLogger struct {
	success *color.Color
	error   *color.Color
	warning *color.Color
	info    *color.Color
	section *color.Color
}

func NewColorLogger() *ColorLogger {
	return &ColorLogger{
		success: color.New(color.FgGreen, color.Bold),
		error:   color.New(color.FgRed, color.Bold),
		warning: color.New(color.FgYellow),
		info:    color.New(color.FgCyan),
		section: color.New(color.FgHiCyan, color.Bold, color.Underline),
	}
}

func (l *ColorLogger) Success(msg string) {
	l.success.Printf("✓ %s\n", msg)
}

func (l *ColorLogger) Error(msg string) {
	l.error.Printf("✗ %s\n", msg)
}

func (l *ColorLogger) Warn(msg string) {
	l.warning.Printf("⚠ %s\n", msg)
}

func (l *ColorLogger) Info(msg string) {
	l.info.Printf("• %s\n", msg)
}

func (l *ColorLogger) Section(title string) {
	l.section.Printf("\n%s\n", title)
}

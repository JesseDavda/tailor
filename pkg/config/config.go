package config

import (
	"fmt"
	"os"
	"strings"
)

var Global = &Config{
	Model: "claude-haiku-4-5",
}

type Config struct {
	// API Configuration
	APIKey string
	Model  string

	// Input Configuration
	ResumePath string
	JobPath    string
	JobText    string

	// Output Configuration
	OutputPath string

	// Cover Letter Configuration
	GenerateCoverLetter   bool
	CoverLetterOutputPath string

	// Runtime Configuration
	DryRun       bool
	ApprovalMode string
}

func (c *Config) Validate() error {
	// Check for API key
	if c.APIKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY is required (set via --api-key flag or ANTHROPIC_API_KEY environment variable)")
	}

	// Check for resume path
	if c.ResumePath == "" {
		return fmt.Errorf("--resume/-r flag is required")
	}

	// Check if resume file exists
	if _, err := os.Stat(c.ResumePath); os.IsNotExist(err) {
		return fmt.Errorf("resume file not found: %s", c.ResumePath)
	}

	// Check for job description (either file or text)
	if c.JobPath == "" && c.JobText == "" {
		return fmt.Errorf("either --job/-j (file path) or --job-text (direct text) is required")
	}

	// If both are provided, prefer job file
	if c.JobPath != "" && c.JobText != "" {
		c.JobText = "" // Clear job text to use path instead
	}

	// Check if job file exists (if path is provided)
	if c.JobPath != "" {
		if _, err := os.Stat(c.JobPath); os.IsNotExist(err) {
			return fmt.Errorf("job description file not found: %s", c.JobPath)
		}
	}

	// Validate approval mode
	validModes := map[string]bool{
		"one-by-one": true,
		"batch":      true,
		"auto-high":  true,
	}
	if c.ApprovalMode == "" {
		c.ApprovalMode = "one-by-one" // Default
	} else if !validModes[c.ApprovalMode] {
		return fmt.Errorf("invalid approval mode: %s (must be one-by-one, batch, or auto-high)", c.ApprovalMode)
	}

	// Auto-generate cover letter output path if not specified
	if c.GenerateCoverLetter && c.CoverLetterOutputPath == "" {
		c.CoverLetterOutputPath = generateCoverLetterPath(c.OutputPath)
	}

	return nil
}

func LoadAPIKey(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return os.Getenv("ANTHROPIC_API_KEY")
}

func generateCoverLetterPath(resumePath string) string {
	base := resumePath
	if idx := strings.LastIndex(base, "."); idx > 0 {
		base = base[:idx]
	}
	return base + "-cover-letter.txt"
}

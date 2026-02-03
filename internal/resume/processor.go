package resume

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadResumeYAML reads a resume YAML file from the given path
func ReadResumeYAML(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read resume file: %w", err)
	}

	content := string(data)

	// Validate it's valid YAML
	if err := ValidateYAML(content); err != nil {
		return "", fmt.Errorf("invalid YAML in resume file: %w", err)
	}

	return content, nil
}

// WriteResumeYAML writes the tailored resume YAML to the given path
func WriteResumeYAML(path string, content string) error {
	// Validate YAML before writing
	if err := ValidateYAML(content); err != nil {
		return fmt.Errorf("refusing to write invalid YAML: %w", err)
	}

	// Write with appropriate permissions (0644 - owner can read/write, others can read)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write resume file: %w", err)
	}

	return nil
}

// ValidateYAML checks if the given string is valid YAML
func ValidateYAML(content string) error {
	var data interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}
	return nil
}

// CountLines returns the number of lines in the content
func CountLines(content string) int {
	if content == "" {
		return 0
	}
	return len(strings.Split(content, "\n"))
}

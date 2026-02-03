package resume

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadResumeYAML(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read resume file: %w", err)
	}

	content := string(data)

	if err := ValidateYAML(content); err != nil {
		return "", fmt.Errorf("invalid YAML in resume file: %w", err)
	}

	return content, nil
}

func WriteResumeYAML(path string, content string) error {
	if err := ValidateYAML(content); err != nil {
		return fmt.Errorf("refusing to write invalid YAML: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write resume file: %w", err)
	}

	return nil
}

func ValidateYAML(content string) error {
	var data interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}
	return nil
}

func CountLines(content string) int {
	if content == "" {
		return 0
	}
	return len(strings.Split(content, "\n"))
}

package job

import (
	"fmt"
	"os"
	"strings"
)

func ReadJobDescription(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read job description file: %w", err)
	}

	content := string(data)
	content = strings.TrimSpace(content)

	if content == "" {
		return "", fmt.Errorf("job description file is empty")
	}

	return content, nil
}

func ParseJobText(text string) (string, error) {
	text = strings.TrimSpace(text)

	if text == "" {
		return "", fmt.Errorf("job description text is empty")
	}

	return text, nil
}

package resume

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadResumeYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "resume.yaml")

	testContent := `cv:
  name: John Doe
  email: john@example.com
`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	content, err := ReadResumeYAML(testFile)
	if err != nil {
		t.Fatalf("ReadResumeYAML() error = %v, want nil", err)
	}

	if !strings.Contains(content, "John Doe") {
		t.Errorf("ReadResumeYAML() content does not contain expected data")
	}
}

func TestReadResumeYAML_FileNotFound(t *testing.T) {
	_, err := ReadResumeYAML("/nonexistent/file.yaml")
	if err == nil {
		t.Error("ReadResumeYAML() expected error for nonexistent file, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to read resume file") {
		t.Errorf("ReadResumeYAML() error should mention file read failure, got: %v", err)
	}
}

func TestReadResumeYAML_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.yaml")

	invalidYAML := `cv:
  name: [invalid
  yaml: syntax
`

	err := os.WriteFile(testFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = ReadResumeYAML(testFile)
	if err == nil {
		t.Error("ReadResumeYAML() expected error for invalid YAML, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid YAML") {
		t.Errorf("ReadResumeYAML() error should mention invalid YAML, got: %v", err)
	}
}

func TestWriteResumeYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.yaml")

	testContent := `cv:
  name: Jane Doe
  email: jane@example.com
`

	err := WriteResumeYAML(testFile, testContent)
	if err != nil {
		t.Fatalf("WriteResumeYAML() error = %v, want nil", err)
	}

	// Verify file was written
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(data) != testContent {
		t.Errorf("Written content does not match expected content")
	}
}

func TestWriteResumeYAML_InvalidPath(t *testing.T) {
	err := WriteResumeYAML("/nonexistent/dir/file.yaml", "cv:\n  name: Test")
	if err == nil {
		t.Error("WriteResumeYAML() expected error for invalid path, got nil")
	}
}

func TestWriteResumeYAML_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.yaml")

	invalidYAML := `cv:
  name: [invalid
  yaml: syntax
`

	err := WriteResumeYAML(testFile, invalidYAML)
	if err == nil {
		t.Error("WriteResumeYAML() expected error for invalid YAML, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "refusing to write invalid YAML") {
		t.Errorf("WriteResumeYAML() error should mention refusing to write, got: %v", err)
	}
}

func TestValidateYAML_Valid(t *testing.T) {
	validYAML := `cv:
  name: John Doe
  sections:
    experience:
      - company: TechCo
        position: Engineer
`

	err := ValidateYAML(validYAML)
	if err != nil {
		t.Errorf("ValidateYAML() error = %v, want nil", err)
	}
}

func TestValidateYAML_Invalid(t *testing.T) {
	invalidYAML := `cv:
  name: [invalid
  yaml: syntax
`

	err := ValidateYAML(invalidYAML)
	if err == nil {
		t.Error("ValidateYAML() expected error for invalid YAML, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "YAML parsing error") {
		t.Errorf("ValidateYAML() error should mention parsing error, got: %v", err)
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "empty string",
			content:  "",
			expected: 0,
		},
		{
			name:     "single line",
			content:  "line 1",
			expected: 1,
		},
		{
			name:     "multiple lines",
			content:  "line 1\nline 2\nline 3",
			expected: 3,
		},
		{
			name:     "trailing newline",
			content:  "line 1\nline 2\n",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CountLines(tt.content)
			if count != tt.expected {
				t.Errorf("CountLines() = %d, want %d", count, tt.expected)
			}
		})
	}
}

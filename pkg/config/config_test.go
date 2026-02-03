package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfig_Validate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")

	// Create test files
	os.WriteFile(resumeFile, []byte("cv:\n  name: Test"), 0644)
	os.WriteFile(jobFile, []byte("Job description"), 0644)

	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: resumeFile,
		JobPath:    jobFile,
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestConfig_Validate_MissingAPIKey(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")

	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:     "",
		ResumePath: resumeFile,
		JobPath:    jobFile,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() expected error for missing API key, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "ANTHROPIC_API_KEY is required") {
		t.Errorf("Validate() error should mention API key, got: %v", err)
	}
}

func TestConfig_Validate_MissingResume(t *testing.T) {
	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: "",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() expected error for missing resume, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "resume") {
		t.Errorf("Validate() error should mention resume, got: %v", err)
	}
}

func TestConfig_Validate_ResumeNotFound(t *testing.T) {
	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: "/nonexistent/resume.yaml",
		JobPath:    "/some/job.txt",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() expected error for nonexistent resume, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "resume file not found") {
		t.Errorf("Validate() error should mention resume not found, got: %v", err)
	}
}

func TestConfig_Validate_MissingJobInput(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	os.WriteFile(resumeFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: resumeFile,
		JobPath:    "",
		JobText:    "",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() expected error for missing job input, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "job") {
		t.Errorf("Validate() error should mention job, got: %v", err)
	}
}

func TestConfig_Validate_JobFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	os.WriteFile(resumeFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: resumeFile,
		JobPath:    "/nonexistent/job.txt",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() expected error for nonexistent job file, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "job description file not found") {
		t.Errorf("Validate() error should mention job file not found, got: %v", err)
	}
}

func TestConfig_Validate_JobTextProvided(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	os.WriteFile(resumeFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: resumeFile,
		JobText:    "Job description text",
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil for job text", err)
	}
}

func TestConfig_Validate_BothJobInputs(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")
	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:     "test-key",
		ResumePath: resumeFile,
		JobPath:    jobFile,
		JobText:    "Job text",
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}

	// JobText should be cleared (prefers JobPath)
	if cfg.JobText != "" {
		t.Error("Validate() should clear JobText when both are provided")
	}
}

func TestConfig_Validate_InvalidApprovalMode(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")
	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:       "test-key",
		ResumePath:   resumeFile,
		JobPath:      jobFile,
		ApprovalMode: "invalid-mode",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() expected error for invalid approval mode, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid approval mode") {
		t.Errorf("Validate() error should mention invalid approval mode, got: %v", err)
	}
}

func TestConfig_Validate_DefaultApprovalMode(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")
	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:       "test-key",
		ResumePath:   resumeFile,
		JobPath:      jobFile,
		ApprovalMode: "",
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}

	if cfg.ApprovalMode != "one-by-one" {
		t.Errorf("ApprovalMode = %s, want one-by-one", cfg.ApprovalMode)
	}
}

func TestConfig_Validate_ValidApprovalModes(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")
	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	modes := []string{"one-by-one", "batch", "auto-high"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			cfg := &Config{
				APIKey:       "test-key",
				ResumePath:   resumeFile,
				JobPath:      jobFile,
				ApprovalMode: mode,
			}

			err := cfg.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v for mode %s, want nil", err, mode)
			}
		})
	}
}

func TestConfig_Validate_CoverLetterAutoPath(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")
	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	cfg := &Config{
		APIKey:              "test-key",
		ResumePath:          resumeFile,
		JobPath:             jobFile,
		OutputPath:          "/path/to/output.yaml",
		GenerateCoverLetter: true,
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}

	if cfg.CoverLetterOutputPath == "" {
		t.Error("CoverLetterOutputPath should be auto-generated")
	}
	if !strings.Contains(cfg.CoverLetterOutputPath, "cover-letter") {
		t.Errorf("CoverLetterOutputPath = %s, should contain 'cover-letter'", cfg.CoverLetterOutputPath)
	}
}

func TestConfig_Validate_CoverLetterManualPath(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.yaml")
	jobFile := filepath.Join(tmpDir, "job.txt")
	os.WriteFile(resumeFile, []byte("test"), 0644)
	os.WriteFile(jobFile, []byte("test"), 0644)

	manualPath := "/custom/cover-letter.txt"
	cfg := &Config{
		APIKey:                "test-key",
		ResumePath:            resumeFile,
		JobPath:               jobFile,
		GenerateCoverLetter:   true,
		CoverLetterOutputPath: manualPath,
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}

	if cfg.CoverLetterOutputPath != manualPath {
		t.Errorf("CoverLetterOutputPath = %s, want %s", cfg.CoverLetterOutputPath, manualPath)
	}
}

func TestGenerateCoverLetterPath(t *testing.T) {
	tests := []struct {
		name       string
		resumePath string
		expected   string
	}{
		{
			name:       "yaml extension",
			resumePath: "/path/to/resume.yaml",
			expected:   "/path/to/resume-cover-letter.txt",
		},
		{
			name:       "yml extension",
			resumePath: "/path/to/resume.yml",
			expected:   "/path/to/resume-cover-letter.txt",
		},
		{
			name:       "no extension",
			resumePath: "/path/to/resume",
			expected:   "/path/to/resume-cover-letter.txt",
		},
		{
			name:       "multiple dots",
			resumePath: "/path/to/resume.backup.yaml",
			expected:   "/path/to/resume.backup-cover-letter.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateCoverLetterPath(tt.resumePath)
			if result != tt.expected {
				t.Errorf("generateCoverLetterPath() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestLoadAPIKey(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		envValue  string
		expected  string
	}{
		{
			name:      "flag provided",
			flagValue: "flag-key",
			envValue:  "env-key",
			expected:  "flag-key",
		},
		{
			name:      "only env provided",
			flagValue: "",
			envValue:  "env-key",
			expected:  "env-key",
		},
		{
			name:      "neither provided",
			flagValue: "",
			envValue:  "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env var
			os.Setenv("ANTHROPIC_API_KEY", tt.envValue)
			defer os.Unsetenv("ANTHROPIC_API_KEY")

			result := LoadAPIKey(tt.flagValue)
			if result != tt.expected {
				t.Errorf("LoadAPIKey() = %s, want %s", result, tt.expected)
			}
		})
	}
}

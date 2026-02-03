package coverletter

import (
	"fmt"
	"os"
	"strings"
	"tailor/internal/changes"
)

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteToFile(path string, data *changes.CoverLetterData) error {
	if data == nil {
		return fmt.Errorf("no cover letter data to write")
	}

	content := w.formatPlainText(data)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write cover letter file: %w", err)
	}

	return nil
}

func (w *Writer) formatPlainText(data *changes.CoverLetterData) string {
	var sb strings.Builder

	sb.WriteString("KEY QUALIFICATIONS\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")

	for i, point := range data.BulletPoints {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, point))
	}

	sb.WriteString("\n" + strings.Repeat("=", 50) + "\n\n")

	sb.WriteString("COVER LETTER\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")
	sb.WriteString(data.FullLetter)
	sb.WriteString("\n")

	return sb.String()
}

func (w *Writer) GetPreview(data *changes.CoverLetterData) string {
	if data == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n=== Cover Letter Preview ===\n\n")

	sb.WriteString("Key Qualifications:\n")
	for i, point := range data.BulletPoints {
		if i >= 5 {
			sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(data.BulletPoints)-5))
			break
		}
		sb.WriteString(fmt.Sprintf("  • %s\n", point))
	}

	sb.WriteString("\nCover Letter (first 400 chars):\n")
	preview := data.FullLetter
	if len(preview) > 400 {
		preview = preview[:400] + "..."
	}
	sb.WriteString(preview)
	sb.WriteString("\n")

	return sb.String()
}

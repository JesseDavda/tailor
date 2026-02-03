package executor

import (
	"fmt"
	"time"

	"tailor/internal/changes"
	"tailor/internal/claude"
	"tailor/internal/coverletter"
	"tailor/internal/interactive"
	"tailor/internal/job"
	"tailor/internal/resume"
)

func (e *Executor) loadResume() (string, error) {
	e.terminal.UpdateSpinner("Reading master resume")

	content, err := resume.ReadResumeYAML(e.config.ResumePath)
	if err != nil {
		e.terminal.Fail("failed to read resume: " + err.Error())
		return "", fmt.Errorf("failed to read resume: %w", err)
	}

	lines := resume.CountLines(content)
	e.terminal.Success(fmt.Sprintf("Read master resume (%d lines)", lines))
	return content, nil
}

func (e *Executor) loadJobDescription() (string, error) {
	e.terminal.UpdateSpinner("Reading job description")

	var content string
	var err error

	if e.config.JobPath != "" {
		content, err = job.ReadJobDescription(e.config.JobPath)
	} else {
		content, err = job.ParseJobText(e.config.JobText)
	}

	if err != nil {
		e.terminal.Fail("failed to load job description: " + err.Error())
		return "", fmt.Errorf("failed to load job description: %w", err)
	}

	e.terminal.Success("Read job description")
	return content, nil
}

func (e *Executor) analyzeResume(masterYAML, jobDesc string) (*changes.ChangeSet, error) {
	e.terminal.UpdateSpinner("Analyzing resume with Claude API")
	start := time.Now()

	client := claude.NewClient(e.config.APIKey, e.config.Model)
	changeSet, err := client.TailorResumeInteractive(
		e.ctx,
		masterYAML,
		jobDesc,
		e.config.GenerateCoverLetter,
	)
	if err != nil {
		e.terminal.Fail("failed to analyze resume: " + err.Error())
		return nil, fmt.Errorf("failed to analyze resume: %w", err)
	}

	elapsed := time.Since(start)
	e.terminal.Success(fmt.Sprintf("Analysis complete (took %v)", elapsed.Round(time.Millisecond)))
	return changeSet, nil
}

func (e *Executor) reviewChanges(changeSet *changes.ChangeSet) error {
	session := e.terminal.PauseForInteractive()

	e.terminal.Section("CHANGE REVIEW")

	approver := interactive.NewApprover(e.config.ApprovalMode, session)
	if err := approver.ReviewChanges(changeSet); err != nil {
		session.Resume("")
		e.terminal.Error("review failed: " + err.Error())
		return fmt.Errorf("review failed: %w", err)
	}

	session.Resume("Review complete")
	return nil
}

func (e *Executor) applyChanges(masterYAML string, changeSet *changes.ChangeSet) (string, error) {
	e.terminal.UpdateSpinner("Parsing YAML document")
	doc, err := resume.ParseYAMLDocument(masterYAML)
	if err != nil {
		e.terminal.Fail("failed to parse YAML: " + err.Error())
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}
	e.terminal.Success("YAML parsed")

	e.terminal.UpdateSpinner("Validating changes")
	applier := changes.NewApplier(e.terminal.GetLogger())
	if err := applier.ValidateBeforeApply(doc, changeSet); err != nil {
		e.terminal.Fail("validation failed: " + err.Error())
		return "", fmt.Errorf("validation failed: %w", err)
	}
	e.terminal.Success("Changes validated")

	e.terminal.UpdateSpinner("Applying approved changes")
	result, err := applier.ApplyChanges(doc, changeSet)
	if err != nil {
		e.terminal.Fail("failed to apply changes: " + err.Error())
		return "", fmt.Errorf("failed to apply changes: %w", err)
	}

	if result.FailureCount > 0 {
		e.terminal.Warn(fmt.Sprintf("%d changes failed to apply", result.FailureCount))
		for _, err := range result.Errors {
			e.terminal.Error(err.Error())
		}
	}

	e.terminal.Success(fmt.Sprintf("Applied %d changes", result.SuccessCount))

	e.terminal.UpdateSpinner("Generating tailored YAML")
	tailoredYAML, err := doc.ToYAML()
	if err != nil {
		e.terminal.Fail("failed to generate YAML: " + err.Error())
		return "", fmt.Errorf("failed to generate YAML: %w", err)
	}
	e.terminal.Success("YAML generated")

	return tailoredYAML, nil
}

func (e *Executor) writeOutput(content string) error {
	if e.config.DryRun {
		e.terminal.UpdateSpinner("Dry run mode - skipping write")
		e.terminal.Section("DRY RUN MODE")
		e.terminal.Info(fmt.Sprintf("Would write to: %s", e.config.OutputPath))
		e.terminal.Info(fmt.Sprintf("Content length: %d lines", resume.CountLines(content)))
		return nil
	}

	e.terminal.UpdateSpinner(fmt.Sprintf("Writing to %s", e.config.OutputPath))
	if err := resume.WriteResumeYAML(e.config.OutputPath, content); err != nil {
		e.terminal.Fail("failed to write output: " + err.Error())
		return fmt.Errorf("failed to write output: %w", err)
	}
	e.terminal.Success("Tailored resume written")

	return nil
}

func (e *Executor) writeCoverLetter(changeSet *changes.ChangeSet) error {
	if !changeSet.HasCoverLetter() {
		return nil
	}

	if e.config.DryRun {
		e.terminal.Section("COVER LETTER PREVIEW (DRY RUN)")
		writer := coverletter.NewWriter()
		preview := writer.GetPreview(changeSet.CoverLetter)
		fmt.Println(preview)
		e.terminal.Info(fmt.Sprintf("Would write to: %s", e.config.CoverLetterOutputPath))
		return nil
	}

	e.terminal.UpdateSpinner(fmt.Sprintf("Writing cover letter to %s", e.config.CoverLetterOutputPath))

	writer := coverletter.NewWriter()
	if err := writer.WriteToFile(e.config.CoverLetterOutputPath, changeSet.CoverLetter); err != nil {
		e.terminal.Fail("failed to write cover letter: " + err.Error())
		return fmt.Errorf("failed to write cover letter: %w", err)
	}

	e.terminal.Success("Cover letter written")
	return nil
}

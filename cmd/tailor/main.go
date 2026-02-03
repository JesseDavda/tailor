package main

import (
	"fmt"
	"os"

	"tailor/internal/executor"
	cfg "tailor/pkg/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tailor",
	Short: "Tailor your RenderCV resume to job descriptions using Claude AI",
	Long:  `Interactive resume tailoring with Claude AI - maintain factual accuracy while optimizing for ATS.`,
	Example: `  # Basic usage
  tailor --resume master.yaml --job job-desc.txt

  # Generate cover letter
  tailor -r master.yaml -j job.txt --cover-letter

  # Custom cover letter output
  tailor -r master.yaml -j job.txt -c --cover-letter-output letter.txt

  # Batch approval by section
  tailor -r master.yaml -j job.txt --approval-mode batch

  # Dry run (preview without writing)
  tailor -r master.yaml -j job.txt --dry-run`,
	RunE: runTailor,
}

func init() {
	// Input flags
	rootCmd.Flags().StringVarP(&cfg.Global.ResumePath, "resume", "r", "",
		"Path to master RenderCV YAML resume (required)")
	rootCmd.Flags().StringVarP(&cfg.Global.JobPath, "job", "j", "",
		"Path to job description file")
	rootCmd.Flags().StringVar(&cfg.Global.JobText, "job-text", "",
		"Direct job description text input")

	// Output flags
	rootCmd.Flags().StringVarP(&cfg.Global.OutputPath, "output", "o",
		"tailored-resume.yaml", "Output path for tailored YAML")

	// Cover letter flags
	rootCmd.Flags().BoolVarP(&cfg.Global.GenerateCoverLetter, "cover-letter", "c", false,
		"Generate cover letter with key qualifications and full text")
	rootCmd.Flags().StringVar(&cfg.Global.CoverLetterOutputPath, "cover-letter-output", "",
		"Custom output path for cover letter (default: derived from resume output)")

	// Configuration flags
	apiKeyFlag := ""
	rootCmd.Flags().StringVar(&apiKeyFlag, "api-key", "",
		"Anthropic API key (or use ANTHROPIC_API_KEY env var)")
	rootCmd.Flags().StringVarP(&cfg.Global.Model, "model", "m", "", "Claude model to use, list of models can be found here: https://platform.claude.com/docs/en/about-claude/models/overview")

	// Runtime flags
	rootCmd.Flags().BoolVar(&cfg.Global.DryRun, "dry-run", false,
		"Preview without writing output")
	rootCmd.Flags().StringVar(&cfg.Global.ApprovalMode, "approval-mode",
		"one-by-one", "Approval mode: one-by-one, batch, auto-high")

	rootCmd.MarkFlagRequired("resume")

	// Load API key after flags are parsed
	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		apiKeyFlagValue, _ := cmd.Flags().GetString("api-key")
		cfg.Global.APIKey = cfg.LoadAPIKey(apiKeyFlagValue)
		return nil
	}
}

func runTailor(cmd *cobra.Command, args []string) error {
	if err := cfg.Global.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	exec := executor.New(cfg.Global)
	return exec.Run()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package claude

import (
	"context"
	"fmt"
	"tailor/internal/changes"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Client wraps the Anthropic API client
type Client struct {
	client *anthropic.Client
	model  string
}

// NewClient creates a new Claude API client
func NewClient(apiKey string, model string) *Client {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &Client{
		client: &client,
		model:  model,
	}
}

// TailorResumeInteractive analyzes the resume and returns a set of proposed changes
func (c *Client) TailorResumeInteractive(ctx context.Context, masterYAML, jobDescription string) (*changes.ChangeSet, error) {
	userPrompt := BuildInteractiveUserPrompt(masterYAML, jobDescription)

	// Create the message request
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 8000,
		System: []anthropic.TextBlockParam{
			{
				Text: SystemPromptInteractive,
				Type: "text",
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	}

	message, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Claude API error: %w", err)
	}

	// Extract the response text
	if len(message.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude API")
	}

	var responseText string
	for _, block := range message.Content {
		switch block.Type {
		case "text":
			responseText += block.Text
		}
	}

	if responseText == "" {
		return nil, fmt.Errorf("no text content in Claude API response")
	}

	// Parse the JSON response into a ChangeSet
	changeSet, err := changes.ParseChangeSet(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse changeset: %w", err)
	}

	return changeSet, nil
}

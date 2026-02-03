package changes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseChangeSet parses Claude's JSON response into a ChangeSet
func ParseChangeSet(jsonResponse string) (*ChangeSet, error) {
	// Clean markdown code blocks if present
	cleaned := cleanMarkdownCodeBlocks(jsonResponse)

	var changeSet ChangeSet
	if err := json.Unmarshal([]byte(cleaned), &changeSet); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Initialize all changes to pending status
	for i := range changeSet.Changes {
		changeSet.Changes[i].Status = StatusPending
	}

	// Validate the changeset
	if err := validateChangeSet(&changeSet); err != nil {
		return nil, fmt.Errorf("invalid changeset: %w", err)
	}

	return &changeSet, nil
}

// cleanMarkdownCodeBlocks removes markdown code block markers if present
func cleanMarkdownCodeBlocks(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}

	content = strings.TrimSuffix(content, "```")

	return strings.TrimSpace(content)
}

// validateChangeSet validates the structure and content of a changeset
func validateChangeSet(cs *ChangeSet) error {
	if len(cs.Changes) == 0 {
		return fmt.Errorf("changeset contains no changes")
	}

	for i, change := range cs.Changes {
		if err := validateChange(&change); err != nil {
			return fmt.Errorf("invalid change at index %d: %w", i, err)
		}
	}

	return nil
}

// validateChange validates a single change
func validateChange(c *Change) error {
	// Validate ID
	if c.ID == "" {
		return fmt.Errorf("change ID is required")
	}

	// Validate type
	switch c.Type {
	case ChangeTypeModify, ChangeTypeReorder, ChangeTypeAdd, ChangeTypeRemove:
		// Valid type
	default:
		return fmt.Errorf("invalid change type: %s", c.Type)
	}

	// Validate path
	if c.Path == "" {
		return fmt.Errorf("change path is required")
	}

	// Validate operation based on type
	switch c.Type {
	case ChangeTypeModify:
		if c.Operation.NewValue == nil {
			return fmt.Errorf("modify operation requires new_value")
		}
	case ChangeTypeReorder:
		if len(c.Operation.NewOrder) == 0 {
			return fmt.Errorf("reorder operation requires new_order")
		}
	case ChangeTypeAdd:
		if c.Operation.NewValue == nil {
			return fmt.Errorf("add operation requires new_value")
		}
	case ChangeTypeRemove:
		// Remove operations are valid with just a path
	}

	// Validate confidence
	confidence := strings.ToLower(c.Confidence)
	if confidence != "high" && confidence != "medium" && confidence != "low" {
		return fmt.Errorf("invalid confidence level: %s (must be high, medium, or low)", c.Confidence)
	}

	// Validate priority
	if c.Priority < 1 || c.Priority > 3 {
		return fmt.Errorf("invalid priority: %d (must be 1, 2, or 3)", c.Priority)
	}

	return nil
}

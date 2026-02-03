package changes

// ChangeType represents the type of change to apply
type ChangeType string

const (
	ChangeTypeModify  ChangeType = "modify"
	ChangeTypeReorder ChangeType = "reorder"
	ChangeTypeAdd     ChangeType = "add"
	ChangeTypeRemove  ChangeType = "remove"
)

// ChangeStatus represents the approval status of a change
type ChangeStatus string

const (
	StatusPending  ChangeStatus = "pending"
	StatusApproved ChangeStatus = "approved"
	StatusRejected ChangeStatus = "rejected"
	StatusSkipped  ChangeStatus = "skipped"
)

// ChangeSet represents a collection of changes with summary
type ChangeSet struct {
	Changes []Change      `json:"changes"`
	Summary ChangeSummary `json:"summary"`
}

// Change represents a single granular change to the resume
type Change struct {
	ID         string       `json:"id"`
	Type       ChangeType   `json:"type"`
	Path       string       `json:"path"`
	Operation  Operation    `json:"operation"`
	Reason     string       `json:"reason"`
	Confidence string       `json:"confidence"` // high, medium, low
	Priority   int          `json:"priority"`   // 1 (critical), 2 (important), 3 (minor)
	Status     ChangeStatus `json:"-"`          // Internal status, not from JSON
}

// Operation contains the details of what to change
type Operation struct {
	OldValue interface{}   `json:"old_value,omitempty"`
	NewValue interface{}   `json:"new_value,omitempty"`
	OldOrder []int         `json:"old_order,omitempty"` // for reorder type
	NewOrder []int         `json:"new_order,omitempty"` // for reorder type
	Position *int          `json:"position,omitempty"`  // for add type (where to insert)
	Index    *int          `json:"index,omitempty"`     // for remove type (what to remove)
}

// ChangeSummary provides high-level overview of the changeset
type ChangeSummary struct {
	TotalChanges     int      `json:"total_changes"`
	SectionsAffected []string `json:"sections_affected"`
	KeyOptimizations []string `json:"key_optimizations"`
}

// IsApproved returns true if the change has been approved
func (c *Change) IsApproved() bool {
	return c.Status == StatusApproved
}

// IsPending returns true if the change is still pending review
func (c *Change) IsPending() bool {
	return c.Status == StatusPending
}

// Approve marks the change as approved
func (c *Change) Approve() {
	c.Status = StatusApproved
}

// Reject marks the change as rejected
func (c *Change) Reject() {
	c.Status = StatusRejected
}

// Skip marks the change as skipped
func (c *Change) Skip() {
	c.Status = StatusSkipped
}

// GetApprovedChanges returns only the approved changes from a changeset
func (cs *ChangeSet) GetApprovedChanges() []Change {
	approved := []Change{}
	for _, change := range cs.Changes {
		if change.IsApproved() {
			approved = append(approved, change)
		}
	}
	return approved
}

// CountByStatus returns counts of changes by status
func (cs *ChangeSet) CountByStatus() map[ChangeStatus]int {
	counts := map[ChangeStatus]int{
		StatusPending:  0,
		StatusApproved: 0,
		StatusRejected: 0,
		StatusSkipped:  0,
	}
	for _, change := range cs.Changes {
		counts[change.Status]++
	}
	return counts
}

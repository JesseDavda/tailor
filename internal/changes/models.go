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

// CoverLetterData contains AI-generated cover letter content
type CoverLetterData struct {
	BulletPoints []string `json:"bullet_points"`
	FullLetter   string   `json:"full_letter"`
}

// ChangeSet represents a collection of changes with summary
type ChangeSet struct {
	Changes     []Change         `json:"changes"`
	Summary     ChangeSummary    `json:"summary"`
	CoverLetter *CoverLetterData `json:"cover_letter,omitempty"`
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
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
	OldOrder []int       `json:"old_order,omitempty"` // for reorder type
	NewOrder []int       `json:"new_order,omitempty"` // for reorder type
	Position *int        `json:"position,omitempty"`  // for add type (where to insert)
	Index    *int        `json:"index,omitempty"`     // for remove type (what to remove)
}

type ChangeSummary struct {
	TotalChanges     int      `json:"total_changes"`
	SectionsAffected []string `json:"sections_affected"`
	KeyOptimizations []string `json:"key_optimizations"`
}

func (c *Change) IsApproved() bool {
	return c.Status == StatusApproved
}

func (c *Change) IsPending() bool {
	return c.Status == StatusPending
}

func (c *Change) Approve() {
	c.Status = StatusApproved
}

func (c *Change) Reject() {
	c.Status = StatusRejected
}

func (c *Change) Skip() {
	c.Status = StatusSkipped
}

func (cs *ChangeSet) GetApprovedChanges() []Change {
	approved := []Change{}
	for _, change := range cs.Changes {
		if change.IsApproved() {
			approved = append(approved, change)
		}
	}
	return approved
}

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

func (cs *ChangeSet) HasCoverLetter() bool {
	return cs.CoverLetter != nil &&
		len(cs.CoverLetter.BulletPoints) > 0 &&
		cs.CoverLetter.FullLetter != ""
}

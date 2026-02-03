package changes

import "testing"

func TestChange_IsApproved(t *testing.T) {
	tests := []struct {
		name     string
		status   ChangeStatus
		expected bool
	}{
		{"approved change", StatusApproved, true},
		{"pending change", StatusPending, false},
		{"rejected change", StatusRejected, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := Change{Status: tt.status}
			if got := change.IsApproved(); got != tt.expected {
				t.Errorf("IsApproved() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestChange_Approve(t *testing.T) {
	change := Change{Status: StatusPending}
	change.Approve()

	if change.Status != StatusApproved {
		t.Errorf("Approve() did not set status to approved, got %v", change.Status)
	}
}

func TestChange_Reject(t *testing.T) {
	change := Change{Status: StatusPending}
	change.Reject()

	if change.Status != StatusRejected {
		t.Errorf("Reject() did not set status to rejected, got %v", change.Status)
	}
}

func TestChangeSet_GetApprovedChanges(t *testing.T) {
	tests := []struct {
		name          string
		changes       []Change
		expectedCount int
	}{
		{
			name: "mixed statuses",
			changes: []Change{
				{ID: "1", Status: StatusApproved},
				{ID: "2", Status: StatusPending},
				{ID: "3", Status: StatusApproved},
				{ID: "4", Status: StatusRejected},
			},
			expectedCount: 2,
		},
		{
			name: "all approved",
			changes: []Change{
				{ID: "1", Status: StatusApproved},
				{ID: "2", Status: StatusApproved},
			},
			expectedCount: 2,
		},
		{
			name: "none approved",
			changes: []Change{
				{ID: "1", Status: StatusPending},
				{ID: "2", Status: StatusRejected},
			},
			expectedCount: 0,
		},
		{
			name:          "empty changeset",
			changes:       []Change{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChangeSet{Changes: tt.changes}
			approved := cs.GetApprovedChanges()

			if len(approved) != tt.expectedCount {
				t.Errorf("GetApprovedChanges() returned %d changes, want %d", len(approved), tt.expectedCount)
			}

			// Verify all returned changes are actually approved
			for _, change := range approved {
				if change.Status != StatusApproved {
					t.Errorf("GetApprovedChanges() returned non-approved change with status %v", change.Status)
				}
			}
		})
	}
}

func TestChangeSet_CountByStatus(t *testing.T) {
	tests := []struct {
		name     string
		changes  []Change
		expected map[ChangeStatus]int
	}{
		{
			name: "mixed statuses",
			changes: []Change{
				{ID: "1", Status: StatusApproved},
				{ID: "2", Status: StatusPending},
				{ID: "3", Status: StatusApproved},
				{ID: "4", Status: StatusRejected},
				{ID: "5", Status: StatusPending},
			},
			expected: map[ChangeStatus]int{
				StatusApproved: 2,
				StatusPending:  2,
				StatusRejected: 1,
			},
		},
		{
			name: "all same status",
			changes: []Change{
				{ID: "1", Status: StatusApproved},
				{ID: "2", Status: StatusApproved},
			},
			expected: map[ChangeStatus]int{
				StatusApproved: 2,
			},
		},
		{
			name:     "empty changeset",
			changes:  []Change{},
			expected: map[ChangeStatus]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChangeSet{Changes: tt.changes}
			counts := cs.CountByStatus()

			// Check all expected counts
			for status, expectedCount := range tt.expected {
				if got := counts[status]; got != expectedCount {
					t.Errorf("CountByStatus()[%v] = %d, want %d", status, got, expectedCount)
				}
			}

			// Ensure no unexpected statuses with non-zero counts
			for status, count := range counts {
				if count > 0 && tt.expected[status] != count {
					t.Errorf("CountByStatus()[%v] = %d, but expected %d", status, count, tt.expected[status])
				}
			}
		})
	}
}

func TestChangeSet_HasCoverLetter(t *testing.T) {
	tests := []struct {
		name        string
		coverLetter *CoverLetterData
		expected    bool
	}{
		{
			name: "has cover letter",
			coverLetter: &CoverLetterData{
				BulletPoints: []string{
					"Point 1",
					"Point 2",
				},
				FullLetter: "This is a full letter",
			},
			expected: true,
		},
		{
			name: "empty bullet points",
			coverLetter: &CoverLetterData{
				BulletPoints: []string{},
				FullLetter:   "This is a full letter",
			},
			expected: false,
		},
		{
			name: "empty full letter",
			coverLetter: &CoverLetterData{
				BulletPoints: []string{"Point 1"},
				FullLetter:   "",
			},
			expected: false,
		},
		{
			name:        "nil cover letter",
			coverLetter: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChangeSet{CoverLetter: tt.coverLetter}
			if got := cs.HasCoverLetter(); got != tt.expected {
				t.Errorf("HasCoverLetter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

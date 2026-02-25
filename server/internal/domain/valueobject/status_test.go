package valueobject

import (
	"testing"
)

func TestNewPostStatus(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    PostStatus
		wantErr bool
	}{
		{"draft", "draft", StatusDraft, false},
		{"published", "published", StatusPublished, false},
		{"empty string defaults to draft", "", StatusDraft, false},
		{"invalid status", "invalid", "", true},
		{"unknown status", "archived", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPostStatus(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostStatus(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPostStatus(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestPostStatusIsPublished(t *testing.T) {
	tests := []struct {
		status   PostStatus
		expected bool
	}{
		{StatusDraft, false},
		{StatusPublished, true},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.IsPublished(); got != tt.expected {
				t.Errorf("IsPublished() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPostStatusIsDraft(t *testing.T) {
	tests := []struct {
		status   PostStatus
		expected bool
	}{
		{StatusDraft, true},
		{StatusPublished, false},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.IsDraft(); got != tt.expected {
				t.Errorf("IsDraft() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPostStatusCanPublish(t *testing.T) {
	tests := []struct {
		status   PostStatus
		expected bool
	}{
		{StatusDraft, true},
		{StatusPublished, false},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.CanPublish(); got != tt.expected {
				t.Errorf("CanPublish() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPostStatusCanUnpublish(t *testing.T) {
	tests := []struct {
		status   PostStatus
		expected bool
	}{
		{StatusDraft, false},
		{StatusPublished, true},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.CanUnpublish(); got != tt.expected {
				t.Errorf("CanUnpublish() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPostStatusString(t *testing.T) {
	if StatusDraft.String() != "draft" {
		t.Errorf("String() = %q, want %q", StatusDraft.String(), "draft")
	}
	if StatusPublished.String() != "published" {
		t.Errorf("String() = %q, want %q", StatusPublished.String(), "published")
	}
}

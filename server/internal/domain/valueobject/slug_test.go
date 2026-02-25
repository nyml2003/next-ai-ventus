package valueobject

import (
	"testing"
)

func TestSlugCreate(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"valid lowercase with hyphen", "hello-world", false},
		{"valid with numbers", "hello-world-123", false},
		{"valid single word", "hello", false},
		{"invalid chars at", "hello@world", true},
		{"invalid space", "hello world", true},
		{"invalid uppercase", "Hello-World", true},
		{"invalid underscore", "hello_world", true},
		{"empty string", "", true},
		{"starts with hyphen", "-hello", true},
		{"ends with hyphen", "hello-", true},
		{"consecutive hyphens", "hello--world", true},
		{"chinese characters", "你好世界", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := NewSlug(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSlug(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
				return
			}
			if !tt.wantErr && slug.String() != tt.raw {
				t.Errorf("NewSlug(%q) = %q, want %q", tt.raw, slug.String(), tt.raw)
			}
		})
	}
}

func TestSlugGenerateFromTitle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		existing []string
		want     string
	}{
		{
			name:     "simple title",
			title:    "Hello World",
			existing: []string{},
			want:     "hello-world",
		},
		{
			name:     "title with special chars",
			title:    "Hello, World! (2024)",
			existing: []string{},
			want:     "hello-world-2024",
		},
		{
			name:     "chinese title",
			title:    "你好世界",
			existing: []string{},
			want:     "ni-hao-shi-jie",
		},
		{
			name:     "no conflict",
			title:    "Hello World",
			existing: []string{"other-slug"},
			want:     "hello-world",
		},
		{
			name:     "conflict once",
			title:    "Hello World",
			existing: []string{"hello-world"},
			want:     "hello-world-2",
		},
		{
			name:     "conflict twice",
			title:    "Hello World",
			existing: []string{"hello-world", "hello-world-2"},
			want:     "hello-world-3",
		},
		{
			name:     "conflict with gap",
			title:    "Hello World",
			existing: []string{"hello-world", "hello-world-2", "hello-world-4"},
			want:     "hello-world-3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := GenerateFromTitle(tt.title, tt.existing)
			if slug.String() != tt.want {
				t.Errorf("GenerateFromTitle(%q, %v) = %q, want %q",
					tt.title, tt.existing, slug.String(), tt.want)
			}
		})
	}
}

func TestSlugEquals(t *testing.T) {
	slug1, _ := NewSlug("hello-world")
	slug2, _ := NewSlug("hello-world")
	slug3, _ := NewSlug("different-slug")

	if !slug1.Equals(slug2) {
		t.Error("slug1 should equal slug2")
	}

	if slug1.Equals(slug3) {
		t.Error("slug1 should not equal slug3")
	}
}

func TestSlugString(t *testing.T) {
	slug, err := NewSlug("test-slug")
	if err != nil {
		t.Fatalf("NewSlug failed: %v", err)
	}

	if slug.String() != "test-slug" {
		t.Errorf("String() = %q, want %q", slug.String(), "test-slug")
	}
}

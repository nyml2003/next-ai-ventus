package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/next-ai-ventus/server/internal/domain/valueobject"
)

func TestNewPost(t *testing.T) {
	slug, _ := valueobject.NewSlug("hello-world")

	tests := []struct {
		name    string
		id      string
		title   string
		slug    valueobject.Slug
		content string
		tags    []valueobject.Tag
		wantErr bool
	}{
		{
			name:    "valid post",
			id:      "2024-06-hello",
			title:   "Hello World",
			slug:    slug,
			content: "# Hello\n\nThis is content.",
			tags:    []valueobject.Tag{},
			wantErr: false,
		},
		{
			name:    "empty title",
			id:      "2024-06-test",
			title:   "",
			slug:    slug,
			content: "Content",
			tags:    []valueobject.Tag{},
			wantErr: true,
		},
		{
			name:    "empty content",
			id:      "2024-06-test",
			title:   "Title",
			slug:    slug,
			content: "",
			tags:    []valueobject.Tag{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(tt.id, tt.title, tt.slug, tt.content, tt.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if post.ID != tt.id {
					t.Errorf("ID = %q, want %q", post.ID, tt.id)
				}
				if post.Title != tt.title {
					t.Errorf("Title = %q, want %q", post.Title, tt.title)
				}
				if post.Status != valueobject.StatusDraft {
					t.Errorf("Status = %v, want Draft", post.Status)
				}
				if post.Version != 1 {
					t.Errorf("Version = %d, want 1", post.Version)
				}
				if post.Excerpt == "" {
					t.Error("Excerpt should be generated")
				}
			}
		})
	}
}

func TestPostPublish(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	post, _ := NewPost("1", "Test", slug, "Content", nil)

	t.Run("publish draft", func(t *testing.T) {
		err := post.Publish()
		if err != nil {
			t.Errorf("Publish() error = %v", err)
		}
		if !post.IsPublished() {
			t.Error("Post should be published")
		}
		if post.PublishedAt == nil {
			t.Error("PublishedAt should be set")
		}
		if post.Version != 2 {
			t.Errorf("Version = %d, want 2", post.Version)
		}
	})

	t.Run("publish already published", func(t *testing.T) {
		err := post.Publish()
		if err != ErrAlreadyPublished {
			t.Errorf("Publish() error = %v, want ErrAlreadyPublished", err)
		}
	})
}

func TestPostUnpublish(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	post, _ := NewPost("1", "Test", slug, "Content", nil)
	post.Publish()

	t.Run("unpublish published", func(t *testing.T) {
		err := post.Unpublish()
		if err != nil {
			t.Errorf("Unpublish() error = %v", err)
		}
		if post.IsPublished() {
			t.Error("Post should be draft")
		}
		if post.PublishedAt != nil {
			t.Error("PublishedAt should be nil")
		}
	})

	t.Run("unpublish draft", func(t *testing.T) {
		err := post.Unpublish()
		if err != ErrNotPublished {
			t.Errorf("Unpublish() error = %v, want ErrNotPublished", err)
		}
	})
}

func TestPostUpdateContent(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	post, _ := NewPost("1", "Test", slug, "Old content", nil)
	oldVersion := post.Version

	t.Run("update content", func(t *testing.T) {
		err := post.UpdateContent("New content")
		if err != nil {
			t.Errorf("UpdateContent() error = %v", err)
		}
		if post.Content != "New content" {
			t.Errorf("Content = %q, want %q", post.Content, "New content")
		}
		if post.Version != oldVersion+1 {
			t.Errorf("Version = %d, want %d", post.Version, oldVersion+1)
		}
	})

	t.Run("update with empty content", func(t *testing.T) {
		err := post.UpdateContent("")
		if err != ErrEmptyContent {
			t.Errorf("UpdateContent() error = %v, want ErrEmptyContent", err)
		}
	})
}

func TestPostUpdateTitle(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	post, _ := NewPost("1", "Old Title", slug, "Content", nil)
	oldVersion := post.Version

	t.Run("update title", func(t *testing.T) {
		err := post.UpdateTitle("New Title")
		if err != nil {
			t.Errorf("UpdateTitle() error = %v", err)
		}
		if post.Title != "New Title" {
			t.Errorf("Title = %q, want %q", post.Title, "New Title")
		}
		if post.Version != oldVersion+1 {
			t.Errorf("Version = %d, want %d", post.Version, oldVersion+1)
		}
	})

	t.Run("update with empty title", func(t *testing.T) {
		err := post.UpdateTitle("")
		if err != ErrEmptyTitle {
			t.Errorf("UpdateTitle() error = %v, want ErrEmptyTitle", err)
		}
	})
}

func TestPostUpdateTags(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	post, _ := NewPost("1", "Test", slug, "Content", nil)
	
	tag1, _ := valueobject.NewTag("go")
	tag2, _ := valueobject.NewTag("rust")
	newTags := []valueobject.Tag{tag1, tag2}

	oldVersion := post.Version
	post.UpdateTags(newTags)

	if len(post.Tags) != 2 {
		t.Errorf("Tags length = %d, want 2", len(post.Tags))
	}
	if post.Version != oldVersion+1 {
		t.Errorf("Version = %d, want %d", post.Version, oldVersion+1)
	}
}

func TestPostGetTagNames(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	tag1, _ := valueobject.NewTag("go")
	tag2, _ := valueobject.NewTag("rust")
	post, _ := NewPost("1", "Test", slug, "Content", []valueobject.Tag{tag1, tag2})

	names := post.GetTagNames()
	if len(names) != 2 {
		t.Errorf("GetTagNames() length = %d, want 2", len(names))
	}
	if names[0] != "go" || names[1] != "rust" {
		t.Errorf("GetTagNames() = %v, want [go rust]", names)
	}
}

func TestPostHasTag(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	tag1, _ := valueobject.NewTag("go")
	post, _ := NewPost("1", "Test", slug, "Content", []valueobject.Tag{tag1})

	if !post.HasTag("go") {
		t.Error("HasTag(go) should be true")
	}
	if post.HasTag("rust") {
		t.Error("HasTag(rust) should be false")
	}
}

func TestPostGenerateExcerpt(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	
	tests := []struct {
		name     string
		content  string
		maxLen   int
		expected string
	}{
		{
			name:     "short content",
			content:  "Short content",
			maxLen:   100,
			expected: "Short content",
		},
		{
			name:     "long content truncated",
			content:  strings.Repeat("a", 250),
			maxLen:   100,
			expected: strings.Repeat("a", 100) + "...",
		},
		{
			name:     "content with markdown",
			content:  "# Title\n\n**Bold** text",
			maxLen:   100,
			expected: "Title\n\nBold text", // 简化实现保留换行
		},
		{
			name:     "content with link",
			content:  "Check [this link](http://example.com) out",
			maxLen:   100,
			expected: "Check this link out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, _ := NewPost("1", "Test", slug, tt.content, nil)
			post.GenerateExcerpt(tt.maxLen)
			if post.Excerpt != tt.expected {
				t.Errorf("Excerpt = %q, want %q", post.Excerpt, tt.expected)
			}
		})
	}
}

func TestPostTimestamps(t *testing.T) {
	slug, _ := valueobject.NewSlug("test-post")
	before := time.Now()
	post, _ := NewPost("1", "Test", slug, "Content", nil)
	after := time.Now()

	if post.CreatedAt.Before(before) || post.CreatedAt.After(after) {
		t.Error("CreatedAt should be set to current time")
	}
	if post.UpdatedAt.Before(before) || post.UpdatedAt.After(after) {
		t.Error("UpdatedAt should be set to current time")
	}
}

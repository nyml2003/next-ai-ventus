package valueobject

import (
	"testing"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"valid lowercase", "go", false},
		{"valid with hyphen", "web-development", false},
		{"valid with number", "go123", false},
		{"invalid uppercase", "Go", true},
		{"invalid space", "go lang", true},
		{"invalid special char", "go@lang", true},
		{"empty string", "", true},
		{"too long", "this-is-a-very-long-tag-name-that-exceeds-the-limit", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, err := NewTag(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTag(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
				return
			}
			if !tt.wantErr && tag.String() != tt.raw {
				t.Errorf("NewTag(%q) = %q, want %q", tt.raw, tag.String(), tt.raw)
			}
		})
	}
}

func TestTagString(t *testing.T) {
	tag, _ := NewTag("golang")
	if tag.String() != "golang" {
		t.Errorf("String() = %q, want %q", tag.String(), "golang")
	}
}

func TestTagEquals(t *testing.T) {
	tag1, _ := NewTag("go")
	tag2, _ := NewTag("go")
	tag3, _ := NewTag("python")

	if !tag1.Equals(tag2) {
		t.Error("tag1 should equal tag2")
	}
	if tag1.Equals(tag3) {
		t.Error("tag1 should not equal tag3")
	}
}

func TestNormalizeTags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "deduplicate",
			input:    []string{"go", "python", "go", "rust"},
			expected: []string{"go", "python", "rust"},
		},
		{
			name:     "sort alphabetically",
			input:    []string{"zebra", "apple", "banana"},
			expected: []string{"apple", "banana", "zebra"},
		},
		{
			name:     "dedup and sort",
			input:    []string{"zebra", "apple", "zebra", "banana", "apple"},
			expected: []string{"apple", "banana", "zebra"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"solo"},
			expected: []string{"solo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("NormalizeTags(%v) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("NormalizeTags(%v) = %v, want %v", tt.input, result, tt.expected)
					return
				}
			}
		})
	}
}

func TestTagValidationError(t *testing.T) {
	_, err := NewTag("InvalidTag")
	if err == nil {
		t.Error("expected error for uppercase tag")
	}
}

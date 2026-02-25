package markdown

import (
	"strings"
	"testing"
)

func TestExtractTOC(t *testing.T) {
	content := `# Title 1
Some content
## Subtitle 1.1
More content
# Title 2
Even more content
### Deep section
`

	toc := ExtractTOC(content)

	if len(toc) != 2 {
		t.Errorf("len(toc) = %d, want 2", len(toc))
	}

	t.Run("first level", func(t *testing.T) {
		if toc[0].Text != "Title 1" {
			t.Errorf("toc[0].Text = %q, want Title 1", toc[0].Text)
		}
		if toc[0].Level != 1 {
			t.Errorf("toc[0].Level = %d, want 1", toc[0].Level)
		}
		if toc[0].Anchor != "title-1" {
			t.Errorf("toc[0].Anchor = %q, want title-1", toc[0].Anchor)
		}
	})

	t.Run("nested items", func(t *testing.T) {
		if len(toc[0].Children) != 1 {
			t.Errorf("len(toc[0].Children) = %d, want 1", len(toc[0].Children))
		}
		if toc[0].Children[0].Text != "Subtitle 1.1" {
			t.Errorf("toc[0].Children[0].Text = %q, want Subtitle 1.1", toc[0].Children[0].Text)
		}
	})
}

func TestGenerateAnchor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Hello-World", "hello-world"},
		{"Hello_World", "helloworld"},
		{"Hello  World", "hello--world"},
		{"Hello123", "hello123"},
		{"Hello!@#", "hello"},
		{"中文标题", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generateAnchor(tt.input)
			if result != tt.expected {
				t.Errorf("generateAnchor(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractPlainText(t *testing.T) {
	content := "# Hello World\n\n" +
		"This is **bold** and *italic* text.\n\n" +
		"```go\n" +
		"func main() {}\n" +
		"```\n\n" +
		"Check [this link](http://example.com)."

	result := extractPlainText(content, 100)

	// 应该移除标题标记和格式标记
	if strings.Contains(result, "#") {
		t.Error("Should remove heading markers")
	}
	if strings.Contains(result, "**") {
		t.Error("Should remove bold markers")
	}
	if strings.Contains(result, "func main") {
		t.Error("Should remove code blocks")
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{"hello world", 2},
		{"hello", 1},
		{"Hello World Test", 3},
		{"你好世界", 4}, // 4 个中文字符
		{"hello 你好", 3}, // 1 个英文单词 + 2 个中文字符
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result := countWords(tt.content)
			if result != tt.expected {
				t.Errorf("countWords(%q) = %d, want %d", tt.content, result, tt.expected)
			}
		})
	}
}

func TestRenderToHTML(t *testing.T) {
	t.Run("heading", func(t *testing.T) {
		content := "# Hello World"
		html := renderToHTML(content)
		if !strings.Contains(html, "<h1 id=\"hello-world\">Hello World</h1>") {
			t.Errorf("HTML = %q, should contain h1", html)
		}
	})

	t.Run("bold text", func(t *testing.T) {
		content := "**bold text**"
		html := renderToHTML(content)
		if !strings.Contains(html, "<strong>bold text</strong>") {
			t.Errorf("HTML = %q, should contain strong", html)
		}
	})

	t.Run("italic text", func(t *testing.T) {
		content := "*italic text*"
		html := renderToHTML(content)
		if !strings.Contains(html, "<em>italic text</em>") {
			t.Errorf("HTML = %q, should contain em", html)
		}
	})

	t.Run("code block", func(t *testing.T) {
		content := "```\ncode\n```"
		html := renderToHTML(content)
		if !strings.Contains(html, "<pre><code>") {
			t.Errorf("HTML should contain pre and code")
		}
	})

	t.Run("inline code", func(t *testing.T) {
		content := "`code`"
		html := renderToHTML(content)
		if !strings.Contains(html, "<code>code</code>") {
			t.Errorf("HTML = %q, should contain inline code", html)
		}
	})

	t.Run("unordered list", func(t *testing.T) {
		content := "- item 1\n- item 2"
		html := renderToHTML(content)
		if !strings.Contains(html, "<ul>") || !strings.Contains(html, "</ul>") {
			t.Errorf("HTML = %q, should contain ul", html)
		}
		if !strings.Contains(html, "<li>item 1</li>") {
			t.Errorf("HTML = %q, should contain li", html)
		}
	})

	t.Run("ordered list", func(t *testing.T) {
		content := "1. item 1\n2. item 2"
		html := renderToHTML(content)
		if !strings.Contains(html, "<ol>") || !strings.Contains(html, "</ol>") {
			t.Errorf("HTML = %q, should contain ol", html)
		}
	})

	t.Run("link", func(t *testing.T) {
		content := "[text](http://example.com)"
		html := renderToHTML(content)
		expected := `<a href="http://example.com">text</a>`
		if !strings.Contains(html, expected) {
			t.Errorf("HTML = %q, should contain %q", html, expected)
		}
	})

	t.Run("escape html", func(t *testing.T) {
		content := "<script>alert('xss')</script>"
		html := renderToHTML(content)
		if strings.Contains(html, "<script>") {
			t.Error("HTML should escape script tags")
		}
	})
}

func TestParse(t *testing.T) {
	content := `# Hello World

This is a test post.

## Section 1

Some content here.

- Item 1
- Item 2
`

	result := Parse(content)

	if result.WordCount == 0 {
		t.Error("WordCount should be > 0")
	}

	if result.Excerpt == "" {
		t.Error("Excerpt should not be empty")
	}

	if len(result.TOC) == 0 {
		t.Error("TOC should not be empty")
	}

	if result.HTML == "" {
		t.Error("HTML should not be empty")
	}
}

func TestBuildTOCTree(t *testing.T) {
	items := []*TOCItem{
		{Level: 1, Text: "H1"},
		{Level: 2, Text: "H2"},
		{Level: 2, Text: "H2-2"},
		{Level: 3, Text: "H3"},
		{Level: 1, Text: "H1-2"},
	}

	tree := buildTOCTree(items)

	if len(tree) != 2 {
		t.Errorf("len(tree) = %d, want 2", len(tree))
	}

	if len(tree[0].Children) != 2 {
		t.Errorf("len(tree[0].Children) = %d, want 2", len(tree[0].Children))
	}

	if len(tree[0].Children[0].Children) != 0 {
		t.Errorf("H2 should have no children")
	}

	if len(tree[0].Children[1].Children) != 1 {
		t.Errorf("H2-2 should have 1 child")
	}
}

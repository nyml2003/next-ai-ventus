package domain

import (
	"errors"
	"time"

	"github.com/next-ai-ventus/server/internal/domain/valueobject"
)

var (
	ErrEmptyTitle   = errors.New("post title cannot be empty")
	ErrEmptyContent = errors.New("post content cannot be empty")
	ErrAlreadyPublished = errors.New("post is already published")
	ErrNotPublished     = errors.New("post is not published")
)

// Post 是博客文章实体
type Post struct {
	ID          string
	Title       string
	Slug        valueobject.Slug
	Content     string
	Excerpt     string
	Tags        []valueobject.Tag
	Status      valueobject.PostStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	Version     int
	Cover       string
}

// NewPost 创建新文章
func NewPost(id, title string, slug valueobject.Slug, content string, tags []valueobject.Tag) (*Post, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	now := time.Now()
	post := &Post{
		ID:        id,
		Title:     title,
		Slug:      slug,
		Content:   content,
		Tags:      tags,
		Status:    valueobject.StatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}

	post.GenerateExcerpt(200)
	return post, nil
}

// Publish 发布文章
func (p *Post) Publish() error {
	if p.Status.IsPublished() {
		return ErrAlreadyPublished
	}

	now := time.Now()
	p.Status = valueobject.StatusPublished
	p.PublishedAt = &now
	p.UpdatedAt = now
	p.Version++
	return nil
}

// Unpublish 取消发布（转为草稿）
func (p *Post) Unpublish() error {
	if p.Status.IsDraft() {
		return ErrNotPublished
	}

	p.Status = valueobject.StatusDraft
	p.PublishedAt = nil
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}

// UpdateContent 更新内容
func (p *Post) UpdateContent(content string) error {
	if content == "" {
		return ErrEmptyContent
	}

	p.Content = content
	p.GenerateExcerpt(200)
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}

// UpdateTitle 更新标题
func (p *Post) UpdateTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}

	p.Title = title
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}

// UpdateTags 更新标签
func (p *Post) UpdateTags(tags []valueobject.Tag) {
	p.Tags = tags
	p.UpdatedAt = time.Now()
	p.Version++
}

// GenerateExcerpt 从内容生成摘要
func (p *Post) GenerateExcerpt(maxLen int) {
	if p.Content == "" {
		p.Excerpt = ""
		return
	}

	// 简单的纯文本提取（移除 Markdown 语法）
	excerpt := extractPlainText(p.Content)
	
	if len(excerpt) > maxLen {
		p.Excerpt = excerpt[:maxLen] + "..."
	} else {
		p.Excerpt = excerpt
	}
}

// GetTagNames 获取标签名称列表
func (p *Post) GetTagNames() []string {
	names := make([]string, len(p.Tags))
	for i, tag := range p.Tags {
		names[i] = tag.String()
	}
	return names
}

// HasTag 检查是否有指定标签
func (p *Post) HasTag(tagName string) bool {
	for _, tag := range p.Tags {
		if tag.String() == tagName {
			return true
		}
	}
	return false
}

// IsPublished 检查文章是否已发布
func (p *Post) IsPublished() bool {
	return p.Status.IsPublished()
}

// extractPlainText 简单提取纯文本（移除 Markdown 标记）
func extractPlainText(markdown string) string {
	// 这是一个简化实现，实际应该使用 markdown 解析器
	result := markdown
	
	// 移除标题标记
	for i := 6; i >= 1; i-- {
		prefix := ""
		for j := 0; j < i; j++ {
			prefix += "#"
		}
		result = removeAll(result, prefix+" ")
	}
	
	// 移除粗体、斜体
	result = removeAll(result, "**")
	result = removeAll(result, "*")
	result = removeAll(result, "__")
	result = removeAll(result, "_")
	
	// 移除代码块标记
	result = removeAll(result, "```")
	result = removeAll(result, "`")
	
	// 移除链接标记，保留文本 [](url) -> 
	result = removeLinks(result)
	
	return result
}

func removeAll(s, substr string) string {
	result := s
	for {
		newResult := ""
		found := false
		for i := 0; i < len(result); i++ {
			if i+len(substr) <= len(result) && result[i:i+len(substr)] == substr {
				found = true
				i += len(substr) - 1
				continue
			}
			newResult += string(result[i])
		}
		result = newResult
		if !found {
			break
		}
	}
	return result
}

func removeLinks(s string) string {
	result := ""
	i := 0
	for i < len(s) {
		if s[i] == '[' {
			// 找到对应的 ] 和 (
			closeBracket := -1
			openParen := -1
			closeParen := -1
			
			for j := i + 1; j < len(s); j++ {
				if s[j] == ']' && closeBracket == -1 {
					closeBracket = j
				} else if s[j] == '(' && closeBracket != -1 && openParen == -1 {
					openParen = j
				} else if s[j] == ')' && openParen != -1 {
					closeParen = j
					break
				}
			}
			
			if closeBracket != -1 && openParen != -1 && closeParen != -1 {
				// 提取链接文本
				result += s[i+1 : closeBracket]
				i = closeParen + 1
				continue
			}
		}
		result += string(s[i])
		i++
	}
	return result
}

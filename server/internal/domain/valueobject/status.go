package valueobject

import (
	"errors"
	"fmt"
)

// PostStatus 表示文章发布状态
type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
)

var (
	ValidStatuses = []PostStatus{StatusDraft, StatusPublished}
	ErrInvalidStatus = errors.New("invalid post status")
)

// NewPostStatus 从字符串创建 PostStatus
func NewPostStatus(raw string) (PostStatus, error) {
	switch raw {
	case "draft", "":
		return StatusDraft, nil
	case "published":
		return StatusPublished, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidStatus, raw)
	}
}

// String 返回状态字符串
func (s PostStatus) String() string {
	return string(s)
}

// IsPublished 检查是否已发布
func (s PostStatus) IsPublished() bool {
	return s == StatusPublished
}

// IsDraft 检查是否为草稿
func (s PostStatus) IsDraft() bool {
	return s == StatusDraft
}

// CanPublish 检查是否可以发布（只有草稿可以发布）
func (s PostStatus) CanPublish() bool {
	return s == StatusDraft
}

// CanUnpublish 检查是否可以取消发布（只有已发布可以取消）
func (s PostStatus) CanUnpublish() bool {
	return s == StatusPublished
}

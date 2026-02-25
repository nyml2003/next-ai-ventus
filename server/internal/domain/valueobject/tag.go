package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
)

var (
	// 只允许小写字母、数字、连字符，最大长度 30
	tagRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	maxTagLen = 30

	ErrInvalidTag = errors.New("invalid tag format")
)

// Tag 是文章标签值对象
type Tag struct {
	name string
}

// NewTag 从字符串创建 Tag，验证格式
func NewTag(raw string) (Tag, error) {
	if raw == "" {
		return Tag{}, ErrInvalidTag
	}

	if len(raw) > maxTagLen {
		return Tag{}, fmt.Errorf("%w: tag too long (max %d chars)", ErrInvalidTag, maxTagLen)
	}

	if !tagRegex.MatchString(raw) {
		return Tag{}, fmt.Errorf("%w: %q", ErrInvalidTag, raw)
	}

	return Tag{name: raw}, nil
}

// String 返回 tag 名称
func (t Tag) String() string {
	return t.name
}

// Equals 比较两个 Tag 是否相等
func (t Tag) Equals(other Tag) bool {
	return t.name == other.name
}

// NormalizeTags 规范化标签列表：去重、排序
func NormalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	// 去重
	seen := make(map[string]bool)
	unique := make([]string, 0, len(tags))
	for _, tag := range tags {
		if !seen[tag] {
			seen[tag] = true
			unique = append(unique, tag)
		}
	}

	// 排序
	sort.Strings(unique)

	return unique
}

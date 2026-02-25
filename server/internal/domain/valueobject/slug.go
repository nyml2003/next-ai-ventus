package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mozillazg/go-slugify"
)

var (
	// 只允许小写字母、数字、连字符
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	ErrInvalidSlug = errors.New("invalid slug format")
)

// Slug 是 URL 友好的短链接值对象
type Slug struct {
	value string
}

// NewSlug 从字符串创建 Slug，验证格式
func NewSlug(raw string) (Slug, error) {
	if raw == "" {
		return Slug{}, ErrInvalidSlug
	}

	if !slugRegex.MatchString(raw) {
		return Slug{}, fmt.Errorf("%w: %q", ErrInvalidSlug, raw)
	}

	return Slug{value: raw}, nil
}

// GenerateFromTitle 从标题自动生成唯一 Slug
func GenerateFromTitle(title string, existing []string) Slug {
	// 使用 go-slugify 将标题转换为 slug
	base := slugify.Slugify(title)

	if base == "" {
		// 如果标题全是特殊字符，使用时间戳
		base = "post"
	}

	// 检查冲突
	if !contains(existing, base) {
		return Slug{value: base}
	}

	// 收集已使用的序号
	prefix := base + "-"
	usedNums := make(map[int]bool)
	for _, s := range existing {
		if strings.HasPrefix(s, prefix) {
			numStr := strings.TrimPrefix(s, prefix)
			if num, err := strconv.Atoi(numStr); err == nil {
				usedNums[num] = true
			}
		}
	}

	// 找到第一个未使用的序号
	nextNum := 2
	for usedNums[nextNum] {
		nextNum++
	}

	return Slug{value: fmt.Sprintf("%s-%d", base, nextNum)}
}

// String 返回 slug 字符串
func (s Slug) String() string {
	return s.value
}

// Equals 比较两个 Slug 是否相等
func (s Slug) Equals(other Slug) bool {
	return s.value == other.value
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

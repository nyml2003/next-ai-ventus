package service

import (
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
	"github.com/next-ai-ventus/server/internal/repository"
)

// SlugService Slug 生成服务
type SlugService struct {
	repo repository.PostRepository
}

// NewSlugService 创建 Slug 服务
func NewSlugService(repo repository.PostRepository) *SlugService {
	return &SlugService{repo: repo}
}

// GenerateUniqueSlug 根据标题生成唯一 Slug
func (s *SlugService) GenerateUniqueSlug(title string) (valueobject.Slug, error) {
	// 获取所有文章来收集 slug（这里可以优化，只获取 slug 列）
	result, err := s.repo.FindAll(repository.ListOptions{
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		// 如果获取失败，尝试无冲突生成
		return s.generateSlugWithoutCheck(title), nil
	}

	// 收集所有 slug
	existingSlugs := make([]string, 0, len(result.Items))
	for _, post := range result.Items {
		existingSlugs = append(existingSlugs, post.Slug.String())
	}

	// 生成唯一 slug
	return valueobject.GenerateFromTitle(title, existingSlugs), nil
}

// CheckConflict 检查 slug 是否冲突
func (s *SlugService) CheckConflict(slug string, excludeID string) (bool, error) {
	exists, err := s.repo.Exists(slug)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// 如果提供了 excludeID，检查是否是同一篇文章
	if excludeID != "" {
		post, err := s.repo.FindBySlug(slug)
		if err != nil {
			// 如果找不到，说明 slug 存在但文章不存在（数据不一致）
			return true, nil
		}
		if post.ID == excludeID {
			return false, nil // 同一篇文章，不算冲突
		}
	}

	return true, nil
}

// ValidateSlug 验证 slug 格式
func (s *SlugService) ValidateSlug(slug string) error {
	_, err := valueobject.NewSlug(slug)
	return err
}

// generateSlugWithoutCheck 无检查生成 slug（用于降级）
func (s *SlugService) generateSlugWithoutCheck(title string) valueobject.Slug {
	return valueobject.GenerateFromTitle(title, []string{})
}

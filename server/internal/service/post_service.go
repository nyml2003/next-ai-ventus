package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
	"github.com/next-ai-ventus/server/internal/repository"
)

var (
	ErrVersionConflict = errors.New("version conflict: post has been modified")
	ErrUnauthorized    = errors.New("unauthorized")
)

// CreatePostInput 创建文章输入
type CreatePostInput struct {
	Title   string
	Content string
	Tags    []string
}

// UpdatePostInput 更新文章输入
type UpdatePostInput struct {
	Title   *string
	Content *string
	Tags    []string
	Status  *string
}

// PostService 文章应用服务
type PostService struct {
	repo        repository.PostRepository
	slugService *SlugService
}

// NewPostService 创建文章服务
func NewPostService(repo repository.PostRepository, slugService *SlugService) *PostService {
	return &PostService{
		repo:        repo,
		slugService: slugService,
	}
}

// CreatePost 创建文章
func (s *PostService) CreatePost(input CreatePostInput) (*domain.Post, error) {
	// 验证输入
	if input.Title == "" {
		return nil, domain.ErrEmptyTitle
	}
	if input.Content == "" {
		return nil, domain.ErrEmptyContent
	}

	// 生成唯一 slug
	slug, err := s.slugService.GenerateUniqueSlug(input.Title)
	if err != nil {
		return nil, fmt.Errorf("generate slug failed: %w", err)
	}

	// 生成文章 ID（格式：YYYY-MM-slug）
	now := time.Now()
	id := fmt.Sprintf("%d-%02d-%s", now.Year(), now.Month(), slug.String())

	// 转换标签
	tags, err := s.parseTags(input.Tags)
	if err != nil {
		return nil, fmt.Errorf("invalid tags: %w", err)
	}

	// 创建文章
	post, err := domain.NewPost(id, input.Title, slug, input.Content, tags)
	if err != nil {
		return nil, err
	}

	// 保存
	if err := s.repo.Save(post); err != nil {
		return nil, fmt.Errorf("save post failed: %w", err)
	}

	return post, nil
}

// UpdatePost 更新文章（带乐观锁）
func (s *PostService) UpdatePost(id string, input UpdatePostInput, expectedVersion int) (*domain.Post, error) {
	// 获取现有文章
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 乐观锁检查
	if post.Version != expectedVersion {
		return nil, ErrVersionConflict
	}

	// 更新标题
	if input.Title != nil {
		if err := post.UpdateTitle(*input.Title); err != nil {
			return nil, err
		}
		// 如果标题变了，重新生成 slug
		if *input.Title != post.Title {
			newSlug, err := s.slugService.GenerateUniqueSlug(*input.Title)
			if err != nil {
				return nil, fmt.Errorf("generate slug failed: %w", err)
			}
			post.Slug = newSlug
		}
	}

	// 更新内容
	if input.Content != nil {
		if err := post.UpdateContent(*input.Content); err != nil {
			return nil, err
		}
	}

	// 更新标签
	if input.Tags != nil {
		tags, err := s.parseTags(input.Tags)
		if err != nil {
			return nil, fmt.Errorf("invalid tags: %w", err)
		}
		post.UpdateTags(tags)
	}

	// 更新状态
	if input.Status != nil {
		switch *input.Status {
		case "published":
			if err := post.Publish(); err != nil {
				return nil, err
			}
		case "draft":
			if err := post.Unpublish(); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid status: %s", *input.Status)
		}
	}

	// 保存
	if err := s.repo.Save(post); err != nil {
		return nil, fmt.Errorf("save post failed: %w", err)
	}

	return post, nil
}

// DeletePost 删除文章
func (s *PostService) DeletePost(id string) error {
	// 检查文章是否存在
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// GetPost 获取文章
func (s *PostService) GetPost(id string) (*domain.Post, error) {
	return s.repo.FindByID(id)
}

// GetPostBySlug 根据 slug 获取文章
func (s *PostService) GetPostBySlug(slug string) (*domain.Post, error) {
	return s.repo.FindBySlug(slug)
}

// ListPosts 列出文章
func (s *PostService) ListPosts(opts repository.ListOptions) (*repository.PaginatedResult, error) {
	return s.repo.FindAll(opts)
}

// GetStats 获取文章统计
func (s *PostService) GetStats() (total, published, draft int, err error) {
	total, err = s.repo.Count(repository.CountOptions{})
	if err != nil {
		return 0, 0, 0, err
	}

	published, err = s.repo.Count(repository.CountOptions{Status: "published"})
	if err != nil {
		return 0, 0, 0, err
	}

	draft, err = s.repo.Count(repository.CountOptions{Status: "draft"})
	if err != nil {
		return 0, 0, 0, err
	}

	return total, published, draft, nil
}

// GetAllTags 获取所有标签
func (s *PostService) GetAllTags() ([]string, error) {
	return s.repo.FindAllTags()
}

// parseTags 解析标签字符串
func (s *PostService) parseTags(tagNames []string) ([]valueobject.Tag, error) {
	tags := make([]valueobject.Tag, 0, len(tagNames))
	for _, name := range tagNames {
		if name == "" {
			continue
		}
		tag, err := valueobject.NewTag(name)
		if err != nil {
			return nil, fmt.Errorf("invalid tag %q: %w", name, err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

package repository

import (
	"errors"
	"sync"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrSlugExists   = errors.New("slug already exists")
)

// ListOptions 文章列表查询选项
type ListOptions struct {
	Page     int
	PageSize int
	Tag      string
	Status   string // "", "draft", "published"
	OrderBy  string // "date_desc", "date_asc"
}

// CountOptions 文章计数选项
type CountOptions struct {
	Status string // "", "draft", "published"
}

// PaginatedResult 分页结果
type PaginatedResult struct {
	Items      []*domain.Post
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

// PostRepository 文章仓库接口
type PostRepository interface {
	// FindByID 根据 ID 查找文章
	FindByID(id string) (*domain.Post, error)

	// FindBySlug 根据 Slug 查找文章
	FindBySlug(slug string) (*domain.Post, error)

	// FindAll 查询文章列表（支持分页、标签、状态筛选）
	FindAll(opts ListOptions) (*PaginatedResult, error)

	// FindByTag 根据标签查找文章（不分页，用于索引）
	FindByTag(tag string) ([]*domain.Post, error)

	// FindAllTags 获取所有标签列表
	FindAllTags() ([]string, error)

	// Save 保存文章（创建或更新）
	Save(post *domain.Post) error

	// Delete 删除文章
	Delete(id string) error

	// Exists 检查 Slug 是否已存在
	Exists(slug string) (bool, error)

	// Count 统计文章数量
	Count(opts CountOptions) (int, error)
}

// MemoryPostRepository 内存实现的 PostRepository（用于测试）
type MemoryPostRepository struct {
	posts     map[string]*domain.Post  // id -> post
	slugIndex map[string]string        // slug -> id
	tagIndex  map[string]map[string]struct{}  // tag -> set(id)
	version   int                        // 用于乐观锁检查
	mu        sync.RWMutex               // 并发安全
}

// NewMemoryPostRepository 创建内存仓库实例
func NewMemoryPostRepository() *MemoryPostRepository {
	return &MemoryPostRepository{
		posts:     make(map[string]*domain.Post),
		slugIndex: make(map[string]string),
		tagIndex:  make(map[string]map[string]struct{}),
		version:   1,
	}
}

// FindByID 根据 ID 查找文章
func (r *MemoryPostRepository) FindByID(id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}
	return post, nil
}

// FindBySlug 根据 Slug 查找文章
func (r *MemoryPostRepository) FindBySlug(slug string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.slugIndex[slug]
	if !ok {
		return nil, ErrPostNotFound
	}
	return r.posts[id], nil
}

// FindAll 查询文章列表
func (r *MemoryPostRepository) FindAll(opts ListOptions) (*PaginatedResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 设置默认值
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}
	if opts.OrderBy == "" {
		opts.OrderBy = "date_desc"
	}

	// 筛选文章
	var filtered []*domain.Post
	for _, post := range r.posts {
		// 状态筛选
		if opts.Status != "" && post.Status.String() != opts.Status {
			continue
		}
		// 标签筛选
		if opts.Tag != "" && !post.HasTag(opts.Tag) {
			continue
		}
		filtered = append(filtered, post)
	}

	// 排序
	if opts.OrderBy == "date_desc" {
		// 按创建时间倒序
		for i := 0; i < len(filtered)-1; i++ {
			for j := i + 1; j < len(filtered); j++ {
				if filtered[i].CreatedAt.Before(filtered[j].CreatedAt) {
					filtered[i], filtered[j] = filtered[j], filtered[i]
				}
			}
		}
	} else if opts.OrderBy == "date_asc" {
		// 按创建时间正序
		for i := 0; i < len(filtered)-1; i++ {
			for j := i + 1; j < len(filtered); j++ {
				if filtered[i].CreatedAt.After(filtered[j].CreatedAt) {
					filtered[i], filtered[j] = filtered[j], filtered[i]
				}
			}
		}
	}

	// 分页
	total := len(filtered)
	totalPages := (total + opts.PageSize - 1) / opts.PageSize
	if totalPages < 1 {
		totalPages = 1
	}

	start := (opts.Page - 1) * opts.PageSize
	end := start + opts.PageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	var items []*domain.Post
	if start < total {
		items = filtered[start:end]
	}

	return &PaginatedResult{
		Items:      items,
		Total:      total,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	}, nil
}

// FindByTag 根据标签查找文章
func (r *MemoryPostRepository) FindByTag(tag string) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.tagIndex[tag]
	if !ok || len(ids) == 0 {
		return []*domain.Post{}, nil
	}

	var posts []*domain.Post
	for id := range ids {
		if post, ok := r.posts[id]; ok {
			posts = append(posts, post)
		}
	}

	// 按时间倒序排序
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt.Before(posts[j].CreatedAt) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}

	return posts, nil
}

// FindAllTags 获取所有标签列表
func (r *MemoryPostRepository) FindAllTags() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tagSet := make(map[string]bool)
	for _, post := range r.posts {
		for _, tag := range post.Tags {
			tagSet[tag.String()] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	// 排序
	for i := 0; i < len(tags)-1; i++ {
		for j := i + 1; j < len(tags); j++ {
			if tags[i] > tags[j] {
				tags[i], tags[j] = tags[j], tags[i]
			}
		}
	}

	return tags, nil
}

// Save 保存文章
func (r *MemoryPostRepository) Save(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查 slug 冲突（排除自身）
	if existingID, exists := r.slugIndex[post.Slug.String()]; exists && existingID != post.ID {
		return ErrSlugExists
	}

	// 更新 slug 和标签索引
	if oldPost, ok := r.posts[post.ID]; ok {
		// 保存旧的 slug 和标签（从存储的旧对象读取）
		oldSlugStr := oldPost.Slug.String()
		oldTags := make([]valueobject.Tag, len(oldPost.Tags))
		copy(oldTags, oldPost.Tags)

		// 删除旧 slug（如果不同）
		if oldSlugStr != post.Slug.String() {
			delete(r.slugIndex, oldSlugStr)
		}
		// 删除旧标签索引
		r.removeFromTagIndex(post.ID, oldTags)
	}
	r.slugIndex[post.Slug.String()] = post.ID

	// 保存文章的副本（避免外部修改影响存储）
	r.posts[post.ID] = copyPost(post)

	// 更新标签索引
	r.addToTagIndex(post.ID, post.Tags)

	r.version++
	return nil
}

// Delete 删除文章
func (r *MemoryPostRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.posts[id]
	if !ok {
		return ErrPostNotFound
	}

	// 删除索引
	delete(r.slugIndex, post.Slug.String())
	r.removeFromTagIndex(id, post.Tags)
	delete(r.posts, id)

	r.version++
	return nil
}

// Exists 检查 Slug 是否已存在
func (r *MemoryPostRepository) Exists(slug string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.slugIndex[slug]
	return exists, nil
}

// Count 统计文章数量
func (r *MemoryPostRepository) Count(opts CountOptions) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if opts.Status == "" {
		return len(r.posts), nil
	}

	count := 0
	for _, post := range r.posts {
		if post.Status.String() == opts.Status {
			count++
		}
	}
	return count, nil
}

// addToTagIndex 添加标签索引
func (r *MemoryPostRepository) addToTagIndex(id string, tags []valueobject.Tag) {
	for _, tag := range tags {
		tagName := tag.String()
		if _, ok := r.tagIndex[tagName]; !ok {
			r.tagIndex[tagName] = make(map[string]struct{})
		}
		r.tagIndex[tagName][id] = struct{}{}
	}
}

// removeFromTagIndex 删除标签索引
func (r *MemoryPostRepository) removeFromTagIndex(id string, tags []valueobject.Tag) {
	for _, tag := range tags {
		tagName := tag.String()
		if ids, ok := r.tagIndex[tagName]; ok {
			delete(ids, id)
			if len(ids) == 0 {
				delete(r.tagIndex, tagName)
			}
		}
	}
}

// copyPost 创建文章的深拷贝
func copyPost(post *domain.Post) *domain.Post {
	// 复制标签
	tags := make([]valueobject.Tag, len(post.Tags))
	copy(tags, post.Tags)

	return &domain.Post{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		Excerpt:     post.Excerpt,
		Tags:        tags,
		Status:      post.Status,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
		Version:     post.Version,
		Cover:       post.Cover,
	}
}

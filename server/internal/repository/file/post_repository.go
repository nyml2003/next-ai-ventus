package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
	"github.com/next-ai-ventus/server/internal/repository"
)

// FilePostRepository 文件系统实现的 PostRepository
type FilePostRepository struct {
	basePath string
	posts    map[string]*domain.Post
	slugMap  map[string]string
	tagMap   map[string]map[string]struct{}
	mu       sync.RWMutex
}

// NewFilePostRepository 创建文件存储仓库
func NewFilePostRepository(basePath string) (*FilePostRepository, error) {
	repo := &FilePostRepository{
		basePath: basePath,
		posts:    make(map[string]*domain.Post),
		slugMap:  make(map[string]string),
		tagMap:   make(map[string]map[string]struct{}),
	}

	// 确保目录存在
	postsDir := filepath.Join(basePath, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		return nil, fmt.Errorf("create posts directory failed: %w", err)
	}

	// 加载已有数据
	if err := repo.LoadIndex(); err != nil {
		return nil, fmt.Errorf("load index failed: %w", err)
	}

	return repo, nil
}

// metaJSON 是 meta.json 的结构
type metaJSON struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Excerpt     string   `json:"excerpt"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	PublishedAt *string  `json:"publishedAt,omitempty"`
	Version     int      `json:"version"`
	Cover       string   `json:"cover,omitempty"`
}

// LoadIndex 从文件系统加载索引
func (r *FilePostRepository) LoadIndex() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	postsDir := filepath.Join(r.basePath, "posts")
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		postID := entry.Name()
		post, err := r.loadPost(postID)
		if err != nil {
			// 跳过损坏的文章，记录错误
			continue
		}

		r.posts[postID] = post
		r.slugMap[post.Slug.String()] = postID

		// 更新标签索引
		for _, tag := range post.Tags {
			tagName := tag.String()
			if _, ok := r.tagMap[tagName]; !ok {
				r.tagMap[tagName] = make(map[string]struct{})
			}
			r.tagMap[tagName][postID] = struct{}{}
		}
	}

	return nil
}

// loadPost 加载单篇文章
func (r *FilePostRepository) loadPost(id string) (*domain.Post, error) {
	postDir := filepath.Join(r.basePath, "posts", id)

	// 读取 meta.json
	metaPath := filepath.Join(postDir, "meta.json")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read meta.json failed: %w", err)
	}

	var meta metaJSON
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return nil, fmt.Errorf("parse meta.json failed: %w", err)
	}

	// 读取 content.md
	contentPath := filepath.Join(postDir, "content.md")
	content, err := os.ReadFile(contentPath)
	if err != nil {
		return nil, fmt.Errorf("read content.md failed: %w", err)
	}

	// 重建 Post
	slug, err := valueobject.NewSlug(meta.Slug)
	if err != nil {
		return nil, fmt.Errorf("invalid slug: %w", err)
	}

	tags := make([]valueobject.Tag, 0, len(meta.Tags))
	for _, tagName := range meta.Tags {
		tag, err := valueobject.NewTag(tagName)
		if err != nil {
			continue // 跳过无效标签
		}
		tags = append(tags, tag)
	}

	status, err := valueobject.NewPostStatus(meta.Status)
	if err != nil {
		status = valueobject.StatusDraft
	}

	createdAt, _ := time.Parse(time.RFC3339, meta.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, meta.UpdatedAt)

	var publishedAt *time.Time
	if meta.PublishedAt != nil {
		pt, _ := time.Parse(time.RFC3339, *meta.PublishedAt)
		publishedAt = &pt
	}

	post := &domain.Post{
		ID:          meta.ID,
		Title:       meta.Title,
		Slug:        slug,
		Content:     string(content),
		Excerpt:     meta.Excerpt,
		Tags:        tags,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		PublishedAt: publishedAt,
		Version:     meta.Version,
		Cover:       meta.Cover,
	}

	return post, nil
}

// savePost 保存单篇文章到文件
func (r *FilePostRepository) savePost(post *domain.Post) error {
	postDir := filepath.Join(r.basePath, "posts", post.ID)

	// 创建目录
	if err := os.MkdirAll(postDir, 0755); err != nil {
		return fmt.Errorf("create post directory failed: %w", err)
	}

	// 准备 meta.json
	tagNames := make([]string, len(post.Tags))
	for i, tag := range post.Tags {
		tagNames[i] = tag.String()
	}

	meta := metaJSON{
		ID:        post.ID,
		Title:     post.Title,
		Slug:      post.Slug.String(),
		Excerpt:   post.Excerpt,
		Tags:      tagNames,
		Status:    post.Status.String(),
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		Version:   post.Version,
		Cover:     post.Cover,
	}

	if post.PublishedAt != nil {
		publishedAtStr := post.PublishedAt.Format(time.RFC3339)
		meta.PublishedAt = &publishedAtStr
	}

	// 写入 meta.json
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal meta.json failed: %w", err)
	}

	metaPath := filepath.Join(postDir, "meta.json")
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return fmt.Errorf("write meta.json failed: %w", err)
	}

	// 写入 content.md
	contentPath := filepath.Join(postDir, "content.md")
	if err := os.WriteFile(contentPath, []byte(post.Content), 0644); err != nil {
		return fmt.Errorf("write content.md failed: %w", err)
	}

	return nil
}

// FindByID 根据 ID 查找文章
func (r *FilePostRepository) FindByID(id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, repository.ErrPostNotFound
	}
	return copyPost(post), nil
}

// FindBySlug 根据 Slug 查找文章
func (r *FilePostRepository) FindBySlug(slug string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.slugMap[slug]
	if !ok {
		return nil, repository.ErrPostNotFound
	}

	post, ok := r.posts[id]
	if !ok {
		return nil, repository.ErrPostNotFound
	}

	return copyPost(post), nil
}

// FindAll 查询文章列表
func (r *FilePostRepository) FindAll(opts repository.ListOptions) (*repository.PaginatedResult, error) {
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
		filtered = append(filtered, copyPost(post))
	}

	// 排序
	if opts.OrderBy == "date_desc" {
		sortPostsByDate(filtered, true)
	} else if opts.OrderBy == "date_asc" {
		sortPostsByDate(filtered, false)
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

	return &repository.PaginatedResult{
		Items:      items,
		Total:      total,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	}, nil
}

// FindByTag 根据标签查找文章
func (r *FilePostRepository) FindByTag(tag string) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.tagMap[tag]
	if !ok || len(ids) == 0 {
		return []*domain.Post{}, nil
	}

	var posts []*domain.Post
	for id := range ids {
		if post, ok := r.posts[id]; ok {
			posts = append(posts, copyPost(post))
		}
	}

	sortPostsByDate(posts, true)
	return posts, nil
}

// FindAllTags 获取所有标签列表
func (r *FilePostRepository) FindAllTags() ([]string, error) {
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

	sortStrings(tags)
	return tags, nil
}

// Save 保存文章
func (r *FilePostRepository) Save(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查 slug 冲突
	if existingID, exists := r.slugMap[post.Slug.String()]; exists && existingID != post.ID {
		return repository.ErrSlugExists
	}

	// 如果是更新，删除旧索引
	if oldPost, ok := r.posts[post.ID]; ok {
		delete(r.slugMap, oldPost.Slug.String())
		r.removeFromTagIndex(post.ID, oldPost.Tags)
	}

	// 保存到文件
	if err := r.savePost(post); err != nil {
		return err
	}

	// 更新内存索引
	r.posts[post.ID] = copyPost(post)
	r.slugMap[post.Slug.String()] = post.ID
	r.addToTagIndex(post.ID, post.Tags)

	return nil
}

// Delete 删除文章
func (r *FilePostRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.posts[id]
	if !ok {
		return repository.ErrPostNotFound
	}

	// 删除目录
	postDir := filepath.Join(r.basePath, "posts", id)
	if err := os.RemoveAll(postDir); err != nil {
		return fmt.Errorf("remove post directory failed: %w", err)
	}

	// 删除索引
	delete(r.slugMap, post.Slug.String())
	r.removeFromTagIndex(id, post.Tags)
	delete(r.posts, id)

	return nil
}

// Exists 检查 Slug 是否已存在
func (r *FilePostRepository) Exists(slug string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.slugMap[slug]
	return exists, nil
}

// Count 统计文章数量
func (r *FilePostRepository) Count(opts repository.CountOptions) (int, error) {
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
func (r *FilePostRepository) addToTagIndex(id string, tags []valueobject.Tag) {
	for _, tag := range tags {
		tagName := tag.String()
		if _, ok := r.tagMap[tagName]; !ok {
			r.tagMap[tagName] = make(map[string]struct{})
		}
		r.tagMap[tagName][id] = struct{}{}
	}
}

// removeFromTagIndex 删除标签索引
func (r *FilePostRepository) removeFromTagIndex(id string, tags []valueobject.Tag) {
	for _, tag := range tags {
		tagName := tag.String()
		if ids, ok := r.tagMap[tagName]; ok {
			delete(ids, id)
			if len(ids) == 0 {
				delete(r.tagMap, tagName)
			}
		}
	}
}

// copyPost 创建文章的深拷贝
func copyPost(post *domain.Post) *domain.Post {
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

// sortPostsByDate 按日期排序文章
func sortPostsByDate(posts []*domain.Post, desc bool) {
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			shouldSwap := false
			if desc {
				shouldSwap = posts[i].CreatedAt.Before(posts[j].CreatedAt)
			} else {
				shouldSwap = posts[i].CreatedAt.After(posts[j].CreatedAt)
			}
			if shouldSwap {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
}

// sortStrings 排序字符串切片
func sortStrings(strs []string) {
	for i := 0; i < len(strs)-1; i++ {
		for j := i + 1; j < len(strs); j++ {
			if strs[i] > strs[j] {
				strs[i], strs[j] = strs[j], strs[i]
			}
		}
	}
}

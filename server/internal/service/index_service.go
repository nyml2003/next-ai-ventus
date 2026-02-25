package service

import (
	"sort"
	"time"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/repository"
)

// Index 是文章索引结构
type Index struct {
	PostIDs   []string            // 按时间倒序的文章ID列表
	SlugToID  map[string]string   // slug -> id
	TagToIDs  map[string][]string // tag -> ids（按时间倒序）
	DateToIDs map[string][]string // "2024-06" -> ids
}

// IndexService 索引管理服务
type IndexService struct {
	repo repository.PostRepository
}

// NewIndexService 创建索引服务
func NewIndexService(repo repository.PostRepository) *IndexService {
	return &IndexService{repo: repo}
}

// BuildIndex 从仓库构建完整索引
func (s *IndexService) BuildIndex() (*Index, error) {
	// 获取所有文章
	result, err := s.repo.FindAll(repository.ListOptions{
		Page:     1,
		PageSize: 10000, // 获取所有
	})
	if err != nil {
		return nil, err
	}

	index := &Index{
		PostIDs:   make([]string, 0),
		SlugToID:  make(map[string]string),
		TagToIDs:  make(map[string][]string),
		DateToIDs: make(map[string][]string),
	}

	// 按时间倒序排序（FindAll 默认已排序）
	for _, post := range result.Items {
		s.addToIndex(index, post)
	}

	return index, nil
}

// addToIndex 添加单篇文章到索引
func (s *IndexService) addToIndex(index *Index, post *domain.Post) {
	// 添加到文章ID列表
	index.PostIDs = append(index.PostIDs, post.ID)

	// 添加 slug 映射
	index.SlugToID[post.Slug.String()] = post.ID

	// 添加标签索引
	for _, tag := range post.Tags {
		tagName := tag.String()
		index.TagToIDs[tagName] = append(index.TagToIDs[tagName], post.ID)
	}

	// 添加日期索引
	dateKey := post.CreatedAt.Format("2006-01")
	index.DateToIDs[dateKey] = append(index.DateToIDs[dateKey], post.ID)
}

// SearchByTag 根据标签搜索文章ID
func (s *IndexService) SearchByTag(index *Index, tag string) []string {
	if ids, ok := index.TagToIDs[tag]; ok {
		// 返回副本避免外部修改
		result := make([]string, len(ids))
		copy(result, ids)
		return result
	}
	return []string{}
}

// SearchByDateRange 根据日期范围搜索文章ID
func (s *IndexService) SearchByDateRange(index *Index, start, end time.Time) []string {
	result := make(map[string]bool)

	// 遍历日期索引
	for dateKey, ids := range index.DateToIDs {
		date, _ := time.Parse("2006-01", dateKey)
		if (date.Equal(start) || date.After(start)) && (date.Before(end) || date.Equal(end)) {
			for _, id := range ids {
				result[id] = true
			}
		}
	}

	// 转换为切片
	var ids []string
	for id := range result {
		ids = append(ids, id)
	}

	// 按原始顺序排序
	sort.Slice(ids, func(i, j int) bool {
		idxI := -1
		idxJ := -1
		for k, id := range index.PostIDs {
			if id == ids[i] {
				idxI = k
			}
			if id == ids[j] {
				idxJ = k
			}
		}
		return idxI < idxJ
	})

	return ids
}

// GetSlugID 根据 slug 获取文章ID
func (s *IndexService) GetSlugID(index *Index, slug string) (string, bool) {
	id, ok := index.SlugToID[slug]
	return id, ok
}

// GetAllTags 获取所有标签列表
func (s *IndexService) GetAllTags(index *Index) []string {
	tags := make([]string, 0, len(index.TagToIDs))
	for tag := range index.TagToIDs {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// GetTagCount 获取标签文章数量
func (s *IndexService) GetTagCount(index *Index, tag string) int {
	if ids, ok := index.TagToIDs[tag]; ok {
		return len(ids)
	}
	return 0
}

// GetArchiveMonths 获取归档月份列表
func (s *IndexService) GetArchiveMonths(index *Index) []string {
	months := make([]string, 0, len(index.DateToIDs))
	for month := range index.DateToIDs {
		months = append(months, month)
	}
	// 倒序排序（最新的在前）
	sort.Sort(sort.Reverse(sort.StringSlice(months)))
	return months
}

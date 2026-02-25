package service

import (
	"testing"
	"time"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
	"github.com/next-ai-ventus/server/internal/repository"
)

func createTestPostWithDate(id, title, slugStr string, date time.Time) *domain.Post {
	slug, _ := valueobject.NewSlug(slugStr)
	post, _ := domain.NewPost(id, title, slug, "Content", nil)
	post.CreatedAt = date
	return post
}

func TestIndexService_BuildIndex(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	
	// 创建测试数据
	post1 := createTestPostWithDate("1", "Go Post", "go-post", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC))
	post2 := createTestPostWithDate("2", "Rust Post", "rust-post", time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC))
	post3 := createTestPostWithDate("3", "Web Post", "web-post", time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC))
	
	tagGo, _ := valueobject.NewTag("go")
	tagWeb, _ := valueobject.NewTag("web")
	
	post1.UpdateTags([]valueobject.Tag{tagGo})
	post2.UpdateTags([]valueobject.Tag{tagGo})
	post3.UpdateTags([]valueobject.Tag{tagWeb})
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	service := NewIndexService(repo)
	index, err := service.BuildIndex()
	
	if err != nil {
		t.Errorf("BuildIndex() error = %v", err)
	}
	if len(index.PostIDs) != 3 {
		t.Errorf("len(PostIDs) = %d, want 3", len(index.PostIDs))
	}
	if len(index.SlugToID) != 3 {
		t.Errorf("len(SlugToID) = %d, want 3", len(index.SlugToID))
	}
	if len(index.TagToIDs) != 2 {
		t.Errorf("len(TagToIDs) = %d, want 2", len(index.TagToIDs))
	}
	if len(index.DateToIDs) != 2 { // 2024-06 and 2024-05
		t.Errorf("len(DateToIDs) = %d, want 2", len(index.DateToIDs))
	}
}

func TestIndexService_SearchByTag(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	
	post1 := createTestPostWithDate("1", "Go 1", "go-1", time.Now())
	post2 := createTestPostWithDate("2", "Go 2", "go-2", time.Now())
	post3 := createTestPostWithDate("3", "Rust", "rust", time.Now())
	
	tagGo, _ := valueobject.NewTag("go")
	tagRust, _ := valueobject.NewTag("rust")
	
	post1.UpdateTags([]valueobject.Tag{tagGo})
	post2.UpdateTags([]valueobject.Tag{tagGo})
	post3.UpdateTags([]valueobject.Tag{tagRust})
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	service := NewIndexService(repo)
	index, _ := service.BuildIndex()

	tests := []struct {
		name     string
		tag      string
		expected int
	}{
		{"search go", "go", 2},
		{"search rust", "rust", 1},
		{"search non-existing", "python", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := service.SearchByTag(index, tt.tag)
			if len(ids) != tt.expected {
				t.Errorf("SearchByTag() len = %d, want %d", len(ids), tt.expected)
			}
		})
	}
}

func TestIndexService_SearchByDateRange(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	
	post1 := createTestPostWithDate("1", "June Post", "june-post", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC))
	post2 := createTestPostWithDate("2", "May Post", "may-post", time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC))
	post3 := createTestPostWithDate("3", "April Post", "april-post", time.Date(2024, 4, 20, 0, 0, 0, 0, time.UTC))
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	service := NewIndexService(repo)
	index, _ := service.BuildIndex()

	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "all months",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
		{
			name:     "may to june",
			start:    time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
		{
			name:     "june only",
			start:    time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := service.SearchByDateRange(index, tt.start, tt.end)
			if len(ids) != tt.expected {
				t.Errorf("SearchByDateRange() len = %d, want %d", len(ids), tt.expected)
			}
		})
	}
}

func TestIndexService_GetSlugID(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	post := createTestPostWithDate("1", "Test", "test-slug", time.Now())
	repo.Save(post)

	service := NewIndexService(repo)
	index, _ := service.BuildIndex()

	t.Run("existing slug", func(t *testing.T) {
		id, ok := service.GetSlugID(index, "test-slug")
		if !ok {
			t.Error("GetSlugID() should return ok=true")
		}
		if id != "1" {
			t.Errorf("GetSlugID() = %s, want 1", id)
		}
	})

	t.Run("non-existing slug", func(t *testing.T) {
		_, ok := service.GetSlugID(index, "non-existing")
		if ok {
			t.Error("GetSlugID() should return ok=false")
		}
	})
}

func TestIndexService_GetAllTags(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	
	post1 := createTestPostWithDate("1", "Test", "test-1", time.Now())
	post2 := createTestPostWithDate("2", "Test", "test-2", time.Now())
	
	tagGo, _ := valueobject.NewTag("go")
	tagRust, _ := valueobject.NewTag("rust")
	
	post1.UpdateTags([]valueobject.Tag{tagGo})
	post2.UpdateTags([]valueobject.Tag{tagRust})
	
	repo.Save(post1)
	repo.Save(post2)

	service := NewIndexService(repo)
	index, _ := service.BuildIndex()

	tags := service.GetAllTags(index)
	if len(tags) != 2 {
		t.Errorf("GetAllTags() len = %d, want 2", len(tags))
	}
	// 应该按字母顺序排序
	if tags[0] != "go" || tags[1] != "rust" {
		t.Errorf("GetAllTags() = %v, want [go rust]", tags)
	}
}

func TestIndexService_GetTagCount(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	
	post1 := createTestPostWithDate("1", "Test", "test-1", time.Now())
	post2 := createTestPostWithDate("2", "Test", "test-2", time.Now())
	
	tagGo, _ := valueobject.NewTag("go")
	
	post1.UpdateTags([]valueobject.Tag{tagGo})
	post2.UpdateTags([]valueobject.Tag{tagGo})
	
	repo.Save(post1)
	repo.Save(post2)

	service := NewIndexService(repo)
	index, _ := service.BuildIndex()

	count := service.GetTagCount(index, "go")
	if count != 2 {
		t.Errorf("GetTagCount() = %d, want 2", count)
	}

	count = service.GetTagCount(index, "non-existing")
	if count != 0 {
		t.Errorf("GetTagCount() = %d, want 0", count)
	}
}

func TestIndexService_GetArchiveMonths(t *testing.T) {
	repo := repository.NewMemoryPostRepository()
	
	post1 := createTestPostWithDate("1", "Test", "test-1", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC))
	post2 := createTestPostWithDate("2", "Test", "test-2", time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC))
	post3 := createTestPostWithDate("3", "Test", "test-3", time.Date(2024, 4, 20, 0, 0, 0, 0, time.UTC))
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	service := NewIndexService(repo)
	index, _ := service.BuildIndex()

	months := service.GetArchiveMonths(index)
	if len(months) != 3 {
		t.Errorf("GetArchiveMonths() len = %d, want 3", len(months))
	}
	// 应该按倒序排列（最新的在前）
	if months[0] != "2024-06" {
		t.Errorf("GetArchiveMonths()[0] = %s, want 2024-06", months[0])
	}
}

func BenchmarkBuildIndex(b *testing.B) {
	repo := repository.NewMemoryPostRepository()
	
	// 创建 1000 篇文章
	for i := 0; i < 1000; i++ {
		id := string(rune('a' + i%26)) + string(rune('0'+i/26))
		slug := "post-" + id
		post := createTestPostWithDate(id, "Title "+id, slug, time.Now())
		repo.Save(post)
	}

	service := NewIndexService(repo)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.BuildIndex()
		if err != nil {
			b.Fatal(err)
		}
	}
}

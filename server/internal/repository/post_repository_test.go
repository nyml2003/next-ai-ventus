package repository

import (
	"sync"
	"testing"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
)

func createTestPost(id, title, slugStr string) *domain.Post {
	slug, _ := valueobject.NewSlug(slugStr)
	post, _ := domain.NewPost(id, title, slug, "Content of "+title, nil)
	return post
}

func TestMemoryPostRepository_FindByID(t *testing.T) {
	repo := NewMemoryPostRepository()
	post := createTestPost("1", "Test", "test-slug")
	repo.Save(post)

	t.Run("existing post", func(t *testing.T) {
		found, err := repo.FindByID("1")
		if err != nil {
			t.Errorf("FindByID() error = %v", err)
		}
		if found.ID != "1" {
			t.Errorf("FindByID() = %v, want ID=1", found)
		}
	})

	t.Run("non-existing post", func(t *testing.T) {
		_, err := repo.FindByID("999")
		if err != ErrPostNotFound {
			t.Errorf("FindByID() error = %v, want ErrPostNotFound", err)
		}
	})
}

func TestMemoryPostRepository_FindBySlug(t *testing.T) {
	repo := NewMemoryPostRepository()
	post := createTestPost("1", "Test", "test-slug")
	repo.Save(post)

	t.Run("existing slug", func(t *testing.T) {
		found, err := repo.FindBySlug("test-slug")
		if err != nil {
			t.Errorf("FindBySlug() error = %v", err)
		}
		if found.Slug.String() != "test-slug" {
			t.Errorf("FindBySlug() slug = %v, want test-slug", found.Slug.String())
		}
	})

	t.Run("non-existing slug", func(t *testing.T) {
		_, err := repo.FindBySlug("non-existent")
		if err != ErrPostNotFound {
			t.Errorf("FindBySlug() error = %v, want ErrPostNotFound", err)
		}
	})
}

func TestMemoryPostRepository_Save(t *testing.T) {
	t.Run("create new post", func(t *testing.T) {
		repo := NewMemoryPostRepository()
		post := createTestPost("1", "Test", "test-slug")
		
		err := repo.Save(post)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// 验证索引
		if id, ok := repo.slugIndex["test-slug"]; !ok || id != "1" {
			t.Error("slug index not updated")
		}
	})

	t.Run("update existing post", func(t *testing.T) {
		repo := NewMemoryPostRepository()
		post := createTestPost("1", "Test", "test-slug")
		repo.Save(post)

		// 更新标题和 slug
		newSlug, _ := valueobject.NewSlug("new-slug")
		post.Slug = newSlug
		post.UpdateTitle("New Title")
		
		err := repo.Save(post)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// 旧 slug 应该被删除
		if _, ok := repo.slugIndex["test-slug"]; ok {
			t.Error("old slug index not removed")
		}
		// 新 slug 应该存在
		if id, ok := repo.slugIndex["new-slug"]; !ok || id != "1" {
			t.Error("new slug index not created")
		}
	})

	t.Run("slug conflict", func(t *testing.T) {
		repo := NewMemoryPostRepository()
		post1 := createTestPost("1", "Test 1", "test-slug")
		post2 := createTestPost("2", "Test 2", "test-slug")
		
		repo.Save(post1)
		err := repo.Save(post2)
		
		if err != ErrSlugExists {
			t.Errorf("Save() error = %v, want ErrSlugExists", err)
		}
	})
}

func TestMemoryPostRepository_Delete(t *testing.T) {
	repo := NewMemoryPostRepository()
	post := createTestPost("1", "Test", "test-slug")
	repo.Save(post)

	t.Run("delete existing", func(t *testing.T) {
		err := repo.Delete("1")
		if err != nil {
			t.Errorf("Delete() error = %v", err)
		}

		_, err = repo.FindByID("1")
		if err != ErrPostNotFound {
			t.Error("post should be deleted")
		}

		if _, ok := repo.slugIndex["test-slug"]; ok {
			t.Error("slug index should be removed")
		}
	})

	t.Run("delete non-existing", func(t *testing.T) {
		err := repo.Delete("999")
		if err != ErrPostNotFound {
			t.Errorf("Delete() error = %v, want ErrPostNotFound", err)
		}
	})
}

func TestMemoryPostRepository_FindAll(t *testing.T) {
	repo := NewMemoryPostRepository()
	
	// 创建测试数据
	post1 := createTestPost("1", "First", "first-slug")
	post2 := createTestPost("2", "Second", "second-slug")
	post3 := createTestPost("3", "Third", "third-slug")
	
	// 给 post1 添加标签
	tag1, _ := valueobject.NewTag("go")
	post1.UpdateTags([]valueobject.Tag{tag1})
	
	// 发布 post1
	post1.Publish()
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	t.Run("pagination", func(t *testing.T) {
		result, err := repo.FindAll(ListOptions{Page: 1, PageSize: 2})
		if err != nil {
			t.Errorf("FindAll() error = %v", err)
		}
		if result.Total != 3 {
			t.Errorf("Total = %d, want 3", result.Total)
		}
		if len(result.Items) != 2 {
			t.Errorf("Items length = %d, want 2", len(result.Items))
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		result, err := repo.FindAll(ListOptions{Status: "published"})
		if err != nil {
			t.Errorf("FindAll() error = %v", err)
		}
		if result.Total != 1 {
			t.Errorf("Total = %d, want 1", result.Total)
		}
		if result.Items[0].ID != "1" {
			t.Errorf("Item ID = %s, want 1", result.Items[0].ID)
		}
	})

	t.Run("filter by tag", func(t *testing.T) {
		result, err := repo.FindAll(ListOptions{Tag: "go"})
		if err != nil {
			t.Errorf("FindAll() error = %v", err)
		}
		if result.Total != 1 {
			t.Errorf("Total = %d, want 1", result.Total)
		}
	})
}

func TestMemoryPostRepository_FindByTag(t *testing.T) {
	repo := NewMemoryPostRepository()
	
	post1 := createTestPost("1", "Go Post", "go-post")
	post2 := createTestPost("2", "Rust Post", "rust-post")
	post3 := createTestPost("3", "Another Go", "another-go")
	
	tagGo, _ := valueobject.NewTag("go")
	tagWeb, _ := valueobject.NewTag("web")
	
	post1.UpdateTags([]valueobject.Tag{tagGo, tagWeb})
	post3.UpdateTags([]valueobject.Tag{tagGo})
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	t.Run("find by go", func(t *testing.T) {
		posts, err := repo.FindByTag("go")
		if err != nil {
			t.Errorf("FindByTag() error = %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("len(posts) = %d, want 2", len(posts))
		}
	})

	t.Run("find by web", func(t *testing.T) {
		posts, err := repo.FindByTag("web")
		if err != nil {
			t.Errorf("FindByTag() error = %v", err)
		}
		if len(posts) != 1 {
			t.Errorf("len(posts) = %d, want 1", len(posts))
		}
	})

	t.Run("find by non-existing", func(t *testing.T) {
		posts, err := repo.FindByTag("python")
		if err != nil {
			t.Errorf("FindByTag() error = %v", err)
		}
		if len(posts) != 0 {
			t.Errorf("len(posts) = %d, want 0", len(posts))
		}
	})
}

func TestMemoryPostRepository_FindAllTags(t *testing.T) {
	repo := NewMemoryPostRepository()
	
	post1 := createTestPost("1", "Post 1", "post-1")
	post2 := createTestPost("2", "Post 2", "post-2")
	
	tagGo, _ := valueobject.NewTag("go")
	tagRust, _ := valueobject.NewTag("rust")
	tagWeb, _ := valueobject.NewTag("web")
	
	post1.UpdateTags([]valueobject.Tag{tagGo, tagWeb})
	post2.UpdateTags([]valueobject.Tag{tagRust, tagWeb})
	
	repo.Save(post1)
	repo.Save(post2)

	tags, err := repo.FindAllTags()
	if err != nil {
		t.Errorf("FindAllTags() error = %v", err)
	}
	if len(tags) != 3 {
		t.Errorf("len(tags) = %d, want 3", len(tags))
	}
	// 应该按字母顺序排序
	if tags[0] != "go" || tags[1] != "rust" || tags[2] != "web" {
		t.Errorf("tags = %v, want [go rust web]", tags)
	}
}

func TestMemoryPostRepository_Exists(t *testing.T) {
	repo := NewMemoryPostRepository()
	post := createTestPost("1", "Test", "test-slug")
	repo.Save(post)

	t.Run("existing slug", func(t *testing.T) {
		exists, err := repo.Exists("test-slug")
		if err != nil {
			t.Errorf("Exists() error = %v", err)
		}
		if !exists {
			t.Error("Exists() should be true")
		}
	})

	t.Run("non-existing slug", func(t *testing.T) {
		exists, err := repo.Exists("non-existent")
		if err != nil {
			t.Errorf("Exists() error = %v", err)
		}
		if exists {
			t.Error("Exists() should be false")
		}
	})
}

func TestMemoryPostRepository_Count(t *testing.T) {
	repo := NewMemoryPostRepository()
	
	post1 := createTestPost("1", "Draft", "draft-slug")
	post2 := createTestPost("2", "Published", "published-slug")
	post2.Publish()
	post3 := createTestPost("3", "Another Draft", "another-draft")
	
	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	t.Run("count all", func(t *testing.T) {
		count, err := repo.Count(CountOptions{})
		if err != nil {
			t.Errorf("Count() error = %v", err)
		}
		if count != 3 {
			t.Errorf("Count() = %d, want 3", count)
		}
	})

	t.Run("count published", func(t *testing.T) {
		count, err := repo.Count(CountOptions{Status: "published"})
		if err != nil {
			t.Errorf("Count() error = %v", err)
		}
		if count != 1 {
			t.Errorf("Count() = %d, want 1", count)
		}
	})

	t.Run("count draft", func(t *testing.T) {
		count, err := repo.Count(CountOptions{Status: "draft"})
		if err != nil {
			t.Errorf("Count() error = %v", err)
		}
		if count != 2 {
			t.Errorf("Count() = %d, want 2", count)
		}
	})
}

func TestMemoryPostRepository_TagIndexUpdate(t *testing.T) {
	repo := NewMemoryPostRepository()
	
	post := createTestPost("1", "Test", "test-slug")
	tagGo, _ := valueobject.NewTag("go")
	tagRust, _ := valueobject.NewTag("rust")
	
	post.UpdateTags([]valueobject.Tag{tagGo})
	repo.Save(post)

	// 更新标签
	post.UpdateTags([]valueobject.Tag{tagRust})
	repo.Save(post)

	// 旧标签应该被移除
	postsWithGo, _ := repo.FindByTag("go")
	if len(postsWithGo) != 0 {
		t.Errorf("go tag should have 0 posts, got %d", len(postsWithGo))
	}

	// 新标签应该存在
	postsWithRust, _ := repo.FindByTag("rust")
	if len(postsWithRust) != 1 {
		t.Errorf("rust tag should have 1 post, got %d", len(postsWithRust))
	}
}

func TestMemoryPostRepository_Concurrent(t *testing.T) {
	repo := NewMemoryPostRepository()
	
	// 并发保存多个文章
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := string(rune('0' + n%10))
			slug := "slug-" + string(rune('a'+n%26))
			post := createTestPost(id, "Title", slug)
			repo.Save(post)
		}(i)
	}
	wg.Wait()

	// 验证总数
	result, _ := repo.FindAll(ListOptions{})
	// 因为有 ID 冲突，总数应该 <= 10
	if result.Total > 10 {
		t.Errorf("Total = %d, expected <= 10 (due to ID conflicts)", result.Total)
	}
}

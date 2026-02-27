package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
	"github.com/next-ai-ventus/server/internal/repository"
)

func setupTestRepo(t *testing.T) (*FilePostRepository, string) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "ventus-test-")
	if err != nil {
		t.Fatalf("create temp dir failed: %v", err)
	}

	repo, err := NewFilePostRepository(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("create repository failed: %v", err)
	}

	// 清理函数
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return repo, tmpDir
}

func createTestPost(id, title, slugStr string) *domain.Post {
	slug, _ := valueobject.NewSlug(slugStr)
	post, _ := domain.NewPost(id, title, slug, "Test content", nil)
	return post
}

func TestFilePostRepository_SaveAndFind(t *testing.T) {
	repo, tmpDir := setupTestRepo(t)

	post := createTestPost("2024-06-hello", "Hello World", "hello-world")

	t.Run("save post", func(t *testing.T) {
		err := repo.Save(post)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// 验证文件是否创建
		metaPath := filepath.Join(tmpDir, "posts", "2024-06-hello", "meta.json")
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			t.Error("meta.json should be created")
		}

		contentPath := filepath.Join(tmpDir, "posts", "2024-06-hello", "content.md")
		if _, err := os.Stat(contentPath); os.IsNotExist(err) {
			t.Error("content.md should be created")
		}
	})

	t.Run("find by id", func(t *testing.T) {
		found, err := repo.FindByID("2024-06-hello")
		if err != nil {
			t.Errorf("FindByID() error = %v", err)
			return
		}
		if found.ID != post.ID {
			t.Errorf("ID = %q, want %q", found.ID, post.ID)
		}
	})

	t.Run("find by slug", func(t *testing.T) {
		found, err := repo.FindBySlug("hello-world")
		if err != nil {
			t.Errorf("FindBySlug() error = %v", err)
			return
		}
		if found.ID != post.ID {
			t.Errorf("ID = %q, want %q", found.ID, post.ID)
		}
	})
}

func TestFilePostRepository_LoadIndex(t *testing.T) {
	repo1, tmpDir := setupTestRepo(t)

	// 创建文章
	post := createTestPost("2024-06-test", "Test Post", "test-post")
	tag, _ := valueobject.NewTag("go")
	post.UpdateTags([]valueobject.Tag{tag})
	repo1.Save(post)

	// 创建新的仓库实例（从文件加载）
	repo2, err := NewFilePostRepository(tmpDir)
	if err != nil {
		t.Fatalf("create repository 2 failed: %v", err)
	}

	// 验证是否能找到文章
	found, err := repo2.FindByID("2024-06-test")
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}
	if found.Title != "Test Post" {
		t.Errorf("Title = %q, want Test Post", found.Title)
	}
	if len(found.Tags) != 1 {
		t.Errorf("Tags length = %d, want 1", len(found.Tags))
	}
}

func TestFilePostRepository_Update(t *testing.T) {
	repo, _ := setupTestRepo(t)

	// 创建文章
	post := createTestPost("2024-06-test", "Original", "original-slug")
	repo.Save(post)

	// 更新文章
	newSlug, _ := valueobject.NewSlug("updated-slug")
	post.Slug = newSlug
	post.UpdateTitle("Updated Title")
	repo.Save(post)

	// 验证旧 slug 不存在
	_, err := repo.FindBySlug("original-slug")
	if err != repository.ErrPostNotFound {
		t.Error("Old slug should not be found")
	}

	// 验证新 slug 存在
	found, err := repo.FindBySlug("updated-slug")
	if err != nil {
		t.Errorf("FindBySlug() error = %v", err)
	}
	if found.Title != "Updated Title" {
		t.Errorf("Title = %q, want Updated Title", found.Title)
	}
}

func TestFilePostRepository_Delete(t *testing.T) {
	repo, tmpDir := setupTestRepo(t)

	// 创建文章
	post := createTestPost("2024-06-test", "To Delete", "to-delete")
	repo.Save(post)

	// 删除文章
	err := repo.Delete("2024-06-test")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// 验证文件已删除
	postDir := filepath.Join(tmpDir, "posts", "2024-06-test")
	if _, err := os.Stat(postDir); !os.IsNotExist(err) {
		t.Error("Post directory should be deleted")
	}

	// 验证内存索引已删除
	_, err = repo.FindByID("2024-06-test")
	if err != repository.ErrPostNotFound {
		t.Error("Post should not be found")
	}
}

func TestFilePostRepository_List(t *testing.T) {
	repo, _ := setupTestRepo(t)

	// 创建多篇文章
	for i := 0; i < 5; i++ {
		id := "2024-06-post-" + string(rune('a'+i))
		slug := "post-" + string(rune('a'+i))
		post := createTestPost(id, "Post", slug)
		repo.Save(post)
	}

	result, err := repo.FindAll(repository.ListOptions{
		Page:     1,
		PageSize: 2,
	})
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
		return
	}

	if result.Total != 5 {
		t.Errorf("Total = %d, want 5", result.Total)
	}
	if len(result.Items) != 2 {
		t.Errorf("Items length = %d, want 2", len(result.Items))
	}
}

func TestFilePostRepository_Count(t *testing.T) {
	repo, _ := setupTestRepo(t)

	// 创建文章
	post1 := createTestPost("2024-06-1", "Post 1", "post-1")
	post2 := createTestPost("2024-06-2", "Post 2", "post-2")
	post2.Publish()
	post3 := createTestPost("2024-06-3", "Post 3", "post-3")

	repo.Save(post1)
	repo.Save(post2)
	repo.Save(post3)

	t.Run("count all", func(t *testing.T) {
		count, err := repo.Count(repository.CountOptions{})
		if err != nil {
			t.Errorf("Count() error = %v", err)
		}
		if count != 3 {
			t.Errorf("Count = %d, want 3", count)
		}
	})

	t.Run("count published", func(t *testing.T) {
		count, err := repo.Count(repository.CountOptions{Status: "published"})
		if err != nil {
			t.Errorf("Count() error = %v", err)
		}
		if count != 1 {
			t.Errorf("Count = %d, want 1", count)
		}
	})
}

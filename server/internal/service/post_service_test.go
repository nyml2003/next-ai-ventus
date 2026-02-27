package service

import (
	"testing"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/repository"
)

func setupTestServices() (*PostService, *repository.MemoryPostRepository) {
	repo := repository.NewMemoryPostRepository()
	slugService := NewSlugService(repo)
	postService := NewPostService(repo, slugService)
	return postService, repo
}

func TestPostService_CreatePost(t *testing.T) {
	service, _ := setupTestServices()

	t.Run("create valid post", func(t *testing.T) {
		input := CreatePostInput{
			Title:   "Hello World",
			Content: "# Hello\n\nThis is content.",
			Tags:    []string{"go", "tutorial"},
		}

		post, err := service.CreatePost(input)
		if err != nil {
			t.Errorf("CreatePost() error = %v", err)
			return
		}

		if post.Title != input.Title {
			t.Errorf("Title = %q, want %q", post.Title, input.Title)
		}

		if post.Slug.String() != "hello-world" {
			t.Errorf("Slug = %q, want hello-world", post.Slug.String())
		}

		if len(post.Tags) != 2 {
			t.Errorf("Tags length = %d, want 2", len(post.Tags))
		}

		if post.IsPublished() {
			t.Error("New post should be draft, not published")
		}
	})

	t.Run("create with empty title", func(t *testing.T) {
		input := CreatePostInput{
			Title:   "",
			Content: "Content",
		}

		_, err := service.CreatePost(input)
		if err != domain.ErrEmptyTitle {
			t.Errorf("CreatePost() error = %v, want ErrEmptyTitle", err)
		}
	})

	t.Run("create with empty content", func(t *testing.T) {
		input := CreatePostInput{
			Title:   "Title",
			Content: "",
		}

		_, err := service.CreatePost(input)
		if err != domain.ErrEmptyContent {
			t.Errorf("CreatePost() error = %v, want ErrEmptyContent", err)
		}
	})

	t.Run("create with invalid tags", func(t *testing.T) {
		input := CreatePostInput{
			Title:   "Title",
			Content: "Content",
			Tags:    []string{"InvalidTag"},
		}

		_, err := service.CreatePost(input)
		if err == nil {
			t.Error("CreatePost() should return error for invalid tag")
		}
	})
}

func TestPostService_CreatePostWithDuplicateTitle(t *testing.T) {
	service, _ := setupTestServices()

	// 创建第一篇文章
	input := CreatePostInput{
		Title:   "Hello World",
		Content: "Content 1",
	}
	post1, err := service.CreatePost(input)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	// 创建第二篇相同标题的文章
	input2 := CreatePostInput{
		Title:   "Hello World",
		Content: "Content 2",
	}
	post2, err := service.CreatePost(input2)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	// Slug 应该不同
	if post1.Slug.String() == post2.Slug.String() {
		t.Error("Duplicate titles should have different slugs")
	}

	// 第二篇应该带序号
	if post2.Slug.String() != "hello-world-2" {
		t.Errorf("Second post slug = %q, want hello-world-2", post2.Slug.String())
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	service, _ := setupTestServices()

	// 创建文章
	post, err := service.CreatePost(CreatePostInput{
		Title:   "Original Title",
		Content: "Original content",
	})
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	t.Run("update title", func(t *testing.T) {
		newTitle := "Updated Title"
		updated, err := service.UpdatePost(post.ID, UpdatePostInput{
			Title: &newTitle,
		}, post.Version)

		if err != nil {
			t.Errorf("UpdatePost() error = %v", err)
			return
		}

		if updated.Title != newTitle {
			t.Errorf("Title = %q, want %q", updated.Title, newTitle)
		}

		if updated.Version != post.Version+1 {
			t.Errorf("Version = %d, want %d", updated.Version, post.Version+1)
		}
	})

	t.Run("update with version conflict", func(t *testing.T) {
		newTitle := "Another Title"
		_, err := service.UpdatePost(post.ID, UpdatePostInput{
			Title: &newTitle,
		}, 1) // 使用过期的版本号

		if err != ErrVersionConflict {
			t.Errorf("UpdatePost() error = %v, want ErrVersionConflict", err)
		}
	})

	t.Run("publish post", func(t *testing.T) {
		status := "published"
		updated, err := service.UpdatePost(post.ID, UpdatePostInput{
			Status: &status,
		}, post.Version+1)

		if err != nil {
			t.Errorf("UpdatePost() error = %v", err)
			return
		}

		if !updated.IsPublished() {
			t.Error("Post should be published")
		}
	})
}

func TestPostService_DeletePost(t *testing.T) {
	service, _ := setupTestServices()

	// 创建文章
	post, err := service.CreatePost(CreatePostInput{
		Title:   "To Delete",
		Content: "Content",
	})
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	t.Run("delete existing", func(t *testing.T) {
		err := service.DeletePost(post.ID)
		if err != nil {
			t.Errorf("DeletePost() error = %v", err)
		}

		// 验证已删除
		_, err = service.GetPost(post.ID)
		if err != repository.ErrPostNotFound {
			t.Error("Post should be deleted")
		}
	})

	t.Run("delete non-existing", func(t *testing.T) {
		err := service.DeletePost("non-existing")
		if err != repository.ErrPostNotFound {
			t.Errorf("DeletePost() error = %v, want ErrPostNotFound", err)
		}
	})
}

func TestPostService_GetPost(t *testing.T) {
	service, _ := setupTestServices()

	// 创建文章
	post, err := service.CreatePost(CreatePostInput{
		Title:   "Test",
		Content: "Content",
	})
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	t.Run("get by id", func(t *testing.T) {
		found, err := service.GetPost(post.ID)
		if err != nil {
			t.Errorf("GetPost() error = %v", err)
			return
		}
		if found.ID != post.ID {
			t.Errorf("ID = %q, want %q", found.ID, post.ID)
		}
	})

	t.Run("get by slug", func(t *testing.T) {
		found, err := service.GetPostBySlug(post.Slug.String())
		if err != nil {
			t.Errorf("GetPostBySlug() error = %v", err)
			return
		}
		if found.ID != post.ID {
			t.Errorf("ID = %q, want %q", found.ID, post.ID)
		}
	})
}

func TestPostService_ListPosts(t *testing.T) {
	service, _ := setupTestServices()

	// 创建多篇文章
	for i := 0; i < 5; i++ {
		_, err := service.CreatePost(CreatePostInput{
			Title:   "Post",
			Content: "Content",
		})
		if err != nil {
			t.Fatalf("CreatePost() error = %v", err)
		}
	}

	t.Run("list all", func(t *testing.T) {
		result, err := service.ListPosts(repository.ListOptions{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Errorf("ListPosts() error = %v", err)
			return
		}
		if result.Total != 5 {
			t.Errorf("Total = %d, want 5", result.Total)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		result, err := service.ListPosts(repository.ListOptions{
			Page:     1,
			PageSize: 2,
		})
		if err != nil {
			t.Errorf("ListPosts() error = %v", err)
			return
		}
		if len(result.Items) != 2 {
			t.Errorf("Items length = %d, want 2", len(result.Items))
		}
		if result.TotalPages != 3 {
			t.Errorf("TotalPages = %d, want 3", result.TotalPages)
		}
	})
}

func TestPostService_GetStats(t *testing.T) {
	service, _ := setupTestServices()

	// 创建文章
	post1, _ := service.CreatePost(CreatePostInput{Title: "Post 1", Content: "Content"})
	post2, _ := service.CreatePost(CreatePostInput{Title: "Post 2", Content: "Content"})
	service.CreatePost(CreatePostInput{Title: "Post 3", Content: "Content"})

	// 发布两篇
	status := "published"
	service.UpdatePost(post1.ID, UpdatePostInput{Status: &status}, post1.Version)
	service.UpdatePost(post2.ID, UpdatePostInput{Status: &status}, post2.Version)

	total, published, draft, err := service.GetStats()
	if err != nil {
		t.Errorf("GetStats() error = %v", err)
		return
	}

	if total != 3 {
		t.Errorf("Total = %d, want 3", total)
	}
	if published != 2 {
		t.Errorf("Published = %d, want 2", published)
	}
	if draft != 1 {
		t.Errorf("Draft = %d, want 1", draft)
	}
}

func TestPostService_GetAllTags(t *testing.T) {
	service, _ := setupTestServices()

	// 创建带标签的文章
	service.CreatePost(CreatePostInput{
		Title:   "Post 1",
		Content: "Content",
		Tags:    []string{"go", "tutorial"},
	})
	service.CreatePost(CreatePostInput{
		Title:   "Post 2",
		Content: "Content",
		Tags:    []string{"rust", "tutorial"},
	})

	tags, err := service.GetAllTags()
	if err != nil {
		t.Errorf("GetAllTags() error = %v", err)
		return
	}

	if len(tags) != 3 {
		t.Errorf("len(tags) = %d, want 3", len(tags))
	}
}

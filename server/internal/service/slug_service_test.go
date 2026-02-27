package service

import (
	"testing"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/domain/valueobject"
	"github.com/next-ai-ventus/server/internal/repository"
)

func setupSlugService() (*SlugService, *repository.MemoryPostRepository) {
	repo := repository.NewMemoryPostRepository()
	return NewSlugService(repo), repo
}

func TestSlugService_GenerateUniqueSlug(t *testing.T) {
	service, repo := setupSlugService()

	t.Run("generate from simple title", func(t *testing.T) {
		slug, err := service.GenerateUniqueSlug("Hello World")
		if err != nil {
			t.Errorf("GenerateUniqueSlug() error = %v", err)
			return
		}
		if slug.String() != "hello-world" {
			t.Errorf("Slug = %q, want hello-world", slug.String())
		}
	})

	t.Run("generate unique slug when conflict", func(t *testing.T) {
		// 先创建一篇文章
		slug1, _ := valueobject.NewSlug("test-post")
		post, _ := domain.NewPost("1", "Test Post", slug1, "Content", nil)
		repo.Save(post)

		slug, err := service.GenerateUniqueSlug("Test Post")
		if err != nil {
			t.Errorf("GenerateUniqueSlug() error = %v", err)
			return
		}
		if slug.String() != "test-post-2" {
			t.Errorf("Slug = %q, want test-post-2", slug.String())
		}
	})

	t.Run("generate from chinese title", func(t *testing.T) {
		slug, err := service.GenerateUniqueSlug("你好世界")
		if err != nil {
			t.Errorf("GenerateUniqueSlug() error = %v", err)
			return
		}
		if slug.String() != "ni-hao-shi-jie" {
			t.Errorf("Slug = %q, want ni-hao-shi-jie", slug.String())
		}
	})
}

func TestSlugService_CheckConflict(t *testing.T) {
	service, repo := setupSlugService()

	// 创建文章
	slug, _ := valueobject.NewSlug("test-slug")
	post, _ := domain.NewPost("1", "Test", slug, "Content", nil)
	repo.Save(post)

	t.Run("existing slug", func(t *testing.T) {
		conflict, err := service.CheckConflict("test-slug", "")
		if err != nil {
			t.Errorf("CheckConflict() error = %v", err)
			return
		}
		if !conflict {
			t.Error("CheckConflict() should return true for existing slug")
		}
	})

	t.Run("existing slug with same id", func(t *testing.T) {
		conflict, err := service.CheckConflict("test-slug", "1")
		if err != nil {
			t.Errorf("CheckConflict() error = %v", err)
			return
		}
		if conflict {
			t.Error("CheckConflict() should return false for same id")
		}
	})

	t.Run("non-existing slug", func(t *testing.T) {
		conflict, err := service.CheckConflict("non-existing", "")
		if err != nil {
			t.Errorf("CheckConflict() error = %v", err)
			return
		}
		if conflict {
			t.Error("CheckConflict() should return false for non-existing slug")
		}
	})
}

func TestSlugService_ValidateSlug(t *testing.T) {
	service, _ := setupSlugService()

	tests := []struct {
		slug    string
		wantErr bool
	}{
		{"hello-world", false},
		{"hello-world-123", false},
		{"Hello-World", true},     // 大写
		{"hello_world", true},     // 下划线
		{"hello world", true},     // 空格
		{"", true},                // 空
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			err := service.ValidateSlug(tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSlug(%q) error = %v, wantErr %v", tt.slug, err, tt.wantErr)
			}
		})
	}
}

func BenchmarkGenerateUniqueSlug(b *testing.B) {
	service, repo := setupSlugService()

	// 预创建 100 篇文章
	for i := 0; i < 100; i++ {
		slug, _ := valueobject.NewSlug("test-slug-" + string(rune('a'+i%26)))
		post, _ := domain.NewPost(string(rune('0'+i%10)), "Title", slug, "Content", nil)
		repo.Save(post)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateUniqueSlug("Test Post")
		if err != nil {
			b.Fatal(err)
		}
	}
}



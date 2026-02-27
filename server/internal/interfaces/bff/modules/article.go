package modules

import (
	"errors"

	"github.com/next-ai-ventus/server/pkg/markdown"
)

// ArticleData Article 模块数据
type ArticleData struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Content     string   `json:"content"`
	HTML        string   `json:"html"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	PublishedAt *string  `json:"publishedAt,omitempty"`
	WordCount   int      `json:"wordCount"`
}

// HandleArticle 处理 Article 模块
func HandleArticle(ctx *ModuleContext) (interface{}, error) {
	// 获取 slug 参数
	slug, ok := ctx.Params["slug"].(string)
	if !ok || slug == "" {
		return nil, errors.New("slug is required")
	}

	// 查询文章
	post, err := ctx.Services.PostService.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}

	// 渲染 Markdown
	mdResult := markdown.Parse(post.Content)

	var publishedAt *string
	if post.PublishedAt != nil {
		formatted := post.PublishedAt.Format("2006-01-02")
		publishedAt = &formatted
	}

	return ArticleData{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug.String(),
		Content:     post.Content,
		HTML:        mdResult.HTML,
		Tags:        post.GetTagNames(),
		Status:      post.Status.String(),
		CreatedAt:   post.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   post.UpdatedAt.Format("2006-01-02"),
		PublishedAt: publishedAt,
		WordCount:   mdResult.WordCount,
	}, nil
}

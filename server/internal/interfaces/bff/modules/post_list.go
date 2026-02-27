package modules

import (
	"fmt"

	"github.com/next-ai-ventus/server/internal/repository"
)

// PostListData PostList 模块数据
type PostListData struct {
	Items []PostItem `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

// PostItem 文章列表项
type PostItem struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	Excerpt   string   `json:"excerpt"`
	Tags      []string `json:"tags"`
	Date      string   `json:"date"`
	Href      string   `json:"href"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// HandlePostList 处理 PostList 模块
func HandlePostList(ctx *ModuleContext) (interface{}, error) {
	// 解析参数
	page := 1
	if p, ok := ctx.Params["page"].(float64); ok {
		page = int(p)
	}

	tag := ""
	if t, ok := ctx.Params["tag"].(string); ok {
		tag = t
	}

	// 查询文章列表
	result, err := ctx.Services.PostService.ListPosts(repository.ListOptions{
		Page:       page,
		PageSize:   10,
		Tag:        tag,
		Status:     "published", // 只显示已发布的文章
		OrderBy:    "date_desc",
	})
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]PostItem, 0, len(result.Items))
	for _, post := range result.Items {
		items = append(items, PostItem{
			ID:      post.ID,
			Title:   post.Title,
			Slug:    post.Slug.String(),
			Excerpt: post.Excerpt,
			Tags:    post.GetTagNames(),
			Date:    post.CreatedAt.Format("2006-01-02"),
			Href:    fmt.Sprintf("/pages/post/index.html?slug=%s", post.Slug.String()),
		})
	}

	return PostListData{
		Items: items,
		Pagination: PaginationInfo{
			Page:       result.Page,
			PageSize:   result.PageSize,
			Total:      result.Total,
			TotalPages: result.TotalPages,
		},
	}, nil
}

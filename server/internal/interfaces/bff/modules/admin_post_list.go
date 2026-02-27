package modules

import (
	"fmt"

	"github.com/next-ai-ventus/server/internal/repository"
)

// AdminPostListData AdminPostList 模块数据
type AdminPostListData struct {
	Stats struct {
		Total     int `json:"total"`
		Published int `json:"published"`
		Draft     int `json:"draft"`
	} `json:"stats"`
	Items []AdminPostItem `json:"items"`
	Pagination AdminPaginationInfo `json:"pagination"`
	NewPostHref string `json:"newPostHref"`
}

// AdminPostItem 管理端文章列表项
type AdminPostItem struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	Status    string   `json:"status"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Href      string   `json:"href"`
}

// AdminPaginationInfo 分页信息
type AdminPaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// HandleAdminPostList 处理 AdminPostList 模块
func HandleAdminPostList(ctx *ModuleContext) (interface{}, error) {
	// 解析参数
	page := 1
	if p, ok := ctx.Params["page"].(float64); ok {
		page = int(p)
	}

	status := ""
	if s, ok := ctx.Params["status"].(string); ok {
		status = s
	}

	tag := ""
	if t, ok := ctx.Params["tag"].(string); ok {
		tag = t
	}

	// 获取统计信息
	total, published, draft, err := ctx.Services.PostService.GetStats()
	if err != nil {
		return nil, err
	}

	// 查询文章列表
	result, err := ctx.Services.PostService.ListPosts(repository.ListOptions{
		Page:       page,
		PageSize:   20,
		Tag:        tag,
		Status:     status,
		OrderBy:    "date_desc",
	})
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]AdminPostItem, 0, len(result.Items))
	for _, post := range result.Items {
		items = append(items, AdminPostItem{
			ID:        post.ID,
			Title:     post.Title,
			Slug:      post.Slug.String(),
			Status:    post.Status.String(),
			Tags:      post.GetTagNames(),
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04"),
			Href:      fmt.Sprintf("/pages/admin-editor/index.html?id=%s", post.ID),
		})
	}

	return AdminPostListData{
		Stats: struct {
			Total     int `json:"total"`
			Published int `json:"published"`
			Draft     int `json:"draft"`
		}{
			Total:     total,
			Published: published,
			Draft:     draft,
		},
		Items: items,
		Pagination: AdminPaginationInfo{
			Page:       result.Page,
			PageSize:   result.PageSize,
			Total:      result.Total,
			TotalPages: result.TotalPages,
		},
		NewPostHref: "/pages/admin-editor/index.html",
	}, nil
}

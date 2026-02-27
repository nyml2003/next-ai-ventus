package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/next-ai-ventus/server/internal/domain"
	"github.com/next-ai-ventus/server/internal/interfaces/bff"
	"github.com/next-ai-ventus/server/internal/interfaces/http/response"
	"github.com/next-ai-ventus/server/internal/repository"
	"github.com/next-ai-ventus/server/internal/service"
)

// APIRequest 统一 API 请求
type APIRequest struct {
	SceneCode string                 `json:"sceneCode" binding:"required"`
	Data      map[string]interface{} `json:"data"`
}

// APIHandler 统一 API 处理器
type APIHandler struct {
	postService *service.PostService
	authService *service.AuthService
	bffHandler  *bff.Handler
}

// NewAPIHandler 创建统一 API 处理器
func NewAPIHandler(
	postService *service.PostService,
	authService *service.AuthService,
	bffHandler *bff.Handler,
) *APIHandler {
	return &APIHandler{
		postService: postService,
		authService: authService,
		bffHandler:  bffHandler,
	}
}

// HandlePublic 处理公开 API（无需认证）
func (h *APIHandler) HandlePublic(c *gin.Context) {
	var req APIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	switch req.SceneCode {
	case "auth.login":
		h.handleLogin(c, req.Data)
	case "page.get":
		h.handlePageGet(c, req.Data)
	case "post.recordView":
		h.handleRecordView(c, req.Data)
	default:
		response.Error(c, response.CodeInvalidParam)
	}
}

// HandleAdmin 处理管理 API（需要认证）
func (h *APIHandler) HandleAdmin(c *gin.Context) {
	var req APIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	switch req.SceneCode {
	case "post.create":
		h.handlePostCreate(c, req.Data)
	case "post.update":
		h.handlePostUpdate(c, req.Data)
	case "post.delete":
		h.handlePostDelete(c, req.Data)
	case "post.get":
		h.handlePostGet(c, req.Data)
	case "post.list":
		h.handlePostList(c, req.Data)
	case "file.upload":
		h.handleFileUpload(c)
	default:
		response.Error(c, response.CodeInvalidParam)
	}
}

// ==================== Auth Handlers ====================

func (h *APIHandler) handleLogin(c *gin.Context, data map[string]interface{}) {
	username, _ := data["username"].(string)
	password, _ := data["password"].(string)

	if username == "" || password == "" {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	if !h.authService.ValidateCredentials(username, password) {
		response.Error(c, response.CodeInvalidCredentials)
		return
	}

	token, err := h.authService.GenerateToken(username)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	// 设置 Cookie
	c.SetCookie("token", token, 86400, "/", "", false, true)

	response.Success(c, gin.H{
		"token": token,
	})
}

// ==================== Post Handlers ====================

func (h *APIHandler) handlePostCreate(c *gin.Context, data map[string]interface{}) {
	title, _ := data["title"].(string)
	content, _ := data["content"].(string)

	if title == "" || content == "" {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	// 解析标签
	var tags []string
	if tagList, ok := data["tags"].([]interface{}); ok {
		for _, t := range tagList {
			if tag, ok := t.(string); ok {
				tags = append(tags, tag)
			}
		}
	}

	post, err := h.postService.CreatePost(service.CreatePostInput{
		Title:   title,
		Content: content,
		Tags:    tags,
	})
	if err != nil {
		mapErrorAndRespond(c, err)
		return
	}

	response.Success(c, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"slug":    post.Slug.String(),
		"status":  post.Status.String(),
		"version": post.Version,
	})
}

func (h *APIHandler) handlePostUpdate(c *gin.Context, data map[string]interface{}) {
	id, _ := data["id"].(string)
	if id == "" {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	versionFloat, _ := data["version"].(float64)
	version := int(versionFloat)

	var input service.UpdatePostInput

	if title, ok := data["title"].(string); ok {
		input.Title = &title
	}
	if content, ok := data["content"].(string); ok {
		input.Content = &content
	}
	if status, ok := data["status"].(string); ok {
		input.Status = &status
	}
	if tagList, ok := data["tags"].([]interface{}); ok {
		for _, t := range tagList {
			if tag, ok := t.(string); ok {
				input.Tags = append(input.Tags, tag)
			}
		}
	}

	post, err := h.postService.UpdatePost(id, input, version)
	if err != nil {
		mapErrorAndRespond(c, err)
		return
	}

	response.Success(c, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"slug":    post.Slug.String(),
		"status":  post.Status.String(),
		"version": post.Version,
	})
}

func (h *APIHandler) handlePostDelete(c *gin.Context, data map[string]interface{}) {
	id, _ := data["id"].(string)
	if id == "" {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	if err := h.postService.DeletePost(id); err != nil {
		mapErrorAndRespond(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *APIHandler) handlePostGet(c *gin.Context, data map[string]interface{}) {
	id, _ := data["id"].(string)
	slug, _ := data["slug"].(string)

	var post interface{}
	var err error

	if id != "" {
		post, err = h.postService.GetPost(id)
	} else if slug != "" {
		post, err = h.postService.GetPostBySlug(slug)
	} else {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	if err != nil {
		mapErrorAndRespond(c, err)
		return
	}

	response.Success(c, post)
}

func (h *APIHandler) handlePostList(c *gin.Context, data map[string]interface{}) {
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		OrderBy:  "date_desc",
	}

	if page, ok := data["page"].(float64); ok {
		opts.Page = int(page)
	}
	if pageSize, ok := data["pageSize"].(float64); ok {
		opts.PageSize = int(pageSize)
	}
	if status, ok := data["status"].(string); ok {
		opts.Status = status
	}
	if tag, ok := data["tag"].(string); ok {
		opts.Tag = tag
	}

	result, err := h.postService.ListPosts(opts)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, result)
}

func (h *APIHandler) handleRecordView(c *gin.Context, data map[string]interface{}) {
	// MVP 版本简化处理
	response.Success(c, gin.H{"success": true})
}

// ==================== BFF Handler ====================

func (h *APIHandler) handlePageGet(c *gin.Context, data map[string]interface{}) {
	page, _ := data["page"].(string)
	modules, _ := data["modules"].([]interface{})
	params, _ := data["params"].(map[string]interface{})

	if page == "" || len(modules) == 0 {
		response.Error(c, response.CodeInvalidParam)
		return
	}

	// 转换 modules 为字符串切片
	moduleNames := make([]string, 0, len(modules))
	for _, m := range modules {
		if name, ok := m.(string); ok {
			moduleNames = append(moduleNames, name)
		}
	}

	// 调用 BFF handler 内部方法
	results := h.bffHandler.ExecuteModules(page, moduleNames, params)
	response.Success(c, gin.H{
		"page":    page,
		"modules": results,
	})
}

// ==================== File Handler ====================

func (h *APIHandler) handleFileUpload(c *gin.Context) {
	// 复用原有的上传逻辑
	handler := NewUploadHandler()
	handler.Upload(c)
}

// ==================== Helper Functions ====================

func mapErrorAndRespond(c *gin.Context, err error) {
	switch err {
	case service.ErrVersionConflict:
		response.Error(c, response.CodeVersionConflict)
	case repository.ErrPostNotFound:
		response.Error(c, response.CodePostNotFound)
	case repository.ErrSlugExists:
		response.Error(c, response.CodeSlugExists)
	case domain.ErrEmptyTitle:
		response.Error(c, response.CodeInvalidTitle)
	case domain.ErrEmptyContent:
		response.Error(c, response.CodeInvalidContent)
	case domain.ErrAlreadyPublished:
		response.Error(c, response.CodeInvalidStatus)
	case domain.ErrNotPublished:
		response.Error(c, response.CodeInvalidStatus)
	default:
		response.ErrorWithMessage(c, response.CodeInternalError, err.Error())
	}
}

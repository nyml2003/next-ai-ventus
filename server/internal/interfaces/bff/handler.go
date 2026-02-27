package bff

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/next-ai-ventus/server/internal/interfaces/bff/modules"
	"github.com/next-ai-ventus/server/internal/service"
)

// Handler BFF 处理器
type Handler struct {
	services *modules.Services
	registry map[string]modules.ModuleHandler
}

// NewHandler 创建 BFF 处理器
func NewHandler(postService *service.PostService, indexService *service.IndexService) *Handler {
	services := &modules.Services{
		PostService:  postService,
		IndexService: indexService,
	}

	return &Handler{
		services: services,
		registry: map[string]modules.ModuleHandler{
			// MVP BFF 模块
			"header":         modules.HandleHeader,
			"footer":         modules.HandleFooter,
			"postList":       modules.HandlePostList,
			"article":        modules.HandleArticle,
			"adminSidebar":   modules.HandleAdminSidebar,
			"adminFilter":    modules.HandleAdminFilter,
			"adminPostList":  modules.HandleAdminPostList,
			"editor":         modules.HandleEditor,
			"editorSettings": modules.HandleEditorSettings,
			// P1 扩展模块（预留）
			// "hero":          modules.HandleHero,
			// "sidebar":       modules.HandleSidebar,
			// "toc":           modules.HandleTOC,
			// "related":       modules.HandleRelated,
		},
	}
}

// PageRequest BFF 请求
type PageRequest struct {
	Page    string                 `json:"page" binding:"required"`
	Modules []string               `json:"modules" binding:"required"`
	Params  map[string]interface{} `json:"params"`
}

// ModuleResult 模块结果
type ModuleResult struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// PageResponse BFF 响应
type PageResponse struct {
	Page    string                    `json:"page"`
	Modules map[string]ModuleResult   `json:"modules"`
}

// Handle 处理 BFF 请求
func (h *Handler) Handle(c *gin.Context) {
	var req PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 并行执行模块
	results := h.ExecuteModules(req.Page, req.Modules, req.Params)

	c.JSON(http.StatusOK, PageResponse{
		Page:    req.Page,
		Modules: results,
	})
}

// ExecuteModules 并行执行模块（导出供 APIHandler 使用）
func (h *Handler) ExecuteModules(page string, moduleNames []string, params map[string]interface{}) map[string]ModuleResult {
	results := make(map[string]ModuleResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, name := range moduleNames {
		handler, ok := h.registry[name]
		if !ok {
			results[name] = ModuleResult{
				Code:  404,
				Error: "module not found: " + name,
			}
			continue
		}

		wg.Add(1)
		go func(moduleName string, handler modules.ModuleHandler) {
			defer wg.Done()

			ctx := &modules.ModuleContext{
				Page:     page,
				Params:   params,
				Services: h.services,
			}

			data, err := handler(ctx)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results[moduleName] = ModuleResult{
					Code:  500,
					Error: err.Error(),
				}
			} else {
				results[moduleName] = ModuleResult{
					Code: 200,
					Data: data,
				}
			}
		}(name, handler)
	}

	wg.Wait()
	return results
}

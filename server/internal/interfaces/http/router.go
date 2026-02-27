package http

import (
	"github.com/gin-gonic/gin"

	"github.com/next-ai-ventus/server/internal/interfaces/bff"
	"github.com/next-ai-ventus/server/internal/interfaces/http/handlers"
	"github.com/next-ai-ventus/server/internal/interfaces/http/middleware"
	"github.com/next-ai-ventus/server/internal/interfaces/http/response"
	"github.com/next-ai-ventus/server/internal/service"
)

// SetupRouter 配置路由
func SetupRouter(
	postService *service.PostService,
	authService *service.AuthService,
	bffHandler *bff.Handler,
) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// 创建统一 API 处理器
	apiHandler := handlers.NewAPIHandler(postService, authService, bffHandler)

	// 公开 API - 统一 POST
	r.POST("/api/public", apiHandler.HandlePublic)

	// 静态文件（上传的图片）
	r.Static("/uploads", "./storage/uploads")

	// 需认证 API - 统一 POST
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth())
	{
		admin.POST("", apiHandler.HandleAdmin)
	}

	// 404 处理
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, response.CodeNotFound)
	})

	return r
}

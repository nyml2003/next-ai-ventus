package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/next-ai-ventus/server/internal/interfaces/bff"
	httpInterface "github.com/next-ai-ventus/server/internal/interfaces/http"
	"github.com/next-ai-ventus/server/internal/repository/file"
	"github.com/next-ai-ventus/server/internal/service"
)

func main() {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 获取配置
	contentPath := getEnv("CONTENT_PATH", "./content")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	port := getEnv("PORT", "8080")

	// 初始化仓库
	repo, err := file.NewFilePostRepository(contentPath)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// 初始化服务
	slugService := service.NewSlugService(repo)
	postService := service.NewPostService(repo, slugService)
	indexService := service.NewIndexService(repo)
	authService := service.NewAuthService(jwtSecret)

	// 初始化 BFF 处理器
	bffHandler := bff.NewHandler(postService, indexService)

	// 设置路由
	router := httpInterface.SetupRouter(postService, authService, bffHandler)

	// 启动服务器
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

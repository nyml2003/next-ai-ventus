package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/next-ai-ventus/server/internal/interfaces/http/response"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	basePath string
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		basePath: "./storage/uploads",
	}
}

// Upload 处理文件上传
func (h *UploadHandler) Upload(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, response.CodeFileNotFound)
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := filepath.Ext(header.Filename)
	if !isAllowedFileType(ext) {
		response.Error(c, response.CodeInvalidFileType)
		return
	}

	// 验证文件大小（最大 5MB）
	if header.Size > 5*1024*1024 {
		response.Error(c, response.CodeFileTooLarge)
		return
	}

	// 生成文件名：时间戳_原始文件名
	now := time.Now()
	filename := fmt.Sprintf("%d%02d%02d_%s",
		now.Year(), now.Month(), now.Day(),
		header.Filename)

	// 创建目录：uploads/YYYY/MM/
	dir := filepath.Join(h.basePath,
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%02d", now.Month()))

	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, response.CodeUploadFailed)
		return
	}

	// 保存文件
	dst := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		response.Error(c, response.CodeUploadFailed)
		return
	}

	// 返回访问 URL
	url := fmt.Sprintf("/uploads/%d/%02d/%s", now.Year(), now.Month(), filename)
	response.Success(c, gin.H{
		"url": url,
	})
}

// isAllowedFileType 检查文件类型是否允许
func isAllowedFileType(ext string) bool {
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}
	return allowed[ext]
}

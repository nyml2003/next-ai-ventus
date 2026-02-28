package modules

import (
	"os"
)

// HandleLogo 处理 Logo 模块
func HandleLogo(ctx *ModuleContext) (interface{}, error) {
	// 从环境变量或配置读取站点信息
	siteName := os.Getenv("SITE_NAME")
	if siteName == "" {
		siteName = "Ventus Blog"
	}

	return map[string]interface{}{
		"siteName": siteName,
		"logo":     "/logo.png", // 默认 Logo 路径
		"href":     "/",
	}, nil
}

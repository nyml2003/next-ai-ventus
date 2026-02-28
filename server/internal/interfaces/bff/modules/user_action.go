package modules

// HandleUserAction 处理用户操作模块（登录/用户信息）
func HandleUserAction(ctx *ModuleContext) (interface{}, error) {
	// MVP 版本简化处理
	// 实际应该从 request context 获取 JWT token 并解析用户信息
	
	return map[string]interface{}{
		"isLoggedIn": false,
		"loginHref":  "/pages/login/index.html",
		"user":       nil,
	}, nil
}

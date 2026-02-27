package modules

import "github.com/next-ai-ventus/server/internal/service"

// ModuleContext BFF 模块上下文
type ModuleContext struct {
	Page     string
	Params   map[string]interface{}
	Services *Services
}

// Services 包含所有应用服务
type Services struct {
	PostService  *service.PostService
	IndexService *service.IndexService
}

// ModuleHandler BFF 模块处理函数类型
type ModuleHandler func(ctx *ModuleContext) (interface{}, error)

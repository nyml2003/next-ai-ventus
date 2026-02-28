package modules

// HandleNav 处理导航模块
func HandleNav(ctx *ModuleContext) (interface{}, error) {
	// 返回导航链接列表
	return map[string]interface{}{
		"links": []map[string]string{
			{"name": "首页", "href": "/"},
			{"name": "技术", "href": "/?tag=tech"},
			{"name": "生活", "href": "/?tag=life"},
		},
	}, nil
}

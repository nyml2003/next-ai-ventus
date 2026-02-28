package modules

// HandleTagCloud 处理标签云模块（P1 功能，MVP 返回空）
func HandleTagCloud(ctx *ModuleContext) (interface{}, error) {
	// MVP 版本返回空数据或热门标签
	// 实际应该从 IndexService 获取标签统计
	
	return map[string]interface{}{
		"tags": []map[string]interface{}{
			{"name": "Go", "count": 5, "href": "/?tag=go"},
			{"name": "React", "count": 3, "href": "/?tag=react"},
			{"name": "架构", "count": 2, "href": "/?tag=架构"},
		},
	}, nil
}

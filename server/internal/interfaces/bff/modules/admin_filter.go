package modules

// AdminFilterData AdminFilter 模块数据
type AdminFilterData struct {
	StatusOptions []FilterOption `json:"statusOptions"`
	AllTags       []string       `json:"allTags"`
}

// FilterOption 筛选选项
type FilterOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// HandleAdminFilter 处理 AdminFilter 模块
func HandleAdminFilter(ctx *ModuleContext) (interface{}, error) {
	// 获取所有标签
	tags, err := ctx.Services.PostService.GetAllTags()
	if err != nil {
		tags = []string{}
	}

	return AdminFilterData{
		StatusOptions: []FilterOption{
			{Value: "", Label: "全部"},
			{Value: "published", Label: "已发布"},
			{Value: "draft", Label: "草稿"},
		},
		AllTags: tags,
	}, nil
}

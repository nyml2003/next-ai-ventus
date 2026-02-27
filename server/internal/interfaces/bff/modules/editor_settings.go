package modules

// EditorSettingsData EditorSettings 模块数据
type EditorSettingsData struct {
	AllTags []string `json:"allTags"`
}

// HandleEditorSettings 处理 EditorSettings 模块
func HandleEditorSettings(ctx *ModuleContext) (interface{}, error) {
	// 获取所有标签供选择
	tags, err := ctx.Services.PostService.GetAllTags()
	if err != nil {
		tags = []string{}
	}

	return EditorSettingsData{
		AllTags: tags,
	}, nil
}

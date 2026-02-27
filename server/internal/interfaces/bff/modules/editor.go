package modules



// EditorData Editor 模块数据
type EditorData struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	Cover       string   `json:"cover"`
	Version     int      `json:"version"`
	IsNew       bool     `json:"isNew"`
}

// HandleEditor 处理 Editor 模块（获取文章编辑数据）
func HandleEditor(ctx *ModuleContext) (interface{}, error) {
	// 获取文章 ID（如果是编辑）
	if id, ok := ctx.Params["id"].(string); ok && id != "" {
		post, err := ctx.Services.PostService.GetPost(id)
		if err != nil {
			return nil, err
		}

		return EditorData{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			Tags:    post.GetTagNames(),
			Status:  post.Status.String(),
			Cover:   post.Cover,
			Version: post.Version,
			IsNew:   false,
		}, nil
	}

	// 新建文章返回默认值
	return EditorData{
		Title:   "",
		Content: "",
		Tags:    []string{},
		Status:  "draft",
		Cover:   "",
		IsNew:   true,
	}, nil
}

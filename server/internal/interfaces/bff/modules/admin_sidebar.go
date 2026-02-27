package modules

// AdminSidebarData AdminSidebar 模块数据
type AdminSidebarData struct {
	User struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
	Menu []MenuItem `json:"menu"`
}

// MenuItem 菜单项
type MenuItem struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Href     string `json:"href"`
	Active   bool   `json:"active"`
}

// HandleAdminSidebar 处理 AdminSidebar 模块
func HandleAdminSidebar(ctx *ModuleContext) (interface{}, error) {
	// 确定当前激活的菜单
	currentPage := ctx.Page

	menu := []MenuItem{
		{
			Name:   "仪表盘",
			Icon:   "dashboard",
			Href:   "/admin",
			Active: currentPage == "admin",
		},
		{
			Name:   "文章管理",
			Icon:   "file-text",
			Href:   "/admin/posts",
			Active: currentPage == "adminPosts",
		},
		{
			Name:   "图片管理",
			Icon:   "image",
			Href:   "/admin/images",
			Active: currentPage == "adminImages",
		},
	}

	return AdminSidebarData{
		User: struct {
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		}{
			Name:   "Admin",
			Avatar: "/avatar.png",
		},
		Menu: menu,
	}, nil
}

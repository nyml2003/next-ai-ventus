package modules

// HeaderData Header 模块数据
type HeaderData struct {
	SiteName string `json:"siteName"`
	Logo     string `json:"logo"`
	NavLinks []struct {
		Name string `json:"name"`
		Href string `json:"href"`
	} `json:"navLinks"`
}

// HandleHeader 处理 Header 模块
func HandleHeader(ctx *ModuleContext) (interface{}, error) {
	// 这里可以从配置中读取站点信息
	// MVP 版本使用硬编码
	return HeaderData{
		SiteName: "Ventus Blog",
		Logo:     "/logo.png",
		NavLinks: []struct {
			Name string `json:"name"`
			Href string `json:"href"`
		}{
			{Name: "首页", Href: "/"},
			{Name: "关于", Href: "/about"},
		},
	}, nil
}

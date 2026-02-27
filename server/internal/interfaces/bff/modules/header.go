package modules

// HeaderData Header 模块数据
type HeaderData struct {
	SiteName  string `json:"siteName"`
	Logo      string `json:"logo"`
	NavLinks  []NavLink `json:"navLinks"`
	LoginHref string `json:"loginHref"`
}

// NavLink 导航链接
type NavLink struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

// HandleHeader 处理 Header 模块
func HandleHeader(ctx *ModuleContext) (interface{}, error) {
	return HeaderData{
		SiteName:  "Ventus Blog",
		Logo:      "/logo.png",
		NavLinks: []NavLink{
			{Name: "首页", Href: "/pages/home/index.html"},
			{Name: "关于", Href: "/pages/about/index.html"},
		},
		LoginHref: "/pages/login/index.html",
	}, nil
}

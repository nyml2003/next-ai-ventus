package modules

// FooterData Footer 模块数据
type FooterData struct {
	Copyright string `json:"copyright"`
	PoweredBy string `json:"poweredBy"`
}

// HandleFooter 处理 Footer 模块
func HandleFooter(ctx *ModuleContext) (interface{}, error) {
	return FooterData{
		Copyright: "© 2024 Ventus Blog",
		PoweredBy: "Powered by Ventus",
	}, nil
}

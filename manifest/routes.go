// manifest/routes.go 解析插件前端路由声明。
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

// Route 前端路由与菜单声明,不绑定具体前端框架。
type Route struct {
	// View 宿主区域:console 控制台 / uc 用户中心 / index 首页。
	View string `json:"view"`
	Path string `json:"path"`
	Name string `json:"name"`
	// App 该菜单所需的插件应用编码;与 Action 成对出现。
	App string `json:"app,omitempty"`
	// Action 该菜单所需的操作名(来自 action.json);声明后按权限过滤菜单。
	Action string `json:"action,omitempty"`
	Menu   *Menu  `json:"menu,omitempty"`
}

// 宿主区域枚举。
const (
	ViewConsole = "console"
	ViewUC      = "uc"
	ViewIndex   = "index"
)

// ValidView 报告区域是否为合法枚举值。
func ValidView(view string) bool {
	switch view {
	case ViewConsole, ViewUC, ViewIndex:
		return true
	}
	return false
}

// Menu 侧边栏菜单项声明。
type Menu struct {
	Title string `json:"title"`
	Icon  string `json:"icon,omitempty"`
	Group string `json:"group,omitempty"`
	Order int    `json:"order,omitempty"`
}

// LoadRoutes 读取并校验 routes.json。
func LoadRoutes(path string) ([]Route, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("路由文件读取失败: %w", err)
	}
	var routes []Route
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("路由 JSON 无效: %w", err)
	}
	if routes == nil {
		return nil, fmt.Errorf("路由须为 JSON 数组")
	}
	seenPaths := make(map[string]bool, len(routes))
	for i, route := range routes {
		if route.Path == "" {
			return nil, fmt.Errorf("routes[%d].path 必填", i)
		}
		if seenPaths[route.Path] {
			return nil, fmt.Errorf("routes[%d].path 重复: %s", i, route.Path)
		}
		seenPaths[route.Path] = true
		if route.Name == "" {
			return nil, fmt.Errorf("routes[%d].name 必填", i)
		}
		if !ValidView(route.View) {
			return nil, fmt.Errorf("routes[%d].view 须为 console/uc/index", i)
		}
		if (route.App == "") != (route.Action == "") {
			return nil, fmt.Errorf("routes[%d] 的 app 与 action 须成对声明", i)
		}
		if route.Menu != nil && route.Menu.Title == "" {
			return nil, fmt.Errorf("routes[%d].menu.title 必填", i)
		}
	}
	return routes, nil
}

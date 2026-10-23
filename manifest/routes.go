// manifest/routes.go 解析插件前端路由声明。
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

// Route 前端路由与菜单声明,不绑定具体前端框架。
type Route struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Menu *Menu  `json:"menu,omitempty"`
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
	for i, route := range routes {
		if route.Path == "" {
			return nil, fmt.Errorf("routes[%d].path 必填", i)
		}
		if route.Name == "" {
			return nil, fmt.Errorf("routes[%d].name 必填", i)
		}
	}
	return routes, nil
}

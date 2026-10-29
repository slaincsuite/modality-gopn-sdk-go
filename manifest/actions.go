// manifest/actions.go 解析并校验插件操作目录(action.json):应用、操作与数据源
// 的单一来源,同时定义插件前端的权限单元。
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ActionFileName 操作目录缺省文件名;plugin.json 的 actions 字段可指向其他路径。
const ActionFileName = "action.json"

// msgKebab 命名约束的校验消息,多处复用。
const msgKebab = "须为 kebab-case"

// App 插件对宿主暴露的应用;操作按读/写/列表/其他四类登记,策略与插件前端共用。
type App struct {
	// Code 应用编码,同时作为 sam 功能的服务段。
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Actions 应用暴露的操作清单。
	Actions []Action `json:"actions"`
	// Resources 应用可访问的数据源清单。
	Resources []Resource `json:"resources,omitempty"`
}

// Action 应用的一个操作;既作策略的动作段,也是插件前端权限的最小单元。
type Action struct {
	// Name 操作名,同时作为 sam 功能的动作段。
	Name        string `json:"name"`
	Title       string `json:"title"`
	Kind        string `json:"kind"` // read|write|list|other
	Description string `json:"description,omitempty"`
	// Resource 所属数据源编码;须在 App.Resources 中声明。
	Resource string `json:"resource,omitempty"`
}

// Resource 应用的数据源声明。
type Resource struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// 操作类别枚举:读操作 / 写操作 / 列表操作 / 其他操作。
const (
	OpKindRead  = "read"
	OpKindWrite = "write"
	OpKindList  = "list"
	OpKindOther = "other"
)

// ValidOperationKind 报告操作类别是否为合法枚举值。
func ValidOperationKind(kind string) bool {
	switch kind {
	case OpKindRead, OpKindWrite, OpKindList, OpKindOther:
		return true
	}
	return false
}

// actionFile action.json 的根结构。
type actionFile struct {
	Apps []App `json:"apps"`
}

// LoadActions 读取并校验 action.json;返回应用声明列表。
func LoadActions(path string) ([]App, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: 读取 %s: %w", ErrInvalidActions, path, err)
	}
	var catalog actionFile
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("%w: 解析 %s: %w", ErrInvalidActions, path, err)
	}
	if fields := validateApps(catalog.Apps); len(fields) > 0 {
		return nil, fmt.Errorf("%w: %w", ErrInvalidActions, &ValidationError{Errors: fields})
	}
	return catalog.Apps, nil
}

// loadActions 加载操作目录:声明 actions 或存在 action.json 时解析并校验;
// 均不存在时 Apps 为空(纯服务插件合法)。
func (m *Manifest) loadActions(dir string) error {
	path := resolveActionsPath(m, dir)
	if path == "" {
		return nil
	}
	apps, err := LoadActions(path)
	if err != nil {
		return err
	}
	m.Apps = apps
	return nil
}

// resolveActionsPath 解析操作目录文件路径:plugin.json 的 actions 字段优先,
// 缺省探测同目录 action.json;均不存在时返回空串。
func resolveActionsPath(m *Manifest, dir string) string {
	if m.Actions != "" {
		return filepath.Join(dir, m.Actions)
	}
	def := filepath.Join(dir, ActionFileName)
	if _, err := os.Stat(def); err == nil {
		return def
	}
	return ""
}

// validateActionsFile 校验 plugin.json 的 actions 指针所指文件存在。
func validateActionsFile(dir, actions string) []FieldError {
	if actions == "" || dir == "" {
		return nil
	}
	if _, err := os.Stat(filepath.Join(dir, actions)); err != nil {
		return []FieldError{{Field: "actions", Message: "文件不存在: " + actions}}
	}
	return nil
}

// validateApps 校验应用声明:应用编码/操作名 kebab-case 且唯一、操作类别合法、
// 数据源编码唯一且被操作引用时必须已声明。
func validateApps(apps []App) []FieldError {
	var fields []FieldError
	seenApps := make(map[string]struct{}, len(apps))
	for i, app := range apps {
		field := fmt.Sprintf("apps[%d]", i)
		seenCode := true
		switch {
		case app.Code == "":
			fields = append(fields, FieldError{Field: field + ".code", Message: "必填"})
			seenCode = false
		case !isKebabCase(app.Code):
			fields = append(fields, FieldError{Field: field + ".code", Message: msgKebab})
			seenCode = false
		case hasKey(seenApps, app.Code):
			fields = append(fields, FieldError{Field: field + ".code", Message: "重复: " + app.Code})
		}
		if seenCode {
			seenApps[app.Code] = struct{}{}
		}
		if app.Name == "" {
			fields = append(fields, FieldError{Field: field + ".name", Message: "必填"})
		}
		resources := make(map[string]struct{}, len(app.Resources))
		for j, res := range app.Resources {
			rf := fmt.Sprintf("%s.resources[%d]", field, j)
			seenResource := true
			switch {
			case res.Code == "":
				fields = append(fields, FieldError{Field: rf + ".code", Message: "必填"})
				seenResource = false
			case !isKebabCase(res.Code):
				fields = append(fields, FieldError{Field: rf + ".code", Message: msgKebab})
				seenResource = false
			case hasKey(resources, res.Code):
				fields = append(fields, FieldError{Field: rf + ".code", Message: "重复: " + res.Code})
			}
			if seenResource {
				resources[res.Code] = struct{}{}
			}
			if res.Name == "" {
				fields = append(fields, FieldError{Field: rf + ".name", Message: "必填"})
			}
		}
		if len(app.Actions) == 0 {
			fields = append(fields, FieldError{Field: field + ".actions", Message: "至少声明一个操作"})
		}
		seenActions := make(map[string]struct{}, len(app.Actions))
		for j, action := range app.Actions {
			af := fmt.Sprintf("%s.actions[%d]", field, j)
			seenAction := true
			switch {
			case action.Name == "":
				fields = append(fields, FieldError{Field: af + ".name", Message: "必填"})
				seenAction = false
			case !isKebabCase(action.Name):
				fields = append(fields, FieldError{Field: af + ".name", Message: msgKebab})
				seenAction = false
			case hasKey(seenActions, action.Name):
				fields = append(fields, FieldError{Field: af + ".name", Message: "重复: " + action.Name})
			}
			if seenAction {
				seenActions[action.Name] = struct{}{}
			}
			if action.Title == "" {
				fields = append(fields, FieldError{Field: af + ".title", Message: "必填"})
			}
			if !ValidOperationKind(action.Kind) {
				fields = append(fields, FieldError{Field: af + ".kind", Message: "须为 read/write/list/other"})
			}
			if action.Resource != "" {
				if _, ok := resources[action.Resource]; !ok {
					fields = append(fields, FieldError{Field: af + ".resource", Message: "未在 resources 声明: " + action.Resource})
				}
			}
		}
	}
	return fields
}

// hasKey 报告集合中是否已登记该键。
func hasKey(set map[string]struct{}, key string) bool {
	_, ok := set[key]
	return ok
}

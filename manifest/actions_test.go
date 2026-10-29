// manifest/actions_test.go 操作目录(action.json)单测:解析、字段路径错误与缺省发现。
package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// validActionsJSON 合法操作目录声明(四类操作 + 数据源,对齐示例插件)。
const validActionsJSON = `{
	"apps": [
		{
			"code": "plugin-starter",
			"name": "示例插件",
			"description": "示例应用",
			"actions": [
				{"name": "example-list", "title": "查询示例列表", "kind": "list", "description": "分页查询示例数据", "resource": "example"},
				{"name": "example-query", "title": "查看示例", "kind": "read", "description": "查看示例数据详情", "resource": "example"},
				{"name": "example-create", "title": "新建示例", "kind": "write", "description": "创建示例数据", "resource": "example"},
				{"name": "example-archive", "title": "归档示例", "kind": "other", "description": "归档示例数据", "resource": "example"}
			],
			"resources": [{"code": "example", "name": "示例数据源", "description": "示例记录"}]
		}
	]
}`

func TestLoadActions_Valid(t *testing.T) {
	path := writeActions(t, validActionsJSON)
	apps, err := LoadActions(path)
	if err != nil {
		t.Fatalf("LoadActions() = %v", err)
	}
	if len(apps) != 1 {
		t.Fatalf("apps 数量 = %d, want 1", len(apps))
	}
	app := apps[0]
	if app.Code != "plugin-starter" || app.Name != "示例插件" {
		t.Fatalf("应用字段不符: %+v", app)
	}
	if len(app.Actions) != 4 || app.Actions[3].Kind != OpKindOther {
		t.Fatalf("操作声明不符: %+v", app.Actions)
	}
	if len(app.Resources) != 1 || app.Resources[0].Code != "example" {
		t.Fatalf("数据源声明不符: %+v", app.Resources)
	}
}

func TestLoadActions_FieldErrors(t *testing.T) {
	cases := []struct {
		name      string
		apps      string
		wantField string
	}{
		{
			name:      "kind 非法",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"a-x","title":"A","kind":"copy"}]}]`,
			wantField: "apps[0].actions[0].kind",
		},
		{
			name:      "kind 缺失",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"a-x","title":"A"}]}]`,
			wantField: "apps[0].actions[0].kind",
		},
		{
			name:      "操作重名",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"a-x","title":"A","kind":"read"},{"name":"a-x","title":"B","kind":"write"}]}]`,
			wantField: "apps[0].actions[1].name",
		},
		{
			name:      "资源引用悬空",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"a-x","title":"A","kind":"read","resource":"gone"}]}]`,
			wantField: "apps[0].actions[0].resource",
		},
		{
			name:      "应用编码非 kebab-case",
			apps:      `[{"code":"Demo_App","name":"示例","actions":[{"name":"a-x","title":"A","kind":"read"}]}]`,
			wantField: "apps[0].code",
		},
		{
			name:      "操作名非 kebab-case",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"A_X","title":"A","kind":"read"}]}]`,
			wantField: "apps[0].actions[0].name",
		},
		{
			name:      "应用无操作",
			apps:      `[{"code":"demo","name":"示例","actions":[]}]`,
			wantField: "apps[0].actions",
		},
		{
			name:      "应用缺名",
			apps:      `[{"code":"demo","actions":[{"name":"a-x","title":"A","kind":"read"}]}]`,
			wantField: "apps[0].name",
		},
		{
			name:      "应用编码重复",
			apps:      `[{"code":"demo","name":"A","actions":[{"name":"a-x","title":"A","kind":"read"}]},{"code":"demo","name":"B","actions":[{"name":"b-x","title":"B","kind":"read"}]}]`,
			wantField: "apps[1].code",
		},
		{
			name:      "操作缺展示名",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"a-x","kind":"read"}]}]`,
			wantField: "apps[0].actions[0].title",
		},
		{
			name:      "数据源编码重复",
			apps:      `[{"code":"demo","name":"示例","actions":[{"name":"a-x","title":"A","kind":"read"}],"resources":[{"code":"r","name":"R"},{"code":"r","name":"R2"}]}]`,
			wantField: "apps[0].resources[1].code",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeActions(t, `{"apps":`+tc.apps+`}`)
			_, err := LoadActions(path)
			if !errors.Is(err, ErrInvalidActions) {
				t.Fatalf("LoadActions() = %v, want ErrInvalidActions", err)
			}
			var validation *ValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("LoadActions() = %v, want ValidationError", err)
			}
			if !hasField(validation, tc.wantField) {
				t.Fatalf("错误缺少字段 %s: %+v", tc.wantField, validation.Errors)
			}
		})
	}
}

func TestLoadActions_Unreadable(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{name: "文件不存在", path: filepath.Join(t.TempDir(), ActionFileName)},
		{name: "JSON 非法", path: writeActions(t, `{"apps":`)},
		{name: "apps 非数组", path: writeActions(t, `{"apps":{}}`)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := LoadActions(tc.path); !errors.Is(err, ErrInvalidActions) {
				t.Fatalf("LoadActions() = %v, want ErrInvalidActions", err)
			}
		})
	}
}

func TestLoad_AutoDiscoversActionFile(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin","services":["demo.v1.Plugin"]}`, true)
	if err := os.WriteFile(filepath.Join(dir, ActionFileName), []byte(validActionsJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(got.Apps) != 1 || got.Apps[0].Code != "plugin-starter" {
		t.Fatalf("apps = %+v", got.Apps)
	}
}

func TestLoad_WithoutActionFile(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin","services":["demo.v1.Plugin"]}`, true)
	got, err := Load(dir)
	if err != nil {
		t.Fatalf("缺省无操作目录应合法: %v", err)
	}
	if len(got.Apps) != 0 {
		t.Fatalf("未声明操作目录时应为空,得到 %+v", got.Apps)
	}
}

func TestLoad_DeclaredActionsFile(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin","services":["demo.v1.Plugin"],"actions":"catalog/ops.json"}`, true)
	if err := os.MkdirAll(filepath.Join(dir, "catalog"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "catalog", "ops.json"), []byte(validActionsJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(got.Apps) != 1 {
		t.Fatalf("apps = %+v", got.Apps)
	}
}

func TestLoad_DeclaredActionsFileMissing(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin","services":["demo.v1.Plugin"],"actions":"ops.json"}`, true)
	_, err := Load(dir)
	if !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("Load() = %v, want ErrInvalidManifest", err)
	}
	var validation *ValidationError
	if !errors.As(err, &validation) || !hasField(validation, "actions") {
		t.Fatalf("错误缺少 actions 字段: %v", err)
	}
}

func TestLoad_InvalidActionsRejected(t *testing.T) {
	dir := writePlugin(t, "demo-plugin", `{"name":"demo-plugin","version":"1.0.0","binary":"demo-plugin","services":["demo.v1.Plugin"]}`, true)
	invalid := `{"apps":[{"code":"demo","name":"示例","actions":[{"name":"a-x","title":"A","kind":"copy"}]}]}`
	if err := os.WriteFile(filepath.Join(dir, ActionFileName), []byte(invalid), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(dir)
	if !errors.Is(err, ErrInvalidManifest) || !errors.Is(err, ErrInvalidActions) {
		t.Fatalf("Load() = %v, want ErrInvalidManifest 且 ErrInvalidActions", err)
	}
	var validation *ValidationError
	if !errors.As(err, &validation) || !hasField(validation, "apps[0].actions[0].kind") {
		t.Fatalf("错误缺少字段路径: %v", err)
	}
}

// writeActions 写入临时 action.json 并返回路径。
func writeActions(t *testing.T, data string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ActionFileName)
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// hasField 报告校验错误中是否包含指定字段路径。
func hasField(validation *ValidationError, field string) bool {
	for _, fieldErr := range validation.Errors {
		if fieldErr.Field == field {
			return true
		}
	}
	return false
}

// Package manifest 解析并校验 modality 插件包清单。
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	semver "github.com/Masterminds/semver/v3"
)

// FileName 插件清单固定文件名。
const FileName = "plugin.json"

// Manifest 对应插件包根目录的 plugin.json。
type Manifest struct {
	Name           string   `json:"name"`
	DisplayName    string   `json:"displayName,omitempty"`
	Version        string   `json:"version"`
	Description    string   `json:"description,omitempty"`
	Author         string   `json:"author,omitempty"`
	Binary         string   `json:"binary"`
	MinHostVersion string   `json:"min_host_version,omitempty"`
	Services       []string `json:"services"`
	// Grants 插件申请的宿主能力列表;宿主逐项重新鉴权,SDK 只做格式校验。
	Grants []string `json:"grants,omitempty"`
	// Extensions 插件贡献的扩展实现;服务名必须出现在 Services。
	Extensions []Extension `json:"extensions,omitempty"`
	UI             *UI      `json:"ui,omitempty"`
	Dir            string   `json:"-"`
}

// Extension 插件扩展实现声明。
type Extension struct {
	// Name 扩展实现名。
	Name string `json:"name"`
	// Point 宿主扩展点名。
	Point string `json:"point"`
	// DisplayName 展示名。
	DisplayName string `json:"displayName,omitempty"`
	// Description 说明。
	Description string `json:"description,omitempty"`
	// Priority 多实例优先级,大者先执行。
	Priority int `json:"priority,omitempty"`
	// ServiceName 实现 RPC 服务名。
	ServiceName string `json:"service_name"`
}

// UI 插件前端资源声明。
type UI struct {
	Entry      string `json:"entry"`
	RoutesFile string `json:"routes_file,omitempty"`
	MountPath  string `json:"mount_path,omitempty"`
}

// FieldError 单个字段的校验错误。
type FieldError struct {
	Field   string
	Message string
}

// Error 返回带字段路径的错误文本。
func (e FieldError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

// ValidationError 汇总清单校验错误,可经 errors.As 提取字段明细。
type ValidationError struct {
	Errors []FieldError
}

// Error 汇总全部字段错误。
func (e *ValidationError) Error() string {
	parts := make([]string, len(e.Errors))
	for i, fieldErr := range e.Errors {
		parts[i] = fieldErr.Error()
	}
	return "manifest: " + strings.Join(parts, "; ")
}

// Load 从插件目录读取并校验 plugin.json。
func Load(dir string) (*Manifest, error) {
	absolute, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("manifest: 解析目录 %s: %w", dir, err)
	}
	path := filepath.Join(absolute, FileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: 读取 %s: %w", ErrInvalidManifest, path, err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("%w: 解析 %s: %w", ErrInvalidManifest, path, err)
	}
	manifest.Dir = absolute
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidManifest, err)
	}
	return &manifest, nil
}

// Discover 扫描 dir 的直接子目录并加载其中的 plugin.json。
func Discover(dir string) ([]*Manifest, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("manifest: 读取目录 %s: %w", dir, err)
	}
	var manifests []*Manifest
	seen := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pluginDir := filepath.Join(dir, entry.Name())
		if _, statErr := os.Stat(filepath.Join(pluginDir, FileName)); statErr != nil {
			continue
		}
		loaded, loadErr := Load(pluginDir)
		if loadErr != nil {
			return nil, loadErr
		}
		if previous, exists := seen[loaded.Name]; exists {
			return nil, fmt.Errorf("%w: %s (%s, %s)", ErrDuplicateName, loaded.Name, previous, pluginDir)
		}
		seen[loaded.Name] = pluginDir
		manifests = append(manifests, loaded)
	}
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].Name < manifests[j].Name
	})
	return manifests, nil
}

// CheckHostVersion 校验目标宿主版本是否满足 min_host_version 约束。
func (m *Manifest) CheckHostVersion(hostVersion string) error {
	if m == nil || m.MinHostVersion == "" {
		return nil
	}
	constraint, err := semver.NewConstraint(m.MinHostVersion)
	if err != nil {
		return fmt.Errorf("manifest: min_host_version 无效: %w", err)
	}
	version, err := semver.NewVersion(hostVersion)
	if err != nil {
		return fmt.Errorf("manifest: 宿主版本无效: %w", err)
	}
	if !constraint.Check(version) {
		return fmt.Errorf("%w: 插件要求 %s,宿主为 %s", ErrHostVersionIncompatible, m.MinHostVersion, hostVersion)
	}
	return nil
}

// Validate 校验清单及插件目录中的文件布局。
func (m *Manifest) Validate() error {
	if m == nil {
		return &ValidationError{Errors: []FieldError{{Field: "manifest", Message: "不能为空"}}}
	}
	var fields []FieldError
	fields = append(fields, validateName(m.Name, m.Dir)...)
	if m.Version == "" {
		fields = append(fields, FieldError{Field: "version", Message: "必填"})
	} else if _, err := semver.StrictNewVersion(m.Version); err != nil {
		fields = append(fields, FieldError{Field: "version", Message: "须为 SemVer"})
	}
	if m.MinHostVersion != "" {
		if _, err := semver.NewConstraint(m.MinHostVersion); err != nil {
			fields = append(fields, FieldError{Field: "min_host_version", Message: "须为 SemVer 范围"})
		}
	}
	fields = append(fields, validateBinary(m.Dir, m.Binary)...)
	fields = append(fields, validateServices(m.Services)...)
	fields = append(fields, validateGrants(m.Grants)...)
	fields = append(fields, validateExtensions(m.Services, m.Extensions)...)
	fields = append(fields, validateUI(m.Dir, m.UI)...)
	if len(fields) > 0 {
		return &ValidationError{Errors: fields}
	}
	return nil
}

func validateName(name, dir string) []FieldError {
	if name == "" {
		return []FieldError{{Field: "name", Message: "必填"}}
	}
	if !isKebabCase(name) {
		return []FieldError{{Field: "name", Message: "须为 kebab-case"}}
	}
	if dir != "" && filepath.Base(dir) != name {
		return []FieldError{{Field: "name", Message: "须与插件目录名一致"}}
	}
	return nil
}

func validateBinary(dir, binary string) []FieldError {
	if binary == "" {
		return []FieldError{{Field: "binary", Message: "必填"}}
	}
	if dir == "" {
		return nil
	}
	path := filepath.Join(dir, binary)
	info, err := os.Stat(path)
	if err != nil {
		return []FieldError{{Field: "binary", Message: "文件不存在: " + binary}}
	}
	if info.IsDir() {
		return []FieldError{{Field: "binary", Message: "路径为目录: " + binary}}
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o111 == 0 {
		return []FieldError{{Field: "binary", Message: "文件不可执行: " + binary}}
	}
	return nil
}

func validateServices(services []string) []FieldError {
	if len(services) == 0 {
		return []FieldError{{Field: "services", Message: "至少声明一个服务"}}
	}
	var fields []FieldError
	seen := make(map[string]struct{}, len(services))
	for i, service := range services {
		field := fmt.Sprintf("services[%d]", i)
		if service == "" {
			fields = append(fields, FieldError{Field: field, Message: "不能为空"})
			continue
		}
		if _, exists := seen[service]; exists {
			fields = append(fields, FieldError{Field: field, Message: "重复: " + service})
		}
		seen[service] = struct{}{}
	}
	return fields
}

func validateGrants(grants []string) []FieldError {
	var fields []FieldError
	seen := make(map[string]struct{}, len(grants))
	for i, grant := range grants {
		field := fmt.Sprintf("grants[%d]", i)
		if grant == "" {
			fields = append(fields, FieldError{Field: field, Message: "不能为空"})
			continue
		}
		if _, exists := seen[grant]; exists {
			fields = append(fields, FieldError{Field: field, Message: "重复: " + grant})
		}
		seen[grant] = struct{}{}
	}
	return fields
}

func validateExtensions(services []string, extensions []Extension) []FieldError {
	var fields []FieldError
	declared := make(map[string]struct{}, len(services))
	for _, s := range services {
		declared[s] = struct{}{}
	}
	seen := make(map[string]struct{}, len(extensions))
	for i, ext := range extensions {
		field := fmt.Sprintf("extensions[%d]", i)
		if ext.Name == "" {
			fields = append(fields, FieldError{Field: field + ".name", Message: "必填"})
		}
		if ext.Point == "" {
			fields = append(fields, FieldError{Field: field + ".point", Message: "必填"})
		}
		if ext.ServiceName == "" {
			fields = append(fields, FieldError{Field: field + ".service_name", Message: "必填"})
		} else if _, ok := declared[ext.ServiceName]; !ok {
			fields = append(fields, FieldError{Field: field + ".service_name", Message: "未在 services 声明: " + ext.ServiceName})
		}
		key := ext.Point + "/" + ext.Name
		if _, exists := seen[key]; exists {
			fields = append(fields, FieldError{Field: field, Message: "重复: " + key})
		}
		seen[key] = struct{}{}
	}
	return fields
}

func validateUI(dir string, ui *UI) []FieldError {
	if ui == nil {
		return nil
	}
	var fields []FieldError
	if ui.Entry == "" {
		fields = append(fields, FieldError{Field: "ui.entry", Message: "ui 存在时必填"})
	} else if dir != "" {
		entry := filepath.Join(dir, ui.Entry)
		info, err := os.Stat(entry)
		if err != nil {
			fields = append(fields, FieldError{Field: "ui.entry", Message: "文件不存在: " + ui.Entry})
		} else if info.IsDir() {
			fields = append(fields, FieldError{Field: "ui.entry", Message: "路径为目录: " + ui.Entry})
		}
	}
	if ui.MountPath != "" && !strings.HasPrefix(ui.MountPath, "/") {
		fields = append(fields, FieldError{Field: "ui.mount_path", Message: "须以 / 开头"})
	}
	if ui.RoutesFile != "" && dir != "" {
		if _, err := LoadRoutes(filepath.Join(dir, ui.RoutesFile)); err != nil {
			fields = append(fields, FieldError{Field: "ui.routes_file", Message: err.Error()})
		}
	}
	return fields
}

func isKebabCase(value string) bool {
	if value == "" || value != strings.ToLower(value) {
		return false
	}
	for _, part := range strings.Split(value, "-") {
		if part == "" {
			return false
		}
		for _, char := range part {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') {
				return false
			}
		}
	}
	return true
}

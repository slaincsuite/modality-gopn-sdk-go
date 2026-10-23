// registry.go 管理插件进程内扩展实现声明。
package extensionpoint

import (
	"fmt"
	"sync"
)

// Registry 是插件内扩展声明注册表;实际启用仍由宿主鉴权与编排。
type Registry struct {
	mu         sync.RWMutex
	extensions []Extension
}

// NewRegistry 创建空注册表。
func NewRegistry() *Registry { return &Registry{} }

// Register 登记扩展实现;同扩展点同名实现不可重复。
func (r *Registry) Register(extension Extension) error {
	if !extension.Valid() {
		return fmt.Errorf("extensionpoint: 扩展声明无效: %s", extension.Name)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.extensions {
		if existing.Point == extension.Point && existing.Name == extension.Name {
			return fmt.Errorf("extensionpoint: %s/%s 已注册", extension.Point, extension.Name)
		}
	}
	r.extensions = append(r.extensions, extension)
	return nil
}

// Extensions 返回扩展声明的只读副本。
func (r *Registry) Extensions() []Extension {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Extension, len(r.extensions))
	copy(out, r.extensions)
	return out
}

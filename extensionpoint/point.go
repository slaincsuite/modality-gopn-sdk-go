// point.go 定义插件扩展点的通用注册语义。
package extensionpoint

// Mode 表示扩展点的实例数量约束。
type Mode string

const (
	// Singleton 同一扩展点仅启用一个实现。
	Singleton Mode = "singleton"
	// MultiInstance 同一扩展点可启用多个实现。
	MultiInstance Mode = "multi-instance"
)

// Definition 描述宿主开放的扩展点。
type Definition struct {
	// Name 扩展点稳定名称,建议格式 module.point。
	Name string
	// Module 提供扩展点的模块名。
	Module string
	// DisplayName 控制台展示名。
	DisplayName string
	// Description 扩展点用途。
	Description string
	// Mode 注册模式。
	Mode Mode
}

// Extension 描述插件贡献的一个扩展实现。
type Extension struct {
	// Name 扩展实现稳定名称。
	Name string
	// Point 扩展点名称。
	Point string
	// DisplayName 控制台展示名。
	DisplayName string
	// Description 扩展实现说明。
	Description string
	// Priority 多实例执行顺序,数值大者先执行。
	Priority int
	// ServiceName 提供 RPC 实现的服务名。
	ServiceName string
}

// Valid 判断扩展点声明是否有效。
func (d Definition) Valid() bool {
	return d.Name != "" && d.Module != "" &&
		(d.Mode == Singleton || d.Mode == MultiInstance)
}

// Valid 判断扩展实现声明是否有效。
func (e Extension) Valid() bool {
	return e.Name != "" && e.Point != "" && e.ServiceName != ""
}

// capability.go 定义宿主开放能力的元数据与访问模式。
package host

// AccessMode 表示能力允许的访问类型。
type AccessMode string

const (
	// AccessRead 仅允许读取,适用于模块元数据、配置快照和注册表。
	AccessRead AccessMode = "read"
	// AccessControlledWrite 允许调用宿主显式声明且逐次鉴权的受控写入。
	AccessControlledWrite AccessMode = "controlled-write"
)

// Capability 描述宿主开放给插件的一个能力。
type Capability struct {
	// Name 能力全名,格式为 module.function,例如 upms.user.list。
	Name string
	// Module 提供能力的模块名。
	Module string
	// Description 面向插件开发者与平台管理员的说明。
	Description string
	// Mode 最小权限访问模式。
	Mode AccessMode
}

// Valid 判断能力声明是否满足最小要素。
func (c Capability) Valid() bool {
	return c.Name != "" && c.Module != "" &&
		(c.Mode == AccessRead || c.Mode == AccessControlledWrite)
}

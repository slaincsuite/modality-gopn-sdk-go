// authorization.go 定义插件调用宿主时携带的授权凭据。
package host

// Authorization 是一次能力授权的描述;宿主仍须重新鉴权,不得信任该结构本身。
type Authorization struct {
	// Plugin 发起调用的插件名。
	Plugin string
	// Capability 已声明的能力。
	Capability Capability
}

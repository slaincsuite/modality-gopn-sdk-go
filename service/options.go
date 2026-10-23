// options.go 定义插件构建选项。
package service

import (
	"github.com/hashicorp/go-plugin"

	"github.com/slaincsuite/modality-gopn-sdk-go/core"
)

// WithLogger 注入插件日志实现;nil 会被忽略。
func WithLogger(logger core.Logger) Option {
	return func(plugin *Plugin) {
		if logger == nil {
			return
		}
		plugin.config.Logger = logger
		plugin.loggerSet = true
	}
}

// WithLogLevel 设置默认日志级别;已显式注入 Logger 时不覆盖实现。
func WithLogLevel(level core.LogLevel) Option {
	return func(plugin *Plugin) {
		plugin.config.LogLevel = level
		if !plugin.loggerSet {
			plugin.config.Logger = core.NewLogger(level, nil)
		}
	}
}

// WithProtocolVersion 覆盖握手协议版本;必须与宿主一致。
func WithProtocolVersion(version int) Option {
	return func(plugin *Plugin) {
		plugin.config.Handshake.ProtocolVersion = version
	}
}

// WithBroker 注册宿主反向调用处理器;由宿主持有相同 broker ID 时使用。
func WithBroker(handler func(*plugin.GRPCBroker) error) Option {
	return func(plugin *Plugin) {
		plugin.broker = handler
	}
}

// WithTestConfig 设置 go-plugin 测试模式;仅限测试使用。
func WithTestConfig(test *plugin.ServeTestConfig) Option {
	return func(plugin *Plugin) {
		plugin.test = test
	}
}

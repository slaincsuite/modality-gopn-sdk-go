// Package modality 提供 modality 插件二进制的组装与运行入口。
package service

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/slaincsuite/modality-gopn-sdk-go/core"
)

// RegisterFunc 注册一个 proto 生成的 gRPC 服务。
// 可直接在闭包中调用 protoc-gen-go-grpc 生成的 RegisterXxxServer。
type RegisterFunc func(registrar grpc.ServiceRegistrar)

// Service 描述插件提供的一个 gRPC 服务。
type Service struct {
	// Name 与 plugin.json 的 services 元素一致。
	Name string
	// Register 注册服务描述符与实现。
	Register RegisterFunc
}

// NewService 创建服务声明。
func NewService(name string, register RegisterFunc) Service {
	return Service{Name: name, Register: register}
}

// Option 修改插件运行配置。
type Option func(*Plugin)

// Plugin 插件进程配置;由 NewPlugin 创建并调用 Serve 阻塞运行。
type Plugin struct {
	mu        sync.Mutex
	services  []Service
	names     map[string]struct{}
	config    core.Config
	loggerSet bool
	broker    func(*plugin.GRPCBroker) error
	test      *plugin.ServeTestConfig
}

// NewPlugin 创建插件实例。
func NewPlugin(options ...Option) *Plugin {
	plugin := &Plugin{
		names:  make(map[string]struct{}),
		config: core.NewConfig(),
	}
	for _, option := range options {
		if option != nil {
			option(plugin)
		}
	}
	return plugin
}

// Register 注册服务;名称或注册函数为空、名称重复时返回错误。
func (p *Plugin) Register(service Service) error {
	if err := validateService(service); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.names[service.Name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateService, service.Name)
	}
	p.names[service.Name] = struct{}{}
	p.services = append(p.services, service)
	return nil
}

// Services 返回服务声明的只读副本。
func (p *Plugin) Services() []Service {
	p.mu.Lock()
	defer p.mu.Unlock()
	services := make([]Service, len(p.services))
	copy(services, p.services)
	return services
}

// Serve 注册服务并阻塞运行;进程退出或 go-plugin 结束时返回。
// 允许零服务:只提供应用与操作目录(action.json)/UI 的插件合法,
// 此时仍以 go-plugin 握手 + 内建健康服务常驻,供宿主管理生命周期。
func (p *Plugin) Serve() error {
	services := p.Services()
	if err := p.config.Validate(); err != nil {
		return err
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  uint(p.config.Handshake.ProtocolVersion),
			MagicCookieKey:   p.config.Handshake.MagicCookieKey,
			MagicCookieValue: p.config.Handshake.MagicCookieValue,
		},
		Plugins: map[string]plugin.Plugin{
			pluginName: &grpcPlugin{
				services: services,
				broker:   p.broker,
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     newHCLogger(p.config.Logger, p.config.LogLevel),
		Test:       p.test,
	})
	return nil
}

// validateService 校验单个服务声明。
func validateService(service Service) error {
	if service.Name == "" {
		return fmt.Errorf("%w: Name 必填", ErrInvalidService)
	}
	if service.Register == nil {
		return fmt.Errorf("%w: %s: Register 必填", ErrInvalidService, service.Name)
	}
	return nil
}

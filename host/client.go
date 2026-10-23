// client.go 提供默认拒绝的宿主只读访问客户端。
package host

import (
	"context"
	"fmt"
)

// Request 是一次宿主访问请求;Payload 由能力契约约定。
type Request struct {
	// Capability 能力全名。
	Capability string
	// Payload 能力输入,仅允许 JSON 兼容数据。
	Payload any
}

// Response 是宿主返回的能力结果;Data 由能力契约约定。
type Response struct {
	// Data 能力输出。
	Data any
}

// Caller 执行已授权调用;传输可由 gRPC broker 或类型化 SDK 适配实现。
type Caller interface {
	// Call 必须把 Authorization 作为不可伪造调用凭据交给宿主。
	Call(ctx context.Context, authorization Authorization, request Request) (Response, error)
}

// Config 插件侧宿主访问配置。
type Config struct {
	// PluginName 插件标识。
	PluginName string
	// Capabilities 宿主声明并授予的能力集。
	Capabilities []Capability
}

// Client 是默认拒绝的宿主访问门面。
type Client struct {
	config Config
	caller Caller
}

// New 创建宿主客户端;caller 为空时所有调用返回宿主不可用。
func New(config Config, caller Caller) *Client {
	return &Client{config: config, caller: caller}
}

// Capabilities 返回能力声明的只读副本。
func (c *Client) Capabilities() []Capability {
	out := make([]Capability, len(c.config.Capabilities))
	copy(out, c.config.Capabilities)
	return out
}

// Authorize 为指定能力生成调用凭据;SDK 侧不授予宿主未声明的能力。
func (c *Client) Authorize(capability string) (Authorization, error) {
	for _, item := range c.config.Capabilities {
		if item.Name != capability {
			continue
		}
		if !item.Valid() {
			return Authorization{}, fmt.Errorf("%w: %s", ErrInvalidCapability, capability)
		}
		return Authorization{Plugin: c.config.PluginName, Capability: item}, nil
	}
	return Authorization{}, fmt.Errorf("%w: %s", ErrForbidden, capability)
}

// Call 执行授权调用;未授权、宿主未配置或能力非法均失败。
func (c *Client) Call(ctx context.Context, capability string, payload any) (Response, error) {
	if c.caller == nil {
		return Response{}, ErrHostUnavailable
	}
	authorization, err := c.Authorize(capability)
	if err != nil {
		return Response{}, err
	}
	return c.caller.Call(ctx, authorization, Request{Capability: capability, Payload: payload})
}

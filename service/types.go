// types.go 统一导出插件服务、日志、扩展与宿主访问公共类型。
package service

import (
	"github.com/slaincsuite/modality-gopn-sdk-go/core"
	"github.com/slaincsuite/modality-gopn-sdk-go/host"
)

// Logger 插件日志接口,与 go-plugin 日志桥接。
type Logger = core.Logger

// LogLevel 日志级别。
type LogLevel = core.LogLevel

const (
	// LogLevelDebug 输出调试及以上日志。
	LogLevelDebug = core.LogLevelDebug
	// LogLevelInfo 输出信息及以上日志。
	LogLevelInfo = core.LogLevelInfo
	// LogLevelWarn 输出告警及以上日志。
	LogLevelWarn = core.LogLevelWarn
	// LogLevelError 只输出错误日志。
	LogLevelError = core.LogLevelError
)

// NopLogger 返回丢弃日志的实现。
func NopLogger() Logger {
	return core.NopLogger()
}

// HostCaller 反向调用宿主开放能力;实现必须自带鉴权上下文与调用审计。
type HostCaller = host.Caller

// Authorization 授权后的宿主调用凭据;不包含长期令牌。
type Authorization = host.Authorization

// Capability 宿主开放能力声明。
type Capability = host.Capability

// AccessMode 能力访问模式。
type AccessMode = host.AccessMode

// HostConfig 宿主连接配置。
type HostConfig = host.Config

const (
	// AccessRead 仅允许读取元数据或配置。
	AccessRead = host.AccessRead
	// AccessControlledWrite 允许写入宿主显式声明的受控端点。
	AccessControlledWrite = host.AccessControlledWrite
)

// NewHostClient 创建宿主只读客户端;未授权调用返回错误。
func NewHostClient(config HostConfig, caller HostCaller) *host.Client {
	return host.New(config, caller)
}

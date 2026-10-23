// core/config.go 汇总 SDK 握手、日志级别与日志实现。
package core

import "fmt"

// Config SDK 运行时配置。
type Config struct {
	Handshake Handshake
	Logger    Logger
	LogLevel  LogLevel
}

// Option 修改 SDK 配置。
type Option func(*Config)

// NewConfig 创建应用默认值后的配置。
func NewConfig(options ...Option) Config {
	cfg := Config{
		Handshake: DefaultHandshake(),
		LogLevel:  LogLevelInfo,
	}
	for _, option := range options {
		if option != nil {
			option(&cfg)
		}
	}
	if cfg.Logger == nil {
		cfg.Logger = NewLogger(cfg.LogLevel, nil)
	}
	return cfg
}

// WithLogger 覆盖插件日志实现;nil 会被忽略。
func WithLogger(logger Logger) Option {
	return func(cfg *Config) {
		if logger != nil {
			cfg.Logger = logger
		}
	}
}

// WithLogLevel 设置默认日志级别。
func WithLogLevel(level LogLevel) Option {
	return func(cfg *Config) {
		cfg.LogLevel = level
	}
}

// WithProtocolVersion 覆盖握手协议版本;必须与宿主配置一致。
func WithProtocolVersion(version int) Option {
	return func(cfg *Config) {
		cfg.Handshake.ProtocolVersion = version
	}
}

// Validate 校验握手、日志级别与日志实现。
func (c Config) Validate() error {
	if err := c.Handshake.Validate(); err != nil {
		return err
	}
	if c.LogLevel < LogLevelDebug || c.LogLevel > LogLevelError {
		return fmt.Errorf("core: 非法日志级别 %d", c.LogLevel)
	}
	if c.Logger == nil {
		return fmt.Errorf("core: Logger 必填")
	}
	return nil
}

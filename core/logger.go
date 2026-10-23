// core/logger.go 提供插件侧日志接口及默认 stderr 实现。
package core

import (
	"fmt"
	"io"
	"log"
	"os"
)

// LogLevel 控制默认日志实现的最低输出级别。
type LogLevel int

const (
	// LogLevelDebug 输出调试、信息、告警和错误日志。
	LogLevelDebug LogLevel = iota
	// LogLevelInfo 输出信息、告警和错误日志。
	LogLevelInfo
	// LogLevelWarn 输出告警和错误日志。
	LogLevelWarn
	// LogLevelError 只输出错误日志。
	LogLevelError
)

// Logger 插件运行时日志接口;实现应并发安全。
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type standardLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger 创建带级别过滤的默认日志器。
func NewLogger(level LogLevel, output io.Writer) Logger {
	if output == nil {
		output = os.Stderr
	}
	return &standardLogger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}
}

// DefaultLogger 返回写入 stderr、级别为 Info 的日志器。
func DefaultLogger() Logger {
	return NewLogger(LogLevelInfo, os.Stderr)
}

func (l *standardLogger) Debugf(format string, args ...any) {
	l.write(LogLevelDebug, "DEBUG", format, args...)
}

func (l *standardLogger) Infof(format string, args ...any) {
	l.write(LogLevelInfo, "INFO", format, args...)
}

func (l *standardLogger) Warnf(format string, args ...any) {
	l.write(LogLevelWarn, "WARN", format, args...)
}

func (l *standardLogger) Errorf(format string, args ...any) {
	l.write(LogLevelError, "ERROR", format, args...)
}

func (l *standardLogger) write(level LogLevel, label, format string, args ...any) {
	if level < l.level {
		return
	}
	l.logger.Printf("["+label+"] "+format, args...)
}

type nopLogger struct{}

// NopLogger 返回丢弃全部日志的实现,便于测试与嵌入式场景。
func NopLogger() Logger { return nopLogger{} }

func (nopLogger) Debugf(string, ...any) {}
func (nopLogger) Infof(string, ...any)  {}
func (nopLogger) Warnf(string, ...any)  {}
func (nopLogger) Errorf(string, ...any) {}

// String 返回日志级别的可读名称。
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return fmt.Sprintf("LogLevel(%d)", int(l))
	}
}

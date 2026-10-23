// core/logger.go 提供插件侧日志接口及默认 stderr 实现。
package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	// With 返回携带结构化键值对的子日志器;实现应保留原有上下文。
	With(kv ...any) Logger
}

type standardLogger struct {
	level  LogLevel
	logger *log.Logger
	fields []any
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

func (l *standardLogger) With(kv ...any) Logger {
	return &standardLogger{
		level:  l.level,
		logger: l.logger,
		fields: append(append([]any{}, l.fields...), kv...),
	}
}

func (l *standardLogger) write(level LogLevel, label, format string, args ...any) {
	if level < l.level {
		return
	}
	msg := "[" + label + "] " + format
	if len(l.fields) > 0 {
		msg += " " + formatPairs(l.fields)
	}
	l.logger.Printf(msg, args...)
}

func formatPairs(kv []any) string {
	parts := make([]string, 0, (len(kv)+1)/2)
	for i := 0; i < len(kv); i += 2 {
		if i+1 < len(kv) {
			parts = append(parts, fmt.Sprintf("%v=%v", kv[i], kv[i+1]))
		} else {
			parts = append(parts, fmt.Sprintf("%v", kv[i]))
		}
	}
	return strings.Join(parts, " ")
}

type nopLogger struct{}

// NopLogger 返回丢弃全部日志的实现,便于测试与嵌入式场景。
func NopLogger() Logger { return nopLogger{} }

func (nopLogger) Debugf(string, ...any) {}
func (nopLogger) Infof(string, ...any)  {}
func (nopLogger) Warnf(string, ...any)  {}
func (nopLogger) Errorf(string, ...any) {}
func (nopLogger) With(...any) Logger    { return nopLogger{} }

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

// logger_bridge.go 将 SDK 日志适配为 go-plugin 的 hclog。
package service

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/slaincsuite/modality-gopn-sdk-go/core"
)

type hclogAdapter struct {
	logger core.Logger
	name   string
	args   []any
	level  core.LogLevel
}

func newHCLogger(logger core.Logger, level core.LogLevel) hclog.Logger {
	return &hclogAdapter{logger: logger, level: level}
}

func (l *hclogAdapter) Log(level hclog.Level, msg string, args ...any) {
	line := formatLine(l.name, msg, l.args, args)
	switch level {
	case hclog.Trace, hclog.Debug:
		l.logger.Debugf("%s", line)
	case hclog.Info:
		l.logger.Infof("%s", line)
	case hclog.Warn:
		l.logger.Warnf("%s", line)
	case hclog.Error:
		l.logger.Errorf("%s", line)
	}
}

func (l *hclogAdapter) Trace(msg string, args ...any) { l.Log(hclog.Trace, msg, args...) }
func (l *hclogAdapter) Debug(msg string, args ...any) { l.Log(hclog.Debug, msg, args...) }
func (l *hclogAdapter) Info(msg string, args ...any)  { l.Log(hclog.Info, msg, args...) }
func (l *hclogAdapter) Warn(msg string, args ...any)  { l.Log(hclog.Warn, msg, args...) }
func (l *hclogAdapter) Error(msg string, args ...any) { l.Log(hclog.Error, msg, args...) }

func (l *hclogAdapter) IsTrace() bool { return false }
func (l *hclogAdapter) IsDebug() bool { return l.level <= core.LogLevelDebug }
func (l *hclogAdapter) IsInfo() bool  { return l.level <= core.LogLevelInfo }
func (l *hclogAdapter) IsWarn() bool  { return l.level <= core.LogLevelWarn }
func (l *hclogAdapter) IsError() bool { return l.level <= core.LogLevelError }

func (l *hclogAdapter) ImpliedArgs() []any {
	return l.args
}

func (l *hclogAdapter) With(args ...any) hclog.Logger {
	merged := make([]any, 0, len(l.args)+len(args))
	merged = append(merged, l.args...)
	merged = append(merged, args...)
	return &hclogAdapter{logger: l.logger, name: l.name, args: merged, level: l.level}
}

func (l *hclogAdapter) Named(name string) hclog.Logger {
	return &hclogAdapter{logger: l.logger, name: l.name + name, args: l.args, level: l.level}
}

func (l *hclogAdapter) Name() string { return l.name }

func (l *hclogAdapter) ResetNamed(name string) hclog.Logger {
	return &hclogAdapter{logger: l.logger, name: name, args: l.args, level: l.level}
}

func (l *hclogAdapter) SetLevel(hclog.Level) {}
func (l *hclogAdapter) GetLevel() hclog.Level {
	switch l.level {
	case core.LogLevelDebug:
		return hclog.Debug
	case core.LogLevelWarn:
		return hclog.Warn
	case core.LogLevelError:
		return hclog.Error
	default:
		return hclog.Info
	}
}
func (l *hclogAdapter) StandardLogger(*hclog.StandardLoggerOptions) *log.Logger {
	return log.New(logWriter{logger: l.logger}, "", 0)
}
func (l *hclogAdapter) StandardWriter(*hclog.StandardLoggerOptions) io.Writer {
	return logWriter{logger: l.logger}
}

type logWriter struct {
	logger core.Logger
}

func (w logWriter) Write(data []byte) (int, error) {
	w.logger.Infof("%s", string(data))
	return len(data), nil
}

func formatLine(name, message string, implied, args []any) string {
	line := message
	if name != "" {
		line = name + ": " + line
	}
	return line + formatArgs(implied) + formatArgs(args)
}

func formatArgs(args []any) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, (len(args)+1)/2)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			parts = append(parts, fmt.Sprint(args[i]))
			break
		}
		parts = append(parts, fmt.Sprintf("%v=%v", args[i], args[i+1]))
	}
	return " " + strings.Join(parts, " ")
}

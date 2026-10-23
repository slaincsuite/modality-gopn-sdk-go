// core 单测:握手、配置校验、日志级别过滤。
package core

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestDefaultHandshake(t *testing.T) {
	handshake := DefaultHandshake()
	if err := handshake.Validate(); err != nil {
		t.Fatalf("DefaultHandshake().Validate() = %v", err)
	}
	if handshake.ProtocolVersion != 1 {
		t.Errorf("ProtocolVersion = %d, want 1", handshake.ProtocolVersion)
	}
	if handshake.MagicCookieKey != "GOPN_PLUGIN" {
		t.Errorf("MagicCookieKey = %q", handshake.MagicCookieKey)
	}
	if handshake.MagicCookieValue != "gopn-plugin-handshake" {
		t.Errorf("MagicCookieValue = %q", handshake.MagicCookieValue)
	}
}

func TestHandshakeValidate(t *testing.T) {
	tests := []struct {
		name      string
		handshake Handshake
	}{
		{name: "version", handshake: Handshake{MagicCookieKey: "k", MagicCookieValue: "v"}},
		{name: "key", handshake: Handshake{ProtocolVersion: 1, MagicCookieValue: "v"}},
		{name: "value", handshake: Handshake{ProtocolVersion: 1, MagicCookieKey: "k"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handshake.Validate()
			if !errors.Is(err, ErrInvalidHandshake) {
				t.Fatalf("Validate() = %v, want ErrInvalidHandshake", err)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewConfig(WithLogLevel(LogLevel(99)))
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid log level error")
	}
}

func TestNewConfigUsesRequestedLevel(t *testing.T) {
	cfg := NewConfig(WithLogLevel(LogLevelDebug))
	if cfg.LogLevel != LogLevelDebug {
		t.Fatalf("LogLevel = %v, want debug", cfg.LogLevel)
	}
	if cfg.Logger == nil {
		t.Fatal("expected default logger")
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var output bytes.Buffer
	logger := NewLogger(LogLevelWarn, &output)
	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warnf("warn")
	logger.Errorf("error")
	text := output.String()
	if strings.Contains(text, "debug") || strings.Contains(text, "info") {
		t.Fatalf("低级别日志不应输出: %q", text)
	}
	if !strings.Contains(text, "warn") || !strings.Contains(text, "error") {
		t.Fatalf("缺少 warn/error 日志: %q", text)
	}
}

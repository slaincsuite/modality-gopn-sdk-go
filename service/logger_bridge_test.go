// logger_test.go 验证 hclog 适配器的级别、字段与前缀。
package service

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/slaincsuite/modality-gopn-sdk-go/core"
)

type recordingLogger struct {
	mu    sync.Mutex
	lines []string
}

func (l *recordingLogger) Debugf(format string, args ...any) {
	l.record("debug", format, args...)
}

func (l *recordingLogger) Infof(format string, args ...any) {
	l.record("info", format, args...)
}

func (l *recordingLogger) Warnf(format string, args ...any) {
	l.record("warn", format, args...)
}

func (l *recordingLogger) Errorf(format string, args ...any) {
	l.record("error", format, args...)
}

func (l *recordingLogger) record(level, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, level+": "+fmt.Sprintf(format, args...))
}

func TestHCLogAdapter(t *testing.T) {
	recorder := &recordingLogger{}
	logger := newHCLogger(recorder, core.LogLevelWarn)
	if logger.IsDebug() {
		t.Fatal("debug should be disabled at warn level")
	}
	if !logger.IsWarn() || !logger.IsError() {
		t.Fatal("warn/error should be enabled at warn level")
	}
	if logger.GetLevel() != hclog.Warn {
		t.Fatalf("GetLevel() = %v, want warn", logger.GetLevel())
	}
	logger.Named("plugin").Warn("starting", "id", 7)
	if len(recorder.lines) != 1 {
		t.Fatalf("lines = %v", recorder.lines)
	}
	line := recorder.lines[0]
	if !strings.Contains(line, "plugin: starting") || !strings.Contains(line, "id=7") {
		t.Fatalf("line = %q", line)
	}
}

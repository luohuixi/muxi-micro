package logx

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
)

func TestNewStdLogger(t *testing.T) {
	l := NewStdLogger()
	if l == nil {
		t.Fatal("expected logger, got nil")
	}
	if l.logger == nil {
		t.Fatal("expected underlying log.Logger initialized")
	}
}

func TestOutputLevels(t *testing.T) {
	l := NewStdLogger()

	fields := logger.Field{"key": "value"}
	l.Info("info msg", fields)
	l.Debug("debug msg")
	l.Warn("warn msg")
	l.Error("error msg")

}

func TestWithFields(t *testing.T) {
	l := NewStdLogger()

	// 添加初始字段
	l2 := l.With(logger.Field{"request_id": "1234"}).(*StdLogger)
	if v, ok := l2.fields["request_id"]; !ok || v != "1234" {
		t.Fatalf("expected field request_id=1234, got %v", l2.fields)
	}

	// 再添加新字段（测试合并）
	l3 := l2.With(logger.Field{"user": "alice"}).(*StdLogger)
	if v, ok := l3.fields["user"]; !ok || v != "alice" {
		t.Fatalf("expected user=alice, got %v", l3.fields)
	}
}

func TestOutputFieldMerging(t *testing.T) {
	l := NewStdLogger()
	buf := &bytes.Buffer{}
	l.logger = log.New(buf, "", 0)

	l.fields["base"] = "yes"
	l.output("INFO", "test", logger.Field{"child": "ok"})

	out := buf.String()
	if !strings.Contains(out, "\"base\":\"yes\"") || !strings.Contains(out, "\"child\":\"ok\"") {
		t.Fatalf("field merge failed: %s", out)
	}
}

func TestSync(t *testing.T) {
	l := NewStdLogger()
	if err := l.Sync(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

package logx

import (
	"encoding/json"
	"fmt"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type StdLogger struct {
	logger *log.Logger
	fields logger.Field
	mu     sync.RWMutex
}

// NewStdLogger 创建一个基于标准库 log 的 Logger 实现，非常轻量，适合简单场景/日志的兜底场景使用
func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", 0),
		fields: make(logger.Field),
	}
}

func (l *StdLogger) Info(msg string, fields ...logger.Field) {
	l.output("INFO", msg, fields...)
}

func (l *StdLogger) Debug(msg string, fields ...logger.Field) {
	l.output("DEBUG", msg, fields...)
}

func (l *StdLogger) Warn(msg string, fields ...logger.Field) {
	l.output("WARN", msg, fields...)
}

func (l *StdLogger) Error(msg string, fields ...logger.Field) {
	l.output("ERROR", msg, fields...)
}

func (l *StdLogger) Fatal(msg string, fields ...logger.Field) {
	l.output("FATAL", msg, fields...)
	os.Exit(1)
}

// With 创建带上下文的 logger
func (l *StdLogger) With(fields ...logger.Field) logger.Logger {
	l.mu.RLock()
	base := make(logger.Field, len(l.fields))
	for k, v := range l.fields {
		base[k] = v
	}
	l.mu.RUnlock()

	for _, f := range fields {
		for k, v := range f {
			base[k] = v
		}
	}

	return &StdLogger{
		logger: l.logger,
		fields: base,
	}
}

// Sync 在标准库中没有缓冲区，这里直接返回 nil
func (l *StdLogger) Sync() error {
	return nil
}

func (l *StdLogger) output(level, msg string, fields ...logger.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 合并全局与传入字段
	merged := make(logger.Field)
	for k, v := range l.fields {
		merged[k] = v
	}
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}

	// 获取时间与调用信息
	now := time.Now().Format("2006-01-02T15:04:05.000-0700")
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		short := file
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			short = file[idx+1:]
		}
		caller = fmt.Sprintf("%s:%d", short, line)
	}

	// 构建 JSON map
	logMap := map[string]interface{}{
		"level":      level,
		"@timestamp": now,
		"caller":     caller,
		"msg":        msg,
	}

	// 合并自定义字段
	for k, v := range merged {
		logMap[k] = v
	}

	// 序列化为 JSON
	jsonBytes, err := json.Marshal(logMap)
	if err != nil {
		// 兜底输出
		l.logger.Printf("[LOGGING_ERROR] level=%s msg=%s err=%v", level, msg, err)
		return
	}

	// 输出单行 JSON 日志
	l.logger.Println(string(jsonBytes))
}

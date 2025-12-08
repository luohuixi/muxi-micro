package log

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"google.golang.org/grpc"
)

const (
	LogIDKey     = "Micro-LogID"
	LoggerKey    = "Micro-Logger"
	NameKey      = "Micro-Name"
	DefaultName  = "unknown"
	DefaultLogID = "without logID"
)

// 生成日志 ID
func genLogID(prefix string) string {
	// 当前时间纳秒 + 随机字节混合
	timeBytes := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))

	//生成随机短id
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	randomBytes := hex.EncodeToString(b)
	combined := append(timeBytes, randomBytes...)

	// SHA1 哈希处理
	hash := sha1.Sum(combined)
	shortHash := hex.EncodeToString(hash[:8]) // 取前8字节（16个字符）

	logID := fmt.Sprintf("%s-%s", prefix, shortHash)
	return logID
}

// 设置到响应中需要
func SetLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, LogIDKey, logID)
}

// 为了保证获取的便利性这里用的是context.Context
func GetLogID(ctx context.Context) string {
	value, ok := ctx.Value(LogIDKey).(string)
	if !ok {
		return DefaultLogID
	}
	return value
}

// 用于设置在上下文中获取携带了特殊信息的日志,主动打印需要
func SetLogger(ctx context.Context, logger logger.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// 用于获取在上下文中获取携带了特殊信息的日志
func GetLogger(ctx context.Context) logger.Logger {
	ginLogger, ok := ctx.Value(LoggerKey).(logger.Logger)
	if !ok {
		return nil // 如果不存在需要返回一个空指针
	}
	return ginLogger
}

func GlobalLoggerServerInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logID := GetLogID(ctx)

		if logID == DefaultLogID {
			logID = genLogID(GetGlobalName(ctx))
		}

		newCtx := SetLogID(ctx, logID)
		l = l.With(
			logger.Field{"logID": logID},
		)

		newCtx2 := SetLogger(newCtx, l)

		return handler(newCtx2, req)
	}
}

func GlobalNameServerInterceptor(name string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx := context.WithValue(ctx, NameKey, name)
		return handler(newCtx, req)
	}
}

func GetGlobalName(ctx context.Context) string {
	value, ok := ctx.Value(NameKey).(string)
	if !ok {
		return DefaultName
	}
	return value
}

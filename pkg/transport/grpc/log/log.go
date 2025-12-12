package log

import (
	"context"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"google.golang.org/grpc"
)

const (
	LogIDKey     = "Micro-LogID"
	LoggerKey    = "Micro-Logger"
	DefaultLogID = "without logID"
)

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
		l = l.With(
			logger.Field{"logID": logID},
		)

		newCtx := SetLogger(ctx, l)
		return handler(newCtx, req)
	}
}

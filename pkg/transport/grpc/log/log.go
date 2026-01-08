package log

import (
	"context"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	LogIDKey     = "Micro-LogID"
	LoggerKey    = "Micro-Logger"
	DefaultLogID = "without logID"
)

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
		logID := DefaultLogID
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if values := md.Get(LogIDKey); len(values) > 0 {
				logID = values[0]
			}
		}
		l = l.With(
			logger.Field{"logID": logID},
		)

		newCtx := SetLogID(ctx, logID)
		newCtx = SetLogger(newCtx, l)
		return handler(newCtx, req)
	}
}

func GlobalLoggerClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		logID := GetLogID(ctx)

		ctx = metadata.AppendToOutgoingContext(ctx, LogIDKey, logID)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

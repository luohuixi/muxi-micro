package infrastructure

import (
	"context"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"google.golang.org/grpc"
)

func (g *Grpc) ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			str := g.Addr + " 发生错误"
			g.logger.Error(str, logger.Error(err))
		}
		return resp, err
	}
}

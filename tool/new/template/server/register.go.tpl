package server

import (
	pb "helloworld/api/local/v1"
	"helloworld/internal/infrastructure"
	"helloworld/internal/service"
	"google.golang.org/grpc"
)

func NewGRPCServer(helloService *service.HelloService, g *infrastructure.Grpc) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(g.ServerInterceptor()),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterHelloServiceServer(s, helloService)
	return s
}

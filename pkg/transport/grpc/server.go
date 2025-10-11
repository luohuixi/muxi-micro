package grpc

import (
	"context"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/registry"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Option func(*GRPCServer)

func WithRegistrationCenter(registrationCenter registry.RegistrationCenter) Option {
	return func(s *GRPCServer) {
		s.registrationCenter = registrationCenter
	}
}

func NewGRPCServer(grpcServer *grpc.Server, opts ...Option) *GRPCServer {
	s := GRPCServer{
		grpcServer: grpcServer,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func (s *GRPCServer) Serve(ctx context.Context, serviceName, host, port string) error {

	// 注册服务到注册中心,如果有的话
	if s.registrationCenter != nil {
		err := s.registrationCenter.Register(ctx, serviceName, host, port)
		if err != nil {
			return err
		}
	}

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server started on :%s", port)
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

type GRPCServer struct {
	registrationCenter registry.RegistrationCenter
	grpcServer         *grpc.Server
}

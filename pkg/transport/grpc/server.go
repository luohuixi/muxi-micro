package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"github.com/muxi-Infra/muxi-micro/pkg/tracer"
	grpclog "github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/log"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/registry"
	"google.golang.org/grpc"
)

type Option func(*GRPCServer)

type GRPCServer struct {
	registrationCenter registry.RegistrationCenter
	grpcServer         *grpc.Server
	name               string
	host               string
	port               string
	timeout            time.Duration
	l                  logger.Logger
	interceptors       []grpc.UnaryServerInterceptor
}

// WithName 设置服务名, 用于服务注册
func WithName(name string) Option {
	return func(s *GRPCServer) {
		s.name = name
	}
}

// WithHost 设置服务主机地址, 用于服务注册，默认会使用本机内网IP
func WithHost(host string) Option {
	return func(s *GRPCServer) {
		s.host = host
	}
}

func WithPort(port string) Option {
	return func(s *GRPCServer) {
		s.port = port
	}
}

func WithTimeout(t time.Duration) Option {
	return func(s *GRPCServer) {
		s.timeout = t
	}
}

// WithGlobalLogger 用于将全局 LogID 传递给这个 Logger
func WithGlobalLogger(l logger.Logger) Option {
	return func(s *GRPCServer) {
		s.l = l
	}
}

func WithServerTracer(t tracer.Tracer) Option {
	return func(s *GRPCServer) {
		s.interceptors = append(s.interceptors, t.ServerInterceptor())
	}
}

func WithExtraServerInterceptor(interceptor ...grpc.UnaryServerInterceptor) Option {
	return func(s *GRPCServer) {
		s.interceptors = append(s.interceptors, interceptor...)
	}
}

func WithRegistrationCenter(registrationCenter registry.RegistrationCenter) Option {
	return func(s *GRPCServer) {
		s.registrationCenter = registrationCenter
	}
}

func NewGRPCServer(opts ...Option) *GRPCServer {
	s := &GRPCServer{
		name:    DefaultName,
		port:    DefaultPort,
		timeout: DefaultTimeout,
		l:       logx.NewStdLogger(),
	}

	for _, o := range opts {
		o(s)
	}

	// 优先执行全局 logger 拦截器
	s.interceptors = append(
		[]grpc.UnaryServerInterceptor{grpclog.GlobalLoggerServerInterceptor(s.l)},
		s.interceptors...,
	)

	s.grpcServer = grpc.NewServer(
		grpc.ConnectionTimeout(s.timeout),
		grpc.ChainUnaryInterceptor(s.interceptors...),
	)

	return s
}

func (s *GRPCServer) Serve(ctx context.Context) error {
	// 注册服务到注册中心,如果有的话
	if s.registrationCenter != nil {
		if s.host == "" {
			var err error
			s.host, err = registry.GetLocalIP()
			if err != nil {
				return err
			}
		}
		err := s.registrationCenter.Register(ctx, s.name, s.host, s.port)
		if err != nil {
			return err
		}
	}

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server started on :%s", s.port)
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

type registerFunc func(*grpc.Server)

func (s *GRPCServer) ProtoRegister(fn ...registerFunc) {
	for _, f := range fn {
		f(s.grpcServer)
	}
}

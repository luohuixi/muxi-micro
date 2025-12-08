package grpc

import (
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"github.com/muxi-Infra/muxi-micro/pkg/tracer"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type Option2 func(*GRPCClient)

type GRPCClient struct {
	addr            string
	name            string
	l               logger.Logger
	discoveryCenter discovery.DiscoverCenter
	interceptors    []grpc.UnaryClientInterceptor
	conn            *grpc.ClientConn
}

func WithRetry(try uint, time time.Duration) Option2 {
	return func(c *GRPCClient) {
		interceptor := retry.UnaryClientInterceptor(
			retry.WithMax(try),
			retry.WithBackoff(retry.BackoffLinear(time)),
		)
		c.interceptors = append(c.interceptors, interceptor)
	}
}

// WithAddress 用于不需要服务发现的情况
func WithAddress(addr string) Option2 {
	return func(c *GRPCClient) {
		c.addr = addr
	}
}

// WithDiscoveryName 设置服务发现的服务名
func WithDiscoveryName(name string) Option2 {
	return func(c *GRPCClient) {
		c.name = name
	}
}

// WithClientLogger 用于记录 resolver 新增或删减的节点
func WithClientLogger(l logger.Logger) Option2 {
	return func(c *GRPCClient) {
		c.l = l
	}
}

func WithServiceDiscovery(discoveryCenter discovery.DiscoverCenter) Option2 {
	return func(s *GRPCClient) {
		s.discoveryCenter = discoveryCenter
	}
}

func WithClientTracer(t tracer.Tracer) Option2 {
	return func(c *GRPCClient) {
		c.interceptors = append(c.interceptors, t.ClientInterceptor())
	}
}

func WithExtraClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) Option2 {
	return func(c *GRPCClient) {
		c.interceptors = append(c.interceptors, interceptors...)
	}
}

// NewGRPCClient 每一个微服务，创建一个client，一个resolver
func NewGRPCClient(opts ...Option2) (*GRPCClient, error) {
	client := &GRPCClient{
		addr: DefaultHost + ":" + DefaultPort,
		name: DefaultName,
		l:    logx.NewStdLogger(),
	}

	for _, opt := range opts {
		opt(client)
	}

	if client.discoveryCenter == nil {
		conn, err := grpc.NewClient(
			client.addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(client.interceptors...),
		)
		if err != nil {
			return nil, err
		}

		client.conn = conn
		return client, nil
	}

	r := discovery.NewResolver(client.name, client.discoveryCenter, client.l)
	resolver.Register(r)
	conn, err := grpc.NewClient(
		"muxi:///"+client.name, // 没用上
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(client.interceptors...),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, err
	}
	client.conn = conn

	return client, nil
}

func (c *GRPCClient) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

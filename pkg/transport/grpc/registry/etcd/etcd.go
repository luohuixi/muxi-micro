package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegistry struct {
	client      *clientv3.Client
	logger      logger.Logger
	endpoints   []string
	dialTimeout time.Duration
	leaseTTL    int64
	namespace   string
	username    string
	password    string
}

type ServerOption func(*EtcdRegistry)

func WithUsername(username string) ServerOption {
	return func(r *EtcdRegistry) {
		r.username = username
	}
}

func WithPassword(password string) ServerOption {
	return func(r *EtcdRegistry) {
		r.password = password
	}
}

func WithEndpoints(endpoints []string) ServerOption {
	return func(r *EtcdRegistry) {
		r.endpoints = endpoints
	}
}

func WithLogger(l logger.Logger) ServerOption {
	return func(r *EtcdRegistry) {
		r.logger = l
	}
}

func WithDialTimeout(timeout time.Duration) ServerOption {
	return func(r *EtcdRegistry) {
		r.dialTimeout = timeout
	}
}

func WithLeaseTTL(ttl int64) ServerOption {
	return func(r *EtcdRegistry) {
		r.leaseTTL = ttl
	}
}

func WithNamespace(ns string) ServerOption {
	return func(r *EtcdRegistry) {
		r.namespace = ns
	}
}

// ===== 构造函数 =====
func NewEtcdRegistry(opts ...ServerOption) (*EtcdRegistry, error) {
	r := &EtcdRegistry{
		endpoints:   []string{"127.0.0.1:2379"},
		dialTimeout: 5 * time.Second,
		leaseTTL:    10,
		namespace:   "/services",
		logger:      logx.NewStdLogger(),
		username:    "",
		password:    "",
	}

	for _, opt := range opts {
		opt(r)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: r.dialTimeout,
		Username:    r.username,
		Password:    r.password,
	})
	if err != nil {
		return nil, err
	}
	r.client = cli
	return r, nil
}

// ===== 核心方法 =====
func (r *EtcdRegistry) Register(ctx context.Context, serviceName, host, port string) error {
	key := fmt.Sprintf("%s/%s/%s:%s", r.namespace, serviceName, host, port)
	val := fmt.Sprintf("%s:%s", host, port)

	lease, err := r.client.Grant(ctx, r.leaseTTL)
	if err != nil {
		return err
	}

	_, err = r.client.Put(ctx, key, val, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	ch, err := r.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ka, ok := <-ch:
				if !ok {
					r.logger.Warn(fmt.Sprintf("keepalive channel closed for %s", key))
					return
				}
				r.logger.Debug(fmt.Sprintf("lease renewed: %v", ka.ID))
			case <-ctx.Done():
				r.logger.Info(fmt.Sprintf("deregistering service: %s", key))
				return
			}
		}
	}()

	r.logger.Info(fmt.Sprintf("✅ Service registered: %s -> %s", key, val))

	return nil
}

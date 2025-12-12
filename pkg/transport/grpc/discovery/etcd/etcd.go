package etcd

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/discovery"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdDiscovery struct {
	client             *clientv3.Client
	logger             logger.Logger
	endpoints          []string
	dialTimeout        time.Duration
	namespace          string
	username           string
	password           string
	watchEventChanSize int

	sync.Mutex
}

type Option func(*EtcdDiscovery)

func WithUsername(username string) Option {
	return func(r *EtcdDiscovery) {
		r.username = username
	}
}

func WithPassword(password string) Option {
	return func(r *EtcdDiscovery) {
		r.password = password
	}
}

func WithEndpoints(endpoints []string) Option {
	return func(r *EtcdDiscovery) {
		r.endpoints = endpoints
	}
}

// WithLogger 用于记录新增或删减节点
func WithLogger(l logger.Logger) Option {
	return func(r *EtcdDiscovery) {
		r.logger = l
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(r *EtcdDiscovery) {
		r.dialTimeout = timeout
	}
}

func WithNamespace(ns string) Option {
	return func(r *EtcdDiscovery) {
		r.namespace = ns
	}
}

// WithWatchEventChanSize 设置监听etcd变更的channel大小, 默认为10
func WithWatchEventChanSize(size int) Option {
	return func(r *EtcdDiscovery) {
		r.watchEventChanSize = size
	}
}

func NewEtcdDiscovery(opts ...Option) (*EtcdDiscovery, error) {
	r := &EtcdDiscovery{
		endpoints:          []string{"127.0.0.1:2379"},
		dialTimeout:        5 * time.Second,
		namespace:          "/services",
		logger:             logx.NewStdLogger(),
		username:           "",
		password:           "",
		watchEventChanSize: 10,
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

func (r *EtcdDiscovery) Discover(ctx context.Context, serviceName string) ([]string, error) {
	key := fmt.Sprintf("%s/%s/", r.namespace, serviceName)

	val, err := r.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var nodes []string
	for _, kv := range val.Kvs {
		nodes = append(nodes, string(kv.Value))
	}

	return nodes, nil
}

// Watch 动态监听etcd变更
func (r *EtcdDiscovery) Watch(ctx context.Context, serviceName string) <-chan *discovery.Event {
	eventCh := make(chan *discovery.Event, r.watchEventChanSize)

	go func() {
		defer func() {
			close(eventCh)
		}()

		key := fmt.Sprintf("%s/%s/", r.namespace, serviceName)
		ch := r.client.Watch(ctx, key, clientv3.WithPrefix())
		for resp := range ch {
			for _, ev := range resp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					eventCh <- &discovery.Event{
						Type:    "PUT",
						Address: string(ev.Kv.Value),
					}
					r.logger.Warn(fmt.Sprintf("service[%s] has a new address: %s add into etcd", serviceName, string(ev.Kv.Value)))
				case clientv3.EventTypeDelete:
					addr := strings.TrimPrefix(string(ev.Kv.Key), key)
					eventCh <- &discovery.Event{
						Type:    "DELETE",
						Address: addr,
					}
					r.logger.Warn(fmt.Sprintf("(service[%s], addr[%s]) has been removed from etcd", serviceName, addr))
				}
			}
		}
	}()

	return eventCh
}

package discovery

import (
	"context"
	"fmt"
	"sync"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"google.golang.org/grpc/resolver"
)

// Resolver 实现了 Resolver 和 Builder 接口
type Resolver struct {
	serviceName string
	discovery   DiscoverCenter
	cc          resolver.ClientConn
	servers     []resolver.Address
	cancel      context.CancelFunc
	l           logger.Logger

	sync.Mutex
}

func NewResolver(serviceName string, discovery DiscoverCenter, l logger.Logger) *Resolver {
	return &Resolver{
		serviceName: serviceName,
		discovery:   discovery,
		l:           l,
	}
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	servers, err := r.discovery.Discover(ctx, r.serviceName)
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no available instance for service[%s]", r.serviceName)
	}

	r.Lock()
	r.servers = toAddrList(servers)
	r.Unlock()

	err = r.cc.UpdateState(resolver.State{Addresses: r.servers})
	if err != nil {
		return nil, err
	}

	go r.watch()

	return r, nil
}

func (r *Resolver) Scheme() string {
	return "muxi"
}

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *Resolver) Close() {
	r.cancel()
}

func (r *Resolver) watch() {
	ch := r.discovery.Watch(context.Background(), r.serviceName)
	for event := range ch {
		r.Lock()
		switch event.Type {
		case "PUT":
			r.servers = append(r.servers, resolver.Address{Addr: event.Address})
		case "DELETE":
			var updated []resolver.Address
			for _, addr := range r.servers {
				if addr.Addr != event.Address {
					updated = append(updated, addr)
				}
			}
			// 如果全部删光了, grpc会报错, updated是空的话是否要更新r.servers？
			if len(updated) == 0 {
				r.l.Error(fmt.Sprintf("all addr of service[%s] have been removed", r.serviceName))
			}
			r.servers = updated
		}
		r.Unlock()
		err := r.cc.UpdateState(resolver.State{Addresses: r.servers})
		if err != nil {
			r.l.Error(fmt.Sprintf("resolver update addr failed (service[%s])", r.serviceName), logger.Field{"err": err})
		}
	}
}

func toAddrList(nodes []string) []resolver.Address {
	result := make([]resolver.Address, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, resolver.Address{Addr: n})
	}
	return result
}

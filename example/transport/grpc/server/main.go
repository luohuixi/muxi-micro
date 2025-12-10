package server

import (
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/registry/etcd"
)

// 和proto一起生成的服务端模板的大致形式
func HelloServer() (*grpc.GRPCServer, error) {
	center, err := etcd.NewEtcdRegistry()
	if err != nil {
		return nil, err
	}

	server := grpc.NewGRPCServer(
		grpc.WithName("test"),
		grpc.WithRegistrationCenter(center),
	)

	return server, nil
}

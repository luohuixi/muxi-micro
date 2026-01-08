{{- define "client" -}}

package client

import (
	"YourPath/{{.Path}}"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/discovery/etcd"
)

func {{.ServiceName}}Client() ({{.Pkg}}.{{.ServiceName}}Client, func(), error) {
	center, err := etcd.NewEtcdDiscovery()
	if err != nil {
		return nil, nil, err
	}

	client, err := grpc.NewGRPCClient(
		grpc.WithServiceDiscovery(center),
		grpc.WithDiscoveryName("muxi-micro-server"),
	)
	if err != nil {
		return nil, nil, err
	}

	c := {{.Pkg}}.New{{.ServiceName}}Client(client.Conn())
	return c, func() { _ = client.Close() }, nil
}

{{- end -}}

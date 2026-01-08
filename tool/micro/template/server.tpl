{{- define "server" -}}

package server

import (
	"YourPath/{{.Path}}"
	mgrpc "github.com/muxi-Infra/muxi-micro/pkg/transport/grpc"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/grpc/registry/etcd"
	"google.golang.org/grpc"
)

func {{.ServiceName}}Server(srv {{.Pkg}}.{{.ServiceName}}Server) (*mgrpc.GRPCServer, error) {
	center, err := etcd.NewEtcdRegistry()
	if err != nil {
		return nil, err
	}

	server := mgrpc.NewGRPCServer(
		mgrpc.WithName("muxi-micro-server"),
		mgrpc.WithRegistrationCenter(center),
	)

	server.ProtoRegister(func(s *grpc.Server) {
		{{.Pkg}}.Register{{.ServiceName}}Server(s, srv)
	})

	return server, nil
}

{{- end -}}

//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"helloworld/internal/infrastructure"
	"helloworld/internal/repository"
	"helloworld/internal/server"
	"helloworld/internal/service"
	"google.golang.org/grpc"
)

func WireApp() (*grpc.Server, func(), error) {
	panic(wire.Build(
		infrastructure.ProvideGrpcInstance,
		infrastructure.ProvideDB,
		infrastructure.ProvideCache,
		repository.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
	))
}

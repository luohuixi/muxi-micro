package main

import (
	"helloworld/internal/infrastructure"
	"helloworld/internal/wire"
	"log"
	"net"
)

func main() {
	grpcInstance, err := infrastructure.ProvideGrpcInstance()
	if err != nil {
		log.Fatal(err)
	}
	s, _, err := wire.WireApp()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen(grpcInstance.Net, grpcInstance.Addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

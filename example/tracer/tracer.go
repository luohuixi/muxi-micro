package main

import (
	"log"

	"github.com/muxi-Infra/muxi-micro/pkg/tracer"
)

func Zipkin() {
	config, err := tracer.NewZipkin("http://localhost:9411/api/v2/spans", "demo_service", "localhost:50052", 1)
	if err != nil {
		log.Fatal(err)
	}

	// 客户端
	//conn, err := grpc.Dial(
	//	"localhost:50051",
	//	grpc.WithInsecure(), //禁用TLS
	//	grpc.WithUnaryInterceptor(
	//		grpcopentracing.UnaryClientInterceptor(),
	//	),
	//)

	// 服务端
	//s := grpc.NewServer(
	//	grpc_middleware.WithUnaryServerChain(
	//		config.ZipkinGrpc(),
	//	),
	//)

	if err := config.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Zipkin()
}

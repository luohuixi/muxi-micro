package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"net"
)

// 简易的 HelloWorld 实现
type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func GrpcServer() {
	// Zipkin
	config, err := tracer.NewZipkin(
		"http://localhost:9411/api/v2/spans",
		"demo_service",
		"localhost:50051",
		1,
	)

	// Jaeger
	//config , err := tracer.NewJaeger(
	//	"http://localhost:14268/api/traces",
	//	"demo_service",
	//	1,
	//	)

	// SkyWalking
	//config , err := tracer.NewSkyWalking(
	//	"localhost:11800",
	//	"demo_service",
	//	"demo_instance", //实例名
	//	1
	//	)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := config.Close(); err != nil {
			log.Println(err)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			config.ServerInterceptor(),
		),
	)

	helloworld.RegisterGreeterServer(s, &server{})

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func GrpcClient() {
	// 同上
	config, err := tracer.NewZipkin(
		"http://localhost:9411/api/v2/spans",
		"demo_client",
		"localhost:50051",
		1,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := config.Close(); err != nil {
			log.Println(err)
		}
	}()

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			config.ClientInterceptor(),
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := helloworld.NewGreeterClient(conn)

	resp, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: "MuXi"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", resp.Message)
}

func GinService() {
	// 同上
	config, err := tracer.NewZipkin(
		"http://localhost:9411/api/v2/spans",
		"demo_gin",
		"localhost:8081",
		1,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := config.Close(); err != nil {
			log.Println(err)
		}
	}()

	r := gin.Default()
	r.Use(config.GinMiddleware(r))

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello world"})
	})
	r.Run("0.0.0.0:8081")
}

func main() {
	// grpc服务端
	//GrpcServer()
	// grpc客户端
	GrpcClient()
	// gin
	//GinService()
}

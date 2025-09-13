package service

import (
	"context"
	pb "helloworld/api/local/v1"
	"helloworld/internal/repository"
)

type HelloService struct {
	pb.UnimplementedHelloServiceServer
	helloRepo repository.UserModels
}

func NewHelloService(helloRepo repository.UserModels) *HelloService {
	return &HelloService{
		helloRepo: helloRepo,
	}
}

func (s *HelloService) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello " + req.Username,
		Code:    200,
	}, nil
}

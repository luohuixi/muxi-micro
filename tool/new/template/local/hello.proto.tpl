syntax = "proto3";

option go_package = "./v1";

message HelloRequest {
  string username = 1;
}

message HelloResponse {
  int64 code = 1;
  string message = 2;
}

service HelloService {
  rpc Hello(HelloRequest) returns (HelloResponse);
}
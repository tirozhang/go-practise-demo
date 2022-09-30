package main

import (
	"context"
	"log"
	"net"

	hello "github.com/tirozhang/go-practise-demo/grpc/hello/gen/v1"
	"google.golang.org/grpc"
)

type Server struct {
	hello.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, &Server{})
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

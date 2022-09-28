package main

import (
	"context"
	hello "github.com/tirozhang/go-practise-demo/grpc/hello/gen/v1"
	"google.golang.org/grpc"
	"log"
)

func main() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer dial.Close()

	c := hello.NewGreeterClient(dial)
	sayHello, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "tiro"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	println(sayHello.Message)

}

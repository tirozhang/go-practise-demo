package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	stream "github.com/tirozhang/go-practise-demo/grpc/stream_grpc_test/gen/v1"

	"google.golang.org/grpc"
)

func main() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(dial *grpc.ClientConn) {
		err := dial.Close()
		if err != nil {
			log.Fatalf("close dial failed: %v", err)
		}
	}(dial)

	c := stream.NewGreeterClient(dial)
	res, err := c.GetStream(context.Background(), &stream.StreamReqData{Data: "tiro"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		res, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Println(res.Data)
	}

	putStream, _ := c.PutStream(context.Background())
	i := 0
	for {
		i++
		err := putStream.Send(&stream.StreamReqData{
			Data: fmt.Sprintf("hello %d", i),
		})
		if err != nil {
			log.Fatalf("send failed: %v", err)
		}
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}

	allStream, _ := c.AllStream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			recv, err := allStream.Recv()
			if err != nil {
				log.Fatalf("recv failed: %v", err)
			}
			log.Println(recv.Data)
		}
	}()
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			err := allStream.Send(&stream.StreamReqData{
				Data: fmt.Sprintf("hello %d", i),
			})
			if err != nil {
				log.Fatalf("send failed: %v", err)
			}
			time.Sleep(time.Second)
			if i > 10 {
				break
			}
		}
	}()
	wg.Wait()
}

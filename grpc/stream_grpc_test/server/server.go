package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	stream "github.com/tirozhang/go-practise-demo/grpc/stream_grpc_test/gen/v1"

	"google.golang.org/grpc"
)

const PORT = ":50051"

type Server struct {
	stream.UnimplementedGreeterServer
}

func (s *Server) GetStream(in *stream.StreamReqData, res stream.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		err := res.Send(&stream.StreamResData{
			Data: fmt.Sprintf("hello %d", i),
		})
		time.Sleep(time.Second)
		if err != nil {
			return err
		}
		if i > 10 {
			break
		}
	}
	return nil
}
func (s *Server) PutStream(cliStr stream.Greeter_PutStreamServer) error {
	for {
		recv, err := cliStr.Recv()
		if err != nil {
			fmt.Printf("recv failed: %v", err)
			break
		}
		log.Println(recv.Data)
	}
	return nil
}

func (s *Server) AllStream(allStr stream.Greeter_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			recv, err := allStr.Recv()
			if err != nil {
				fmt.Printf("recv failed: %v", err)
				break
			}
			log.Println(recv.Data)
		}
	}()

	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			err := allStr.Send(&stream.StreamResData{
				Data: fmt.Sprintf("hello %d", i),
			})
			time.Sleep(time.Second)
			if err != nil {
				fmt.Printf("send failed: %v", err)
				break
			}
			if i > 10 {
				break
			}
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	stream.RegisterGreeterServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

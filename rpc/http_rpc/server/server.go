package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

func main() {
	// 1. 注册处理逻辑 handler
	err := rpc.RegisterName("HelloService", new(HelloService))
	if err != nil {
		panic(err)
	}
	// 2.实例化一个server
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.Reader
			io.Closer
		}{
			Reader: r.Body,
			Writer: w,
		}
		err := rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
		if err != nil {
			panic(err)
		}
	})

	// 3. 启动服务
	err = http.ListenAndServe(":1234", nil)
	if err != nil {
		panic(err)
	}
}

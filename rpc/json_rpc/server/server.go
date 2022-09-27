package main

import (
	"net"
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

	// 1.实例化一个server
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	// 2. 注册处理逻辑 handler
	err = rpc.RegisterName("HelloService", new(HelloService))
	if err != nil {
		panic(err)
	}
	// 3. 启动服务
	for {
		conn, err := listener.Accept() //当一个新的连接进来的时候，会返回一个conn
		if err != nil {
			panic(err)
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}

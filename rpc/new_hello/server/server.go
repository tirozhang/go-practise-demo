package main

import (
	"github.com/tirozhang/go-practise-demo/rpc/new_hello/hanlder"
	"github.com/tirozhang/go-practise-demo/rpc/new_hello/server_proxy"
	"net"
	"net/rpc"
)

func main() {

	// 1.实例化一个server
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	// 2. 注册处理逻辑 handler
	err = server_proxy.RegisterHelloService(new(hanlder.HelloService))
	if err != nil {
		panic(err)
	}
	// 3. 启动服务
	for {
		conn, err := listener.Accept() //当一个新的连接进来的时候，会返回一个conn
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(conn)
	}

}

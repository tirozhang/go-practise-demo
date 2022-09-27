package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// 1. 建立连接
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}

	// 2. 调用远程方法
	var reply string
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	err = client.Call("HelloService.Hello", "hello1234", &reply)
	if err != nil {
		panic(err)
	}
	println(reply)
}

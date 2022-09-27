package main

import "net/rpc"

func main() {
	// 1. 建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}

	// 2. 调用远程方法
	var reply string
	err = client.Call("HelloService.Hello", "hello1234", &reply)
	if err != nil {
		panic(err)
	}
	println(reply)
}

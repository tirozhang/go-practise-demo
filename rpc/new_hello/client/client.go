package main

import (
	"github.com/tirozhang/go-practise-demo/rpc/new_hello/client_proxy"
)

func main() {
	var reply string

	c := client_proxy.NewHelloServiceClient()
	err := c.Hello("hello1234", &reply)
	if err != nil {
		panic(err)
	}
	println(reply)
}

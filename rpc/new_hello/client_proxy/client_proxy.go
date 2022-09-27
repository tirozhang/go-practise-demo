package client_proxy

import (
	"github.com/tirozhang/go-practise-demo/rpc/new_hello/hanlder"
	"net/rpc"
)

type HelloServiceStub struct {
	*rpc.Client
}

func NewHelloServiceClient() *HelloServiceStub {
	conn, _ := rpc.Dial("tcp", "localhost:1234")

	return &HelloServiceStub{Client: conn}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {
	err := c.Call(hanlder.HelloServiceName+".Hello", request, reply)

	if err != nil {
		return err
	}
	return nil
}

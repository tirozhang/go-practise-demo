package server_proxy

import (
	"github.com/tirozhang/go-practise-demo/rpc/new_hello/hanlder"
	"net/rpc"
)

type HelloService interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(srv HelloService) error {
	return rpc.RegisterName(hanlder.HelloServiceName, srv)
}

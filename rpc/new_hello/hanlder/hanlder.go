package hanlder

import "fmt"

const HelloServiceName = "handler.HelloService"

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	fmt.Println("request:", request)
	*reply = "hello:" + request
	return nil
}

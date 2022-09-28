package main

import (
	"encoding/json"
	"fmt"
	hello "github.com/tirozhang/go-practise-demo/protobuf/hello"
	"google.golang.org/protobuf/proto"
)

// protoc -I ./hello hello.proto --go_out=:./hello 无插件
// protoc -I ./hello hello.proto --go_out=plugins=grpc:./hello

func main() {
	req := hello.HelloRequest{
		Name:    "test",
		Age:     18,
		Courses: []string{"math", "english"},
	}
	marshal, err := proto.Marshal(&req)
	if err != nil {
		return
	}
	fmt.Println(marshal)
	fmt.Printf("%X\n", marshal)
	fmt.Println(string(marshal))

	b, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}

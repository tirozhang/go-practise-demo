package main

import (
	"encoding/json"
	"fmt"

	hello "github.com/tirozhang/go-practise-demo/protobuf/hello"
	"google.golang.org/protobuf/proto"
)

// github.com/golang/protobuf 已经准备作废
// protoc -I ./hello hello.proto --go_out=:./hello 无插件
// protoc -I ./hello hello.proto --go_out=plugins=grpc:./hello

// go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
// protoc -I . hello.proto --go_out=.
// protoc -I . hello.proto --go_out=paths=source_relative:gen/v1 生成pb
// protoc --go_out=.  --go-grpc_out=.   hello.proto 生成pb+grpc
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto

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

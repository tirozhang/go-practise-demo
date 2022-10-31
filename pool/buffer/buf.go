package main

import (
	"bytes"
	"fmt"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}

func main() {
	buf := GetBuffer()
	buf.WriteString("hello")
	println(buf.String())
	PutBuffer(buf)
	fmt.Println(buf.String())
}

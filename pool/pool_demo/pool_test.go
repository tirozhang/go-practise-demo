package main

import "testing"

// go  test -bench . -benchmem

var buf = []byte(`{"Name":"zhangsan","Age":18,"Remark":"hello"}`)

func BenchmarkStepOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stepOne(buf)
	}
}

func BenchmarkStepTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stepTwo(buf)
	}
}

package main

import (
	"encoding/json"
	"sync"
)

type Student struct {
	Name   string
	Age    int
	Remark [1024]byte
}

var studentPool = sync.Pool{
	New: func() interface{} {
		return new(Student)
	},
}

func main() {
	student := Student{
		Name: "zhangsan",
		Age:  18,
	}
	studentJson, _ := json.Marshal(student)
	stepOne(studentJson)

}

func stepOne(buf []byte) {
	stu := &Student{}
	_ = json.Unmarshal(buf, &stu)
}

func stepTwo(buf []byte) {
	stu := studentPool.Get().(*Student)
	json.Unmarshal(buf, stu)
	studentPool.Put(stu)
}

//参考文章
//https://geektutu.com/post/hpg-sync-pool.html
//https://geektutu.com/post/hpg-benchmark.html

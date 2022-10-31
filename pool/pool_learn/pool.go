package main

import "sync"

type Student struct {
	Name string
	Age  int
}

func main() {
	var studentPool = sync.Pool{
		New: func() interface{} {
			return new(Student)
		},
	}
	stu := studentPool.Get().(*Student)
	println(stu.Name, stu.Age)
	stu.Name = "zhangsan"
	stu.Age = 18
	studentPool.Put(stu)

	stu2 := studentPool.Get().(*Student)
	println(stu2.Name, stu2.Age)

	var studentPoolNil = sync.Pool{} // New is nil
	stu3 := studentPoolNil.Get()     // 返回nil
	println(stu3)

}

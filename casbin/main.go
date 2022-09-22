package main

// https://darjun.github.io/2020/06/12/godailylib/casbin/

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
)

type Object struct {
	Name  string
	Owner string
}

type Subject struct {
	Name string
	Hour int
}

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	fmt.Printf("%s, %s, %s : %v\n", sub, act, obj, ok)
}

func checkDomain(e *casbin.Enforcer, sub, domain, obj, act string) {
	ok, _ := e.Enforce(sub, domain, obj, act)
	fmt.Printf("%s, %s,%s, %s : %v\n", sub, domain, act, obj, ok)
}

func checkABAC(e *casbin.Enforcer, sub Subject, obj Object, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s at %d:00\n", sub.Name, act, obj.Name, sub.Hour)
	} else {
		fmt.Printf("%s CANNOT %s %s at %d:00\n", sub.Name, act, obj.Name, sub.Hour)
	}
}

func demo(sub, obj, act string) {
	// 1.load 模型和策略
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	// 2. 判断权限
	ok, _ := e.Enforce(sub, obj, act)
	fmt.Printf("%s, %s, %s : %v\n", sub, act, obj, ok)
}

func main() {

	// 1.ACL model
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	//check(e, "doctorA", "data1", "read")
	//check(e, "doctorA", "data1", "write")
	//check(e, "doctorB", "data2", "read")
	//check(e, "doctorB", "data2", "write")

	// 2.RBAC model
	e, err = casbin.NewEnforcer("./model.rbac.conf", "./policy.rbac.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	//check(e, "doctorA", "data", "read")
	//check(e, "doctorA", "data", "write")
	//check(e, "assistA", "data", "read")
	//check(e, "assistA", "data", "write")

	// 3.RBAC1 model
	e, err = casbin.NewEnforcer("./model.rbac.conf", "./policy.rbac.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	//check(e, "doctorA", "data2", "read")
	//check(e, "doctorB", "data2", "read")

	// 4.RBAC domain
	e, err = casbin.NewEnforcer("./model.rbac.domain.conf", "./policy.rbac.domain.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	//checkDomain(e, "user_01", "team1", "drafts", "create")
	//checkDomain(e, "user_02", "team1", "drafts", "create")
	//checkDomain(e, "user_01", "team2", "drafts", "create")

	// 5.ABAC
	e, err = casbin.NewEnforcer("./model.abac.conf", "./policy.abac.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	o := Object{"data", "doctor"}
	s1 := Subject{"doctor", 10}

	checkABAC(e, s1, o, "read")

	s2 := Subject{"assist", 10}
	checkABAC(e, s2, o, "read")

	s3 := Subject{"doctor", 20}
	checkABAC(e, s3, o, "read")

	s4 := Subject{"assist", 20}
	checkABAC(e, s4, o, "read")

}

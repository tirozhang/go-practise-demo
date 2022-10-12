package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/tirozhang/go-practise-demo/net-http/framework/middleware"

	"github.com/tirozhang/go-practise-demo/net-http/controller"

	"github.com/tirozhang/go-practise-demo/net-http/framework"

	"github.com/golang/glog"
)

func main() {
	//var fooHandler *FooHandler = &FooHandler{}
	//fmt.Println(new(FooHandler))
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	defer glog.Flush()

	go NewCoreHttp()
	go NewHandleHttp()
	NewFileServer()
}

func NewCoreHttp() {
	glog.Infoln("new core http")
	core := framework.NewCore()
	core.Use(
		middleware.Test1(),
		middleware.Test2(),
	)
	core.Use(middleware.Timeout(1 * time.Second))
	RegisterRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    "localhost:8080",
	}
	glog.Fatal(server.ListenAndServe())
}

func NewHandleHttp() {
	// http://localhost:8081/demo
	http.Handle("/demo", new(controller.Handler))

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		if err != nil {
			glog.Errorf("write response failed, err:%v", err)
		}
	})

	http.HandleFunc("/foo", controller.NewFoo())

	glog.Fatal(http.ListenAndServe(":8081", nil))
}

func NewFileServer() {
	http.Handle("/file", http.FileServer(http.Dir("/tmp")))
	glog.Fatal(http.ListenAndServe(":8082", nil))
}

// RegisterRouter 注册规则路由
func RegisterRouter(core *framework.Core) {

	core.Get("/foo", controller.FooControllerHandler)

	// 静态路由+HTTP方法匹配
	core.Get("/user/login", framework.TimeoutHandler(controller.UserLoginController, 1*time.Second))

	//批量通用前缀
	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Test1())
		// 动态路由
		subjectApi.Delete("/:id", controller.SubjectDelController)

		subjectApi.Put("/:id", controller.SubjectUpdateController)
		subjectApi.Get("/:id", middleware.Test2(), controller.SubjectGetController)
		subjectApi.Get("/list/all", controller.SubjectListController)

		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", controller.SubjectNameController)
		}
	}
}

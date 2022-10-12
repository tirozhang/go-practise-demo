package framework

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

type Core struct {
	router      map[string]*Tree    // all routers
	middlewares []ControllerHandler // 从core这边设置的中间件

}

func NewCore() *Core {

	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

// Use 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Get Method=GET
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	handlers = c.combineHandlers(handlers...)
	if err := c.router["GET"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
	PrintTree(0, c.router["GET"].root)
	//log.Println("add router success", url)
}

// Post Method=POST
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	handlers = c.combineHandlers(handlers...)
	if err := c.router["POST"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Put Method=Put
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	handlers = c.combineHandlers(handlers...)
	if err := c.router["PUT"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Delete Method=Delete
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	handlers = c.combineHandlers(handlers...)
	if err := c.router["DELETE"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("core.serveHTTP")

	ctx := NewContext(w, r)

	handlers := c.FindRouteByRequest(r)

	if handlers == nil {
		ctx.Json(http.StatusNotFound, "404 not found")
		return
	}
	glog.Infoln("core.router")

	ctx.SetHandlers(handlers)

	if err := ctx.Next(); err != nil {
		ctx.Json(http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) FindRouteByRequest(r *http.Request) []ControllerHandler {
	glog.Infoln("core.findRouteByRequest")

	method := strings.ToUpper(r.Method)
	url := strings.ToUpper(r.URL.Path)

	glog.Infoln("core.findRouteByRequest", method, url)

	router := c.router[method]

	if router == nil {
		return nil
	}

	handler := router.FindHandler(url)

	if handler == nil {
		return nil
	}

	return handler
}

func (c *Core) combineHandlers(handlers ...ControllerHandler) []ControllerHandler {
	finalSize := len(c.middlewares) + len(handlers)
	mergedHandlers := make([]ControllerHandler, finalSize)
	copy(mergedHandlers, c.middlewares)
	copy(mergedHandlers[len(c.middlewares):], handlers)
	return mergedHandlers
}

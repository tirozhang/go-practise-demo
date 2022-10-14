package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/tirozhang/go-practise-demo/net-http/framework"

	"github.com/golang/glog"
)

func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), 1*time.Second)
	defer cancel()

	obj := map[string]interface{}{
		"data": nil,
	}

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// panic(errors.New("panic"))
		// time.Sleep(10 * time.Second)

		foo, _ := c.QueryString("foo", "empty")
		obj["data"] = foo
		c.SetStatus(http.StatusOK).Json(obj)
		finish <- struct{}{}
	}()
	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		glog.Errorf("%v", p)
		c.SetStatus(http.StatusInternalServerError).Json(p)
	case <-finish:
		glog.Infoln("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.SetStatus(http.StatusInternalServerError).Json("time out")
		c.SetHasTimeout()
	}

	return nil
}

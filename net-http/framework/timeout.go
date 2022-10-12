package framework

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/glog"
)

func TimeoutHandler(handler ControllerHandler, d time.Duration) ControllerHandler {
	return func(c *Context) error {
		finish := make(chan error, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			finish <- handler(c)
		}()
		select {
		case p := <-panicChan:
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			glog.Errorf("%v", p)
			c.Json(http.StatusInternalServerError, p)
		case err := <-finish:
			return err
		case <-durationCtx.Done():
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			c.Json(http.StatusInternalServerError, "time out")
			c.SetHasTimeout()
		}
		return nil
	}
}

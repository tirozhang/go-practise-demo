package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tirozhang/go-practise-demo/net-http/framework"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
		finish := make(chan error, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			finish <- c.Next()
		}()

		// 执行业务逻辑后操作
		select {
		case p := <-panicChan:
			c.SetStatus(500).Json("time out")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.SetStatus(500).Json("time out")
		}

		return nil
	}
}

package controller

import "github.com/tirozhang/go-practise-demo/net-http/framework"

func UserLoginController(c *framework.Context) error {
	c.Json(200, "ok, UserLoginController")
	return nil
}

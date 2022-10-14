package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
)

type Person struct {
	Age  int    `uri:"age" binding:"required,gt=10"` // age必须大于10
	Name string `uri:"name" binding:"required"`      // name必须
}

func main() {

	router := gin.Default() //	default Use(Logger(), Recovery()) // 从engine这边设置的中间件

	router.Use(gin.Logger(), gin.Recovery()) // 从router这边设置的中间件

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)

	routerGroup := router.Group("/v1")
	{
		routerGroup.Use(MyLogger())
		routerGroup.GET("/someGet", getting)
		routerGroup.POST("/somePost", posting)
		routerGroup.GET("/goods/:id/:action", goods)
		routerGroup.GET("/user/:name/:age", user)

	}

	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // 是 c.Request.URL.Query().Get("lastname") 的简写

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	router.GET("/moreJSON", func(c *gin.Context) {
		// 你也可以使用一个结构体
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}

		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123

		// 注意 msg.Name 在 JSON 中变成了 "user"，因为在 msg 结构体中定义了这个字段的 tag
		c.JSON(http.StatusOK, msg)
	})

	router.GET("/moreProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"

		// 你也可以使用 protobuf 数据格式
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		// 注意这里我们使用 `ProtoBuf` 作为渲染器
		c.ProtoBuf(http.StatusOK, data)
	})

	router.POST("/loginJSON", func(c *gin.Context) {
		var login LoginFrom
		if err := c.ShouldBind(&login); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if login.User != "manu" || login.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	router.POST("/signup", func(c *gin.Context) {
		var json SignUpFrom
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.Name != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func user(context *gin.Context) {
	var person Person
	if err := context.ShouldBindUri(&person); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"name": person.Name,
		"age":  person.Age,
	})
}

func goods(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"id":     context.Param("id"),
		"action": context.Param("action"),
	})
}

func posting(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"method": "POST",
	})
}

func getting(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"method": "GET",
	})
}

type LoginFrom struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required,min=3,max=10"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// SignUpFrom https://github.com/go-playground/validator/tree/master
type SignUpFrom struct {
	Age        int    `form:"age" json:"age" xml:"age" binding:"required,gt=10"`
	Name       string `form:"name" json:"name" xml:"name" binding:"required,min=3,max=10"`
	Email      string `form:"email" json:"email" xml:"email" binding:"required,email"`
	Password   string `form:"password" json:"password" xml:"password" binding:"required"`
	RePassword string `form:"repassword" json:"repassword" xml:"repassword" binding:"required,eqfield=Password"`
}

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 请求前

		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Print(latency)

		// 获取状态
		status := c.Writer.Status()
		log.Println(status)
	}
}

func TokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		// 从 header 中取出 token
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}
		// 验证 token
		if token != "123" {
			//  token 错误
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

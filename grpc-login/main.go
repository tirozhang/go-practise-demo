package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tirozhang/go-practise-demo/login/token"

	"github.com/tirozhang/go-practise-demo/login/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/tirozhang/go-practise-demo/login/auth"
	"github.com/tirozhang/go-practise-demo/login/wechat"

	authpb "github.com/tirozhang/go-practise-demo/login/api/gen/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	logger, err := newZapLogger()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	initDB()

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "wx4f4bc4dec97d474b",
			AppSecret: "70b4f7f8b1f0b4e6a0b4f7f8b1f0b4e6",
		},
		Logger:         logger,
		TokenGenerator: token.NewJWTTokenGen("auth server", "123456"),
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen: %v", zap.Error(err))
	}
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", zap.Error(err))
	}
}

func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	return cfg.Build()
}

func initDB() {

	var c = config.MysqlConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Name:     "db_auth",
		User:     "root",
		Password: "12345678",
	}

	// 初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)

	// 全局模式
	var err error
	config.DbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return
}

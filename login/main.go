package main

import (
	"log"
	"net"

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
	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "wx4f4bc4dec97d474b",
			AppSecret: "70b4f7f8b1f0b4e6a0b4f7f8b1f0b4e6",
		},
		Logger: logger,
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

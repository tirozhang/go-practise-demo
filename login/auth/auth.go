package auth

import (
	"context"

	authpb "github.com/tirozhang/go-practise-demo/login/api/gen/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	OpenIDResolver OpenIDResolver
	Logger         *zap.Logger
	authpb.UnimplementedAuthServiceServer
}

func (s *Service) Login(ctx context.Context, in *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("login received code", zap.String("code", in.Code))

	openID, err := s.OpenIDResolver.Resolve(ctx, in.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "can't resolve openID: %v", err)
	}
	s.Logger.Info("login resolved openID", zap.String("openID", openID))
	return &authpb.LoginResponse{}, nil
}

type OpenIDResolver interface {
	Resolve(ctx context.Context, code string) (string, error)
}

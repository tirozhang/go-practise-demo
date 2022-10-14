package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/tirozhang/go-practise-demo/login/model"

	authpb "github.com/tirozhang/go-practise-demo/login/api/gen/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	OpenIDResolver OpenIDResolver
	TokenGenerator TokenGenerator
	Logger         *zap.Logger
	authpb.UnimplementedAuthServiceServer
}

type OpenIDResolver interface {
	Resolve(ctx context.Context, code string) (string, error)
}

// TokenGenerator generates a token for the specified account.
type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

func (s *Service) Login(ctx context.Context, in *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("login received code", zap.String("code", in.Code))

	openID, err := s.OpenIDResolver.Resolve(ctx, in.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "can't resolve openID: %v", err)
	}

	userID, err := model.GetAuthInstance().ResolveUserID(ctx, openID)
	if err != nil {
		s.Logger.Error("can't resolve user id", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}
	token, err := s.TokenGenerator.GenerateToken(strconv.Itoa(int(userID)), time.Hour)
	if err != nil {
		s.Logger.Error("can't generate token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}
	return &authpb.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int32(time.Hour.Seconds()),
	}, nil
}

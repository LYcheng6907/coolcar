// 最外层auth文件夹代表auth 微服务，本层auth代表auth逻辑
package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	OpenIDResolver OpenIDResolver
	TokenGenerator TokenGenerator
	TokenExpire    time.Duration
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

// 返回openid
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// 为accountID生成 Token
type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

// 实现Login接口
func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// s.Logger.Info("received code", zap.String("code", req.Code))
	openID, err := s.OpenIDResolver.Resolve(req.Code)
	if err != nil { // 显示在客户端，显示标准化错误给用户看
		return nil, status.Errorf(codes.Unavailable,
			"cannot resolve openid: %v", err)
	}
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))

		return nil, status.Error(codes.Internal, "")
	}

	token, err := s.TokenGenerator.GenerateToken(accountID.String(), s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: token,                          //"token for open id:" + accountID,
		ExpiresIn:   int32(s.TokenExpire.Seconds()), // 过期时间
	}, nil
}

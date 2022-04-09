// 处理拦截器的请求
package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"coolcar/shared/id"

	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer"
)

type tokenVerifier interface {
	Verify(token string) (string, error)
}
type interceptor struct {
	verifier tokenVerifier
}

func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	f, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file %v", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key file %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key %v", err)
	}

	i := &interceptor{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
	}
	return i.HandleReq, nil
}

func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		// 小程序端中的header中有没有authorization
		return nil, status.Error(codes.Unauthenticated, "")
	}
	accountID, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid %d", err)
	}
	// 把accountID塞进ContextWithAccountID（公共，目的是用来知晓请求是哪个accountID发过来的），再去调用request
	return handler(ContextWithAccountID(ctx, id.AccountID(accountID)), req)
}

func tokenFromContext(c context.Context) (string, error) {
	unauthenticated := status.Error(codes.Unauthenticated, "")

	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", unauthenticated
	}

	tkn := ""

	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}

	if tkn == "" {
		return "", unauthenticated
	}

	return tkn, nil
}

type accountIDKey struct{}

// 创建context并赋值
func ContextWithAccountID(c context.Context, accountID id.AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, accountID)
}

// 获得accountid，返回 unauthenticated error if no accountID is availd
func AccountIDFromcontext(c context.Context) (id.AccountID, error) {
	v := c.Value(accountIDKey{})
	aid, ok := v.(id.AccountID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return aid, nil
}

package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 实现 GeneratorToken 是一个JWTTokenGen
type JWTTokenGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

// 构造函数，外部使用本函数来构造JWTToken
func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

// 不知道函数写的对不对，写一段单元测试
func (t *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	// nowSec := time.Now().Unix() 对测试不利，不希望现在和10分钟以后的值不同
	nowSec := t.nowFunc().Unix()
	// 此时的token相当于 header + playload,还要进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer, // 谁签名（jwt签名者）"coolcar/auth"
		IssuedAt:  nowSec,   // 转换成秒
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID, // 颁发给谁

	})

	// 签名
	return token.SignedString(t.privateKey)

}

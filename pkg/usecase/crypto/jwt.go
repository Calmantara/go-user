package crypto

import (
	"context"

	"github.com/kataras/jwt"
)

var (
	sharedKey []byte = []byte("tH1s1ss3cr3tk3yforjwT")
)

type Jwt interface {
	CreateJWT(ctx context.Context, claim any) (string, error)
	VerifyJWT(ctx context.Context, token string) (claims any)
}

type JwtKataraImpl struct{}

func NewJwt() Jwt {
	return &JwtKataraImpl{}
}

func (j *JwtKataraImpl) CreateJWT(ctx context.Context, claim any) (string, error) {
	token, err := jwt.Sign(jwt.HS256, sharedKey, claim)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func (j *JwtKataraImpl) VerifyJWT(ctx context.Context, token string) (claims any) {
	verifiedToken, err := jwt.Verify(jwt.HS256, sharedKey, []byte(token))
	if err != nil {
		panic(err)
	}

	err = verifiedToken.Claims(&claims)
	if err != nil {
		panic(err)
	}
	return claims
}

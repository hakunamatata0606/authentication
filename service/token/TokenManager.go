package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimMap map[string]interface{}

type TokenManager interface {
	CreateToken(ClaimMap, int64) (string, error)
	ParseToken(string) (ClaimMap, error)
}

type JwtTokenManager struct {
	secret []byte
}

func NewJwtTokenManager(secret string) TokenManager {
	return &JwtTokenManager{
		secret: []byte(secret),
	}
}

func (jwtManager *JwtTokenManager) CreateToken(claims ClaimMap, timeout int64) (string, error) {
	timeout = timeout + time.Now().Unix()
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}
	jwtClaims["__exp"] = timeout
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(jwtManager.secret)
}

func (jwtManager *JwtTokenManager) ParseToken(tokenStr string) (ClaimMap, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtManager.secret, nil
	})

	if err != nil {
		return nil, err
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected error could not get claims")
	}

	exp, ok := jwtClaims["__exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("unexpected error could not get exp")
	}
	if exp < (float64)(time.Now().Unix()) {
		return nil, fmt.Errorf("token timeout")
	}

	claims := ClaimMap{}
	for k, v := range jwtClaims {
		claims[k] = v
	}
	claims["token"] = tokenStr

	return claims, nil
}

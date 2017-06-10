package providers

import (
	"fmt"

	// 	"github.com/alioygur/gocart/engine"
	jwtgo "github.com/dgrijalva/jwt-go"
)

func NewJWT() engine.JWTSignParser {
	return &jwt{}
}

type jwt struct{}

func (j *jwt) Sign(claims map[string]interface{}, secret string) (string, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwt) Parse(tokenStr string, secret string) (map[string]interface{}, error) {
	token, err := jwtgo.Parse(tokenStr, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		e, ok := err.(*jwtgo.ValidationError)
		if ok {
			return nil, newTokenErr(e)
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, fmt.Errorf("can't get map claims from: %s", tokenStr)
	}
	return claims, nil
}

func newTokenErr(err *jwtgo.ValidationError) error {
	if err.Errors == jwtgo.ValidationErrorExpired {
		return engine.NewTokenErr(err.Error(), true)
	}

	return engine.NewTokenErr(err.Error(), false)
}

package service

import "github.com/golang-jwt/jwt"

type JWTToken struct {
	jwt.Token
}

func ValidateJWT(tokenStr string) (JWTToken, error) {

}

package jwt

import (
	"os"

	"github.com/golang-jwt/jwt"
)

type AccessToken struct {
	*jwt.Token
	jwt.StandardClaims
	UserId uint
}

func NewAccessToken(userId uint) *AccessToken {
	return &AccessToken{
		Token:  jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{}),
		UserId: userId,
	}
}

func (at *AccessToken) SignedString() (string, error) {
	secret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	return at.Token.SignedString(secret)
}

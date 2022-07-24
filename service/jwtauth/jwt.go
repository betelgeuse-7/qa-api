package jwtauth

import (
	"errors"
	"time"

	"github.com/betelgeuse-7/qa/config"
	"github.com/golang-jwt/jwt"
)

const (
	AT_EXPIRY = (time.Hour * 24) * 2 // 2 days
)

type TokenRepo struct {
	cfg *config.ConfigJwt
}

func NewTokenRepo(cfg *config.ConfigJwt) *TokenRepo {
	return &TokenRepo{cfg: cfg}
}

type AccessToken struct {
	*jwt.StandardClaims
	UserId int64 `json:"user_id"`
}

type tokenBuilderFn func(*TokenRepo, int64) (string, error)

func (t *TokenRepo) NewToken(userId int64, fn tokenBuilderFn) (string, error) {
	return fn(t, userId)
}

func NewAccessToken(tr *TokenRepo, userId int64) (string, error) {
	return newToken(tr, "access", userId)
}

func newToken(tr *TokenRepo, type_ string, userId int64) (string, error) {
	switch type_ {
	case "access":
		atClaims := &AccessToken{StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AT_EXPIRY).Unix(),
		}, UserId: userId}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		return t.SignedString(tr.cfg.SecretKey)
	}
	return "", errors.New("invalid token type: '" + type_ + "\n")
}

// use *AccessToken as the claims. RefreshToken, and AccessToken is the same basically
func (tr *TokenRepo) ParseToken(raw string) (*jwt.Token, *AccessToken, error) {
	claims := &AccessToken{}
	// !
	// If there's no 'exp' field in JWT payload, this function leads to nil pointer dereference
	// error, and the server crashes. Shouldn't this function return an error, in case a required
	// field is non-existent?
	// We can easily crash the server, by sending a JWT with a missing 'exp' field. This is stupid.
	tok, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
		return tr.cfg.SecretKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return tok, claims, nil
}

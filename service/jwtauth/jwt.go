package jwtauth

import (
	"errors"
	"time"

	"github.com/betelgeuse-7/qa/config"
	"github.com/golang-jwt/jwt"
)

const (
	AT_EXPIRY = time.Hour * 2
	// 3 days
	RT_EXPIRY = (time.Hour * 24) * 3
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

type RefreshToken struct {
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

func NewRefreshToken(tr *TokenRepo, userId int64) (string, error) {
	return newToken(tr, "refresh", userId)
}

func newToken(tr *TokenRepo, type_ string, userId int64) (string, error) {
	switch type_ {
	case "access":
		atClaims := &AccessToken{StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AT_EXPIRY).Unix(),
		}, UserId: userId}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		return t.SignedString(tr.cfg.SecretKey)
	case "refresh":
		rtClaims := &RefreshToken{StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RT_EXPIRY).Unix(),
		}, UserId: userId}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
		return t.SignedString(tr.cfg.SecretKey)
	}
	return "", errors.New("invalid token type: '" + type_ + "\n")
}

// use *AccessToken as the claims. RefreshToken, and AccessToken is the same basically
func (tr *TokenRepo) ParseToken(raw string) (*jwt.Token, *AccessToken, error) {
	claims := &AccessToken{}
	tok, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(tr.cfg.SecretKey), nil
	})
	if err != nil {
		return nil, nil, err
	}
	return tok, claims, nil
}

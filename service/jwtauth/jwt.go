package jwtauth

import (
	"errors"
	"time"

	"github.com/betelgeuse-7/qa/config"
	"github.com/golang-jwt/jwt"
)

const (
	_AT_EXPIRY = time.Hour * 2
	// 3 days
	_RT_EXPIRY = (time.Hour * 24) * 3
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
			ExpiresAt: _AT_EXPIRY.Milliseconds(),
		}, UserId: userId}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		return t.SignedString(tr.cfg.SecretKey)
	case "refresh":
		rtClaims := &RefreshToken{StandardClaims: &jwt.StandardClaims{
			ExpiresAt: _RT_EXPIRY.Milliseconds(),
		}, UserId: userId}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
		return t.SignedString(tr.cfg.SecretKey)
	}
	return "", errors.New("invalid token type: '" + type_ + "\n")
}

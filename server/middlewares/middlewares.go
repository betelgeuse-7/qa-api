package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type rContextKey int

const (
	JwtTokenKey rContextKey = iota
)

func JWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get Authorization header
		authHeader := r.Header.Get("Authorization")
		// check if authHeader starts with 'Bearer '
		authHeaderSplit := strings.Split(authHeader, " ")
		if authHeaderSplit[0] != "Bearer " {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized :("))
			return
		}
		// get token from authHeader

		//jwtToken := authHeaderSplit[1]

		token, err := "", errors.New("middlewares.go#30") //jwt.ValidateJWT(jwtToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), JwtTokenKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqMethod := r.Method

		if reqMethod == http.MethodPost || reqMethod == http.MethodPut || reqMethod == http.MethodPatch {

			ctHeader := r.Header.Get("Content-Type")

			if ctHeader == "" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte("a Content-Type header must be present"))
				return
			}

			if ctHeader != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte("Content-Type header's value must be set to 'application/json', and the body must contain valid JSON data"))
				return
			}

			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

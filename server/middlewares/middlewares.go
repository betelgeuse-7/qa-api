package middleware

import (
	"context"
	"net/http"
	"qa/service"
	"strings"
)

func JWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get Authorization header
		authHeader := r.Header.Get("Authorization")
		// check if authHeader starts with 'Bearer '
		authHeaderSplit := strings.Split(authHeader, " ")
		if authHeaderSplit[0] != "Bearer " {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// get token from authHeader
		jwtToken := authHeaderSplit[1]

		token, err := service.ValidateJWT(jwtToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

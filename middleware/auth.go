package middleware

import (
	"context"
	"final-project/helper"
	"final-project/helper/response"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = response.New[any](response.Authentication)

		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Warn("no token provided")
			resp.Error(helper.ErrNotLoggedIn).Code(http.StatusUnauthorized).Send(w)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			logger.Warn("invalid token format")
			resp.Error(helper.ErrNotLoggedIn).Code(http.StatusUnauthorized).Send(w)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := helper.VerifyJWT(token)
		if err != nil {
			logger.Error("failed to verify token", "error", err.Error())
			resp.Error(helper.ErrNotLoggedIn).Code(http.StatusUnauthorized).Send(w)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), helper.UserIDKey, claims["sub"]))

		next.ServeHTTP(w, r)
	})
}

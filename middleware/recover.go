package middleware

import (
	"final-project/helper"
	"final-project/helper/response"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = response.New[any](response.Default)
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "cause", r)
				resp.Error(helper.ErrInternal).Code(http.StatusInternalServerError).Send(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

package auth

import (
	"net/http"

	"go.uber.org/zap"
)

// Middleware enforces bearer auth on incoming requests.
func Middleware(validator *Validator, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := validator.Parse(r.Header.Get("Authorization"))
			if err != nil {
				log.Warn("auth failed", zap.Error(err))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

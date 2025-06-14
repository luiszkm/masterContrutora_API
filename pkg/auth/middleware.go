// file: pkg/auth/middleware.go
package auth

import (
	"context"
	"net/http"
)

func (s *JWTService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Tenta ler o cookie da requisição.
		cookie, err := r.Cookie("jwt-token")
		if err != nil {
			// Se o cookie não existir, o usuário não está autenticado.
			if err == http.ErrNoCookie {
				http.Error(w, "Token de autorização ausente", http.StatusUnauthorized)
				return
			}
			// Outro erro qualquer.
			http.Error(w, "Requisição inválida", http.StatusBadRequest)
			return
		}

		// 2. Pega o valor do token do cookie.
		tokenStr := cookie.Value

		claims, err := s.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims["sub"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

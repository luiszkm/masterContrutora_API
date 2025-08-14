// file: pkg/auth/middleware.go
package auth

import (
	"context"
	"net/http"
)

const PermissoesContextKey = contextKey("permissoes")

func (s *JWTService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		// 1. Tenta ler o cookie da requisição.
		cookie, err := r.Cookie("jwt-token")
		if err == nil {
			// Cookie encontrado
			tokenStr = cookie.Value
		} else if err == http.ErrNoCookie {
			// 2. Se cookie não existir, tenta o header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token de autorização ausente", http.StatusUnauthorized)
				return
			}
			// Verifica se o header tem o formato "Bearer <token>"
			const bearerPrefix = "Bearer "
			if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
				http.Error(w, "Formato de autorização inválido", http.StatusUnauthorized)
				return
			}
			tokenStr = authHeader[len(bearerPrefix):]
		} else {
			// Outro erro qualquer.
			http.Error(w, "Requisição inválida", http.StatusBadRequest)
			return
		}

		claims, err := s.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		// --- INÍCIO DA MUDANÇA ---
		ctx := context.WithValue(r.Context(), UserContextKey, claims["sub"])

		// Extrai as permissões (claims) e as adiciona ao contexto.
		if permissoes, ok := claims["permissoes"].([]interface{}); ok {
			permissoesStr := make([]string, len(permissoes))
			for i, v := range permissoes {
				permissoesStr[i] = v.(string)
			}
			ctx = context.WithValue(ctx, PermissoesContextKey, permissoesStr)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func Authorize(permissaoRequerida string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pega as permissões do usuário que o AuthMiddleware já colocou no contexto.
			permissoes, ok := r.Context().Value(PermissoesContextKey).([]string)
			if !ok {
				http.Error(w, "Permissões não encontradas no token", http.StatusForbidden)
				return
			}

			// Verifica se a permissão necessária está na lista de permissões do usuário.
			for _, p := range permissoes {
				if p == permissaoRequerida {
					// Permissão encontrada, pode prosseguir.
					next.ServeHTTP(w, r)
					return
				}
			}

			// Se o loop terminar e a permissão não for encontrada, o acesso é negado.
			http.Error(w, "Acesso negado: permissão insuficiente", http.StatusForbidden)
		})
	}
}

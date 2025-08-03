package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/pkg/auth"
	"github.com/luiszkm/masterCostrutora/pkg/logging"
)

// contextKey é o tipo para chaves do contexto
type contextKey string

const (
	RequestIDKey contextKey = "requestId"
	UserIDKey    contextKey = "userId"
	StartTimeKey contextKey = "startTime"
)

// ResponseWriter customizado para capturar status code e response size
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
	written      bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(data)
	w.responseSize += n
	return n, err
}

// RequestLogger middleware que registra todas as requisições
func RequestLogger(logger *logging.AppLogger, jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Gerar Request ID único
			requestID := uuid.New().String()

			// Adicionar informações ao contexto
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			ctx = context.WithValue(ctx, StartTimeKey, time.Now())

			// Extrair User ID do JWT se disponível
			userID := extractUserIDFromRequest(r, jwtService)
			if userID != "" {
				ctx = context.WithValue(ctx, UserIDKey, userID)
			}

			// Wrapper para capturar response
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Adicionar headers de rastreamento
			rw.Header().Set("X-Request-ID", requestID)

			// Executar próximo handler
			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)

			// Calcular duração
			startTime := ctx.Value(StartTimeKey).(time.Time)
			duration := time.Since(startTime)

			// Registrar requisição
			extra := map[string]interface{}{
				"requestId":    requestID,
				"userAgent":    r.UserAgent(),
				"remoteAddr":   r.RemoteAddr,
				"responseSize": rw.responseSize,
				"referer":      r.Referer(),
			}

			// Adicionar parâmetros de query se existirem
			if r.URL.RawQuery != "" {
				extra["queryParams"] = r.URL.RawQuery
			}

			logger.LogHTTPRequest(ctx, r.Method, r.URL.Path, rw.statusCode, duration, userID, extra)
		})
	}
}

// ErrorRecovery middleware que captura panics e erros não tratados
func ErrorRecovery(logger *logging.AppLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Capturar stack trace completo
					stackTrace := string(debug.Stack())

					// Criar erro estruturado
					errorMsg := fmt.Sprintf("Panic recovered: %v", err)

					extra := map[string]interface{}{
						"panic":      true,
						"method":     r.Method,
						"path":       r.URL.Path,
						"userAgent":  r.UserAgent(),
						"remoteAddr": r.RemoteAddr,
						"stackTrace": stackTrace,
					}

					// Adicionar Request ID se disponível
					if requestID := r.Context().Value(RequestIDKey); requestID != nil {
						extra["requestId"] = requestID
					}

					// Log do erro
					logger.LogError(r.Context(), errorMsg, fmt.Errorf("%v", err), extra)

					// Resposta de erro genérica
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// ErrorLogger middleware que intercepta erros HTTP de status >= 400
func ErrorLogger(logger *logging.AppLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrapper para interceptar WriteHeader
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			// Log erros HTTP (status >= 400)
			if rw.statusCode >= 400 {
				extra := map[string]interface{}{
					"method":       r.Method,
					"path":         r.URL.Path,
					"statusCode":   rw.statusCode,
					"userAgent":    r.UserAgent(),
					"remoteAddr":   r.RemoteAddr,
					"responseSize": rw.responseSize,
				}

				// Adicionar Request ID se disponível
				if requestID := r.Context().Value(RequestIDKey); requestID != nil {
					extra["requestId"] = requestID
				}

				// Adicionar User ID se disponível
				if userID := r.Context().Value(UserIDKey); userID != nil {
					extra["userId"] = userID
				}

				// Determinar se é erro do cliente ou servidor
				var level logging.LogLevel = logging.LevelWarn
				if rw.statusCode >= 500 {
					level = logging.LevelError
				}

				message := fmt.Sprintf("HTTP %d: %s %s", rw.statusCode, r.Method, r.URL.Path)

				if level == logging.LevelError {
					// Para erros 5xx, usar LogError
					logger.LogError(r.Context(), message, fmt.Errorf("HTTP %d error", rw.statusCode), extra)
				} else {
					// Para erros 4xx, usar LogAudit
					logger.LogAudit(r.Context(), message, extra)
				}
			}
		})
	}
}

// DatabaseErrorLogger captura erros de banco de dados
func DatabaseErrorLogger(logger *logging.AppLogger) func(operation string, err error, extra map[string]interface{}) {
	return func(operation string, err error, extra map[string]interface{}) {
		if err == nil {
			return
		}

		if extra == nil {
			extra = make(map[string]interface{})
		}

		extra["operation"] = operation
		extra["database"] = true

		logger.LogError(context.Background(), fmt.Sprintf("Database error in %s", operation), err, extra)
	}
}

// ServiceErrorLogger captura erros de serviços
func ServiceErrorLogger(logger *logging.AppLogger) func(ctx context.Context, service string, method string, err error, extra map[string]interface{}) {
	return func(ctx context.Context, service string, method string, err error, extra map[string]interface{}) {
		if err == nil {
			return
		}

		if extra == nil {
			extra = make(map[string]interface{})
		}

		extra["service"] = service
		extra["method"] = method
		extra["serviceError"] = true

		message := fmt.Sprintf("Service error in %s.%s", service, method)
		logger.LogError(ctx, message, err, extra)
	}
}

// extractUserIDFromRequest extrai o ID do usuário do token JWT
func extractUserIDFromRequest(r *http.Request, jwtService *auth.JWTService) string {
	if jwtService == nil {
		return ""
	}

	// Tentar extrair do header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Tentar extrair do cookie se não há header
		cookie, err := r.Cookie("token")
		if err != nil {
			return ""
		}
		authHeader = "Bearer " + cookie.Value
	}

	// Formato: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	token := parts[1]
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		return ""
	}

	// Extrair user ID do claims
	if sub, ok := claims["sub"].(string); ok {
		return sub
	}

	return ""
}

// RequestID middleware que adiciona Request ID ao contexto
func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

// RealIP middleware que adiciona o IP real do cliente
func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}

// Timeout middleware que adiciona timeout às requisições
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

// Heartbeat middleware para health checks
func Heartbeat(endpoint string) func(http.Handler) http.Handler {
	return middleware.Heartbeat(endpoint)
}

// Compress middleware para compressão gzip
func Compress(level int) func(http.Handler) http.Handler {
	return middleware.Compress(level, "text/html", "text/css", "text/plain", "text/javascript", "application/javascript", "application/json", "application/xml")
}

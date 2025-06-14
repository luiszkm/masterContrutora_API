// file: pkg/auth/jwt.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserContextKey = contextKey("userID")

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secretKey: []byte(secret)}
}

func (s *JWTService) GenerateToken(userID uuid.UUID, permissoes []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":        userID.String(),
		"permissoes": permissoes, // Adicionando permissões ao token
		"exp":        time.Now().Add(time.Hour * 8).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de assinatura inesperado: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("falha ao fazer parse do token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

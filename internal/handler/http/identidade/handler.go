package identidade

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/domain/identidade"
	dto "github.com/luiszkm/masterCostrutora/internal/service/identidade/dtos"
)

type registrarRequest struct {
	Nome           string `json:"nome"`
	Email          string `json:"email"`
	Senha          string `json:"senha"`
	ConfirmarSenha string `json:"confirmarSenha"`
}

type usuarioResponse struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type loginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
}

type Service interface {
	Registrar(ctx context.Context, input dto.RegistrarUsuarioInput) (*identidade.Usuario, error)
	Login(ctx context.Context, input dto.LoginInput) (string, error)
}
type Handler struct {
	service Service
	logger  *slog.Logger
}

// NovoIdentidadeHandler é o construtor que será chamado no main.go
func NovoIdentidadeHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) HandleRegistrar(w http.ResponseWriter, r *http.Request) {
	var req registrarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Senha != req.ConfirmarSenha {
		http.Error(w, "As senhas não conferem", http.StatusBadRequest)
		return
	}

	input := dto.RegistrarUsuarioInput{
		Nome:  req.Nome,
		Email: req.Email,
		Senha: req.Senha,
	}

	usuario, err := h.service.Registrar(r.Context(), input)
	if err != nil {
		// TODO: Tratar erros específicos, como email já existente (409 Conflict).
		h.logger.ErrorContext(r.Context(), "falha ao registrar usuário", "erro", err)
		http.Error(w, "Não foi possível registrar o usuário", http.StatusInternalServerError)
		return
	}

	resp := usuarioResponse{
		ID:    usuario.ID,
		Nome:  usuario.Nome,
		Email: usuario.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// HandleLogin trata a autenticação do usuário e retorna um token JWT.
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.LoginInput{
		Email: req.Email,
		Senha: req.Senha,
	}

	tokenString, err := h.service.Login(r.Context(), input)
	if err != nil {
		// TODO: Tratar erro de credenciais inválidas com status 401.
		h.logger.WarnContext(r.Context(), "tentativa de login falhou", "email", req.Email, "erro", err)
		http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
		return
	}
	// 1. Criamos o cookie com o token.
	cookie := http.Cookie{
		Name:     "jwt-token", // Nome do cookie
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 8), // Duração do cookie
		HttpOnly: true,                          // Impede o acesso via JavaScript (CRUCIAL para segurança)
		Secure:   true,                          // Garante que o cookie só seja enviado via HTTPS
		SameSite: http.SameSiteLaxMode,          // Ajuda a proteger contra ataques CSRF
		Path:     "/",                           // O cookie será válido para todo o site
	}

	// 2. Definimos o cookie na resposta.
	http.SetCookie(w, &cookie)

	// 3. Retornamos uma resposta de sucesso sem corpo.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login bem-sucedido"}`))

}

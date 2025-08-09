// file: internal/service/identidade/service.go
package identidade

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/authz"
	"github.com/luiszkm/masterCostrutora/internal/domain/identidade"
	dto "github.com/luiszkm/masterCostrutora/internal/service/identidade/dtos"
)

var (
	ErrCredenciaisInvalidas = fmt.Errorf("credenciais inválidas") // Define um erro customizado para credenciais inválidas
)

// Interfaces para as dependências externas que serão injetadas.
type JWTService interface {
	GenerateToken(userID uuid.UUID, permissoes []string) (string, error)
}

type Hasher interface {
	Hash(senha string) (string, error)
	Checar(senha, hash string) bool
}

// Service depende das interfaces, não das implementações concretas.
type Service struct {
	repo       identidade.UsuarioRepository
	hasher     Hasher
	jwtService JWTService
	logger     *slog.Logger
}

// NovoServico agora está alinhado com as interfaces.
func NovoServico(repo identidade.UsuarioRepository, hasher Hasher, jwtService JWTService, logger *slog.Logger) *Service {
	return &Service{
		repo:       repo,
		hasher:     hasher,
		jwtService: jwtService,
		logger:     logger,
	}
}

// Registrar agora aceita o DTO de entrada do serviço.
func (s *Service) Registrar(ctx context.Context, input dto.RegistrarUsuarioInput) (*identidade.Usuario, error) {
	const op = "service.identidade.Registrar"

	hash, err := s.hasher.Hash(input.Senha)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	permissoes := authz.GetPermissoesParaPapel(authz.PapelAdmin)

	novoUsuario := &identidade.Usuario{
		ID:         uuid.NewString(),
		Nome:       input.Nome,
		Email:      input.Email,
		SenhaHash:  hash,
		Permissoes: permissoes,
		Ativo:      true,
	}

	if err := s.repo.Salvar(ctx, novoUsuario); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return novoUsuario, nil
}

// Login agora aceita o DTO de entrada do serviço, corrigindo o erro de compilação.
func (s *Service) Login(ctx context.Context, input dto.LoginInput) (string, error) {
	const op = "service.identidade.Login"

	usuario, err := s.repo.BuscarPorEmail(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !s.hasher.Checar(input.Senha, usuario.SenhaHash) {
		return "", ErrCredenciaisInvalidas // Define este erro customizado no pacote
	}

	userID, err := uuid.Parse(usuario.ID)
	if err != nil {
		return "", fmt.Errorf("%s: id de usuário inválido: %w", op, err)
	}

	return s.jwtService.GenerateToken(userID, usuario.Permissoes)
}

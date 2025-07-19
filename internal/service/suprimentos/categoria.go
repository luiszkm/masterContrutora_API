package suprimentos

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

func (s *Service) CriarCategoria(ctx context.Context, input dto.CriarCategoriaInput) (*suprimentos.Categoria, error) {
	const op = "service.suprimentos.CriarCategoria"

	novaCategoria := &suprimentos.Categoria{
		ID:        uuid.NewString(),
		Nome:      input.Nome,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.categoriaRepo.Salvar(ctx, novaCategoria); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return novaCategoria, nil
}

func (s *Service) ListarCategorias(ctx context.Context) ([]*suprimentos.Categoria, error) {
	return s.categoriaRepo.ListarTodas(ctx)
}

func (s *Service) AtualizarCategoria(ctx context.Context, id string, input dto.AtualizarCategoriaInput) (*suprimentos.Categoria, error) {
	const op = "service.suprimentos.AtualizarCategoria"

	categoria, err := s.categoriaRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	categoria.Nome = input.Nome
	categoria.UpdatedAt = time.Now()

	if err := s.categoriaRepo.Atualizar(ctx, categoria); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return categoria, nil
}

func (s *Service) DeletarCategoria(ctx context.Context, id string) error {
	const op = "service.suprimentos.DeletarCategoria"
	// Primeiro, verifica se a categoria existe.
	if _, err := s.categoriaRepo.BuscarPorID(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return s.categoriaRepo.Deletar(ctx, id)
}

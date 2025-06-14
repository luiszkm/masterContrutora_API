// file: internal/service/suprimentos/service.go
package suprimentos

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

type Service struct {
	fornecedorRepo suprimentos.FornecedorRepository
	materialRepo   suprimentos.MaterialRepository
	logger         *slog.Logger
}

func NovoServico(fRepo suprimentos.FornecedorRepository, mRepo suprimentos.MaterialRepository, logger *slog.Logger) *Service {
	return &Service{fornecedorRepo: fRepo, materialRepo: mRepo, logger: logger}
}

func (s *Service) CadastrarFornecedor(ctx context.Context, input dto.CadastrarFornecedorInput) (*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.CadastrarFornecedor"

	// TODO: Adicionar validações de negócio (ex: formato do CNPJ).

	novoFornecedor := &suprimentos.Fornecedor{
		ID:        uuid.NewString(),
		Nome:      input.Nome,
		CNPJ:      input.CNPJ,
		Categoria: input.Categoria,
		Contato:   input.Contato,
		Email:     input.Email,
		Status:    "Ativo",
	}

	if err := s.fornecedorRepo.Salvar(ctx, novoFornecedor); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "novo fornecedor cadastrado", "fornecedor_id", novoFornecedor.ID)
	return novoFornecedor, nil
}

func (s *Service) ListarFornecedores(ctx context.Context) ([]*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.ListarFornecedores"
	fornecedores, err := s.fornecedorRepo.ListarTodos(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return fornecedores, nil
}

func (s *Service) CadastrarMaterial(ctx context.Context, input dto.CadastrarMaterialInput) (*suprimentos.Material, error) {
	novoMaterial := &suprimentos.Material{
		ID:              uuid.NewString(),
		Nome:            input.Nome,
		Descricao:       input.Descricao,
		UnidadeDeMedida: input.UnidadeDeMedida,
		Categoria:       input.Categoria,
	}
	if err := s.materialRepo.Salvar(ctx, novoMaterial); err != nil {
		return nil, fmt.Errorf("falha ao salvar material: %w", err)
	}
	return novoMaterial, nil
}
func (s *Service) ListarMateriais(ctx context.Context) ([]*suprimentos.Material, error) {
	return s.materialRepo.ListarTodos(ctx)
}

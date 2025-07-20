// file: internal/service/suprimentos/service.go
package suprimentos

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

var (
	ErrCategoriaExistente = fmt.Errorf("categoria já existe")
	ErrFornecedorInativo  = fmt.Errorf("fornecedor está inativo")
)

type EventPublisher interface {
	Publicar(ctx context.Context, evento bus.Evento)
}

type EtapaFinder interface {
	BuscarPorID(ctx context.Context, id string) (*obras.Etapa, error)
}
type FornecedorFinder interface {
	BuscarPorID(ctx context.Context, id string) (*suprimentos.Fornecedor, error)
}
type MaterialFinder interface {
	BuscarPorID(ctx context.Context, id string) (*suprimentos.Produto, error)
}
type Service struct {
	fornecedorRepo   suprimentos.FornecedorRepository
	produtoRepo      suprimentos.ProdutoRepository
	orcamentoRepo    suprimentos.OrcamentoRepository
	categoriaRepo    suprimentos.CategoriaRepository
	etapaFinder      EtapaFinder
	fornecedorFinder FornecedorFinder
	materialFinder   MaterialFinder
	eventBus         EventPublisher
	logger           *slog.Logger
}

func NovoServico(
	fRepo suprimentos.FornecedorRepository,
	mRepo suprimentos.ProdutoRepository,
	oRepo suprimentos.OrcamentoRepository,
	catRepo suprimentos.CategoriaRepository,
	eFinder EtapaFinder,
	fFinder FornecedorFinder,
	mFinder MaterialFinder,
	eventBus EventPublisher,
	logger *slog.Logger,
) *Service {
	return &Service{
		fornecedorRepo:   fRepo,
		produtoRepo:      mRepo,
		orcamentoRepo:    oRepo,
		categoriaRepo:    catRepo,
		etapaFinder:      eFinder,
		fornecedorFinder: fFinder,
		materialFinder:   mFinder,
		eventBus:         eventBus,
		logger:           logger,
	}
}

func (s *Service) CadastrarFornecedor(ctx context.Context, input dto.CadastrarFornecedorInput) (*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.CadastrarFornecedor"

	// TODO: Adicionar validações de negócio (ex: formato do CNPJ).

	novoFornecedor := &suprimentos.Fornecedor{
		ID:          uuid.NewString(),
		Nome:        input.Nome,
		CNPJ:        input.CNPJ,
		Contato:     &input.Contato,
		Email:       &input.Email,
		Status:      "Ativo",
		Endereco:    input.Endereco,    // NOVO
		Observacoes: input.Observacoes, // NOVO
		Avaliacao:   input.Avaliacao,
	}

	if err := s.fornecedorRepo.Salvar(ctx, novoFornecedor, input.CategoriaIDs); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "novo fornecedor cadastrado", "fornecedor_id", novoFornecedor.ID)
	return s.fornecedorRepo.BuscarPorID(ctx, novoFornecedor.ID)
}

func (s *Service) ListarFornecedores(ctx context.Context) ([]*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.ListarFornecedores"
	fornecedores, err := s.fornecedorRepo.ListarTodos(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return fornecedores, nil
}

func (s *Service) AtualizarFornecedor(ctx context.Context, id string, input dto.AtualizarFornecedorInput) (*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.AtualizarFornecedor"

	fornecedor, err := s.fornecedorRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: fornecedor com id [%s] não encontrado: %w", op, id, err)
	}

	// Somente altera os campos enviados
	if input.Nome != nil {
		fornecedor.Nome = *input.Nome
	}
	if input.CNPJ != nil {
		fornecedor.CNPJ = *input.CNPJ
	}
	if input.Contato != nil {
		fornecedor.Contato = input.Contato
	}
	if input.Email != nil {
		fornecedor.Email = input.Email
	}
	if input.Status != nil {
		fornecedor.Status = *input.Status
	}
	if input.Endereco != nil {
		fornecedor.Endereco = input.Endereco
	}
	if input.Avaliacao != nil {
		fornecedor.Avaliacao = input.Avaliacao
	}
	if input.Observacoes != nil {
		fornecedor.Observacoes = input.Observacoes
	}

	// Atualização de categorias somente se input.CategoriaIDs != nil
	err = s.fornecedorRepo.Atualizar(ctx, fornecedor, input.CategoriaIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao atualizar fornecedor: %w", op, err)
	}

	s.logger.InfoContext(ctx, "fornecedor atualizado", "fornecedor_id", fornecedor.ID)
	return s.fornecedorRepo.BuscarPorID(ctx, id)
}

func (s *Service) DeletarFornecedor(ctx context.Context, id string) error {
	const op = "service.suprimentos.DeletarFornecedor"
	fornecedor, err := s.fornecedorRepo.BuscarPorID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: fornecedor com id [%s] não encontrado: %w", op, id, err)
	}
	if fornecedor.Status == "Inativo" {
		return fmt.Errorf("%s: fornecedor com id [%s] já está inativo", op, id)
	}
	if err := s.fornecedorRepo.Deletar(ctx, id); err != nil {
		return fmt.Errorf("%s: falha ao deletar fornecedor: %w", op, err)
	}
	s.logger.InfoContext(ctx, "fornecedor excluído (soft delete)", "fornecedor_id", id)
	return nil
}

func (s *Service) BuscarPorID(ctx context.Context, id string) (*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.BuscarPorID"
	fornecedor, err := s.fornecedorRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: fornecedor com id [%s] não encontrado: %w", op, id, err)
	}
	s.logger.InfoContext(ctx, "fornecedor encontrado", "fornecedor_id", fornecedor.ID)
	return fornecedor, nil
}

func (s *Service) CadastrarMaterial(ctx context.Context, input dto.CadastrarProdutoInput) (*suprimentos.Produto, error) {
	novoMaterial := &suprimentos.Produto{
		ID:              uuid.NewString(),
		Nome:            input.Nome,
		Descricao:       input.Descricao,
		UnidadeDeMedida: input.UnidadeDeMedida,
		Categoria:       input.Categoria,
	}
	if err := s.produtoRepo.Salvar(ctx, novoMaterial); err != nil {
		return nil, fmt.Errorf("falha ao salvar material: %w", err)
	}
	return novoMaterial, nil
}
func (s *Service) ListarMateriais(ctx context.Context) ([]*suprimentos.Produto, error) {
	return s.produtoRepo.ListarTodos(ctx)
}

// file: internal/service/suprimentos/service.go
package suprimentos

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
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
	BuscarPorID(ctx context.Context, id string) (*suprimentos.Material, error)
}
type Service struct {
	fornecedorRepo   suprimentos.FornecedorRepository
	materialRepo     suprimentos.MaterialRepository
	orcamentoRepo    suprimentos.OrcamentoRepository
	etapaFinder      EtapaFinder
	fornecedorFinder FornecedorFinder
	materialFinder   MaterialFinder
	eventBus         EventPublisher
	logger           *slog.Logger
}

func NovoServico(
	fRepo suprimentos.FornecedorRepository,
	mRepo suprimentos.MaterialRepository,
	oRepo suprimentos.OrcamentoRepository,
	eFinder EtapaFinder,
	fFinder FornecedorFinder,
	mFinder MaterialFinder,
	eventBus EventPublisher,
	logger *slog.Logger,
) *Service {
	return &Service{
		fornecedorRepo:   fRepo,
		materialRepo:     mRepo,
		orcamentoRepo:    oRepo,
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

func (s *Service) AtualizarFornecedor(ctx context.Context, id string, input dto.AtualizarFornecedorInput) (*suprimentos.Fornecedor, error) {
	const op = "service.suprimentos.AtualizarFornecedor"
	fornecedor, err := s.fornecedorRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: fornecedor com id [%s] não encontrado: %w", op, id, err)
	}
	fornecedor.Nome = input.Nome
	fornecedor.CNPJ = input.CNPJ
	fornecedor.Categoria = input.Categoria
	fornecedor.Contato = input.Contato
	fornecedor.Email = input.Email
	fornecedor.Status = input.Status
	if err := s.fornecedorRepo.Atualizar(ctx, fornecedor); err != nil {
		return nil, fmt.Errorf("%s: falha ao atualizar fornecedor: %w", op, err)
	}
	s.logger.InfoContext(ctx, "fornecedor atualizado", "fornecedor_id", fornecedor.ID)
	return fornecedor, nil
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

func (s *Service) CriarOrcamento(ctx context.Context, etapaID string, input dto.CriarOrcamentoInput) (*suprimentos.Orcamento, error) {
	const op = "service.suprimentos.CriarOrcamento"

	// 1. Validações de Existência (colaboração entre contextos)
	if _, err := s.etapaFinder.BuscarPorID(ctx, etapaID); err != nil {
		return nil, fmt.Errorf("%s: etapa com id [%s] não encontrada: %w", op, etapaID, err)
	}
	if _, err := s.fornecedorFinder.BuscarPorID(ctx, input.FornecedorID); err != nil {
		return nil, fmt.Errorf("%s: fornecedor com id [%s] não encontrado: %w", op, input.FornecedorID, err)
	}
	for _, item := range input.Itens {
		if _, err := s.materialFinder.BuscarPorID(ctx, item.MaterialID); err != nil {
			return nil, fmt.Errorf("%s: material com id [%s] não encontrado: %w", op, item.MaterialID, err)
		}
	}

	// 2. Lógica de Negócio e Criação do Agregado
	valorTotal := 0.0
	orcamentoID := uuid.NewString()
	itensOrcamento := make([]suprimentos.ItemOrcamento, len(input.Itens))

	for i, itemInput := range input.Itens {
		valorTotal += itemInput.Quantidade * itemInput.ValorUnitario
		itensOrcamento[i] = suprimentos.ItemOrcamento{
			ID:            uuid.NewString(),
			OrcamentoID:   orcamentoID,
			MaterialID:    itemInput.MaterialID,
			Quantidade:    itemInput.Quantidade,
			ValorUnitario: itemInput.ValorUnitario,
		}
	}

	orcamento := &suprimentos.Orcamento{
		ID:           orcamentoID,
		Numero:       input.Numero,
		EtapaID:      etapaID,
		FornecedorID: input.FornecedorID,
		Itens:        itensOrcamento,
		ValorTotal:   valorTotal,
		Status:       "Em Aberto",
		DataEmissao:  time.Now(),
	}

	// 3. Persistência Atômica através do repositório
	if err := s.orcamentoRepo.Salvar(ctx, orcamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "novo orçamento cadastrado", "orcamento_id", orcamento.ID, "etapa_id", etapaID)
	return orcamento, nil
}

func (s *Service) AtualizarStatusOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarStatusOrcamentoInput) (*suprimentos.Orcamento, error) {
	const op = "service.suprimentos.AtualizarStatusOrcamento"

	orcamento, err := s.orcamentoRepo.BuscarPorID(ctx, orcamentoID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Adicionar validações de transição de status.
	// Ex: um orçamento 'Rejeitado' não pode ser 'Pago'.
	orcamento.Status = input.Status

	if err := s.orcamentoRepo.Atualizar(ctx, orcamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	payload := events.OrcamentoStatusAtualizadoPayload{
		OrcamentoID: orcamento.ID,
		EtapaID:     orcamento.EtapaID,
		NovoStatus:  orcamento.Status,
		Valor:       orcamento.ValorTotal,
	}
	evento := bus.Evento{
		Nome:    events.OrcamentoStatusAtualizado,
		Payload: payload,
	}
	s.eventBus.Publicar(ctx, evento)
	// --- FIM DA MUDANÇA ---

	s.logger.InfoContext(ctx, "status do orçamento atualizado e evento publicado", "orcamento_id", orcamentoID)
	return orcamento, nil
}

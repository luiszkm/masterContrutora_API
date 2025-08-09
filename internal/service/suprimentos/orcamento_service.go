package suprimentos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

func (s *Service) ListarOrcamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.OrcamentoListItemDTO], error) {
	orcamentos, paginacao, err := s.orcamentoRepo.ListarOrcamentos(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar orçamentos: %w", err) // TODO: Mudar para o DTO correto
	}
	return &common.RespostaPaginada[*dto.OrcamentoListItemDTO]{
		Dados:     orcamentos,
		Paginacao: *paginacao,
	}, nil
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

	agora := time.Now()
	ano, mes := agora.Year(), agora.Month()
	count, err := s.orcamentoRepo.ContarPorMesAno(ctx, ano, mes)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao contar orçamentos para numeração: %w", op, err)
	}
	sequencial := count + 1
	mesAbrev := strings.ToUpper(mes.String()[0:3])
	numeroFormatado := fmt.Sprintf("ORC-%d-%s-%03d", ano, mesAbrev, sequencial)

	// 2. Lógica de Negócio e Criação do Agregado
	valorTotal := 0.0
	orcamentoID := uuid.NewString()
	itensOrcamento := make([]suprimentos.ItemOrcamento, len(input.Itens))

	for i, itemInput := range input.Itens {
		// Tenta encontrar o produto pelo nome
		produto, err := s.produtoRepo.BuscarPorNome(ctx, itemInput.NomeProduto)

		// Se não encontrar, cria um novo
		if err != nil {
			if errors.Is(err, postgres.ErrNaoEncontrado) {
				produto = &suprimentos.Produto{
					ID:              uuid.NewString(),
					Nome:            itemInput.NomeProduto,
					UnidadeDeMedida: itemInput.UnidadeDeMedida,
					Categoria:       itemInput.Categoria,
					CreatedAt:       agora,
					UpdatedAt:       agora,
				}
				if err := s.produtoRepo.Salvar(ctx, produto); err != nil {
					return nil, fmt.Errorf("%s: falha ao criar novo produto '%s': %w", op, itemInput.NomeProduto, err)
				}
			} else {
				// Outro erro de banco de dados
				return nil, fmt.Errorf("%s: falha ao buscar produto '%s': %w", op, itemInput.NomeProduto, err)
			}
		}

		// Monta o item do orçamento com o ID do produto (encontrado ou recém-criado)
		valorTotal += itemInput.Quantidade * itemInput.ValorUnitario
		itensOrcamento[i] = suprimentos.ItemOrcamento{
			ID:            uuid.NewString(),
			OrcamentoID:   orcamentoID,
			ProdutoID:     produto.ID, // Usa o ID correto
			Quantidade:    itemInput.Quantidade,
			ValorUnitario: itemInput.ValorUnitario,
		}
	}

	orcamento := &suprimentos.Orcamento{
		ID:           orcamentoID,
		Numero:       numeroFormatado,
		EtapaID:      etapaID,
		FornecedorID: input.FornecedorID,
		Itens:        itensOrcamento,
		ValorTotal:   valorTotal,
		Status:       "Em Aberto",
		DataEmissao:  agora,
		CreatedAt:    agora,
		UpdatedAt:    agora,
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
	
	// Capturar status anterior antes da mudança
	statusAnterior := orcamento.Status
	orcamento.Status = input.Status

	if orcamento.Status == "Aprovado" {
		now := time.Now()
		orcamento.DataAprovacao = &now
	}

	if err := s.orcamentoRepo.Atualizar(ctx, orcamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	payload := events.OrcamentoStatusAtualizadoPayload{
		OrcamentoID:    orcamento.ID,
		EtapaID:        orcamento.EtapaID,
		StatusAnterior: statusAnterior,
		NovoStatus:     orcamento.Status,
		Valor:          orcamento.ValorTotal,
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

func (s *Service) BuscarOrcamentoPorID(ctx context.Context, id string) (*dto.OrcamentoDetalhadoDTO, error) {
	orcamento, err := s.orcamentoRepo.BuscarPorDetalhesID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar orçamento detalhado: %w", err)
	}
	return orcamento, nil
}
func (s *Service) AtualizarOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarOrcamentoInput) (*dto.OrcamentoDetalhadoDTO, error) {
	const op = "service.suprimentos.AtualizarOrcamento"

	// 1. Busca o agregado de domínio COMPLETO do banco de dados.
	// É crucial que este método retorne a entidade `suprimentos.Orcamento` e não um DTO.
	orcamento, err := s.orcamentoRepo.BuscarPorID(ctx, orcamentoID)
	if err != nil {
		return nil, err // Retorna erro 404 se não for encontrado
	}

	// 2. Valida se os novos IDs de etapa e fornecedor existem.
	if _, err := s.etapaFinder.BuscarPorID(ctx, input.EtapaID); err != nil {
		return nil, fmt.Errorf("%s: etapa de destino com id [%s] não encontrada: %w", op, input.EtapaID, err)
	}
	if _, err := s.fornecedorFinder.BuscarPorID(ctx, input.FornecedorID); err != nil {
		return nil, fmt.Errorf("%s: fornecedor de destino com id [%s] não encontrado: %w", op, input.FornecedorID, err)
	}

	// 3. Lógica de "Upsert" de produtos e montagem da nova lista de itens.
	valorTotal := 0.0
	novosItensOrcamento := make([]suprimentos.ItemOrcamento, len(input.Itens))
	for i, itemInput := range input.Itens {
		produto, err := s.produtoRepo.BuscarPorNome(ctx, itemInput.NomeProduto)
		if err != nil {
			if errors.Is(err, postgres.ErrNaoEncontrado) {
				produto = &suprimentos.Produto{
					ID:              uuid.NewString(),
					Nome:            itemInput.NomeProduto,
					UnidadeDeMedida: itemInput.UnidadeDeMedida,
					Categoria:       itemInput.Categoria,
				}
				if err := s.produtoRepo.Salvar(ctx, produto); err != nil {
					return nil, fmt.Errorf("%s: falha ao criar novo produto '%s': %w", op, itemInput.NomeProduto, err)
				}
			} else {
				return nil, fmt.Errorf("%s: falha ao buscar produto '%s': %w", op, itemInput.NomeProduto, err)
			}
		}

		valorTotal += itemInput.Quantidade * itemInput.ValorUnitario
		novosItensOrcamento[i] = suprimentos.ItemOrcamento{
			ID:            uuid.NewString(),
			OrcamentoID:   orcamentoID,
			ProdutoID:     produto.ID,
			Quantidade:    itemInput.Quantidade,
			ValorUnitario: itemInput.ValorUnitario,
		}
	}

	// 4. MODIFICA o objeto de domínio em memória com os novos dados.
	// Esta é a etapa que estava a falhar silenciosamente.
	orcamento.EtapaID = input.EtapaID
	orcamento.FornecedorID = input.FornecedorID
	orcamento.Observacoes = input.Observacoes
	orcamento.CondicoesPagamento = input.CondicoesPagamento
	orcamento.ValorTotal = valorTotal
	orcamento.Itens = novosItensOrcamento
	// Poderíamos adicionar lógica para alterar o status aqui se necessário, ex:
	// orcamento.Status = "Revisado"

	// 5. Persiste o agregado COMPLETO e MODIFICADO.
	if err := s.orcamentoRepo.Atualizar(ctx, orcamento); err != nil {
		return nil, fmt.Errorf("%s: falha ao persistir atualização no repositório: %w", op, err)
	}

	// 6. Retorna a visão detalhada do orçamento para confirmar as alterações.
	s.logger.InfoContext(ctx, "orçamento atualizado com sucesso", "orcamento_id", orcamentoID)
	return s.orcamentoRepo.BuscarPorDetalhesID(ctx, orcamentoID)
}

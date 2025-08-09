// file: internal/service/financeiro/service.go
package financeiro

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
)

// Finders para buscar entidades de outros contextos
type PagamentoRepository interface {
	// A interface agora espera um DBTX
	Salvar(ctx context.Context, dbtx db.DBTX, pagamento *financeiro.RegistroDePagamento) error
	ListarPagamentos(ctx context.Context, filtros common.ListarFiltros) ([]*financeiro.RegistroDePagamento, *common.PaginacaoInfo, error)
}
type FuncionarioFinder interface {
	BuscarPorID(ctx context.Context, id string) (*pessoal.Funcionario, error)
}
type ObraFinder interface {
	BuscarPorID(ctx context.Context, id string) (*obras.Obra, error)
}
type EventPublisher interface {
	Publicar(ctx context.Context, evento bus.Evento)
}
type ApontamentoRepository interface {
	// Precisamos de uma forma de buscar para atualizar e salvar
	BuscarPorID(ctx context.Context, id string) (*pessoal.ApontamentoQuinzenal, error)
	Atualizar(ctx context.Context, dbtx db.DBTX, apontamento *pessoal.ApontamentoQuinzenal) error
}

// Erros de negócio customizados
var ErrFuncionarioInativo = errors.New("não é possível registrar pagamento para um funcionário inativo")

type Service struct {
	pagamentoRepo     PagamentoRepository
	apontamentoRepo   ApontamentoRepository
	funcionarioFinder FuncionarioFinder
	obraFinder        ObraFinder
	eventBus          EventPublisher
	dbpool            *pgxpool.Pool
	logger            *slog.Logger
}

func NovoServico(
	pRepo PagamentoRepository,
	aRepo ApontamentoRepository,
	fFinder FuncionarioFinder,
	oFinder ObraFinder,
	bus EventPublisher,
	dbpool *pgxpool.Pool,
	logger *slog.Logger,
) *Service {
	return &Service{
		pagamentoRepo:     pRepo,
		apontamentoRepo:   aRepo,
		funcionarioFinder: fFinder,
		obraFinder:        oFinder,
		eventBus:          bus,
		dbpool:            dbpool,
		logger:            logger,
	}
}

func (s *Service) RegistrarPagamento(ctx context.Context, input dto.RegistrarPagamentoInput) (*financeiro.RegistroDePagamento, error) {
	const op = "service.financeiro.RegistrarPagamento"

	// Validação de existência
	funcionario, err := s.funcionarioFinder.BuscarPorID(ctx, input.FuncionarioID)
	if err != nil {
		return nil, fmt.Errorf("%s: funcionário não encontrado: %w", op, err)
	}
	if _, err := s.obraFinder.BuscarPorID(ctx, input.ObraID); err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	// Regra de Negócio
	if funcionario.Status != "Ativo" {
		return nil, ErrFuncionarioInativo
	}

	novoPagamento := &financeiro.RegistroDePagamento{
		ID:                uuid.NewString(),
		FuncionarioID:     input.FuncionarioID,
		ObraID:            input.ObraID,
		PeriodoReferencia: input.PeriodoReferencia,
		ValorCalculado:    input.ValorCalculado,
		DataDeEfetivacao:  time.Now(),
		ContaBancariaID:   input.ContaBancariaID,
	}

	// --- AQUI ESTÁ A MUDANÇA PRINCIPAL ---
	// Em vez de um 'tx' indefinido, passamos o pool de conexões principal (s.dbpool).
	// O *pgxpool.Pool (s.dbpool) satisfaz a interface DBTX,
	// então a chamada ao método do repositório, que agora espera um DBTX, funciona perfeitamente.
	if err := s.pagamentoRepo.Salvar(ctx, s.dbpool, novoPagamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "novo pagamento registrado", "pagamento_id", novoPagamento.ID)
	return novoPagamento, nil
}

func (s *Service) RegistrarPagamentosEmLote(ctx context.Context, input dto.RegistrarPagamentoEmLoteInput) (*dto.ResultadoExecucaoLote, error) {
	const op = "service.financeiro.RegistrarPagamentosEmLote"

	dataEfetivacao, err := time.Parse("2006-01-02", input.DataDeEfetivacao)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de efetivação inválido: %w", op, err)
	}

	resultado := &dto.ResultadoExecucaoLote{
		Resumo:   dto.ResumoExecucao{TotalSolicitado: len(input.ApontamentoIDs)},
		Sucessos: make([]dto.DetalheSucesso, 0),
		Falhas:   make([]dto.DetalheFalha, 0),
	}

	for _, apontamentoID := range input.ApontamentoIDs {
		tx, err := s.dbpool.Begin(ctx)
		if err != nil {
			s.logger.ErrorContext(ctx, "falha ao iniciar transação para pagamento em lote", "apontamentoId", apontamentoID, "erro", err)
			resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{ApontamentoID: apontamentoID, Motivo: "Erro interno: não foi possível iniciar a transação."})
			resultado.Resumo.TotalFalha++
			continue
		}

		// Usamos um closure para facilitar o gerenciamento do rollback com defer.
		err = func(tx pgx.Tx) error {
			apontamento, err := s.apontamentoRepo.BuscarPorID(ctx, apontamentoID)
			if err != nil {
				return errors.New("Apontamento não encontrado.")
			}

			if err := apontamento.AprovarEPagar(); err != nil {
				return err
			}

			// CORREÇÃO: Passamos `tx` para o método do repositório.
			if err := s.apontamentoRepo.Atualizar(ctx, tx, apontamento); err != nil {
				return errors.New("Erro ao salvar atualização do apontamento.")
			}

			novoPagamento := &financeiro.RegistroDePagamento{
				ID:                uuid.NewString(),
				FuncionarioID:     apontamento.FuncionarioID,
				ObraID:            apontamento.ObraID,
				PeriodoReferencia: fmt.Sprintf("%s a %s", apontamento.PeriodoInicio.Format("02/01"), apontamento.PeriodoFim.Format("02/01/2006")),
				ValorCalculado:    apontamento.ValorTotalCalculado,
				DataDeEfetivacao:  dataEfetivacao,
				ContaBancariaID:   input.ContaBancariaID,
			}

			// CORREÇÃO: Passamos `tx` para o método do repositório.
			if err := s.pagamentoRepo.Salvar(ctx, tx, novoPagamento); err != nil {
				return errors.New("Erro ao salvar registro de pagamento.")
			}

			// Publica o evento (ainda antes do commit, mas só se tudo deu certo até aqui)
			s.eventBus.Publicar(ctx, bus.Evento{ /* ... payload do evento ... */ })

			// Se tudo deu certo, retorna nil para o closure, permitindo o commit.
			return nil
		}(tx)

		if err != nil {
			tx.Rollback(ctx) // Garante o rollback
			resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{ApontamentoID: apontamentoID, Motivo: err.Error()})
			resultado.Resumo.TotalFalha++
		} else {
			if err := tx.Commit(ctx); err != nil {
				// Se o commit falhar, é um erro grave.
				resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{ApontamentoID: apontamentoID, Motivo: "Erro interno ao finalizar o pagamento."})
				resultado.Resumo.TotalFalha++
			} else {
				// Sucesso!
				resultado.Sucessos = append(resultado.Sucessos, dto.DetalheSucesso{ApontamentoID: apontamentoID /* ... */})
				resultado.Resumo.TotalSucesso++
			}
		}
	}

	return resultado, nil
}

func (s *Service) ListarPagamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*financeiro.RegistroDePagamento], error) {
	const op = "service.financeiro.ListarPagamentos"

	pagamentos, paginacao, err := s.pagamentoRepo.ListarPagamentos(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resposta := &common.RespostaPaginada[*financeiro.RegistroDePagamento]{
		Dados:     pagamentos,
		Paginacao: *paginacao,
	}

	return resposta, nil
}

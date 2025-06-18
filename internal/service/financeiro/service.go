// file: internal/service/financeiro/service.go
package financeiro

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
)

// Finders para buscar entidades de outros contextos
type FuncionarioFinder interface {
	BuscarPorID(ctx context.Context, id string) (*pessoal.Funcionario, error)
}
type ObraFinder interface {
	BuscarPorID(ctx context.Context, id string) (*obras.Obra, error)
}

// Erros de negócio customizados
var ErrFuncionarioInativo = errors.New("não é possível registrar pagamento para um funcionário inativo")

type Service struct {
	repo              financeiro.Repository
	funcionarioFinder FuncionarioFinder
	obraFinder        ObraFinder
	logger            *slog.Logger
}

func NovoServico(repo financeiro.Repository, fFinder FuncionarioFinder, oFinder ObraFinder, logger *slog.Logger) *Service {
	return &Service{
		repo:              repo,
		funcionarioFinder: fFinder,
		obraFinder:        oFinder,
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

	if err := s.repo.Salvar(ctx, novoPagamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "novo pagamento registrado", "pagamento_id", novoPagamento.ID)
	return novoPagamento, nil
}

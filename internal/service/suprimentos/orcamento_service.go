package suprimentos

import (
	"context"
	"fmt"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
)

func (s *Service) ListarOrcamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*suprimentos.Orcamento], error) {
	orcamentos, paginacao, err := s.orcamentoRepo.ListarTodos(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar orçamentos: %w", err)
	}
	return &common.RespostaPaginada[*suprimentos.Orcamento]{
		Dados:     orcamentos,
		Paginacao: *paginacao,
	}, nil
}

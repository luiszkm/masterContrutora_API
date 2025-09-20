package suprimentos

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

func (s *Service) CriarServico(ctx context.Context, input dto.CriarServicoInput) (*suprimentos.Servico, error) {
	const op = "service.suprimentos.CriarServico"

	novoServico := &suprimentos.Servico{
		ID:        uuid.NewString(),
		Nome:      input.Nome,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.servicoRepo.Salvar(ctx, novoServico); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return novoServico, nil
}

func (s *Service) ListarServicos(ctx context.Context) ([]*suprimentos.Servico, error) {
	return s.servicoRepo.ListarTodos(ctx)
}

func (s *Service) ListarServicosPaginado(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*suprimentos.Servico], error) {
	const op = "service.suprimentos.ListarServicosPaginado"

	servicos, paginacaoInfo, err := s.servicoRepo.Listar(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resposta := &common.RespostaPaginada[*suprimentos.Servico]{
		Dados:     servicos,
		Paginacao: *paginacaoInfo,
	}

	return resposta, nil
}

func (s *Service) AtualizarServico(ctx context.Context, id string, input dto.AtualizarServicoInput) (*suprimentos.Servico, error) {
	const op = "service.suprimentos.AtualizarServico"

	servico, err := s.servicoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	servico.Nome = input.Nome
	servico.UpdatedAt = time.Now()

	if err := s.servicoRepo.Atualizar(ctx, servico); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return servico, nil
}

func (s *Service) BuscarServico(ctx context.Context, id string) (*suprimentos.Servico, error) {
	const op = "service.suprimentos.BuscarServico"
	servico, err := s.servicoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return servico, nil
}

func (s *Service) DeletarServico(ctx context.Context, id string) error {
	const op = "service.suprimentos.DeletarServico"
	if _, err := s.servicoRepo.BuscarPorID(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return s.servicoRepo.Deletar(ctx, id)
}
package pessoal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
)

var (
	ErrFuncionarioSemApontamentoAnterior = errors.New("funcionário não possui um apontamento anterior para ser usado como template")
)

func (s *Service) ListarComUltimoApontamento(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ListagemFuncionarioDTO, *common.PaginacaoInfo, error) {
	const op = "service.pessoal.ListarComUltimoApontamento"

	funcionarios, paginacao, err := s.querier.ListarComUltimoApontamento(ctx, filtros)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "lista de funcionários com último apontamento recuperada", "total", len(funcionarios))
	return funcionarios, paginacao, nil

}

func (s *Service) CriarApontamento(ctx context.Context, input dto.CriarApontamentoInput) (*pessoal.ApontamentoQuinzenal, error) {
	const op = "service.pessoal.CriarApontamento"

	// Validações
	if _, err := s.repo.BuscarPorID(ctx, input.FuncionarioID); err != nil {
		return nil, fmt.Errorf("%s: funcionário com id [%s] não encontrado: %w", op, input.FuncionarioID, err)
	}
	if _, err := s.obraFinder.BuscarPorID(ctx, input.ObraID); err != nil {
		return nil, fmt.Errorf("%s: obra com id [%s] não encontrada: %w", op, input.ObraID, err)
	}

	inicio, err := time.Parse("2006-01-02", input.PeriodoInicio)
	if err != nil {
		return nil, fmt.Errorf("%s: data de início inválida: %w", op, err)
	}
	fim, err := time.Parse("2006-01-02", input.PeriodoFim)
	if err != nil {
		return nil, fmt.Errorf("%s: data de fim inválida: %w", op, err)
	}
	ValorTotalCalculado := (input.Diaria * float64(input.DiasTrabalhados)) + input.ValorAdicional - input.Descontos - input.Adiantamento

	apontamento := &pessoal.ApontamentoQuinzenal{
		ID:                  uuid.NewString(),
		FuncionarioID:       input.FuncionarioID,
		ObraID:              input.ObraID,
		PeriodoInicio:       inicio,
		PeriodoFim:          fim,
		Status:              pessoal.StatusApontamentoEmAberto,
		Diaria:              input.Diaria,
		DiasTrabalhados:     input.DiasTrabalhados,
		Adicionais:          input.ValorAdicional,
		Descontos:           input.Descontos,
		Adiantamentos:       input.Adiantamento,
		ValorTotalCalculado: ValorTotalCalculado,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.apontamentoRepo.Salvar(ctx, apontamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "novo apontamento quinzenal criado", "apontamento_id", apontamento.ID)
	return apontamento, nil
}
func (s *Service) AprovarApontamento(ctx context.Context, apontamentoID string) (*pessoal.ApontamentoQuinzenal, error) {
	const op = "service.pessoal.AprovarApontamento"

	// 1. Carrega o agregado do banco.
	apontamento, err := s.apontamentoRepo.BuscarPorID(ctx, apontamentoID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// 2. Executa o método de negócio do próprio agregado (Rich Domain Model).
	// Toda a lógica e validação de estado estão encapsuladas aqui!
	if err := apontamento.Aprovar(); err != nil {
		return nil, fmt.Errorf("%s: regra de negócio violada: %w", op, err)
	}

	// 3. Persiste o estado atualizado do agregado.
	if err := s.apontamentoRepo.Atualizar(ctx, apontamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "apontamento aprovado com sucesso", "apontamento_id", apontamento.ID)
	return apontamento, nil
}
func (s *Service) RegistrarPagamentoApontamento(ctx context.Context, apontamentoID string, contaPagamentoID string) (*pessoal.ApontamentoQuinzenal, error) {
	const op = "service.pessoal.RegistrarPagamentoApontamento"

	apontamento, err := s.apontamentoRepo.BuscarPorID(ctx, apontamentoID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Usa o método do nosso Rich Domain Model.
	if err := apontamento.RegistrarPagamento(); err != nil {
		return nil, fmt.Errorf("%s: regra de negócio violada: %w", op, err)
	}

	if err := s.apontamentoRepo.Atualizar(ctx, apontamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Publica o evento para que o contexto Financeiro possa agir.
	payload := events.PagamentoApontamentoRealizadoPayload{
		FuncionarioID:     apontamento.FuncionarioID,
		ObraID:            apontamento.ObraID,
		PeriodoReferencia: fmt.Sprintf("%s a %s", apontamento.PeriodoInicio.Format("02/01"), apontamento.PeriodoFim.Format("02/01/2006")),
		ValorCalculado:    apontamento.ValorTotalCalculado,
		DataDeEfetivacao:  time.Now(),
		ContaBancariaID:   contaPagamentoID,
	}
	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.PagamentoApontamentoRealizado,
		Payload: payload,
	})

	s.logger.InfoContext(ctx, "pagamento de apontamento registrado e evento publicado", "apontamento_id", apontamentoID)
	return apontamento, nil
}

func (s *Service) ListarApontamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*pessoal.ApontamentoQuinzenal], error) {
	apontamentos, paginacao, err := s.apontamentoRepo.Listar(ctx, filtros)
	if err != nil {
		return nil, err
	}
	return &common.RespostaPaginada[*pessoal.ApontamentoQuinzenal]{Dados: apontamentos, Paginacao: *paginacao}, nil
}

func (s *Service) ListarApontamentosPorFuncionario(ctx context.Context, funcionarioID string, filtros common.ListarFiltros) (*common.RespostaPaginada[*pessoal.ApontamentoQuinzenal], error) {
	apontamentos, paginacao, err := s.apontamentoRepo.ListarPorFuncionarioID(ctx, funcionarioID, filtros)
	if err != nil {
		return nil, err
	}
	return &common.RespostaPaginada[*pessoal.ApontamentoQuinzenal]{Dados: apontamentos, Paginacao: *paginacao}, nil
}

func (s *Service) AtualizarApontamento(ctx context.Context, id string, input dto.AtualizarApontamentoInput) (*pessoal.ApontamentoQuinzenal, error) {
	const op = "service.pessoal.AtualizarApontamento"

	// 1. Busca o agregado que será modificado.
	apontamento, err := s.apontamentoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Parse PeriodoInicio and PeriodoFim from string to time.Time
	periodoInicio, err := time.Parse("2006-01-02", input.PeriodoInicio)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de início inválido: %w", op, err)
	}
	periodoFim, err := time.Parse("2006-01-02", input.PeriodoFim)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de fim inválido: %w", op, err)
	}

	// 2. Executa o método de negócio do agregado. A diária não é mais necessária aqui.
	err = apontamento.AtualizarValores(
		input.DiasTrabalhados,
		input.Diaria,
		input.Descontos,
		input.Adiantamento,
		input.ValorAdicional,
		periodoInicio,
		periodoFim,
		input.ObraID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: regra de negócio violada: %w", op, err)
	}

	// 3. Persiste o estado atualizado.
	if err := s.apontamentoRepo.Atualizar(ctx, apontamento); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "apontamento atualizado com sucesso", "apontamento_id", id)
	return apontamento, nil
}

func (s *Service) ReplicarParaProximaQuinzena(ctx context.Context, input dto.ReplicarApontamentosInput) (*dto.ResultadoReplicacao, error) {
	const op = "service.pessoal.ReplicarParaProximaQuinzena"

	resultado := &dto.ResultadoReplicacao{
		Resumo: dto.ResumoReplicacao{
			TotalSolicitado: len(input.FuncionarioIDs),
		},
		Sucessos: make([]dto.DetalheSucesso, 0),
		Falhas:   make([]dto.DetalheFalha, 0),
	}

	for _, funcID := range input.FuncionarioIDs {
		// ADR-011: A validação ocorre no serviço de aplicação, antes de criar o agregado.
		existeEmAberto, err := s.apontamentoRepo.ExisteApontamentoEmAberto(ctx, funcID)
		if err != nil {
			s.logger.ErrorContext(ctx, "falha ao verificar apontamento em aberto", "funcionarioId", funcID, "erro", err)
			resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{FuncionarioID: funcID, Motivo: "Erro interno ao verificar apontamentos."})
			resultado.Resumo.TotalFalha++
			continue
		}
		if existeEmAberto {
			motivo := "Este funcionário já possui um apontamento com o status 'EM_ABERTO'."
			resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{FuncionarioID: funcID, Motivo: motivo})
			resultado.Resumo.TotalFalha++
			continue
		}

		// Lógica de "Template" (Item 2.2 do Adendo V3)
		ultimoApontamento, err := s.apontamentoRepo.BuscarUltimoPorFuncionarioID(ctx, funcID)
		if err != nil {
			motivo := ""
			if errors.Is(err, ErrFuncionarioSemApontamentoAnterior) { // Usando um erro mais específico que pode vir do repo
				motivo = ErrFuncionarioSemApontamentoAnterior.Error()
			} else {
				s.logger.ErrorContext(ctx, "falha ao buscar ultimo apontamento", "funcionarioId", funcID, "erro", err)
				motivo = "Erro interno ao buscar histórico."
			}
			resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{FuncionarioID: funcID, Motivo: motivo})
			resultado.Resumo.TotalFalha++
			continue
		}

		// Cria o novo apontamento
		novoApontamento := criarApontamentoAPartirDeTemplate(ultimoApontamento)

		if err := s.apontamentoRepo.Salvar(ctx, novoApontamento); err != nil {
			s.logger.ErrorContext(ctx, "falha ao salvar novo apontamento replicado", "funcionarioId", funcID, "erro", err)
			resultado.Falhas = append(resultado.Falhas, dto.DetalheFalha{FuncionarioID: funcID, Motivo: "Erro interno ao salvar novo apontamento."})
			resultado.Resumo.TotalFalha++
			continue
		}

		// Sucesso para este funcionário
		resultado.Sucessos = append(resultado.Sucessos, dto.DetalheSucesso{
			FuncionarioID:     funcID,
			NovoApontamentoID: novoApontamento.ID,
		})
		resultado.Resumo.TotalSucesso++
	}

	s.logger.InfoContext(ctx, "operação de replicação de apontamentos finalizada", "solicitados", resultado.Resumo.TotalSolicitado, "sucessos", resultado.Resumo.TotalSucesso, "falhas", resultado.Resumo.TotalFalha)
	return resultado, nil
}

// criarApontamentoAPartirDeTemplate é uma função helper para a lógica de template.
func criarApontamentoAPartirDeTemplate(template *pessoal.ApontamentoQuinzenal) *pessoal.ApontamentoQuinzenal {
	// Calcula o novo período
	novoInicio := template.PeriodoFim.AddDate(0, 0, 1)
	novoFim := novoInicio.AddDate(0, 0, 14) // Próximos 15 dias

	return &pessoal.ApontamentoQuinzenal{
		ID:            uuid.NewString(),
		FuncionarioID: template.FuncionarioID,
		ObraID:        template.ObraID, // Copia o contexto
		Diaria:        template.Diaria, // Mantém o valor da diária
		PeriodoInicio: novoInicio,
		PeriodoFim:    novoFim,
		// Zera os campos transacionais
		DiasTrabalhados:     0,
		Adicionais:          0,
		Descontos:           0,
		Adiantamentos:       0,
		ValorTotalCalculado: 0,
		Status:              pessoal.StatusApontamentoEmAberto, // Define o estado inicial
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

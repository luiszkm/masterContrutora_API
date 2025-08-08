package financeiro

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
)


// ContaReceberService encapsula a lógica de negócio para contas a receber
type ContaReceberService struct {
	contaReceberRepo financeiro.ContaReceberRepository
	eventBus         EventPublisher
	logger           *slog.Logger
}

func NovoContaReceberService(
	contaReceberRepo financeiro.ContaReceberRepository,
	eventBus EventPublisher,
	logger *slog.Logger,
) *ContaReceberService {
	return &ContaReceberService{
		contaReceberRepo: contaReceberRepo,
		eventBus:         eventBus,
		logger:           logger.With("service", "ContaReceber"),
	}
}

// CriarConta cria uma nova conta a receber
func (s *ContaReceberService) CriarConta(ctx context.Context, input dto.CriarContaReceberInput) (*dto.ContaReceberOutput, error) {
	const op = "service.financeiro.conta_receber.CriarConta"

	conta := &financeiro.ContaReceber{
		ID:                      uuid.NewString(),
		ObraID:                  input.ObraID,
		CronogramaRecebimentoID: input.CronogramaRecebimentoID,
		Cliente:                 input.Cliente,
		TipoContaReceber:        input.TipoContaReceber,
		Descricao:               input.Descricao,
		ValorOriginal:           input.ValorOriginal,
		ValorRecebido:           0,
		DataVencimento:          input.DataVencimento,
		Status:                  financeiro.StatusContaReceberPendente,
		NumeroDocumento:         input.NumeroDocumento,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	// Validar
	if err := conta.Validar(); err != nil {
		return nil, fmt.Errorf("%s: dados inválidos: %w", op, err)
	}

	// Salvar no banco
	if err := s.contaReceberRepo.Salvar(ctx, nil, conta); err != nil {
		return nil, fmt.Errorf("%s: falha ao salvar conta: %w", op, err)
	}

	// Publicar evento
	payload := events.ContaReceberCriadaPayload{
		ContaReceberID:          conta.ID,
		ObraID:                  conta.ObraID,
		CronogramaRecebimentoID: conta.CronogramaRecebimentoID,
		Cliente:                 conta.Cliente,
		TipoContaReceber:        conta.TipoContaReceber,
		Descricao:               conta.Descricao,
		ValorOriginal:           conta.ValorOriginal,
		DataVencimento:          conta.DataVencimento,
		NumeroDocumento:         conta.NumeroDocumento,
		UsuarioID:               "system", // TODO: pegar do contexto
	}

	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.ContaReceberCriada,
		Payload: payload,
	})

	s.logger.InfoContext(ctx, "conta a receber criada", "conta_id", conta.ID, "cliente", conta.Cliente)

	return s.toOutput(conta), nil
}

// RegistrarRecebimento registra um recebimento em uma conta
func (s *ContaReceberService) RegistrarRecebimento(ctx context.Context, contaID string, input dto.RegistrarRecebimentoContaInput) (*dto.ContaReceberOutput, error) {
	const op = "service.financeiro.conta_receber.RegistrarRecebimento"

	// Buscar conta
	conta, err := s.contaReceberRepo.BuscarPorID(ctx, contaID)
	if err != nil {
		return nil, fmt.Errorf("%s: conta não encontrada: %w", op, err)
	}

	// Registrar recebimento
	if err := conta.RegistrarRecebimento(input.Valor, input.FormaPagamento, input.Observacoes); err != nil {
		return nil, fmt.Errorf("%s: falha ao registrar recebimento: %w", op, err)
	}

	// Atualizar no banco
	if err := s.contaReceberRepo.Atualizar(ctx, conta); err != nil {
		return nil, fmt.Errorf("%s: falha ao atualizar conta: %w", op, err)
	}

	// Publicar evento
	payload := events.ContaReceberPagaPayload{
		ContaReceberID:     conta.ID,
		ObraID:             conta.ObraID,
		Cliente:            conta.Cliente,
		ValorRecebido:      input.Valor,
		ValorTotalRecebido: conta.ValorRecebido,
		ValorOriginal:      conta.ValorOriginal,
		ValorSaldo:         conta.ValorSaldo(),
		DataRecebimento:    time.Now(),
		FormaPagamento:     input.FormaPagamento,
		Status:             conta.Status,
		ContaBancariaID:    input.ContaBancariaID,
		UsuarioID:          "system", // TODO: pegar do contexto
	}

	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.ContaReceberPaga,
		Payload: payload,
	})

	s.logger.InfoContext(ctx, "recebimento registrado", 
		"conta_id", conta.ID, 
		"valor", input.Valor, 
		"status", conta.Status)

	return s.toOutput(conta), nil
}

// BuscarPorID busca uma conta por ID
func (s *ContaReceberService) BuscarPorID(ctx context.Context, id string) (*dto.ContaReceberOutput, error) {
	const op = "service.financeiro.conta_receber.BuscarPorID"

	conta, err := s.contaReceberRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.toOutput(conta), nil
}

// ListarPorObraID lista contas de uma obra
func (s *ContaReceberService) ListarPorObraID(ctx context.Context, obraID string) ([]*dto.ContaReceberOutput, error) {
	const op = "service.financeiro.conta_receber.ListarPorObraID"

	contas, err := s.contaReceberRepo.ListarPorObraID(ctx, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaReceberOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	return outputs, nil
}

// ListarVencidas lista contas vencidas
func (s *ContaReceberService) ListarVencidas(ctx context.Context) ([]*dto.ContaReceberOutput, error) {
	const op = "service.financeiro.conta_receber.ListarVencidas"

	contas, err := s.contaReceberRepo.ListarVencidas(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaReceberOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	return outputs, nil
}

// Listar lista contas com filtros e paginação
func (s *ContaReceberService) Listar(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ContaReceberOutput], error) {
	const op = "service.financeiro.conta_receber.Listar"

	contas, paginacao, err := s.contaReceberRepo.Listar(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaReceberOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	resposta := &common.RespostaPaginada[*dto.ContaReceberOutput]{
		Dados:     outputs,
		Paginacao: *paginacao,
	}

	return resposta, nil
}

// VerificarContasVencidas verifica e marca contas vencidas
func (s *ContaReceberService) VerificarContasVencidas(ctx context.Context) error {
	const op = "service.financeiro.conta_receber.VerificarContasVencidas"

	// Buscar contas que venceram
	contas, err := s.contaReceberRepo.ListarVencidas(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao buscar contas vencidas: %w", op, err)
	}

	for _, conta := range contas {
		if conta.EstaVencido() && conta.Status == financeiro.StatusContaReceberPendente {
			conta.MarcarComoVencido()
			
			if err := s.contaReceberRepo.Atualizar(ctx, conta); err != nil {
				s.logger.ErrorContext(ctx, "falha ao marcar conta como vencida", 
					"conta_id", conta.ID, "erro", err)
				continue
			}

			// Publicar evento de vencimento
			payload := events.ContaReceberVencidaPayload{
				ContaReceberID:   conta.ID,
				ObraID:           conta.ObraID,
				Cliente:          conta.Cliente,
				Descricao:        conta.Descricao,
				ValorOriginal:    conta.ValorOriginal,
				ValorSaldo:       conta.ValorSaldo(),
				DataVencimento:   conta.DataVencimento,
				DiasVencidos:     conta.DiasVencimento(),
				TipoContaReceber: conta.TipoContaReceber,
			}

			s.eventBus.Publicar(ctx, bus.Evento{
				Nome:    events.ContaReceberVencida,
				Payload: payload,
			})
		}
	}

	s.logger.InfoContext(ctx, "verificação de contas vencidas concluída", 
		"contas_processadas", len(contas))

	return nil
}

// ObterResumo obtém resumo das contas a receber
func (s *ContaReceberService) ObterResumo(ctx context.Context, filtros dto.FiltrosContaReceberInput) (*dto.ResumoContasReceberOutput, error) {
	const op = "service.financeiro.conta_receber.ObterResumo"

	// Converter filtros para common.ListarFiltros
	commonFiltros := common.ListarFiltros{
		TamanhoPagina: 1000000, // Buscar tudo para o resumo
		Pagina:        1,
	}

	contas, _, err := s.contaReceberRepo.Listar(ctx, commonFiltros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resumo := &dto.ResumoContasReceberOutput{}

	for _, conta := range contas {
		resumo.TotalContas++
		resumo.TotalValorOriginal += conta.ValorOriginal
		resumo.TotalValorRecebido += conta.ValorRecebido
		resumo.TotalValorSaldo += conta.ValorSaldo()

		switch conta.Status {
		case financeiro.StatusContaReceberPendente:
			resumo.ContasPendentes++
		case financeiro.StatusContaReceberVencido:
			resumo.ContasVencidas++
		case financeiro.StatusContaReceberRecebido:
			resumo.ContasRecebidas++
		}
	}

	if resumo.TotalValorOriginal > 0 {
		resumo.PercentualRecebimento = (resumo.TotalValorRecebido / resumo.TotalValorOriginal) * 100
	}

	return resumo, nil
}

// toOutput converte entidade para DTO de output
func (s *ContaReceberService) toOutput(conta *financeiro.ContaReceber) *dto.ContaReceberOutput {
	return &dto.ContaReceberOutput{
		ID:                      conta.ID,
		ObraID:                  conta.ObraID,
		CronogramaRecebimentoID: conta.CronogramaRecebimentoID,
		Cliente:                 conta.Cliente,
		TipoContaReceber:        conta.TipoContaReceber,
		Descricao:               conta.Descricao,
		ValorOriginal:           conta.ValorOriginal,
		ValorRecebido:           conta.ValorRecebido,
		ValorSaldo:              conta.ValorSaldo(),
		PercentualRecebido:      conta.PercentualRecebido(),
		DataVencimento:          conta.DataVencimento,
		DataRecebimento:         conta.DataRecebimento,
		Status:                  conta.Status,
		FormaPagamento:          conta.FormaPagamento,
		Observacoes:             conta.Observacoes,
		NumeroDocumento:         conta.NumeroDocumento,
		EstaVencido:             conta.EstaVencido(),
		DiasVencimento:          conta.DiasVencimento(),
		CreatedAt:               conta.CreatedAt,
		UpdatedAt:               conta.UpdatedAt,
	}
}
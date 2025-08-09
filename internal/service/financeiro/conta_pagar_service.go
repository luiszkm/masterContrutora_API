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

// ContaPagarService encapsula a lógica de negócio para contas a pagar
type ContaPagarService struct {
	contaPagarRepo financeiro.ContaPagarRepository
	eventBus       EventPublisher
	logger         *slog.Logger
}

func NovoContaPagarService(
	contaPagarRepo financeiro.ContaPagarRepository,
	eventBus EventPublisher,
	logger *slog.Logger,
) *ContaPagarService {
	return &ContaPagarService{
		contaPagarRepo: contaPagarRepo,
		eventBus:       eventBus,
		logger:         logger.With("service", "ContaPagar"),
	}
}

// CriarConta cria uma nova conta a pagar
func (s *ContaPagarService) CriarConta(ctx context.Context, input dto.CriarContaPagarInput) (*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.CriarConta"

	conta := &financeiro.ContaPagar{
		ID:              uuid.NewString(),
		FornecedorID:    input.FornecedorID,
		ObraID:          input.ObraID,
		OrcamentoID:     input.OrcamentoID,
		FornecedorNome:  input.FornecedorNome,
		TipoContaPagar:  input.TipoContaPagar,
		Descricao:       input.Descricao,
		ValorOriginal:   input.ValorOriginal,
		ValorPago:       0,
		DataVencimento:  input.DataVencimento,
		Status:          financeiro.StatusContaPagarPendente,
		NumeroDocumento: input.NumeroDocumento,
		NumeroCompraNF:  input.NumeroCompraNF,
		Observacoes:     input.Observacoes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Validar
	if err := conta.Validar(); err != nil {
		return nil, fmt.Errorf("%s: dados inválidos: %w", op, err)
	}

	// Salvar no banco
	if err := s.contaPagarRepo.Salvar(ctx, nil, conta); err != nil {
		return nil, fmt.Errorf("%s: falha ao salvar conta: %w", op, err)
	}

	// Publicar evento de criação de conta a pagar
	s.publicarEventoContaCriada(ctx, conta)

	s.logger.InfoContext(ctx, "conta a pagar criada", "conta_id", conta.ID, "fornecedor", conta.FornecedorNome)

	return s.toOutput(conta), nil
}

// CriarContaDeOrcamento cria uma conta a pagar a partir de um orçamento aprovado
func (s *ContaPagarService) CriarContaDeOrcamento(ctx context.Context, input dto.CriarContaPagarDeOrcamentoInput, orcamento interface{}) (*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.CriarContaDeOrcamento"

	// TODO: Implementar quando tivermos interface para buscar orçamento
	// Por enquanto, criamos uma conta genérica
	
	contaInput := dto.CriarContaPagarInput{
		OrcamentoID:     &input.OrcamentoID,
		FornecedorNome:  "Fornecedor do Orçamento", // TODO: buscar do orçamento
		TipoContaPagar:  "MATERIAL",
		Descricao:       "Conta gerada automaticamente do orçamento " + input.OrcamentoID,
		ValorOriginal:   1000.00, // TODO: buscar valor do orçamento
		DataVencimento:  input.DataVencimento,
		NumeroDocumento: input.NumeroDocumento,
		NumeroCompraNF:  input.NumeroCompraNF,
		Observacoes:     input.Observacoes,
	}

	return s.CriarConta(ctx, contaInput)
}

// RegistrarPagamento registra um pagamento em uma conta
func (s *ContaPagarService) RegistrarPagamento(ctx context.Context, contaID string, input dto.RegistrarPagamentoContaPagarInput) (*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.RegistrarPagamento"

	// Buscar conta
	conta, err := s.contaPagarRepo.BuscarPorID(ctx, contaID)
	if err != nil {
		return nil, fmt.Errorf("%s: conta não encontrada: %w", op, err)
	}

	// Registrar pagamento
	if err := conta.RegistrarPagamento(input.Valor, input.FormaPagamento, input.Observacoes); err != nil {
		return nil, fmt.Errorf("%s: falha ao registrar pagamento: %w", op, err)
	}

	// Atualizar no banco
	if err := s.contaPagarRepo.Atualizar(ctx, conta); err != nil {
		return nil, fmt.Errorf("%s: falha ao atualizar conta: %w", op, err)
	}

	// Publicar evento de pagamento
	s.publicarEventoPagamentoRealizado(ctx, conta, input.Valor, input.ContaBancariaID)

	s.logger.InfoContext(ctx, "pagamento registrado", 
		"conta_id", conta.ID, 
		"valor", input.Valor, 
		"status", conta.Status)

	return s.toOutput(conta), nil
}

// BuscarPorID busca uma conta por ID
func (s *ContaPagarService) BuscarPorID(ctx context.Context, id string) (*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.BuscarPorID"

	conta, err := s.contaPagarRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.toOutput(conta), nil
}

// ListarPorObraID lista contas de uma obra
func (s *ContaPagarService) ListarPorObraID(ctx context.Context, obraID string) ([]*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.ListarPorObraID"

	contas, err := s.contaPagarRepo.ListarPorObraID(ctx, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaPagarOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	return outputs, nil
}

// ListarPorFornecedorID lista contas de um fornecedor
func (s *ContaPagarService) ListarPorFornecedorID(ctx context.Context, fornecedorID string) ([]*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.ListarPorFornecedorID"

	contas, err := s.contaPagarRepo.ListarPorFornecedorID(ctx, fornecedorID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaPagarOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	return outputs, nil
}

// ListarVencidas lista contas vencidas
func (s *ContaPagarService) ListarVencidas(ctx context.Context) ([]*dto.ContaPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.ListarVencidas"

	contas, err := s.contaPagarRepo.ListarVencidas(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaPagarOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	return outputs, nil
}

// Listar lista contas com filtros e paginação
func (s *ContaPagarService) Listar(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ContaPagarOutput], error) {
	const op = "service.financeiro.conta_pagar.Listar"

	contas, paginacao, err := s.contaPagarRepo.Listar(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.ContaPagarOutput
	for _, conta := range contas {
		outputs = append(outputs, s.toOutput(conta))
	}

	resposta := &common.RespostaPaginada[*dto.ContaPagarOutput]{
		Dados:     outputs,
		Paginacao: *paginacao,
	}

	return resposta, nil
}

// VerificarContasVencidas verifica e marca contas vencidas
func (s *ContaPagarService) VerificarContasVencidas(ctx context.Context) error {
	const op = "service.financeiro.conta_pagar.VerificarContasVencidas"

	// Buscar contas que venceram
	contas, err := s.contaPagarRepo.ListarVencidas(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao buscar contas vencidas: %w", op, err)
	}

	for _, conta := range contas {
		if conta.EstaVencido() && conta.Status == financeiro.StatusContaPagarPendente {
			conta.MarcarComoVencido()
			
			if err := s.contaPagarRepo.Atualizar(ctx, conta); err != nil {
				s.logger.ErrorContext(ctx, "falha ao marcar conta como vencida", 
					"conta_id", conta.ID, "erro", err)
				continue
			}

			// Publicar evento de vencimento
			s.publicarEventoContaVencida(ctx, conta)
		}
	}

	s.logger.InfoContext(ctx, "verificação de contas vencidas concluída", 
		"contas_processadas", len(contas))

	return nil
}

// ObterResumo obtém resumo das contas a pagar
func (s *ContaPagarService) ObterResumo(ctx context.Context, filtros dto.FiltrosContaPagarInput) (*dto.ResumoContasPagarOutput, error) {
	const op = "service.financeiro.conta_pagar.ObterResumo"

	// Converter filtros para common.ListarFiltros
	commonFiltros := common.ListarFiltros{
		TamanhoPagina: 1000000, // Buscar tudo para o resumo
		Pagina:        1,
	}

	contas, _, err := s.contaPagarRepo.Listar(ctx, commonFiltros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resumo := &dto.ResumoContasPagarOutput{}

	for _, conta := range contas {
		resumo.TotalContas++
		resumo.TotalValorOriginal += conta.ValorOriginal
		resumo.TotalValorPago += conta.ValorPago
		resumo.TotalValorSaldo += conta.ValorSaldo()

		switch conta.Status {
		case financeiro.StatusContaPagarPendente:
			resumo.ContasPendentes++
		case financeiro.StatusContaPagarVencido:
			resumo.ContasVencidas++
		case financeiro.StatusContaPagarPago:
			resumo.ContasPagas++
		}
	}

	if resumo.TotalValorOriginal > 0 {
		resumo.PercentualPago = (resumo.TotalValorPago / resumo.TotalValorOriginal) * 100
	}

	return resumo, nil
}

// CancelarContaDeOrcamento cancela uma conta a pagar baseada no orçamento
func (s *ContaPagarService) CancelarContaDeOrcamento(ctx context.Context, orcamentoID string) error {
	const op = "service.financeiro.conta_pagar.CancelarContaDeOrcamento"

	// Buscar conta pelo orçamento ID de forma eficiente
	contas, err := s.contaPagarRepo.ListarPorOrcamentoID(ctx, orcamentoID)
	if err != nil {
		return fmt.Errorf("%s: falha ao buscar contas por orçamento: %w", op, err)
	}

	if len(contas) == 0 {
		s.logger.WarnContext(ctx, "nenhuma conta a pagar encontrada para o orçamento", 
			"orcamento_id", orcamentoID)
		return nil // Não é erro se não existe conta
	}

	// Pegar a primeira conta (deve ser única por orçamento)
	contaEncontrada := contas[0]

	// Verificar se a conta pode ser cancelada (não pode ter pagamentos)
	if contaEncontrada.ValorPago > 0 {
		return fmt.Errorf("%s: conta não pode ser cancelada pois já possui pagamentos (valor pago: %.2f)", 
			op, contaEncontrada.ValorPago)
	}

	// Marcar como cancelada
	contaEncontrada.Status = financeiro.StatusContaPagarCancelado
	contaEncontrada.Observacoes = func() *string { 
		obs := "Cancelada automaticamente devido ao cancelamento do orçamento"
		if contaEncontrada.Observacoes != nil {
			obs = *contaEncontrada.Observacoes + " | " + obs
		}
		return &obs 
	}()
	contaEncontrada.UpdatedAt = time.Now()

	// Salvar alteração
	if err := s.contaPagarRepo.Atualizar(ctx, contaEncontrada); err != nil {
		return fmt.Errorf("%s: falha ao cancelar conta: %w", op, err)
	}

	// Publicar evento de cancelamento
	s.publicarEventoContaCancelada(ctx, contaEncontrada, orcamentoID)

	s.logger.InfoContext(ctx, "conta a pagar cancelada devido ao cancelamento do orçamento", 
		"conta_id", contaEncontrada.ID,
		"orcamento_id", orcamentoID,
		"valor_original", contaEncontrada.ValorOriginal)

	return nil
}

// toOutput converte entidade para DTO de output
func (s *ContaPagarService) toOutput(conta *financeiro.ContaPagar) *dto.ContaPagarOutput {
	return &dto.ContaPagarOutput{
		ID:              conta.ID,
		FornecedorID:    conta.FornecedorID,
		ObraID:          conta.ObraID,
		OrcamentoID:     conta.OrcamentoID,
		FornecedorNome:  conta.FornecedorNome,
		TipoContaPagar:  conta.TipoContaPagar,
		Descricao:       conta.Descricao,
		ValorOriginal:   conta.ValorOriginal,
		ValorPago:       conta.ValorPago,
		ValorSaldo:      conta.ValorSaldo(),
		PercentualPago:  conta.PercentualPago(),
		DataVencimento:  conta.DataVencimento,
		DataPagamento:   conta.DataPagamento,
		Status:          conta.Status,
		FormaPagamento:  conta.FormaPagamento,
		Observacoes:     conta.Observacoes,
		NumeroDocumento: conta.NumeroDocumento,
		NumeroCompraNF:  conta.NumeroCompraNF,
		EstaVencido:     conta.EstaVencido(),
		DiasVencimento:  conta.DiasVencimento(),
		CreatedAt:       conta.CreatedAt,
		UpdatedAt:       conta.UpdatedAt,
	}
}

// Métodos auxiliares para publicar eventos
func (s *ContaPagarService) publicarEventoContaCriada(ctx context.Context, conta *financeiro.ContaPagar) {
	// TODO: Definir evento para conta a pagar criada se necessário
	s.logger.InfoContext(ctx, "conta a pagar criada", "conta_id", conta.ID)
}

func (s *ContaPagarService) publicarEventoPagamentoRealizado(ctx context.Context, conta *financeiro.ContaPagar, valorPago float64, contaBancariaID *string) {
	// Publicar evento de movimentação financeira (saída)
	payload := events.MovimentacaoFinanceiraRegistradaPayload{
		MovimentacaoID:   uuid.NewString(),
		ContaBancariaID:  *contaBancariaID, // TODO: tratar caso seja nil
		TipoMovimentacao: "SAIDA",
		Valor:            valorPago,
		DataMovimentacao: time.Now(),
		DataCompetencia:  time.Now(),
		Descricao:        fmt.Sprintf("Pagamento - %s", conta.Descricao),
		DocumentoID:      &conta.ID,
		DocumentoTipo:    func() *string { s := "CONTA_PAGAR"; return &s }(),
		Status:           "REALIZADO",
		UsuarioID:        "system", // TODO: pegar do contexto
	}

	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.MovimentacaoFinanceiraRegistrada,
		Payload: payload,
	})
}

func (s *ContaPagarService) publicarEventoContaVencida(ctx context.Context, conta *financeiro.ContaPagar) {
	// TODO: Definir evento para conta a pagar vencida se necessário
	s.logger.WarnContext(ctx, "conta a pagar vencida", 
		"conta_id", conta.ID, 
		"fornecedor", conta.FornecedorNome,
		"dias_vencidos", conta.DiasVencimento())
}

func (s *ContaPagarService) publicarEventoContaCancelada(ctx context.Context, conta *financeiro.ContaPagar, orcamentoID string) {
	// TODO: Definir evento específico para conta cancelada se necessário
	s.logger.InfoContext(ctx, "conta a pagar cancelada", 
		"conta_id", conta.ID,
		"orcamento_id", orcamentoID,
		"fornecedor", conta.FornecedorNome,
		"valor_original", conta.ValorOriginal)
}
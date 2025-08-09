// file: internal/handler/http/financeiro/handler.go
package financeiro

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	financeiro_service "github.com/luiszkm/masterCostrutora/internal/service/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
)

type Service interface {
	RegistrarPagamento(ctx context.Context, input dto.RegistrarPagamentoInput) (*financeiro.RegistroDePagamento, error)
	RegistrarPagamentosEmLote(ctx context.Context, input dto.RegistrarPagamentoEmLoteInput) (*dto.ResultadoExecucaoLote, error)
	ListarPagamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*financeiro.RegistroDePagamento], error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NovoFinanceiroHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

type registrarPagamentoRequest struct {
	FuncionarioID     string  `json:"funcionarioId"`
	ObraID            string  `json:"obraId"`
	PeriodoReferencia string  `json:"periodoReferencia"`
	ValorCalculado    float64 `json:"valorCalculado"`
	ContaBancariaID   string  `json:"contaBancariaId"`
}

func (h *Handler) HandleRegistrarPagamento(w http.ResponseWriter, r *http.Request) {
	var req registrarPagamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.RegistrarPagamentoInput{
		FuncionarioID:     req.FuncionarioID,
		ObraID:            req.ObraID,
		PeriodoReferencia: req.PeriodoReferencia,
		ValorCalculado:    req.ValorCalculado,
		ContaBancariaID:   req.ContaBancariaID,
	}

	pagamento, err := h.service.RegistrarPagamento(r.Context(), input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "RECURSO_NAO_ENCONTRADO", "Funcionário ou Obra não encontrado(a)", http.StatusNotFound)
			return
		}
		if errors.Is(err, financeiro_service.ErrFuncionarioInativo) {
			web.RespondError(w, r, "REGRA_NEGOCIO_VIOLADA", err.Error(), http.StatusUnprocessableEntity) // 422
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao registrar pagamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao registrar pagamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, pagamento, http.StatusCreated)
}
func (h *Handler) HandlePagamentoDeApontamentoRealizado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.PagamentoApontamentoRealizadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de pagamento inválido", "evento", evento.Nome)
		return
	}

	h.logger.Info("EVENTO RECEBIDO PELO CONTEXTO FINANCEIRO!", "funcionario_id", payload.FuncionarioID, "valor", payload.ValorCalculado)

	input := dto.RegistrarPagamentoInput{
		FuncionarioID:     payload.FuncionarioID,
		ObraID:            payload.ObraID,
		PeriodoReferencia: payload.PeriodoReferencia,
		ValorCalculado:    payload.ValorCalculado,
		ContaBancariaID:   payload.ContaBancariaID,
	}

	// Chama o próprio serviço para criar o registro de pagamento.
	if _, err := h.service.RegistrarPagamento(ctx, input); err != nil {
		h.logger.ErrorContext(ctx, "falha ao processar evento de pagamento", "erro", err)
		// Aqui entraria a lógica de retentativas e DLQ do ADR-007.
	}
}

func (h *Handler) HandleRegistrarPagamentosEmLote(w http.ResponseWriter, r *http.Request) {
	var input dto.RegistrarPagamentoEmLoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload da requisição é inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas do payload
	if len(input.ApontamentoIDs) == 0 || input.ContaBancariaID == "" || input.DataDeEfetivacao == "" {
		web.RespondError(w, r, "DADOS_OBRIGATORIOS", "apontamentoIds, contaBancariaId e dataDeEfetivacao são obrigatórios.", http.StatusBadRequest)
		return
	}

	resultado, err := h.service.RegistrarPagamentosEmLote(r.Context(), input)
	if err != nil {
		// Este erro só deve ocorrer para falhas inesperadas na camada de serviço (ex: data inválida)
		h.logger.ErrorContext(r.Context(), "falha na execução de pagamentos em lote", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", err.Error(), http.StatusInternalServerError)
		return
	}

	// Conforme a V4, a resposta de sucesso parcial deve ser 207 Multi-Status
	web.Respond(w, r, resultado, http.StatusMultiStatus)
}

func (h *Handler) HandleListarPagamentos(w http.ResponseWriter, r *http.Request) {
	filtros := web.ParseFiltros(r)
	
	pagamentos, err := h.service.ListarPagamentos(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar pagamentos", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar pagamentos", http.StatusInternalServerError)
		return
	}
	
	web.Respond(w, r, pagamentos, http.StatusOK)
}

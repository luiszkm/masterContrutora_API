package obras

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

// EventPublisher interface para publicar eventos
type EventPublisher interface {
	Publicar(ctx context.Context, evento bus.Evento)
}

// CronogramaService encapsula a lógica de negócio para cronogramas de recebimento
type CronogramaService struct {
	cronogramaRepo obras.CronogramaRecebimentoRepository
	obraRepo       obras.ObrasRepository
	eventBus       EventPublisher
	logger         *slog.Logger
	dbpool         *pgxpool.Pool
}

func NovoCronogramaService(
	cronogramaRepo obras.CronogramaRecebimentoRepository,
	obraRepo obras.ObrasRepository,
	eventBus EventPublisher,
	logger *slog.Logger,
	dbpool *pgxpool.Pool,
) *CronogramaService {
	return &CronogramaService{
		cronogramaRepo: cronogramaRepo,
		obraRepo:       obraRepo,
		eventBus:       eventBus,
		logger:         logger.With("service", "CronogramaRecebimento"),
		dbpool:         dbpool,
	}
}

// CriarCronograma cria um único cronograma de recebimento
func (s *CronogramaService) CriarCronograma(ctx context.Context, input dto.CriarCronogramaRecebimentoInput) (*dto.CronogramaRecebimentoOutput, error) {
	const op = "service.obras.cronograma.CriarCronograma"

	// Validar se a obra existe
	obra, err := s.obraRepo.BuscarPorID(ctx, input.ObraID)
	if err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	// Criar cronograma
	cronograma := &obras.CronogramaRecebimento{
		ID:             uuid.NewString(),
		ObraID:         input.ObraID,
		NumeroEtapa:    input.NumeroEtapa,
		DescricaoEtapa: input.DescricaoEtapa,
		ValorPrevisto:  input.ValorPrevisto,
		DataVencimento: input.DataVencimento,
		Status:         obras.StatusRecebimentoPendente,
		ValorRecebido:  0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Validar
	if err := cronograma.Validar(); err != nil {
		return nil, fmt.Errorf("%s: dados inválidos: %w", op, err)
	}

	// Salvar no banco
	if err := s.cronogramaRepo.Salvar(ctx, s.dbpool, cronograma); err != nil {
		return nil, fmt.Errorf("%s: falha ao salvar cronograma: %w", op, err)
	}

	// Publicar evento
	payload := events.CronogramaRecebimentoCriadoPayload{
		ObraID:             cronograma.ObraID,
		ObraNome:           obra.Nome,
		Cliente:            obra.Cliente,
		CronogramasIds:     []string{cronograma.ID},
		ValorTotalPrevisto: cronograma.ValorPrevisto,
		QuantidadeEtapas:   1,
		PrimeiroVencimento: cronograma.DataVencimento,
		UsuarioID:          "system", // TODO: pegar do contexto
	}

	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.CronogramaRecebimentoCriado,
		Payload: payload,
	})

	s.logger.InfoContext(ctx, "cronograma de recebimento criado", "cronograma_id", cronograma.ID, "obra_id", obra.ID)

	return s.toOutput(cronograma), nil
}

// CriarCronogramaEmLote cria múltiplos cronogramas de uma vez
func (s *CronogramaService) CriarCronogramaEmLote(ctx context.Context, input dto.CriarCronogramaEmLoteInput) ([]*dto.CronogramaRecebimentoOutput, error) {
	const op = "service.obras.cronograma.CriarCronogramaEmLote"

	// Validar se a obra existe
	obra, err := s.obraRepo.BuscarPorID(ctx, input.ObraID)
	if err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	// Iniciar transação
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Se solicitado, remover cronogramas existentes
	if input.SubstituirExistente {
		cronogramasExistentes, err := s.cronogramaRepo.ListarPorObraID(ctx, input.ObraID)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao listar cronogramas existentes: %w", op, err)
		}

		for _, cronograma := range cronogramasExistentes {
			if err := s.cronogramaRepo.Deletar(ctx, cronograma.ID); err != nil {
				return nil, fmt.Errorf("%s: falha ao deletar cronograma existente: %w", op, err)
			}
		}
	}

	// Criar cronogramas
	var cronogramas []*obras.CronogramaRecebimento
	var cronogramasIds []string
	var valorTotalPrevisto float64
	var primeiroVencimento time.Time

	for i, inputCronograma := range input.Cronogramas {
		cronograma := &obras.CronogramaRecebimento{
			ID:             uuid.NewString(),
			ObraID:         input.ObraID,
			NumeroEtapa:    inputCronograma.NumeroEtapa,
			DescricaoEtapa: inputCronograma.DescricaoEtapa,
			ValorPrevisto:  inputCronograma.ValorPrevisto,
			DataVencimento: inputCronograma.DataVencimento,
			Status:         obras.StatusRecebimentoPendente,
			ValorRecebido:  0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Validar
		if err := cronograma.Validar(); err != nil {
			return nil, fmt.Errorf("%s: dados inválidos no cronograma %d: %w", op, i+1, err)
		}

		cronogramas = append(cronogramas, cronograma)
		cronogramasIds = append(cronogramasIds, cronograma.ID)
		valorTotalPrevisto += cronograma.ValorPrevisto

		// Definir primeiro vencimento
		if i == 0 || cronograma.DataVencimento.Before(primeiroVencimento) {
			primeiroVencimento = cronograma.DataVencimento
		}
	}

	// Salvar em lote
	if err := s.cronogramaRepo.SalvarMuitos(ctx, tx, cronogramas); err != nil {
		return nil, fmt.Errorf("%s: falha ao salvar cronogramas: %w", op, err)
	}

	// Commit da transação
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: falha ao fazer commit: %w", op, err)
	}

	// Publicar evento
	payload := events.CronogramaRecebimentoCriadoPayload{
		ObraID:             input.ObraID,
		ObraNome:           obra.Nome,
		Cliente:            obra.Cliente,
		CronogramasIds:     cronogramasIds,
		ValorTotalPrevisto: valorTotalPrevisto,
		QuantidadeEtapas:   len(cronogramas),
		PrimeiroVencimento: primeiroVencimento,
		UsuarioID:          "system", // TODO: pegar do contexto
	}

	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.CronogramaRecebimentoCriado,
		Payload: payload,
	})

	s.logger.InfoContext(ctx, "cronogramas criados em lote", 
		"obra_id", obra.ID, 
		"quantidade", len(cronogramas), 
		"valor_total", valorTotalPrevisto)

	// Converter para output
	var outputs []*dto.CronogramaRecebimentoOutput
	for _, cronograma := range cronogramas {
		outputs = append(outputs, s.toOutput(cronograma))
	}

	return outputs, nil
}

// RegistrarRecebimento registra um recebimento em um cronograma
func (s *CronogramaService) RegistrarRecebimento(ctx context.Context, cronogramaID string, input dto.RegistrarRecebimentoInput) (*dto.CronogramaRecebimentoOutput, error) {
	const op = "service.obras.cronograma.RegistrarRecebimento"

	// Buscar cronograma
	cronograma, err := s.cronogramaRepo.BuscarPorID(ctx, cronogramaID)
	if err != nil {
		return nil, fmt.Errorf("%s: cronograma não encontrado: %w", op, err)
	}

	// Buscar obra para dados do evento
	obra, err := s.obraRepo.BuscarPorID(ctx, cronograma.ObraID)
	if err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	// Registrar recebimento
	if err := cronograma.RegistrarRecebimento(input.Valor, input.Observacoes); err != nil {
		return nil, fmt.Errorf("%s: falha ao registrar recebimento: %w", op, err)
	}

	// Atualizar no banco
	if err := s.cronogramaRepo.Atualizar(ctx, cronograma); err != nil {
		return nil, fmt.Errorf("%s: falha ao atualizar cronograma: %w", op, err)
	}

	// Atualizar valor recebido na obra
	obra.ValorRecebido += input.Valor
	if err := s.obraRepo.Atualizar(ctx, obra); err != nil {
		s.logger.WarnContext(ctx, "falha ao atualizar valor recebido na obra", "obra_id", obra.ID, "erro", err)
	}

	// Publicar evento de recebimento
	payload := events.RecebimentoRealizadoPayload{
		CronogramaRecebimentoID: &cronograma.ID,
		ObraID:                  &cronograma.ObraID,
		ObraNome:                &obra.Nome,
		Cliente:                 obra.Cliente,
		ValorRecebido:           input.Valor,
		DataRecebimento:         time.Now(),
		Descricao:               cronograma.DescricaoEtapa,
		UsuarioID:               "system", // TODO: pegar do contexto
	}

	s.eventBus.Publicar(ctx, bus.Evento{
		Nome:    events.RecebimentoRealizado,
		Payload: payload,
	})

	s.logger.InfoContext(ctx, "recebimento registrado", 
		"cronograma_id", cronograma.ID, 
		"valor", input.Valor, 
		"status", cronograma.Status)

	return s.toOutput(cronograma), nil
}

// ListarPorObraID lista cronogramas de uma obra
func (s *CronogramaService) ListarPorObraID(ctx context.Context, obraID string) ([]*dto.CronogramaRecebimentoOutput, error) {
	const op = "service.obras.cronograma.ListarPorObraID"

	cronogramas, err := s.cronogramaRepo.ListarPorObraID(ctx, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var outputs []*dto.CronogramaRecebimentoOutput
	for _, cronograma := range cronogramas {
		outputs = append(outputs, s.toOutput(cronograma))
	}

	return outputs, nil
}

// BuscarPorID busca um cronograma por ID
func (s *CronogramaService) BuscarPorID(ctx context.Context, id string) (*dto.CronogramaRecebimentoOutput, error) {
	const op = "service.obras.cronograma.BuscarPorID"

	cronograma, err := s.cronogramaRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.toOutput(cronograma), nil
}

// VerificarEtapasVencidas verifica e marca etapas vencidas
func (s *CronogramaService) VerificarEtapasVencidas(ctx context.Context) error {
	const op = "service.obras.cronograma.VerificarEtapasVencidas"

	// Buscar etapas que venceram até hoje
	hoje := time.Now()
	inicio := hoje.AddDate(0, 0, -30) // Últimos 30 dias para não processar tudo sempre

	cronogramasVencidos, err := s.cronogramaRepo.ListarVencidosPorPeriodo(ctx, inicio, hoje)
	if err != nil {
		return fmt.Errorf("%s: falha ao buscar cronogramas vencidos: %w", op, err)
	}

	for _, cronograma := range cronogramasVencidos {
		if cronograma.EstaVencido() && cronograma.Status == obras.StatusRecebimentoPendente {
			cronograma.MarcarComoVencido()
			
			if err := s.cronogramaRepo.Atualizar(ctx, cronograma); err != nil {
				s.logger.ErrorContext(ctx, "falha ao marcar cronograma como vencido", 
					"cronograma_id", cronograma.ID, "erro", err)
				continue
			}

			// Buscar dados da obra para o evento
			obra, err := s.obraRepo.BuscarPorID(ctx, cronograma.ObraID)
			if err != nil {
				s.logger.ErrorContext(ctx, "falha ao buscar obra para evento vencido", 
					"obra_id", cronograma.ObraID, "erro", err)
				continue
			}

			// Publicar evento de vencimento
			payload := events.EtapaRecebimentoVencidaPayload{
				CronogramaRecebimentoID: cronograma.ID,
				ObraID:                  cronograma.ObraID,
				ObraNome:                obra.Nome,
				Cliente:                 obra.Cliente,
				NumeroEtapa:             cronograma.NumeroEtapa,
				DescricaoEtapa:          cronograma.DescricaoEtapa,
				ValorPrevisto:           cronograma.ValorPrevisto,
				ValorSaldo:              cronograma.ValorSaldo(),
				DataVencimento:          cronograma.DataVencimento,
				DiasVencidos:            int(hoje.Sub(cronograma.DataVencimento).Hours() / 24),
			}

			s.eventBus.Publicar(ctx, bus.Evento{
				Nome:    events.EtapaRecebimentoVencida,
				Payload: payload,
			})
		}
	}

	s.logger.InfoContext(ctx, "verificação de etapas vencidas concluída", 
		"cronogramas_processados", len(cronogramasVencidos))

	return nil
}

// toOutput converte entidade para DTO de output
func (s *CronogramaService) toOutput(cronograma *obras.CronogramaRecebimento) *dto.CronogramaRecebimentoOutput {
	return &dto.CronogramaRecebimentoOutput{
		ID:                     cronograma.ID,
		ObraID:                 cronograma.ObraID,
		NumeroEtapa:            cronograma.NumeroEtapa,
		DescricaoEtapa:         cronograma.DescricaoEtapa,
		ValorPrevisto:          cronograma.ValorPrevisto,
		DataVencimento:         cronograma.DataVencimento,
		Status:                 cronograma.Status,
		DataRecebimento:        cronograma.DataRecebimento,
		ValorRecebido:          cronograma.ValorRecebido,
		ValorSaldo:             cronograma.ValorSaldo(),
		PercentualRecebido:     cronograma.PercentualRecebido(),
		ObservacoesRecebimento: cronograma.ObservacoesRecebimento,
		EstaVencido:            cronograma.EstaVencido(),
		CreatedAt:              cronograma.CreatedAt,
		UpdatedAt:              cronograma.UpdatedAt,
	}
}
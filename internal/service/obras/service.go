// file: internal/service/obras/service.go
package obras

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid" // Usaremos UUID para os IDs.
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto" // Importa o pacote de DTO
	// Importa o pacote de DTO
)

type Querier interface {
	BuscarDashboardPorID(ctx context.Context, id string) (*dto.ObraDashboard, error)
}
type EtapaRepository interface {
	Salvar(ctx context.Context, etapa *obras.Etapa) error
	BuscarPorID(ctx context.Context, etapaID string) (*obras.Etapa, error) // NOVO
	Atualizar(ctx context.Context, etapa *obras.Etapa) error               // NOVO
}

// Service encapsula a lógica de negócio para o contexto de Obras.
type Service struct {
	obraRepo  obras.Repository
	etapaRepo EtapaRepository
	querier   Querier
	logger    *slog.Logger
}

// NovoServico é o construtor para o serviço de obras.
func NovoServico(obraRepo obras.Repository, etapaRepo EtapaRepository, querier Querier, logger *slog.Logger) *Service {
	return &Service{
		obraRepo:  obraRepo,
		etapaRepo: etapaRepo,
		querier:   querier,
		logger:    logger,
	}
}

// CriarNovaObra é o caso de uso para registrar uma nova construção.
func (s *Service) CriarNovaObra(ctx context.Context, input dto.CriarNovaObraInput) (*obras.Obra, error) {
	const op = "service.obras.CriarNovaObra"

	// Validação básica de entrada
	if input.Nome == "" || input.Cliente == "" || input.Endereco == "" {
		return nil, fmt.Errorf("%s: nome, cliente e endereço são obrigatórios", op)
	}

	dataInicio, err := time.Parse("2006-01-02", input.DataInicio)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de início inválido: %w", op, err)
	}

	novaObra := &obras.Obra{
		ID:         uuid.NewString(),
		Nome:       input.Nome,
		Cliente:    input.Cliente,
		Endereco:   input.Endereco,
		DataInicio: dataInicio,
		Status:     obras.StatusEmPlanejamento, // Status inicial padrão
	}

	// Delega a persistência para o repositório
	if err := s.obraRepo.Salvar(ctx, novaObra); err != nil {
		// Adiciona contexto ao erro original usando %w
		return nil, fmt.Errorf("%s: falha ao salvar nova obra: %w", op, err)
	}

	s.logger.InfoContext(ctx, "nova obra criada com sucesso", "obra_id", novaObra.ID, "obra_nome", novaObra.Nome)

	return novaObra, nil
}
func (s *Service) BuscarDashboard(ctx context.Context, id string) (*dto.ObraDashboard, error) {
	const op = "service.obras.BuscarDashboard"

	// A lógica agora usa a dependência correta.
	dashboard, err := s.querier.BuscarDashboardPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dashboard, nil
}
func (s *Service) AdicionarEtapa(ctx context.Context, obraID string, input dto.AdicionarEtapaInput) (*obras.Etapa, error) {
	const op = "service.obras.AdicionarEtapa"

	// TODO: Antes de adicionar, poderíamos validar se a Obra com o 'obraID' realmente existe
	// usando s.obraRepo.BuscarPorID(ctx, obraID).

	inicio, err := time.Parse("2006-01-02", input.DataInicioPrevista)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de início inválido: %w", op, err)
	}
	fim, err := time.Parse("2006-01-02", input.DataFimPrevista)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de fim inválido: %w", op, err)
	}

	novaEtapa := &obras.Etapa{
		ID:                 uuid.NewString(),
		ObraID:             obraID,
		Nome:               input.Nome,
		DataInicioPrevista: inicio,
		DataFimPrevista:    fim,
		Status:             "Pendente", // Status inicial padrão
	}

	if err := s.etapaRepo.Salvar(ctx, novaEtapa); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "etapa adicionada com sucesso", "etapa_id", novaEtapa.ID, "obra_id", obraID)

	return novaEtapa, nil
}

func (s *Service) AtualizarStatusEtapa(ctx context.Context, etapaID string, input dto.AtualizarStatusEtapaInput) (*obras.Etapa, error) {
	const op = "service.obras.AtualizarStatusEtapa"

	// 1. Buscar a etapa que queremos modificar
	etapa, err := s.etapaRepo.BuscarPorID(ctx, etapaID)
	if err != nil {
		// Propaga o erro "não encontrado" ou qualquer outro erro do repositório
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// 2. Aplicar a lógica de negócio/validação
	// TODO: Adicionar validações de transição de status.
	// Por exemplo, uma etapa não pode ir de "Concluída" de volta para "Pendente".
	etapa.Status = obras.StatusEtapa(input.Status)

	// 3. Salvar a etapa atualizada
	if err := s.etapaRepo.Atualizar(ctx, etapa); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "status da etapa atualizado", "etapa_id", etapa.ID, "novo_status", etapa.Status)
	return etapa, nil
}

// file: internal/service/obras/service.go
package obras

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"                                     // Usaremos UUID para os IDs.
	"github.com/luiszkm/masterCostrutora/internal/domain/common" // Importa o pacote de filtros e paginação
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto" // Importa o pacote de DTO
	// Importa o pacote de DTO
)

type ObrasQuerier interface {
	BuscarDashboardPorID(ctx context.Context, id string) (*dto.ObraDashboard, error)
	ListarObras(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ObraListItemDTO, *common.PaginacaoInfo, error)
}
type PessoalFinder interface {
	BuscarPorID(ctx context.Context, funcionarioID string) (*pessoal.Funcionario, error)
}

// Service encapsula a lógica de negócio para o contexto de Obras.
type Service struct {
	obraRepo      obras.ObrasRepository
	etapaRepo     obras.EtapaRepository
	alocacaoRepo  obras.AlocacaoRepository
	pessoalFinder PessoalFinder
	obrasQuerier  ObrasQuerier
	logger        *slog.Logger
}

// ListarObras implements obras.Service.
func (s *Service) ListarObras(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ObraListItemDTO], error) {
	panic("unimplemented")
}

func NovoServico(obraRepo obras.ObrasRepository, etapaRepo obras.EtapaRepository,
	alocacaoRepo obras.AlocacaoRepository, pessoalFinder PessoalFinder, obrasQuerier ObrasQuerier, logger *slog.Logger) *Service {
	return &Service{
		alocacaoRepo:  alocacaoRepo,
		pessoalFinder: pessoalFinder,
		obraRepo:      obraRepo,
		etapaRepo:     etapaRepo,
		obrasQuerier:  obrasQuerier,
		logger:        logger,
	}
}

// func (s *Service) ListarObras(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ObraListItemDTO], error) {
// 	const op = "service.obras.ListarObras"

// 	obras, paginacao, err := s.obrasQuerier.ListarObras(ctx, filtros)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return &common.RespostaPaginada[*dto.ObraListItemDTO]{
// 		Dados:     obras,
// 		Paginacao: *paginacao,
// 	}, nil
// }

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
	dashboard, err := s.obrasQuerier.BuscarDashboardPorID(ctx, id)
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
func (s *Service) AlocarFuncionario(ctx context.Context, obraID string, input dto.AlocarFuncionarioInput) (*obras.Alocacao, error) {
	const op = "service.obras.AlocarFuncionario"

	// --- Lógica de Negócio e Validação ---
	// 1. Verifica se a obra existe (podemos usar o querier para isso)
	_, err := s.obrasQuerier.BuscarDashboardPorID(ctx, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}
	// 2. Verifica se o funcionário existe (colaboração entre contextos!)
	_, err = s.pessoalFinder.BuscarPorID(ctx, input.FuncionarioID)
	if err != nil {
		return nil, fmt.Errorf("%s: funcionário não encontrado: %w", op, err)
	}

	inicio, err := time.Parse("2006-01-02", input.DataInicioAlocacao)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data inválido: %w", op, err)
	}

	novaAlocacao := &obras.Alocacao{
		ID:                 uuid.NewString(),
		ObraID:             obraID,
		FuncionarioID:      input.FuncionarioID,
		DataInicioAlocacao: inicio,
		// DataFimAlocacao fica nula por padrão
	}

	if err := s.alocacaoRepo.Salvar(ctx, novaAlocacao); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "funcionário alocado com sucesso", "alocacao_id", novaAlocacao.ID)
	return novaAlocacao, nil
}
func (s *Service) DeletarObra(ctx context.Context, id string) error {
	const op = "service.obras.DeletarObra"
	// TODO: Adicionar lógica de negócio aqui. Ex: não se pode deletar uma obra com pagamentos pendentes.
	if err := s.obraRepo.Deletar(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "obra movida para a lixeira", "obra_id", id)
	return nil
}

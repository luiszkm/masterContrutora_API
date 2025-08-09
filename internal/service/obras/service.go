// file: internal/service/obras/service.go
package obras

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid" // Usaremos UUID para os IDs.
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common" // Importa o pacote de filtros e paginação
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto" // Importa o pacote de DTO
	// Importa o pacote de DTO
)

type ObrasQuerier interface {
	BuscarDashboardPorID(ctx context.Context, id string) (*dto.ObraDashboard, error)
	ListarObras(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ObraListItemDTO, *common.PaginacaoInfo, error)
	BuscarDetalhesPorID(ctx context.Context, obraID string) (*dto.ObraDetalhadaDTO, error)
}
type PessoalFinder interface {
	BuscarPorID(ctx context.Context, funcionarioID string) (*pessoal.Funcionario, error)
}

// Service encapsula a lógica de negócio para o contexto de Obras.
type Service struct {
	obraRepo        obras.ObrasRepository
	etapaRepo       obras.EtapaRepository
	alocacaoRepo    obras.AlocacaoRepository
	etapaPadraoRepo obras.EtapaPadraoRepository
	pessoalFinder   PessoalFinder
	obrasQuerier    ObrasQuerier
	logger          *slog.Logger
	dbpool          *pgxpool.Pool // NOVO

}

func NovoServico(obraRepo obras.ObrasRepository, etapaRepo obras.EtapaRepository,
	etapaPadraoRepo obras.EtapaPadraoRepository,
	alocacaoRepo obras.AlocacaoRepository, pessoalFinder PessoalFinder, obrasQuerier ObrasQuerier,
	logger *slog.Logger, dbpool *pgxpool.Pool) *Service {
	return &Service{
		alocacaoRepo:    alocacaoRepo,
		pessoalFinder:   pessoalFinder,
		obraRepo:        obraRepo,
		etapaRepo:       etapaRepo,
		etapaPadraoRepo: etapaPadraoRepo,
		obrasQuerier:    obrasQuerier,
		logger:          logger,
		dbpool:          dbpool, // NOVO
	}
}

// ListarObras implements obras.Service.
func (s *Service) ListarObras(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ObraListItemDTO], error) {
	const op = "service.obras.ListarObras"

	obras, paginacao, err := s.obrasQuerier.ListarObras(ctx, filtros)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &common.RespostaPaginada[*dto.ObraListItemDTO]{
		Dados:     obras,
		Paginacao: *paginacao,
	}, nil
}

// CriarNovaObra é o caso de uso para registrar uma nova construção.
func (s *Service) CriarNovaObra(ctx context.Context, input dto.CriarNovaObraInput) (*obras.Obra, error) {
	const op = "service.obras.CriarNovaObra"

	// Inicia a transação
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx) // Garante o rollback em caso de erro

	// 1. Cria a entidade Obra principal
	dataInicio, err := time.Parse("2006-01-02", input.DataInicio)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data inválido: %w", op, err)
	}
	novaObra := &obras.Obra{
		ID:         uuid.NewString(),
		Nome:       input.Nome,
		Cliente:    input.Cliente,
		Endereco:   input.Endereco,
		DataInicio: dataInicio,
		Status:     obras.StatusEmPlanejamento,
	}
	if err := s.obraRepo.Salvar(ctx, tx, novaObra); err != nil {
		return nil, fmt.Errorf("%s: falha ao salvar nova obra: %w", op, err)
	}

	// 2. Busca todas as etapas padrão do catálogo
	etapasPadrao, err := s.etapaPadraoRepo.ListarTodas(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao buscar catálogo de etapas: %w", op, err)
	}

	// 3. Cria uma instância de cada etapa padrão para a nova obra
	for _, etapaPadrao := range etapasPadrao {
		// Regra de negócio: A primeira etapa ("Fundações") começa em planejamento, as outras pendentes.
		status := obras.StatusEtapaPendente
		if etapaPadrao.Nome == "Fundações" {
			status = obras.StatusEtapaEmAndamento // Corrigido para "Em Andamento" para ser mais realista
		}

		novaEtapa := &obras.Etapa{
			ID:                 uuid.NewString(),
			ObraID:             novaObra.ID,
			Nome:               etapaPadrao.Nome,
			DataInicioPrevista: &time.Time{}, // Data padrão, pode ser ajustada depois
			DataFimPrevista:    &time.Time{}, // Previsão de 1 mês
			Status:             status,
		}
		if err := s.etapaRepo.Salvar(ctx, tx, novaEtapa); err != nil {
			return nil, fmt.Errorf("%s: falha ao salvar etapa '%s': %w", op, novaEtapa.Nome, err)
		}
	}

	// Se tudo deu certo, comita a transação
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: falha ao comitar transação: %w", op, err)
	}

	s.logger.InfoContext(ctx, "nova obra criada com etapas padrão", "obra_id", novaObra.ID, "etapas_criadas", len(etapasPadrao))
	return novaObra, nil
}
func (s *Service) ListarEtapasPadrao(ctx context.Context) ([]*obras.EtapaPadrao, error) {
	return s.etapaPadraoRepo.ListarTodas(ctx)
}

func (s *Service) CriarEtapaPadrao(ctx context.Context, input dto.CriarEtapaPadraoInput) (*obras.EtapaPadrao, error) {
	const op = "service.obras.CriarEtapaPadrao"

	novaEtapa := &obras.EtapaPadrao{
		ID:        uuid.NewString(),
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Ordem:     input.Ordem,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.etapaPadraoRepo.Salvar(ctx, novaEtapa); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return novaEtapa, nil
}

func (s *Service) BuscarEtapaPadrao(ctx context.Context, id string) (*obras.EtapaPadrao, error) {
	const op = "service.obras.BuscarEtapaPadrao"
	etapa, err := s.etapaPadraoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return etapa, nil
}

func (s *Service) AtualizarEtapaPadrao(ctx context.Context, id string, input dto.AtualizarEtapaPadraoInput) (*obras.EtapaPadrao, error) {
	const op = "service.obras.AtualizarEtapaPadrao"

	etapa, err := s.etapaPadraoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	etapa.Nome = input.Nome
	etapa.Descricao = input.Descricao
	etapa.Ordem = input.Ordem
	etapa.UpdatedAt = time.Now()

	if err := s.etapaPadraoRepo.Atualizar(ctx, etapa); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return etapa, nil
}

func (s *Service) DeletarEtapaPadrao(ctx context.Context, id string) error {
	const op = "service.obras.DeletarEtapaPadrao"
	// Primeiro, verifica se a etapa padrão existe.
	if _, err := s.etapaPadraoRepo.BuscarPorID(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return s.etapaPadraoRepo.Deletar(ctx, id)
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

	// Inicia a transação
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx) // Garante o rollback em caso de erro
	// 1. Valida se a obra de destino existe.
	if _, err := s.obraRepo.BuscarPorID(ctx, obraID); err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	// 2. Busca a Etapa Padrão no catálogo para obter o nome e outras informações.
	etapaPadrao, err := s.etapaPadraoRepo.BuscarPorID(ctx, input.EtapaPadraoID)
	if err != nil {
		return nil, fmt.Errorf("%s: etapa padrão não encontrada no catálogo: %w", op, err)
	}

	// 3. Valida e converte as datas da requisição.
	inicio, err := time.Parse("2006-01-02", input.DataInicioPrevista)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de início inválido: %w", op, err)
	}
	fim, err := time.Parse("2006-01-02", input.DataFimPrevista)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de fim inválido: %w", op, err)
	}

	// 4. Cria a nova entidade 'Etapa' (a instância na obra).
	novaEtapa := &obras.Etapa{
		ID:                 uuid.NewString(),
		ObraID:             obraID,
		Nome:               etapaPadrao.Nome, // O nome vem do catálogo
		DataInicioPrevista: &inicio,
		DataFimPrevista:    &fim,
		Status:             obras.StatusEtapaPendente, // Status inicial padrão
	}

	// 5. Salva a nova etapa no banco de dados.
	if err := s.etapaRepo.Salvar(ctx, tx, novaEtapa); err != nil {
		return nil, fmt.Errorf("%s: falha ao salvar nova etapa na obra: %w", op, err)
	}

	s.logger.InfoContext(ctx, "etapa adicionada à obra com sucesso", "etapa_id", novaEtapa.ID, "obra_id", obraID, "nome_etapa", novaEtapa.Nome)
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
func (s *Service) AlocarFuncionarios(ctx context.Context, obraID string, input dto.AlocarFuncionariosInput) ([]*obras.Alocacao, error) {
	const op = "service.obras.AlocarFuncionarios"

	// --- Lógica de Negócio e Validação ---
	// 1. Valida a obra uma única vez.
	_, err := s.obrasQuerier.BuscarDashboardPorID(ctx, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	// 2. Valida a data de início uma única vez.
	inicio, err := time.Parse("2006-01-02", input.DataInicioAlocacao)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data inválido: %w", op, err)
	}

	var alocacoesParaSalvar []*obras.Alocacao
	// 3. Itera sobre cada funcionário para validar e criar a alocação.
	for _, funcionarioID := range input.FuncionarioIDs {
		// Verifica se o funcionário existe.
		_, err = s.pessoalFinder.BuscarPorID(ctx, funcionarioID)
		if err != nil {
			return nil, fmt.Errorf("%s: funcionário com ID %s não encontrado: %w", op, funcionarioID, err)
		}

		novaAlocacao := &obras.Alocacao{
			ID:                 uuid.NewString(),
			ObraID:             obraID,
			FuncionarioID:      funcionarioID,
			DataInicioAlocacao: inicio,
		}
		alocacoesParaSalvar = append(alocacoesParaSalvar, novaAlocacao)
	}

	// 4. Salva todas as alocações de uma vez.
	if err := s.alocacaoRepo.SalvarMuitos(ctx, alocacoesParaSalvar); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "funcionários alocados com sucesso", "obra_id", obraID, "quantidade", len(alocacoesParaSalvar))
	return alocacoesParaSalvar, nil
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

func (s *Service) BuscarDetalhesPorID(ctx context.Context, obraID string) (*dto.ObraDetalhadaDTO, error) {
	return s.obrasQuerier.BuscarDetalhesPorID(ctx, obraID)
}

func (s *Service) AtualizarObra(ctx context.Context, obraID string, input dto.AtualizarObraInput) (*obras.Obra, error) {
	const op = "service.obras.AtualizarObra"

	if input.Nome == "" || input.Cliente == "" || input.Endereco == "" {
		return nil, fmt.Errorf("%s: nome, cliente e endereço são obrigatórios", op)
	}

	dataInicio, err := time.Parse("2006-01-02", input.DataInicio)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de início inválido: %w", op, err)
	}
	dataFim, err := time.Parse("2006-01-02", input.DataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: formato de data de fim inválido: %w", op, err)
	}

	obraAtualizada := &obras.Obra{
		ID:         obraID,
		Nome:       input.Nome,
		Cliente:    input.Cliente,
		Endereco:   input.Endereco,
		DataInicio: dataInicio,
		DataFim:    &dataFim,
		Status:     obras.Status(input.Status),
		Descricao:  &input.Descricao,
	}

	if err := s.obraRepo.Atualizar(ctx, obraAtualizada); err != nil {
		return nil, fmt.Errorf("%s: falha ao atualizar obra: %w", op, err)
	}

	s.logger.InfoContext(ctx, "obra atualizada com sucesso", "obra_id", obraAtualizada.ID)

	return obraAtualizada, nil
}

func (s *Service) ListarEtapasPorObra(ctx context.Context, obraID string) ([]*obras.Etapa, error) {
	const op = "service.obras.ListarEtapasPorObra"

	// Valida se a obra existe antes de buscar suas etapas
	if _, err := s.obraRepo.BuscarPorID(ctx, obraID); err != nil {
		return nil, fmt.Errorf("%s: obra não encontrada: %w", op, err)
	}

	return s.etapaRepo.ListarPorObraID(ctx, obraID)
}

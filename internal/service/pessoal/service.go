// file: internal/service/pessoal/service.go
package pessoal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
)

var (
	ErrFuncionarioAlocado = errors.New("não é possível excluir um funcionário que está alocado em uma obra ativa")
)

type EventPublisher interface {
	Publicar(ctx context.Context, evento bus.Evento)
}
type ObraFinder interface {
	BuscarPorID(ctx context.Context, id string) (*obras.Obra, error)
}
type Service struct {
	repo            pessoal.FuncionarioRepository
	apontamentoRepo pessoal.ApontamentoRepository // NOVA DEPENDÊNCIA
	alocacaoFinder  AlocacaoFinder
	obraFinder      ObraFinder // NOVA DEPENDÊNCIA
	logger          *slog.Logger
	eventBus        EventPublisher
}

func NovoServico(repo pessoal.FuncionarioRepository, apontamentoRepo pessoal.ApontamentoRepository, alocacaoFinder AlocacaoFinder, obraFinder ObraFinder, eventBus EventPublisher, logger *slog.Logger) *Service {
	return &Service{
		repo:            repo,
		apontamentoRepo: apontamentoRepo,
		alocacaoFinder:  alocacaoFinder,
		obraFinder:      obraFinder,
		eventBus:        eventBus,
		logger:          logger,
	}
}

type AlocacaoFinder interface {
	ExistemAlocacoesAtivasParaFuncionario(ctx context.Context, funcionarioID string) (bool, error)
}

func (s *Service) CadastrarFuncionario(ctx context.Context, nome, cpf, cargo, departamento string) (*pessoal.Funcionario, error) {
	const op = "service.pessoal.CadastrarFuncionario"

	novoFuncionario := &pessoal.Funcionario{
		ID:              uuid.NewString(),
		Nome:            nome,
		CPF:             cpf,
		Cargo:           cargo,
		Departamento:    departamento,
		DataContratacao: time.Now(),
		Status:          "Ativo",
	}

	if err := s.repo.Salvar(ctx, novoFuncionario); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "novo funcionário cadastrado", "funcionario_id", novoFuncionario.ID)
	return novoFuncionario, nil
}

func (s *Service) DeletarFuncionario(ctx context.Context, id string) error {
	const op = "service.pessoal.DeletarFuncionario"

	// Regra de Negócio: Verificar se o funcionário tem alocações ativas.
	alocado, err := s.alocacaoFinder.ExistemAlocacoesAtivasParaFuncionario(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: falha ao verificar alocações: %w", op, err)
	}
	if alocado {
		return ErrFuncionarioAlocado
	}

	if err := s.repo.Deletar(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.logger.InfoContext(ctx, "funcionário excluído (soft delete)", "funcionario_id", id)
	return nil
}

func (s *Service) ListarFuncionarios(ctx context.Context) ([]*pessoal.Funcionario, error) {
	const op = "service.pessoal.ListarFuncionarios"
	funcionarios, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: erro ao listar funcionários: %w", op, err)
	}
	s.logger.InfoContext(ctx, "lista de funcionários recuperada", "total", len(funcionarios))
	return funcionarios, nil
}

func (s *Service) AtualizarFuncionario(ctx context.Context, funcionario *pessoal.Funcionario) error {
	const op = "service.pessoal.AtualizarFuncionario"

	// Verifica se o funcionário existe
	existente, err := s.repo.BuscarPorID(ctx, funcionario.ID)
	if err != nil {
		return fmt.Errorf("%s: erro ao buscar funcionário: %w", op, err)
	}
	if existente == nil {
		return fmt.Errorf("%s: funcionário não encontrado com ID %s", op, funcionario.ID)
	}

	// Atualiza os dados do funcionário
	if err := s.repo.Atualizar(ctx, funcionario); err != nil {
		return fmt.Errorf("%s: erro ao atualizar funcionário: %w", op, err)
	}
	s.logger.InfoContext(ctx, "funcionário atualizado", "funcionario_id", funcionario.ID)
	return nil
}

func (s *Service) BuscarPorID(ctx context.Context, id string) (*pessoal.Funcionario, error) {
	const op = "service.pessoal.BuscarPorID"

	funcionario, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: erro ao buscar funcionário: %w", op, err)
	}
	if funcionario == nil {
		return nil, fmt.Errorf("%s: funcionário não encontrado com ID %s", op, id)
	}
	s.logger.InfoContext(ctx, "funcionário encontrado", "funcionario_id", funcionario.ID)
	return funcionario, nil
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

	apontamento := &pessoal.ApontamentoQuinzenal{
		ID:            uuid.NewString(),
		FuncionarioID: input.FuncionarioID,
		ObraID:        input.ObraID,
		PeriodoInicio: inicio,
		PeriodoFim:    fim,
		Status:        pessoal.StatusApontamentoEmAberto, // Usando a constante do domínio
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
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

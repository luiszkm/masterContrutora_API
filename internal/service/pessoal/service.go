// file: internal/service/pessoal/service.go
package pessoal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
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
type PessoalQuerier interface {
	ListarComUltimoApontamento(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ListagemFuncionarioDTO, *common.PaginacaoInfo, error)
}

type Service struct {
	repo            pessoal.FuncionarioRepository
	apontamentoRepo pessoal.ApontamentoRepository // NOVA DEPENDÊNCIA
	alocacaoFinder  AlocacaoFinder
	obraFinder      ObraFinder     // NOVA DEPENDÊNCIA
	querier         PessoalQuerier // NOVA DEPENDÊNCIA
	logger          *slog.Logger
	eventBus        EventPublisher
	dbpool          *pgxpool.Pool // NOVA DEPENDÊNCIA

}

// ListarComUltimoApontamento implements pessoal.Service.

func NovoServico(
	repo pessoal.FuncionarioRepository,
	apontamentoRepo pessoal.ApontamentoRepository,
	alocacaoFinder AlocacaoFinder,
	obraFinder ObraFinder,
	eventBus EventPublisher,
	querier PessoalQuerier,
	logger *slog.Logger,
	dbpool *pgxpool.Pool, // NOVA DEPENDÊNCIA
) *Service {
	return &Service{
		repo:            repo,
		apontamentoRepo: apontamentoRepo,
		alocacaoFinder:  alocacaoFinder,
		obraFinder:      obraFinder,
		eventBus:        eventBus,
		querier:         querier,
		logger:          logger,
		dbpool:          dbpool, // Atribui a dependência
	}
}

type AlocacaoFinder interface {
	ExistemAlocacoesAtivasParaFuncionario(ctx context.Context, funcionarioID string) (bool, error)
}

func (s *Service) CadastrarFuncionario(ctx context.Context, nome, cpf, cargo, departamento, telefone, ChavePix string, diaria float64) (*pessoal.Funcionario, error) {
	const op = "service.pessoal.CadastrarFuncionario"

	novoFuncionario := &pessoal.Funcionario{
		ID:              uuid.NewString(),
		Nome:            nome,
		CPF:             cpf,
		Cargo:           cargo,
		Departamento:    departamento,
		DataContratacao: time.Now(),
		Status:          "Ativo",
		ChavePix:        ChavePix, // Inicialmente vazio, pode ser atualizado depois
		Telefone:        telefone, // Inicialmente vazio, pode ser atualizado depois
		Diaria:          diaria,
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

func (s *Service) AtualizarFuncionario(ctx context.Context, id string, input dto.AtualizarFuncionarioInput) (*pessoal.Funcionario, error) {
	const op = "service.pessoal.AtualizarFuncionario"

	// 1. Busca o funcionário existente para garantir que ele existe.
	funcionario, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: erro ao buscar funcionário para atualização: %w", op, err)
	}

	// 2. Atualiza apenas os campos que foram fornecidos (não são nulos).
	if input.Nome != nil {
		funcionario.Nome = *input.Nome
	}
	if input.Email != nil {
		funcionario.Email = *input.Email
	}
	if input.AvaliacaoDesempenho != nil {
		funcionario.AvaliacaoDesempenho = *input.AvaliacaoDesempenho
	}
	if input.Observacoes != nil {
		funcionario.Observacoes = *input.Observacoes
	}

	if input.CPF != nil {
		funcionario.CPF = *input.CPF
	}
	if input.Telefone != nil {
		funcionario.Telefone = *input.Telefone
	}
	if input.Cargo != nil {
		funcionario.Cargo = *input.Cargo
	}
	if input.Departamento != nil {
		funcionario.Departamento = *input.Departamento
	}
	if input.ValorDiaria != nil {
		funcionario.ValorDiaria = *input.ValorDiaria
	}
	if input.ChavePix != nil {
		funcionario.ChavePix = *input.ChavePix
	}
	if input.Status != nil {
		funcionario.Status = *input.Status
	}
	if input.MotivoDesligamento != nil {
		funcionario.MotivoDesligamento = *input.MotivoDesligamento
	}
	if input.DesligamentoData != nil {
		if *input.DesligamentoData == "" {
			funcionario.DesligamentoData = nil // Permite limpar a data
		} else {
			data, err := time.Parse("2006-01-02", *input.DesligamentoData)
			if err != nil {
				return nil, fmt.Errorf("%s: formato de data de desligamento inválido: %w", op, err)
			}
			funcionario.DesligamentoData = &data
		}
	}

	funcionario.UpdatedAt = time.Now()

	// 3. Persiste o objeto funcionário completo e atualizado.
	if err := s.repo.Atualizar(ctx, funcionario); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.InfoContext(ctx, "funcionário atualizado com sucesso", "funcionario_id", id)
	return funcionario, nil
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

func (s *Service) AtivarFuncionario(ctx context.Context, id string) error {
	const op = "service.pessoal.Ativar"

	// 1. Busca o funcionário existente para garantir que ele existe.
	err := s.repo.AtivarFuncionario(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: erro ao ativar funcionário: %w", op, err)
	}
	s.logger.InfoContext(ctx, "funcionário ativado com sucesso", "funcionario_id", id)
	return nil

}

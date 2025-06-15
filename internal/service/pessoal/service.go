// file: internal/service/pessoal/service.go
package pessoal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
)

var (
	ErrFuncionarioAlocado = errors.New("não é possível excluir um funcionário que está alocado em uma obra ativa")
)

type Service struct {
	repo           pessoal.FuncionarioRepository
	alocacaoFinder AlocacaoFinder // NOVA DEPENDÊNCIA
	logger         *slog.Logger
}
type AlocacaoFinder interface {
	ExistemAlocacoesAtivasParaFuncionario(ctx context.Context, funcionarioID string) (bool, error)
}

func NovoServico(repo pessoal.FuncionarioRepository, alocacaoFinder AlocacaoFinder, logger *slog.Logger) *Service {
	return &Service{repo: repo, alocacaoFinder: alocacaoFinder, logger: logger}
}

func (s *Service) CadastrarFuncionario(ctx context.Context, nome, cpf, cargo string, salario, diaria float64) (*pessoal.Funcionario, error) {
	const op = "service.pessoal.CadastrarFuncionario"

	novoFuncionario := &pessoal.Funcionario{
		ID:              uuid.NewString(),
		Nome:            nome,
		CPF:             cpf,
		Cargo:           cargo,
		DataContratacao: time.Now(),
		Salario:         salario,
		Diaria:          diaria,
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

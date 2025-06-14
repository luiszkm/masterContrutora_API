// file: internal/service/pessoal/service.go
package pessoal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
)

type Service struct {
	repo   pessoal.Repository
	logger *slog.Logger
}

func NovoServico(repo pessoal.Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
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

package pessoal

import (
	"context"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// Repository define o contrato para a persistência de Funcionarios.
type FuncionarioRepository interface {
	Salvar(ctx context.Context, funcionario *Funcionario) error
	BuscarPorID(ctx context.Context, funcionarioID string) (*Funcionario, error)
	Deletar(ctx context.Context, id string) error // NOVO
	Listar(ctx context.Context) ([]*Funcionario, error)
	Atualizar(ctx context.Context, funcionario *Funcionario) error // NOVO
	AtivarFuncionario(ctx context.Context, id string) error
}

type ApontamentoRepository interface {
	Salvar(ctx context.Context, db db.DBTX, apontamento *ApontamentoQuinzenal) error // Modificado
	BuscarPorID(ctx context.Context, id string) (*ApontamentoQuinzenal, error)
	Atualizar(ctx context.Context, db db.DBTX, apontamento *ApontamentoQuinzenal) error // Modificado
	Listar(ctx context.Context, filtros common.ListarFiltros) ([]*ApontamentoQuinzenal, *common.PaginacaoInfo, error)
	ListarPorFuncionarioID(ctx context.Context, funcionarioID string, filtros common.ListarFiltros) ([]*ApontamentoQuinzenal, *common.PaginacaoInfo, error)
	ExisteApontamentoEmAberto(ctx context.Context, funcionarioID string) (bool, error)
	BuscarUltimoPorFuncionarioID(ctx context.Context, funcionarioID string) (*ApontamentoQuinzenal, error)
}

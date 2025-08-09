package obras

import (
	"context"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

type ObrasRepository interface {
	Salvar(ctx context.Context, db db.DBTX, obra *Obra) error // Modificado
	BuscarPorID(ctx context.Context, id string) (*Obra, error)
	Deletar(ctx context.Context, id string) error
	Atualizar(ctx context.Context, obra *Obra) error
}

type AlocacaoRepository interface {
	Salvar(ctx context.Context, alocacao *Alocacao) error
	SalvarMuitos(ctx context.Context, alocacoes []*Alocacao) error
	ExistemAlocacoesAtivasParaFuncionario(ctx context.Context, funcionarioID string) (bool, error) // NOVO

}

type EtapaRepository interface {
	Salvar(ctx context.Context, db db.DBTX, etapa *Etapa) error           // Modificado
	BuscarPorID(ctx context.Context, etapaID string) (*Etapa, error)      // NOVO
	Atualizar(ctx context.Context, etapa *Etapa) error                    // NOVO
	ListarPorObraID(ctx context.Context, obraID string) ([]*Etapa, error) // NOVO MÉTODO

}

type EtapaPadraoRepository interface {
	Salvar(ctx context.Context, etapa *EtapaPadrao) error
	Atualizar(ctx context.Context, etapa *EtapaPadrao) error
	BuscarPorID(ctx context.Context, id string) (*EtapaPadrao, error)
	ListarTodas(ctx context.Context) ([]*EtapaPadrao, error)
	Deletar(ctx context.Context, id string) error
}

type CronogramaRecebimentoRepository interface {
	Salvar(ctx context.Context, db db.DBTX, cronograma *CronogramaRecebimento) error
	SalvarMuitos(ctx context.Context, db db.DBTX, cronogramas []*CronogramaRecebimento) error
	Atualizar(ctx context.Context, cronograma *CronogramaRecebimento) error
	BuscarPorID(ctx context.Context, id string) (*CronogramaRecebimento, error)
	ListarPorObraID(ctx context.Context, obraID string) ([]*CronogramaRecebimento, error)
	ListarVencidosPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]*CronogramaRecebimento, error)
	Deletar(ctx context.Context, id string) error
}

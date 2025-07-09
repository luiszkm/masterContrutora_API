package obras

import (
	"context"
)

type ObrasRepository interface {
	Salvar(ctx context.Context, obra *Obra) error
	BuscarPorID(ctx context.Context, id string) (*Obra, error)
	Deletar(ctx context.Context, id string) error
}

type AlocacaoRepository interface {
	Salvar(ctx context.Context, alocacao *Alocacao) error
	SalvarMuitos(ctx context.Context, alocacoes []*Alocacao) error
	ExistemAlocacoesAtivasParaFuncionario(ctx context.Context, funcionarioID string) (bool, error) // NOVO

}

type EtapaRepository interface {
	Salvar(ctx context.Context, etapa *Etapa) error
	BuscarPorID(ctx context.Context, etapaID string) (*Etapa, error) // NOVO
	Atualizar(ctx context.Context, etapa *Etapa) error               // NOVO

}

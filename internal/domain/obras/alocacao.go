// file: internal/domain/obras/alocacao.go
package obras

import (
	"time"
)

// Alocacao representa o agregado que vincula um funcionário a uma obra
// por um determinado período.
type Alocacao struct {
	ID                 string
	ObraID             string
	FuncionarioID      string
	DataInicioAlocacao time.Time
	DataFimAlocacao    *time.Time
}

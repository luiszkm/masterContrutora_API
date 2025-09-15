// file: internal/domain/obras/etapa.go
package obras

import (
	"time"
)

// enum de status de etapa
type StatusEtapa string

const (
	StatusEtapaPendente    StatusEtapa = "Pendente"
	StatusEtapaEmAndamento StatusEtapa = "Em Andamento"
	StatusEtapaConcluida   StatusEtapa = "Completa"
)

type Etapa struct {
	ID                 string
	ObraID             string
	Nome               string
	Ordem              int        `json:"ordem"`
	DataInicioPrevista *time.Time `json:"data_inicio_prevista"`
	DataFimPrevista    *time.Time `json:"data_fim_prevista"`
	Status             StatusEtapa
}

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
	StatusEtapaConcluida   StatusEtapa = "Concluída"
)

type Etapa struct {
	ID                 string
	ObraID             string
	Nome               string
	DataInicioPrevista time.Time
	DataFimPrevista    time.Time
	Status             StatusEtapa
}

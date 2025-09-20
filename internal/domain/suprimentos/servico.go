package suprimentos

import "time"

type Servico struct {
	ID        string
	Nome      string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
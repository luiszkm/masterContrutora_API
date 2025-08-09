package suprimentos

import "time"

type Categoria struct {
	ID        string
	Nome      string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

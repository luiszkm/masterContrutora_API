// file: internal/service/obras/dto/obra_list_item_dto.go
package dto

// ObraListItemDTO representa os dados de uma única obra em uma lista.
type ObraListItemDTO struct {
	ID      string `json:"id"`
	Nome    string `json:"nome"`
	Cliente string `json:"cliente"`
	Status  string `json:"status"`
}

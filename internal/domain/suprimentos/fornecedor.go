// file: internal/domain/suprimentos/fornecedor.go
package suprimentos

type Fornecedor struct {
	ID            string
	Nome          string
	CNPJ          string
	Categoria     string
	Contato       string
	Email         string
	Status        string // "Ativo", "Inativo"
	Website       *string
	Endereco      *string
	NomeAtendente *string
	Avaliacao     *float64 // Avaliação de 0 a 5
	Observacoes   *string
}

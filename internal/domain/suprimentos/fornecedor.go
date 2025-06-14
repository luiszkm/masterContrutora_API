// file: internal/domain/suprimentos/fornecedor.go
package suprimentos

type Fornecedor struct {
	ID        string
	Nome      string
	CNPJ      string
	Categoria string
	Contato   string
	Email     string
	Status    string // "Ativo", "Inativo"
}

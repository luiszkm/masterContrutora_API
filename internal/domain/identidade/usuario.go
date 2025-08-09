// file: internal/domain/identidade/usuario.go
package identidade

// Remova o import "golang.org/x/crypto/bcrypt"

// A struct Usuario agora é um simples contêiner de dados, sem dependências externas.
type Usuario struct {
	ID         string
	Nome       string
	Email      string
	SenhaHash  string
	Permissoes []string
	Ativo      bool
}

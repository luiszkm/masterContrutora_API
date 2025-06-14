// file: pkg/security/password.go
package security

import "golang.org/x/crypto/bcrypt"

// Hasher define a interface para o nosso serviço de hashing de senhas.
// A aplicação dependerá desta interface, não de uma implementação específica.
type Hasher interface {
	Hash(senha string) (string, error)
	Checar(senha, hash string) bool
}

// BcryptHasher é a implementação concreta que usa a biblioteca bcrypt.
type BcryptHasher struct{}

// NewBcryptHasher cria uma nova instância do nosso hasher.
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Hash gera um hash a partir de uma senha em texto plano.
func (b *BcryptHasher) Hash(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	return string(bytes), err
}

// Checar compara uma senha em texto plano com o hash.
func (b *BcryptHasher) Checar(senha, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
	return err == nil
}

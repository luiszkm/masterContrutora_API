// file: internal/repository/postgres/usuario_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/identidade"
)

// UsuarioRepositoryPostgres implementa a persistência para o agregado Usuario.
type UsuarioRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewUsuarioRepository é o construtor para o repositório de usuários.
func NewUsuarioRepository(db *pgxpool.Pool, logger *slog.Logger) identidade.UsuarioRepository {
	return &UsuarioRepositoryPostgres{
		db:     db,
		logger: logger,
	}
}

// Salvar insere um novo usuário no banco de dados.
func (r *UsuarioRepositoryPostgres) Salvar(ctx context.Context, usuario *identidade.Usuario) error {
	const op = "repository.postgres.usuario.Salvar"

	query := `
		INSERT INTO usuarios (id, nome, email, senha_hash, permissoes, ativo)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		usuario.ID,
		usuario.Nome,
		usuario.Email,
		usuario.SenhaHash,
		usuario.Permissoes,
		usuario.Ativo,
	)

	if err != nil {
		// TODO: Adicionar tratamento para erro de violação de constraint UNIQUE do email.
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// BuscarPorEmail encontra um usuário pelo seu endereço de e-mail.
func (r *UsuarioRepositoryPostgres) BuscarPorEmail(ctx context.Context, email string) (*identidade.Usuario, error) {
	const op = "repository.postgres.usuario.BuscarPorEmail"

	query := `SELECT id, nome, email, senha_hash, permissoes, ativo FROM usuarios WHERE email = $1 AND ativo = TRUE`
	row := r.db.QueryRow(ctx, query, email)

	var u identidade.Usuario
	err := row.Scan(
		&u.ID,
		&u.Nome,
		&u.Email,
		&u.SenhaHash,
		&u.Permissoes,
		&u.Ativo,
	)

	if err != nil {
		// Traduz o erro do driver para um erro de domínio conhecido.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

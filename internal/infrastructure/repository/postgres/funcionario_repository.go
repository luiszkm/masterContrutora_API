// file: internal/repository/postgres/funcionario_repository.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	pessoal_dto "github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
)

type FuncionarioRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoFuncionarioRepository(db *pgxpool.Pool, logger *slog.Logger) *FuncionarioRepositoryPostgres {
	return &FuncionarioRepositoryPostgres{db: db, logger: logger}
}

// BuscarPorID implements obras.PessoalFinder.
func (r *FuncionarioRepositoryPostgres) BuscarPorID(ctx context.Context, funcionarioID string) (*pessoal.Funcionario, error) {

	const op = "repository.postgres.funcionario.BuscarPorID"
	query := `
		SELECT 
		    id, nome, cpf, telefone, cargo, departamento, data_contratacao, 
		    valor_diaria, chave_pix, status, desligamento_data, motivo_desligamento, 
			created_at, updated_at
		FROM funcionarios
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, funcionarioID)

	var f pessoal.Funcionario
	var desligamento_data sql.NullTime

	err := row.Scan(
		&f.ID, &f.Nome, &f.CPF, &f.Telefone, &f.Cargo, &f.Departamento, &f.DataContratacao,
		&f.ValorDiaria, &f.ChavePix, &f.Status, &desligamento_data, &f.MotivoDesligamento,
		&f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: funcionário não encontrado com ID %s", op, funcionarioID)
		}
		return nil, fmt.Errorf("%s: erro ao buscar funcionário: %w", op, err)
	}
	if desligamento_data.Valid {
		f.DesligamentoData = &desligamento_data.Time
	}
	return &f, nil
}

func (r *FuncionarioRepositoryPostgres) Salvar(ctx context.Context, f *pessoal.Funcionario) error {
	const op = "repository.postgres.funcionario.Salvar"
	query := `
		INSERT INTO funcionarios 
		    (id, nome, cpf, telefone, cargo, departamento, data_contratacao, valor_diaria, chave_pix, status, created_at, updated_at)
		VALUES 
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`
	_, err := r.db.Exec(ctx, query,
		f.ID, f.Nome, f.CPF, f.Telefone, f.Cargo, f.Departamento,
		f.DataContratacao, f.ValorDiaria, f.ChavePix, f.Status,
	)
	if err != nil {
		// TODO: Tratar erro de violação de constraint UNIQUE do CPF
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (r *FuncionarioRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.funcionario.Deletar"
	query := `UPDATE funcionarios SET desligamento_data = NOW() WHERE id = $1 AND desligamento_data IS NULL`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (r *FuncionarioRepositoryPostgres) Atualizar(ctx context.Context, f *pessoal.Funcionario) error {
	const op = "repository.postgres.funcionario.Atualizar"
	query := `
		UPDATE funcionarios
		SET 
		    nome = $1, cpf = $2, telefone = $3, cargo = $4, departamento = $5, 
		    valor_diaria = $6, chave_pix = $7, status = $8, updated_at = NOW()
		WHERE id = $9 AND desligamento_data IS NULL
	`
	cmd, err := r.db.Exec(ctx, query,
		f.Nome, f.CPF, f.Telefone, f.Cargo, f.Departamento,
		f.ValorDiaria, f.ChavePix, f.Status, f.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
func (r *FuncionarioRepositoryPostgres) Listar(ctx context.Context) ([]*pessoal.Funcionario, error) {
	const op = "repository.postgres.funcionario.Listar"
	query := `
		SELECT id, nome, cpf, cargo, departamento, status, valor_diaria
		FROM funcionarios
		ORDER BY nome ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var funcionarios []*pessoal.Funcionario
	for rows.Next() {
		var f pessoal.Funcionario
		if err := rows.Scan(
			&f.ID,
			&f.Nome,
			&f.CPF,
			&f.Cargo,
			&f.Departamento,
			&f.Status,
			&f.ValorDiaria,
		); err != nil {
			return nil, fmt.Errorf("%s: erro ao ler linha: %w", op, err)
		}
		funcionarios = append(funcionarios, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre linhas: %w", op, err)
	}

	return funcionarios, nil
}

func (r *FuncionarioRepositoryPostgres) ListarComUltimoApontamento(ctx context.Context, filtros common.ListarFiltros) ([]*pessoal_dto.ListagemFuncionarioDTO, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.funcionario.ListarComUltimoApontamento"

	args := pgx.NamedArgs{}
	baseQuery := `
		FROM
			funcionarios f
		LEFT JOIN LATERAL (
			SELECT * FROM apontamentos_quinzenais aq
			WHERE aq.funcionario_id = f.id
			ORDER BY aq.periodo_fim DESC
			LIMIT 1
		) a ON true
		WHERE f.desligamento_data IS NULL
	`
	// ... (Lógica de construção de query para filtros e contagem)
	whereClauses := []string{"f.desligamento_data IS NULL"}
	if filtros.Status != "" {
		whereClauses = append(whereClauses, "f.status = @status")
		args["status"] = filtros.Status
	}
	whereString := " WHERE " + strings.Join(whereClauses, " AND ")

	countQueryBuilder := strings.Builder{}
	countQueryBuilder.WriteString("SELECT COUNT(f.id) FROM funcionarios f")
	countQueryBuilder.WriteString(whereString)

	var totalItens int
	err := r.db.QueryRow(ctx, countQueryBuilder.String(), args).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar funcionários: %w", op, err)
	}
	paginacao := common.NewPaginacaoInfo(totalItens, filtros.Pagina, filtros.TamanhoPagina)
	if totalItens == 0 {
		return []*pessoal_dto.ListagemFuncionarioDTO{}, paginacao, nil
	}

	// --- INÍCIO DA CORREÇÃO NA QUERY PRINCIPAL ---
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT
			f.id, f.nome, f.cargo, f.departamento, f.data_contratacao,
			a.diaria, -- A diária vem da tabela de apontamentos (a)
			f.chave_pix,
			a.dias_trabalhados, a.adicionais, a.descontos, a.adiantamentos, a.status
	`)
	queryBuilder.WriteString(baseQuery)
	queryBuilder.WriteString(strings.Replace(whereString, " WHERE ", " AND ", 1)) // Adiciona os filtros à query principal

	offset := (filtros.Pagina - 1) * filtros.TamanhoPagina
	queryBuilder.WriteString(" ORDER BY f.nome ASC LIMIT @limit OFFSET @offset")
	args["limit"] = filtros.TamanhoPagina
	args["offset"] = offset
	// --- FIM DA CORREÇÃO NA QUERY PRINCIPAL ---

	rows, err := r.db.Query(ctx, queryBuilder.String(), args)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var lista []*pessoal_dto.ListagemFuncionarioDTO
	for rows.Next() {
		var dto pessoal_dto.ListagemFuncionarioDTO
		// Variáveis para receber valores que podem ser nulos
		var departamento, chavePix, statusApontamento sql.NullString
		var diasTrabalhados sql.NullInt32
		var adicionais, descontos, adiantamento, diariaApontamento sql.NullFloat64

		// --- INÍCIO DA CORREÇÃO NO SCAN ---
		// A ordem e quantidade dos campos agora correspondem ao SELECT
		err := rows.Scan(
			&dto.ID, &dto.Nome, &dto.Cargo, &departamento, &dto.DataContratacao,
			&diariaApontamento, &chavePix, &diasTrabalhados, &adicionais,
			&descontos, &adiantamento, &statusApontamento,
		)
		// --- FIM DA CORREÇÃO NO SCAN ---

		if err != nil {
			return nil, nil, fmt.Errorf("%s: falha ao escanear linha: %w", op, err)
		}

		// Converte os tipos Null* para ponteiros no DTO final
		if departamento.Valid {
			dto.Departamento = &departamento.String
		}
		if chavePix.Valid {
			dto.ChavePix = &chavePix.String
		}
		if statusApontamento.Valid {
			dto.StatusApontamento = &statusApontamento.String
		}
		if diasTrabalhados.Valid {
			v := int(diasTrabalhados.Int32)
			dto.DiasTrabalhados = &v
		}
		if adicionais.Valid {
			dto.ValorAdicional = &adicionais.Float64
		}
		if descontos.Valid {
			dto.Descontos = &descontos.Float64
		}
		if adiantamento.Valid {
			dto.Adiantamento = &adiantamento.Float64
		}
		// O campo 'ValorDiaria' do DTO agora é preenchido com a diária do apontamento.
		if diariaApontamento.Valid {
			dto.Diaria = diariaApontamento.Float64
		}

		lista = append(lista, &dto)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("%s: erro ao iterar sobre as linhas: %w", op, err)
	}

	return lista, paginacao, nil
}

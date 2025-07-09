// file: cmd/seeder/main.go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
)

const (
	defaultObraName = "Sem Obra"
)

func main() {
	// 1. Configuração do Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Carregamento das Variáveis de Ambiente
	if err := godotenv.Load(); err != nil {
		logger.Warn("arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("a variável de ambiente DATABASE_URL é obrigatória")
		os.Exit(1)
	}

	// 3. Conexão com o Banco de Dados
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("não foi possível conectar ao banco de dados", "erro", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("conexão com o PostgreSQL estabelecida com sucesso")

	// 4. Execução do Seeder
	if err := seedDefaultObra(ctx, dbpool, logger); err != nil {
		logger.Error("falha ao executar o seeder", "erro", err)
		os.Exit(1)
	}
	if err := seedCategoriasMaterial(ctx, dbpool, logger); err != nil {
		logger.Error("falha ao inserir categorias de material", "erro", err)
		os.Exit(1)
	}

	logger.Info("seeder executado com sucesso!")
}

// seedDefaultObra verifica e cria a obra e etapas padrão se não existirem.
func seedDefaultObra(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) error {
	const op = "seeder.seedDefaultObra"

	// Passo 1: Verificar se a "Sem Obra" já existe para garantir a idempotência.
	var exists bool
	queryExists := `SELECT EXISTS(SELECT 1 FROM obras WHERE nome = $1)`
	err := db.QueryRow(ctx, queryExists, defaultObraName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: falha ao verificar existência da obra: %w", op, err)
	}

	if exists {
		logger.Info("a obra padrão 'Sem Obra' já existe. Nenhuma ação necessária.")
		return nil
	}

	// Se não existe, vamos criar a obra e as etapas dentro de uma transação.
	logger.Info("obra padrão 'Sem Obra' não encontrada. Iniciando processo de criação...")

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx) // Garante o rollback em caso de pânico ou erro

	// Passo 2: Criar a Obra "Sem Obra"
	obraPadrao := obras.Obra{
		ID:         uuid.NewString(),
		Nome:       defaultObraName,
		Cliente:    "Interno",
		Endereco:   "N/A",
		DataInicio: time.Now(),
		Status:     obras.StatusEmAndamento,
		Descricao:  "Obra padrão criada pelo seeder para inicialização do sistema.",
	}

	// O método Salvar do repositório precisa ser adaptado para aceitar a transação (tx)
	// Vamos usar tx.Exec diretamente para simplificar o seeder.
	queryObra := `INSERT INTO obras (id, nome, cliente, endereco, data_inicio, status, descricao) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.Exec(ctx, queryObra, obraPadrao.ID,
		obraPadrao.Nome,
		obraPadrao.Cliente,
		obraPadrao.Endereco,
		obraPadrao.DataInicio,
		obraPadrao.Status,
		obraPadrao.Descricao,
	); err != nil {
		return fmt.Errorf("%s: falha ao salvar obra padrão: %w", op, err)
	}
	logger.Info("obra padrão criada com sucesso", "obra_id", obraPadrao.ID)

	// Passo 3: Criar as Etapas Padrão
	// Usei nomes comuns para as etapas 5 e 6 que estavam faltando.
	nomesEtapas := []string{"Fundações", "Estrutura", "Alvenaria", "Instalações", "Acabamentos", "Pintura"}
	queryEtapa := `INSERT INTO etapas (id, obra_id, nome, data_inicio_prevista, data_fim_prevista, status) VALUES ($1, $2, $3, $4, $5, $6)`

	batch := &pgx.Batch{}
	for _, nome := range nomesEtapas {
		etapa := obras.Etapa{
			ID:                 uuid.NewString(),
			ObraID:             obraPadrao.ID,
			Nome:               nome,
			DataInicioPrevista: time.Now(),
			DataFimPrevista:    time.Now().AddDate(0, 1, 0), // Previsão de 1 mês
			Status:             obras.StatusEtapaPendente,
		}
		batch.Queue(queryEtapa, etapa.ID, etapa.ObraID, etapa.Nome, etapa.DataInicioPrevista, etapa.DataFimPrevista, etapa.Status)
		logger.Info("etapa padrão adicionada ao lote", "etapa_nome", etapa.Nome)
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		// The second return value of br.Exec() is pgconn.CommandTag, which is not an error.
		// We only care about the error returned by the first return value.
		return fmt.Errorf("%s: falha ao executar lote de inserção de etapas: %w", op, err)
	}
	if err := br.Close(); err != nil {
		return fmt.Errorf("%s: falha ao fechar lote de inserção: %w", op, err)
	}
	logger.Info("etapas padrão criadas com sucesso.")

	// Se tudo correu bem, comita a transação.
	return tx.Commit(ctx)

}
func seedCategoriasMaterial(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) error {
	const op = "seeder.seedCategoriasMaterial"

	categoriasPadrao := []string{
		"Cimento", "Alvenaria", "Agregados", "Aço", "Acabamentos",
		"Pintura", "Elétrica", "Hidráulica", "Madeiras", "Outros",
	}

	logger.Info("iniciando seeder para categorias de material...")

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Passo 1: Inserir cada categoria se ela não existir.
	queryInsert := `
		INSERT INTO categorias_materiais (id, nome)
		VALUES ($1, $2)
		ON CONFLICT (nome) DO NOTHING
	`
	batch := &pgx.Batch{}
	for _, nomeCategoria := range categoriasPadrao {
		batch.Queue(queryInsert, uuid.NewString(), nomeCategoria)
	}

	br := tx.SendBatch(ctx, batch)
	totalInserted := 0
	for i := 0; i < len(categoriasPadrao); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("%s: falha ao executar inserção de categoria: %w", op, err)
		}
		totalInserted++
	}
	if err != nil {
		return fmt.Errorf("%s: falha ao executar lote de inserção de categorias: %w", op, err) // This error is for communication failures, not constraint violations.
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("%s: falha ao fechar lote de inserção: %w", op, err)
	}

	if totalInserted > 0 {
		logger.Info("novas categorias de material foram inseridas", "quantidade", totalInserted)
	} else {
		logger.Info("todas as categorias de material padrão já existem. Nenhuma ação necessária.")
	}

	// Se tudo correu bem, comita a transação.
	return tx.Commit(ctx)
}

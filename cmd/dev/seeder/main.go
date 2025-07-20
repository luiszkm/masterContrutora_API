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
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/pkg/security"
)

// Estrutura para manter os IDs criados e reutilizá-los
type SeedData struct {
	UserIDs         []string
	CategoriaIDs    map[string]string // Mapeia nome da categoria para ID
	ProdutoIDs      []string
	FornecedorIDs   []string
	ObraIDs         []string
	EtapaIDsPorObra map[string][]string // Mapeia ID da obra para lista de IDs de suas etapas
	FuncionarioIDs  []string            // NOVO

}

func main() {
	// 1. Configuração
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		logger.Warn("arquivo .env não encontrado")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("variável de ambiente DATABASE_URL é obrigatória")
		os.Exit(1)
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("não foi possível conectar ao banco de dados", "erro", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("conexão com o PostgreSQL estabelecida com sucesso")

	// Estrutura para passar os IDs entre as funções do seeder
	seedData := &SeedData{
		CategoriaIDs:    make(map[string]string),
		EtapaIDsPorObra: make(map[string][]string),
	}

	// 2. Execução dos Seeders em Sequência Lógica
	runSeeder(ctx, dbpool, logger, "Usuários Padrão", seedUsers, seedData)
	runSeeder(ctx, dbpool, logger, "Categorias Padrão", seedCategorias, seedData)
	runSeeder(ctx, dbpool, logger, "Produtos Padrão", seedProdutos, seedData)
	runSeeder(ctx, dbpool, logger, "Fornecedores Padrão", seedFornecedores, seedData)
	runSeeder(ctx, dbpool, logger, "Obras e Etapas Padrão", seedObrasEEtapas, seedData)
	runSeeder(ctx, dbpool, logger, "Funcionários Padrão", seedFuncionarios, seedData) // NOVO
	runSeeder(ctx, dbpool, logger, "Apontamentos Padrão", seedApontamentos, seedData) // NOVO
	runSeeder(ctx, dbpool, logger, "Orçamentos Padrão", seedOrcamentos, seedData)
	runSeeder(ctx, dbpool, logger, "Orçamentos Padrão", seedOrcamentos, seedData)

	logger.Info("todos os seeders foram executados com sucesso!")
}

// Helper para executar e logar cada seeder
func runSeeder(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, name string, seederFunc func(context.Context, *pgxpool.Pool, *slog.Logger, *SeedData) error, data *SeedData) {
	if err := seederFunc(ctx, db, logger, data); err != nil {
		logger.Error("falha ao executar seeder", "nome", name, "erro", err)
		os.Exit(1)
	}
}

// --- Funções de Seeding Individuais ---

func seedUsers(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	hasher := security.NewBcryptHasher()
	senhaHash, _ := hasher.Hash("senha_forte_123")
	email := "admin@construtora.com"

	var userID string
	queryCheck := `SELECT id FROM usuarios WHERE email = $1`
	err := db.QueryRow(ctx, queryCheck, email).Scan(&userID)

	if err == pgx.ErrNoRows {
		userID = uuid.NewString()
		queryInsert := `INSERT INTO usuarios (id, nome, email, senha_hash, permissoes, ativo) VALUES ($1, $2, $3, $4, $5, $6)`
		// Simplesmente damos todas as permissões para o usuário admin do seeder
		permissoes := []string{"obras:ler", "obras:escrever", "suprimentos:ler", "suprimentos:escrever"}
		_, err := db.Exec(ctx, queryInsert, userID, "Admin Seeder", email, senhaHash, permissoes, true)
		if err != nil {
			return err
		}
		logger.Info("usuário admin padrão criado", "id", userID)
	} else if err != nil {
		return err
	} else {
		logger.Info("usuário admin padrão já existe")
	}
	data.UserIDs = append(data.UserIDs, userID)
	return nil
}

func seedCategorias(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	categorias := []string{"Cimento", "Aço", "Areia e Brita", "Madeira", "Acabamentos", "Elétrica", "Hidráulica"}
	query := `INSERT INTO categorias (id, nome) VALUES ($1, $2) ON CONFLICT (nome) DO NOTHING RETURNING id, nome`

	for _, nome := range categorias {
		var id, nomeRetornado string
		err := db.QueryRow(ctx, query, uuid.NewString(), nome).Scan(&id, &nomeRetornado)
		if err == pgx.ErrNoRows {
			// Categoria já existia, precisamos buscar o ID dela
			querySelect := `SELECT id FROM categorias WHERE nome = $1`
			db.QueryRow(ctx, querySelect, nome).Scan(&id)
		} else if err != nil {
			return err
		}
		data.CategoriaIDs[nome] = id
	}
	logger.Info("seeder de categorias finalizado", "total_categorias", len(data.CategoriaIDs))
	return nil
}

func seedProdutos(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	produtos := []suprimentos.Produto{
		{Nome: "Cimento CPII 50kg", UnidadeDeMedida: "saco", Categoria: "Cimento"},
		{Nome: "Vergalhão CA-50 10mm", UnidadeDeMedida: "barra", Categoria: "Aço"},
		{Nome: "Areia Média Lavada", UnidadeDeMedida: "m³", Categoria: "Areia e Brita"},
		{Nome: "Tábua de Pinus 30cm", UnidadeDeMedida: "unidade", Categoria: "Madeira"},
		{Nome: "Piso Porcelanato 80x80", UnidadeDeMedida: "m²", Categoria: "Acabamentos"},
		{Nome: "Cabo Flexível 2.5mm", UnidadeDeMedida: "rolo 100m", Categoria: "Elétrica"},
		{Nome: "Tubo PVC Esgoto 100mm", UnidadeDeMedida: "barra 6m", Categoria: "Hidráulica"},
	}
	query := `INSERT INTO produtos (id, nome, unidade_de_medida, categoria) VALUES ($1, $2, $3, $4) ON CONFLICT (nome) DO NOTHING RETURNING id`

	for _, p := range produtos {
		var id string
		err := db.QueryRow(ctx, query, uuid.NewString(), p.Nome, p.UnidadeDeMedida, p.Categoria).Scan(&id)
		if err == pgx.ErrNoRows {
			querySelect := `SELECT id FROM produtos WHERE nome = $1`
			db.QueryRow(ctx, querySelect, p.Nome).Scan(&id)
		} else if err != nil {
			return err
		}
		data.ProdutoIDs = append(data.ProdutoIDs, id)
	}
	logger.Info("seeder de produtos finalizado", "total_produtos", len(data.ProdutoIDs))
	return nil
}

func seedFornecedores(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	fornecedores := []map[string]interface{}{
		{"nome": "Depósito Constrular", "cnpj": "11.111.111/0001-11", "categorias": []string{"Cimento", "Areia e Brita"}},
		{"nome": "Aços & Ferros Brasil", "cnpj": "22.222.222/0001-22", "categorias": []string{"Aço"}},
		{"nome": "Madeireira Pinus Forte", "cnpj": "33.333.333/0001-33", "categorias": []string{"Madeira", "Acabamentos"}},
	}
	queryFornecedor := `INSERT INTO fornecedores (id, nome, cnpj, status) VALUES ($1, $2, $3, 'Ativo') ON CONFLICT (cnpj) DO NOTHING RETURNING id`
	queryCategoria := `INSERT INTO fornecedor_categorias (fornecedor_id, categoria_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	for _, f := range fornecedores {
		var id string
		err := db.QueryRow(ctx, queryFornecedor, uuid.NewString(), f["nome"], f["cnpj"]).Scan(&id)
		if err == pgx.ErrNoRows {
			querySelect := `SELECT id FROM fornecedores WHERE cnpj = $1`
			db.QueryRow(ctx, querySelect, f["cnpj"]).Scan(&id)
		} else if err != nil {
			return err
		}
		data.FornecedorIDs = append(data.FornecedorIDs, id)

		// Associa categorias
		for _, catNome := range f["categorias"].([]string) {
			catID := data.CategoriaIDs[catNome]
			db.Exec(ctx, queryCategoria, id, catID)
		}
	}
	logger.Info("seeder de fornecedores finalizado", "total_fornecedores", len(data.FornecedorIDs))
	return nil
}

func seedObrasEEtapas(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	obras := []map[string]string{
		{"nome": "Residencial Jardins", "cliente": "MRV Engenharia"},
		{"nome": "Edifício Comercial Central", "cliente": "WTorre"},
		{"nome": "Galpão Logístico Sul", "cliente": "Loggi"},
	}
	etapasPadrao := []string{"Fundações", "Estrutura", "Alvenaria", "Instalações", "Acabamentos"}
	queryObra := `INSERT INTO obras (id, nome, cliente, endereco, data_inicio, status) VALUES ($1, $2, $3, 'Endereço Padrão', NOW(), 'Em Andamento') ON CONFLICT (nome) DO NOTHING RETURNING id`
	queryEtapa := `
		INSERT INTO etapas (id, obra_id, nome, status, data_inicio_prevista, data_fim_prevista)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, o := range obras {
		var id string
		err := db.QueryRow(ctx, queryObra, uuid.NewString(), o["nome"], o["cliente"]).Scan(&id)
		if err == pgx.ErrNoRows {
			querySelect := `SELECT id FROM obras WHERE nome = $1`
			db.QueryRow(ctx, querySelect, o["nome"]).Scan(&id)
		} else if err != nil {
			return err
		}
		data.ObraIDs = append(data.ObraIDs, id)

		// Cria as etapas para esta obra com datas de exemplo
		dataInicio := time.Now()
		for i, nomeEtapa := range etapasPadrao {
			status := "Pendente"
			if i == 0 {
				status = "Em Andamento"
			}

			// Define datas sequenciais para cada etapa
			dataFim := dataInicio.AddDate(0, 1, 0) // Adiciona 1 mês
			etapaID := uuid.NewString()

			// CORREÇÃO: Passa as datas para a query
			_, err := db.Exec(ctx, queryEtapa, etapaID, id, nomeEtapa, status, dataInicio, dataFim)
			if err != nil {
				return err
			}
			data.EtapaIDsPorObra[id] = append(data.EtapaIDsPorObra[id], etapaID)

			// A data de início da próxima etapa é a data de fim desta
			dataInicio = dataFim
		}
	}
	logger.Info("seeder de obras e etapas finalizado", "total_obras", len(data.ObraIDs))
	return nil
}

// file: cmd/seeder/main.go

// ... (outras funções do seeder)

func seedOrcamentos(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	if len(data.UserIDs) == 0 || len(data.ObraIDs) < 2 || len(data.ProdutoIDs) < 4 || len(data.FornecedorIDs) < 3 {
		logger.Warn("não há dados prévios suficientes (usuários, obras, produtos, fornecedores) para criar orçamentos de exemplo")
		return nil
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryCheck := `SELECT id FROM orcamentos WHERE numero = $1`
	queryOrcamento := `INSERT INTO orcamentos (id, numero, etapa_id, fornecedor_id, valor_total, status, criado_por_usuario_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	queryItem := `INSERT INTO orcamento_itens (id, orcamento_id, produto_id, quantidade, valor_unitario) VALUES ($1, $2, $3, $4, $5)`

	// --- Orçamento 1 ---
	numero1 := fmt.Sprintf("ORC-%d-JUL-001", time.Now().Year())
	var orcamentoID1 string

	// 1. VERIFICA se o orçamento já existe
	err = tx.QueryRow(ctx, queryCheck, numero1).Scan(&orcamentoID1)
	if err == pgx.ErrNoRows { // 2. Se NÃO EXISTE, cria
		orcamentoID1 = uuid.NewString()
		obraID1 := data.ObraIDs[0]
		etapaID1 := data.EtapaIDsPorObra[obraID1][0] // Pega a primeira etapa ("Fundações")

		_, err = tx.Exec(ctx, queryOrcamento, orcamentoID1, numero1, etapaID1, data.FornecedorIDs[0], 15000.00, "Aprovado", data.UserIDs[0])
		if err != nil {
			return err
		}

		// Cria os itens associados
		_, err = tx.Exec(ctx, queryItem, uuid.NewString(), orcamentoID1, data.ProdutoIDs[0], 200, 50.0) // Cimento
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, queryItem, uuid.NewString(), orcamentoID1, data.ProdutoIDs[1], 50, 100.0) // Vergalhão
		if err != nil {
			return err
		}
		logger.Info("orçamento de exemplo criado", "numero", numero1)
	} else if err != nil {
		return err // Outro erro de banco
	}

	// --- Orçamento 2 ---
	numero2 := fmt.Sprintf("ORC-%d-JUL-002", time.Now().Year())
	var orcamentoID2 string

	// 1. VERIFICA se o orçamento já existe
	err = tx.QueryRow(ctx, queryCheck, numero2).Scan(&orcamentoID2)
	if err == pgx.ErrNoRows { // 2. Se NÃO EXISTE, cria
		orcamentoID2 = uuid.NewString()
		obraID2 := data.ObraIDs[1]
		etapaID2 := data.EtapaIDsPorObra[obraID2][1] // Pega a segunda etapa ("Estrutura")

		_, err = tx.Exec(ctx, queryOrcamento, orcamentoID2, numero2, etapaID2, data.FornecedorIDs[2], 8000.00, "Em Aberto", data.UserIDs[0])
		if err != nil {
			return err
		}

		// Cria o item associado
		_, err = tx.Exec(ctx, queryItem, uuid.NewString(), orcamentoID2, data.ProdutoIDs[3], 100, 80.0) // Tábua de Pinus
		if err != nil {
			return err
		}
		logger.Info("orçamento de exemplo criado", "numero", numero2)
	} else if err != nil {
		return err
	}

	logger.Info("seeder de orçamentos finalizado")
	return tx.Commit(ctx)
}

func seedFuncionarios(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	funcionarios := []map[string]interface{}{
		{"nome": "Carlos Pereira", "cpf": "111.111.111-11", "cargo": "Pedreiro", "valor_diaria": 180.50},
		{"nome": "Mariana Costa", "cpf": "222.222.222-22", "cargo": "Engenheira Civil", "valor_diaria": 350.00},
		{"nome": "Roberto Silva", "cpf": "333.333.333-33", "cargo": "Eletricista", "valor_diaria": 220.00},
		{"nome": "Juliana Almeida", "cpf": "444.444.444-44", "cargo": "Servente", "valor_diaria": 120.00},
	}

	query := `
		INSERT INTO funcionarios (id, nome, cpf, cargo, valor_diaria, data_contratacao, status)
		VALUES ($1, $2, $3, $4, $5, NOW(), 'Ativo')
		ON CONFLICT (cpf) DO NOTHING RETURNING id
	`

	for _, f := range funcionarios {
		var id string
		err := db.QueryRow(ctx, query, uuid.NewString(), f["nome"], f["cpf"], f["cargo"], f["valor_diaria"]).Scan(&id)

		if err == pgx.ErrNoRows {
			querySelect := `SELECT id FROM funcionarios WHERE cpf = $1`
			db.QueryRow(ctx, querySelect, f["cpf"]).Scan(&id)
		} else if err != nil {
			return err
		}
		data.FuncionarioIDs = append(data.FuncionarioIDs, id)
	}

	logger.Info("seeder de funcionários finalizado", "total_funcionarios", len(data.FuncionarioIDs))
	return nil
}

func seedApontamentos(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger, data *SeedData) error {
	if len(data.FuncionarioIDs) < 2 || len(data.ObraIDs) < 1 {
		logger.Warn("não há funcionários ou obras suficientes para criar apontamentos de exemplo")
		return nil
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO apontamentos_quinzenais (id, funcionario_id, obra_id, periodo_inicio, periodo_fim, dias_trabalhados, diaria, valor_total_calculado, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (funcionario_id, periodo_inicio, periodo_fim) DO NOTHING
	`

	// Apontamento 1 (Pago)
	diaria1 := 180.50
	dias1 := 14
	valor1 := diaria1 * float64(dias1)
	_, err = tx.Exec(ctx, query, uuid.NewString(), data.FuncionarioIDs[0], data.ObraIDs[0], "2025-06-01", "2025-06-15", dias1, diaria1, valor1, "PAGO")
	if err != nil {
		return err
	}

	// Apontamento 2 (Aprovado)
	diaria2 := 350.00
	dias2 := 15
	valor2 := diaria2 * float64(dias2)
	_, err = tx.Exec(ctx, query, uuid.NewString(), data.FuncionarioIDs[1], data.ObraIDs[0], "2025-06-16", "2025-06-30", dias2, diaria2, valor2, "APROVADO_PARA_PAGAMENTO")
	if err != nil {
		return err
	}

	// Apontamento 3 (Em Aberto)
	diaria3 := 180.50
	dias3 := 10
	valor3 := diaria3 * float64(dias3)
	_, err = tx.Exec(ctx, query, uuid.NewString(), data.FuncionarioIDs[0], data.ObraIDs[1], "2025-07-01", "2025-07-15", dias3, diaria3, valor3, "EM_ABERTO")
	if err != nil {
		return err
	}

	logger.Info("seeder de apontamentos finalizado")
	return tx.Commit(ctx)
}

package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/dashboard"
	"github.com/luiszkm/masterCostrutora/internal/service/dashboard/dto"
)

// DashboardQuerierPostgres implementa a interface dashboard.Querier para PostgreSQL
type DashboardQuerierPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NovoDashboardQuerier cria uma nova instância do querier de dashboard
func NovoDashboardQuerier(db *pgxpool.Pool, logger *slog.Logger) dashboard.Querier {
	return &DashboardQuerierPostgres{
		db:     db,
		logger: logger,
	}
}

// ObterFluxoCaixa implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterFluxoCaixa(ctx context.Context, dataInicio, dataFim time.Time) ([]*dto.FluxoCaixaDTO, error) {
	const op = "repository.postgres.dashboard.ObterFluxoCaixa"

	query := `
		WITH periodos AS (
			SELECT date_trunc('month', generate_series($1::date, $2::date, '1 month'::interval)) as periodo
		),
		entradas AS (
			-- ENTRADAS REAIS: Recebimentos de obras (receitas efetivas)
			SELECT 
				date_trunc('month', cr.data_recebimento) as periodo,
				COALESCE(SUM(cr.valor_recebido), 0) as valor
			FROM contas_receber cr
			WHERE cr.data_recebimento BETWEEN $1 AND $2 
				AND cr.status = 'RECEBIDO'
			GROUP BY date_trunc('month', cr.data_recebimento)
		),
		saidas AS (
			-- SAÍDAS REAIS: Pagamentos efetivamente realizados
			SELECT periodo, SUM(valor) as valor FROM (
				-- Pagamentos de funcionários (registros_pagamento)
				SELECT 
					date_trunc('month', rp.data_de_efetivacao) as periodo,
					COALESCE(SUM(rp.valor_calculado), 0) as valor
				FROM registros_pagamento rp
				WHERE rp.data_de_efetivacao BETWEEN $1 AND $2
				GROUP BY date_trunc('month', rp.data_de_efetivacao)
				
				-- Pagamentos de fornecedores (contas a pagar)
				UNION ALL
				SELECT 
				    date_trunc('month', cp.data_pagamento) as periodo,
				    COALESCE(SUM(cp.valor_pago), 0) as valor
				FROM contas_pagar cp
				WHERE cp.data_pagamento BETWEEN $1 AND $2
				    AND cp.status IN ('PAGO', 'PARCIAL')
				    AND cp.data_pagamento IS NOT NULL
				GROUP BY date_trunc('month', cp.data_pagamento)
			) todas_saidas
			GROUP BY periodo
		)
		SELECT 
			p.periodo,
			COALESCE(e.valor, 0) as entradas,
			COALESCE(s.valor, 0) as saidas,
			COALESCE(e.valor, 0) - COALESCE(s.valor, 0) as saldo_liquido
		FROM periodos p
		LEFT JOIN entradas e ON p.periodo = e.periodo
		LEFT JOIN saidas s ON p.periodo = s.periodo
		ORDER BY p.periodo`

	rows, err := q.db.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	fluxos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.FluxoCaixaDTO])
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear fluxo de caixa: %w", op, err)
	}

	return fluxos, nil
}

// ObterFluxoCaixaResumo implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterFluxoCaixaResumo(ctx context.Context, dataInicio, dataFim time.Time) (*dto.FluxoCaixaResumoDTO, error) {
	const op = "repository.postgres.dashboard.ObterFluxoCaixaResumo"

	// Primeiro obter o fluxo detalhado
	fluxos, err := q.ObterFluxoCaixa(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resumo := &dto.FluxoCaixaResumoDTO{
		FluxoPorPeriodo: fluxos,
	}

	// Calcular totais
	for _, fluxo := range fluxos {
		resumo.TotalEntradas += fluxo.Entradas
		resumo.TotalSaidas += fluxo.Saidas
	}
	resumo.SaldoAtual = resumo.TotalEntradas - resumo.TotalSaidas

	// Determinar tendência (comparando últimos 3 meses com 3 anteriores)
	if len(fluxos) >= 6 {
		ultimos3 := fluxos[len(fluxos)-3:]
		anteriores3 := fluxos[len(fluxos)-6 : len(fluxos)-3]
		
		var saldoUltimos, saldoAnteriores float64
		for _, f := range ultimos3 {
			saldoUltimos += f.SaldoLiquido
		}
		for _, f := range anteriores3 {
			saldoAnteriores += f.SaldoLiquido
		}
		
		if saldoUltimos > saldoAnteriores*1.1 {
			resumo.TendenciaMensal = "crescente"
		} else if saldoUltimos < saldoAnteriores*0.9 {
			resumo.TendenciaMensal = "decrescente"
		} else {
			resumo.TendenciaMensal = "estavel"
		}
	} else {
		resumo.TendenciaMensal = "indeterminada"
	}

	return resumo, nil
}

// ObterDistribuicaoDespesas implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterDistribuicaoDespesas(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DistribuicaoDespesasDTO, error) {
	const op = "repository.postgres.dashboard.ObterDistribuicaoDespesas"

	query := `
		WITH despesas_materiais AS (
			SELECT 
				COALESCE(p.categoria, 'Sem Categoria') as categoria,
				SUM(oi.quantidade * oi.valor_unitario) as valor,
				COUNT(oi.id) as quantidade_itens
			FROM orcamentos o
			JOIN orcamento_itens oi ON o.id = oi.orcamento_id
			JOIN produtos p ON oi.produto_id = p.id
			WHERE o.data_aprovacao BETWEEN $1 AND $2
				AND o.status = 'Aprovado'
				AND o.deleted_at IS NULL
				AND p.deleted_at IS NULL
			GROUP BY p.categoria
		),
		despesas_mao_obra AS (
			SELECT 
				'Mão de Obra' as categoria,
				SUM(valor_total_calculado) as valor,
				COUNT(*) as quantidade_itens
			FROM apontamentos_quinzenais
			WHERE created_at BETWEEN $1 AND $2
				AND status = 'Aprovado'
		),
		todas_despesas AS (
			SELECT categoria, valor, quantidade_itens FROM despesas_materiais
			UNION ALL
			SELECT categoria, valor, quantidade_itens FROM despesas_mao_obra
		)
		SELECT 
			categoria,
			SUM(valor) as valor,
			SUM(quantidade_itens) as quantidade_itens
		FROM todas_despesas
		WHERE valor > 0
		GROUP BY categoria
		ORDER BY SUM(valor) DESC`

	rows, err := q.db.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var itens []*dto.DistribuicaoDespesasItemDTO
	for rows.Next() {
		var item dto.DistribuicaoDespesasItemDTO
		err := rows.Scan(&item.Categoria, &item.Valor, &item.QuantidadeItens)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear distribuição de despesas: %w", op, err)
		}
		itens = append(itens, &item)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	distribuicao := &dto.DistribuicaoDespesasDTO{
		Distribuicao: itens,
	}

	// Calcular totais e percentuais
	var totalGasto float64
	for _, item := range itens {
		totalGasto += item.Valor
	}
	distribuicao.TotalGasto = totalGasto

	// Calcular percentuais e identificar maior categoria
	if totalGasto > 0 && len(itens) > 0 {
		distribuicao.MaiorCategoria = itens[0].Categoria
		distribuicao.ValorMaiorCategoria = itens[0].Valor
		
		for _, item := range itens {
			item.Percentual = (item.Valor / totalGasto) * 100
		}
	}

	return distribuicao, nil
}

// ObterProgressoObras implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterProgressoObras(ctx context.Context) (*dto.ProgressoObrasDTO, error) {
	const op = "repository.postgres.dashboard.ObterProgressoObras"

	query := `
		WITH obras_progresso AS (
			SELECT 
				o.id as obra_id,
				o.nome as nome_obra,
				o.status,
				o.data_inicio,
				o.data_fim as data_fim_prevista,
				COUNT(e.id) as etapas_total,
				COUNT(CASE WHEN e.status = 'Concluída' THEN 1 END) as etapas_concluidas,
				CASE 
					WHEN COUNT(e.id) > 0 THEN 
						(COUNT(CASE WHEN e.status = 'Concluída' THEN 1 END)::float / COUNT(e.id)::float) * 100
					ELSE 0 
				END as percentual_concluido
			FROM obras o
			LEFT JOIN etapas e ON o.id = e.obra_id
			WHERE o.deleted_at IS NULL
			GROUP BY o.id, o.nome, o.status, o.data_inicio, o.data_fim
		)
		SELECT 
			obra_id,
			nome_obra,
			percentual_concluido,
			etapas_concluidas,
			etapas_total,
			status,
			data_inicio,
			data_fim_prevista
		FROM obras_progresso
		ORDER BY percentual_concluido DESC, nome_obra`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	obras, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.ProgressoObraItemDTO])
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear progresso das obras: %w", op, err)
	}

	progresso := &dto.ProgressoObrasDTO{
		ProgressoPorObra: obras,
	}

	// Calcular estatísticas agregadas
	var somaPercentual float64
	for _, obra := range obras {
		somaPercentual += obra.PercentualConcluido
		
		switch obra.Status {
		case "Em Andamento":
			progresso.ObrasEmAndamento++
		case "Concluída":
			progresso.ObrasConcluidas++
		}
	}

	progresso.TotalObras = len(obras)
	if progresso.TotalObras > 0 {
		progresso.ProgressoMedio = somaPercentual / float64(progresso.TotalObras)
	}

	return progresso, nil
}

// ObterDistribuicaoObras implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterDistribuicaoObras(ctx context.Context) (*dto.DistribuicaoObrasDTO, error) {
	const op = "repository.postgres.dashboard.ObterDistribuicaoObras"

	query := `
		WITH distribuicao AS (
			SELECT 
				o.status,
				COUNT(*) as quantidade,
				-- Estimativa de valor total baseada em orçamentos aprovados
				COALESCE(SUM(orcamentos.valor_total), 0) as valor_total
			FROM obras o
			LEFT JOIN etapas e ON o.id = e.obra_id
			LEFT JOIN orcamentos ON e.id = orcamentos.etapa_id 
				AND orcamentos.status = 'Aprovado' 
				AND orcamentos.deleted_at IS NULL
			WHERE o.deleted_at IS NULL
			GROUP BY o.status
		)
		SELECT 
			status,
			quantidade,
			valor_total
		FROM distribuicao
		ORDER BY quantidade DESC`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var itens []*dto.DistribuicaoObraItemDTO
	for rows.Next() {
		var item dto.DistribuicaoObraItemDTO
		err := rows.Scan(&item.Status, &item.Quantidade, &item.ValorTotal)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear distribuição de obras: %w", op, err)
		}
		itens = append(itens, &item)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	distribuicao := &dto.DistribuicaoObrasDTO{
		DistribuicaoPorStatus: itens,
	}

	// Calcular totais e percentuais
	var totalObras int
	for _, item := range itens {
		totalObras += item.Quantidade
	}
	distribuicao.TotalObras = totalObras

	// Calcular percentuais e identificar status mais comum
	if totalObras > 0 && len(itens) > 0 {
		distribuicao.StatusMaisComum = itens[0].Status
		
		for _, item := range itens {
			item.Percentual = (float64(item.Quantidade) / float64(totalObras)) * 100
		}
	}

	return distribuicao, nil
}

// ObterTendenciasObras implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterTendenciasObras(ctx context.Context, mesesAtras int) (*dto.TendenciasObrasDTO, error) {
	const op = "repository.postgres.dashboard.ObterTendenciasObras"

	dataInicio := time.Now().AddDate(0, -mesesAtras, 0)

	// Query para tendências mensais
	queryTendencias := `
		WITH periodos AS (
			SELECT date_trunc('month', generate_series($1::date, CURRENT_DATE, '1 month'::interval)) as periodo
		)
		SELECT 
			p.periodo,
			COUNT(CASE WHEN date_trunc('month', o.data_inicio) = p.periodo THEN 1 END) as obras_iniciadas,
			COUNT(CASE WHEN date_trunc('month', o.data_fim) = p.periodo THEN 1 END) as obras_concluidas,
			-- Obras em atraso são aquelas com data_fim_prevista no passado mas ainda em andamento
			COUNT(CASE 
				WHEN o.status = 'Em Andamento' 
					AND o.data_fim < CURRENT_DATE 
					AND date_trunc('month', o.data_fim) = p.periodo 
				THEN 1 
			END) as obras_em_atraso
		FROM periodos p
		LEFT JOIN obras o ON (
			date_trunc('month', o.data_inicio) = p.periodo OR 
			date_trunc('month', o.data_fim) = p.periodo
		) AND o.deleted_at IS NULL
		GROUP BY p.periodo
		ORDER BY p.periodo`

	rows, err := q.db.Query(ctx, queryTendencias, dataInicio)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	tendenciasMensais, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.TendenciaObraItemDTO])
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear tendências das obras: %w", op, err)
	}

	// Query para estatísticas gerais de atraso
	queryAtraso := `
		SELECT 
			COUNT(CASE WHEN status = 'Em Andamento' AND data_fim < CURRENT_DATE THEN 1 END) as obras_em_atraso,
			COUNT(CASE WHEN status = 'Em Andamento' AND data_fim >= CURRENT_DATE THEN 1 END) as obras_no_prazo,
			COUNT(CASE WHEN status IN ('Em Andamento', 'Concluída') AND data_fim BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days' THEN 1 END) as previsao_conclusao_mes
		FROM obras 
		WHERE deleted_at IS NULL`

	var obrasEmAtraso, obrasNoPrazo, previsaoConclusaoMes int
	err = q.db.QueryRow(ctx, queryAtraso).Scan(&obrasEmAtraso, &obrasNoPrazo, &previsaoConclusaoMes)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter estatísticas de atraso: %w", op, err)
	}

	tendencias := &dto.TendenciasObrasDTO{
		ObrasEmAtraso:         obrasEmAtraso,
		ObrasNoPrazo:          obrasNoPrazo,
		TendenciaMensal:       tendenciasMensais,
		PrevisaoConclusaoMes:  previsaoConclusaoMes,
	}

	// Calcular percentual de atraso
	totalObrasAvaliadas := obrasEmAtraso + obrasNoPrazo
	if totalObrasAvaliadas > 0 {
		tendencias.PercentualAtraso = (float64(obrasEmAtraso) / float64(totalObrasAvaliadas)) * 100
	}

	// Determinar tendência geral (comparando últimos 3 meses com anteriores)
	if len(tendenciasMensais) >= 6 {
		var atrasosRecentes, atrasosAnteriores int
		for i := len(tendenciasMensais) - 3; i < len(tendenciasMensais); i++ {
			atrasosRecentes += tendenciasMensais[i].ObrasEmAtraso
		}
		for i := len(tendenciasMensais) - 6; i < len(tendenciasMensais) - 3; i++ {
			atrasosAnteriores += tendenciasMensais[i].ObrasEmAtraso
		}

		if atrasosRecentes < atrasosAnteriores {
			tendencias.TendenciaGeral = "melhorando"
		} else if atrasosRecentes > atrasosAnteriores {
			tendencias.TendenciaGeral = "piorando"
		} else {
			tendencias.TendenciaGeral = "estavel"
		}
	} else {
		tendencias.TendenciaGeral = "indeterminada"
	}

	return tendencias, nil
}

// ObterProdutividadeFuncionarios implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterProdutividadeFuncionarios(ctx context.Context) (*dto.ProdutividadeFuncionariosDTO, error) {
	const op = "repository.postgres.dashboard.ObterProdutividadeFuncionarios"

	query := `
		WITH produtividade AS (
			SELECT 
				f.id as funcionario_id,
				f.nome as nome_funcionario,
				f.cargo,
				f.status,
				COUNT(aq.id) as periodos_trabalho,
				SUM(aq.dias_trabalhados) as dias_trabalhados,
				AVG(aq.dias_trabalhados) as media_dias_por_periodo,
				COUNT(DISTINCT al.obra_id) as obras_alocadas
			FROM funcionarios f
			LEFT JOIN apontamentos_quinzenais aq ON f.id = aq.funcionario_id 
				AND aq.created_at >= CURRENT_DATE - INTERVAL '6 months'
			LEFT JOIN alocacoes al ON f.id = al.funcionario_id
			GROUP BY f.id, f.nome, f.cargo, f.status
		)
		SELECT 
			funcionario_id,
			nome_funcionario,
			cargo,
			COALESCE(dias_trabalhados, 0) as dias_trabalhados,
			COALESCE(media_dias_por_periodo, 0) as media_dias_por_periodo,
			COALESCE(obras_alocadas, 0) as obras_alocadas
		FROM produtividade
		WHERE status = 'Ativo'
		ORDER BY media_dias_por_periodo DESC, dias_trabalhados DESC`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var funcionarios []*dto.ProdutividadeFuncionarioItemDTO
	for rows.Next() {
		var funcionario dto.ProdutividadeFuncionarioItemDTO
		err := rows.Scan(&funcionario.FuncionarioID, &funcionario.NomeFuncionario, &funcionario.Cargo,
			&funcionario.DiasTrabalhados, &funcionario.MediaDiasPorPeriodo, &funcionario.ObrasAlocadas)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear produtividade dos funcionários: %w", op, err)
		}
		funcionarios = append(funcionarios, &funcionario)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	produtividade := &dto.ProdutividadeFuncionariosDTO{
		ProdutividadePorFuncionario: funcionarios,
	}

	// Calcular estatísticas agregadas
	var somaMediaDias float64
	var funcionariosAtivos int

	for _, f := range funcionarios {
		somaMediaDias += f.MediaDiasPorPeriodo
		funcionariosAtivos++
		
		// Calcular índice de produtividade (normalizado de 0 a 100)
		// Considerando que 15 dias por quinzena é 100%
		f.IndiceProdutividade = (f.MediaDiasPorPeriodo / 15.0) * 100
		if f.IndiceProdutividade > 100 {
			f.IndiceProdutividade = 100
		}
	}

	produtividade.TotalFuncionarios = len(funcionarios)
	produtividade.FuncionariosAtivos = funcionariosAtivos
	if funcionariosAtivos > 0 {
		produtividade.MediaGeralProdutividade = somaMediaDias / float64(funcionariosAtivos)
	}

	// Top 5 mais produtivos
	top5Count := 5
	if len(funcionarios) < 5 {
		top5Count = len(funcionarios)
	}
	produtividade.Top5Produtivos = funcionarios[:top5Count]

	return produtividade, nil
}

// ObterCustosMaoObra implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterCustosMaoObra(ctx context.Context, dataInicio, dataFim time.Time) (*dto.CustosMaoObraDTO, error) {
	const op = "repository.postgres.dashboard.ObterCustosMaoObra"

	// Custos por funcionário
	queryFuncionarios := `
		SELECT 
			f.id as funcionario_id,
			f.nome as nome_funcionario,
			f.cargo,
			f.valor_diaria,
			COUNT(aq.id) as periodos_trabalho,
			SUM(aq.valor_total_calculado) as custo_total,
			AVG(aq.valor_total_calculado) as custo_medio
		FROM funcionarios f
		LEFT JOIN apontamentos_quinzenais aq ON f.id = aq.funcionario_id 
			AND aq.created_at BETWEEN $1 AND $2
		GROUP BY f.id, f.nome, f.cargo, f.valor_diaria
		HAVING SUM(aq.valor_total_calculado) > 0
		ORDER BY SUM(aq.valor_total_calculado) DESC`

	rows, err := q.db.Query(ctx, queryFuncionarios, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	custosFuncionarios, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.CustoMaoObraFuncionarioDTO])
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear custos por funcionário: %w", op, err)
	}

	// Custos por obra
	queryObras := `
		SELECT 
			o.id as obra_id,
			o.nome as nome_obra,
			COUNT(DISTINCT aq.funcionario_id) as num_funcionarios,
			SUM(aq.valor_total_calculado) as custo_total,
			AVG(aq.valor_total_calculado) as custo_medio
		FROM obras o
		JOIN apontamentos_quinzenais aq ON o.id = aq.obra_id
		WHERE aq.created_at BETWEEN $1 AND $2
			AND o.deleted_at IS NULL
		GROUP BY o.id, o.nome
		HAVING SUM(aq.valor_total_calculado) > 0
		ORDER BY SUM(aq.valor_total_calculado) DESC`

	rowsObras, err := q.db.Query(ctx, queryObras, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsObras.Close()

	custosObras, err := pgx.CollectRows(rowsObras, pgx.RowToAddrOfStructByName[dto.CustoMaoObraPorObraDTO])
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear custos por obra: %w", op, err)
	}

	custos := &dto.CustosMaoObraDTO{
		CustosPorFuncionario: custosFuncionarios,
		CustosPorObra:        custosObras,
	}

	// Calcular totais
	var custoTotal float64
	for _, cf := range custosFuncionarios {
		custoTotal += cf.CustoTotal
	}
	custos.CustoTotalMaoObra = custoTotal

	if len(custosFuncionarios) > 0 {
		custos.CustoMedioFuncionario = custoTotal / float64(len(custosFuncionarios))
	}

	var custoTotalObras float64
	for _, co := range custosObras {
		custoTotalObras += co.CustoTotal
	}
	if len(custosObras) > 0 {
		custos.CustoMedioObra = custoTotalObras / float64(len(custosObras))
	}

	return custos, nil
}

// ObterTopFuncionarios implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterTopFuncionarios(ctx context.Context, limite int) (*dto.TopFuncionariosDTO, error) {
	const op = "repository.postgres.dashboard.ObterTopFuncionarios"

	query := `
		SELECT 
			f.id as funcionario_id,
			f.nome as nome_funcionario,
			f.cargo,
			f.avaliacao_desempenho,
			f.data_contratacao,
			COUNT(DISTINCT aq.id) as dias_trabalhados_total,
			COUNT(DISTINCT al.obra_id) as obras_participadas
		FROM funcionarios f
		LEFT JOIN apontamentos_quinzenais aq ON f.id = aq.funcionario_id
		LEFT JOIN alocacoes al ON f.id = al.funcionario_id
		WHERE f.status = 'Ativo'
			AND f.avaliacao_desempenho IS NOT NULL 
			AND f.avaliacao_desempenho != ''
		GROUP BY f.id, f.nome, f.cargo, f.avaliacao_desempenho, f.data_contratacao
		ORDER BY 
			-- Critério de avaliação: tempo de empresa + produtividade
			((CURRENT_DATE - f.data_contratacao)::integer / 365.0) * 0.3 +
			COUNT(DISTINCT aq.id) * 0.7 DESC
		LIMIT $1`

	rows, err := q.db.Query(ctx, query, limite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var funcionarios []*dto.TopFuncionarioDTO
	for rows.Next() {
		var funcionario dto.TopFuncionarioDTO
		err := rows.Scan(&funcionario.FuncionarioID, &funcionario.NomeFuncionario, &funcionario.Cargo,
			&funcionario.AvaliacaoDesempenho, &funcionario.DataContratacao, 
			&funcionario.DiasTrabalhadosTotal, &funcionario.ObrasParticipadas)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear top funcionários: %w", op, err)
		}
		funcionarios = append(funcionarios, &funcionario)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	// Calcular nota de avaliação baseada no texto (simulação)
	for _, f := range funcionarios {
		// Aqui você pode implementar uma lógica mais sofisticada
		// baseada no conteúdo da avaliação textual
		switch strings.ToLower(f.AvaliacaoDesempenho) {
		case "excelente":
			f.NotaAvaliacao = 9.5
		case "muito bom":
			f.NotaAvaliacao = 8.5
		case "bom":
			f.NotaAvaliacao = 7.5
		case "regular":
			f.NotaAvaliacao = 6.0
		default:
			f.NotaAvaliacao = 7.0 // valor padrão
		}
	}

	return &dto.TopFuncionariosDTO{
		Top5Funcionarios:  funcionarios,
		CriterioAvaliacao: "Tempo de empresa + Produtividade",
	}, nil
}

// ObterFornecedoresPorCategoria implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterFornecedoresPorCategoria(ctx context.Context) (*dto.FornecedoresPorCategoriaDTO, error) {
	const op = "repository.postgres.dashboard.ObterFornecedoresPorCategoria"

	query := `
		SELECT 
			c.id as categoria_id,
			c.nome as categoria_nome,
			COUNT(DISTINCT f.id) as quantidade_fornecedores,
			COALESCE(AVG(f.avaliacao), 0.0) as avaliacao_media
		FROM categorias c
		LEFT JOIN fornecedor_categorias fc ON c.id = fc.categoria_id
		LEFT JOIN fornecedores f ON fc.fornecedor_id = f.id 
			AND f.deleted_at IS NULL 
			AND f.status = 'Ativo'
		GROUP BY c.id, c.nome
		HAVING COUNT(DISTINCT f.id) > 0
		ORDER BY COUNT(DISTINCT f.id) DESC`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var categorias []*dto.FornecedorPorCategoriaItemDTO
	for rows.Next() {
		var categoria dto.FornecedorPorCategoriaItemDTO
		err := rows.Scan(&categoria.CategoriaID, &categoria.CategoriaNome, &categoria.QuantidadeFornecedores, &categoria.AvaliacaoMedia)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear fornecedores por categoria: %w", op, err)
		}
		categorias = append(categorias, &categoria)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	distribuicao := &dto.FornecedoresPorCategoriaDTO{
		DistribuicaoPorCategoria: categorias,
		TotalCategorias:          len(categorias),
	}

	// Calcular totais e percentuais
	var totalFornecedores int
	var melhorAvaliacao float64
	for _, cat := range categorias {
		totalFornecedores += cat.QuantidadeFornecedores
		if cat.AvaliacaoMedia > melhorAvaliacao {
			melhorAvaliacao = cat.AvaliacaoMedia
			distribuicao.CategoriaComMelhorAvaliacao = cat.CategoriaNome
		}
	}
	distribuicao.TotalFornecedores = totalFornecedores

	// Calcular percentuais e identificar categoria mais popular
	if totalFornecedores > 0 && len(categorias) > 0 {
		distribuicao.CategoriaMaisPopular = categorias[0].CategoriaNome
		
		for _, cat := range categorias {
			cat.Percentual = (float64(cat.QuantidadeFornecedores) / float64(totalFornecedores)) * 100
		}
	}

	return distribuicao, nil
}

// ObterTopFornecedores implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterTopFornecedores(ctx context.Context, limite int) (*dto.TopFornecedoresDTO, error) {
	const op = "repository.postgres.dashboard.ObterTopFornecedores"

	query := `
		SELECT 
			f.id as fornecedor_id,
			f.nome as nome_fornecedor,
			f.cnpj,
			COALESCE(f.avaliacao, 0.0) as avaliacao,
			f.status,
			COUNT(o.id) as total_orcamentos,
			SUM(o.valor_total) as valor_total_gasto,
			MAX(o.data_emissao) as ultimo_orcamento
		FROM fornecedores f
		LEFT JOIN orcamentos o ON f.id = o.fornecedor_id 
			AND o.deleted_at IS NULL
		WHERE f.deleted_at IS NULL 
			AND f.status = 'Ativo'
			AND f.avaliacao IS NOT NULL
		GROUP BY f.id, f.nome, f.cnpj, f.avaliacao, f.status
		ORDER BY f.avaliacao DESC, SUM(o.valor_total) DESC
		LIMIT $1`

	rows, err := q.db.Query(ctx, query, limite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var fornecedores []*dto.TopFornecedorDTO
	for rows.Next() {
		var fornecedor dto.TopFornecedorDTO
		err := rows.Scan(&fornecedor.FornecedorID, &fornecedor.NomeFornecedor, &fornecedor.CNPJ,
			&fornecedor.Avaliacao, &fornecedor.Status, &fornecedor.TotalOrcamentos,
			&fornecedor.ValorTotalGasto, &fornecedor.UltimoOrcamento)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear top fornecedores: %w", op, err)
		}
		fornecedores = append(fornecedores, &fornecedor)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	// Buscar categorias para cada fornecedor
	for _, f := range fornecedores {
		queryCategorias := `
			SELECT c.nome 
			FROM categorias c
			JOIN fornecedor_categorias fc ON c.id = fc.categoria_id
			WHERE fc.fornecedor_id = $1`
		
		rowsCat, err := q.db.Query(ctx, queryCategorias, f.FornecedorID)
		if err != nil {
			continue // Não falha se não conseguir buscar categorias
		}

		categorias, err := pgx.CollectRows(rowsCat, pgx.RowTo[string])
		rowsCat.Close()
		if err == nil {
			f.Categorias = categorias
		}
	}

	// Calcular estatísticas
	var somaAvaliacoes float64
	var fornecedoresAtivos int
	
	// Contar total de fornecedores ativos
	err = q.db.QueryRow(ctx, `SELECT COUNT(*) FROM fornecedores WHERE status = 'Ativo' AND deleted_at IS NULL`).Scan(&fornecedoresAtivos)
	if err != nil {
		fornecedoresAtivos = len(fornecedores) // fallback
	}

	for _, f := range fornecedores {
		somaAvaliacoes += f.Avaliacao
	}

	var avaliacaoMedia float64
	if len(fornecedores) > 0 {
		avaliacaoMedia = somaAvaliacoes / float64(len(fornecedores))
	}

	return &dto.TopFornecedoresDTO{
		Top5Fornecedores:   fornecedores,
		CriterioAvaliacao:  "Avaliação + Volume de negócios",
		AvaliacaoMedia:     avaliacaoMedia,
		FornecedoresAtivos: fornecedoresAtivos,
	}, nil
}

// ObterGastosFornecedores implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterGastosFornecedores(ctx context.Context, dataInicio, dataFim time.Time, limite int) (*dto.GastosFornecedoresDTO, error) {
	const op = "repository.postgres.dashboard.ObterGastosFornecedores"

	query := `
		SELECT 
			f.id as fornecedor_id,
			f.nome as nome_fornecedor,
			COALESCE(f.avaliacao, 0.0) as avaliacao,
			COUNT(o.id) as quantidade_orcamentos,
			SUM(o.valor_total) as valor_total_gasto,
			MAX(o.data_emissao) as ultimo_orcamento
		FROM fornecedores f
		JOIN orcamentos o ON f.id = o.fornecedor_id
		WHERE o.data_emissao BETWEEN $1 AND $2
			AND o.deleted_at IS NULL
			AND f.deleted_at IS NULL
			AND o.status = 'Aprovado'
		GROUP BY f.id, f.nome, f.avaliacao
		HAVING SUM(o.valor_total) > 0
		ORDER BY SUM(o.valor_total) DESC
		LIMIT $3`

	rows, err := q.db.Query(ctx, query, dataInicio, dataFim, limite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var fornecedores []*dto.GastoFornecedorItemDTO
	for rows.Next() {
		var fornecedor dto.GastoFornecedorItemDTO
		err := rows.Scan(&fornecedor.FornecedorID, &fornecedor.NomeFornecedor, &fornecedor.Avaliacao, 
			&fornecedor.QuantidadeOrcamentos, &fornecedor.ValorTotalGasto, &fornecedor.UltimoOrcamento)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear gastos com fornecedores: %w", op, err)
		}
		fornecedores = append(fornecedores, &fornecedor)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre os resultados: %w", op, err)
	}

	gastos := &dto.GastosFornecedoresDTO{
		Top10Gastos: fornecedores,
	}

	// Calcular totais
	var totalGasto float64
	for _, f := range fornecedores {
		totalGasto += f.ValorTotalGasto
		// Calcular ticket médio
		if f.QuantidadeOrcamentos > 0 {
			f.TicketMedio = f.ValorTotalGasto / float64(f.QuantidadeOrcamentos)
		}
	}
	gastos.TotalGastoFornecedores = totalGasto

	if len(fornecedores) > 0 {
		gastos.GastoMedioFornecedor = totalGasto / float64(len(fornecedores))
		gastos.FornecedorMaiorGasto = fornecedores[0].NomeFornecedor
		gastos.ValorMaiorGasto = fornecedores[0].ValorTotalGasto
		
		// Calcular percentuais
		for _, f := range fornecedores {
			if totalGasto > 0 {
				f.Percentual = (f.ValorTotalGasto / totalGasto) * 100
			}
		}
	}

	return gastos, nil
}

// ObterEstatisticasGeraisFornecedores implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterEstatisticasGeraisFornecedores(ctx context.Context) (*dto.EstatisticasGeraisFornecedoresDTO, error) {
	const op = "repository.postgres.dashboard.ObterEstatisticasGeraisFornecedores"

	query := `
		SELECT 
			COUNT(*) as total_fornecedores,
			COUNT(CASE WHEN status = 'Ativo' THEN 1 END) as fornecedores_ativos,
			COUNT(CASE WHEN status = 'Inativo' THEN 1 END) as fornecedores_inativos,
			AVG(CASE WHEN avaliacao IS NOT NULL THEN avaliacao END) as avaliacao_media_geral,
			AVG(EXTRACT(days FROM CURRENT_DATE - created_at)) as tempo_medio_contrato
		FROM fornecedores 
		WHERE deleted_at IS NULL`

	var stats dto.EstatisticasGeraisFornecedoresDTO
	err := q.db.QueryRow(ctx, query).Scan(
		&stats.TotalFornecedores,
		&stats.FornecedoresAtivos,
		&stats.FornecedoresInativos,
		&stats.AvaliacaoMediaGeral,
		&stats.TempoMedioContrato,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &stats, nil
}

// ObterResumoGeral implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterResumoGeral(ctx context.Context) (*dto.ResumoGeralDTO, error) {
	const op = "repository.postgres.dashboard.ObterResumoGeral"

	query := `
		SELECT 
			(SELECT COUNT(*) FROM obras WHERE deleted_at IS NULL) as total_obras,
			(SELECT COUNT(*) FROM obras WHERE status = 'Em Andamento' AND deleted_at IS NULL) as obras_em_andamento,
			(SELECT COUNT(*) FROM funcionarios WHERE status = 'Ativo') as funcionarios_ativos,
			(SELECT COUNT(*) FROM funcionarios) as total_funcionarios,
			(SELECT COUNT(*) FROM fornecedores WHERE status = 'Ativo' AND deleted_at IS NULL) as fornecedores_ativos,
			(SELECT COUNT(*) FROM fornecedores WHERE deleted_at IS NULL) as total_fornecedores,
			(SELECT COALESCE(SUM(valor_total), 0) FROM orcamentos WHERE status = 'Aprovado' AND deleted_at IS NULL) as total_investido,
			(SELECT COUNT(*) FROM obras WHERE status = 'Em Andamento' AND data_fim < CURRENT_DATE AND deleted_at IS NULL) as obras_em_atraso`

	var resumo dto.ResumoGeralDTO
	err := q.db.QueryRow(ctx, query).Scan(
		&resumo.TotalObras,
		&resumo.ObrasEmAndamento,
		&resumo.FuncionariosAtivos,
		&resumo.TotalFuncionarios,
		&resumo.FornecedoresAtivos,
		&resumo.TotalFornecedores,
		&resumo.TotalInvestido,
		&resumo.ObrasEmAtraso,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Calcular percentual de atraso
	if resumo.ObrasEmAndamento > 0 {
		resumo.PercentualAtraso = (float64(resumo.ObrasEmAtraso) / float64(resumo.ObrasEmAndamento)) * 100
	}

	// Calcular progresso médio (simulação simples)
	var progressoMedio float64
	err = q.db.QueryRow(ctx, `
		SELECT AVG(
			CASE 
				WHEN etapas_total.total > 0 THEN 
					(etapas_concluidas.concluidas::float / etapas_total.total::float) * 100
				ELSE 0 
			END
		)
		FROM obras o
		LEFT JOIN (
			SELECT obra_id, COUNT(*) as total
			FROM etapas 
			GROUP BY obra_id
		) etapas_total ON o.id = etapas_total.obra_id
		LEFT JOIN (
			SELECT obra_id, COUNT(*) as concluidas
			FROM etapas 
			WHERE status = 'Concluída'
			GROUP BY obra_id
		) etapas_concluidas ON o.id = etapas_concluidas.obra_id
		WHERE o.deleted_at IS NULL AND o.status = 'Em Andamento'`).Scan(&progressoMedio)
	
	if err == nil {
		resumo.ProgressoMedioObras = progressoMedio
	}

	// Calcular saldo financeiro (simulação - entradas menos saídas dos últimos 30 dias)
	var saldoFinanceiro float64
	err = q.db.QueryRow(ctx, `
		SELECT 
			COALESCE(SUM(
				CASE 
					WHEN o.status = 'Aprovado' THEN o.valor_total 
					ELSE 0 
				END
			), 0) - 
			COALESCE((
				SELECT SUM(valor_total_calculado) 
				FROM apontamentos_quinzenais 
				WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
					AND status = 'Aprovado'
			), 0)
		FROM orcamentos o
		WHERE o.data_aprovacao >= CURRENT_DATE - INTERVAL '30 days'
			AND o.deleted_at IS NULL`).Scan(&saldoFinanceiro)
	
	if err == nil {
		resumo.SaldoFinanceiroAtual = saldoFinanceiro
	}

	return &resumo, nil
}

// ObterAlertas implementa dashboard.Querier
func (q *DashboardQuerierPostgres) ObterAlertas(ctx context.Context) (*dto.AlertasDTO, error) {
	const op = "repository.postgres.dashboard.ObterAlertas"

	alertas := &dto.AlertasDTO{}

	// Obras com atraso
	queryObrasAtraso := `
		SELECT nome 
		FROM obras 
		WHERE status = 'Em Andamento' 
			AND data_fim < CURRENT_DATE 
			AND deleted_at IS NULL
		ORDER BY data_fim
		LIMIT 10`

	rows, err := q.db.Query(ctx, queryObrasAtraso)
	if err == nil {
		obrasAtraso, _ := pgx.CollectRows(rows, pgx.RowTo[string])
		alertas.ObrasComAtraso = obrasAtraso
		rows.Close()
	}

	// Fornecedores inativos recentemente
	queryFornecedoresInativos := `
		SELECT nome 
		FROM fornecedores 
		WHERE status = 'Inativo' 
			AND updated_at >= CURRENT_DATE - INTERVAL '30 days'
			AND deleted_at IS NULL
		LIMIT 10`

	rows, err = q.db.Query(ctx, queryFornecedoresInativos)
	if err == nil {
		fornecedoresInativos, _ := pgx.CollectRows(rows, pgx.RowTo[string])
		alertas.FornecedoresInativos = fornecedoresInativos
		rows.Close()
	}

	// Funcionários sem apontamento recente
	queryFuncionariosSemApontamento := `
		SELECT f.nome 
		FROM funcionarios f
		LEFT JOIN apontamentos_quinzenais aq ON f.id = aq.funcionario_id 
			AND aq.created_at >= CURRENT_DATE - INTERVAL '30 days'
		WHERE f.status = 'Ativo' 
			AND aq.id IS NULL
		LIMIT 10`

	rows, err = q.db.Query(ctx, queryFuncionariosSemApontamento)
	if err == nil {
		funcionariosSemApontamento, _ := pgx.CollectRows(rows, pgx.RowTo[string])
		alertas.FuncionariosSemApontamento = funcionariosSemApontamento
		rows.Close()
	}

	// Orçamentos pendentes
	err = q.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM orcamentos 
		WHERE status = 'Em Aberto' 
			AND deleted_at IS NULL`).Scan(&alertas.OrcamentosPendentes)

	// Pagamentos pendentes (apontamentos aprovados mas não pagos)
	err = q.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM apontamentos_quinzenais 
		WHERE status = 'Aprovado' 
			AND id NOT IN (
				SELECT DISTINCT apontamento_id 
				FROM registros_pagamento 
				WHERE apontamento_id IS NOT NULL
			)`).Scan(&alertas.PagamentosPendentes)

	return alertas, nil
}
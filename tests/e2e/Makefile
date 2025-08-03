# Makefile para testes E2E do dashboard

.PHONY: test-e2e test-dashboard test-connection clean-test

# Executa todos os testes E2E
test-e2e:
	@echo "=== Executando testes E2E do Dashboard ==="
	go test -v ./tests/e2e/... -timeout=30s

# Executa apenas testes do dashboard
test-dashboard:
	@echo "=== Testando seções do Dashboard ==="
	go test -v ./tests/e2e/ -run TestDashboardSections -timeout=30s

# Testa conexão com banco
test-connection:
	@echo "=== Testando conexão com banco ==="
	go test -v ./tests/e2e/ -run TestConnection -timeout=10s

# Executa testes com mais detalhes
test-debug:
	@echo "=== Debug dos testes E2E ==="
	go test -v ./tests/e2e/... -timeout=60s -count=1

# Limpa cache de testes
clean-test:
	go clean -testcache
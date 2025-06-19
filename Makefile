# file: Makefile

# Carrega as variáveis do arquivo .env e as exporta para o ambiente deste Makefile
# Isso garante que qualquer comando executado pelo make tenha acesso a elas.
include .env
export

.PHONY: test

## test: Roda todos os testes da aplicação em modo verbose.
test:
	@echo "==> Rodando testes..."
	@go test -v ./...

## up: Sobe o contêiner do banco de dados em background.
up:
    @echo "==> Iniciando contêineres..."
    @docker-compose up -d

## down: Para e remove os contêineres.
down:
    @echo "==> Parando contêineres..."
    @docker-compose down

## down-v: Para os contêineres e apaga os volumes do banco (reset).
down-v:
    @echo "==> Parando contêineres e resetando o banco de dados..."
    @docker-compose down -v
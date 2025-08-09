
# --- Estágio 1: Builder ---
# Usamos uma imagem oficial do Go para compilar nosso código.
FROM golang:1.23-alpine as builder

# Define o diretório de trabalho dentro do contêiner.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências.
# Isso aproveita o cache do Docker. As dependências só serão baixadas novamente
# se o go.mod ou go.sum mudarem.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código fonte.
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 cria um binário estático, que não depende de bibliotecas C.
# GOOS=linux garante que o binário é compilado para o ambiente Linux do contêiner.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go


# --- Estágio 2: Final ---
# Usamos uma imagem mínima para a versão final, apenas com o necessário para rodar.
# 'alpine' é pequena e contém um shell, útil para depuração.
FROM alpine:latest

# Copia apenas o binário compilado do estágio 'builder'.
# Isso resulta em uma imagem final muito menor e mais segura.
COPY --from=builder /server /server

# Expõe a porta que nosso servidor HTTP escuta.
EXPOSE 8080

# O comando para executar quando o contêiner iniciar.
CMD ["/server"]
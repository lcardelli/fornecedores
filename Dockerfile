# Estágio de build
FROM golang:1.22.3 AS builder

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Copie os arquivos de dependência
COPY go.mod go.sum ./
# {{ edit_1 }}
RUN go mod download || go get -u ./...
# Copie todo o código fonte
COPY . .

# Compile o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Estágio final
FROM alpine:latest

# Instale dependências necessárias
RUN apk add --no-cache ca-certificates tzdata

# Defina o diretório de trabalho
WORKDIR /app

# Copie o binário compilado e os arquivos estáticos
COPY --from=builder /app/main /app/
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

# Garanta que o binário tem permissão de execução
RUN chmod +x /app/main

# Exponha a porta que o aplicativo usa
EXPOSE 8080

# Use o comando completo com o caminho absoluto
CMD ["/app/main"]

COPY .env .
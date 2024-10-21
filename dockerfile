FROM golang:1.21-alpine AS builder

# Instala as dependências necessárias
RUN apk add --no-cache git

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixa todas as dependências
RUN go mod download

# Copia o código-fonte
COPY . .

# Compila o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Começa uma nova etapa para criar uma imagem mínima
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o executável compilado do estágio anterior
COPY --from=builder /app/main .

# Copia os arquivos estáticos e templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Expõe a porta que o aplicativo usa
EXPOSE 8080

# Comando para executar o aplicativo
CMD ["./main"]



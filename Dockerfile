# Etapa 1: Constructor
# Usamos 'latest' para cumplir el requisito de versión >= 1.24
FROM golang:latest AS builder

WORKDIR /app

# Descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código
COPY . .

# --- LA MAGIA ESTÁ AQUÍ ---
# 1. CGO_ENABLED=0: Le dice a Go que no use librerías externas de C
# 2. GOOS=linux: Asegura que sea para Linux
# Esto crea un binario "estático" que funciona en cualquier Linux (incluido Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Etapa 2: Ejecutor
FROM alpine:latest

# Instalamos certificados de seguridad (útil si tu app hace peticiones HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiamos el binario desde la etapa 1
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
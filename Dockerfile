# Stage 1: Build stage (Simplest + Cache Clean)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source code (including pre-generated files)
COPY . .

# --- Limpar o cache de módulos Go DEPOIS de copiar tudo ---
RUN go clean -modcache

# Build the application (assuming code is already generated)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server/main.go

# Stage 2: Runtime stage (Mostly unchanged)
FROM alpine:latest

WORKDIR /app

# Install only necessary runtime dependencies (CA certs, psql client, migrate CLI)
RUN apk update && apk add --no-cache ca-certificates postgresql-client curl && \
    echo "Downloading migrate..." && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz && \
    echo "Moving migrate..." && \
    mv migrate /usr/local/bin/migrate && \
    echo "Setting migrate permissions..." && \
    chmod +x /usr/local/bin/migrate && \
    echo "Migrate installed."

# Copy the compiled binary from the build stage
COPY --from=builder /server /app/server
# Copy migrations
COPY ./internal/infra/database/migrations ./migrations

# Expose ports
EXPOSE 8080
EXPOSE 50051

# Environment variables
ENV DATABASE_URL="postgresql://user:password@db:5432/ordersdb?sslmode=disable"
ENV WEB_SERVER_PORT=8080
ENV GRPC_SERVER_PORT=50051
ENV GRAPHQL_SERVER_PATH=/query

# Entrypoint command
CMD ["sh", "-c", "echo 'Waiting for DB...' && sleep 5 && echo 'Running migrations...' && migrate -path /app/migrations -database ${DATABASE_URL} up && echo 'Starting server...' && /app/server"]

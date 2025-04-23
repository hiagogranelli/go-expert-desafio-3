FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go clean -modcache

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates postgresql-client curl && \
    echo "Downloading migrate..." && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz && \
    echo "Moving migrate..." && \
    mv migrate /usr/local/bin/migrate && \
    echo "Setting migrate permissions..." && \
    chmod +x /usr/local/bin/migrate && \
    echo "Migrate installed."

COPY --from=builder /server /app/server
COPY ./internal/infra/database/migrations ./migrations

EXPOSE 8080
EXPOSE 50051

ENV DATABASE_URL="postgresql://user:password@db:5432/ordersdb?sslmode=disable"
ENV WEB_SERVER_PORT=8080
ENV GRPC_SERVER_PORT=50051
ENV GRAPHQL_SERVER_PATH=/query

CMD ["sh", "-c", "echo 'Waiting for DB...' && sleep 5 && echo 'Running migrations...' && migrate -path /app/migrations -database ${DATABASE_URL} up && echo 'Starting server...' && /app/server"]

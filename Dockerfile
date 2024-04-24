FROM golang:1.22.2-bookworm

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /app.bin ./cmd/server/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

EXPOSE 5000

CMD goose -dir db/migrations postgres 'postgresql://postgres:postgres@localhost:5432/db_dashboards?sslmode=disable' up && /app.bin
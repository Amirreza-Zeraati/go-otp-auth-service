FROM golang:1.21 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init

RUN go build -o main .

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/templates ./templates

EXPOSE 3000

CMD ["./main"]

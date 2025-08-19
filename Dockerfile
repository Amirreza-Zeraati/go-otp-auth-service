FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

FROM alpine:3.18

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/docs ./docs
COPY --from=builder --chown=appuser:appuser /app/templates ./templates

USER appuser

EXPOSE 3000

CMD ["./main"]

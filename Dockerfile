FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

FROM alpine:3.18

WORKDIR /app

RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /app/main .
COPY .env .

EXPOSE 3000

CMD ["./main"]

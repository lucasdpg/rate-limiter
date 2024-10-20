FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/rate-limiter ./cmd/server/main.go

FROM alpine:3.18 AS release

WORKDIR /app

COPY --from=builder /app/rate-limiter /app/rate-limiter
COPY .env /app/.env
COPY static/ /app/static/

CMD ["/app/rate-limiter"]

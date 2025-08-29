# Dockerfile para SSE Server em Go

FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o sse-server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sse-server ./sse-server
EXPOSE 8080
ENV AUTH_TOKEN=seu_token_fixo_aqui
CMD ["./sse-server"]

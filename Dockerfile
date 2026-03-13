FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /fsocket ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /fsocket .

EXPOSE 8080

ARG AUTH_TOKEN
ENV AUTH_TOKEN=${AUTH_TOKEN}

CMD ["./fsocket"]

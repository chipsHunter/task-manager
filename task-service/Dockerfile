FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./


RUN apk add --no-cache gcc musl-dev

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal



RUN go build -o task-service ./cmd/task-service

FROM alpine:latest

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

WORKDIR /app

COPY --from=builder /app/task-service .
COPY ./config ./config

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 8082

CMD ["./task-service"]
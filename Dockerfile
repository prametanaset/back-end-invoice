FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY configs ./configs

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080
CMD ["./app"]

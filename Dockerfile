FROM golang:1.21.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /order-service ./cmd/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /order-service .

COPY .env .

EXPOSE 8004

CMD ["./order-service"]

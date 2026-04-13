FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init -g cmd/gateway/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway_bin ./cmd/gateway/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /gateway_bin .

ENV HTTP_PORT=8080
ENV COLLECTOR_ADDR=collector:8888

EXPOSE 8080

CMD ["./gateway_bin"]

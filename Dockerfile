FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o /medods ./cmd/medods/main.go

FROM alpine:latest

COPY --from=builder /medods .

COPY configs/config.yml . 
COPY .env .


EXPOSE 8080

CMD ["./medods"]

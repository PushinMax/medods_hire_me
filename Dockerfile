FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /medods cmd/medods/main.go

FROM alpine:latest

COPY --from=builder /medods /medods

# EXPOSE 8080

CMD ["/medods"]
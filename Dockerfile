FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/api/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .
COPY database.db .
EXPOSE 8080

CMD ["./main"]

# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-server ./cmd

# Final slim image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/auth-server .

EXPOSE 8000

CMD ["./auth-server"]

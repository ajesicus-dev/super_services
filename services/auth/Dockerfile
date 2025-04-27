# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd

# Final slim image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/app .

EXPOSE 8000

CMD ["./app"]

# Build stage - Using Debian-based Go image
FROM golang:latest AS builder

# Install git using apt (Debian/Ubuntu)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Enable Go toolchain auto-download
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/rate-limiter .

EXPOSE 1234

CMD ["./rate-limiter"]
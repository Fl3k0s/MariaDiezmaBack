# Stage 1: Build binary
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api

# Stage 2: Final lightweight image
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy compiled binary from builder
COPY --from=builder /app/api /app/api
COPY --from=builder /app/.env.example /app/.env.example

EXPOSE 8080

ENTRYPOINT ["/app/api"]

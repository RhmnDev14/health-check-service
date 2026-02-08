# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (simplified to reduce memory usage)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Copy web assets
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Run
CMD ["./main"]

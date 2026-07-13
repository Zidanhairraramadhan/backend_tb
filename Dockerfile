# ── Stage 1: Builder ──
FROM golang:1.22-alpine AS builder

# Install required system packages
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy module files first (for better Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build the binary (static, no CGO for Alpine compatibility)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o musiclink-server ./

# ── Stage 2: Final Image ──
FROM alpine:3.19

# Install CA certificates for HTTPS connections (e.g., to Supabase)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/musiclink-server .

# Expose default port
EXPOSE 5000

# Run the binary
CMD ["./musiclink-server"]

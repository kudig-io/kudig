# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o kudig ./cmd/kudig

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    bash \
    curl \
    jq

# Copy binary from builder
COPY --from=builder /app/kudig /usr/local/bin/kudig

# Copy default rules (if exists)
COPY rules/ /app/rules/

# Set environment
ENV PATH="/app:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD kudig --version || exit 1

# Default command
ENTRYPOINT ["kudig"]
CMD ["--help"]

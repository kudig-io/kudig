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

# Build binary with version info
ARG VERSION=2.0.0
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -o kudig ./cmd/kudig

# Prepare rules (ensure dir exists even if empty)
RUN mkdir -p /app/rules-out && cp /app/rules/*.yaml /app/rules-out/ 2>/dev/null || true

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies and create non-root user
RUN apk add --no-cache \
    ca-certificates \
    bash \
    curl \
    jq && \
    addgroup -S kudig && adduser -S kudig -G kudig

# Copy binary from builder
COPY --from=builder /app/kudig /usr/local/bin/kudig

# Copy default rules
COPY --from=builder /app/rules-out/ /app/rules/

# Set environment
ENV PATH="/app:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD kudig --version || exit 1

# Run as non-root user
USER kudig

# Default command
ENTRYPOINT ["kudig"]
CMD ["--help"]

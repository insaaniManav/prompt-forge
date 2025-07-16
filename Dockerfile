# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git, ca-certificates, and build tools (needed for go modules, HTTPS, and CGO/SQLite)
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Copy go mod files
COPY api/go.mod api/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY api/ .

# Build the application with SQLite compatibility
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -tags="sqlite_omit_load_extension" -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS, sqlite for database, and wget for health checks
RUN apk --no-cache add ca-certificates sqlite wget

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy frontend files
COPY frontend/ ./frontend/

# Create directory for database
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DATABASE_PATH=/data/promptforge.db

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --method=GET http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./main"] 

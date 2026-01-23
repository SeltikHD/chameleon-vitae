# =========================================================
# Stage 1: Builder
# =========================================================
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install git and SSL certificates (needed for dependencies)
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the static binary targeting the correct entry point
# -ldflags="-s -w": Strips debug info for smaller size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server/main.go

# =========================================================
# Stage 2: Runner
# =========================================================
FROM alpine:3

WORKDIR /root/

# Install CA certificates to enable HTTPS calls (Firebase/Groq)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Create the secrets directory
RUN mkdir -p /etc/secrets

# Expose the port
EXPOSE 8080

# Start the application
CMD ["./main"]

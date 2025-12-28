# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Debug: List what was copied
RUN ls -la && ls -la cmd/ || echo "cmd/ not found"

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bambustatus ./cmd/bambustatus

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bambustatus .

# Copy web assets
COPY --from=builder /app/web ./web

# Expose HTTP port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./bambustatus"]

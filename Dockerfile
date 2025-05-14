# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o balancer ./cmd/balancer

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/balancer .
COPY --from=builder /app/configs/config.yaml ./configs/

# Create logs directory
RUN mkdir -p logs

EXPOSE 8080

CMD ["./balancer"]
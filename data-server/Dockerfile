FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o data-server ./cmd/data-server

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/data-server .

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./data-server"]

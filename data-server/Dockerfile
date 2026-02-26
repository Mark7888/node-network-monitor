FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Accept version as build argument (defaults to "development")
ARG VERSION=development

# Build application with version
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-X mark7888/speedtest-data-server/internal/version.Version=${VERSION}" -o data-server ./cmd/data-server

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

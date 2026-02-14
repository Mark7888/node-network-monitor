FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o speedtest-node ./cmd/speedtest-node

FROM alpine:latest

# Install speedtest CLI
RUN apk add --no-cache ca-certificates curl && \
    curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.alpine.sh | sh && \
    apk add speedtest

WORKDIR /app
COPY --from=builder /app/speedtest-node .

RUN mkdir -p /app/data /app/logs

VOLUME ["/app/data", "/app/logs"]

CMD ["./speedtest-node"]

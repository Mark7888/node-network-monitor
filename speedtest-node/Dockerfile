FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y gcc libsqlite3-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o speedtest-node ./cmd/speedtest-node/main.go

FROM debian:bookworm-slim

# Install speedtest CLI
RUN apt-get update && \
    apt-get install -y curl ca-certificates && \
    curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash && \
    apt-get install -y speedtest && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/speedtest-node .

RUN mkdir -p /app/data /app/logs

CMD ["./speedtest-node"]

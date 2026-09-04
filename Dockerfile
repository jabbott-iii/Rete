FROM golang:1.26-alpine AS builder

WORKDIR /src

# Install build deps for CGO sqlite3 driver
RUN apk add --no-cache build-base

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/salus .

FROM alpine:3.22

# Runtime deps: sqlite shared lib + CA certs
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app
COPY --from=builder /out/salus /usr/local/bin/salus

# Persist sqlite database file (rete.db)
VOLUME ["/app/data"]

ENV SALUS_DB_PATH=/app/data/salus.db

# This app is an interactive TUI/CLI
ENTRYPOINT ["salus"]

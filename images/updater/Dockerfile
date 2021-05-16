FROM golang:1.16.4-buster AS builder

ENV CGO_ENABLED=0

# build updater
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o updater cmd/updater/main.go

FROM debian:buster-20210511-slim
RUN \
    apt-get update \
    && apt-get install -y --no-install-suggests --no-install-recommends \
        ca-certificates

COPY --from=builder /app/checker /usr/local/bin/bottlerocket-updater


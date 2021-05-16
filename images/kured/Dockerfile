FROM golang:1.16.4-buster AS builder

ENV CGO_ENABLED=0

WORKDIR /app

# build checker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o checker cmd/checker/main.go

# download kured 1.6.1
RUN go install github.com/weaveworks/kured/cmd/kured@56e2c12

FROM debian:buster-20210511-slim
RUN \
    apt-get update \
    && apt-get install -y --no-install-suggests --no-install-recommends \
        ca-certificates

COPY --from=builder /go/bin/kured /usr/local/bin/kured
COPY --from=builder /app/checker /usr/local/bin/bottlerocket-checker

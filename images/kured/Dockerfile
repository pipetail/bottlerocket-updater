FROM golang:1.16.4-buster AS builder

ENV CGO_ENABLED=0

# download kured from commit 42e4c31 since this commit
# supports guard command
RUN go install github.com/weaveworks/kured/cmd/kured@42e4c31

# build checker
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o checker cmd/checker/main.go
RUN go build -o reboot cmd/reboot/main.go

FROM debian:buster-20210511-slim
RUN \
    apt-get update \
    && apt-get install -y --no-install-suggests --no-install-recommends \
        ca-certificates

COPY --from=builder /go/bin/kured /usr/local/bin/kured
COPY --from=builder /app/checker /usr/local/bin/bottlerocket-checker
COPY --from=builder /app/reboot /usr/local/bin/bottlerocket-reboot

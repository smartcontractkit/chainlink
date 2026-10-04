FROM golang:1.26@sha256:6c2a5538f964f1c82f97ad14988bf05de100d922d159d0e398b54c7b0ca0c6c9 AS builder

ENV GOPRIVATE=github.com/smartcontractkit/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ../.. .
RUN CGO_ENABLED=0 GOOS=linux go build -o /fake main.go

FROM alpine:3@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
# Unpinned on purpose: always take the latest patched ca-certificates.
# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates
COPY --from=builder /fake /fake
EXPOSE 9111
CMD ["/fake"]

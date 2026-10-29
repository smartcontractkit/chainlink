FROM golang:1.27@sha256:957861185dbf29d93668ace41e011d091c6e779953d47f0316ef2cc58e28fe2b AS builder

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

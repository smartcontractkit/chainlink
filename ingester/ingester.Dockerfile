FROM golang:1.13-alpine as builder

RUN apk add --no-cache make curl git g++ gcc musl-dev linux-headers

WORKDIR /usr/local/src/chainlink/ingester

# Do go mod download in a cacheable step
ADD ingester/go.mod ingester/go.sum ./
RUN go mod download

# Build the ingester binary
ADD ingester .
RUN go build

# Copy into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/src/chainlink/ingester/ingester /usr/local/bin/

ENTRYPOINT ["ingester"]

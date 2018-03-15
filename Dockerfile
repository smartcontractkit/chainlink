# Build Chainlink
FROM golang:1.10-alpine as builder

RUN apk add --no-cache make curl git gcc musl-dev linux-headers
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD . /go/src/github.com/smartcontractkit/chainlink
RUN cd /go/src/github.com/smartcontractkit/chainlink && make build

# Copy Chainlink into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/smartcontractkit/chainlink/chainlink /usr/local/bin/

EXPOSE 6688
ENTRYPOINT ["chainlink"]

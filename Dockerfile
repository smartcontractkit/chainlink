# Build Chainlink
FROM golang:1.10-alpine as builder

RUN apk add --no-cache make curl git gcc g++ musl-dev linux-headers yarn python
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD . /go/src/github.com/smartcontractkit/chainlink
RUN cd /go/src/github.com/smartcontractkit/chainlink && make build

# Copy Chainlink into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/smartcontractkit/chainlink/chainlink /usr/local/bin/
COPY --from=builder /go/src/github.com/smartcontractkit/chainlink /go/src/github.com/smartcontractkit/chainlink

EXPOSE 6688
EXPOSE 6689
ENTRYPOINT ["chainlink"]

FROM golang:1.12-alpine as builder

WORKDIR /go/src/github.com/smartcontractkit/aggregator-monitor

RUN apk add --no-cache make curl git g++ gcc musl-dev linux-headers
RUN go get -u github.com/gobuffalo/packr/v2/packr2

ENV GO111MODULE=on
ADD . .
RUN go install
RUN go build

# Copy into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/smartcontractkit/aggregator-monitor/aggregator-monitor /usr/local/bin/

ENTRYPOINT ["aggregator-monitor"]

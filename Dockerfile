# Build Chainlink
FROM smartcontract/builder:1.0.2 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ARG COMMIT_SHA
ARG ENVIRONMENT

ADD . /go/src/github.com/smartcontractkit/chainlink
WORKDIR /go/src/github.com/smartcontractkit/chainlink
RUN make build

# Final layer: ubuntu with chainlink binary
FROM ubuntu:16.04

COPY --from=builder \
  /go/src/github.com/smartcontractkit/chainlink/chainlink \
  /root

WORKDIR /root
EXPOSE 6688
EXPOSE 6689
ENTRYPOINT ["./chainlink"]
CMD ["node"]

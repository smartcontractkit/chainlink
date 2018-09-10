# Build Chainlink
FROM smartcontract/builder:1.0.9 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ARG COMMIT_SHA
ARG ENVIRONMENT

ADD . /go/src/github.com/smartcontractkit/chainlink
WORKDIR /go/src/github.com/smartcontractkit/chainlink
RUN make build

# Final layer: ubuntu with chainlink binary
FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder \
  /go/src/github.com/smartcontractkit/chainlink/chainlink \
  /usr/local/bin/

EXPOSE 6688
ENTRYPOINT ["chainlink"]
CMD ["node"]

# Build Chainlink
FROM smartcontract/builder:1.0.2 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ARG COMMIT_SHA
ARG ENVIRONMENT

# Do dependency installs first, since these will change less than the full
# source tree and can get cached
WORKDIR /go/src/github.com/smartcontractkit/chainlink
ADD Gopkg.* /go/src/github.com/smartcontractkit/chainlink/
RUN dep ensure -vendor-only
ADD package.json yarn.lock /go/src/github.com/smartcontractkit/chainlink/
RUN yarn install

ADD . /go/src/github.com/smartcontractkit/chainlink
RUN make chainlink

# Final layer: ubuntu with chainlink binary
FROM ubuntu:16.04

COPY --from=builder \
  /go/src/github.com/smartcontractkit/chainlink/chainlink \
  /usr/local/bin/
COPY --from=builder \
  /go/src/github.com/smartcontractkit/chainlink/gui/dist \
  /go/src/github.com/smartcontractkit/chainlink/gui/dist

EXPOSE 6688
ENTRYPOINT ["chainlink"]
CMD ["node"]

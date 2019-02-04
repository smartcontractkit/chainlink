# Build Chainlink
FROM smartcontract/builder:1.0.16 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ARG COMMIT_SHA
ARG ENVIRONMENT

WORKDIR /go/src/github.com/smartcontractkit/chainlink
ADD . ./
RUN make install

# Final layer: ubuntu with chainlink binary
FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /root

COPY --from=builder /go/bin/chainlink /usr/local/bin/

COPY --from=builder \
  /go/src/github.com/smartcontractkit/chainlink/chainlink-launcher.sh \
  /root/

RUN chmod +x ./chainlink-launcher.sh

EXPOSE 6688
ENTRYPOINT ["./chainlink-launcher.sh"]
CMD ["node"]

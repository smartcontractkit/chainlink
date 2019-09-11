# Build Chainlink
FROM smartcontract/builder:1.0.22 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ARG COMMIT_SHA
ARG ENVIRONMENT

WORKDIR /chainlink
COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Do dep ensure in a cacheable step
ADD go.* ./
RUN go mod download

# And yarn likewise
COPY yarn.lock package.json ./
COPY explorer/client/package.json ./explorer/client/
COPY explorer/package.json ./explorer/
COPY operator_ui/package.json ./operator_ui/
COPY styleguide/package.json ./styleguide/
COPY tools/prettier-config/package.json ./tools/prettier-config/
RUN make yarndep

# Install chainlink
ADD . ./
RUN make install-chainlink

# Final layer: ubuntu with chainlink binary
FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /root

COPY --from=builder /go/bin/chainlink /usr/local/bin/

EXPOSE 6688
ENTRYPOINT ["chainlink"]
CMD ["local", "node"]

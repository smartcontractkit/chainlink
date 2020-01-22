# Build Chainlink
FROM smartcontract/builder:1.0.25 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /chainlink
COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Install yarn dependencies
COPY yarn.lock package.json ./
COPY explorer/client/package.json ./explorer/client/
COPY explorer/package.json ./explorer/
COPY operator_ui/package.json ./operator_ui/
COPY feeds_ui/package.json ./feeds_ui/
COPY styleguide/package.json ./styleguide/
COPY tools/json-api-client/package.json ./tools/json-api-client/
COPY tools/local-storage/package.json ./tools/local-storage/
COPY tools/redux/package.json ./tools/redux/
COPY tools/ts-test-helpers/package.json ./tools/ts-test-helpers/
COPY evm/v0.5/package.json ./evm/v0.5/
COPY evm/package.json ./evm/
RUN make yarndep

# Do go mod download in a cacheable step
ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA
ARG ENVIRONMENT

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

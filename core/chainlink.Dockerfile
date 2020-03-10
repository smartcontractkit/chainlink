# Build Chainlink
FROM smartcontract/builder:1.0.30 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /chainlink
COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Install yarn dependencies
COPY yarn.lock package.json .yarnrc ./
COPY .yarn .yarn
COPY operator_ui/package.json ./operator_ui/
COPY styleguide/package.json ./styleguide/
COPY tools/json-api-client/package.json ./tools/json-api-client/
COPY tools/local-storage/package.json ./tools/local-storage/
COPY tools/redux/package.json ./tools/redux/
COPY tools/ts-helpers/package.json ./tools/ts-helpers/
COPY belt/package.json ./belt/
COPY belt/bin ./belt/bin
COPY evm-test-helpers/package.json ./evm-test-helpers/
COPY evm-contracts/package.json ./evm-contracts/
RUN make yarndep

# Do go mod download in a cacheable step
ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA
ARG ENVIRONMENT

# Install chainlink
COPY tsconfig.cjs.json tsconfig.es6.json ./
COPY operator_ui ./operator_ui
COPY styleguide ./styleguide
COPY tools/json-api-client ./tools/json-api-client
COPY tools/local-storage ./tools/local-storage
COPY tools/redux ./tools/redux
COPY tools/ts-helpers ./tools/ts-helpers
COPY belt ./belt
COPY belt/bin ./belt/bin
COPY evm-test-helpers ./evm-test-helpers
COPY evm-contracts ./evm-contracts
COPY core core
COPY packr packr

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

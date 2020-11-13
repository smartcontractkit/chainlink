# MAKE ALL CHANGES WITHIN THE DEFAULT WORKDIR FOR YARN AND GO DEP CACHE HITS

ARG BUILDER=smartcontract/builder
FROM ${BUILDER}:1.0.39
WORKDIR /chainlink
# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags tools/bin/ldflags
ARG COMMIT_SHA

# Install yarn dependencies
COPY yarn.lock package.json .yarnrc ./
COPY patches patches
COPY solc_bin solc_bin
COPY .yarn .yarn
COPY operator_ui/package.json ./operator_ui/
COPY belt/package.json ./belt/
COPY belt/bin ./belt/bin
COPY evm-test-helpers/package.json ./evm-test-helpers/
COPY evm-contracts/package.json ./evm-contracts/
COPY tools/bin/restore-solc-cache ./tools/bin/restore-solc-cache
RUN make yarndep


COPY tsconfig.cjs.json tsconfig.es6.json ./
COPY operator_ui ./operator_ui
COPY belt ./belt
COPY belt/bin ./belt/bin
COPY evm-test-helpers ./evm-test-helpers
COPY evm-contracts ./evm-contracts

# Build operator-ui and the smart contracts
RUN make contracts-operator-ui-build

# Build the golang binary

FROM ${BUILDER}:1.0.39
WORKDIR /chainlink

# Have to reintroduce ENV vars from builder image
ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Env vars needed for chainlink build
ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA
ARG ENVIRONMENT

COPY --from=0 /chainlink/evm-contracts/abi ./evm-contracts/abi
COPY --from=0 /chainlink/operator_ui/dist ./operator_ui/dist
COPY core core
COPY packr packr

RUN make chainlink-build

# Final layer: ubuntu with chainlink binary
FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates curl

WORKDIR /root

COPY --from=1 /go/bin/chainlink /usr/local/bin/

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688 || exit 1

CMD ["local", "node"]

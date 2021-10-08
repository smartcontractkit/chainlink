# Build Chainlink
FROM smartcontract/builder:1.0.42

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

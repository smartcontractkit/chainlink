# Build Chainlink
FROM smartcontract/builder:1.0.28

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

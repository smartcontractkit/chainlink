# Build the plugin binary
FROM golang:1.20-buster as solana
WORKDIR /chainlink

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

COPY core core

# Build plugin
RUN make chainlink-solana-build

# Final layer: ubuntu with chainlink and solana binaries
FROM smartcontract/chainlink:plugin-base

# Install solana plugin
COPY --from=solana /go/bin/chainlink-solana /usr/local/bin/

CMD ["local", "node"]

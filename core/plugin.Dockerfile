# Build the plugin binary
FROM golang:1.20-buster as build
WORKDIR /chainlink

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

COPY core core

# Build plugins
RUN make chainlink-solana-install
RUN make chainlink-median-install

# Final layer: ubuntu with chainlink and plugin binaries
FROM smartcontract/chainlink:plugin-base

# Install plugins
COPY --from=build /go/bin/chainlink-solana /usr/local/bin/
ENV CL_SOLANA chainlink-solana
COPY --from=build /go/bin/chainlink-median /usr/local/bin/
ENV CL_MEDIAN chainlink-median

CMD ["local", "node"]

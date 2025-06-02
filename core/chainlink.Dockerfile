# Build image: Chainlink binary
FROM golang:1.24-bullseye AS buildgo
RUN go version
WORKDIR /chainlink

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

# Build chainlink bin with cover flag https://go.dev/doc/build-cover#FAQ
ARG GO_COVER_FLAG=false

COPY . .

RUN apt-get update && apt-get install -y jq

# Build the golang binary
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  if [ "$GO_COVER_FLAG" = "true" ]; then \
        make install-chainlink-cover; \
    else \
        make install-chainlink; \
    fi

# Link LOOP Plugin source dirs with simple names
RUN --mount=type=cache,target=/go/pkg/mod \
  mkdir -p /chainlink-feeds && \
  GO_PACKAGE_PATH=$(go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-feeds) && \
  cp -r $GO_PACKAGE_PATH/* /chainlink-feeds/

RUN --mount=type=cache,target=/go/pkg/mod \
  mkdir -p /chainlink-solana && \
  GO_PACKAGE_PATH=$(go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-solana) && \
  cp -r $GO_PACKAGE_PATH/* /chainlink-solana/

# Build image: Plugins
FROM golang:1.24-bullseye AS buildplugins
RUN go version

WORKDIR /chainlink-feeds
COPY --from=buildgo /chainlink-feeds .
RUN go install ./cmd/chainlink-feeds

WORKDIR /chainlink-solana
COPY --from=buildgo /chainlink-solana .
RUN go install ./pkg/solana/cmd/chainlink-solana

# Final image: ubuntu with chainlink binary
FROM ubuntu:24.04

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && rm -rf /var/lib/apt/lists/*

COPY --from=buildgo /go/bin/chainlink /usr/local/bin/

# Install (but don't enable) feeds LOOP Plugin
COPY --from=buildplugins /go/bin/chainlink-feeds /usr/local/bin/
# Install Solana LOOP Plugin
COPY --from=buildplugins /go/bin/chainlink-solana /usr/local/bin/
# Optionally enable the Solana LOOP Plugin
ARG CL_SOLANA_CMD=chainlink-solana
ENV CL_SOLANA_CMD=${CL_SOLANA_CMD}

# CCIP specific
COPY ./cci[p]/confi[g] /ccip-config
ARG CL_CHAIN_DEFAULTS
ENV CL_CHAIN_DEFAULTS=${CL_CHAIN_DEFAULTS}

RUN if [ ${CHAINLINK_USER} != root ]; then \
  useradd --uid 14933 --create-home ${CHAINLINK_USER}; \
  fi
USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}
# explicit set the cache dir. needed so both root and non-root user has an explicit location
ENV XDG_CACHE_HOME=/home/${CHAINLINK_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

# Set up env and dir for go coverage profiling https://go.dev/doc/build-cover#FAQ
ARG GO_COVER_DIR="/var/tmp/go-coverage"
ENV GOCOVERDIR=${GO_COVER_DIR}
RUN mkdir -p $GO_COVER_DIR

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]

# Build image: Chainlink binary
FROM golang:1.22-bullseye as buildgo
RUN go version
WORKDIR /chainlink

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

# Build chainlink bin with cover flag https://go.dev/doc/build-cover#FAQ
ARG GO_COVER_FLAG=false

COPY . .

RUN apt-get update && apt-get install -y jq

# Build the golang binary
RUN if [ "$GO_COVER_FLAG" = "true" ]; then \
        make install-chainlink-cover; \
    else \
        make install-chainlink-delve; \
    fi

# Link LOOP Plugin source dirs with simple names
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-feeds | xargs -I % ln -s % /chainlink-feeds
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-solana | xargs -I % ln -s % /chainlink-solana

# Build image: Plugins
FROM golang:1.22-bullseye as buildplugins
RUN go version

WORKDIR /chainlink-feeds
COPY --from=buildgo /chainlink-feeds .
RUN go install ./cmd/chainlink-feeds

WORKDIR /chainlink-solana
COPY --from=buildgo /chainlink-solana .
RUN go install ./pkg/solana/cmd/chainlink-solana

# Final image: ubuntu with chainlink binary
FROM golang:1.22-bullseye

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && apt-get clean all \
  && rm -rf /var/lib/apt/lists/*

COPY --from=buildgo /go/bin/chainlink /usr/local/bin/

# Install (but don't enable) LOOP Plugins
COPY --from=buildplugins /go/bin/chainlink-feeds /usr/local/bin/
COPY --from=buildplugins /go/bin/chainlink-solana /usr/local/bin/

# Dependency of CosmWasm/wasmd
COPY --from=buildgo /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.*.so /usr/lib/
RUN chmod 755 /usr/lib/libwasmvm.*.so

RUN if [ ${CHAINLINK_USER} != root ]; then \
  useradd --uid 14933 --create-home ${CHAINLINK_USER}; \
  fi
USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}
RUN mkdir -p go
# explicit set the cache dir. needed so both root and non-root user has an explicit location
ENV XDG_CACHE_HOME /home/${CHAINLINK_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

# Set up env and dir for go coverage profiling https://go.dev/doc/build-cover#FAQ
ARG GO_COVER_DIR="/var/tmp/go-coverage"
ENV GOCOVERDIR=${GO_COVER_DIR}
RUN mkdir -p $GO_COVER_DIR

# Install dlv
ENV GOPATH=/home/${CHAINLINK_USER}/go
ENV PATH=$PATH:$GOPATH/bin
RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

# Delve port
EXPOSE 40000

CMD ["dlv", "exec", "/usr/local/bin/chainlink", "--accept-multiclient", "--headless", "--listen=0.0.0.0:40000", "--api-version=2", "--", "local", "node"]

# Build image: Chainlink binary
FROM golang:1.21-bullseye as buildgo
RUN go version
WORKDIR /chainlink

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

COPY . .

# Build the golang binary
RUN make install-chainlink

# Link LOOP Plugin source dirs with simple names
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-feeds | xargs -I % ln -s % /chainlink-feeds
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-solana | xargs -I % ln -s % /chainlink-solana

# Build image: Plugins
FROM golang:1.21-bullseye as buildplugins
RUN go version

WORKDIR /chainlink-feeds
COPY --from=buildgo /chainlink-feeds .
RUN go install ./cmd/chainlink-feeds

WORKDIR /chainlink-solana
COPY --from=buildgo /chainlink-solana .
RUN go install ./pkg/solana/cmd/chainlink-solana

# Final image: ubuntu with chainlink binary
FROM ubuntu:20.04

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-15 \
  && apt-get clean all

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
# explicit set the cache dir. needed so both root and non-root user has an explicit location
ENV XDG_CACHE_HOME /home/${CHAINLINK_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]

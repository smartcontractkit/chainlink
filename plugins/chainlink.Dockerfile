# Build image: Chainlink binary
FROM golang:1.21-bullseye as buildgo
RUN go version
WORKDIR /chainlink

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

COPY . .

RUN apt-get update && apt-get install -y jq

# Build the golang binaries
RUN make install-chainlink

# Install medianpoc binary
RUN make install-medianpoc

# Install ocr3-capability binary
RUN make install-ocr3-capability

# Link LOOP Plugin source dirs with simple names
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-feeds | xargs -I % ln -s % /chainlink-feeds
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-data-streams | xargs -I % ln -s % /chainlink-data-streams
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-solana | xargs -I % ln -s % /chainlink-solana
RUN mkdir /chainlink-starknet
RUN go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-starknet/relayer | xargs -I % ln -s % /chainlink-starknet/relayer

# Build image: Plugins
FROM golang:1.21-bullseye as buildplugins
RUN go version

WORKDIR /chainlink-feeds
COPY --from=buildgo /chainlink-feeds .
RUN go install ./cmd/chainlink-feeds

WORKDIR /chainlink-data-streams
COPY --from=buildgo /chainlink-data-streams .
RUN go install ./mercury/cmd/chainlink-mercury

WORKDIR /chainlink-solana
COPY --from=buildgo /chainlink-solana .
RUN go install ./pkg/solana/cmd/chainlink-solana

WORKDIR /chainlink-starknet/relayer
COPY --from=buildgo /chainlink-starknet/relayer .
RUN go install ./pkg/chainlink/cmd/chainlink-starknet

# Final image: ubuntu with chainlink binary
FROM ubuntu:20.04

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && apt-get clean all

COPY --from=buildgo /go/bin/chainlink /usr/local/bin/
COPY --from=buildgo /go/bin/chainlink-medianpoc /usr/local/bin/
COPY --from=buildgo /go/bin/chainlink-ocr3-capability /usr/local/bin/

COPY --from=buildplugins /go/bin/chainlink-feeds /usr/local/bin/
ENV CL_MEDIAN_CMD chainlink-feeds
COPY --from=buildplugins /go/bin/chainlink-mercury /usr/local/bin/
ENV CL_MERCURY_CMD chainlink-mercury
COPY --from=buildplugins /go/bin/chainlink-solana /usr/local/bin/
ENV CL_SOLANA_CMD chainlink-solana
COPY --from=buildplugins /go/bin/chainlink-starknet /usr/local/bin/
ENV CL_STARKNET_CMD chainlink-starknet

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

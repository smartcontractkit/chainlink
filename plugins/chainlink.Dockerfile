# Build image: Chainlink binary
FROM golang:1.24-bullseye AS buildgo
RUN go version
WORKDIR /chainlink

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA
ARG APTOS_RELAYER_GIT_REF
ARG COSMOS_SHA
ARG STARKNET_SHA
# Flag to control installation of private plugins (default: false)
ARG CL_INSTALL_PRIVATE_PLUGINS=false

# Flags for Go Delve debugger
ARG GO_GCFLAGS

COPY . .

# Used to authenticate with GitHub to fetch private dependencies.
RUN --mount=type=secret,id=GIT_AUTH_TOKEN ./plugins/scripts/setup_git_auth.sh

RUN apt-get update && apt-get install -y jq

# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Build the golang binaries
RUN make GO_GCFLAGS="${GO_GCFLAGS}" install-chainlink

# Install medianpoc binary
RUN make install-medianpoc

# Install ocr3-capability binary
RUN make install-ocr3-capability

# Install LOOP Plugins
RUN make install-plugins \
  APTOS_RELAYER_GIT_REF=${APTOS_RELAYER_GIT_REF} \
  COSMOS_SHA=${COSMOS_SHA} \
  STARKNET_SHA=${STARKNET_SHA} \
  CL_INSTALL_PRIVATE_PLUGINS=${CL_INSTALL_PRIVATE_PLUGINS}

# -----------------------------------------------------------------------------
# Final image: common base stage for the final images
FROM ubuntu:24.04 AS final-base

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && apt-get clean all

# Copy Delve debugger from build stage
COPY --from=buildgo /go/bin/dlv /usr/local/bin/dlv

COPY --from=buildgo /go/bin/chainlink /usr/local/bin/
COPY --from=buildgo /go/bin/chainlink-medianpoc /usr/local/bin/
COPY --from=buildgo /go/bin/chainlink-ocr3-capability /usr/local/bin/
COPY --from=buildgo /go/bin/chainlink-feeds /usr/local/bin/
ENV CL_MEDIAN_CMD=chainlink-feeds
COPY --from=buildgo /go/bin/chainlink-mercury /usr/local/bin/
ENV CL_MERCURY_CMD=chainlink-mercury
COPY --from=buildgo /go/bin/chainlink-cosmos /usr/local/bin/
COPY --from=buildgo /go/bin/chainlink-solana /usr/local/bin/
ENV CL_SOLANA_CMD=chainlink-solana
COPY --from=buildgo /go/bin/chainlink-starknet /usr/local/bin/

# Dependency of CosmWasm/wasmd
COPY --from=buildgo /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.*.so /usr/lib/
RUN chmod 755 /usr/lib/libwasmvm.*.so

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

EXPOSE 6688

# -----------------------------------------------------------------------------
# Final image with private plugins (placed earlier so it's not the default target/stage)
FROM final-base AS final-private-plugins
COPY --from=buildgo /go/bin/chainlink-aptos /usr/local/bin/
ENV CL_APTOS_CMD=chainlink-aptos
ENTRYPOINT ["chainlink"]
HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1
CMD ["local", "node"]

# -----------------------------------------------------------------------------
# Final image without private plugins (this is the last stage, so it will be built by default)
FROM final-base AS final
ENTRYPOINT ["chainlink"]
HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1
CMD ["local", "node"]

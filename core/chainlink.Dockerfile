##
# Build image: Chainlink binary with plugins.
##
FROM golang:1.24-bullseye AS buildgo
RUN go version
RUN apt-get update && apt-get install -y jq && rm -rf /var/lib/apt/lists/*

WORKDIR /chainlink

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .

# Flag to control installation of private plugins (default: true).
ARG CL_INSTALL_PRIVATE_PLUGINS=true
# Flag to control installation of testing plugins (default: false).
ARG CL_INSTALL_TESTING_PLUGINS=false
# Env vars needed for chainlink build
ARG COMMIT_SHA
# Build chainlink bin with cover flag https://go.dev/doc/build-cover#FAQ
ARG GO_COVER_FLAG=false

ENV CL_LOOPINSTALL_OUTPUT_DIR=/tmp/loopinstall-output
RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    ./plugins/scripts/setup_git_auth.sh && \
    mkdir -p /gobins && mkdir -p "${CL_LOOPINSTALL_OUTPUT_DIR}" && \
    GOBIN=/go/bin make install-loopinstall && \
    GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-local install-plugins-public && \
    if [ "${CL_INSTALL_PRIVATE_PLUGINS}" = "true" ]; then \
        GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-private; \
    fi && \
    if [ "${CL_INSTALL_TESTING_PLUGINS}" = "true" ]; then \
        GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-testing; \
    fi

# Copy any shared libraries.
RUN --mount=type=cache,target=/go/pkg/mod \
    mkdir -p /tmp/lib && \
    ./plugins/scripts/copy_loopinstall_libs.sh \
    "$CL_LOOPINSTALL_OUTPUT_DIR" \
    /tmp/lib

# Build chainlink.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    if [ "$GO_COVER_FLAG" = "true" ]; then \
          GOBIN=/gobins make install-chainlink-cover; \
      else \
          GOBIN=/gobins make install-chainlink; \
      fi

##
# Final Image
##
FROM ubuntu:24.04

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl && rm -rf /var/lib/apt/lists/*

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && rm -rf /var/lib/apt/lists/*

RUN if [ ${CHAINLINK_USER} != root ]; then useradd --uid 14933 --create-home ${CHAINLINK_USER}; fi
USER ${CHAINLINK_USER}

# Set plugin environment variable configuration.
ENV CL_SOLANA_CMD=chainlink-solana

ARG CL_MEDIAN_CMD
ENV CL_MEDIAN_CMD=${CL_MEDIAN_CMD}
ARG CL_EVM_CMD
ENV CL_EVM_CMD=${CL_EVM_CMD}

# CCIP specific
COPY ./cci[p]/confi[g] /ccip-config
ARG CL_CHAIN_DEFAULTS
ENV CL_CHAIN_DEFAULTS=${CL_CHAIN_DEFAULTS}

# Copy the binaries from the build stage (plugins + chainlink).
COPY --from=buildgo /gobins/ /usr/local/bin/
# Copy shared libraries from the build stage.
COPY --from=buildgo /tmp/lib /usr/lib/

WORKDIR /home/${CHAINLINK_USER}

# Explicitly set the cache dir. Needed so both root and non-root user has an explicit location.
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

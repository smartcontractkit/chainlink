##
# Build image: Chainlink binary with plugins.
##

# Stage: deps-base — module downloads, no source tree.
# Stages that don't need the full source (remote plugins, delve) branch from
# here so that source-only changes never invalidate their layer cache.
FROM golang:1.25.7-bookworm AS deps-base
RUN go version
RUN apt-get update && apt-get install -y jq && rm -rf /var/lib/apt/lists/*

WORKDIR /chainlink

ADD go.mod go.sum ./
RUN go mod download

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

# Stage: deps — full source tree for stages that compile chainlink code.
FROM deps-base AS deps
COPY . .

# Stage: Delve debugger (no source needed, branches from deps-base)
FROM deps-base AS build-delve
RUN go install github.com/go-delve/delve/cmd/dlv@v1.24.2

# Stage: Remote plugins — only manifest YAMLs, no source tree.
# Cached as long as go.mod/go.sum and plugin manifests are unchanged,
# so typical source-only PRs skip the entire ~160s remote plugin build.
# Uses `go tool loopinstall` via the Makefile (resolved from the `tool`
# directive in go.mod). If this fails without the full source tree, fall back
# to installing loopinstall standalone:
#   RUN go install github.com/smartcontractkit/chainlink-common/pkg/loop/cmd/loopinstall@v0.11.1
# and invoke `loopinstall` directly instead of `make install-plugins-*`.
FROM deps-base AS build-remote-plugins
ARG CL_INSTALL_PRIVATE_PLUGINS=true
ARG CL_INSTALL_TESTING_PLUGINS=false

COPY plugins/plugins.public.yaml plugins/plugins.private.yaml plugins/plugins.testing.yaml ./plugins/
COPY plugins/scripts/ ./plugins/scripts/

ENV CL_LOOPINSTALL_OUTPUT_DIR=/tmp/loopinstall-output \
    GIT_CONFIG_GLOBAL=/tmp/gitconfig-github-token
RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
    set -e && \
    trap 'rm -f "$GIT_CONFIG_GLOBAL"' EXIT && \
    ./plugins/scripts/setup_git_auth.sh && \
    mkdir -p /gobins "${CL_LOOPINSTALL_OUTPUT_DIR}" && \
    GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-public && \
    if [ "${CL_INSTALL_PRIVATE_PLUGINS}" = "true" ]; then \
        GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-private; \
    fi && \
    if [ "${CL_INSTALL_TESTING_PLUGINS}" = "true" ]; then \
        GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-testing; \
    fi

RUN mkdir -p /tmp/lib && \
    ./plugins/scripts/copy_loopinstall_libs.sh \
    "$CL_LOOPINSTALL_OUTPUT_DIR" \
    /tmp/lib

# Stage: Local plugins (needs source tree for ./plugins/cmd/...)
FROM deps AS build-local-plugins
RUN --mount=type=cache,target=/root/.cache/go-build,id=go-build-local-plugins \
    mkdir -p /gobins && \
    GOBIN=/gobins make install-plugins-local

# Stage: Chainlink binary (needs source tree)
FROM deps AS build-chainlink
ARG COMMIT_SHA
ARG VERSION_TAG
ARG CL_IS_PROD_BUILD=true

RUN --mount=type=cache,target=/root/.cache/go-build,id=go-build-chainlink \
    mkdir -p /gobins && \
    if [ "$CL_IS_PROD_BUILD" = "false" ]; then \
          GOBIN=/gobins make install-chainlink-dev; \
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

# Expose image metadata to the running node.
ARG CL_AUTO_DOCKER_TAG=unset
ENV CL_DOCKER_TAG=${CL_AUTO_DOCKER_TAG}

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

# Copy binaries from the parallel build stages.
COPY --from=build-remote-plugins /gobins/ /usr/local/bin/
COPY --from=build-local-plugins /gobins/ /usr/local/bin/
COPY --from=build-chainlink /gobins/ /usr/local/bin/
# Copy shared libraries from the remote plugins build stage.
COPY --from=build-remote-plugins /tmp/lib /usr/lib/
COPY --from=build-delve /go/bin/dlv /usr/local/bin/


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

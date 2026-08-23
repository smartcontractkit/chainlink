##
# Build image: Chainlink binary with plugins.
##

# Stage: deps-base — module downloads, no source tree.
# Stages that don't need the full source (remote plugins, delve) branch from
# here so that source-only changes never invalidate their layer cache.
FROM golang:1.26.6-bookworm AS deps-base
RUN go version
RUN apt-get update && apt-get install -y --no-install-recommends jq=1.6-2.1+deb12u2 && rm -rf /var/lib/apt/lists/*

WORKDIR /chainlink

COPY go.mod go.sum ./
COPY plugins/scripts/setup_git_auth.sh ./plugins/scripts/

# CL_GOPRIVATE: set to "github.com/smartcontractkit/*" when building images
# that depend on private Go modules (e.g. chainlink-internal-solana). When
# empty (the default), go mod download uses the public module proxy as usual.
ARG CL_GOPRIVATE=""
ENV GOPRIVATE="${CL_GOPRIVATE}" \
    GOCACHE=/root/.cache/go-build \
    GOMODCACHE=/go/pkg/mod
RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
    --mount=type=cache,id=go-mod-cache,target=/go/pkg/mod \
    --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build \
    set -e \
    && export GIT_CONFIG_GLOBAL=/tmp/gitconfig-go-mod-download \
    && trap 'rm -f "$GIT_CONFIG_GLOBAL"' EXIT \
    && ./plugins/scripts/setup_git_auth.sh \
    && go mod download

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/


# Stage: deps — Go source tree for stages that compile chainlink code.
# `COPY . .` relies on .dockerignore excluding everything not needed for the
# build (docs, CI config, test infra, non-Go tooling) so unrelated workspace
# churn doesn't invalidate the layer.
FROM deps-base AS deps
COPY . .

# Stage: Delve debugger (no source needed, branches from deps-base)
FROM deps-base AS build-delve
RUN --mount=type=cache,id=go-mod-cache,target=/go/pkg/mod \
    --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build \
    go install github.com/go-delve/delve/cmd/dlv@v1.24.2

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
    --mount=type=cache,id=go-mod-cache,target=/go/pkg/mod \
    --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build \
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
    fi && \
    mkdir -p /tmp/lib && \
    ./plugins/scripts/copy_loopinstall_libs.sh \
    "$CL_LOOPINSTALL_OUTPUT_DIR" \
    /tmp/lib

# Stage: Local plugins (needs source tree for ./plugins/cmd/...)
FROM deps AS build-local-plugins
RUN --mount=type=cache,id=go-mod-cache,target=/go/pkg/mod \
    --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build \
    mkdir -p /gobins \
    && GOBIN=/gobins make install-plugins-local

# Stage: Chainlink binary (needs source tree)
FROM deps AS build-chainlink
ARG COMMIT_SHA
ARG VERSION_TAG
ARG CL_IS_PROD_BUILD=true

RUN --mount=type=cache,id=go-mod-cache,target=/go/pkg/mod \
    --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build \
    mkdir -p /gobins && \
    if [ "$CL_IS_PROD_BUILD" = "false" ]; then \
          GOBIN=/gobins make install-chainlink-dev; \
    else \
          GOBIN=/gobins make install-chainlink; \
    fi

##
# Final Image
##
FROM ubuntu:24.04 AS final
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates=20260601~24.04.1 \
    gnupg=2.4.4-2ubuntu17.4 \
    lsb-release=12.0-2 \
    curl=8.5.0-2ubuntu10.12 \
    && rm -rf /var/lib/apt/lists/*

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc \
    | gpg --dearmor -o /usr/share/keyrings/postgresql-archive-keyring.gpg \
    && gpg --no-default-keyring \
        --keyring /usr/share/keyrings/postgresql-archive-keyring.gpg \
        --fingerprint B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8 \
    && echo "deb [signed-by=/usr/share/keyrings/postgresql-archive-keyring.gpg] \
    https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" \
        >/etc/apt/sources.list.d/pgdg.list \
    && apt-get update && apt-get install -y --no-install-recommends postgresql-client-17=17.10-1.pgdg24.04+1 \
    && rm -rf /var/lib/apt/lists/*

# Prod images (CHAINLINK_USER=chainlink) run as UID:GID 14933:14933 for deterministic
# ownership on bind mounts in Docker/Compose deployments without K8s securityContext overrides.
RUN if [ ${CHAINLINK_USER} != root ]; then \
      groupadd --gid 14933 ${CHAINLINK_USER} && \
      useradd --uid 14933 --gid 14933 --create-home ${CHAINLINK_USER}; \
    fi
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

FROM final AS debug

COPY --from=build-delve /go/bin/dlv /usr/local/bin/

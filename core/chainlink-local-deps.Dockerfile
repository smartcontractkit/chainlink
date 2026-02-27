##
# Build image: Chainlink binary with plugins, using local chainlink-common (go.mod replace).
# Use when chainlink/go.mod has: replace github.com/smartcontractkit/chainlink-common => ../chainlink-common
#
# Build context MUST be the parent directory containing both chainlink/ and chainlink-common/
# (e.g. docker_ctx = ".." when running from chainlink repo). When CL_USE_LOCAL_CAPABILITIES=true,
# context must also contain capabilities/ (sibling of chainlink/) so private plugins are built from
# local source instead of GitHub (no GIT_AUTH_TOKEN needed).
##
FROM golang:1.25.7-bookworm AS buildgo
RUN go version
RUN apt-get update && apt-get install -y jq && rm -rf /var/lib/apt/lists/*

WORKDIR /chainlink

# Satisfy go.mod replace ../chainlink-common before go mod download
COPY chainlink-common /chainlink-common

# When CL_USE_LOCAL_CAPABILITIES=true, context must include capabilities/ (script ensures it exists).
COPY capabilities /capabilities

COPY chainlink/GNUmakefile chainlink/package.json ./
COPY chainlink/tools/bin/ldflags ./tools/bin/

# ARG early so we can apply replace before first go mod download (avoids fetching capabilities from network).
ARG CL_USE_LOCAL_CAPABILITIES=false
ARG CL_USE_LOCAL_APTOS=false

ADD chainlink/go.mod chainlink/go.sum ./
# With local capabilities, add replace before any download so the module cache never hits the network for capabilities.
RUN if [ "${CL_USE_LOCAL_CAPABILITIES}" = "true" ] && [ -f /capabilities/go.mod ]; then \
    go mod edit -replace=github.com/smartcontractkit/capabilities=../capabilities; \
fi
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY chainlink/ .

# Ensure go.mod is consistent for plugin builds (avoids "go: updates to go.mod needed; go mod tidy" during install-plugins-local).
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod tidy

# When using local capabilities: re-apply replace (COPY chainlink/ overwrote go.mod) and create go.work so capability
# submodules (e.g. libs) resolve when building plugins from local paths like /capabilities/cron.
# Optional: CL_USE_LOCAL_APTOS=true and context containing chainlink-aptos/ use local repo for
# capabilities/chain_capabilities/aptos (write-target work).
RUN if [ "${CL_USE_LOCAL_CAPABILITIES}" = "true" ] && [ -f /capabilities/go.mod ]; then \
    go mod edit -replace=github.com/smartcontractkit/capabilities=../capabilities; \
    printf '%s\n' \
      'go 1.25.5' '' 'use (' \
      './chainlink' \
      './capabilities/libs' \
      './capabilities/cron' \
      './capabilities/readcontract' \
      './capabilities/consensus' \
      './capabilities/http_action' \
      './capabilities/http_trigger' \
      './capabilities/chain_capabilities/evm' \
      './capabilities/chain_capabilities/solana' \
      './capabilities/chain_capabilities/aptos' \
      './capabilities/mock' \
      './capabilities/kvstore' \
      './capabilities/workflowevent' \
      ')' '' \
      'replace github.com/fbsobreira/gotron-sdk => github.com/smartcontractkit/chainlink-tron/relayer/gotron-sdk v0.0.5-0.20251014124537-af6b1684fe15' \
      > /go.work; \
fi
# When CL_USE_LOCAL_APTOS=true, copy chainlink-aptos from build context (excluding build artifacts
# to keep step fast) and add to go.work so capabilities/chain_capabilities/aptos builds against local
# chainlink-aptos (no conditional COPY).
RUN --mount=type=bind,source=.,target=/ctx \
  if [ "${CL_USE_LOCAL_APTOS}" = "true" ] && [ -d /ctx/chainlink-aptos ] && [ -f /ctx/chainlink-aptos/go.mod ]; then \
    apt-get update -qq && apt-get install -y -qq rsync && rm -rf /var/lib/apt/lists/* && \
    mkdir -p /chainlink-aptos && \
    rsync -a \
      --exclude='build/' --exclude='target/' \
      --exclude='*.hex' --exclude='*.zip' \
      /ctx/chainlink-aptos/ /chainlink-aptos/ && \
    sed -i '/^)$/i\  ./chainlink-aptos' /go.work; \
  fi

# Install Delve for debugging with cache mounts
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/go-delve/delve/cmd/dlv@v1.24.2

# Flag to control installation of private plugins (default: true).
ARG CL_INSTALL_PRIVATE_PLUGINS=true
# Flag to control installation of testing plugins (default: false).
ARG CL_INSTALL_TESTING_PLUGINS=false
# Env vars needed for chainlink build
ARG COMMIT_SHA
ARG VERSION_TAG
# Flag to control whether this is a prod build (default: true).
ARG CL_IS_PROD_BUILD=true

ENV CL_LOOPINSTALL_OUTPUT_DIR=/tmp/loopinstall-output \
    GIT_CONFIG_GLOBAL=/tmp/gitconfig-github-token
# Secret must be provided by the build (use a dummy empty file when CL_USE_LOCAL_CAPABILITIES=true and no token).
# When CL_USE_LOCAL_CAPABILITIES=true, set GOWORK only for install-plugins-local (and -private) so capabilities
# resolve from the workspace; clear cached capabilities. install-plugins-public must run without GOWORK so
# plugins built from GOMODCACHE are not required to be in go.work.
RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    set -e && \
    trap 'rm -f "$GIT_CONFIG_GLOBAL"' EXIT && \
    ./plugins/scripts/setup_git_auth.sh && \
    mkdir -p /gobins && mkdir -p "${CL_LOOPINSTALL_OUTPUT_DIR}" && \
    if [ "${CL_USE_LOCAL_CAPABILITIES}" = "true" ] && [ -f /capabilities/go.mod ]; then \
      rm -rf /go/pkg/mod/github.com/smartcontractkit/capabilities* /go/pkg/mod/cache/vcs/* 2>/dev/null || true; \
      export GOWORK=/go.work; \
    fi && \
    GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-local && \
    GOWORK=off GOBIN=/gobins CL_LOOPINSTALL_OUTPUT_DIR=${CL_LOOPINSTALL_OUTPUT_DIR} make install-plugins-public && \
    if [ "${CL_INSTALL_PRIVATE_PLUGINS}" = "true" ]; then \
        if [ "${CL_USE_LOCAL_CAPABILITIES}" = "true" ] && [ -f /capabilities/go.mod ]; then \
            cp plugins/plugins.private.local.yaml plugins/plugins.private.yaml; \
            export GOWORK=/go.work; \
        fi; \
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

# Set plugin environment variable configuration.
ENV CL_SOLANA_CMD=chainlink-solana

ARG CL_MEDIAN_CMD
ENV CL_MEDIAN_CMD=${CL_MEDIAN_CMD}
ARG CL_EVM_CMD
ENV CL_EVM_CMD=${CL_EVM_CMD}

# CCIP specific (path relative to context = parent; chainlink/ccip/config)
COPY chainlink/ccip/config /ccip-config
ARG CL_CHAIN_DEFAULTS
ENV CL_CHAIN_DEFAULTS=${CL_CHAIN_DEFAULTS}

# Copy the binaries from the build stage (plugins + chainlink).
COPY --from=buildgo /gobins/ /usr/local/bin/
# Copy shared libraries from the build stage.
COPY --from=buildgo /tmp/lib /usr/lib/
# Copy dlv (Delve debugger) from the build stage.
COPY --from=buildgo /go/bin/dlv /usr/local/bin/


WORKDIR /home/${CHAINLINK_USER}

# So capability_defaults.toml binary_path "./binaries/<name>" resolve when using this image without host copy.
RUN mkdir -p binaries && \
    ln -sf /usr/local/bin/chainlink-aptos binaries/aptos 2>/dev/null || true && \
    ln -sf /usr/local/bin/chainlink-cron binaries/cron 2>/dev/null || true

# Explicitly set the cache dir. Needed so both root and non-root user has an explicit location.
ENV XDG_CACHE_HOME=/home/${CHAINLINK_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

# Set up env and dir for go coverage profiling https://go.dev/doc/build-cover#FAQ
ENV GOCOVERDIR=/var/tmp/go-coverage
RUN mkdir -p /var/tmp/go-coverage

EXPOSE 6688
ENTRYPOINT ["chainlink"]
HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1
CMD ["local", "node"]

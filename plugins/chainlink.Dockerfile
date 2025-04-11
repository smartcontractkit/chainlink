# Build image: Chainlink binary with plugins.
FROM golang:1.24-bullseye AS buildgo
RUN go version
RUN apt-get update && apt-get install -y jq && rm -rf /var/lib/apt/lists/*

WORKDIR /chainlink

COPY GNUmakefile package.json ./
COPY tools/bin/ldflags ./tools/bin/

ADD go.mod go.sum ./
RUN go mod download
COPY . .

# Install Delve for debugging.
RUN go install github.com/go-delve/delve/cmd/dlv@v1.24.2

# Flag to control installation of private plugins (default: false).
ARG CL_INSTALL_PRIVATE_PLUGINS=false
# Flags for Go Delve debugger
ARG GO_GCFLAGS
# Env vars needed for chainlink build
ARG COMMIT_SHA
ARG COSMOS_SHA
ARG STARKNET_SHA

# Flags for Go Delve debugger.
ARG GO_GCFLAGS

# Install plugins to a specific directory to make it easier to copy to final image.
RUN GOBIN=/go/bin make install-loopinstall

RUN apt-get update && apt-get install -y jq && rm -rf /var/lib/apt/lists/*

RUN --mount=type=secret,id=GIT_AUTH_TOKEN ./plugins/scripts/setup_git_auth.sh && \
    mkdir -p /gobins && \
    GOBIN=/gobins make install-plugins-local install-plugins-public && \
    if [ "${CL_INSTALL_PRIVATE_PLUGINS}" = "true" ]; then \
        GOBIN=/gobins make install-plugins-private; \
    fi

# Build the golang binaries.
RUN GOBIN=/gobins make GO_GCFLAGS="${GO_GCFLAGS}" install-chainlink

# TODO: name build-manifest to account for different plugin paths to support multiple plugin files.
# Copy any additional files specified in the build manifest.
RUN if [ -f "./plugins/docker/output_manifests/build-manifest.json" ]; then \
        echo "Processing build manifest for additional files..." && \
        jq -r '.plugins | to_entries[] | select(.value.additionalFiles != null) | .value.additionalFiles[] | "\(.src):\(.dest)"' \
        ./plugins/docker/output_manifests/build-manifest.json > /tmp/additional_files.txt && \
        if [ -s "/tmp/additional_files.txt" ]; then \
            cat /tmp/additional_files.txt && \
            ./plugins/scripts/copy_additional_files.sh /tmp/additional_files.txt || true; \
        else \
            echo "No additional files to copy"; \
        fi \
    fi

# Final image: ubuntu with chainlink binary
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

# Copy Delve debugger from build stage.
COPY --from=buildgo /go/bin/dlv /usr/local/bin/dlv

# Set plugin environment variable configuration.
ENV CL_MEDIAN_CMD=chainlink-feeds
ENV CL_MERCURY_CMD=chainlink-mercury
ENV CL_SOLANA_CMD=chainlink-solana
ARG CL_APTOS_CMD
ENV CL_APTOS_CMD=${CL_APTOS_CMD}

# Copy the binaries from the build stage (plugins + chainlink).
COPY --from=buildgo --chown=${CHAINLINK_USER}:${CHAINLINK_USER} /gobins /usr/local/bin/
# Copy the additional libs based on the manifests from the build stage.
COPY --from=buildgo --chown=${CHAINLINK_USER}:${CHAINLINK_USER} /usr/lib/ /usr/lib/



USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}
# explicit set the cache dir. needed so both root and non-root user has an explicit location
ENV XDG_CACHE_HOME=/home/${CHAINLINK_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]

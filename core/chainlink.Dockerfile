# MAKE ALL CHANGES WITHIN THE DEFAULT WORKDIR FOR YARN AND GO DEP CACHE HITS
FROM node:16-buster
WORKDIR /chainlink

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags tools/bin/ldflags
ARG COMMIT_SHA

# Install yarn dependencies
COPY yarn.lock package.json .yarnrc ./
COPY .yarn .yarn
COPY operator_ui/package.json ./operator_ui/
COPY contracts/package.json ./contracts/
RUN make yarndep

COPY contracts ./contracts
COPY tsconfig.cjs.json ./
COPY core/web/schema core/web/schema
COPY operator_ui ./operator_ui

# Create the directory that the operator-ui build assets will be placed in.
RUN mkdir -p core/web

# Build operator-ui and the smart contracts
RUN make contracts-operator-ui-build

# Build the golang binary

FROM golang:1.18-buster
WORKDIR /chainlink

COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Env vars needed for chainlink build
ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink build
ARG COMMIT_SHA

COPY core core
# Copy over operator-ui build assets to the web module so that we embed them correctly
COPY --from=0 /chainlink/core/web/assets ./core/web/assets

RUN make chainlink-build

# Final layer: ubuntu with chainlink binary
FROM ubuntu:20.04

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
  && curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-14 \
  && apt-get clean all

COPY --from=1 /go/bin/chainlink /usr/local/bin/

# dependency of terra-money/core
COPY --from=1 /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/api/libwasmvm.so /usr/lib/libwasmvm.so
RUN chmod 755 /usr/lib/libwasmvm.so

RUN if [ ${CHAINLINK_USER} != root ]; then \
  useradd --uid 14933 --create-home ${CHAINLINK_USER}; \
  fi
USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]

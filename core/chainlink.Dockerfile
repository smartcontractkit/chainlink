# Build image: Chainlink binary
FROM golang:1.20-buster as buildgo
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

# Final image: ubuntu with chainlink binary
FROM ubuntu:20.04

ARG CHAINLINK_USER=root
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-14 \
  && apt-get clean all

COPY --from=buildgo /go/bin/chainlink /usr/local/bin/

# Dependency of CosmWasm/wasmd
COPY --from=buildgo /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.*.so /usr/lib/
RUN chmod 755 /usr/lib/libwasmvm.*.so

RUN if [ ${CHAINLINK_USER} != root ]; then \
  useradd --uid 14933 --create-home ${CHAINLINK_USER}; \
  fi
USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]

# This will replace chainlink.Dockerfile once all builds are migrated to goreleaser

# Final image: ubuntu with chainlink binary
FROM ubuntu:24.04

ARG CHAINLINK_USER=root
ARG TARGETARCH
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-16 \
  && apt-get clean all \
  && rm -rf /var/lib/apt/lists/*

COPY ./chainlink /usr/local/bin/

# Copy native libs if cgo is enabled
COPY ./tmp/libs /usr/local/bin/libs

# Copy plugins if exist and enable them
# https://stackoverflow.com/questions/70096208/dockerfile-copy-folder-if-it-exists-conditional-copy/70096420#70096420
COPY ./tm[p]/plugin[s]/ /usr/local/bin/

# Allow individual plugins to be enabled by supplying their path 
ARG CL_MEDIAN_CMD
ARG CL_MERCURY_CMD
ARG CL_SOLANA_CMD
ARG CL_STARKNET_CMD
ENV CL_MEDIAN_CMD=${CL_MEDIAN_CMD} \
  CL_MERCURY_CMD=${CL_MERCURY_CMD} \
  CL_SOLANA_CMD=${CL_SOLANA_CMD} \
  CL_STARKNET_CMD=${CL_STARKNET_CMD}

# CCIP specific
COPY ./cci[p]/confi[g] /chainlink/ccip-config
ARG CL_CHAIN_DEFAULTS
ENV CL_CHAIN_DEFAULTS=${CL_CHAIN_DEFAULTS}

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

# This will replace chainlink.Dockerfile once all builds are migrated to goreleaser

# Final image: ubuntu with chainlink binary
FROM ubuntu:20.04

ARG CHAINLINK_USER=root
ARG TARGETARCH
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates gnupg lsb-release curl

# Install Postgres for CLI tools, needed specifically for DB backups
RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |tee /etc/apt/sources.list.d/pgdg.list \
  && apt-get update && apt-get install -y postgresql-client-14 \
  && apt-get clean all

COPY . /usr/local/bin/
# Copy native libs if cgo is enabled
# COPY ./tmp/linux_${TARGETARCH}/libs /usr/local/bin/libs

RUN if [ ${CHAINLINK_USER} != root ]; then \
  useradd --uid 14933 --create-home ${CHAINLINK_USER}; \
  fi
USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}

EXPOSE 6688
ENTRYPOINT ["chainlink"]

HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1

CMD ["local", "node"]

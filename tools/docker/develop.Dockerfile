FROM ubuntu:20.04

# Add the PostgreSQL PGP key & repository
RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

# Install deps
RUN apt-get update && apt-get install -y postgresql postgresql-contrib direnv build-essential cmake libudev-dev unzip

# Install additional tooling
RUN mkdir -p ~/.local/bin/
ENV PATH="/root/.local/bin:${PATH}"
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/google/gofuzz
RUN pnpm install -g ganache-cli
RUN pip3 install web3 slither-analyzer crytic-compile
RUN curl -L https://github.com/crytic/echidna/releases/download/v1.5.1/echidna-test-v1.5.1-Ubuntu-18.04.tar.gz | tar -xz -C ~/.local/bin
RUN curl -L https://github.com/openethereum/openethereum/releases/download/v3.2.4/openethereum-linux-v3.2.4.zip --output openethereum.zip
RUN unzip openethereum.zip -d ~/.local/bin/ && rm openethereum.zip
RUN chmod +x ~/.local/bin/*

# Setup direnv
RUN echo 'eval "$(direnv hook bash)"' > /root/.bashrc

# Setup postgres
USER postgres
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/10/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/10/main/postgresql.conf
RUN /etc/init.d/postgresql start &&\
  createdb chainlink_test &&\
  createdb node_dev &&\
  createuser --superuser --no-password root &&\
  psql -c "ALTER USER postgres PASSWORD 'node';"

USER root

# add init file - this file starts postgres and keeps container alive after started
RUN touch ~/init
RUN chmod +x ~/init
RUN echo "#!/usr/local/bin/dumb-init /bin/sh" >> ~/init
RUN echo "/etc/init.d/postgresql start" >> ~/init
RUN echo "while true; do sleep 1; done" >> ~/init

ARG SRCROOT=/root/chainlink
WORKDIR ${SRCROOT}

EXPOSE 5432
EXPOSE 6688
EXPOSE 6689
EXPOSE 3000
EXPOSE 3001
EXPOSE 8545
EXPOSE 8546

# Default env setup for testing
ENV CHAINLINK_DB_NAME chainlink_test
ENV CHAINLINK_PGPASSWORD=thispasswordislongenough
ENV CL_DATABASE_URL=postgresql://postgres:$CHAINLINK_PGPASSWORD@localhost:5432/$CHAINLINK_DB_NAME?sslmode=disable
ENV TYPEORM_USERNAME=postgres
ENV TYPEORM_PASSWORD=node
ENV ETH_CHAIN_ID=1337
ENV CHAINLINK_DEV=true
ENV CHAINLINK_TLS_PORT=0
ENV SECURE_COOKIES=false

ENTRYPOINT [ "/root/init" ]

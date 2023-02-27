FROM ubuntu:20.04

# Add the PostgreSQL PGP key & repository
RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN echo "xsBNBGKItdQBCADWmKTNZEYWgXy73FvKFY5fRro4tGNa4Be4TZW3wZpct9Cj8EjykU7S9EPoJ3EdKpxFltHRu7QbDi6LWSNA4XxwnudQrYGxnxx6Ru1KBHFxHhLfWsvFcGMwit/znpxtIt9UzqCm2YTEW5NUnzQ4rXYqVQK2FLG4weYJ5bKwkY+ZsnRJpzxdHGJ0pBiqwkMT8bfQdJymUBown+SeuQ2HEqfjVMsIRe0dweD2PHWeWo9fTXsz1Q5abiGckyOVyoN9//DgSvLUocUcZsrWvYPaN+o8lXTO3GYFGNVsx069rxarkeCjOpiQOWrQmywXISQudcusSgmmgfsRZYW7FDBy5MQrABEBAAHNUVJhcHR1cmUgQXV0b21hdGljIFNpZ25pbmcgS2V5IChjbG91ZC1yYXB0dXJlLXNpZ25pbmcta2V5LTIwMjItMDMtMDctMDhfMDFfMDEucHViKcLAYgQTAQgAFgUCYoi11AkQtT3IDRPt7wUCGwMCGQEAAMGoCAB8QBNIIN3Q2D3aahrfkb6axd55zOwR0tnriuJRoPHoNuorOpCv9aWMMvQACNWkxsvJxEF8OUbzhSYjAR534RDigjTetjK2i2wKLz/kJjZbuF4ZXMynCm40eVm1XZqU63U9XR2RxmXppyNpMqQO9LrzGEnNJuh23icaZY6no12axymxcle/+SCmda8oDAfa0iyA2iyg/eU05buZv54MC6RB13QtS+8vOrKDGr7RYp/VYvQzYWm+ck6DvlaVX6VB51BkLl23SQknyZIJBVPm8ttU65EyrrgG1jLLHFXDUqJ/RpNKq+PCzWiyt4uy3AfXK89RczLu3uxiD0CQI0T31u/IzsBNBGKItdQBCADIMMJdRcg0Phv7+CrZz3xRE8Fbz8AN+YCLigQeH0B9lijxkjAFr+thB0IrOu7ruwNY+mvdP6dAewUur+pJaIjEe+4s8JBEFb4BxJfBBPuEbGSxbi4OPEJuwT53TMJMEs7+gIxCCmwioTggTBp6JzDsT/cdBeyWCusCQwDWpqoYCoUWJLrUQ6dOlI7s6p+iIUNIamtyBCwb4izs27HdEpX8gvO9rEdtcb7399HyO3oD4gHgcuFiuZTpvWHdn9WYwPGM6npJNG7crtLnctTR0cP9KutSPNzpySeAniHx8L9ebdD9tNPCWC+OtOcGRrcBeEznkYh1C4kzdP1ORm5upnknABEBAAHCwF8EGAEIABMFAmKItdQJELU9yA0T7e8FAhsMAABJmAgAhRPk/dFj71bU/UTXrkEkZZzE9JzUgan/ttyRrV6QbFZABByf4pYjBj+yLKw3280//JWurKox2uzEq1hdXPedRHICRuh1Fjd00otaQ+wGF3kY74zlWivB6Wp6tnL9STQ1oVYBUv7HhSHoJ5shELyedxxHxurUgFAD+pbFXIiK8cnAHfXTJMcrmPpC+YWEC/DeqIyEcNPkzRhtRSuERXcq1n+KJvMUAKMD/tezwvujzBaaSWapmdnGmtRjjL7IxUeGamVWOwLQbUr+34MwzdeJdcL8fav5LA8Uk0ulyeXdwiAK8FKQsixI+xZvz7HUs8ln4pZwGw/TpvO9cMkHogtgzQ==" \
  | base64 -d | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add

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
  createdb explorer_dev &&\
  createdb explorer_test &&\
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
ENV CHAINLINK_PGPASSWORD=node
ENV CL_DATABASE_URL=postgresql://postgres:$CHAINLINK_PGPASSWORD@localhost:5432/$CHAINLINK_DB_NAME?sslmode=disable
ENV TYPEORM_USERNAME=postgres
ENV TYPEORM_PASSWORD=node
ENV ETH_CHAIN_ID=1337
ENV CHAINLINK_DEV=true
ENV CHAINLINK_TLS_PORT=0
ENV SECURE_COOKIES=false

ENTRYPOINT [ "/root/init" ]

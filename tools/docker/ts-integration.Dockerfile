FROM node:12.18

# Install docker and docker compose
RUN apt-get update \
    # Install Docker CE CLI
    && apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common lsb-release \
    && curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | apt-key add - 2>/dev/null \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" \
    && apt-get update \
    && apt-get install -y docker-ce-cli \
    #
    # Install Docker Compose
    && curl -sSL "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod +x /usr/local/bin/docker-compose

ENV PATH=/chainlink/tools/bin:./node_modules/.bin:$PATH

# Copy only what we neeed
ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

# copy over all our dependencies
COPY yarn.lock package.json .yarnrc tsconfig.cjs.json tsconfig.es6.json ./
COPY patches patches
COPY solc_bin solc_bin
COPY tools/bin/restore-solc-cache tools/bin/restore-solc-cache
COPY .yarn .yarn
COPY belt belt
COPY evm-test-helpers evm-test-helpers
COPY evm-contracts evm-contracts
# TODO remove @types import
# https://www.pivotaltracker.com/story/show/171715396
COPY operator_ui/@types operator_ui/@types/
COPY tools/ci-ts tools/ci-ts

# install deps
RUN yarn install
RUN tools/bin/restore-solc-cache

# setup contracts
RUN yarn setup:contracts

FROM smartcontract/builder:1.0.40

RUN wget --quiet -O - https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - 2>/dev/null

# Install docker and docker compose
RUN apt-get update \
    #
    # Install Docker CE CLI
    && apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common lsb-release \
    && curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | apt-key add - 2>/dev/null \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" \
    && apt-get update \
    && apt-get install -y docker-ce-cli \
    #
    # Install Docker Compose
    && curl -sSL "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod +x /usr/local/bin/docker-compose \
    #
    # Install jq
    && apt-get -y install jq \
    # Install dependencies needed to run cypress with chrome
    && apt-get install -y xvfb libgtk-3-dev libnotify-dev libgconf-2-4 libnss3 libxss1 libasound2 fonts-liberation libappindicator3-1 xdg-utils
RUN wget -O ./chrome79.deb https://chainlink-monorepo-public-s3.s3-us-west-2.amazonaws.com/integration-tests/chrome79.deb
RUN dpkg -i chrome79.deb

ENV PATH=/chainlink/tools/bin:./node_modules/.bin:$PATH

# Copy only what we neeed
ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

COPY yarn.lock package.json .yarnrc ./
COPY patches patches
COPY solc_bin solc_bin
COPY tools/bin/restore-solc-cache tools/bin/restore-solc-cache
COPY .yarn .yarn
COPY belt/package.json ./belt/
COPY belt/bin ./belt/bin
COPY evm-test-helpers/package.json evm-test-helpers/
COPY contracts/package.json ./contracts/
COPY integration/package.json integration/
COPY integration-scripts/package.json integration-scripts/

# install deps for our integration scripts
RUN yarn
RUN tools/bin/restore-solc-cache
# copy our CI test
COPY tools/ci/ethereum_test tools/ci/
COPY tools/docker tools/docker/

# copy over all our dependencies
COPY tsconfig.cjs.json tsconfig.es6.json ./
COPY belt belt
COPY evm-test-helpers evm-test-helpers
COPY contracts contracts
COPY integration integration
COPY integration-scripts integration-scripts


# setup our integration testing scripts
RUN yarn setup:integration

ENTRYPOINT [ "tools/ci/ethereum_test" ]

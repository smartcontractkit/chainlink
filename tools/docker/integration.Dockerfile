FROM j16r/integration

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
COPY evm-contracts/package.json ./evm-contracts/
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
COPY evm-contracts evm-contracts
COPY integration integration
COPY integration-scripts integration-scripts


# setup our integration testing scripts
RUN yarn setup:integration

ENTRYPOINT [ "tools/ci/ethereum_test" ]

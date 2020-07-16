FROM smartcontract/builder:1.0.34

WORKDIR /chainlink
COPY go.mod go.sum yarn.lock package.json .yarnrc GNUmakefile ./
COPY tools/bin/ldflags tools/bin/ldflags
COPY tools/bin/restore-solc-cache tools/bin/restore-solc-cache
COPY .git .git
COPY VERSION VERSION
COPY .yarn .yarn 
COPY patches patches
COPY solc_bin solc_bin
COPY belt/package.json belt/package.json
COPY belt/bin ./belt/bin
COPY evm-contracts/package.json evm-contracts/package.json
COPY evm-test-helpers/package.json evm-test-helpers/package.json
COPY explorer/client/package.json explorer/client/package.json
COPY explorer/package.json explorer/package.json
COPY feeds/package.json feeds/package.json
COPY integration/package.json integration/package.json
COPY integration-scripts/package.json integration-scripts/package.json
COPY operator_ui/package.json operator_ui/package.json
COPY styleguide/package.json styleguide/package.json
COPY tools/ci-ts/package.json tools/ci-ts/package.json
COPY tools/cypress-job-server/package.json tools/cypress-job-server/package.json
COPY tools/echo-server/package.json tools/echo-server/package.json
COPY tools/external-adapter/package.json tools/external-adapter/package.json
COPY tools/json-api-client/package.json tools/json-api-client/package.json
COPY tools/local-storage/package.json tools/local-storage/package.json
COPY tools/package.json tools/package.json
COPY tools/redux/package.json tools/redux/package.json
COPY tools/ts-helpers/package.json tools/ts-helpers/package.json

RUN make gen-builder-cache

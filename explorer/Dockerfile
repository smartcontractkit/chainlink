FROM henrynguyen5/base:1.0.3 as builder

ENV PATH /chainlink/node_modules/.bin:$PATH

COPY yarn.lock package.json .yarnrc ./
COPY .yarn .yarn
COPY explorer/package.json ./explorer/
COPY explorer/client/package.json ./explorer/client/
COPY tools/package.json ./tools/
COPY tools/json-api-client/package.json ./tools/json-api-client/
COPY tools/local-storage/package.json ./tools/local-storage/
COPY tools/redux/package.json ./tools/redux/
COPY tools/ts-helpers/package.json ./tools/ts-helpers/
COPY styleguide/package.json ./styleguide/
RUN yarn install --frozen-lockfile

COPY tsconfig.*.json babel.config.js ./
COPY explorer ./explorer
COPY tools ./tools
COPY styleguide ./styleguide
RUN yarn setup:explorer
COPY .git .git
RUN yarn workspace @chainlink/explorer-client build
RUN yarn workspace @chainlink/explorer write-version

FROM node:10.16.3-alpine

WORKDIR /app

COPY --from=builder /chainlink/package.json package.json
COPY --from=builder /chainlink/yarn.lock yarn.lock
COPY --from=builder /chainlink/node_modules node_modules
COPY --from=builder /chainlink/explorer/VERSION.json explorer/VERSION.json
COPY --from=builder /chainlink/explorer/node_modules explorer/node_modules
COPY --from=builder /chainlink/explorer/package.json explorer/package.json
COPY --from=builder /chainlink/explorer/src explorer/src
COPY --from=builder /chainlink/explorer/tsconfig.json explorer/tsconfig.json
COPY --from=builder /chainlink/tsconfig.cjs.json ./
COPY --from=builder /chainlink/explorer/ormconfig.json explorer/ormconfig.json
COPY --from=builder /chainlink/explorer/client/public explorer/client/public
COPY --from=builder /chainlink/explorer/client/build explorer/client/build
COPY --from=builder /chainlink/styleguide styleguide
COPY --from=builder /chainlink/tools/json-api-client tools/json-api-client
COPY --from=builder /chainlink/tools/redux tools/redux
RUN yarn workspace @chainlink/explorer read-version
ENV NODE_ENV production
ENTRYPOINT [ "yarn", "workspace", "@chainlink/explorer", "run", "prod" ]

FROM node:10.16.3 as builder

RUN apt-get update && apt-get install -y libudev-dev libusb-dev libusb-1.0-0

RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

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

COPY --from=builder /usr/src/app/package.json package.json
COPY --from=builder /usr/src/app/yarn.lock yarn.lock
COPY --from=builder /usr/src/app/node_modules node_modules
COPY --from=builder /usr/src/app/explorer/VERSION.json explorer/VERSION.json
COPY --from=builder /usr/src/app/explorer/node_modules explorer/node_modules
COPY --from=builder /usr/src/app/explorer/package.json explorer/package.json
COPY --from=builder /usr/src/app/explorer/src explorer/src
COPY --from=builder /usr/src/app/explorer/tsconfig.json explorer/tsconfig.json
COPY --from=builder /usr/src/app/tsconfig.cjs.json ./
COPY --from=builder /usr/src/app/explorer/ormconfig.json explorer/ormconfig.json
COPY --from=builder /usr/src/app/explorer/client/public explorer/client/public
COPY --from=builder /usr/src/app/explorer/client/build explorer/client/build
COPY --from=builder /usr/src/app/styleguide styleguide
COPY --from=builder /usr/src/app/tools/json-api-client tools/json-api-client
COPY --from=builder /usr/src/app/tools/redux tools/redux
RUN yarn workspace @chainlink/explorer read-version
ENV NODE_ENV production
ENTRYPOINT [ "yarn", "workspace", "@chainlink/explorer", "run", "prod" ]

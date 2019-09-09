FROM node:10.15 as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

ADD . .

RUN yarn install
RUN cd explorer \
      && yarn install \
      && yarn build

FROM node:10.15-alpine

WORKDIR /app

COPY --from=builder /usr/src/app/package.json package.json
COPY --from=builder /usr/src/app/yarn.lock yarn.lock
COPY --from=builder /usr/src/app/node_modules node_modules
COPY --from=builder /usr/src/app/explorer/node_modules explorer/node_modules
COPY --from=builder /usr/src/app/explorer/package.json explorer/package.json
COPY --from=builder /usr/src/app/explorer/src explorer/src
COPY --from=builder /usr/src/app/explorer/tsconfig.json explorer/tsconfig.json
COPY --from=builder /usr/src/app/explorer/ormconfig.json explorer/ormconfig.json
COPY --from=builder /usr/src/app/explorer/client/public explorer/client/public
COPY --from=builder /usr/src/app/explorer/client/build explorer/client/build
COPY --from=builder /usr/src/app/styleguide styleguide

ENV NODE_ENV production
ENTRYPOINT [ "yarn", "workspace", "@chainlink/explorer", "run", "prod" ]

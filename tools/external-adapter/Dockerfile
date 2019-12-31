FROM node:10.15-alpine

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}
COPY yarn.lock package.json ./
COPY tools/external-adapter/package.json tools/external-adapter/

# install deps for our integration scripts
RUN yarn install

# copy over all our dependencies
COPY tools/external-adapter tools/external-adapter

# build project
RUN yarn workspace @chainlink/external-adapter build

ENV PORT 6644
EXPOSE 6644

ENTRYPOINT [ "yarn", "workspace", "@chainlink/external-adapter", "start" ]

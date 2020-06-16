FROM node:12.18.0-alpine

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}
COPY yarn.lock .yarnrc package.json tsconfig.cjs.json ./
COPY .yarn .yarn
COPY tools/cypress-job-server/package.json tools/cypress-job-server/

# install deps for our integration scripts
RUN yarn install --frozen-lockfile 

# copy over all our dependencies
COPY tools/cypress-job-server tools/cypress-job-server

# setup project
RUN yarn workspace @chainlink/cypress-job-server setup

ENV JOB_SERVER_PORT 6692
EXPOSE 6692

ENTRYPOINT ["yarn", "workspace", "@chainlink/cypress-job-server", "start"]

FROM node:10.16.3-alpine

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}
COPY yarn.lock package.json ./
COPY tools/cypress-job-server/package.json tools/cypress-job-server/

# install deps for our integration scripts
RUN yarn

# copy over all our dependencies
COPY tools/cypress-job-server tools/cypress-job-server

# build project
RUN yarn workspace @chainlink/cypress-job-server build

ENV JOB_SERVER_PORT 6692
EXPOSE 6692

ENTRYPOINT ["yarn", "workspace", "@chainlink/cypress-job-server", "start"]
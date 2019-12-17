FROM node:10.15 

# Copy only what we neeed
ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}
COPY yarn.lock package.json ./
COPY integration/package.json integration/

# install deps for our integration scripts
RUN yarn

# copy over all our dependencies
COPY integration integration

# setup our integration scripts
RUN yarn workspace @chainlink/integration setup

ENV JOB_SERVER_PORT 6692
EXPOSE 6692

ENTRYPOINT ["yarn", "workspace", "@chainlink/integration", "cypressJobServer"]
FROM node:10.15 

# Copy only what we neeed
WORKDIR /chainlink
COPY yarn.lock package.json ./
COPY evm/package.json evm/
COPY evm/v0.5/package.json evm/v0.5/
COPY integration-scripts/package.json integration-scripts/

# install deps for our integration scripts
RUN yarn

# copy over all our dependencies
COPY evm evm
COPY integration-scripts integration-scripts

# setup our integration scripts
RUN yarn workspace chainlinkv0.5 setup
RUN yarn workspace chainlink setup
RUN yarn workspace @chainlink/integration-scripts setup

ENV PORT 6690
EXPOSE 6690

ENTRYPOINT [ "yarn", "workspace", "@chainlink/integration-scripts", "start-echo-server" ]

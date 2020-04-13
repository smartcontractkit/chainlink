# evm-contracts
FROM node:12.15 as evm-contracts

WORKDIR /belt
COPY ./belt/package.json ./package.json
COPY ./belt .

WORKDIR /evm-contracts
COPY ./evm-contracts/package.json ./package.json
ENV NODE_ENV=development
RUN npm install
RUN npm run install:local

COPY ./evm-contracts .
RUN npm run compile
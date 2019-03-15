FROM node:11-alpine as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

COPY package.json package.json
COPY client/package.json client/package.json
RUN yarn install && cd client/ && yarn install

ADD . .
RUN yarn build

ENTRYPOINT [ "yarn", "prod" ]
FROM node:11.12-alpine as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

COPY package.json package.json
COPY client/package.json client/package.json
RUN yarn install && cd client/ && yarn install

ADD . .
RUN yarn build

ENV NODE_ENV production
ENTRYPOINT [ "yarn", "prod" ]

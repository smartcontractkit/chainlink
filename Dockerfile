FROM node:11-alpine as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

ENV PATH /usr/src/app/node_modules/.bin:$PATH

COPY package.json /usr/src/app/package.json
RUN yarn install

ADD . /usr/src/app

ENTRYPOINT [ "yarn", "prod" ]
FROM node:10.15-alpine as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

COPY package.json yarn.lock ./
COPY client/package.json client/yarn.lock client/
RUN yarn autoinstall

ADD . .
RUN yarn build

ENV NODE_ENV production
ENTRYPOINT [ "yarn", "prod" ]

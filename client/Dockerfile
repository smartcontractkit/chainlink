# Perform a clean install of the app and then build it
FROM node:11-alpine as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

ENV PATH /usr/src/app/node_modules/.bin:$PATH

RUN npm install -g @craco/craco

COPY package.json /usr/src/app/package.json
RUN yarn install

ADD . /usr/src/app

RUN yarn build

# Copy build artifacts to a light nginx container
FROM nginx:1.15-alpine

COPY --from=builder \
  /usr/src/app/build \
  /usr/share/nginx/html

ADD nginx.conf /etc/nginx/conf.d/default.conf
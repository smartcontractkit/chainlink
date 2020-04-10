#FROM smartcontract/builder:1.0.31
# Go install
FROM golang:1.14-alpine3.11
RUN apk add --no-cache --upgrade bash
ENV GOPATH=/go
ENV GOROOT=/usr/local/go
ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

#RUN mkdir -p $GOPATH/bin
#RUN wget -qO- https://raw.githubusercontent.com/golang/dep/master/install.sh | bash
# Node install
ENV NODE_VERSION=12.15.0-r1
RUN apk add --no-cache --upgrade "nodejs=${NODE_VERSION}" "npm=${NODE_VERSION}"

# Yarn install
ENV YARN_VERSION=1.22.4
RUN npm install --global yarn@$YARN_VERSION

#
RUN apk add --no-cache --upgrade bash wget make git python

COPY . .
RUN yarn
RUN go mod download
RUN make install
#-chainlink

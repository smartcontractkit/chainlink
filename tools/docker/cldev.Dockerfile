FROM golang:1.20.5-buster

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

ADD go.* ./
RUN go mod download
RUN mkdir -p tools/bin

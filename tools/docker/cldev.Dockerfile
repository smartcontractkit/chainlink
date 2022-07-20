FROM golang:1.17-buster

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

# Do dep ensure in a cacheable step
ADD go.* ./
RUN go mod download
RUN mkdir -p tools/bin

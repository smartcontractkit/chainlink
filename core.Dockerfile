# Build
FROM golang:1.14-alpine3.11 as go-image
RUN apk add --no-cache --upgrade bash
RUN apk add --update alpine-sdk
RUN apk add --update linux-headers

ENV GOPATH=/go
ENV GOROOT=/usr/local/go
ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN mkdir /chainlink
WORKDIR /chainlink
COPY package.json ./package.json
COPY GNUmakefile ./GNUmakefile
COPY tools/bin/ldflags ./tools/bin/

COPY evm-contracts ./evm-contracts
COPY packr ./packr
COPY core ./core

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN make install-chainlink

# Copy binary
FROM alpine:3.11 as alpine-image
COPY --from=go-image /go/bin/chainlink /bin/chainlink
# Run
EXPOSE 6688
ENTRYPOINT ["chainlink"]
CMD ["local" "node"]
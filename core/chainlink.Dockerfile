# Build image: Chainlink binary
FROM golang:1.21.9-alpine as buildgo

WORKDIR /chainlink
COPY GNUmakefile package.json go.mod go.sum ./
COPY tools/bin/ldflags ./tools/bin/
RUN go mod download

# Install make, jq, and clear apk cache in the same layer
ARG COMMIT_SHA
COPY . .
RUN apk update && \
    apk add --no-cache git curl make jq bash && \
    rm -rf /var/cache/apk/*

# Build the golang binary and plugins
RUN make install-chainlink && \
    go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-feeds | xargs -I % ln -s % /chainlink-feeds && \
    go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-solana | xargs -I % ln -s % /chainlink-solana

# Build plugins
FROM golang:1.21.9-alpine as buildplugins
WORKDIR /chainlink-feeds
COPY --from=buildgo /chainlink-feeds .
RUN go install ./cmd/chainlink-feeds

WORKDIR /chainlink-solana
COPY --from=buildgo /chainlink-solana .
RUN go install ./pkg/solana/cmd/chainlink-solana

# Use a smaller, more secure final base image
FROM alpine:latest

ARG CHAINLINK_USER=root

RUN apk update && \
    apk add --no-cache ca-certificates curl postgresql-client && \
    rm -rf /var/cache/apk/* && \
    if [ "${CHAINLINK_USER}" != "root" ]; then adduser -D -u 14933 ${CHAINLINK_USER}; fi

# Copy binaries from build stage
COPY --from=buildgo /go/bin/chainlink /usr/local/bin/
COPY --from=buildplugins /go/bin/chainlink-feeds /usr/local/bin/
COPY --from=buildplugins /go/bin/chainlink-solana /usr/local/bin/
COPY --from=buildgo /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.*.so /usr/lib/
RUN chmod 755 /usr/lib/libwasmvm.*.so

USER ${CHAINLINK_USER}
WORKDIR /home/${CHAINLINK_USER}
ENV XDG_CACHE_HOME /home/${CHAINLINK_USER}/.cache
RUN mkdir -p ${XDG_CACHE_HOME}

EXPOSE 6688
ENTRYPOINT ["chainlink"]
HEALTHCHECK CMD curl -f http://localhost:6688/health || exit 1
CMD ["local", "node"]

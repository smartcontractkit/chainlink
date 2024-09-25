ARG BASE_IMAGE
ARG IMAGE_VERSION=latest
FROM ${BASE_IMAGE}:${IMAGE_VERSION} AS build-env

WORKDIR /go/testdir
RUN mkdir -p /go/testdir/integration-tests/load
COPY go.mod go.sum  ./
COPY integration-tests/go.mod integration-tests/go.sum ./integration-tests/
COPY integration-tests/load/go.mod integration-tests/load/go.sum ./integration-tests/load/
RUN cd integration-tests && go mod download
RUN cd integration-tests/load && go mod download

COPY . .

ARG SUITES=chaos soak benchmark load ccip-load

RUN /go/testdir/integration-tests/scripts/buildTests "${SUITES}"

FROM ${BASE_IMAGE}:${IMAGE_VERSION}

RUN mkdir -p /go/testdir/integration-tests/scripts
# Dependency of CosmWasm/wasmd
COPY --from=build-env /go/pkg/mod/github.com/\!cosm\!wasm/wasmvm@v*/internal/api/libwasmvm.*.so /usr/lib/
RUN chmod 755 /usr/lib/libwasmvm.*.so
COPY --from=build-env /go/testdir/integration-tests/*.test /go/testdir/integration-tests/
COPY --from=build-env /go/testdir/integration-tests/ccip-tests/*.test /go/testdir/integration-tests/
COPY --from=build-env /go/testdir/integration-tests/scripts /go/testdir/integration-tests/scripts/
RUN echo "All tests"
RUN ls -l /go/testdir/integration-tests/*.test

ENTRYPOINT ["/go/testdir/integration-tests/scripts/entrypoint"]

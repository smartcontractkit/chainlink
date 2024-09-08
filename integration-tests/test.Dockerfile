ARG BASE_IMAGE
ARG IMAGE_VERSION=latest
FROM ${BASE_IMAGE}:${IMAGE_VERSION} AS build-env

ARG SUITES=chaos migration performance reorg smoke soak benchmark

COPY . testdir/
WORKDIR /go/testdir
RUN /go/testdir/integration-tests/scripts/buildTests "${SUITES}"

FROM ${BASE_IMAGE}:${IMAGE_VERSION}

RUN mkdir -p /go/testdir/integration-tests/scripts
COPY --from=build-env /go/pkg /go/pkg
COPY --from=build-env /go/testdir/integration-tests/*.test /go/testdir/integration-tests/
COPY --from=build-env /go/testdir/integration-tests/scripts /go/testdir/integration-tests/scripts/

ENTRYPOINT ["/go/testdir/integration-tests/scripts/entrypoint"]

ARG BASE_IMAGE
ARG IMAGE_VERSION=latest
FROM ${BASE_IMAGE}:${IMAGE_VERSION} as builder

ARG SUITES=chaos migration performance reorg smoke soak benchmark

COPY . testdir/
WORKDIR /go/testdir
RUN /go/testdir/integration-tests/scripts/buildTests "${SUITES}"

# Now pull in th repo and the build run executables only
FROM ${BASE_IMAGE}:${IMAGE_VERSION}
COPY . testdir/
COPY --from=builder /go/testdir/integration-tests/*.test /go/testdir/integration-tests/
ENTRYPOINT ["/go/testdir/integration-tests/scripts/entrypoint"]

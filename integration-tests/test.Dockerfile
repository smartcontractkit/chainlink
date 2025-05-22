ARG BASE_IMAGE
ARG IMAGE_VERSION=latest
FROM ${BASE_IMAGE}:${IMAGE_VERSION} AS build-env

WORKDIR /go/testdir
# Deplyment module uses a local replace for latest code
COPY deployment/go.mod deployment/go.sum /go/testdir/deployment/
COPY go.mod go.sum  ./
COPY integration-tests/go.mod integration-tests/go.sum ./integration-tests/
RUN cd integration-tests && go mod download

COPY . .

# Get the SHA of the current commit and save it to sha.txt
RUN git rev-parse HEAD > /go/testdir/sha.txt

ARG SUITES=chaos soak benchmark ccip-tests/load

RUN /go/testdir/integration-tests/scripts/buildTests "${SUITES}"

FROM ${BASE_IMAGE}:${IMAGE_VERSION}

RUN mkdir -p /go/testdir/integration-tests/scripts
COPY --from=build-env /go/testdir/integration-tests/*.test /go/testdir/integration-tests/
COPY --from=build-env /go/testdir/integration-tests/ccip-tests/*.test /go/testdir/integration-tests/
COPY --from=build-env /go/testdir/integration-tests/scripts /go/testdir/integration-tests/scripts/
COPY --from=build-env /go/testdir/sha.txt /go/testdir/sha.txt

RUN echo "chainlink SHA used:"
RUN cat /go/testdir/sha.txt

RUN echo "All tests"
RUN ls -l /go/testdir/integration-tests/*.test

ENTRYPOINT ["/go/testdir/integration-tests/scripts/entrypoint"]

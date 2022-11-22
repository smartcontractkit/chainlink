ARG BASE_IMAGE
FROM $BASE_IMAGE

# suites example: ./integration-tests/smoke ./integration-tests/soak
ARG SUITES
ARG GINKGO_VERSION=2.5.1

COPY . testdir/
WORKDIR /go/testdir

RUN set -ex \
    && go install github.com/onsi/ginkgo/v2/ginkgo@v$GINKGO_VERSION \
    && ginkgo build $SUITES \
    && ls ./integration-tests/smoke

ENTRYPOINT ["/go/testdir/integration-tests/scripts/entrypoint"]

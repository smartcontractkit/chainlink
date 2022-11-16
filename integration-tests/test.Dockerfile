ARG BASE_IMAGE
FROM $BASE_IMAGE
COPY . testdir/
WORKDIR testdir
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.5.0 
RUN ls
RUN pwd
ENTRYPOINT ["/go/testdir/integration-tests/scripts/entrypoint"]

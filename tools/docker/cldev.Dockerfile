FROM smartcontract/builder:1.0.42

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

# Do dep ensure in a cacheable step
ADD go.* ./
RUN go mod download
RUN mkdir -p tools/bin

ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:${SRCROOT}/tools/bin

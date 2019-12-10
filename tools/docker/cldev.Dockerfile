FROM smartcontract/builder:1.0.22 

WORKDIR /chainlink

# Do dep ensure in a cacheable step
ADD go.* ./
RUN go mod download
RUN mkdir -p tools/bin

ENV PATH /go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/chainlink/tools/bin

FROM smartcontract/builder:1.0.1

# The integration tests need jq
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y jq

# Have to reintroduce ENV vars from Baidu's SGX SDK image
ENV PATH /root/.cargo/bin:/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/sgxsdk/bin:/opt/sgxsdk/bin/x64
ENV LD_LIBRARY_PATH /opt/sgxsdk/sdk_libs
ENV SGX_SDK /opt/sgxsdk

# Install goverage to capture go test coverage
RUN go get -u github.com/haya14busa/goverage

# Create the project working directory in the full GOPATH
RUN mkdir -p /go/src/github.com/smartcontractkit/chainlink/
WORKDIR /go/src/github.com/smartcontractkit/chainlink

# Do dependency installs first, since these will change less than the full
# source tree and can get cached
ADD Gopkg.* ./
RUN dep ensure -vendor-only
ADD package.json yarn.lock ./
RUN yarn install

# Copy in full source
ADD . .

CMD ["/bin/bash"]

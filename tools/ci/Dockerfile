FROM smartcontract/builder:1.0.15

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

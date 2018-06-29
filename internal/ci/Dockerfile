FROM smartcontract/builder:1.0.1

RUN apt-get update && apt-get install -y jq

# Have to reintroduce ENV vars from Baidu's SGX SDK image
ENV PATH /root/.cargo/bin:/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/sgxsdk/bin:/opt/sgxsdk/bin/x64
ENV LD_LIBRARY_PATH /opt/sgxsdk/sdk_libs
ENV SGX_SDK /opt/sgxsdk

ARG COMMIT_SHA
ARG ENVIRONMENT
ARG SGX_ENABLED
ARG SGX_SIMULATION

ADD . /go/src/github.com/smartcontractkit/chainlink
WORKDIR /go/src/github.com/smartcontractkit/chainlink

RUN dep ensure -vendor-only
RUN go get -u github.com/haya14busa/goverage
RUN yarn install

CMD ["/bin/bash"]

# Build Chainlink with SGX
FROM smartcontract/builder:1.0.0 as builder

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
RUN make build

# Final layer: ubuntu with aesm and chainlink binaries (executable + enclave)
FROM ubuntu:16.04

# Install AESM
RUN apt-get update && \
  apt-get install -y \
    curl \
    libssl-dev \
    libcurl4-openssl-dev \
    libprotobuf-dev \
    kmod \
    libprotobuf-c0-dev \
    libxml2-dev \
    libssl1.0.0

RUN /usr/sbin/useradd aesmd 2>/dev/null

RUN mkdir -p /var/opt/aesmd && chown aesmd.aesmd /var/opt/aesmd
RUN mkdir -p /var/run/aesmd && chown aesmd.aesmd /var/run/aesmd

COPY --from=builder /opt/sgxsdk/lib64/libsgx*.so /usr/lib/
COPY --from=builder /opt/intel/ /opt/intel/

# Copy chainlink enclave+stub from build image

# XXX: This is structured like so because wheb SGX_ENABLED is 'no', no sgx
# targets get produced, so here we specify the SGX targets as a wildcard with
# chainlink so COPY doesn't fail when they're missing. Then move them all into
# place.
ARG ENVIRONMENT
COPY --from=builder \
  /go/src/github.com/smartcontractkit/chainlink/chainlink \
  /go/src/github.com/smartcontractkit/chainlink/sgx/target/$ENVIRONMENT/*.so \
  /tmp/
RUN mv /tmp/chainlink /root/ && \
  mv /tmp/enclave.signed.so /root/ 2>/dev/null || true && \
  mv /tmp/libadapters.so /usr/lib/ 2>/dev/null || true && \
  rm -f /tmp/*.so

# Launch chainlink via a small script that watches AESM + Chainlink
ARG SGX_SIMULATION
ENV SGX_SIMULATION $SGX_SIMULATION
WORKDIR /root
COPY ./chainlink-launcher.sh /root
RUN chmod +x ./chainlink-launcher.sh

EXPOSE 6688
EXPOSE 6689
ENTRYPOINT ["./chainlink-launcher.sh"]
CMD ["node"]

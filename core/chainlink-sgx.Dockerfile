# Build Chainlink with SGX
FROM smartcontract/builder:1.0.25 as builder

# Have to reintroduce ENV vars from builder image
ENV PATH /root/.cargo/bin:/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/sgxsdk/bin:/opt/sgxsdk/bin/x64
ENV LD_LIBRARY_PATH /opt/sgxsdk/sdk_libs
ENV SGX_SDK /opt/sgxsdk

WORKDIR /chainlink
COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Install yarn dependencies
COPY yarn.lock package.json ./
COPY explorer/client/package.json ./explorer/client/
COPY explorer/package.json ./explorer/
COPY operator_ui/package.json ./operator_ui/
COPY feeds_ui/package.json ./feeds_ui/
COPY styleguide/package.json ./styleguide/
COPY tools/json-api-client/package.json ./tools/json-api-client/
COPY tools/local-storage/package.json ./tools/local-storage/
COPY tools/redux/package.json ./tools/redux/
COPY tools/ts-test-helpers/package.json ./tools/ts-test-helpers/
COPY evm/v0.5/package.json ./evm/v0.5/
COPY evm/package.json ./evm/
RUN make yarndep

# Do go mod download in a cacheable step
ADD go.mod go.sum ./
RUN go mod download

# Env vars needed for chainlink sgx build
ARG COMMIT_SHA
ARG ENVIRONMENT
ENV SGX_ENABLED yes
ARG SGX_SIMULATION

# Install chainlink
ADD . ./
RUN make install-chainlink

# Final layer: ubuntu with aesm and chainlink binaries (executable + enclave)
FROM ubuntu:18.04

# Install AESM
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && \
  apt-get install -y \
    ca-certificates \
    curl \
    kmod \
    libcurl4-openssl-dev \
    libprotobuf-c0-dev \
    libprotobuf-dev \
    libssl-dev \
    libssl1.0.0 \
    libxml2-dev

RUN /usr/sbin/useradd aesmd 2>/dev/null

RUN mkdir -p /var/opt/aesmd && chown aesmd.aesmd /var/opt/aesmd
RUN mkdir -p /var/run/aesmd && chown aesmd.aesmd /var/run/aesmd

COPY --from=builder /opt/sgxsdk/lib64/libsgx*.so /usr/lib/
COPY --from=builder /opt/intel/ /opt/intel/

# Copy chainlink enclave+stub from build image
ARG ENVIRONMENT
COPY --from=builder /go/bin/chainlink /usr/local/bin/
COPY --from=builder \
  /chainlink/core/sgx/target/$ENVIRONMENT/libadapters.so \
  /usr/lib/
COPY --from=builder \
  /chainlink/core/sgx/target/$ENVIRONMENT/enclave.signed.so \
  /root/

# Launch chainlink via a small script that watches AESM + Chainlink
ARG SGX_SIMULATION
ENV SGX_SIMULATION $SGX_SIMULATION
WORKDIR /root
COPY core/chainlink-launcher-sgx.sh .
RUN chmod +x ./chainlink-launcher-sgx.sh

EXPOSE 6688
ENTRYPOINT ["./chainlink-launcher-sgx.sh"]
CMD ["local", "node"]

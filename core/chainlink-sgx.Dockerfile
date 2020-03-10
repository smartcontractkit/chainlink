# Build Chainlink with SGX
FROM smartcontract/builder:1.0.30 as builder

WORKDIR /chainlink
COPY GNUmakefile VERSION ./
COPY tools/bin/ldflags ./tools/bin/

# Do dep ensure in a cacheable step
ADD go.* ./
RUN go mod download

# And yarn likewise
COPY yarn.lock package.json .yarnrc ./
COPY .yarn .yarn
COPY operator_ui/package.json ./operator_ui/
COPY styleguide/package.json ./styleguide/
COPY tools/json-api-client/package.json ./tools/json-api-client/
COPY tools/local-storage/package.json ./tools/local-storage/
COPY tools/redux/package.json ./tools/redux/
COPY tools/ts-helpers/package.json ./tools/ts-helpers/
COPY belt/package.json ./belt/
COPY belt/bin ./belt/bin
COPY evm-test-helpers/package.json ./evm-test-helpers/
COPY evm-contracts/package.json ./evm-contracts/
RUN make yarndep

# Env vars needed for chainlink sgx build
ARG COMMIT_SHA
ARG ENVIRONMENT
ENV SGX_ENABLED yes
ARG SGX_SIMULATION

# Install chainlink
COPY tsconfig.cjs.json tsconfig.es6.json ./
COPY operator_ui ./operator_ui
COPY styleguide ./styleguide
COPY tools/json-api-client ./tools/json-api-client
COPY tools/local-storage ./tools/local-storage
COPY tools/redux ./tools/redux
COPY tools/ts-helpers ./tools/ts-helpers
COPY belt ./belt
COPY belt/bin ./belt/bin
COPY evm-test-helpers ./evm-test-helpers
COPY evm-contracts ./evm-contracts
COPY core core
COPY packr packr

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

FROM ubuntu:20.04
WORKDIR /solana

ADD https://github.com/solana-labs/solana/releases/download/v1.16.23/solana-release-x86_64-unknown-linux-gnu.tar.bz2 .
RUN tar jxf solana-release-x86_64-unknown-linux-gnu.tar.bz2
ENV PATH=/solana/solana-release/bin:$PATH

ENTRYPOINT [ "solana-test-validator" ]

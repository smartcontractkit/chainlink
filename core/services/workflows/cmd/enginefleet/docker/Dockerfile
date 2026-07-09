# Runtime image for the enginefleet binary.
#
# The binary is compiled separately (see build.sh) and copied in here, rather
# than compiled in-image: the chainlink monorepo builds fully offline from the
# local Go module cache, so building via a throwaway golang container that
# mounts that cache is far faster than an in-image `go mod download` and needs
# no private-module auth.
#
# The WASM host (wasmtime) is CGO and statically links a Rust archive, so the
# only runtime deps are glibc/libgcc, both present in bookworm-slim.

# Stage: build jemalloc WITH heap profiling. Debian's libjemalloc2 is compiled
# without --enable-prof, so we build it from source to get jeprof profiling of
# native (CGo/wasmtime) malloc allocations.
FROM debian:bookworm-slim AS jemalloc
RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential ca-certificates wget bzip2 \
    && rm -rf /var/lib/apt/lists/*
ARG JEMALLOC_VERSION=5.3.0
RUN wget -qO /tmp/jemalloc.tar.bz2 \
        "https://github.com/jemalloc/jemalloc/releases/download/${JEMALLOC_VERSION}/jemalloc-${JEMALLOC_VERSION}.tar.bz2" \
    && mkdir -p /tmp/jemalloc \
    && tar xjf /tmp/jemalloc.tar.bz2 -C /tmp/jemalloc --strip-components=1 \
    && cd /tmp/jemalloc \
    && ./configure --enable-prof --prefix=/opt/jemalloc \
    && make -j"$(nproc)" && make install

# Final image. Extra tools for debugging native (C/CGo/wasmtime) memory, which
# Go's pprof cannot see:
#   - jemalloc + jeprof: native heap profiling via LD_PRELOAD (see run.sh -j)
#   - binutils/perl:     symbolization for jeprof
#   - graphviz:          jeprof --svg/--pdf output
#   - procps:            pmap for inspecting mmap'd regions (wasmtime linear memory)
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        perl \
        binutils \
        graphviz \
        procps \
    && rm -rf /var/lib/apt/lists/*

COPY --from=jemalloc /opt/jemalloc/lib/libjemalloc.so.2 /usr/local/lib/
COPY --from=jemalloc /opt/jemalloc/bin/jeprof /usr/local/bin/jeprof
RUN ldconfig

COPY enginefleet /usr/local/bin/enginefleet

ENTRYPOINT ["/usr/local/bin/enginefleet"]

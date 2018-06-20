# Start from the sgx-rust base image
FROM baiduxlab/sgx-rust:0.9.8

# Add all the things we need to build chainlink
RUN apt-get update && apt-get install -y curl git gcc libssl1.0.0 build-essential

# Install go 1.10.1
RUN cd /usr/local && \
  curl -sS https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz | tar -xz
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
RUN mkdir -p $GOPATH/bin && mkdir $GOPATH/src

# Dep is used to manage dependencies
RUN curl -sS https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

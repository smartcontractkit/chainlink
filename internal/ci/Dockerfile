FROM ubuntu:16.04

# Go
RUN apt-get update && apt-get install -y curl git gcc libssl1.0.0 build-essential

RUN cd /usr/local && \
  curl -sS https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz | tar -xz

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
RUN mkdir -p $GOPATH/bin && mkdir $GOPATH/src

# Go dep
RUN curl -sS https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Goverage
RUN go get -u github.com/haya14busa/goverage

# Yarn
RUN apt-key adv --fetch-keys http://dl.yarnpkg.com/debian/pubkey.gpg && \
    apt-get update -qq \
    echo "deb http://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
    apt-get install -y -qq yarn

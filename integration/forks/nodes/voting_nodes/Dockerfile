FROM ethereum/client-go:alltools-latest
RUN apk add curl bash

COPY . ./root
RUN chmod +x /root/bin/*

EXPOSE 30303
EXPOSE 30303/udp

ENV PATH="$PATH:/root/bin"

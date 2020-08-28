FROM smartcontract/chainlink:latest

ENV CHAINLINK_TLS_PORT 0
ENV CHAINLINK_PORT 6688

RUN mkdir /clroot
RUN chmod -R 700 /clroot

COPY apicredentials /
COPY password.txt /
COPY start.sh /
WORKDIR /

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./start.sh"]

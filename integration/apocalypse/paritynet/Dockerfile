FROM parity/parity:v2.7.2-stable
COPY . /paritynet
USER root
RUN chown -R parity /paritynet
USER parity
WORKDIR /paritynet
ENTRYPOINT /paritynet/start.sh

EXPOSE 8546
EXPOSE 8545
EXPOSE 8180
EXPOSE 8080
EXPOSE 30303
EXPOSE 30303/udp

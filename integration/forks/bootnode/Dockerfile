FROM ethereum/client-go:alltools-latest

COPY . ./src

EXPOSE 30303
EXPOSE 30303/udp

CMD bootnode -nodekey src/boot.key -addr :30303

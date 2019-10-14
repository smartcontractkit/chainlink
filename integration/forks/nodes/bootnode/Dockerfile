FROM ethereum/client-go:alltools-latest
RUN apk add curl bash

COPY . ./src

EXPOSE 30303
EXPOSE 30303/udp

CMD bootnode -nodekey src/boot.key --verbosity 9 -addr :30303

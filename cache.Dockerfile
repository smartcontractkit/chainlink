FROM smartcontract/builder:1.0.30

WORKDIR /chainlink
COPY . . 
RUN make gen-builder-cache

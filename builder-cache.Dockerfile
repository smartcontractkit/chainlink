FROM smartcontract/builder:1.0.34

WORKDIR /chainlink
COPY . . 
RUN make gen-builder-cache

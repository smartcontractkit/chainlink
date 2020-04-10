FROM smartcontract/builder:1.0.31

WORKDIR /chainlink
COPY . . 
RUN make gen-builder-cache

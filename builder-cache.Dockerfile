FROM smartcontract/builder:1.0.33

WORKDIR /chainlink
COPY . . 
RUN make gen-builder-cache

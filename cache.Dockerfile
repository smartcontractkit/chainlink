FROM smartcontract/builder:1.0.30

WORKDIR /chainlink
COPY . . 
RUN yarn
RUN go mod download

FROM smartcontract/builder:1.0.31

COPY . .
RUN yarn
RUN go mod download
RUN make install-chainlink

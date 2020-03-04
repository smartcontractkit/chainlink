FROM smartcontract/builder:1.0.30
COPY . . 
RUN yarn
RUN go mod download

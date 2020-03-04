FROM smartcontract/builder:1.0.30
RUN apt-get update && apt-get install -y libudev-dev libusb-dev libusb-1.0-0

WORKDIR /chainlink
COPY . . 
RUN yarn
RUN go mod download

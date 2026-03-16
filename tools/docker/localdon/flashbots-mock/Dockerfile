FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod main.go ./
RUN go build -o flashbots-mock .

FROM alpine:3.19
COPY --from=build /app/flashbots-mock /usr/local/bin/
ENTRYPOINT ["flashbots-mock"]

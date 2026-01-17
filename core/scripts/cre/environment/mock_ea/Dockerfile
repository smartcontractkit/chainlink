# Mock External Adapter for Streams Trigger E2E Testing
# 
# Build:
#   docker build -t mock-ea:latest .
#   docker tag mock-ea:latest localhost:5001/mock-ea:latest
#   docker push localhost:5001/mock-ea:latest
#
# Run:
#   docker run -p 8080:8080 \
#     -e BTC_USD_PRICE=65000 \
#     -e ETH_USD_PRICE=3500 \
#     -e LINK_USD_PRICE=15 \
#     mock-ea:latest

FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copy source
COPY main.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-ea .

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/mock-ea .

# Default environment variables
# TEST/USD price is hardcoded to MAGIC_NUMBER (424242) in code
ENV PORT=8080
ENV LOG_LEVEL=info

EXPOSE 8080

CMD ["./mock-ea"]

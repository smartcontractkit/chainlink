# Gateway Handler Load Test

Load test for the gateway HTTP handler (`v2.NewGatewayHandler`). Measures how many outbound HTTP requests the handler can process in parallel via the full `HandleNodeMessage` -> `makeOutgoingRequest` -> `httpClient.Send` path.

## Build

```bash
cd core/scripts
go build -o gateway-loadtest-client ./gateway-loadtest/cmd/loadtest/
go build -o gateway-loadtest-server ./gateway-loadtest/cmd/destserver/
```

## Usage

### 1. Start the destination server

```bash
./gateway-loadtest-server -port 8080 -latency 50 -body-size 256
```

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | 8080 | Listen port |
| `-latency` | 0 | Simulated response latency (ms) |
| `-body-size` | 256 | Response body size (bytes) |

### 2. Run the load test

```bash
./gateway-loadtest-client -n 500 -url http://127.0.0.1:8080 -dest-port 8080
```

| Flag | Default | Description |
|------|---------|-------------|
| `-n` | 100 | Number of parallel requests |
| `-url` | `http://127.0.0.1:8080` | Destination server URL |
| `-timeout` | 30000 | Per-request timeout (ms) |
| `-dest-port` | 8080 | Destination port (added to HTTP client AllowedPorts) |

Increasing the limit of conns on a mac: `sudo sysctl -w kern.ipc.somaxconn=20000`

### Example output

```
=== Load Test Results ===
Total requests:   500
Errors:           0
Total duration:   120.5ms
Throughput:       4149.4 req/s

--- Latency Distribution ---
Min:  12.3ms
P50:  55.1ms
P90:  98.2ms
P95:  105.7ms
P99:  118.4ms
Max:  120.1ms
Mean: 58.9ms
```

## Test scenarios

**Baseline (no latency):**
```bash
./gateway-loadtest-server -latency 0 &
./gateway-loadtest-client -n 1000
```

**Slow server:**
```bash
./gateway-loadtest-server -latency 500 &
./gateway-loadtest-client -n 1000
```

**Large responses:**
```bash
./gateway-loadtest-server -latency 10 -body-size 20000 &
./gateway-loadtest-client -n 500
```

# Bridge Status Reporter Documentation

## Overview

The Bridge Status Reporter is a service that continuously monitors External Adapter health by polling their status endpoints and emitting telemetry events. This provides visibility into bridge connectivity, performance, and configuration across your Chainlink network.

## Configuration

### Node Configuration

Add the following to your node's TOML configuration:

```toml
[BridgeStatusReporter]
Enabled = true                    # Enable the service
StatusPath = "/status"            # Path to append to bridge URLs  
PollingInterval = "5m"            # How often to poll bridges
IgnoreInvalidBridges = true       # Skip bridges with HTTP errors
IgnoreJoblessBridges = false      # Skip bridges with no associated jobs
```

### External Adapter Requirements

Your External Adapters must implement a `/status` endpoint (or the path specified in `StatusPath`) that returns bridge status information.

#### Example Status Endpoint Response

```json
{
  "bridge_name": "my-bridge",
  "adapter_name": "crypto-price-adapter",
  "adapter_version": "1.2.3",
  "adapter_uptime_seconds": 86400.5,
  "default_endpoint": "crypto",
  "runtime": {
    "node_version": "18.19.0",
    "platform": "linux",
    "architecture": "x64",
    "hostname": "adapter-server-01"
  },
  "metrics": {
    "enabled": true
  },
  "endpoints": [
    {
      "name": "crypto",
      "aliases": ["price", "market"],
      "transports": ["http", "https"]
    }
  ],
  "configuration": [
    {
      "name": "API_KEY",
      "value": "[REDACTED]",
      "type": "string",
      "description": "API key for data provider",
      "required": true,
      "default_value": "",
      "custom_setting": false,
      "env_default_override": "CRYPTO_API_KEY"
    }
  ]
}
```

## Bridge Registration

Bridges are automatically discovered from your node's bridge registry. The service will:

1. **Query all registered bridges** from the database
2. **Filter active bridges** (unless `IgnoreJoblessBridges` is false)
3. **Poll each bridge's status endpoint** at the configured interval
4. **Emit telemetry events** for successful responses
5. **Log errors** for failed requests (optionally ignored with `IgnoreInvalidBridges`)

## Telemetry Events

The service emits `BridgeStatusEvent` protobuf messages containing:

### Bridge Identification
- `bridge_name` - Name from bridge registry
- `adapter_name` - External adapter identifier  
- `adapter_version` - Version string
- `default_endpoint` - Primary endpoint name

### Runtime Information
- `node_version` - Runtime version (Node.js, etc.)
- `platform` - Operating system
- `architecture` - CPU architecture
- `hostname` - Server hostname

### Operational Data
- `adapter_uptime_seconds` - Adapter uptime in seconds
- `endpoints` - Available endpoints with aliases and transports
- `configuration` - Configuration parameters (values may be redacted)
- `jobs` - Associated Chainlink jobs using this bridge
- `metrics` - Metrics collection status

## Troubleshooting

### Common Issues

**Bridge Not Appearing in Telemetry**
- Verify bridge is registered: Check `/v2/bridges` API endpoint
- Check bridge has associated jobs (if `IgnoreJoblessBridges = false`)
- Ensure bridge URL is accessible from the node and returning data
- Check bridge is correctly configured and logging errors
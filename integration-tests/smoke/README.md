## Smoke tests (local environments)

These products are using local `testcontainers-go` environments:
- RunLog (Direct request)
- Cron
- Flux
- VRFv1
- VRFv2

### Usage
```
go test -v -run ${TestName}
```
### Re-using environments
Configuration is still WIP, but you can make your environment re-usable by providing JSON config.

Create `test_env.json` in the same dir
```
export TEST_ENV_CONFIG_PATH=test_env.json
```

Here is an example for 3 nodes cluster
```
{
  "networks": [
    "epic"
  ],
  "mockserver": {
    "container_name": "mockserver",
    "external_adapters_mock_urls": [
      "/epico1"
    ]
  },
  "geth": {
    "container_name": "geth"
  },
  "nodes": [
    {
      "container_name": "cl-node-0",
      "db_container_name": "cl-db-0"
    },
    {
      "container_name": "cl-node-1",
      "db_container_name": "cl-db-1"
    },
    {
      "container_name": "cl-node-2",
      "db_container_name": "cl-db-2"
    }
  ]
}
```

### Running against Live Testnets
1. Prepare your `overrides.toml` file with selected network and CL image name and version and save anywhere inside `integration-tests` folder.
```toml
[ChainlinkImage]
image="your-image"
version="your-version"

[Network]
selected_networks=["polygon_mumbai"]

[Network.RpcHttpUrls]
polygon_mumbai=["https://http.endpoint.com"]

[Network.RpcWsUrls]
polygon_mumbai=["wss://ws.endpoint.com"]

[Network.WalletKeys]
polygon_mumbai=["my_so_private_key"]
```
Then execute:
```bash
go test -v -run ${TestName}
```



### Debugging CL client API calls
```bash
export CL_CLIENT_DEBUG=true
```
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

```
SELECTED_NETWORKS=<Chain Name> \
<Chain Name>_KEYS= \
<Chain Name>_URLS= \
<Chain Name>_HTTP_URLS= \
go test -v -run ${TestName}
```



### Debugging CL client API calls
```
export CL_CLIENT_DEBUG=true
```
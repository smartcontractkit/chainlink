## Examples
This directory shows some examples on how to assemble different `Chainlink` services, connect and test them

You can use [direnv](https://direnv.net/) or raw `.envrc` files to set up common vars
```
export CTF_CONFIGS=smoke.toml
export CTF_LOG_LEVEL=debug
export CTF_USE_CACHED_OUTPUTS=true
export CTF_LOKI_STREAM=true
export LOKI_TENANT_ID=promtail
export LOKI_URL=http://host.docker.internal:3030/loki/api/v1/push
export TESTCONTAINERS_RYUK_DISABLED=true
export RESTY_DEBUG=false
```

### CLI
```
go get github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && go install github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && mv ~/go/bin/cmd ~/go/bin/ctf
```

### Local observability stack
```
ctf obs up
```

### Multi-node, Multi-network example
```
go test -v -run TestDON
```

### Overriding configs
You can override any configuration by providing more `TOML` files
```
export CTF_CONFIGS=smoke.toml,smoke-another-network.toml
```
Changes will be applied right to left

### Caching
You can re-use already deployed environment and contracts like this
1. Run your test once
2. Change the configuration to cache
```
export CTF_CONFIGS=smoke-cache.toml
export CTF_USE_CACHED_OUTPUTS=true
```
3. Develop your test on cached or external environment, you can override `.out` fields in cached config to connect to any other environment, staging ,etc
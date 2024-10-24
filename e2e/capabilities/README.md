## Examples
This directory shows some examples on how to assemble different `Chainlink` services, connect and test them

You can use [direnv](https://direnv.net/) or raw `.envrc` files to set up common vars
```
export CTF_LOG_LEVEL=info
export CTF_LOKI_STREAM=true
export LOKI_TENANT_ID=promtail
export LOKI_URL=http://host.docker.internal:3030/loki/api/v1/push
export TESTCONTAINERS_RYUK_DISABLED=true
export RESTY_DEBUG=false
```
You can read more in [docs](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/README.md)

### CLI
```
go get github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && go install github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && mv ~/go/bin/cmd ~/go/bin/ctf
```

### Local observability stack
```
ctf obs up
```

### DON + Anvil example
Add env vars to your `.envrc` and run
```
export CTF_CONFIGS=smoke.toml
export PRIVATE_KEY="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

go test -v -run TestDON
```

### DON + Avalanche Fuji example
Add env vars to your `.envrc` and run
```
export CTF_CONFIGS=smoke.toml,smoke-fuji.toml
export PRIVATE_KEY="..."

go test -v -run TestDON
```

### Overriding configs
You can override any configuration by providing more `TOML` files
```
export CTF_CONFIGS=smoke.toml,smoke-another-network.toml
```
Changes will be applied right to left

## Default CLNode credentials
UI login/password:
```
notreal@fakeemail.ch
fj293fbBnlQ!f9vNs
```

### Caching
You can re-use already deployed environment and contracts like this
1. Run your test once
2. Change the configuration to cache
```
export CTF_CONFIGS=smoke-cache.toml
```
3. Develop your test on cached or external environment, you can override `.out` fields in cached config to connect to any other environment, staging ,etc
4. You can control caching of each component by changing `use_cache = true|false`
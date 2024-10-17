## Examples
This directory shows some examples on how to assemble different `Chainlink` services, connect and test them

You can use [direnv](https://direnv.net/) or raw `.envrc` files to set up common vars
```
export CTF_LOG_LEVEL=debug
export CTF_CONFIGS=smoke.toml
export CTF_USE_CACHED_OUTPUTS=true
export CTF_LOKI_STREAM=false
export TESTCONTAINERS_RYUK_DISABLED=true
```

### Multi-node, Multi-network example
```
go test -v -run TestMultiNodeMultiNetwork
```
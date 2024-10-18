## Examples
This directory shows some examples on how to assemble different `Chainlink` services, connect and test them

You can use [direnv](https://direnv.net/) or raw `.envrc` files to set up common vars
```
export CTF_CONFIGS=smoke.toml,smoke-another-network.toml
export CTF_LOG_LEVEL=info
export CTF_USE_CACHED_OUTPUTS=false
export CTF_LOKI_STREAM=true
export LOKI_TENANT_ID=promtail
export LOKI_URL=http://host.docker.internal:3030/loki/api/v1/push
export TESTCONTAINERS_RYUK_DISABLED=true
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
go test -v -run TestMultiNodeMultiNetwork
```
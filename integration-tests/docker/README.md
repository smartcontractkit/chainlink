## Docker environment
This folder contains Chainlink cluster environment created with `testcontainers-go`

### CLI for Local Testing Environment

The command-line interface (CLI) located at `./integration-tests/docker/cmd/test_env.go` can be utilized to initiate a local testing environment. It is intended to replace Docker Compose in the near future.


Example: 
```
# Set required envs
export CHAINLINK_IMAGE="<chainlink_node_docker_image_path>"
export CHAINLINK_VERSION="<chainlink_node_docker_image_version>" 
# Stream logs to Loki
export LOKI_TOKEN=...
export LOKI_URL=https://${loki_host}/loki/api/v1/push

cd ./integration-tests/docker/cmd

go run test_env.go start-env cl-cluster
```
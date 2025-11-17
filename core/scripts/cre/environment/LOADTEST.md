# Local Workflow Environment Setup Guide

This guide walks you through setting up a local environment to test workflows with mock capabilities. It includes steps to build the necessary Docker image, configure the environment, compile and register a workflow, and trigger it locally.

Note that this is still a work in progress!
---

## Prerequisites

Ensure you have the following installed:

- [Brotli](https://github.com/google/brotli) Install with Homebrew: `brew install brotli`
- GitHub token with read access to the [`capabilities`](https://github.com/smartcontractkit/capabilities) repo
- [`docker-mac-net-connect`](https://github.com/chipmk/docker-mac-net-connect) for macOS routing to Docker IPs. Install with Homebrew `brew install chipmk/tap/docker-mac-net-connect` and start as a service `sudo brew services start chipmk/tap/docker-mac-net-connect`

---

## 1. Build the Docker Plugin-testing Image

We’ll use the `-testing` image variant, which includes the required capabilities and mock support. Note that this installations relays on the fact that the mock capability is already installed in the image.

```bash
CL_INSTALL_PRIVATE_PLUGINS=true \
CL_INSTALL_TESTING_PLUGINS=true \
GITHUB_TOKEN=<YOUR_GITHUB_TOKEN> \
make docker-plugins
```

---

## 2. Configure the Environment

Navigate to the environment directory:

```bash
cd core/scripts/cre/environment
```

Update the node image in the CTF config for all the nodesets and also make sure you have the [job-distributor](https://github.com/smartcontractkit/job-distributor) image:

```toml
# File: configs/workflow-load.toml
[nodesets.node_specs.node]
image = "<CL NODE IMAGE WE JUST BUILT>"

# Update the job distributor image as well
[jd]
image = "<JOB DISTRIBUTOR IMAGE>"
```

> We are deploying locally via Docker. There are known CRIB issues that will be addressed separately at a later date.

---

## 3. Start the Environment

```bash
CTF_CONFIGS=./configs/workflow-load.toml \
go run main.go env start --topology=mock --extra-allowed-gateway-ports=16000 --with-beholder
```

---

## 4. Build the Workflow Binary

We’ll use the `fetchtrueusd` workflow.

1. Clone or pull the [capabilities repo](https://github.com/smartcontractkit/capabilities)
2. Navigate to the workflow command directory:

```bash
cd workflows/fetchtrueusd/cmd
```

3. Compile the binary:

```bash
GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0 go build -o fetchtrueusd
```

4. Compress with Brotli:

```bash
brotli -v fetchtrueusd
```

5. Encode to base64:

```bash
cat fetchtrueusd.br | base64 > fetchtrueusd.br.base64
```

> Eventually, a single command will automate this process.

---

## 5. Upload Assets to MinIO

1. Switch back to the environment directory:

```bash
cd core/scripts/cre/environment
```

2. Copy the base64-encoded binary and rename:

```bash
cp <PATH TO BINARY>/fetchtrueusd.br.base64 fetchtrueusd.br
```

3. Create an empty YAML config:

```bash
touch empty.yaml
```

4. Upload both to MinIO:

```bash
go run main.go minio upload fetchtrueusd.br empty.yaml
```

---

## 6. Register the Workflow

```bash
go run main.go workflow register \
  --binary-url="http://minio:16000/default/fetchtrueusd.br" \
  --config-url="http://minio:16000/default/empty.yaml" \
  --secrets-url="http://minio:16000/default/empty.yaml" \
  --id="0089c0071e8c5b535ebeab3f5102091f3a657daf2bb1778eea67f4a12b82c2cb" \
  --name="fetchtrueusd"
```

---

## 7. Register Required Capabilities

The workflow uses:

- `cron-trigger@1.0.0`
- `write_ethereum-testnet-sepolia@1.0.0`

### Register in the On-Chain Registry

```bash
go run main.go registry create --name="cron-trigger" --version="1.0.0" --type="trigger" --don-id=2
go run main.go registry create --name="write_ethereum-testnet-sepolia" --version="1.0.0" --type="target" --don-id=2
```

---

## 8. Register Mock Capabilities on Nodes

You’ll need to communicate with Docker container IPs directly.

### For macOS: Install Network Bridge

Note! make sure that docker-mac-net-connect is installed and running
```bash
brew install chipmk/tap/docker-mac-net-connect
brew services start chipmk/tap/docker-mac-net-connect
```

At the moment we have to fetch the capabilities-node container image manually, best way to do it is by running 
```base
docker inspect capabilities-node0 | grep "IPAddress"
docker inspect capabilities-node1 | grep "IPAddress"
docker inspect capabilities-node1 | grep "IPAddress"
```
These IP addresses are required in order to connect to the mock capability, the cli command will expect a list of IP:PORT values. The mock capability service is exposed on port 7777 so that will be the port for all IPs

### Register Mock Implementations

```bash
go run main.go mock create \
  --id="write_ethereum-testnet-sepolia@1.0.0" \
  --description="mock target" \
  --type="target" \
  --addresses="<ADDRESSES FROM THE PREVIOUS STEP IN FORMAT: IP:PORT,IP:PORT,IP:PORT>"


go run main.go mock create \
  --id="cron-trigger@1.0.0" \
  --description="mock trigger" \
  --type="trigger" \
  --addresses="<ADDRESSES FROM THE PREVIOUS STEP IN FORMAT: IP:PORT,IP:PORT,IP:PORT>"
 
# This is an example command to see the expected format of the addresses
#  go run main.go mock create \
#  --id="cron-trigger@1.0.0" \
#  --description="mock trigger" \
#  --type="trigger" \
#  --addresses="192.168.48.13:7777,192.168.48.14:7777,192.168.48.15:7777"
```

---

## 9. ⏱Trigger the Workflow

Trigger the cron capability to start workflow execution:

```bash
go run main.go mock trigger \
  --id="cron-trigger@1.0.0" \
  --type="cron" \
  --frequency="30s" \
  --duration="10m" \
  --addresses="<ADDRESSES FROM THE PREVIOUS STEP IN FORMAT: IP:PORT,IP:PORT,IP:PORT>"
```

---

## ✅ You're Done

You now have a local environment with:

- A running mock topology
- A registered `fetchtrueusd` workflow
- Required capabilities installed
- Periodic triggers in place

You're ready to observe and iterate on workflow execution locally.

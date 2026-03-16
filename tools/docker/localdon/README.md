# Local DON (4-node OCR2 SVR Median)

A Docker-based 4-node Chainlink DON running an OCR2 median feed with SVR dual transmission on a local Anvil chain. Useful for local development, testing version upgrades, and validating SVR fixes.

## Prerequisites

- Docker (with compose v2)
- Go 1.23+ (for the setup script and local-dev mode)
- [Foundry](https://book.getfoundry.sh/getting-started/installation) (`cast` CLI, for monitoring)

## Quick start (all Docker)

For **dual transmission** (primary + secondary txs), use `setup-full` so secondary ETH keys exist before nodes start (TXMv2 only initializes stores for keys present at startup):

```bash
# Start services, create secondary keys, restart nodes, then deploy and configure
make setup-full
```

Or step by step:

```bash
make up
cd ../../../core/scripts && go run ./setup-localdon --keys-only && cd -
docker compose restart node-1 node-2 node-3 node-4
# Wait for nodes to be ready, then:
make setup
```

**Verify:**

```bash
# Primary: on-chain rounds (address printed by setup)
cast call <OCR2_ADDRESS> 'latestAnswer()(int256)' --rpc-url http://localhost:8545

# Secondary: mock should report increasing count
curl http://localhost:9090/health   # {"status":"ok","secondaryTxCount":N}

make logs
```

If you already have a DON with secondary keys (e.g. after a previous `setup-full`) or only need primary transmission, you can use `make up` then `make setup`.

## Local-dev mode (node-4 from source)

Run node-4 from local source for fast iteration and debugging. The other 3 nodes + infra remain in Docker.

```bash
# Start Docker services (postgres, anvil, flashbots mock, 3 nodes; node-4 disabled)
make local-dev-up

# In a separate terminal: start node-4 from local source
make local-dev-node

# Once node-4 is healthy, deploy contracts and create jobs
make local-dev-setup
```

To restart just node-4 after code changes, Ctrl-C the `local-dev-node` process and re-run `make local-dev-node`. No need to re-deploy contracts or restart Docker services - just delete the node-4 job and re-run setup (or do a full `make clean && make local-dev-up` for a clean slate).

## Makefile targets

| Target | Description |
|--------|-------------|
| `make up` | Start all services (4 Docker nodes), wait for health checks |
| `make setup` | Deploy contracts, fund nodes, create jobs (uses existing keys; for dual tx use setup-full first time) |
| `make setup-full` | Up + create secondary ETH keys + restart nodes + setup (required for dual transmission to work) |
| `make down` | Stop all Docker services |
| `make clean` | Stop all Docker services and remove volumes |
| `make restart` | Full restart: clean + setup-full |
| `make logs` | Follow logs from Docker services |
| `make build` | Build chainlink image from local source (optional) |
| `make local-dev-up` | Start Docker services with node-4 disabled + P2P port exposed |
| `make local-dev-setup` | Run setup script in local-dev mode |
| `make local-dev-node` | Run node-4 locally from source (`go run`) |

## Architecture

```
docker-compose.yml
+-- postgres        (1 instance, 4 databases)
+-- anvil           (local EVM, chain ID 31337, 1s block time)
+-- flashbots-mock  (mock OFA endpoint for secondary transmissions)
+-- node-1          (chainlink 2.36.0, acts as P2P bootstrap peer)
+-- node-2          (chainlink 2.36.0)
+-- node-3          (chainlink 2.36.0)
+-- node-4          (chainlink 2.36.0, or local source in local-dev mode)

setup script (Go, runs on host)
+-- optional: --keys-only creates secondary ETH key per node (then restart nodes so TXMv2 sees them)
+-- deploys LINK token + OCR2Aggregator to Anvil
+-- gathers node keys (primary + secondary if present) via Chainlink API
+-- deploys AuthorizedForwarder per node (both EOAs as senders)
+-- registers forwarders on nodes via API
+-- funds nodes (both primary + secondary EOAs)
+-- calls SetPayees + SetConfig (using forwarder addresses)
+-- creates oracle jobs with dual transmission config
```

### SVR dual transmission flow

Each node has two EOA keys (primary + secondary) and a forwarder contract. The OCR2 job is configured with `enableDualTransmission = true`:

1. **Primary transmission**: Node signs with primary EOA, sends through forwarder to the OCR2Aggregator on-chain (via the regular RPC/Anvil).
2. **Secondary transmission**: Node signs with secondary EOA, sends through forwarder to the same contract but via the flashbots mock endpoint (simulating MEV-Share/OFA).

The flashbots mock is a minimal JSON-RPC server: it responds to `eth_getTransactionCount` with `"0x0"` (so TXMv2 can init nonces) and to `eth_sendRawTransaction` with a fake tx hash. It does **not** submit anything on-chain. Verify dual transmission by checking on-chain `latestAnswer` and `curl http://localhost:9090/health` (response includes `secondaryTxCount`, which should increase over time).

### Node config highlights

- **TXMv2 enabled** with `DualBroadcast = true` and `CustomURL` pointing at the flashbots mock
- **ForwardersEnabled** so that transmissions go through forwarder contracts
- Each node's forwarder authorizes both primary and secondary EOAs as senders

## Swapping a node to a different published version

To test a different Chainlink version on node-4, edit `docker-compose.override.yml`:

```yaml
services:
  node-4:
    image: smartcontract/chainlink:2.37.0
```

Then do a full restart (required because SetConfig on-chain encodes all nodes' keys):

```bash
make restart   # clean + setup-full
```

To revert all nodes to the same version, delete or comment out the override file and `make restart` again.

## Ports

| Service | Host port | Description |
|---------|-----------|-------------|
| Anvil | 8545 | JSON-RPC (HTTP + WS) |
| Postgres | 5433 | PostgreSQL (note: not default 5432) |
| Flashbots mock | 9090 | Mock OFA endpoint (secondary txs) |
| Node 1 | 6688 | Chainlink API |
| Node 1 P2P | 6700 | P2P (only in local-dev mode) |
| Node 2 | 6689 | Chainlink API |
| Node 3 | 6690 | Chainlink API |
| Node 4 | 6691 | Chainlink API (Docker or local) |
| Node 4 P2P | 6694 | P2P (local-dev mode only, on host) |

## API credentials (test-only)

- Email: `local-test@chainlink.test`
- Password: `localdon-testing-only-not-a-real-password`

## Files

```
docker-compose.yml              Main compose file (4 nodes + infra)
docker-compose.override.yml     Optional: override node-4 image for Docker mode
docker-compose.local-dev.yml    Override for local-dev mode (disables node-4, exposes P2P)
Makefile                        Convenience targets
node-config.toml                Shared node config (TXMv2, DualBroadcast, forwarders)
flashbots-mock/
  main.go                       Mock JSON-RPC server (eth_getTransactionCount, eth_sendRawTransaction; health returns secondaryTxCount)
  Dockerfile                    Builds the mock server image
init-db.sh                      Creates 4 postgres databases
local-dev/
  node-config.toml              Node-4 config for host (localhost URLs, P2P announce)
  secrets.toml                  Node-4 secrets for host (localhost DB URL)
secrets/
  node{1-4}-secrets.toml        Per-node DB URLs
  password.txt                  Keystore password
  apicredentials                API login credentials
```

The setup script lives at `core/scripts/setup-localdon/main.go`.

## P2P networking (local-dev mode)

In local-dev mode, node-4 runs on the host and needs to communicate with Docker nodes via P2P:

- **Host -> Docker**: Node-1's P2P port is exposed on host port 6700. Node-4's job spec uses `localhost:6700` as the bootstrapper.
- **Docker -> Host**: Node-4 announces its P2P address as `0.250.250.254:6694` (OrbStack's `host.docker.internal` IP). If not using OrbStack, find your gateway IP with: `docker exec localdon-node-1 getent hosts host.docker.internal` and update `local-dev/node-config.toml`.

## Future improvements

- **DualAggregator contract**: Currently uses the standard `OCR2Aggregator` for both primary and secondary contract addresses. A proper SVR setup would use the `DualAggregator` contract from the `svr-contracts` repo, which handles primary/secondary transmission differentiation at the contract level (with cutoff time, sync iterations, and a secondary proxy). See the `OEV-699-secondary-transmission-integration-test` branch for the contract binding (`dual_aggregator.go`).

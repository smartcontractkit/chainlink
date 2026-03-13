# Local DON (4-node OCR2 Median)

A Docker-based 4-node Chainlink DON running an OCR2 median feed on a local Anvil chain. Useful for local development, testing version upgrades, and validating fixes.

## Prerequisites

- Docker (with compose v2)
- Go 1.23+ (for the setup script and local-dev mode)
- [Foundry](https://book.getfoundry.sh/getting-started/installation) (`cast` CLI, for monitoring)

## Quick start (all Docker)

```bash
# Start all services (postgres, anvil, 4 chainlink nodes)
make up

# Deploy contracts and create oracle jobs
make setup

# Watch node logs
make logs

# Monitor on-chain rounds (address is printed by setup)
cast call <OCR2_ADDRESS> 'latestAnswer()(int256)' --rpc-url http://localhost:8545
```

## Local-dev mode (node-4 from source)

Run node-4 from local source for fast iteration and debugging. The other 3 nodes + infra remain in Docker.

```bash
# Start Docker services (postgres, anvil, 3 nodes; node-4 disabled)
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
| `make setup` | Deploy contracts, create keys, fund nodes, create jobs |
| `make down` | Stop all Docker services |
| `make clean` | Stop all Docker services and remove volumes |
| `make restart` | Full restart: clean + up + setup |
| `make logs` | Follow logs from Docker services |
| `make build` | Build chainlink image from local source (optional) |
| `make local-dev-up` | Start Docker services with node-4 disabled + P2P port exposed |
| `make local-dev-setup` | Run setup script in local-dev mode |
| `make local-dev-node` | Run node-4 locally from source (`go run`) |

## Architecture

```
docker-compose.yml
+-- postgres     (1 instance, 4 databases)
+-- anvil        (local EVM, chain ID 31337, 1s block time)
+-- node-1       (chainlink 2.36.0, acts as P2P bootstrap peer)
+-- node-2       (chainlink 2.36.0)
+-- node-3       (chainlink 2.36.0)
+-- node-4       (chainlink 2.36.0, or local source in local-dev mode)

setup script (Go, runs on host)
+-- deploys LINK token + OCR2Aggregator to Anvil
+-- gathers node keys via Chainlink API
+-- funds nodes, calls SetPayees + SetConfig
+-- creates oracle jobs on all 4 nodes
```

Each node observes pseudo-random prices (5 hardcoded values, randomly selected each round via the `any` pipeline task) and the DON produces a median on-chain.

## Swapping a node to a different published version

To test a different Chainlink version on node-4, edit `docker-compose.override.yml`:

```yaml
services:
  node-4:
    image: smartcontract/chainlink:2.37.0
```

Then do a full restart (required because SetConfig on-chain encodes all nodes' keys):

```bash
make restart
```

To revert all nodes to the same version, delete or comment out the override file and `make restart` again.

## Ports

| Service | Host port | Description |
|---------|-----------|-------------|
| Anvil | 8545 | JSON-RPC (HTTP + WS) |
| Postgres | 5433 | PostgreSQL (note: not default 5432) |
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
docker-compose.yml              Main compose file (4 nodes)
docker-compose.override.yml     Optional: override node-4 image for Docker mode
docker-compose.local-dev.yml    Override for local-dev mode (disables node-4, exposes P2P)
Makefile                        Convenience targets
node-config.toml                Shared node config for Docker nodes
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

- **Host → Docker**: Node-1's P2P port is exposed on host port 6700. Node-4's job spec uses `localhost:6700` as the bootstrapper.
- **Docker → Host**: Node-4 announces its P2P address as `0.250.250.254:6694` (OrbStack's `host.docker.internal` IP). If not using OrbStack, find your gateway IP with: `docker exec localdon-node-1 getent hosts host.docker.internal` and update `local-dev/node-config.toml`.

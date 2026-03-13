# Local DON (4-node OCR2 Median)

A Docker-based 4-node Chainlink DON running an OCR2 median feed on a local Anvil chain. Useful for local development, testing version upgrades, and validating fixes.

## Prerequisites

- Docker (with compose v2)
- Go 1.23+ (for the setup script)
- [Foundry](https://book.getfoundry.sh/getting-started/installation) (`cast` CLI, for monitoring)

## Quick start

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

## Makefile targets

| Target | Description |
|--------|-------------|
| `make up` | Start all services, wait for health checks |
| `make setup` | Deploy contracts, create keys, fund nodes, create jobs |
| `make down` | Stop all services |
| `make clean` | Stop all services and remove volumes |
| `make restart` | Full restart: clean + up + setup |
| `make logs` | Follow logs from all services |
| `make build` | Build chainlink image from local source (optional) |

## Architecture

```
docker-compose.yml
+-- postgres     (1 instance, 4 databases)
+-- anvil        (local EVM, chain ID 31337, 1s block time)
+-- node-1       (chainlink 2.36.0, acts as P2P bootstrap peer)
+-- node-2       (chainlink 2.36.0)
+-- node-3       (chainlink 2.36.0)
+-- node-4       (chainlink 2.36.0, or overridden version)

setup script (Go, runs on host)
+-- deploys LINK token + OCR2Aggregator to Anvil
+-- gathers node keys via Chainlink API
+-- funds nodes, calls SetPayees + SetConfig
+-- creates oracle jobs on all 4 nodes
```

Each node observes pseudo-random prices (5 hardcoded values, randomly selected each round via the `any` pipeline task) and the DON produces a median on-chain.

## Swapping a node to a different version

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

This takes about 60 seconds. To revert all nodes to the same version, delete or comment out the override file and `make restart` again.

## Ports

| Service | Host port | Description |
|---------|-----------|-------------|
| Anvil | 8545 | JSON-RPC (HTTP + WS) |
| Postgres | 5433 | PostgreSQL (note: not default 5432) |
| Node 1 | 6688 | Chainlink API |
| Node 2 | 6689 | Chainlink API |
| Node 3 | 6690 | Chainlink API |
| Node 4 | 6691 | Chainlink API |

## API credentials (test-only)

- Email: `local-test@chainlink.test`
- Password: `localdon-testing-only-not-a-real-password`

## Files

```
docker-compose.yml              Main compose file
docker-compose.override.yml     Optional: override node-4 image
Makefile                        Convenience targets
node-config.toml                Shared node config (OCR2, P2P, chain 31337)
init-db.sh                      Creates 4 postgres databases
secrets/
  node{1-4}-secrets.toml        Per-node DB URLs
  password.txt                  Keystore password
  apicredentials                API login credentials
```

The setup script lives at `core/scripts/setup-localdon/main.go`.

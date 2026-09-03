# LinkSwarm

**Chainlink Functions-gated high-throughput settlement orchestration**

LinkSwarm is a production-shaped settlement engine designed for enterprise payroll and machine-to-machine payments.  
Nothing is released until **Chainlink Functions** has verified the batch commitment.

---

## Table of Contents

- [Overview](#overview)
- [Why LinkSwarm](#why-linkswarm)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Functions Integration](#functions-integration)
- [Reports & Audit Trail](#reports--audit-trail)
- [Security Model](#security-model)
- [Project Structure](#project-structure)
- [Development](#development)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

---

## Overview

LinkSwarm treats Chainlink Functions as the mandatory verification gate for every settlement batch.

1. Deterministic payloads are built for every worker.
2. An Entropy-3.12 commitment root is computed over the entire batch.
3. The commitment (plus batch metadata) is submitted to Chainlink Functions.
4. Only after a successful, verified Functions response does delivery begin.
5. Each worker is paid according to their preferred network and payment method.
6. An immutable private ledger entry is written that includes the Functions request ID and commitment.

This gives Chainlink a concrete, high-volume use case while giving enterprises a clear audit trail and worker autonomy.

---

## Why LinkSwarm

| Capability | Description |
|------------|-------------|
| Functions as gate | No settlement occurs until Functions verifies the commitment |
| High throughput | Designed for 10 000+ workers per batch with controlled concurrency |
| Private ledgers | Per-user immutable history with tax, fee, and metadata |
| Card lock | Instant transaction blocking when a wallet is locked |
| Multi-network delivery | EVM, Solana, Cosmos, UTXO after verification |
| USDC + fiat bridge | Automatic conversion and dual-hash (USDC + FIAT) accounting |
| Idempotent requests | Every Functions call carries an idempotency key |
| Simulation + Live | Fully offline simulation mode and live HTTP client for real gateways |
| Full reports | JSON + CSV with `chainlinkVerified` and `functionsRequestId` on every row |

---

## Architecture

Worker Batch
│
▼
Entropy-3.12 Commitment
│
▼
Chainlink Functions  ←── mandatory verification gate
│
▼ (only on success)
Delivery Adapters (EVM / Solana / Cosmos / UTXO)
│
▼
Private Ledgers + Card-Lock checks + USDC/Fiat bridge
│
▼
JSON / CSV Reports

text

**Primary trust root**: Chainlink Functions  
**Secondary**: delivery networks and bridges

---

## Quick Start

```bash
git clone https://github.com/Pray4Love1/linkswarm.git
cd linkswarm
go mod tidy
go run engine.go

Default mode is simulation. It processes 10 000 workers, calls the Functions simulation client, and writes:





linkswarm_report.json



linkswarm_report.csv

Live mode

Bash

export LINKSWARM_MODE=live
export CHAINLINK_FUNCTIONS_ROUTER=https://your-functions-gateway.example.com
export CHAINLINK_SUBSCRIPTION_ID=123
export CHAINLINK_DON_ID=fun-ethereum-sepolia-1
go run engine.go

Configuration

All configuration is via environment variables:



| Variable | Default | Description |
| --- | --- | --- |
| LINKSWARM_MODE | simulation | simulation or live |
| CHAINLINK_FUNCTIONS_ROUTER | — | Live Functions HTTP endpoint |
| CHAINLINK_SUBSCRIPTION_ID | 0 | Functions subscription ID |
| CHAINLINK_DON_ID | fun-ethereum-sepolia-1 | DON identifier |
| WORKER_COUNT | 10000 | Number of workers in the e2e batch |
| MAX_CONCURRENCY | 512 | Max parallel settlement goroutines |
| REQUEST_TIMEOUT_SEC | 60 | Functions request timeout |
| RETRY_ATTEMPTS | 3 | Functions retry count |
| CONVERSION_FEE_BPS | 10 | USDC conversion fee (basis points) |
| ENABLE_CARD_LOCK | true | Enforce card-lock checks |
| REPORT_DIR | . | Directory for report output |

Functions Integration

Simulation client

Fully offline. Produces deterministic request IDs and verified responses so the entire e2e path can be exercised without credentials.

Live client

HTTP client that posts to:

text

POST {CHAINLINK_FUNCTIONS_ROUTER}/functions/request

Payload includes source, args, subscriptionId, donId, gasLimit, and idempotencyKey.

Source code executed inside Functions

The default source validates the commitment and returns a structured verification payload. It can be replaced with richer business logic as needed.

Reports & Audit Trail

Every settlement row contains:





chainlinkVerified



functionsRequestId



network, payment method



USDC and fiat hashes (when applicable)



full timestamp and epoch



private-ledger metadata linking back to the commitment and request ID

This satisfies enterprise audit and regulatory requirements while remaining portable (JSON + CSV).

Security Model





Functions gate: no funds move until the commitment is verified.



Idempotency keys: prevent duplicate settlement on retries.



Card lock: instantaneous blocking of individual wallets.



Private ledgers: append-only, per-user, with running balances.



Entropy-3.12: deterministic yet non-replayable batch commitments.



Offline-friendly construction: private keys never need to touch the network for the commitment step.

Project Structure

text

linkswarm/
├── engine.go          # Complete production engine + e2e runner
├── go.mod
├── LICENSE            # Apache License 2.0
├── NOTICE
├── README.md
├── CONTRIBUTING.md
├── SECURITY.md
└── ROADMAP.md

Development

Bash

go mod tidy
go run engine.go
go build -o linkswarm engine.go
./linkswarm

Roadmap

See ROADMAP.md.

Contributing

See CONTRIBUTING.md.

License

Apache License 2.0 — see LICENSE.

Contact





Author: Jon Sanders (@Pray4Love1)



Repository: https://github.com/Pray4Love1/linkswarm
# LinksSwarm

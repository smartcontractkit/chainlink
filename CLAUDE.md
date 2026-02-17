# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build and install chainlink binary
make install

# Build dev version (disables TLS requirements)
make chainlink-dev

# Run all tests (requires postgres - see Database Setup below)
go test ./...

# Run tests without database dependency
go test -short ./...

# Run a single test
go test -run TestFunctionName ./path/to/package

# Run tests with race detection
GORACE="log_path=$PWD/race" go test -race ./core/path/to/pkg -count 10

# Generate code (mocks, protobuf, etc.)
make generate

# Tidy all Go modules
make gomodtidy

# Run linter
make golangci-lint

# Fix lint issues automatically
make lint-fix

# Install mockery for generating mocks
make mockery
```

## Database Setup

Tests require PostgreSQL. Set `CL_DATABASE_URL` environment variable:
```bash
make setup-testdb      # One-time setup, creates .dbenv file
source .dbenv          # Load the env var
make testdb            # Prepare test database (run after migrations)
```

## Project Structure

This is the Chainlink oracle node - a Go monorepo with multiple modules:

```
chainlink/
├── core/              # Main node application
│   ├── services/      # Long-running services (OCR, VRF, Keeper, etc.)
│   ├── capabilities/  # Composable capability system
│   ├── cmd/           # CLI commands
│   ├── web/           # REST API
│   └── config/        # Configuration management
├── plugins/           # LOOP (Local Out-Of-Process) plugins
├── deployment/        # Deployment infrastructure (changesets, environments)
├── integration-tests/ # End-to-end tests
├── ccip/              # Cross-chain interoperability protocol configs
└── common/            # Shared chain abstractions
```

## Go Modules

Three main modules with relative replaces:
- `github.com/smartcontractkit/chainlink/v2` (root)
- `github.com/smartcontractkit/chainlink/integration-tests`
- `github.com/smartcontractkit/chainlink/core/scripts`

Use `make gomodtidy` to tidy all modules together.

## Key Architecture Patterns

### Services
All long-running components implement `services.Service`:
```go
type Service interface {
    Start(ctx context.Context) error
    Close() error
}
```

### Job System
- Jobs are managed by a Spawner with type-specific Delegates
- Job types: OCR, OCR2, DirectRequest, Keeper, VRF, Cron, Webhook, etc.
- Jobs use the Pipeline engine for task execution

### Capabilities
Modern composable framework in `core/capabilities/`:
- Registry-based capability discovery
- Supports triggers, actions, consensus, targets
- DON (Decentralized Oracle Network) aware

### Relayers
Chain-agnostic abstraction in `core/services/relay/`:
- Supports EVM, Cosmos, Solana, Aptos, etc.
- Plugin provider pattern for extensibility

### Deployment
Code-driven deployment in `deployment/`:
- **Changesets**: Functions describing infrastructure changes
- **Environment**: Represents existing on-chain/off-chain state
- **Address Book**: Version-controlled on-chain addresses

## Code Style

### Imports
Local imports grouped separately (enforced by goimports):
```go
import (
    "context"

    "github.com/external/pkg"

    "github.com/smartcontractkit/chainlink/v2/core/..."
)
```

### Logging
Use structured logging with key-value pairs:
```go
lggr.Infow("Message", "key1", value1, "key2", value2)
```

### Comments
Do not add comments that simply describe what the code is doing. Code should be self-documenting through clear naming and structure. Comments should explain *why* something is done, not *what* is being done. Only add comments when:
- Explaining non-obvious business logic or requirements
- Documenting complex algorithms or performance considerations
- Clarifying workarounds for external bugs or limitations
- Providing context that cannot be expressed in code

### Dependencies to Avoid
- `github.com/gofrs/uuid` - use `github.com/google/uuid`
- `go.uber.org/multierr` - use `errors.Join` from stdlib
- `github.com/go-gorm/gorm` - use `github.com/jmoiron/sqlx` directly

## Git Commits

When creating commits:
- Do NOT include `Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>` tags
- Use clear, concise commit messages that describe the change

## Testing

### Test Flags
- `-short`: Skip database-dependent tests
- `-p 1`: Limit parallel packages (reduces RAM usage)
- `-race`: Enable race detector
- `-count N`: Run tests N times

### Integration Tests
Located in `integration-tests/`. Require Docker and additional setup - see `integration-tests/README.md`.

## Changesets

PRs that modify code should include a changeset:
```bash
pnpm install    # First time only
pnpm changeset  # Create changeset entry
```

## Vault

The Vault system provides secure, decentralized secrets management for workflows. Secrets are encrypted using threshold cryptography (TDH2) and require consensus from a DON (Decentralized Oracle Network) for all operations.

### Architecture

The Vault system spans three layers:

```
core/capabilities/vault/       # Capability layer - exposes secrets as workflow capability
core/services/ocr2/plugins/vault/  # Consensus layer - OCR3 plugin for DON consensus
core/services/gateway/handlers/vault/  # Gateway layer - routes external requests to DON
```

### Components

**Capability** (`core/capabilities/vault/capability.go`):
- Implements `capabilities.ExecutableCapability` interface
- Exposes `vault.secrets.*` methods to workflows
- Validates and authorizes requests before forwarding to OCR handler
- Methods: `CreateSecrets`, `UpdateSecrets`, `DeleteSecrets`, `GetSecrets`, `ListSecretIdentifiers`, `GetPublicKey`

**Gateway Handler** (`core/capabilities/vault/gw_handler.go`):
- Implements `connector.GatewayConnectorHandler`
- Handles JSON-RPC requests from gateway connector
- Routes to the capability's secrets service methods

**Request Authorizer** (`core/capabilities/vault/request_authorizer.go`):
- Validates request signatures against workflow registry allowlist
- Ensures users can only access their own secrets (owner-based isolation)

**OCR2 Plugin** (`core/services/ocr2/plugins/vault/plugin.go`):
- Implements `ocr3_1types.ReportingPlugin` for consensus
- Uses TDH2 threshold decryption - secrets encrypted with DON public key
- Stores secrets in OCR3's KeyValue state (via `KVStore`)
- Batches requests for efficient consensus rounds

**Gateway Handler** (`core/services/gateway/handlers/vault/handler.go`):
- Implements `gwhandlers.Handler` for the gateway service
- Fans out requests to all Vault DON nodes
- Aggregates responses and waits for quorum
- Caches public key responses

### Data Model

Secrets are identified by: `owner::namespace::key`
- **Owner**: Address of the secret owner (enforced by authorization)
- **Namespace**: Logical grouping (defaults to "main")
- **Key**: Secret identifier within namespace

### Methods

| Method | Description |
|--------|-------------|
| `vault.secrets.create` | Create new encrypted secrets |
| `vault.secrets.update` | Update existing secrets |
| `vault.secrets.delete` | Delete secrets |
| `vault.secrets.list` | List secret identifiers for an owner |
| `vault.secrets.get` | Retrieve secrets (internal/dev only) |
| `vault.publicKey.get` | Get DON's TDH2 public key for encryption |

### Request Flow

1. External client encrypts secrets with DON's TDH2 public key
2. Request sent to Gateway with signed authorization
3. Gateway handler validates and fans out to DON nodes
4. Each node processes via capability → OCR plugin
5. OCR3 consensus produces signed report
6. Gateway aggregates responses, returns on quorum

### OCR2 Plugin Details (`core/services/ocr2/plugins/vault/`)

The plugin implements the OCR3 protocol for decentralized consensus on secrets operations.

**Files:**

| File | Purpose |
|------|---------|
| `plugin.go` | Main `ReportingPlugin` implementation with OCR3 lifecycle methods |
| `config.go` | Plugin configuration (request expiry, DKG settings) |
| `kvstore.go` | Key-value store abstraction over OCR3's `KeyValueStateReader/Writer` |
| `orm.go` | Database ORM for DKG result packages (stores threshold keys) |
| `transmitter.go` | Transmits signed reports back to the request handler |
| `disk_monitor.go` | Monitors vault storage directory size for metrics |
| `tdh2_helpers.go` | Helpers for converting between TDH2 key formats |

**OCR3 Protocol Methods in `plugin.go`:**

| Method | Purpose |
|--------|---------|
| `Query` | Returns empty (no leader-driven queries) |
| `Observation` | Each node observes pending requests, generates per-request observations |
| `ValidateObservation` | Validates observation structure and blob availability |
| `ObservationQuorum` | Returns true when N-F observations received |
| `StateTransition` | Aggregates observations, writes to KV store, produces outcomes |
| `Reports` | Converts outcomes to signed reports (JSON or Protobuf) |
| `ShouldAcceptAttestedReport` | Always returns true |
| `ShouldTransmitAcceptedReport` | Always returns true |

**Consensus Thresholds:**
- **GetSecrets**: Requires 2F+1 observations (needs F+1 decryption shares)
- **Create/Update/Delete/List**: Requires F+1 observations (at least 1 honest node)

**Key Data Structures:**

```
// KVStore keys
"Key::{owner}::{namespace}::{key}"      // Encrypted secret data
"Metadata::{owner}"                      // List of secret identifiers per owner
"PendingQueue::Index"                    // Pending queue length
"PendingQueue::Item::{n}"               // Individual pending items
```

**Deterministic Pending Queue:**
When `EnableDeterministicPendingQueue` is true:
- Pending requests stored in KV store for cross-node consistency
- Nodes broadcast local queue items as blobs
- Random nonces aggregated for deterministic sorting
- Ensures all honest nodes process same requests in same order

**TDH2 Threshold Decryption Flow (GetSecrets):**
1. Each node retrieves encrypted secret from KV store
2. Node generates decryption share using its private key share
3. Share encrypted to each workflow DON member's encryption key
4. Shares aggregated across 2F+1 nodes in StateTransition
5. Workflow DON can reconstruct secret with F+1 shares

### Notable ReportingPlugin Method

**`Observation`**: Fetches a batch of pending requests and generates observations for each.

- *GetSecrets*: Validates secret identifier (owner, namespace, key format and length limits). Checks secret exists in KV store. Unmarshals and verifies ciphertext against DON public key. Verifies ciphertext has correct label (owner address). Generates decryption share using node's private key share. For each passed-in encryption key (workflow DON members), encrypts the decryption share using `box.SealAnonymous`.

- *CreateSecrets*: Validates each secret identifier. Checks for duplicate IDs within the batch. Decodes and validates ciphertext (hex format, size limit). Verifies ciphertext against DON public key. Does NOT check if key exists (deferred to StateTransition for batch ordering).

- *UpdateSecrets*: Same validation as CreateSecrets. Does NOT check if key exists (deferred to StateTransition).

- *DeleteSecrets*: Validates secret identifier. Checks for duplicate IDs in batch. Verifies secret exists in KV store.

- *ListSecretIdentifiers*: Validates owner is non-empty. Reads metadata for owner. Returns sorted list, optionally filtered by namespace.

- *Deterministic Queue Mode*: Also observes local pending queue items not in shared queue. Broadcasts them as blobs. Generates random 32-byte nonce for deterministic sorting.

**`ValidateObservation`**: Validates observation from another node.

- Unmarshals and validates each observation structure
- Checks no duplicate observation IDs
- Validates request/response pair consistency (same number of items)
- Checks for duplicate secret IDs within batch requests
- *Deterministic Queue Mode*: Verifies observation count matches pending queue length. Verifies all pending queue IDs have observations. Fetches all blob handles to ensure availability.

**`StateTransition`**: Aggregates observations and writes state changes. Two phases:

*Phase 1 - Process pending requests*: Groups observations by request ID, then by SHA hash.
- *GetSecrets*: Requires 2F+1 matching observations. Aggregates decryption shares across observations. Merges encrypted shares per encryption key.
- *CreateSecrets*: Requires F+1 matching observations. Checks key doesn't already exist. Checks owner hasn't exceeded max secrets limit. Writes to KV store.
- *UpdateSecrets*: Requires F+1 matching observations. Checks key exists. Overwrites in KV store.
- *DeleteSecrets*: Requires F+1 matching observations. Deletes from KV store and owner metadata.
- *ListSecretIdentifiers*: Requires F+1 matching observations. Takes first matching response.

*Phase 2 - Process pending queue (deterministic mode only)*: Collects items from all observations via blob fetches. Groups by ID then SHA. Keeps items with F+1 consensus. Concatenates all nonces for shared sort key. Sorts by `SHA256(id || aggregated_nonce)`. Truncates to max batch size. Writes new pending queue to KV store.

**`Reports`**: Converts outcomes to signed reports.

- *GetSecrets*: Marshals as Protobuf (contains encrypted shares)
- *Create/Update/Delete/List*: Marshals as canonical JSON
- All reports include report info (ID, request type, format) and key bundle name ("evm")

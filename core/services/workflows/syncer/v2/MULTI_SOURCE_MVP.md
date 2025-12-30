# Multi-Source Workflow Registry MVP

This document describes the MVP implementation for reading workflow metadata from multiple sources (contract + file-based).

## Overview

The workflow registry syncer now supports multiple sources of workflow metadata:

1. **ContractWorkflowSource** (primary): Reads from the on-chain workflow registry contract
2. **FileWorkflowSource** (supplementary): Reads from a local JSON file

Both sources are aggregated by `MultiSourceWorkflowAggregator` and workflows from all sources are reconciled together.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    WorkflowRegistry Syncer                       │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │            MultiSourceWorkflowAggregator                  │   │
│  │                                                           │   │
│  │   ┌─────────────────────┐  ┌─────────────────────┐       │   │
│  │   │ ContractWorkflow    │  │ FileWorkflow        │       │   │
│  │   │ Source              │  │ Source              │       │   │
│  │   │                     │  │                     │       │   │
│  │   │ (on-chain contract) │  │ (/tmp/workflows_   │       │   │
│  │   │                     │  │  metadata.json)     │       │   │
│  │   └─────────────────────┘  └─────────────────────┘       │   │
│  └──────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│                    []WorkflowMetadataView                        │
│                              │                                   │
│                              ▼                                   │
│               generateReconciliationEvents()                     │
│                              │                                   │
│                              ▼                                   │
│                      Event Handler                               │
│                              │                                   │
│                              ▼                                   │
│                     Engine Registry                              │
└─────────────────────────────────────────────────────────────────┘
```

## File Source Format

The file source reads from `/tmp/workflows_metadata.json` (hardcoded for MVP).

### JSON Schema

```json
{
  "workflows": [
    {
      "workflow_id": "<32-byte hex string without 0x prefix>",
      "owner": "<owner address hex without 0x prefix>",
      "created_at": <unix timestamp>,
      "status": <0=active, 1=paused>,
      "workflow_name": "<name>",
      "binary_url": "<URL to fetch binary - same format as contract>",
      "config_url": "<URL to fetch config - same format as contract>",
      "tag": "<version tag>",
      "attributes": "<optional JSON string>",
      "don_family": "<DON family name>"
    }
  ]
}
```

### Example

```json
{
  "workflows": [
    {
      "workflow_id": "0102030405060708091011121314151617181920212223242526272829303132",
      "owner": "f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "created_at": 1733250000,
      "status": 0,
      "workflow_name": "my-file-workflow",
      "binary_url": "file:///home/chainlink/workflows/my_workflow.wasm",
      "config_url": "file:///home/chainlink/workflows/my_config.json",
      "tag": "v1.0.0",
      "don_family": "workflow"
    }
  ]
}
```

## Testing with Local CRE

### Prerequisites

1. Local CRE environment set up (see `core/scripts/cre/environment/README.md`)
2. Docker running
3. Go toolchain installed

### Helper Tool: generate_file_source

A helper tool is provided to generate the workflow metadata JSON with correct workflowID:

```bash
cd core/scripts/cre/environment
go run ./cmd/generate_file_source \
  --binary /path/to/workflow.wasm \
  --config /path/to/config.json \
  --name my-workflow \
  --owner f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
  --output /tmp/workflows_metadata.json \
  --don-family workflow
```

### Test Scenario 1: Contract-Only Workflow

This verifies existing functionality still works.

```bash
# Start the environment
cd core/scripts/cre/environment
go run . env start --auto-setup

# Deploy a workflow via contract
go run . workflow deploy -w ./examples/workflows/v2/cron/main.go -n cron_example

# Verify workflow is running (check logs or trigger if using http-trigger)
```

### Test Scenario 2: File-Source Workflow (Complete Walkthrough)

This tests the new file-based workflow source with an existing workflow.

```bash
# 1. Start the environment
cd core/scripts/cre/environment
go run . env start --auto-setup

# 2. Deploy a workflow via contract first (this creates the binary in containers)
go run . workflow deploy -w ./examples/workflows/v2/cron/main.go -n cron_contract

# 3. Find the compiled workflow binary (created during deploy)
# The binary will be in /home/chainlink/workflows/ in the container

# 4. Get the existing workflow binary from a container
docker cp workflow-node1:/home/chainlink/workflows/cron_contract.wasm /tmp/cron_contract.wasm

# 5. Generate the file source metadata with a DIFFERENT workflow name
go run ./cmd/generate_file_source \
  --binary /tmp/cron_contract.wasm \
  --name file_source_cron \
  --owner f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
  --output /tmp/workflows_metadata.json \
  --don-family workflow \
  --binary-url-prefix "file:///home/chainlink/workflows/" \
  --config-url-prefix "file:///home/chainlink/workflows/"

# 6. Copy the binary to containers with new name
docker cp /tmp/cron_contract.wasm workflow-node1:/home/chainlink/workflows/file_source_workflow.wasm
docker cp /tmp/cron_contract.wasm workflow-node2:/home/chainlink/workflows/file_source_workflow.wasm
docker cp /tmp/cron_contract.wasm workflow-node3:/home/chainlink/workflows/file_source_workflow.wasm
docker cp /tmp/cron_contract.wasm workflow-node4:/home/chainlink/workflows/file_source_workflow.wasm
docker cp /tmp/cron_contract.wasm workflow-node5:/home/chainlink/workflows/file_source_workflow.wasm

# 7. Create an empty config file
echo '{}' > /tmp/file_source_config.json
docker cp /tmp/file_source_config.json workflow-node1:/home/chainlink/workflows/file_source_config.json
docker cp /tmp/file_source_config.json workflow-node2:/home/chainlink/workflows/file_source_config.json
docker cp /tmp/file_source_config.json workflow-node3:/home/chainlink/workflows/file_source_config.json
docker cp /tmp/file_source_config.json workflow-node4:/home/chainlink/workflows/file_source_config.json
docker cp /tmp/file_source_config.json workflow-node5:/home/chainlink/workflows/file_source_config.json

# 8. Copy the metadata file to all nodes
docker cp /tmp/workflows_metadata.json workflow-node1:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata.json workflow-node2:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata.json workflow-node3:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata.json workflow-node4:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata.json workflow-node5:/tmp/workflows_metadata.json

# 9. Wait for the syncer to pick up the workflow (default 12 second interval)
# Check logs for "Loaded workflows from file" messages
docker logs workflow-node1 2>&1 | grep -i "file"

# 10. Verify both workflows are running (contract and file source)
docker logs workflow-node1 2>&1 | grep -i "workflow engine"
```

### Test Scenario 3: Mixed Sources

Test both contract and file sources together.

```bash
# 1. Deploy workflow via contract
go run . workflow deploy -w ./examples/workflows/v2/cron/main.go -n contract_workflow

# 2. Add a different workflow via file source (follow steps 3-8 from Scenario 2)

# 3. Verify both workflows are running
# You should see two workflow engines running
docker logs workflow-node1 2>&1 | grep -i "Aggregated workflows from all sources"
# Should show totalWorkflows: 2
```

### Test Scenario 4: Pause/Delete from File Source

```bash
# 1. Start with both contract and file-source workflows running (as above)

# 2. Pause the file-source workflow by changing status to 1
cat > /tmp/workflows_metadata_paused.json << 'EOF'
{
  "workflows": [
    {
      "workflow_id": "<YOUR_WORKFLOW_ID>",
      "owner": "f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "status": 1,
      "workflow_name": "file_source_cron",
      "binary_url": "file:///home/chainlink/workflows/file_source_workflow.wasm",
      "config_url": "file:///home/chainlink/workflows/file_source_config.json",
      "don_family": "workflow"
    }
  ]
}
EOF

# Copy to all nodes
docker cp /tmp/workflows_metadata_paused.json workflow-node1:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata_paused.json workflow-node2:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata_paused.json workflow-node3:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata_paused.json workflow-node4:/tmp/workflows_metadata.json
docker cp /tmp/workflows_metadata_paused.json workflow-node5:/tmp/workflows_metadata.json

# 3. Wait for syncer to detect the change and check logs
docker logs workflow-node1 2>&1 | grep -i "paused"

# 4. Delete by removing from file
echo '{"workflows":[]}' > /tmp/empty_metadata.json
docker cp /tmp/empty_metadata.json workflow-node1:/tmp/workflows_metadata.json
docker cp /tmp/empty_metadata.json workflow-node2:/tmp/workflows_metadata.json
docker cp /tmp/empty_metadata.json workflow-node3:/tmp/workflows_metadata.json
docker cp /tmp/empty_metadata.json workflow-node4:/tmp/workflows_metadata.json
docker cp /tmp/empty_metadata.json workflow-node5:/tmp/workflows_metadata.json

# 5. Contract workflow should still be running, file-source workflow should be removed
```

### Verifying Multi-Source Works

Check the logs for these messages:

```bash
# See aggregation from multiple sources
docker logs workflow-node1 2>&1 | grep "Aggregated workflows from all sources"

# See file source loading
docker logs workflow-node1 2>&1 | grep "Loaded workflows from file"

# See contract source loading
docker logs workflow-node1 2>&1 | grep "ContractWorkflowSource"
```

## Key Behaviors

### Source Aggregation
- Workflows from all sources are merged into a single list
- The contract source's blockchain head is used for reconciliation
- If one source fails, others continue to work (graceful degradation)

### Workflow ID Collisions
- **MVP Assumption**: WorkflowID collisions are handled externally
- If the same workflowID appears in multiple sources, both entries will be present
- This may cause issues - discovery of such edge cases is a goal of this MVP

### File Source Characteristics
- File is read on every sync interval (default 12 seconds)
- Missing file = empty workflow list (not an error)
- Invalid JSON entries are skipped with a warning
- File source is always "ready" (unlike contract source which needs initialization)

## Implementation Files

| File | Description |
|------|-------------|
| `types.go` | `WorkflowMetadataSource` interface definition |
| `file_workflow_source.go` | File-based source implementation |
| `contract_workflow_source.go` | Contract-based source implementation |
| `multi_source.go` | Aggregator that combines multiple sources |
| `workflow_registry.go` | Updated to use multi-source aggregator |
| `file_workflow_source_test.go` | Unit tests for file source |
| `multi_source_test.go` | Unit tests for aggregator |

## Known Limitations (MVP)

1. **Hardcoded file path**: `/tmp/workflows_metadata.json` is not configurable
2. **No atomic updates**: File changes may be read partially if written during sync
3. **No persistence**: File must be created manually on each node
4. **No validation**: WorkflowID hash is not verified against artifacts
5. **Same DON family**: All workflows in file must match one of the DON's families

## Future Improvements

1. Configurable file path via TOML
2. S3/HTTP-based source implementations
3. WorkflowID collision detection and resolution
4. Source provenance tracking in engine registry
5. File watch for instant updates (instead of polling)
6. Kubernetes ConfigMap/Secret support for CRIB deployments

## Debugging

### Check if file source is being read

Look for these log messages:
- `"Loaded workflows from file"` - File was successfully read
- `"Workflow metadata file does not exist"` - File doesn't exist (normal if not using file source)
- `"Source not ready, skipping"` - Contract source not yet initialized

### Check aggregated workflows

Look for:
- `"Aggregated workflows from all sources"` with `totalWorkflows` count
- `"fetching workflow metadata from all sources"` - Sync is running

### Verify workflow engine started

Look for:
- `"Creating Workflow Engine for workflow spec"` 
- Check the engine registry in metrics


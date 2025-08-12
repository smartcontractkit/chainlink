# Tron PoR Test

This directory contains a Tron-based Proof of Reserve (PoR) test that validates the EVM relayer refactor for Tron support.

## Overview

The test (`por_tron_test.go`) is based on the existing PoR test but adapted to use Tron instead of Anvil. It validates that:

1. **EVM Relayer works with Tron**: The refactored EVM relayer can handle Tron chains with TXM injection
2. **E2E CRE flow works**: The complete workflow including write target capability functions correctly
3. **Tron simulated backend integration**: Uses a Tron simulated backend instead of Ethereum-based ones

## Files

- `por_tron_test.go` - Main test file with Tron-specific PoR test
- `tron_backend.go` - Tron simulated backend management
- `environment-one-don-tron.toml` - Tron-specific configuration

## Prerequisites

- Docker must be running
- Go 1.21+ 
- The project dependencies must be installed

## Running the Test

### Simple Go Test

You can run the test directly like any other Go test:

```bash
cd system-tests/tests/smoke/cre

# Run the Tron PoR test
go test -v -run TestCRE_OCR3_PoR_Workflow_SingleDon_Tron_MockedPrice -timeout 30m
```

### With Environment Variables

For more control, you can set environment variables:

```bash
cd system-tests/tests/smoke/cre

export CTF_CONFIGS="environment-one-don-tron.toml"
export PRIVATE_KEY="da146374a75310b9666e834ee4ad0866d6f4035967bfc76217c5a495fff9f0d0"

go test -v -run TestCRE_OCR3_PoR_Workflow_SingleDon_Tron_MockedPrice -timeout 30m
```

## Test Architecture

### Tron Backend (`tron_backend.go`)
- Manages Tron simulated backend lifecycle
- Creates Docker containers for Tron node and postgres
- Provides JSON-RPC endpoints for the EVM relayer
- Handles cleanup automatically

### Test Flow
1. **Setup**: Creates Tron backend and starts containers
2. **Configuration**: Modifies blockchain config to use Tron endpoints  
3. **Relayer Setup**: Uses existing PoR test infrastructure with Tron configuration
4. **Execution**: Runs workflow and validates write target capability
5. **Cleanup**: Automatically stops and removes containers

### Key Differences from Anvil Test
- Uses Tron chain ID (728126428) instead of Ethereum chain IDs
- Connects to Tron JSON-RPC endpoints instead of Anvil
- Uses Tron-compatible private key
- Leverages EVM relayer's Tron TXM injection for write operations

## Configuration

The test uses `environment-one-don-tron.toml` which specifies:
- Tron chain ID and type
- Single DON configuration  
- Same workflow configurations as original test
- Docker infrastructure setup

## Validation

The test validates:
- ✅ Tron backend starts successfully
- ✅ EVM relayer connects to Tron endpoints
- ✅ Chainlink nodes can read from Tron
- ✅ Write target capability works with Tron TXM
- ✅ Complete E2E workflow executes successfully
- ✅ Feed updates are written to Tron contracts

## Troubleshooting

### Docker Issues
- Ensure Docker is running: `docker info`
- Check for port conflicts on ports 16671, 16672, 5432
- Clean up existing containers: `docker stop tron-node tron-postgres; docker rm tron-node tron-postgres`

### Test Failures
- Check Docker logs for Tron node: `docker logs tron-node`
- Verify Tron endpoints are accessible: `curl http://localhost:16671/jsonrpc`
- Review test logs for detailed error information

### Configuration Issues
- Ensure `environment-one-don-tron.toml` exists and is valid
- Verify private key is set correctly
- Check that all required dependencies are available 
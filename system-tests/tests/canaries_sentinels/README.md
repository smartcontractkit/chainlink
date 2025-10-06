# Canaries Sentinels

## Introduction

The name "Canary testing" comes from the old practice of using canaries in coal mines to detect dangerous gases early—when the canary stopped singing or died, miners knew to evacuate immediately. In our context, we call these recurring workflows "Canaries Sentinels" because they continuously monitor system health by running at regular intervals. When they fail repeatedly (e.g., 3 failures within 5 minutes), it serves as an early warning system that something critical is wrong with our infrastructure—just like the canary warning miners of danger, but with the added vigilance of a sentinel standing guard.

These workflows serve as continuous health checks for various Chainlink capabilities, providing early detection of issues in production-like environments before they impact critical systems.

## Coverage Matrix

This matrix shows which capabilities are covered by each canary test, making it easy to identify gaps in coverage.

| Capability | Proof of Reserve (Cron) | Future Test 1 | Future Test 2 | Coverage Status |
|------------|-------------------------|---------------|---------------|-----------------|
| **Triggers** |  |  |  |  |
| Cron Scheduler | ✅ | | | Covered |
| **Chain Capabilities** |  |  |  |  |
| EVM Balance Reading | ✅ | | | Covered |
| EVM Contract Calls | ✅ | | | Covered |
| EVM Report Writing | ✅ | | | Covered |
| **HTTP Action** |  |  |  |  |
| Send Request | ✅ | | | Covered |
| **Consensus** |  |  |  |  |
| RunInNodeMode | ✅ | | | Covered |


### Test Details

| Test Name | Type | Schedule | Chain | Owner | Status |
|-----------|------|----------|-------|-------|--------|
| **Proof of Reserve (Cron)** | Cron Trigger | Every minute | Geth Testnet | TBD | ✅ Active |

## Local Development

### Prerequisites
- having your local CRE ready


### Setup and Running

from the `core/scripts/cre/environment` directory:

1. **Start the environment with contracts:**
   ```bash
   go run . env start --with-contracts-version v2
   ```

2. **Start the observability stack:**
   ```bash
   ./ctf obs up
   ```

3. **Start beholder:**
   ```bash
   go run . env beholder start
   ```

4. **Deploy required contracts:**

   Deploy data reader:
   ```bash
   go run . examples contracts deploy-balance-reader
   ```

   Deploy a Fake Price Provider:
   ```bash
   go run . examples deploy-fake-price-provider
   ```

5. **Deploy the workflow locally:**

    Create a `config.yaml` file with the following content (update addresses and parameters as needed, this should work for the proof-of-reserve canary):
   ```yaml
    auth_key_secret_name: ""
    chain_selector: 3379446385462418246
    chain_family: evm
    chain_id: "1337"
    balancereaderconfig:
    balance_reader_address: 0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44
    addresses_to_read:
        - 0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44
        - 0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44
    computeconfig:
    feed_id: 0x12345678901234567890123456789012
    url: http://host.docker.internal:80/fake/api/price
    auth_key_secret_name: test-auth-key
    consumer_address: 0x9E545E3C0baAB3E08CdfD552C960A1050f373042
    write_target_name: write_geth-testnet@1.0.0
   ```

   Deploy the workflow using the following command:
   ```bash
   go run . env workflow deploy \
     -w ./../../../../system-tests/tests/canaries_sentinels/proof-of-reserve/cron-based/main.go \
     -n canary_proof_of_reserve \
     --compile \
     -c ./../../../../system-tests/tests/canaries_sentinels/proof-of-reserve/cron-based/config.yaml \
     --with-contracts-version v2
   ```

## Monitoring and Alerting

### Prerequisites

1. **Local Grafana**: Started locally on `http://localhost:3000` when running `./ctf obs up`
2. **CRE Observability Stack**: Clone and setup the chainlink-observability repository

### Setting Up CRE Dashboards

1. **Clone the chainlink-observability repository:**
   ```bash
   git clone https://github.com/smartcontractkit/chainlink-observability
   cd chainlink-observability
   ```

2. **Deploy CRE local dashboards:**
   ```bash
   ./deploy-cre-local.sh
   ```

   This script will:
   - Configure Grafana datasources
   - Import CRE-specific dashboards
   - Set up monitoring for workflow executions
   - Configure alerting rules for canary failures

3. **Access the CRE Dashboard:**
   - Open `http://localhost:3000` in your browser
   - Navigate to the CRE Local Dashboard
   - You should see workflow execution metrics and canary health status

## Contributing

When adding new canary workflows:

1. Create a new directory under the appropriate category
2. Implement the workflow following the existing patterns
3. Update the Coverage Matrix in this README
4. Add configuration examples
5. Include local development instructions
6. Set up appropriate monitoring and alerting

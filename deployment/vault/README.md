# Vault

Vault provides [changesets](https://github.com/smartcontractkit/chainlink/tree/develop/deployment#changsets) for managing treasury and balance monitoring operations using MCMS. It provides various operations including native token transfers, balance monitoring, and treasury management across multiple EVM chains.

## Features

- **Batch Native Token Transfer**: Transfer different amounts of native tokens (ETH, BNB, etc.) to multiple addresses across multiple chains
- **Balance Monitoring**: Monitor and track native token balances (ETH, BNB, etc.) across multiple chains
- **Address Whitelisting**: Use datastore for approved destination addresses
- **MCMS Integration**: Full integration with MCMS for secure multi-sig governance
- **Cross-Chain Support**: Execute operations across multiple EVM chains simultaneously

## Durable pipeline (prod_mainnet) – RBAC Timelock and MCMS ownership

Contracts are deployed by KMS via CI/CD. After deployment and config, ownership of the ManyChainMultiSig contracts (Bypasser, Canceller, Proposer; **CallProxy is left as-is**) must be transferred to the RBAC Timelock, and the deployer must renounce its ADMIN role on the Timelock so that only the Timelock owns the MCMS contracts and is the sole admin of itself.

### Flow 1: New chain (CCIP-style, one pipeline run)

Use when adding a **new** chain selector. One YAML, one pipeline run.

| Step | Changeset | What happens |
|------|-----------|--------------|
| 1 | **deploy_timelock** | KMS deploys RBAC Timelock and MCMS (Bypasser, Canceller, Proposer, CallProxy). Signers and config come from the deploy payload. |
| 2 | **transfer_mcms_ownership_to_timelock** | KMS transfers ownership of Bypasser/Canceller/Proposer to the Timelock and the changeset **builds the accept-ownership MCMS proposal** (output artifact). |
| 3 | **renounce_timelock_deployer** | KMS renounces its ADMIN role on the RBAC Timelock (same run). |

After the run: CI opens a **proposal PR**. Sign the proposal (e.g. Ledger) and execute it (e.g. `execute-proposal` label). **No further pipeline run** — renounce already ran in step 3.

→ Use template: [New chain (deploy + transfer + renounce)](#1-new-chain-deploy--transfer--renounce).

### Flow 2: Migration (existing chains, two pipeline runs)

Use when chain selectors are **already in the datastore** (timelock + MCMS deployed and configured; ownership not yet handed over). Do **not** redeploy.

| Run | What you run | What happens |
|-----|--------------|--------------|
| 1 | **transfer_mcms_ownership_to_timelock** | KMS transfers Bypasser/Canceller/Proposer to the Timelock and the changeset builds the accept-ownership proposal. |
| (human) | Sign and execute the proposal (proposal PR from CI) | Timelock accepts ownership on-chain. |
| 2 | **renounce_timelock_deployer** (same `chainSelectors`) | KMS renounces its ADMIN role on the Timelock for those chains. |

→ Use templates: [Migration – transfer](#2-migration--transfer-only), then [Migration – renounce](#3-migration--renounce-only).

Pipeline payloads are resolved via the vault resolvers in **chainlink-deployments** (`environment`, `chainSelectors`, and optionally `timelockIdentifier` for transfer). Input files live under `domains/vault/<environment>/durable_pipelines/inputs/<name>.yaml`.

---

### YAML templates

Copy one of the following into `domains/vault/prod_mainnet/durable_pipelines/inputs/<filename>.yaml`. Replace chain selectors and signer addresses with your values.

#### 1. New chain (deploy + transfer + renounce)

**When:** New chain selector. One pipeline run; then sign and execute the proposal.

**Suggested filename:** `new_chain_deploy_transfer_renounce.yaml` (or similar).

```yaml
environment: prod_mainnet
domain: vault
changesets:
  - deploy_timelock:
      payload:
        # One entry per chain selector. Replace with your chain selectors and signers.
        5009297550715157269:   # e.g. Ethereum mainnet
          proposer:
            quorum: 1
            signers: ["0x..."]   # replace with proposer signer address(es)
          bypasser:
            quorum: 1
            signers: ["0x..."]   # replace with bypasser signer address(es)
          canceller:
            quorum: 1
            signers: ["0x..."]   # replace with canceller signer address(es)
          timelockmindelay: 0   # e.g. 86400 for 24h
  - transfer_mcms_ownership_to_timelock:
      payload:
        environment: prod_mainnet
        chainSelectors: [5009297550715157269]   # same as in deploy_timelock
  - renounce_timelock_deployer:
      payload:
        environment: prod_mainnet
        chainSelectors: [5009297550715157269]   # same as above
```

#### 2. Migration – transfer only

**When:** Existing chains; first step of migration. Run this pipeline; then sign and execute the proposal; then run the [renounce](#3-migration--renounce-only) pipeline.

**Suggested filename:** `transfer_mcms_ownership.yaml`.

```yaml
environment: prod_mainnet
domain: vault
changesets:
  - transfer_mcms_ownership_to_timelock:
      payload:
        environment: prod_mainnet
        chainSelectors:
          - 5009297550715157269   # add all existing chain selectors to migrate
          # - 12345678901234567890
# Optional under payload:
# timelockIdentifier: "default"   # omit or use "default" for default timelock qualifier
```

#### 3. Migration – renounce only

**When:** Run **after** the accept-ownership proposal from [transfer](#2-migration--transfer-only) has been executed. Use the same `chainSelectors` as in the transfer run.

**Suggested filename:** `renounce_timelock_deployer.yaml`.

```yaml
environment: prod_mainnet
domain: vault
changesets:
  - renounce_timelock_deployer:
      payload:
        environment: prod_mainnet
        chainSelectors:
          - 5009297550715157269   # same list as in transfer_mcms_ownership_to_timelock
          # - 12345678901234567890
```

---

### How to run

1. **Create the input YAML** in `domains/vault/prod_mainnet/durable_pipelines/inputs/` using one of the [YAML templates](#yaml-templates) above.

2. **Open a PR** against `main` in chainlink-deployments with the new/updated YAML.

3. **Trigger execution:** Comment **`/run-pipelines`** on the PR. CI runs the pipeline and persists artifacts.

4. **Proposal PR:** CI opens a **proposal PR** with the accept-ownership MCMS proposal. Then:
   - Check out the proposal PR branch, sign the proposal with Ledger (see [Signing proposals](https://docs.cld.cldev.sh/guides/mcms/signing-proposals/)), push your signature.
   - Add the **`execute-proposal`** label on that PR to execute on-chain.
   - **New chain (three-step YAML):** Renounce already ran; no further pipeline run.
   - **Migration:** After the proposal is executed, run the **renounce** pipeline ([template 3](#3-migration--renounce-only)) with the same `chainSelectors`.

5. **Local run (optional)** from chainlink-deployments repo root:
   ```bash
   go run . durable-pipeline run --environment prod_mainnet --input-file transfer_mcms_ownership.yaml --changeset transfer_mcms_ownership_to_timelock
   go run . durable-pipeline run --environment prod_mainnet --input-file renounce_timelock_deployer.yaml --changeset renounce_timelock_deployer
   ```

---

## Changesets

Vault currently provides the following changesets, with more planned for future releases:

### 1. BatchNativeTransfer

Executes batch native token transfers from timelock-owned funds to whitelisted addresses.

**Execution Order:**

For a successful batch transfer, the following steps should be completed in order:

1. **SetWhitelist** - First, set up the whitelist of approved destination addresses
2. **Fund Timelock Contracts** - Ensure timelock contracts have sufficient native token balances (FundTimelock changeset is for testing only - in production, fund timelocks through appropriate governance processes)
3. **BatchNativeTransfer** - Execute the actual batch transfers

**Internal Workflow:**

The BatchNativeTransfer changeset follows this internal sequence:

1. **Validation Phase** - Validates all transfers against whitelist and checks timelock balances
2. **Execution Phase** - Either executes transfers directly or generates MCMS proposals depending on configuration

**Configuration:**

```go
config := types.BatchNativeTransferConfig{
    TransfersByChain: map[uint64][]types.NativeTransfer{
        16015286601757825753: {{To: "0x742d35cc64ca395db82e2e3e8fa8bc6d1b7c0832", Amount: big.NewInt(10000000000000000)}, {To: "0x892d35cc64ca395db82e2e3e8fa8bc6d1b7c0842", Amount: big.NewInt(1000000000000000)}}, // Sepolia
        13264668187771770619: {{To: "0x123456789012345678901234567890123456789a", Amount: big.NewInt(20000000000000000)}}, // BSC Testnet
    },
    MCMSConfig: &proposalutils.TimelockConfig{
        MinDelay: 86400, // 24 hour delay
    },
    Description:    "Monthly team payments",
}

output, err := BatchNativeTransferChangeset.Apply(env, config)
```

### 2. FundTimelock

Funds timelock contracts with native tokens for future transfers. **Note: This changeset is intended for testing purposes only. In production environments, timelock contracts should be funded through appropriate governance processes.**

**Configuration:**

```go
config := types.FundTimelockConfig{
    FundingByChain: map[uint64]*big.Int{
        16015286601757825753: big.NewInt(5000000000000000000), // 5 ETH (Sepolia)
        13264668187771770619: big.NewInt(10000000000000000000), // 10 BNB (BSC Testnet)
    },
}

output, err := FundTimelockChangeset.Apply(env, config)
```

### 3. SetWhitelist

Sets whitelist state for approved destination addresses using datastore.

**Setting Whitelist:**

```go
config := types.SetWhitelistConfig{
	WhitelistByChain: map[uint64][]types.WhitelistAddress{
		16015286601757825753: { // Sepolia
			{
				Address:     "0x742d35cc64ca395db82e2e3e8fa8bc6d1b7c0832",
				Description: "Team A",
				Labels:      []string{"team", "monthly_payment"},
			},
			{
				Address:     "0x892d35cc64ca395db82e2e3e8fa8bc6d1b7c0842",
				Description: "Team C",
				Labels:      []string{"team", "monthly_payment"},
			},
		},
		13264668187771770619: { // BSC Testnet
			{
				Address:     "0x123456789012345678901234567890123456789a",
				Description: "Team B",
				Labels:      []string{"team", "monthly_payment"},
			},
		},
	},
}
output, err := SetWhitelistChangeset.Apply(env, config)
```

### 4. Transfer MCMS ownership to Timelock (pipeline: `transfer_mcms_ownership_to_timelock`)

Transfers ownership of Bypasser, Canceller, and Proposer ManyChainMultiSig contracts to the RBAC Timelock (CallProxy excluded), and builds an MCMS proposal for `acceptOwnership`. Used for migration of existing chains and for new chains after `deploy_timelock` + `set_mcms_config`. See [Durable pipeline (prod_mainnet)](#durable-pipeline-prod_mainnet--rbac-timelock-and-mcms-ownership) for the full flow.

### 5. Renounce Timelock deployer (pipeline: `renounce_timelock_deployer`)

Renounces the deployer/KMS ADMIN role on the RBAC Timelock for the given chain selectors. Run after executing the accept-ownership proposal produced by `transfer_mcms_ownership_to_timelock`. See [Durable pipeline (prod_mainnet)](#durable-pipeline-prod_mainnet--rbac-timelock-and-mcms-ownership) for the full flow.

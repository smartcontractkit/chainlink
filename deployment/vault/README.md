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

### Correct flow for **new** chain selectors

1. **deploy_timelock** – Deploy RBAC Timelock and MCMS contracts (Bypasser, Canceller, Proposer, CallProxy) for the new chain(s).
2. **set_mcms_config** – Set MCMS config (signers, thresholds, etc.) for the new chain(s). Leave CallProxy config as-is.
3. **transfer_mcms_ownership_to_timelock** – KMS transfers ownership of Bypasser/Canceller/Proposer MCMS to the RBAC Timelock and produces an MCMS proposal for `acceptOwnership`.
4. **Execute the MCMS proposal** – Execute the accept-ownership proposal through the normal MCMS flow (e.g. sign and execute).
5. **renounce_timelock_deployer** – KMS renounces its ADMIN role on the RBAC Timelock. After this, only the Timelock owns the MCMS contracts and is the sole admin of itself.

### Migration for **existing** chain selectors (already deployed, no redeploy)

For chain selectors already in the datastore (contracts already deployed, config already set), do **not** redeploy. Run only the ownership handover:

1. **transfer_mcms_ownership_to_timelock** – With payload listing the existing `chainSelectors`. This performs the transfer and produces the accept-ownership proposal.
2. **Execute the MCMS proposal** – As above.
3. **renounce_timelock_deployer** – With the same `chainSelectors`. KMS renounces ADMIN on the Timelock for those chains.

Pipeline payloads are resolved via the vault resolvers in `chainlink-deployments` (e.g. `environment`, `chainSelectors`, and for transfer optionally `timelockIdentifier`, `onlyAcceptOwnership`). In prod_mainnet, **set_mcms_config** is configured to produce an MCMS proposal (rather than sending directly), so it works both before and after ownership migration; run the pipeline, then sign and execute the proposal as with other MCMS proposals.

### YAML templates (payload only)

Use these under `payload:` for the corresponding changeset in your pipeline input file (see [How to run](#how-to-run) below).

**transfer_mcms_ownership_to_timelock**

```yaml
environment: prod_mainnet
chainSelectors:
  - 5009297550715157269   # e.g. Ethereum mainnet; add all target chain selectors
# Optional:
# timelockIdentifier: ""   # Omit or use "default" for default timelock qualifier
# onlyAcceptOwnership: false   # Set true to only build accept-ownership proposal (skip transfer)
```

**renounce_timelock_deployer**

```yaml
environment: prod_mainnet
chainSelectors:
  - 5009297550715157269   # Same list as used for transfer_mcms_ownership_to_timelock
```

### How to run

Input files live in **chainlink-deployments**: `domains/vault/<environment>/durable_pipelines/inputs/<name>.yaml`.

1. **Create the input YAML** in `domains/vault/prod_mainnet/durable_pipelines/inputs/`. Example for migration (transfer then renounce):

   **Example: `transfer_mcms_ownership.yaml`**
   ```yaml
   environment: prod_mainnet
   domain: vault
   changesets:
     - transfer_mcms_ownership_to_timelock:
         payload:
           environment: prod_mainnet
           chainSelectors:
             - 5009297550715157269
             # add other chain selectors as needed
   ```

   **Example: `renounce_timelock_deployer.yaml`** (run after executing the accept-ownership proposal)
   ```yaml
   environment: prod_mainnet
   domain: vault
   changesets:
     - renounce_timelock_deployer:
         payload:
           environment: prod_mainnet
           chainSelectors:
             - 5009297550715157269
             # same chain selectors as transfer
   ```

2. **Open a PR** against `main` in chainlink-deployments with the new/updated YAML.

3. **Trigger execution**: Comment **`/run-pipelines`** on the PR. CI will run the pipeline and persist artifacts.

4. **If you ran `transfer_mcms_ownership_to_timelock`**: CI will open a **proposal PR** with the accept-ownership MCMS proposal. Then:
   - Check out the proposal PR branch, sign the proposal with Ledger (see [Signing proposals](https://docs.cld.cldev.sh/guides/mcms/signing-proposals/)), push your signature.
   - Add the **`execute-proposal`** label on that PR to execute on-chain.
   - After the proposal is executed, run the **`renounce_timelock_deployer`** pipeline (new input file or new PR) with the same `chainSelectors`.

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

# Vault

Vault provides [changesets](https://github.com/smartcontractkit/chainlink/tree/develop/deployment#changsets) for managing treasury and balance monitoring operations using MCMS. It provides various operations including native token transfers, balance monitoring, and treasury management across multiple EVM chains.

## Features

- **Batch Native Token Transfer**: Transfer different amounts of native tokens (ETH, BNB, etc.) to multiple addresses across multiple chains
- **Balance Monitoring**: Monitor and track native token balances (ETH, BNB, etc.) across multiple chains
- **Address Whitelisting**: Use datastore for approved destination addresses
- **MCMS Integration**: Full integration with MCMS for secure multi-sig governance
- **Cross-Chain Support**: Execute operations across multiple EVM chains simultaneously

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

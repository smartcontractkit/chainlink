# CCIP Migrations for Solana

## Table of Contents

- [Deploy Chain Contracts (and IDL Upload)](#deploy-chain-contracts)
  - [Use Remote Artifacts](#use-remote-artifacts)
  - [Build Locally](#build-locally)
  - [Upgrade Contracts](#upgrade-contracts)
- [Set OCR3 Config on OffRamp](#set-ocr3-config-on-offramp)
- [Add Remote Chain](#add-remote-chain)
  - [Add Remote Chain to Router](#add-remote-chain-to-router)
  - [Add Remote Chain to FeeQuoter](#add-remote-chain-to-feequoter)
  - [Add Remote Chain to OffRamp](#add-remote-chain-to-offramp)
  - [Disable Remote Chain](#disable-remote-chain)
  - [Add Solana EVM Lane](#add-solana--evm-lane)
- [Billing](#billing)
  - [Add Billing Token](#add-billing-token)
  - [Update Billing Token](#update-billing-token)
  - [Withdraw Billed Funds](#withdraw-billed-funds)
- [Token Operations](#token-operations)
  - [Deploy Token](#deploy-token)
  - [Upload Token Metadata](#upload-token-metadata)
  - [Mint Token](#mint-token)
  - [Create ATA](#create-ata)
  - [Set Token Authority](#set-token-authority)
- [Token Pool Operations](#token-pool-operations)
  - [Deploying Token Pool Executables](#deploy-new-token-pool-executables)
    - [Deploy New Token Pool Executable (CLL)](#deploy-new-token-pool-executable-cll)
    - [Deploy New Token Pool Executable (Parnters)](#deploy-new-token-pool-executable-parnters)
  - [Adding Token Pools for Solana Tokens (Deploying Token Pool PDAs)](#adding-token-pools-for-solana-tokens-deploying-token-pool-pdas)
  - [Updating Token Pool Rate Limits for a EVM <> Solana Lane](#updating-token-pool-rate-limits-for-a-evm--solana-lane)
- [Ownership/Authority Transfer](#ownershipauthority-transfer)
- [Partner Token Pool E2E](#partner-token-pool-e2e)
- [Verify Contracts](#verify-contracts)

## Deploy Chain Contracts

### Use Remote Artifacts
```golang
// deploy chain contracts on solana
registry.Add("0029_deploy_chain_contracts_on_solana",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.DeployChainContractsChangeset(e, ccipchangesetsolana.DeployChainContractsConfig{
        HomeChainSelector: shared.HomeChainMainnet.Selector,
        ChainSelector:     chainsel.SOLANA_MAINNET.Selector,
        BuildConfig: &ccipchangesetsolana.BuildSolanaConfig{
          DestinationDir: e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].ProgramsPath,
          // there should be artifacts available with this commit sha here
          // https://github.com/smartcontractkit/chainlink-ccip/releases
          GitCommitSha:   "be8d09930aaa",
        },
        MCMSWithTimelockConfig: &mcmsConfigs,
        ContractParamsPerChain: ccipchangesetsolana.ChainContractParams{
          FeeQuoterParams: shared.DeriveSolanaFeeQuoterInitParams(
            shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].LinkToken,
            nil,
          ),
          OffRampParams: ccipchangesetsolana.OffRampParams{
            // https://smartcontract-it.atlassian.net/browse/NONEVM-1717
            EnableExecutionAfter: int64((20 * time.Minute).Seconds()),
          },
        },
      })
    },
  ).With(struct{}{}))

// upload idl
registry.Add("0030_deploy_idl_for_chain_contracts_on_solana",
  migrations.ConfigureLegacy(ccipchangesetsolana.UploadIDL).
    With(ccipchangesetsolana.IDLConfig{
      ChainSelector:    chainsel.SOLANA_MAINNET.Selector,
      GitCommitSha:     "be8d09930aaa",
      Router:           true,
      FeeQuoter:        true,
      OffRamp:          true,
      RMNRemote:        true,
      AccessController: true,
      MCM:              true,
      Timelock:         true,
      BurnMintTokenPoolMetadata: []string{
        ccipshared.CLLMetadata,
      },
      LockReleaseTokenPoolMetadata: []string{
        ccipshared.CLLMetadata,
      },
    }))
```

### Build Locally

```golang
// deploy chain contracts by building locally
registry.Add("0019_deploy_solana_ccip_contracts",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
    return ccipchangesetsolana.DeployChainContractsChangeset(e, ccipchangesetsolana.DeployChainContractsConfig{
      HomeChainSelector: shared.HomeChainMainnet.Selector,
      ChainSelector:     chainsel.SOLANA_MAINNET.Selector,
      MCMSWithTimelockConfig: &McmsConfigs,
      BuildConfig: &ccipchangesetsolana.BuildSolanaConfig{
        DestinationDir: e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].ProgramsPath,
        GitCommitSha:   "be8d09930aaa",
        LocalBuild: ccipchangesetsolana.LocalBuildConfig{
          CleanDestinationDir: true,
          GenerateVanityKeys:  true,
          BuildLocally:        true,
        },
      },
      ContractParamsPerChain: ccipchangesetsolana.ChainContractParams{
        FeeQuoterParams: shared.DeriveSolanaFeeQuoterInitParams(
          shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].LinkToken,
          nil,
        ),
        OffRampParams: ccipchangesetsolana.OffRampParams{
          EnableExecutionAfter: int64((20 * time.Minute).Seconds()),
        },
      },
    })
  },
).With(struct{}{}))
```

### Upgrade Contracts

```golang
// upgrade contracts in place using remote resolution
registry.Add("0077_upgrade_chain_contracts_on_solana",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.DeployChainContractsChangeset(e, ccipchangesetsolana.DeployChainContractsConfig{
        HomeChainSelector: shared.HomeChainMainnet.Selector,
        ChainSelector:     chainsel.SOLANA_MAINNET.Selector,
        BuildConfig: &ccipchangesetsolana.BuildSolanaConfig{
          DestinationDir: e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].ProgramsPath,
          GitCommitSha:   "0ee732e80586",
        },
        UpgradeConfig: ccipchangesetsolana.UpgradeConfig{
          NewFeeQuoterVersion:            &deployment.Version1_0_0,
          NewRouterVersion:               &deployment.Version1_0_0,
          NewBurnMintTokenPoolVersion:    &deployment.Version1_0_0,
          NewLockReleaseTokenPoolVersion: &deployment.Version1_0_0,
          NewOffRampVersion:              &deployment.Version1_0_0,
          UpgradeAuthority:               e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].DeployerKey.PublicKey(),
          SpillAddress:                   e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].DeployerKey.PublicKey(),
          MCMS:                           timelockConfig,
        },
      })
    },
  ).With(struct{}{})

// upload idl for upgraded contracts
registry.Add("0079_upgrade_idl_for_chain_contracts_on_solana",
  migrations.ConfigureLegacy(ccipchangesetsolana.UpgradeIDL).
    With(ccipchangesetsolana.IDLConfig{
      ChainSelector: chainsel.SOLANA_MAINNET.Selector,
      GitCommitSha:  "0ee732e80586",
      Router:        true,
      FeeQuoter:     true,
      OffRamp:       true,
      RMNRemote:     true,
      // AccessController:     true,
      // MCM:                  true,
      // Timelock:             true,
      BurnMintTokenPoolMetadata: []string{
        ccipshared.CLLMetadata,
      },
      LockReleaseTokenPoolMetadata: []string{
        ccipshared.CLLMetadata,
      },
    }))

// upgrade partner token pool contracts using local build on CI
// this will do a local build for the partner token pool using the upgrade commit sha
// but it will retain the program_id of that token pool (by syncing the key with the upgraded rust file)
registry.Add("0305_update_ohm_token_pool_as_deployer",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.DeployChainContractsChangeset(e, ccipchangesetsolana.DeployChainContractsConfig{
        HomeChainSelector: shared.HomeChainTestnet.Selector,
        ChainSelector:     chainsel.SOLANA_DEVNET.Selector,
        BuildConfig: &ccipchangesetsolana.BuildSolanaConfig{
          DestinationDir: e.BlockChains.SolanaChains()[chainsel.SOLANA_DEVNET.Selector].ProgramsPath,
          GitCommitSha:   "0ee732e80586",
          LocalBuild: ccipchangesetsolana.LocalBuildConfig{
            BuildLocally: true,
            UpgradeKeys: map[cldf.ContractType]string{
              // partner token pool address
              ccipshared.BurnMintTokenPool: solana.MustPublicKeyFromBase58("GEu7tWV9LjtmwtJF661Ugv4wvx38aoEvyQiZzXHEnhRP").String(),
            },
          },
        },
        BurnMintTokenPoolMetadata: shared.OlympusWhitegloveMetadata,
        UpgradeConfig: ccipchangesetsolana.UpgradeConfig{
          NewBurnMintTokenPoolVersion: &deployment.Version1_0_0,
          SpillAddress:                e.BlockChains.SolanaChains()[chainsel.SOLANA_DEVNET.Selector].DeployerKey.PublicKey(),
          UpgradeAuthority:            e.BlockChains.SolanaChains()[chainsel.SOLANA_DEVNET.Selector].DeployerKey.PublicKey(),
        },
      })
    },
  ).With(struct{}{}))

```

## Set OCR3 Config on OffRamp
```golang
// this typically happens after setCandidate and promoteCandidate changesets
registry.Add("0036_set_ocr_configs_in_offramps_for_solana",
  migrations.ConfigureLegacy(ccipchangesetsolana.SetOCR3ConfigSolana).
    With(v1_6.SetOCR3OffRampConfig{
      HomeChainSel: shared.HomeChainMainnet.Selector,
      RemoteChainSels: []uint64{
        chainsel.SOLANA_MAINNET.Selector,
      },
      CCIPHomeConfigType: globals.ConfigTypeActive,
    }))
```


## Add Remote Chain

Use the below three changesets to add a remote chain to Solana

### Add Remote Chain to Router

```golang
registry.Add("0040_add_remote_chains_solana_router",
		migrations.ConfigureLegacy(
			func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
				return ccipchangesetsolana.AddRemoteChainToRouter(e, ccipchangesetsolana.AddRemoteChainToRouterConfig{
					ChainSelector: chainsel.SOLANA_MAINNET.Selector,
					UpdatesByChain: map[uint64]*ccipchangesetsolana.RouterConfig{
						chainsel.ETHEREUM_MAINNET.Selector: {
							RouterDestinationConfig: ccip_router.DestChainConfig{
								AllowedSenders:   [],
								AllowListEnabled: true,
							},
						},
					},
          MCMS: mcmsConfig, // if router owned by timelock
				})
			},
		).With(struct{}{}))
```

### Add Remote Chain to FeeQuoter

```golang
	registry.Add("0041_add_remote_chains_solana_fq",
		migrations.ConfigureLegacy(ccipchangesetsolana.AddRemoteChainToFeeQuoter).
			With(
				ccipchangesetsolana.AddRemoteChainToFeeQuoterConfig{
					ChainSelector: chainsel.SOLANA_MAINNET.Selector,
					UpdatesByChain: map[uint64]*ccipchangesetsolana.FeeQuoterConfig{
						chainsel.ETHEREUM_MAINNET.Selector: {
							FeeQuoterDestinationConfig: shared.FeeQuoterSVMDestChainConfigByDestChain[chainsel.ETHEREUM_MAINNET.Selector],
						},
					},
          MCMS: mcmsConfig, // if feeQuoter owned by timelock
				},
			))
```

### Add Remote Chain to OffRamp

```golang
	registry.Add("0042_add_remote_chains_solana_offramp",
		migrations.ConfigureLegacy(ccipchangesetsolana.AddRemoteChainToOffRamp).
			With(
				ccipchangesetsolana.AddRemoteChainToOffRampConfig{
					ChainSelector: chainsel.SOLANA_MAINNET.Selector,
					UpdatesByChain: map[uint64]*ccipchangesetsolana.OffRampConfig{
						chainsel.ETHEREUM_MAINNET.Selector: {
							EnabledAsSource: true,
						},
					},
          MCMS: mcmsConfig, // if offRamp owned by timelock
				},
			))
```

### Disable Remote Chain
```golang
// use this changeset to disable a remote chain on solana
var _ cldf.ChangeSet[DisableRemoteChainConfig] = DisableRemoteChain
```

### Add Solana <> EVM Lane

```golang
registry.Add("0127_enable_solana_sonic_lane",
  migrations.Configure(crossfamily.AddEVMAndSolanaLaneChangeset).With(
    crossfamily.AddMultiEVMSolanaLaneConfig{
      MCMSConfig: mcmsConfig,
      SolanaChainSelector: chainsel.SOLANA_MAINNET.Selector,
      Configs: []crossfamily.AddRemoteChainE2EConfig{
        {
          EVMChainSelector:                      chainsel.SONIC_MAINNET.Selector,
          IsTestRouter:                          true,
          EVMFeeQuoterDestChainInput:            shared.DeriveFeeQuoterDestChainConfigEVMToSolana(nil),
          InitialSolanaGasPriceForEVMFeeQuoter:  chainConfigs[chainsel.SOLANA_MAINNET.Selector].GasPrice,
          InitialEVMTokenPricesForEVMFeeQuoter:  chainConfigs[chainsel.SONIC_MAINNET.Selector].TokenPrices,
          IsRMNVerificationDisabledOnEVMOffRamp: true,
          SolanaRouterConfig: ccipchangesetsolana.RouterConfig{
            RouterDestinationConfig: ccip_router.DestChainConfig{
              AllowedSenders:   []solana.PublicKey{},
              AllowListEnabled: false,
            },
          },
          SolanaOffRampConfig: ccipchangesetsolana.OffRampConfig{
            EnabledAsSource: true,
          },
          SolanaFeeQuoterConfig: ccipchangesetsolana.FeeQuoterConfig{
            FeeQuoterDestinationConfig: shared.FeeQuoterSVMDestChainConfigByDestChain[chainsel.SONIC_MAINNET.Selector],
          },
        },
      },
    }))
```


## Billing

### Add Billing Token

We always add WSOL and LINK as billing tokens during initial deployment as part of `DeployChainContractsChangeset`

```golang
// use this changeset to add a billing token to solana
var _ cldf.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset
```

### Add Token Transfer Fee

```golang
// use this changeset to add a token transfer fee for a remote chain to solana (used for very specific cases)
var _ cldf.ChangeSet[TokenTransferFeeForRemoteChainConfig] = AddTokenTransferFeeForRemoteChain
```

### Withdraw Billed Funds

```golang
// use this changeset to withdraw billed funds on solana
var _ cldf.ChangeSet[WithdrawBilledFundsConfig] = WithdrawBilledFunds
```

## Token Operations

CLI Docs: https://spl.solana.com/token#example-creating-your-own-fungible-token
You can use the cli if you have access to the deployer key of your env to make these token operations easier

### Deploy Token

```golang
// lets deploy -> create ATA -> mint
registry.Add("0107_deploy_solvbtc_token_solana",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      allowedSenders := shared.GetAllowedSendersSolToEVM(e)
      allowedATAs := make([]string, 0, len(allowedSenders))
      mintAmounts := make(map[string]uint64)
      for _, allowedSender := range allowedSenders {
        allowedATAs = append(allowedATAs, allowedSender.String())
      }

      return ccipchangesetsolana.DeploySolanaToken(e,
        ccipchangesetsolana.DeploySolanaTokenConfig{
          ChainSelector:       chainsel.SOLANA_MAINNET.Selector,
          TokenProgramName:    ccipshared.SPLTokens,
          TokenDecimals:       8,
          TokenSymbol:         shared.SolvBTCToken,
          MintPrivateKey:      solana.MustPrivateKeyFromBase58("xxx"), // ignore if vanity address not required
          ATAList:             allowedATAs,
          MintAmountToAddress: mintAmounts,
        },
      )
    },
  ).With(struct{}{}))
```

### Upload Token Metadata

NOTE:

- This needs to be done before deploying a BnM token pool for the token.
- After you deploy the BnM Token pool the mint authority of the token will be changed to the pool and you wont be able to upload metadata because the metadata uploader needs to be the current mint authority of the pool.
  - And at the time of writing this our current BnM pool implementation does not support multiSig
- Once you upload the metadata, you can deploy the BnM token pool and still come back to update the metadata
  - Because at that point the metadata will have its own authority (original mint authority)

```json
{
  "name": "SolvBTC",
  "symbol": "SolvBTC",
  "uri": "https://raw.githubusercontent.com/solv-finance/solv-btc-metadata-solana/refs/heads/main/SolvBTC/metadata.json"
}
```

```golang
// enable metaboss in domains/ccip/ci-dependencies.yaml
// you have to create another pr to disable it after that
// if you dont, then all other migrations inside domains/ccip will install metaboss
// initial upload
registry.Add("0116_add_solvbtc_token_metadata",
  migrations.ConfigureLegacy(ccipchangesetsolana.UploadTokenMetadata).
    With(ccipchangesetsolana.UploadTokenMetadataConfig{
      ChainSelector: chainsel.SOLANA_MAINNET.Selector,
      TokenMetadata: []ccipchangesetsolana.TokenMetadata{
        {
          TokenPubkey:      shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].SolvBTCToken,
          MetadataJSONPath: "solvbtc.json", // dump this in domains/ccip/mainnet/inputs
        },
      },
    }))
```

```golang
// use below inputs to UploadTokenMetadata to update metadata
type TokenMetadata struct {
	TokenPubkey solana.PublicKey
	// https://metaboss.dev/create.html#metadata
	// only to be provided on initial upload, it takes in name, symbol, uri
	// after initial upload, those fields can be updated using the update inputs
	// put the json in ccip/env/input dir in CLD
	MetadataJSONPath string
	UpdateAuthority  solana.PublicKey // used to set update authority of the token metadata PDA after initial upload
	// https://metaboss.dev/update.html#update-name
	UpdateName string // used to update the name of the token metadata PDA after initial upload
	// https://metaboss.dev/update.html#update-symbol
	UpdateSymbol string // used to update the symbol of the token metadata PDA after initial upload
	// https://metaboss.dev/update.html#update-uri
	UpdateURI string // used to update the uri of the token metadata PDA after initial upload
}
```

### Mint Token
```golang
// use this changeset to mint the token to an address
var _ cldf.ChangeSet[MintSolanaTokenConfig] = MintSolanaToken
```

### Create ATA

```golang
// use this changeset to create ATAs for a token
var _ cldf.ChangeSet[CreateSolanaTokenATAConfig] = CreateSolanaTokenATA
```

### Set Token Authority
```golang
// use this changeset to set the authority of a token
var _ cldf.ChangeSet[SetTokenAuthorityConfig] = SetTokenAuthority
```

## Token Pool Operations

### Deploying Token Pool Executables

##### Deploy New Token Pool Executable (CLL)

- We deploy the BnM and LnR token pool executables as part of DeployChainContractsChangeset for Chainlink

##### Deploy New Token Pool Executable (Parnters)

- For partners, you need to deploy a new executable like below
- At the time of writing this, there are talks around using the same executable for partners as well (so something to keep in mind for later)
- We have already deployed a few partner token pools, so how we will upgrade all of them is TBD

```golang
// this can deploy one token pool at a time based on the metadata you provide
// it will not re-deploy any other contracts
registry.Add("0091_deploy_partner_token_pool_for_solana_solv_bnm_zeus_lnr",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.DeployChainContractsChangeset(e, ccipchangesetsolana.DeployChainContractsConfig{
        HomeChainSelector: shared.HomeChainMainnet.Selector,
        ChainSelector:     chainsel.SOLANA_MAINNET.Selector,
        BuildConfig: &ccipchangesetsolana.BuildSolanaConfig{
          DestinationDir: e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].ProgramsPath,
          GitCommitSha:   "0ee732e80586",
          LocalBuild: ccipchangesetsolana.LocalBuildConfig{
            BuildLocally: true,
          },
        },
        BurnMintTokenPoolMetadata:    shared.SolvWhitegloveMetadata,
      })
    },
  ).With(struct{}{}))
```

### Adding Token Pools for Solana Tokens (Deploying Token Pool PDAs)

You can use one changeset to
- Deploy token pool pda
- Deploy lookup table
- Register Token to Token Admin Registry
- Accept ownership of the Registry
- SetPool on the Token Admin Registry
- Configure EVM Pools On Solana (you need to make sure if the EVM pools were not deployed in CLD, you import their addresses to the address book)
- Configure Solana Pools on EVM (you need to make sure if the EVM pools were not deployed in CLD, you import their addresses to the address book)


```golang
registry.Add("0321_deploy_CCIP_TEST_solana_token_pool",
		migrations.ConfigureLegacy(ccipchangesetsolana.E2ETokenPoolv2).
			With(ccipchangesetsolana.E2ETokenPoolConfigv2{
				ChainSelector: chainsel.SOLANA_DEVNET.Selector,
				MCMS:          mcmsConfigForCS,
				E2ETokens: []ccipchangesetsolana.E2ETokenConfig{
					{
						TokenPubKey: solana.MustPublicKeyFromBase58("3kffP9DNcWKBZFUqJkx5pgNMbHj1q73kzkgiuUiZJpq8"),
						Metadata:    shared.SoylanaManlettWhitegloveMetadata,
						PoolType:    shared.SolanaBnMTokenPoolEnumPtr,
						SolanaToEVMRemoteConfigs: map[uint64]ccipchangesetsolana.EVMRemoteConfig{
							chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector: {
								TokenSymbol:       ccipshared.TokenSymbol("SOYMAN"),
								PoolType:          ccipshared.BurnMintTokenPool,
								PoolVersion:       ccipshared.CurrentTokenPoolVersion,
								RateLimiterConfig: shared.DefaultRateLimiterConfigForTestTokensSolana,
							},
							chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
								TokenSymbol:       ccipshared.TokenSymbol("SOYMAN"),
								PoolType:          ccipshared.BurnMintTokenPool,
								PoolVersion:       ccipshared.CurrentTokenPoolVersion,
								RateLimiterConfig: shared.DefaultRateLimiterConfigForTestTokensSolana,
							},
						},
						EVMToSolanaRemoteConfigs: v1_5_1.ConfigureTokenPoolContractsConfig{
							MCMS:        mcmsConfigForCS,
							TokenSymbol: ccipshared.TokenSymbol("SOYMAN"),
							PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
								chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
									Type:    ccipshared.LockReleaseTokenPool,
									Version: deployment.Version1_5_1,
									SolChainUpdates: map[uint64]v1_5_1.SolChainUpdate{
										chainsel.SOLANA_DEVNET.Selector: {
											TokenAddress:      shared.SolanaTokenAddress[chainsel.SOLANA_DEVNET.Selector].SolvBTCToken.String(),
											Type:              ccipshared.BurnMintTokenPool,
											Metadata:          shared.SolvWhitegloveMetadata,
											RateLimiterConfig: shared.SolvRateLimitConfigEvmToSolana,
										},
									},
								},
							},
						},
					},
				},
			}),
		migrations.OnlyLoadChainsFor(chainsel.SOLANA_DEVNET.Selector, chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector, chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector),
	)
```

### Updating Token Pool Rate Limits for a EVM <> Solana Lane

This changeset allows you to generate one proposal for EVM + Solana rate limit changes for a particular token bridge.

```golang
  registry.Add("0172_configure_pepe_pool_solana",
		migrations.ConfigureLegacy(
			func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
				timelockSignerPDA, err := ccipchangesetsolana.FetchTimelockSigner(e, chainsel.SOLANA_MAINNET.Selector)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
				}

				return ccipchangesetsolana.E2ETokenPool(e, ccipchangesetsolana.E2ETokenPoolConfig{
					RemoteChainTokenPool: []ccipchangesetsolana.SetupTokenPoolForRemoteChainConfig{
						{
							SolChainSelector: chainsel.SOLANA_MAINNET.Selector,
							RemoteTokenPoolConfigs: []ccipchangesetsolana.RemoteChainTokenPoolConfig{
								{
									SolTokenPubKey: shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].PepeToken,
									SolPoolType:    shared.SolanaBnMTokenPoolEnumPtr,
									Metadata:       shared.PepeWhitegloveMetadata,
									EVMRemoteConfigs: map[uint64]ccipchangesetsolana.EVMRemoteConfig{
										chainsel.ETHEREUM_MAINNET.Selector: {
											TokenSymbol:       ccipshared.TokenSymbol(shared.PepeToken),
											PoolType:          ccipshared.LockReleaseTokenPool,
											PoolVersion:       ccipshared.CurrentTokenPoolVersion,
											RateLimiterConfig: shared.DefaultRateLimiterConfigForTestTokensSolana, // circumvent bug
										},
									},
								},
								{
									SolTokenPubKey: shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].PepeToken,
									SolPoolType:    shared.SolanaBnMTokenPoolEnumPtr,
									Metadata:       shared.PepeWhitegloveMetadata,
									EVMRemoteConfigs: map[uint64]ccipchangesetsolana.EVMRemoteConfig{
										chainsel.ETHEREUM_MAINNET.Selector: {
											TokenSymbol:       ccipshared.TokenSymbol(shared.PepeToken),
											PoolType:          ccipshared.LockReleaseTokenPool,
											PoolVersion:       ccipshared.CurrentTokenPoolVersion,
											RateLimiterConfig: shared.PepeRateLimitConfigSolanaToEvm, // set actual value
										},
									},
								},
							},
						},
					},
					MCMS: timelockConfig,
				})
			},
		).With(struct{}{}))
```

## Ownership/Authority Transfer

```golang
// transfer contract owner to timelock
// owner here is the onchain authority that is enforced using rust code
	registry.Add("0143_transfer_to_timelock_solana",
		migrations.ConfigureLegacy(
			func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
				return ccipchangesetsolana.TransferCCIPToMCMSWithTimelockSolana(e, ccipchangesetsolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
					MCMSCfg: *timelockConfig,
					ContractsByChain: map[uint64]ccipchangesetsolana.CCIPContractsToTransfer{
						chainsel.SOLANA_MAINNET.Selector: {
							Router:    true,
							FeeQuoter: true,
							OffRamp:   true,
							BurnMintTokenPools: map[string][]solana.PublicKey{
								ccipshared.CLLMetadata: {shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].LinkToken},
								shared.SolvWhitegloveMetadata: {
									shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].SolvBTCToken,
									shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].XSolvBTCToken,
									shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].SolvBTCJUPToken,
								},
								shared.MapleWhitegloveMetadata: {shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].SyrupUSDCToken},
							},
							LockReleaseTokenPools: map[string][]solana.PublicKey{
								shared.ZeusWhitegloveMetadata: {shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].ZBTCToken},
							},
						},
					},
				})
			},
		).With(struct{}{}))

// transfer upgrade authority of contract executable to timelock
// upgrade authority here is the authority allowed to perform an upgrade to the executable
registry.Add(
  "0274_change_upgrade_authority_to_timelock_for_solana",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      timelockSignerPDA, err := ccipchangesetsolana.FetchTimelockSigner(e, chainsel.SOLANA_DEVNET.Selector)
      if err != nil {
        return cldf.ChangesetOutput{}, err
      }

      return ccipchangesetsolana.SetUpgradeAuthorityChangeset(e,
        ccipchangesetsolana.SetUpgradeAuthorityConfig{
          ChainSelector:         chainsel.SOLANA_DEVNET.Selector,
          NewUpgradeAuthority:   timelockSignerPDA,
          SetAfterInitialDeploy: true,
          SetOffRamp:            true,
          SetMCMSPrograms:       true,
          TransferKeys: []solana.PublicKey{
							// token pools go here
							solana.MustPublicKeyFromBase58("GEu7tWV9LjtmwtJF661Ugv4wvx38aoEvyQiZzXHEnhRP"),
						},
        },
      )
    },
  ).With(struct{}{}))

// transfer back upgrade authority of contract executable to deployer key
registry.Add(
		"0304_change_upgrade_authority_to_deployer_for_solana_ohm_pool",
		migrations.ConfigureLegacy(
			func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
				return ccipchangesetsolana.SetUpgradeAuthorityChangeset(e,
					ccipchangesetsolana.SetUpgradeAuthorityConfig{
						ChainSelector:       chainsel.SOLANA_DEVNET.Selector,
						NewUpgradeAuthority: e.BlockChains.SolanaChains()[chainsel.SOLANA_DEVNET.Selector].DeployerKey.PublicKey(),
						TransferKeys: []solana.PublicKey{
							// ohm token pool
							solana.MustPublicKeyFromBase58("GEu7tWV9LjtmwtJF661Ugv4wvx38aoEvyQiZzXHEnhRP"),
						},
						MCMS: mcmsConfigForCS,
					},
				)
			},
		).With(struct{}{}))

// transfer IDL authority to timelock
registry.Add(
  "0275_test_change_idl_authority_to_timelock_for_solana",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.SetAuthorityIDL(e,
        ccipchangesetsolana.IDLConfig{
          ChainSelector:    chainsel.SOLANA_DEVNET.Selector,
          Router:           true,
          OffRamp:          true,
          FeeQuoter:        true,
          RMNRemote:        true,
          MCM:              true,
          Timelock:         true,
          AccessController: true,
          BurnMintTokenPoolMetadata: []string{
            ccipshared.CLLMetadata,
          },
          LockReleaseTokenPoolMetadata: []string{
            ccipshared.CLLMetadata,
          },
          MCMS: mcmsConfigForCS,
        },
      )
    },
  ).With(struct{}{}))
```

## Partner Token Pool E2E

```golang
// deploy token
registry.Add("0107_deploy_solvbtc_token_solana",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.DeploySolanaToken(e,
        ccipchangesetsolana.DeploySolanaTokenConfig{
          ChainSelector:       chainsel.SOLANA_MAINNET.Selector,
          TokenProgramName:    ccipshared.SPLTokens,
          TokenDecimals:       8,
          TokenSymbol:         shared.SolvBTCToken,
          MintPrivateKey:      solana.MustPrivateKeyFromBase58("xxx"), // ignore if vanity address not required
        },
      )
    },
  ).With(struct{}{}))

// enable metaboss in domains/ccip/ci-dependencies.yaml
// upload metadata
registry.Add("0116_add_solvbtc_token_metadata",
  migrations.ConfigureLegacy(ccipchangesetsolana.UploadTokenMetadata).
    With(ccipchangesetsolana.UploadTokenMetadataConfig{
      ChainSelector: chainsel.SOLANA_MAINNET.Selector,
      TokenMetadata: []ccipchangesetsolana.TokenMetadata{
        {
          TokenPubkey:      shared.SolanaTokenAddress[chainsel.SOLANA_MAINNET.Selector].SolvBTCToken,
          MetadataJSONPath: "solvbtc.json", // dump this in domains/ccip/mainnet/inputs
        },
      },
    }))

// deploy token pool executable
registry.Add("0091_deploy_partner_token_pool_for_solana_solv_bnm_zeus_lnr",
  migrations.ConfigureLegacy(
    func(e cldf.Environment, _ any) (cldf.ChangesetOutput, error) {
      return ccipchangesetsolana.DeployChainContractsChangeset(e, ccipchangesetsolana.DeployChainContractsConfig{
        HomeChainSelector: shared.HomeChainMainnet.Selector,
        ChainSelector:     chainsel.SOLANA_MAINNET.Selector,
        BuildConfig: &ccipchangesetsolana.BuildSolanaConfig{
          DestinationDir: e.BlockChains.SolanaChains()[chainsel.SOLANA_MAINNET.Selector].ProgramsPath,
          GitCommitSha:   "0ee732e80586",
          LocalBuild: ccipchangesetsolana.LocalBuildConfig{
            BuildLocally: true,
          },
        },
        BurnMintTokenPoolMetadata:    shared.SolvWhitegloveMetadata,
      })
    },
  ).With(struct{}{}))

// setup token pool solana <> lane
registry.Add("0321_deploy_CCIP_TEST_solana_token_pool",
		migrations.ConfigureLegacy(ccipchangesetsolana.E2ETokenPoolv2).
			With(ccipchangesetsolana.E2ETokenPoolConfigv2{
				ChainSelector: chainsel.SOLANA_DEVNET.Selector,
				MCMS:          mcmsConfigForCS,
				E2ETokens: []ccipchangesetsolana.E2ETokenConfig{
					{
						TokenPubKey: solana.MustPublicKeyFromBase58("3kffP9DNcWKBZFUqJkx5pgNMbHj1q73kzkgiuUiZJpq8"),
						Metadata:    shared.SoylanaManlettWhitegloveMetadata,
						PoolType:    shared.SolanaBnMTokenPoolEnumPtr,
						SolanaToEVMRemoteConfigs: map[uint64]ccipchangesetsolana.EVMRemoteConfig{
							chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector: {
								TokenSymbol:       ccipshared.TokenSymbol("SOYMAN"),
								PoolType:          ccipshared.BurnMintTokenPool,
								PoolVersion:       ccipshared.CurrentTokenPoolVersion,
								RateLimiterConfig: shared.DefaultRateLimiterConfigForTestTokensSolana,
							},
							chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
								TokenSymbol:       ccipshared.TokenSymbol("SOYMAN"),
								PoolType:          ccipshared.BurnMintTokenPool,
								PoolVersion:       ccipshared.CurrentTokenPoolVersion,
								RateLimiterConfig: shared.DefaultRateLimiterConfigForTestTokensSolana,
							},
						},
						EVMToSolanaRemoteConfigs: v1_5_1.ConfigureTokenPoolContractsConfig{
							MCMS:        mcmsConfigForCS,
							TokenSymbol: ccipshared.TokenSymbol("SOYMAN"),
							PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
								chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
									Type:    ccipshared.LockReleaseTokenPool,
									Version: deployment.Version1_5_1,
									SolChainUpdates: map[uint64]v1_5_1.SolChainUpdate{
										chainsel.SOLANA_DEVNET.Selector: {
											TokenAddress:      shared.SolanaTokenAddress[chainsel.SOLANA_DEVNET.Selector].SolvBTCToken.String(),
											Type:              ccipshared.BurnMintTokenPool,
											Metadata:          shared.SolvWhitegloveMetadata,
											RateLimiterConfig: shared.SolvRateLimitConfigEvmToSolana,
										},
									},
								},
							},
						},
					},
				},
			}),
		migrations.OnlyLoadChainsFor(chainsel.SOLANA_DEVNET.Selector, chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector, chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector),
	)
```

## Verify Contracts

```golang
// initially verify programs afer deploy
registry.Add("0086_verify_solana_programs",
  migrations.ConfigureLegacy(ccipchangesetsolana.VerifyBuild).
    With(ccipchangesetsolana.VerifyBuildConfig{
      ChainSelector:          chainsel.SOLANA_MAINNET.Selector,
      GitCommitSha:           "0ee732e80586",
      VerifyFeeQuoter:        true,
      BurnMintTokenPoolMetadata: []string{
        ccipshared.CLLMetadata,
      },
      LockReleaseTokenPoolMetadata: []string{
        ccipshared.CLLMetadata,
      },
    }))

// 2 step process for verification after mcms transfer
// that is UpgradeAuthority of contract is transferred to mcms

// Step 1
// register verify ix on timelock
// this will spit out a proposal
// get it signed and executed
registry.Add("0086_verify_solana_programs",
  migrations.ConfigureLegacy(ccipchangesetsolana.VerifyBuild).
    With(ccipchangesetsolana.VerifyBuildConfig{
      ChainSelector:          chainsel.SOLANA_MAINNET.Selector,
      GitCommitSha:           "0ee732e80586",
      VerifyFeeQuoter:        true,
      MCMS: &proposalutils.TimelockConfig{
        // if possible, we should make this very small, ideally 5 mins
        // because otherwise we will non verified updated contracts for > 3 hours
		    MinDelay: minDelay, 
	    }
      UpgradeAuthority: timelockSignerPDA,
    }))

// Step 2
registry.Add("0086_verify_solana_programs",
  migrations.ConfigureLegacy(ccipchangesetsolana.VerifyBuild).
    With(ccipchangesetsolana.VerifyBuildConfig{
      ChainSelector:          chainsel.SOLANA_MAINNET.Selector,
      GitCommitSha:           "0ee732e80586",
      VerifyFeeQuoter:        true,
      MCMS: &proposalutils.TimelockConfig{
        // if possible, we should make this very small, ideally 5 mins
        // because otherwise we will non verified updated contracts for > 3 hours
		    MinDelay: minDelay, 
	    },
      UpgradeAuthority: timelockSignerPDA,
      RemoteVerification: true // this is KEY in step 2
    }))
```
### Example Verification PRs
1. https://github.com/smartcontractkit/chainlink-deployments/pull/3907 (Step 1)
2. https://github.com/smartcontractkit/chainlink-deployments/pull/3944 (Step 2)
3. https://github.com/smartcontractkit/chainlink-deployments/actions/runs/15774291241/job/44465124148 (Successful e2e remote verification via MCMs)
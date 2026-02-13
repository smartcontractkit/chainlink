package v1_6_1

import (
	"math/big"
	"strings"
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/transparent_upgradeable_proxy"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/1_5_0/burn_mint_erc20_transparent"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

func TestTransparentUpgradeableProxyValidate(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	tests := []struct {
		name    string
		cfg     TransparentUpgradeableProxyChangesetConfig
		wantErr bool
	}{
		{
			name: "valid config - initialize",
			cfg: TransparentUpgradeableProxyChangesetConfig{
				Tokens: map[uint64]map[string]TransparentUpgradeableProxy{
					chain1: {
						"TEST_1": {
							Symbol:     "TEST_1",
							Decimals:   18,
							MaxSupply:  big.NewInt(0),
							PreMint:    big.NewInt(0),
							Initialize: true,
						},
					},
				},
			},
		},
		{
			name: "invalid config - missing implementation address",
			cfg: TransparentUpgradeableProxyChangesetConfig{
				Tokens: map[uint64]map[string]TransparentUpgradeableProxy{
					chain2: {
						"TEST_2": {
							Symbol:    "TEST_2",
							Decimals:  18,
							MaxSupply: big.NewInt(0),
							PreMint:   big.NewInt(0),
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for chainSelector, configs := range tc.cfg.Tokens {
				chain := e.Env.BlockChains.EVMChains()[chainSelector]

				for _, config := range configs {
					if config.Initialize {
						_, tx, contract, err := burn_mint_erc20_transparent.DeployBurnMintERC20Transparent(chain.DeployerKey, chain.Client)
						require.NoError(t, err)

						_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
						require.NoError(t, err)

						err = e.Env.ExistingAddresses.Save(chainSelector, contract.Address().String(), cldf.NewTypeAndVersion(shared.BurnMintERC20TransparentToken, deployment.Version1_6_1))
						require.NoError(t, err)

						config.BurnMintERC20Transparent = contract.Address()
						tc.cfg.Tokens[chainSelector][config.Symbol] = config
					}
				}
			}

			err := tc.cfg.Validate(e.Env)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTransparentUpgradeableProxyDeploy(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain := evmSelectors[0]

	const initABIString = `[{
		"name": "initialize",
		"type": "function",
		"inputs": [
			{"name": "name", "type": "string"},
			{"name": "symbol", "type": "string"},
			{"name": "decimals", "type": "uint8"},
			{"name": "maxSupply", "type": "uint256"},
			{"name": "preMint", "type": "uint256"},
			{"name": "defaultAdmin", "type": "address"}
		]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(initABIString))
	require.NoError(t, err)

	tests := []struct {
		name                   string
		deployERC20Transparent bool
		cfg                    TransparentUpgradeableProxyChangesetConfig
		wantErr                bool
	}{
		{
			name:                   "valid config",
			deployERC20Transparent: true,
			cfg: TransparentUpgradeableProxyChangesetConfig{
				Tokens: map[uint64]map[string]TransparentUpgradeableProxy{
					chain: {
						"TEST_1": {
							Symbol:     "TEST_1",
							Decimals:   18,
							MaxSupply:  big.NewInt(0),
							PreMint:    big.NewInt(0),
							Initialize: true,
						},
					},
				},
			},
		},
		{
			name: "missing implementation address",
			cfg: TransparentUpgradeableProxyChangesetConfig{
				Tokens: map[uint64]map[string]TransparentUpgradeableProxy{
					chain: {
						"TEST_2": {
							Symbol:     "TEST_2",
							Decimals:   18,
							MaxSupply:  big.NewInt(0),
							PreMint:    big.NewInt(0),
							Initialize: true,
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for chainSelector, configs := range tc.cfg.Tokens {
				chain := e.Env.BlockChains.EVMChains()[chainSelector]

				var implementation *burn_mint_erc20_transparent.BurnMintERC20Transparent

				if tc.deployERC20Transparent {
					_, tx, contract, err := burn_mint_erc20_transparent.DeployBurnMintERC20Transparent(chain.DeployerKey, chain.Client)
					require.NoError(t, err)

					_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
					require.NoError(t, err)

					implementation = contract

					err = e.Env.ExistingAddresses.Save(chainSelector, implementation.Address().String(), cldf.NewTypeAndVersion(shared.BurnMintERC20TransparentToken, deployment.Version1_6_1))
					require.NoError(t, err)
				}

				for name, config := range configs {
					if implementation != nil {
						config.BurnMintERC20Transparent = implementation.Address()
						tc.cfg.Tokens[chainSelector][name] = config
					}

					initData, err := parsedABI.Pack(
						"initialize",
						name,
						config.Symbol,
						config.Decimals,
						config.MaxSupply,
						config.PreMint,
						e.Env.BlockChains.EVMChains()[chainSelector].DeployerKey.From,
					)
					require.NoError(t, err)

					err = tc.cfg.Validate(e.Env)
					if tc.wantErr {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}

					if implementation != nil {
						_, tx, proxy, err := transparent_upgradeable_proxy.DeployTransparentUpgradeableProxy(
							chain.DeployerKey, chain.Client, config.BurnMintERC20Transparent, chain.DeployerKey.From, initData,
						)
						require.NoError(t, err)

						_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
						require.NoError(t, err)

						err = e.Env.ExistingAddresses.Save(chainSelector, proxy.Address().String(), cldf.NewTypeAndVersion(shared.TransparentUpgradeableProxy, deployment.Version1_6_1))
						require.NoError(t, err)
					}
				}
			}
		})
	}
}

func TestTransparentUpgradeableProxyGrantRoleChangesetConfig(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	tests := []struct {
		name    string
		token   string
		cfg     TransparentUpgradeableProxyGrantRoleChangesetConfig
		wantErr bool
	}{
		{
			name:  "valid config",
			token: "TEST_1",
			cfg: TransparentUpgradeableProxyGrantRoleChangesetConfig{
				Tokens: map[uint64]map[string][]TransparentUpgradeableProxyGrantRole{
					chain1: {
						"TEST_1": {
							{
								Symbol:  "TEST_1",
								Role:    TransparentUpgradeableProxyBurnerRole,
								Account: common.HexToAddress("0x000000000000000000000000000000000000dead"),
							},
							{
								Symbol:  "TEST_1",
								Role:    TransparentUpgradeableProxyMinterRole,
								Account: common.HexToAddress("0x000000000000000000000000000000000000dead"),
							},
						},
					},
				},
			},
		},
		{
			name:  "invalid config - invalid role",
			token: "TEST_2",
			cfg: TransparentUpgradeableProxyGrantRoleChangesetConfig{
				Tokens: map[uint64]map[string][]TransparentUpgradeableProxyGrantRole{
					chain2: {
						"TEST_2": {
							{
								Symbol:  "TEST_2",
								Role:    TransparentUpgradeableProxyBurnerRole,
								Account: common.HexToAddress("0x000000000000000000000000000000000000dead"),
							},
							{
								Symbol:  "TEST_2",
								Role:    "INVALID_ROLE",
								Account: common.HexToAddress("0x000000000000000000000000000000000000dead"),
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for chainSelector := range tc.cfg.Tokens {
				chain := e.Env.BlockChains.EVMChains()[chainSelector]

				_, tx, contract, err := burn_mint_erc20_transparent.DeployBurnMintERC20Transparent(chain.DeployerKey, chain.Client)
				require.NoError(t, err)

				_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
				require.NoError(t, err)

				err = e.Env.ExistingAddresses.Save(chainSelector, contract.Address().String(), cldf.NewTypeAndVersion(shared.BurnMintERC20TransparentToken, deployment.Version1_6_1))
				require.NoError(t, err)

				const initABIString = `[{
					"name": "initialize",
					"type": "function",
					"inputs": [
						{"name": "name", "type": "string"},
						{"name": "symbol", "type": "string"},
						{"name": "decimals", "type": "uint8"},
						{"name": "maxSupply", "type": "uint256"},
						{"name": "preMint", "type": "uint256"},
						{"name": "defaultAdmin", "type": "address"}
					]
				}]`

				parsedABI, err := abi.JSON(strings.NewReader(initABIString))
				require.NoError(t, err)

				initData, err := parsedABI.Pack(
					"initialize",
					tc.token,
					tc.token,
					uint8(18),
					big.NewInt(0),
					big.NewInt(0),
					e.Env.BlockChains.EVMChains()[chainSelector].DeployerKey.From,
				)
				require.NoError(t, err)

				_, tx, proxy, err := transparent_upgradeable_proxy.DeployTransparentUpgradeableProxy(
					chain.DeployerKey, chain.Client, contract.Address(), chain.DeployerKey.From, initData,
				)
				require.NoError(t, err)

				_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
				require.NoError(t, err)

				err = e.Env.ExistingAddresses.Save(chainSelector, proxy.Address().String(), cldf.NewTypeAndVersion(shared.TransparentUpgradeableProxy, deployment.Version1_6_1))
				require.NoError(t, err)
			}

			err := tc.cfg.Validate(e.Env)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBeginDefaultAdminTransferTransparentUpgradeableProxy(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	tests := []struct {
		name         string
		cfg          TransparentUpgradeableProxyChangesetConfig
		wantNewAdmin bool
	}{
		{
			name: "valid address to transfer ownership to",
			cfg: TransparentUpgradeableProxyChangesetConfig{
				Tokens: map[uint64]map[string]TransparentUpgradeableProxy{
					chain1: {
						"TEST_1": {
							Symbol:    "TEST_1",
							Decimals:  18,
							MaxSupply: big.NewInt(0),
							PreMint:   big.NewInt(0),
							NewAdmin:  common.HexToAddress("0x000000000000000000000000000000000000dead"),
						},
					},
				},
			},
			wantNewAdmin: true,
		},
		{
			name: "no new admin specified",
			cfg: TransparentUpgradeableProxyChangesetConfig{
				Tokens: map[uint64]map[string]TransparentUpgradeableProxy{
					chain2: {
						"TEST_2": {
							Symbol:    "TEST_2",
							Decimals:  18,
							MaxSupply: big.NewInt(0),
							PreMint:   big.NewInt(0),
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for chainSelector, tokens := range tc.cfg.Tokens {
				chain := e.Env.BlockChains.EVMChains()[chainSelector]

				for token, config := range tokens {
					_, tx, implementation, err := burn_mint_erc20_transparent.DeployBurnMintERC20Transparent(chain.DeployerKey, chain.Client)
					require.NoError(t, err)

					_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
					require.NoError(t, err)

					err = e.Env.ExistingAddresses.Save(chainSelector, implementation.Address().String(), cldf.NewTypeAndVersion(shared.BurnMintERC20TransparentToken, deployment.Version1_6_1))
					require.NoError(t, err)

					config.BurnMintERC20Transparent = implementation.Address()
					tc.cfg.Tokens[chainSelector][config.Symbol] = config

					const initABIString = `[{
						"name": "initialize",
						"type": "function",
						"inputs": [
							{"name": "name", "type": "string"},
							{"name": "symbol", "type": "string"},
							{"name": "decimals", "type": "uint8"},
							{"name": "maxSupply", "type": "uint256"},
							{"name": "preMint", "type": "uint256"},
							{"name": "defaultAdmin", "type": "address"}
						]
					}]`

					parsedABI, err := abi.JSON(strings.NewReader(initABIString))
					require.NoError(t, err)

					initData, err := parsedABI.Pack(
						"initialize",
						token,
						config.Symbol,
						uint8(18),
						big.NewInt(0),
						big.NewInt(0),
						e.Env.BlockChains.EVMChains()[chainSelector].DeployerKey.From,
					)
					require.NoError(t, err)

					_, tx, proxy, err := transparent_upgradeable_proxy.DeployTransparentUpgradeableProxy(
						chain.DeployerKey, chain.Client, implementation.Address(), chain.DeployerKey.From, initData,
					)
					require.NoError(t, err)

					_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
					require.NoError(t, err)

					err = e.Env.ExistingAddresses.Save(chainSelector, proxy.Address().String(), cldf.NewTypeAndVersion(shared.TransparentUpgradeableProxy, deployment.Version1_6_1))
					require.NoError(t, err)

					if tc.wantNewAdmin {
						implementation, err = burn_mint_erc20_transparent.NewBurnMintERC20Transparent(proxy.Address(), chain.Client)
						require.NoError(t, err)

						tx, err = implementation.BeginDefaultAdminTransfer(chain.DeployerKey, config.NewAdmin)
						require.NoError(t, err)

						_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
						require.NoError(t, err)
					}

					r, err := implementation.PendingDefaultAdmin(nil)
					require.NoError(t, err)

					require.Equal(t, config.NewAdmin, r.NewAdmin)
				}
			}
		})
	}
}

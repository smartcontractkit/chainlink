package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestValidateSyncUSDCDomainsWithChainsConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Msg        string
		Input      func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig
		ErrStr     string
		DeployUSDC bool
	}{
		{
			Msg: "Domain mapping not defined",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{}
			},
			ErrStr: "chain selector to usdc domain must be defined",
		},
		{
			Msg: "Chain selector is not valid",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain: map[uint64]semver.Version{
						0: deployment.Version1_5_1,
					},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
				}
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain: map[uint64]semver.Version{
						5009297550715157269: deployment.Version1_5_1,
					},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
				}
			},
			ErrStr: "does not exist in environment",
		},
		{
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain: map[uint64]semver.Version{
						selector: deployment.Version1_5_1,
					},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
				}
			},
			ErrStr: "does not define any USDC token pools, config should be removed",
		},
		{
			Msg: "No USDC token pool found with version",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain: map[uint64]semver.Version{
						selector: deployment.Version1_0_0,
					},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
				}
			},
			DeployUSDC: true,
			ErrStr:     "no USDC token pool found",
		},
		{
			Msg: "Not owned by expected owner",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain: map[uint64]semver.Version{
						selector: deployment.Version1_5_1,
					},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
					MCMS:                      &proposalutils.TimelockConfig{MinDelay: 0 * time.Second},
				}
			},
			DeployUSDC: true,
			ErrStr:     "failed ownership validation",
		},
		{
			Msg: "No domain ID found for selector",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain: map[uint64]semver.Version{
						selector: deployment.Version1_5_1,
					},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
				}
			},
			DeployUSDC: true,
			ErrStr:     "no USDC domain ID defined for chain with selector",
		},
		{
			Msg: "Missing USDC in input",
			Input: func(selector uint64) v1_5_1.SyncUSDCDomainsWithChainsConfig {
				return v1_5_1.SyncUSDCDomainsWithChainsConfig{
					USDCVersionByChain:        map[uint64]semver.Version{},
					ChainSelectorToUSDCDomain: map[uint64]uint32{},
				}
			},
			DeployUSDC: true,
			ErrStr:     "which does support USDC",
		},
	}

	for _, test := range testCases {
		t.Run(test.Msg, func(t *testing.T) {
			if t.Name() == "TestValidateSyncUSDCDomainsWithChainsConfig/Domain_mapping_not_defined" {
				tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-113")
			}
			if t.Name() == "TestValidateSyncUSDCDomainsWithChainsConfig/Chain_selector_is_not_valid" {
				tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-195")
			}
			deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
				testCfg.Chains = 2
				testCfg.PrerequisiteDeploymentOnly = true
				testCfg.IsUSDC = test.DeployUSDC
			})
			e := deployedEnvironment.Env
			selectors := deployedEnvironment.Env.AllChainSelectors()

			if test.DeployUSDC {
				var err error
				e, err = commoncs.Apply(t, e, nil,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
						v1_5_1.ConfigureTokenPoolContractsConfig{
							PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
								selectors[0]: {
									ChainUpdates: v1_5_1.RateLimiterPerChain{
										selectors[1]: testhelpers.CreateSymmetricRateLimits(0, 0),
									},
									Type:    changeset.USDCTokenPool,
									Version: deployment.Version1_5_1,
								},
								selectors[1]: {
									ChainUpdates: v1_5_1.RateLimiterPerChain{
										selectors[0]: testhelpers.CreateSymmetricRateLimits(0, 0),
									},
									Type:    changeset.USDCTokenPool,
									Version: deployment.Version1_5_1,
								},
							},
							TokenSymbol: "USDC",
						},
					),
				)
				require.NoError(t, err)
			}

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			err = test.Input(selectors[0]).Validate(e, state)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestSyncUSDCDomainsWithChainsChangeset(t *testing.T) {
	t.Parallel()

	for _, mcmsConfig := range []*proposalutils.TimelockConfig{nil, {MinDelay: 0 * time.Second}} {
		msg := "Sync domains without MCMS"
		if mcmsConfig != nil {
			msg = "Sync domains with MCMS"
		}

		t.Run(msg, func(t *testing.T) {
			if t.Name() == "TestSyncUSDCDomainsWithChainsChangeset/Sync_domains_without_MCMS" {
				tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-112")
			}
			deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
				testCfg.Chains = 2
				testCfg.PrerequisiteDeploymentOnly = true
				testCfg.IsUSDC = true
			})
			e := deployedEnvironment.Env
			selectors := e.AllChainSelectors()

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(selectors))
			timelockOwnedContractsByChain := make(map[uint64][]common.Address, 1)
			for _, selector := range selectors {
				// Assemble map of addresses required for Timelock scheduling & execution
				timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
					Timelock:  state.Chains[selector].Timelock,
					CallProxy: state.Chains[selector].CallProxy,
				}
				// We would only need the token pool owned by timelock in these tests (if mcms config is provided)
				timelockOwnedContractsByChain[selector] = []common.Address{state.Chains[selector].USDCTokenPools[deployment.Version1_5_1].Address()}
			}

			if mcmsConfig != nil {
				// Transfer ownership of token pools to timelock
				e, err = commoncs.Apply(t, e, timelockContracts,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
						commoncs.TransferToMCMSWithTimelockConfig{
							ContractsByChain: timelockOwnedContractsByChain,
							MCMSConfig:       *mcmsConfig,
						},
					),
				)
				require.NoError(t, err)
			}

			e, err = commoncs.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
					v1_5_1.ConfigureTokenPoolContractsConfig{
						MCMS: mcmsConfig,
						PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
							selectors[0]: {
								ChainUpdates: v1_5_1.RateLimiterPerChain{
									selectors[1]: testhelpers.CreateSymmetricRateLimits(0, 0),
								},
								Type:    changeset.USDCTokenPool,
								Version: deployment.Version1_5_1,
							},
							selectors[1]: {
								ChainUpdates: v1_5_1.RateLimiterPerChain{
									selectors[0]: testhelpers.CreateSymmetricRateLimits(0, 0),
								},
								Type:    changeset.USDCTokenPool,
								Version: deployment.Version1_5_1,
							},
						},
						TokenSymbol: "USDC",
					},
				),
			)
			require.NoError(t, err)

			e, err = commoncs.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.SyncUSDCDomainsWithChainsChangeset),
					v1_5_1.SyncUSDCDomainsWithChainsConfig{
						MCMS: mcmsConfig,
						USDCVersionByChain: map[uint64]semver.Version{
							selectors[0]: deployment.Version1_5_1,
							selectors[1]: deployment.Version1_5_1,
						},
						ChainSelectorToUSDCDomain: map[uint64]uint32{
							selectors[0]: 1,
							selectors[1]: 2,
						},
					},
				),
			)
			require.NoError(t, err)

			state, err = changeset.LoadOnchainState(e)
			require.NoError(t, err)

			for i, selector := range selectors {
				remoteSelector := selectors[0]
				if i == 0 {
					remoteSelector = selectors[1]
				}
				remoteDomain := uint32(1)
				if i == 0 {
					remoteDomain = 2
				}
				usdcTokenPool := state.Chains[selector].USDCTokenPools[deployment.Version1_5_1]
				remoteUsdcTokenPool := state.Chains[remoteSelector].USDCTokenPools[deployment.Version1_5_1]
				domain, err := usdcTokenPool.GetDomain(nil, remoteSelector)
				allowedCaller := make([]byte, 32)
				bytesCopied := copy(allowedCaller, domain.AllowedCaller[:])
				require.Equal(t, 32, bytesCopied)
				require.NoError(t, err)
				require.True(t, domain.Enabled)
				require.Equal(t, remoteDomain, domain.DomainIdentifier)
				require.Equal(t, remoteUsdcTokenPool.Address(), common.BytesToAddress(allowedCaller))
			}

			// Idempotency check
			output, err := v1_5_1.SyncUSDCDomainsWithChainsChangeset(e, v1_5_1.SyncUSDCDomainsWithChainsConfig{
				MCMS: mcmsConfig,
				USDCVersionByChain: map[uint64]semver.Version{
					selectors[0]: deployment.Version1_5_1,
					selectors[1]: deployment.Version1_5_1,
				},
				ChainSelectorToUSDCDomain: map[uint64]uint32{
					selectors[0]: 1,
					selectors[1]: 2,
				},
			})
			require.NoError(t, err)
			require.Empty(t, output.Proposals) //nolint:staticcheck //SA1019 ignoring deprecated field for compatibility; we don't have tools to generate the new field
		})
	}
}

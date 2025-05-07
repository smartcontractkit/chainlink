package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/stretchr/testify/require"

	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	crossfamily "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/crossfamily/v1_6"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestAddEVMSolanaLaneBidirectional(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
			e := tenv.Env
			solChains := tenv.Env.AllChainSelectorsSolana()
			require.NotEmpty(t, solChains)
			evmChains := tenv.Env.AllChainSelectors()
			require.NotEmpty(t, evmChains)
			solChain := solChains[0]
			evmChain := evmChains[0]
			evmState, err := ccipchangeset.LoadOnchainState(e)
			require.NoError(t, err)
			var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
			if tc.mcmsEnabled {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: 1 * time.Second,
					},
					RouterOwnedByTimelock:    true,
					FeeQuoterOwnedByTimelock: true,
					OffRampOwnedByTimelock:   true,
				}
				testhelpers.TransferToTimelock(t, tenv, evmState, []uint64{evmChain})
			}

			// Add EVM and Solana lane
			evmChainState := evmState.Chains[evmChain]
			evmSolanaLaneCSInput := &crossfamily.AddRemoteChainE2EConfig{
				SolanaChainSelector:                  solChain,
				EVMChainSelector:                     evmChain,
				IsTestRouter:                         true,
				EVMOnRampAllowListEnabled:            false,
				EVMFeeQuoterDestChainInput:           v1_6.DefaultFeeQuoterDestChainConfig(true, solChain),
				InitialSolanaGasPriceForEVMFeeQuoter: testhelpers.DefaultGasPrice,
				InitialEVMTokenPricesForEVMFeeQuoter: map[common.Address]*big.Int{
					evmChainState.LinkToken.Address(): testhelpers.DefaultLinkPrice,
					evmChainState.Weth9.Address():     testhelpers.DefaultWethPrice,
				},
				IsRMNVerificationEnabledOnEVMOffRamp: false,
				SolanaRouterConfig: ccipChangesetSolana.RouterConfig{
					RouterDestinationConfig: solRouter.DestChainConfig{
						AllowListEnabled: true,
						AllowedSenders:   []solana.PublicKey{e.SolChains[solChain].DeployerKey.PublicKey()},
					},
				},
				SolanaOffRampConfig: ccipChangesetSolana.OffRampConfig{
					EnabledAsSource: true,
				},
				SolanaFeeQuoterConfig: ccipChangesetSolana.FeeQuoterConfig{
					FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
						IsEnabled:                   true,
						DefaultTxGasLimit:           200000,
						MaxPerMsgGasLimit:           3000000,
						MaxDataBytes:                30000,
						MaxNumberOfTokensPerMsg:     5,
						DefaultTokenDestGasOverhead: 90000,
						DestGasOverhead:             90000,
						// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
						ChainFamilySelector: [4]uint8{40, 18, 213, 44},
					},
				},
				MCMSConfig: mcmsConfig,
			}

			// run the changeset
			e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(crossfamily.AddEVMAndSolanaLaneChangeset, evmSolanaLaneCSInput),
			})
			require.NoError(t, err)

			// Check that the changeset was applied
			evmState, err = ccipchangeset.LoadOnchainState(e)
			require.NoError(t, err)
		})
	}
}

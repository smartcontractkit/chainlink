package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	crossfamily "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/crossfamily/v1_6"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestAddEVMSolanaLaneBidirectional(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
			e := tenv.Env
			solChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
			require.NotEmpty(t, solChains)
			evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
			require.NotEmpty(t, evmChains)
			solChain := solChains[0]
			evmChain1 := evmChains[0]
			evmChain2 := evmChains[1]
			evmState, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)
			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				}
				testhelpers.TransferToTimelock(t, tenv, evmState, []uint64{evmChain1, evmChain2}, false)
			}

			// Add EVM and Solana lane
			evmChainState, _ := evmState.EVMChainState(evmChain1)
			evmChain2State, _ := evmState.EVMChainState(evmChain2)
			feeQCfgSolana := solFeeQuoter.DestChainConfig{
				IsEnabled:                   true,
				DefaultTxGasLimit:           200000,
				MaxPerMsgGasLimit:           3000000,
				MaxDataBytes:                30000,
				MaxNumberOfTokensPerMsg:     5,
				DefaultTokenDestGasOverhead: 90000,
				DestGasOverhead:             90000,
				// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
				ChainFamilySelector: [4]uint8{40, 18, 213, 44},
			}
			feeQCfgEVM := v1_6.DefaultFeeQuoterDestChainConfig(true, solChain)
			evmSolanaLaneCSInput := crossfamily.AddMultiEVMSolanaLaneConfig{
				SolanaChainSelector: solChain,
				MCMSConfig:          mcmsConfig,
				Configs: []crossfamily.AddRemoteChainE2EConfig{
					{
						EVMChainSelector:                     evmChain1,
						IsTestRouter:                         true,
						EVMOnRampAllowListEnabled:            false,
						EVMFeeQuoterDestChainInput:           feeQCfgEVM,
						InitialSolanaGasPriceForEVMFeeQuoter: testhelpers.DefaultGasPrice,
						InitialEVMTokenPricesForEVMFeeQuoter: map[common.Address]*big.Int{
							evmChainState.LinkToken.Address(): testhelpers.DefaultLinkPrice,
							evmChainState.Weth9.Address():     testhelpers.DefaultWethPrice,
						},
						IsRMNVerificationDisabledOnEVMOffRamp: true,
						SolanaRouterConfig: ccipChangesetSolana.RouterConfig{
							RouterDestinationConfig: solRouter.DestChainConfig{
								AllowListEnabled: true,
								AllowedSenders:   []solana.PublicKey{e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()},
							},
						},
						SolanaOffRampConfig: ccipChangesetSolana.OffRampConfig{
							EnabledAsSource: true,
						},
						SolanaFeeQuoterConfig: ccipChangesetSolana.FeeQuoterConfig{
							FeeQuoterDestinationConfig: feeQCfgSolana,
						},
					},
					{
						EVMChainSelector:                     evmChain2,
						IsTestRouter:                         true,
						EVMOnRampAllowListEnabled:            false,
						EVMFeeQuoterDestChainInput:           feeQCfgEVM,
						InitialSolanaGasPriceForEVMFeeQuoter: testhelpers.DefaultGasPrice,
						InitialEVMTokenPricesForEVMFeeQuoter: map[common.Address]*big.Int{
							evmChain2State.LinkToken.Address(): testhelpers.DefaultLinkPrice,
							evmChain2State.Weth9.Address():     testhelpers.DefaultWethPrice,
						},
						IsRMNVerificationDisabledOnEVMOffRamp: true,
						SolanaRouterConfig: ccipChangesetSolana.RouterConfig{
							RouterDestinationConfig: solRouter.DestChainConfig{
								AllowListEnabled: true,
								AllowedSenders:   []solana.PublicKey{e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()},
							},
						},
						SolanaOffRampConfig: ccipChangesetSolana.OffRampConfig{
							EnabledAsSource: true,
						},
						SolanaFeeQuoterConfig: ccipChangesetSolana.FeeQuoterConfig{
							FeeQuoterDestinationConfig: feeQCfgSolana,
						},
					},
				},
			}

			// run the changeset
			out, err := crossfamily.AddEVMAndSolanaLaneChangeset.Apply(e, evmSolanaLaneCSInput)
			require.NoError(t, err)

			// if MCMS is enabled, we need to run the proposal
			if tc.mcmsEnabled {
				for _, prop := range out.MCMSTimelockProposals {
					mcmProp := proposalutils.SignMCMSTimelockProposal(t, e, &prop)
					// return the error so devs can ensure expected reversions
					err = proposalutils.ExecuteMCMSProposalV2(t, e, mcmProp)
					require.NoError(t, err)
					err = proposalutils.ExecuteMCMSTimelockProposalV2(t, e, &prop)
					require.NoError(t, err)
				}
			}

			// Check that the changeset was applied
			evmState, err = stateview.LoadOnchainState(e)
			require.NoError(t, err)

			solanaState, err := stateview.LoadOnchainStateSolana(e)
			require.NoError(t, err)

			// evm changes
			for _, evmChain := range evmChains {
				evmChainState, _ = evmState.EVMChainState(evmChain)

				destCfg, err := evmChainState.OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, solChain)
				require.NoError(t, err)
				require.Equal(t, evmChainState.TestRouter.Address(), destCfg.Router)
				require.False(t, destCfg.AllowlistEnabled)

				srcCfg, err := evmChainState.OffRamp.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, solChain)
				require.NoError(t, err)
				require.Equal(t, evmChainState.TestRouter.Address(), destCfg.Router)
				require.True(t, srcCfg.IsRMNVerificationDisabled)
				require.True(t, srcCfg.IsEnabled)
				expOnRamp, err := evmState.GetOnRampAddressBytes(solChain)
				require.NoError(t, err)
				require.Equal(t, expOnRamp, srcCfg.OnRamp)

				fqDestCfg, err := evmChainState.FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, solChain)
				require.NoError(t, err)
				testhelpers.AssertEqualFeeConfig(t, feeQCfgEVM, fqDestCfg)
			}
			// solana changes
			var offRampSourceChain solOffRamp.SourceChain
			var destChainStateAccount solRouter.DestChain
			var destChainFqAccount solFeeQuoter.DestChain
			var offRampEvmSourceChainPDA solana.PublicKey
			var evmDestChainStatePDA solana.PublicKey
			var fqEvmDestChainPDA solana.PublicKey
			for _, evmChain := range evmChains {
				offRampEvmSourceChainPDA, _, _ = solState.FindOfframpSourceChainPDA(evmChain, solanaState.SolChains[solChain].OffRamp)
				err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
				require.NoError(t, err)
				require.True(t, offRampSourceChain.Config.IsEnabled)

				fqEvmDestChainPDA, _, _ = solState.FindFqDestChainPDA(evmChain, solanaState.SolChains[solChain].FeeQuoter)
				err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
				require.NoError(t, err, "failed to get account info")
				require.Equal(t, solFeeQuoter.TimestampedPackedU224{}, destChainFqAccount.State.UsdPerUnitGas)
				require.True(t, destChainFqAccount.Config.IsEnabled)
				require.Equal(t, feeQCfgSolana, destChainFqAccount.Config)

				evmDestChainStatePDA, _ = solState.FindDestChainStatePDA(evmChain, solanaState.SolChains[solChain].Router)
				err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
				require.NoError(t, err)
				require.NotEmpty(t, destChainStateAccount.Config.AllowedSenders)
				require.True(t, destChainStateAccount.Config.AllowListEnabled)
			}
		})
	}
}

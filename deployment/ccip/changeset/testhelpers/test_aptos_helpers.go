package testhelpers

import (
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/relayer/codec"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func DeployChainContractsToAptosCS(t *testing.T, e DeployedEnv, chainSelector uint64) commonchangeset.ConfiguredChangeSet {
	// Set mock link token address on Address book (to skip deploying)
	err := e.Env.ExistingAddresses.Save(chainSelector, aptoscs.MockLinkAddress, cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
	require.NoError(t, err)

	aptTokenAddress := aptos.AccountAddress{}
	aptTokenAddress.ParseStringRelaxed("0xa")
	ccipConfig := config.DeployAptosChainConfig{
		ContractParamsPerChain: map[uint64]config.ChainContractParams{
			chainSelector: {
				FeeQuoterParams: config.FeeQuoterParams{
					MaxFeeJuelsPerMsg:            new(big.Int).Mul(big.NewInt(100_000_000), big.NewInt(1e18)), // 100M LINK @ 18 decimals
					TokenPriceStalenessThreshold: 24 * 60 * 60,
					FeeTokens:                    []aptos.AccountAddress{aptTokenAddress},
					PremiumMultiplierWeiPerEthByFeeToken: map[shared.TokenSymbol]uint64{
						shared.APTSymbol:  1,
						shared.LinkSymbol: 1,
					},
				},
				OffRampParams: config.OffRampParams{
					ChainSelector:                    chainSelector,
					PermissionlessExecutionThreshold: uint32(globals.PermissionLessExecutionThreshold.Seconds()),
					IsRMNVerificationDisabled:        nil,
					SourceChainSelectors:             nil,
					SourceChainIsEnabled:             nil,
					SourceChainsOnRamp:               nil,
				},
				OnRampParams: config.OnRampParams{
					ChainSelector:  chainSelector,
					AllowlistAdmin: e.Env.BlockChains.AptosChains()[chainSelector].DeployerSigner.AccountAddress(),
					FeeAggregator:  e.Env.BlockChains.AptosChains()[chainSelector].DeployerSigner.AccountAddress(),
				},
			},
		},
		MCMSDeployConfigPerChain: map[uint64]commontypes.MCMSWithTimelockConfigV2{
			chainSelector: {
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(1),
			},
		},
	}

	return commonchangeset.Configure(aptoscs.DeployAptosChain{}, ccipConfig)
}

// GetFungibleAssetBalance queries the given account's balance of the fungible asset in its primary fungible store.
func GetFungibleAssetBalance(
	client aptos.AptosRpcClient,
	account aptos.AccountAddress,
	faMetadataAddress aptos.AccountAddress,
) (uint64, error) {
	bc := bind.NewBoundContract(
		aptos.AccountOne,
		"std",
		"primary_fungible_store",
		client,
	)
	module, function, typeTags, args, err := bc.Encode(
		"balance",
		[]string{
			"0x1::fungible_asset::Metadata",
		},
		[]string{
			"address",
			"address",
		}, []any{
			account,
			faMetadataAddress,
		})
	if err != nil {
		return 0, err
	}
	callData, err := bc.Call(nil, module, function, typeTags, args)
	if err != nil {
		return 0, err
	}

	var balance uint64
	if err := codec.DecodeAptosJsonArray(callData, &balance); err != nil {
		return 0, err
	}
	return balance, nil
}

// MakeBCSEVMExtraArgsV2 makes the BCS encoded extra args for a message sent from an Move based chain that is destined for an EVM chain.
// The extra args are used to specify the gas limit and allow out of order flag for the message.
func MakeBCSEVMExtraArgsV2(gasLimit *big.Int, allowOOO bool) []byte {
	s := &bcs.Serializer{}
	s.U256(*gasLimit)
	s.Bool(allowOOO)
	return append(hexutil.MustDecode(GenericExtraArgsV2Tag), s.ToBytes()...)
}

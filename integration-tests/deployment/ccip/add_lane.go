package ccipdeployment

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var (
	InitialLinkPrice = deployment.E18Mult(20)
	InitialWethPrice = deployment.E18Mult(4000)
	InitialGasPrice  = big.NewInt(2e12)
)

func AddLane(e deployment.Environment, state CCIPOnChainState, from, to uint64) error {
	// TODO: Batch
	tx, err := state.Chains[from].Router.ApplyRampUpdates(e.Chains[from].DeployerKey, []router.RouterOnRamp{
		{
			DestChainSelector: to,
			OnRamp:            state.Chains[from].OnRamp.Address(),
		},
	}, []router.RouterOffRamp{}, []router.RouterOffRamp{})
	if _, err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}
	tx, err = state.Chains[from].OnRamp.ApplyDestChainConfigUpdates(e.Chains[from].DeployerKey,
		[]onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: to,
				Router:            state.Chains[from].Router.Address(),
			},
		})
	if _, err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	_, err = state.Chains[from].FeeQuoter.UpdatePrices(
		e.Chains[from].DeployerKey, fee_quoter.InternalPriceUpdates{
			TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{
				{
					SourceToken: state.Chains[from].LinkToken.Address(),
					UsdPerToken: InitialLinkPrice,
				},
				{
					SourceToken: state.Chains[from].Weth9.Address(),
					UsdPerToken: InitialWethPrice,
				},
			},
			GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
				{
					DestChainSelector: to,
					UsdPerUnitGas:     InitialGasPrice,
				},
			}})
	if _, err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	// Enable dest in fee quoter
	tx, err = state.Chains[from].FeeQuoter.ApplyDestChainConfigUpdates(e.Chains[from].DeployerKey,
		[]fee_quoter.FeeQuoterDestChainConfigArgs{
			{
				DestChainSelector: to,
				DestChainConfig:   defaultFeeQuoterDestChainConfig(),
			},
		})
	if _, err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	tx, err = state.Chains[to].OffRamp.ApplySourceChainConfigUpdates(e.Chains[to].DeployerKey,
		[]offramp.OffRampSourceChainConfigArgs{
			{
				Router:              state.Chains[to].Router.Address(),
				SourceChainSelector: from,
				IsEnabled:           true,
				OnRamp:              common.LeftPadBytes(state.Chains[from].OnRamp.Address().Bytes(), 32),
			},
		})
	if _, err := deployment.ConfirmIfNoError(e.Chains[to], tx, err); err != nil {
		return err
	}
	tx, err = state.Chains[to].Router.ApplyRampUpdates(e.Chains[to].DeployerKey, []router.RouterOnRamp{}, []router.RouterOffRamp{}, []router.RouterOffRamp{
		{
			SourceChainSelector: from,
			OffRamp:             state.Chains[to].OffRamp.Address(),
		},
	})
	_, err = deployment.ConfirmIfNoError(e.Chains[to], tx, err)
	return err
}

func defaultFeeQuoterDestChainConfig() fee_quoter.FeeQuoterDestChainConfig {
	// https://github.com/smartcontractkit/ccip/blob/c4856b64bd766f1ddbaf5d13b42d3c4b12efde3a/contracts/src/v0.8/ccip/libraries/Internal.sol#L337-L337
	/*
		```Solidity
			// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
			bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
		```
	*/
	evmFamilySelector, _ := hex.DecodeString("2812d52c")
	return fee_quoter.FeeQuoterDestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      256,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   50_000,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByte:             10,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       125_000,
		//DefaultTokenDestBytesOverhead:     32,
		DefaultTxGasLimit:      200_000,
		GasMultiplierWeiPerEth: 1,
		NetworkFeeUSDCents:     1,
		ChainFamilySelector:    [4]byte(evmFamilySelector),
	}
}

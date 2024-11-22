package changeset

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var _ deployment.ChangeSet[AddLanesConfig] = AddLanesWithTestRouter

type InitialPrices struct {
	LinkPrice *big.Int // USD to the power of 18 (e18) per LINK
	WethPrice *big.Int // USD to the power of 18 (e18) per WETH
	GasPrice  *big.Int // uint224 packed gas price in USD (112 for exec // 112 for da)
}

type Lane struct {
	From uint64
	To   uint64
}

type AddLanesConfig struct {
	FeeQuoterDestChainConfigArgs map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig
	InitialPricesByChain         map[uint64]InitialPrices
	ChainPairs                   []Lane
}

func (c AddLanesConfig) Validate() error {
	for _, pair := range c.ChainPairs {
		if pair.From == pair.To {
			return fmt.Errorf("cannot add lane to the same chain")
		}
		if _, ok := c.InitialPricesByChain[pair.From]; !ok {
			return fmt.Errorf("missing initial prices for chain %d", pair.From)
		}
		if _, ok := c.FeeQuoterDestChainConfigArgs[pair.From]; !ok {
			// TODO: add more FeeQuoterDestChainConfigArgs validation
			return fmt.Errorf("missing fee quoter dest chain config for chain %d", pair.To)
		} else {
			if _, ok := c.FeeQuoterDestChainConfigArgs[pair.From][pair.To]; !ok {
				return fmt.Errorf("missing fee quoter dest chain config for lane %d->%d", pair.From, pair.To)
			}
		}
	}
	return nil
}

func AddLanesWithTestRouter(e deployment.Environment, cfg AddLanesConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid AddLanesConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err := addLanes(e, cfg)
	if err != nil {
		e.Logger.Errorw("Failed to add lanes", "err", err)
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: newAddresses,
		JobSpecs:    nil,
	}, nil
}

var DefaultInitialPrices = InitialPrices{
	LinkPrice: deployment.E18Mult(20),
	WethPrice: deployment.E18Mult(4000),
	GasPrice:  ToPackedFee(big.NewInt(8e14), big.NewInt(0)),
}

func addLanes(e deployment.Environment, cfg AddLanesConfig) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for _, pair := range cfg.ChainPairs {
		e.Logger.Infow("Enabling lane with test router", "from", pair.From, "to", pair.To)
		if err := AddLane(e, state, pair.From, pair.To, cfg.InitialPricesByChain[pair.From], true, cfg.FeeQuoterDestChainConfigArgs[pair.From][pair.To]); err != nil {
			return err
		}
	}
	return nil
}

func AddLaneWithDefaultPrices(e deployment.Environment, state CCIPOnChainState, from, to uint64, isTestRouter bool) error {
	return AddLane(e, state, from, to, DefaultInitialPrices, isTestRouter, DefaultFeeQuoterDestChainConfig())
}

func AddLane(e deployment.Environment, state CCIPOnChainState, from, to uint64, initialPrices InitialPrices, isTestRouter bool, feeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig) error {
	// TODO: Batch
	var fromRouter *router.Router
	var toRouter *router.Router
	if isTestRouter {
		fromRouter = state.Chains[from].TestRouter
		toRouter = state.Chains[to].TestRouter
	} else {
		fromRouter = state.Chains[from].Router
		toRouter = state.Chains[to].Router
	}
	tx, err := fromRouter.ApplyRampUpdates(e.Chains[from].DeployerKey, []router.RouterOnRamp{
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
				Router:            fromRouter.Address(),
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
					UsdPerToken: initialPrices.LinkPrice,
				},
				{
					SourceToken: state.Chains[from].Weth9.Address(),
					UsdPerToken: initialPrices.WethPrice,
				},
			},
			GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
				{
					DestChainSelector: to,
					UsdPerUnitGas:     initialPrices.GasPrice,
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
				DestChainConfig:   feeQuoterDestChainConfig,
			},
		})
	if _, err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	tx, err = state.Chains[to].OffRamp.ApplySourceChainConfigUpdates(e.Chains[to].DeployerKey,
		[]offramp.OffRampSourceChainConfigArgs{
			{
				Router:              toRouter.Address(),
				SourceChainSelector: from,
				IsEnabled:           true,
				OnRamp:              common.LeftPadBytes(state.Chains[from].OnRamp.Address().Bytes(), 32),
			},
		})
	if _, err := deployment.ConfirmIfNoError(e.Chains[to], tx, err); err != nil {
		return err
	}
	tx, err = toRouter.ApplyRampUpdates(e.Chains[to].DeployerKey, []router.RouterOnRamp{}, []router.RouterOffRamp{}, []router.RouterOffRamp{
		{
			SourceChainSelector: from,
			OffRamp:             state.Chains[to].OffRamp.Address(),
		},
	})
	_, err = deployment.ConfirmIfNoError(e.Chains[to], tx, err)
	return err
}

func DefaultFeeQuoterDestChainConfig() fee_quoter.FeeQuoterDestChainConfig {
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
		DestGasOverhead:                   ccipevm.DestGasOverhead,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByte:             ccipevm.CalldataGasPerByte,
		DestDataAvailabilityOverheadGas:   100,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       125_000,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            11e17, // Gas multiplier in wei per eth is scaled by 1e18, so 11e17 is 1.1 = 110%
		NetworkFeeUSDCents:                1,
		ChainFamilySelector:               [4]byte(evmFamilySelector),
	}
}

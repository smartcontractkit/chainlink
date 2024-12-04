package changeset

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var _ deployment.ChangeSet[AddLanesConfig] = AddLanes

type InitialPrices struct {
	LinkPrice *big.Int // USD to the power of 18 (e18) per LINK
	WethPrice *big.Int // USD to the power of 18 (e18) per WETH
	GasPrice  *big.Int // uint224 packed gas price in USD (112 for exec // 112 for da)
}

func (p InitialPrices) Validate() error {
	if p.LinkPrice == nil {
		return fmt.Errorf("missing link price")
	}
	if p.WethPrice == nil {
		return fmt.Errorf("missing weth price")
	}
	if p.GasPrice == nil {
		return fmt.Errorf("missing gas price")
	}
	return nil
}

type LaneConfig struct {
	SourceSelector        uint64
	DestSelector          uint64
	InitialPricesBySource InitialPrices
	FeeQuoterDestChain    fee_quoter.FeeQuoterDestChainConfig
	TestRouter            bool
}

type AddLanesConfig struct {
	LaneConfigs []LaneConfig
}

func (c AddLanesConfig) Validate() error {
	for _, pair := range c.LaneConfigs {
		if pair.SourceSelector == pair.DestSelector {
			return fmt.Errorf("cannot add lane to the same chain")
		}
		if err := pair.InitialPricesBySource.Validate(); err != nil {
			return fmt.Errorf("error in validating initial prices for chain %d : %w", pair.SourceSelector, err)
		}
		// TODO: add more FeeQuoterDestChainConfigArgs validation
		if pair.FeeQuoterDestChain == (fee_quoter.FeeQuoterDestChainConfig{}) {
			return fmt.Errorf("missing fee quoter dest chain config")
		}
	}
	return nil
}

// AddLanes adds lanes between chains.
// AddLanes is run while the contracts are still owned by the deployer.
// This is useful to test the initial deployment to enable lanes between chains.
// If the testrouter is enabled, the lanes can be used to send messages between chains with testrouter.
// On successful verification with testrouter, the lanes can be enabled with the main router with different addLane ChangeSet.
func AddLanes(e deployment.Environment, cfg AddLanesConfig) (deployment.ChangesetOutput, error) {
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
	for _, laneCfg := range cfg.LaneConfigs {
		e.Logger.Infow("Enabling lane with test router", "from", laneCfg.SourceSelector, "to", laneCfg.DestSelector)
		if err := addLane(e, state, laneCfg, laneCfg.TestRouter); err != nil {
			return err
		}
	}
	return nil
}

func addLane(e deployment.Environment, state CCIPOnChainState, config LaneConfig, isTestRouter bool) error {
	// TODO: Batch
	var fromRouter *router.Router
	var toRouter *router.Router
	from := config.SourceSelector
	to := config.DestSelector
	feeQuoterDestChainConfig := config.FeeQuoterDestChain
	initialPrices := config.InitialPricesBySource
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

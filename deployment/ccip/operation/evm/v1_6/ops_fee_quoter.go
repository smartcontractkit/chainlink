package v1_6

import (
	"errors"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type DeployFeeQInput struct {
	Chain         uint64
	Params        FeeQuoterParams
	LinkAddr      common.Address
	WethAddr      common.Address
	PriceUpdaters []common.Address
}

type ApplyTokenTransferFeeConfigUpdatesConfigPerChain struct {
	TokenTransferFeeConfigs       []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
	TokenTransferFeeConfigsRemove []fee_quoter.FeeQuoterTokenTransferFeeConfigRemoveArgs
}

type ApplyFeeTokensUpdatesInput struct {
	FeeTokensToAdd    []common.Address
	FeeTokensToRemove []common.Address
}

var (
	DeployFeeQuoterOp = opsutil.NewEVMDeployOperation(
		"DeployFeeQuoter",
		semver.MustParse("1.0.0"),
		"Deploys FeeQuoter 1.6.x contract on the specified evm chain",
		shared.FeeQuoter,
		fee_quoter.FeeQuoterMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_6_3, // defaults to v1_6_3, but can be overwritten by input params.FeeQuoterOpts
			EVMBytecode:      common.FromHex(fee_quoter.FeeQuoterBin),
			ZkSyncVMBytecode: fee_quoter.ZkBytecode,
		},
		func(input DeployFeeQInput) []any {
			return []any{
				fee_quoter.FeeQuoterStaticConfig{
					MaxFeeJuelsPerMsg:            input.Params.MaxFeeJuelsPerMsg,
					LinkToken:                    input.LinkAddr,
					TokenPriceStalenessThreshold: input.Params.TokenPriceStalenessThreshold,
				},
				input.PriceUpdaters,
				[]common.Address{input.WethAddr, input.LinkAddr}, // fee tokens
				input.Params.TokenPriceFeedUpdates,
				input.Params.TokenTransferFeeConfigArgs,
				append([]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
					{
						PremiumMultiplierWeiPerEth: input.Params.LinkPremiumMultiplierWeiPerEth,
						Token:                      input.LinkAddr,
					},
					{
						PremiumMultiplierWeiPerEth: input.Params.WethPremiumMultiplierWeiPerEth,
						Token:                      input.WethAddr,
					},
				}, input.Params.MorePremiumMultiplierWeiPerEth...),
				[]fee_quoter.FeeQuoterDestChainConfigArgs{},
			}
		},
	)

	FeeQApplyAuthorizedCallerOp = opsutil.NewEVMCallOperation(
		"FeeQApplyAuthorizedCallerOp",
		semver.MustParse("1.0.0"),
		"Apply authorized caller to FeeQuoter 1.6 contract on the specified evm chain",
		fee_quoter.FeeQuoterABI,
		shared.FeeQuoter,
		fee_quoter.NewFeeQuoter,
		func(feeQuoter *fee_quoter.FeeQuoter, opts *bind.TransactOpts, input fee_quoter.AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
			return feeQuoter.ApplyAuthorizedCallerUpdates(opts, input)
		},
	)

	FeeQuoterApplyDestChainConfigUpdatesOp = opsutil.NewEVMCallOperation(
		"FeeQuoterApplyDestChainConfigUpdatesOp",
		semver.MustParse("1.0.0"),
		"Apply updates to destination chain configs on the FeeQuoter 1.6.0 contract",
		fee_quoter.FeeQuoterABI,
		shared.FeeQuoter,
		fee_quoter.NewFeeQuoter,
		func(feeQuoter *fee_quoter.FeeQuoter, opts *bind.TransactOpts, input []fee_quoter.FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
			return feeQuoter.ApplyDestChainConfigUpdates(opts, input)
		},
	)

	FeeQuoterUpdatePricesOp = opsutil.NewEVMCallOperation(
		"FeeQuoterUpdatePricesOp",
		semver.MustParse("1.0.0"),
		"Update token and gas prices on the FeeQuoter 1.6.0 contract",
		fee_quoter.FeeQuoterABI,
		shared.FeeQuoter,
		fee_quoter.NewFeeQuoter,
		func(feeQuoter *fee_quoter.FeeQuoter, opts *bind.TransactOpts, input fee_quoter.InternalPriceUpdates) (*types.Transaction, error) {
			return feeQuoter.UpdatePrices(opts, input)
		},
	)
	FeeQuoterApplyTokenTransferFeeCfgOp = opsutil.NewEVMCallOperation(
		"FeeQuoterApplyTokenTransferFeeCfgOp",
		semver.MustParse("1.0.0"),
		"Update or Remove token transfer Fee Configs on the FeeQuoter 1.6.0 contract",
		fee_quoter.FeeQuoterABI,
		shared.FeeQuoter,
		fee_quoter.NewFeeQuoter,
		func(feeQuoter *fee_quoter.FeeQuoter, opts *bind.TransactOpts, input ApplyTokenTransferFeeConfigUpdatesConfigPerChain) (*types.Transaction, error) {
			return feeQuoter.ApplyTokenTransferFeeConfigUpdates(opts, input.TokenTransferFeeConfigs, input.TokenTransferFeeConfigsRemove)
		},
	)

	FeeQuoterApplyFeeTokensUpdatesOp = opsutil.NewEVMCallOperation(
		"FeeQuoterApplyFeeTokensUpdatesOp",
		semver.MustParse("1.0.0"),
		"Add or Remove supported fee tokens FeeQuoter 1.6.0 contract",
		fee_quoter.FeeQuoterABI,
		shared.FeeQuoter,
		fee_quoter.NewFeeQuoter,
		func(feeQuoter *fee_quoter.FeeQuoter, opts *bind.TransactOpts, input ApplyFeeTokensUpdatesInput) (*types.Transaction, error) {
			return feeQuoter.ApplyFeeTokensUpdates(opts, input.FeeTokensToRemove, input.FeeTokensToAdd)
		},
	)

	FeeQApplyPremiumMultiplierWeiPerEthUpdateOp = opsutil.NewEVMCallOperation(
		"FeeQApplyPremiumMultiplierWeiPerEthUpdateOp",
		semver.MustParse("1.0.0"),
		"Applies premiumMultiplierWeiPerEth for tokens in FeeQuoter 1.6.0 contract",
		fee_quoter.FeeQuoterABI,
		shared.FeeQuoter,
		fee_quoter.NewFeeQuoter,
		func(feeQuoter *fee_quoter.FeeQuoter, opts *bind.TransactOpts, input []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
			return feeQuoter.ApplyPremiumMultiplierWeiPerEthUpdates(opts, input)
		},
	)
)

type FeeQuoterParams struct {
	MaxFeeJuelsPerMsg              *big.Int
	TokenPriceStalenessThreshold   uint32
	LinkPremiumMultiplierWeiPerEth uint64
	WethPremiumMultiplierWeiPerEth uint64
	MorePremiumMultiplierWeiPerEth []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
	TokenPriceFeedUpdates          []fee_quoter.FeeQuoterTokenPriceFeedUpdate
	TokenTransferFeeConfigArgs     []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
	DestChainConfigArgs            []fee_quoter.FeeQuoterDestChainConfigArgs
}

func (c FeeQuoterParams) Validate() error {
	if c.MaxFeeJuelsPerMsg == nil {
		return errors.New("MaxFeeJuelsPerMsg is nil")
	}
	if c.MaxFeeJuelsPerMsg.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("MaxFeeJuelsPerMsg must be positive")
	}
	if c.TokenPriceStalenessThreshold == 0 {
		return errors.New("TokenPriceStalenessThreshold can't be 0")
	}
	return nil
}

func DefaultFeeQuoterParams() FeeQuoterParams {
	return FeeQuoterParams{
		MaxFeeJuelsPerMsg:              big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
		TokenPriceStalenessThreshold:   uint32(24 * 60 * 60),
		LinkPremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
		WethPremiumMultiplierWeiPerEth: 1e18, // 1.0 ETH
		TokenPriceFeedUpdates:          []fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
		TokenTransferFeeConfigArgs:     []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{},
		MorePremiumMultiplierWeiPerEth: []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{},
		DestChainConfigArgs:            []fee_quoter.FeeQuoterDestChainConfigArgs{},
	}
}

const (
	// https://github.com/smartcontractkit/chainlink/blob/1423e2581e8640d9e5cd06f745c6067bb2893af2/contracts/src/v0.8/ccip/libraries/Internal.sol#L275-L279
	/*
				```Solidity
					// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
					bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
					// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		  		bytes4 public constant CHAIN_FAMILY_SELECTOR_SVM = 0x1e10bdc4;
				```
	*/
	EVMFamilySelector   = "2812d52c"
	SVMFamilySelector   = "1e10bdc4"
	AptosFamilySelector = "ac77ffec"
)

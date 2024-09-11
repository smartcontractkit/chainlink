package feequoter1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
)

type FeeQuoter struct {
	*fee_quoter.FeeQuoter
}

func New(
	fq *fee_quoter.FeeQuoter,
) *FeeQuoter {
	return &FeeQuoter{
		FeeQuoter: fq,
	}
}

func (fq *FeeQuoter) GetStaticConfig(opts *bind.CallOpts) (view.FeeQuoterStaticConfig, error) {
	config, err := fq.FeeQuoter.GetStaticConfig(opts)
	if err != nil {
		return view.FeeQuoterStaticConfig{}, err
	}
	return view.FeeQuoterStaticConfig{
		MaxFeeJuelsPerMsg:  config.MaxFeeJuelsPerMsg.String(),
		LinkToken:          config.LinkToken.Hex(),
		StalenessThreshold: config.StalenessThreshold,
	}, nil
}

func (fq *FeeQuoter) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (view.FeeQuoterDestChainConfig, error) {
	config, err := fq.FeeQuoter.GetDestChainConfig(opts, destChainSelector)
	if err != nil {
		return view.FeeQuoterDestChainConfig{}, err
	}
	return view.FeeQuoterDestChainConfig{
		IsEnabled:                         config.IsEnabled,
		MaxNumberOfTokensPerMsg:           config.MaxNumberOfTokensPerMsg,
		MaxDataBytes:                      config.MaxDataBytes,
		MaxPerMsgGasLimit:                 config.MaxPerMsgGasLimit,
		DestGasOverhead:                   config.DestGasOverhead,
		DestGasPerPayloadByte:             config.DestGasPerPayloadByte,
		DestDataAvailabilityOverheadGas:   config.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    config.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: config.DestDataAvailabilityMultiplierBps,
		DefaultTokenFeeUSDCents:           config.DefaultTokenFeeUSDCents,
		DefaultTokenDestGasOverhead:       config.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:                 config.DefaultTxGasLimit,
		GasMultiplierWeiPerEth:            config.GasMultiplierWeiPerEth,
		NetworkFeeUSDCents:                config.NetworkFeeUSDCents,
		EnforceOutOfOrder:                 config.EnforceOutOfOrder,
		ChainFamilySelector:               fmt.Sprintf("%x", config.ChainFamilySelector),
	}, nil
}

func (fq *FeeQuoter) GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (view.FeeQuoterTokenPriceFeedConfig, error) {
	config, err := fq.FeeQuoter.GetTokenPriceFeedConfig(opts, token)
	if err != nil {
		return view.FeeQuoterTokenPriceFeedConfig{}, err
	}
	return view.FeeQuoterTokenPriceFeedConfig{
		DataFeedAddress: config.DataFeedAddress.Hex(),
		TokenDecimals:   config.TokenDecimals,
	}, nil
}

func (fq *FeeQuoter) GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (view.FeeQuoterTokenTransferFeeConfig, error) {
	config, err := fq.FeeQuoter.GetTokenTransferFeeConfig(opts, destChainSelector, token)
	if err != nil {
		return view.FeeQuoterTokenTransferFeeConfig{}, err
	}
	return view.FeeQuoterTokenTransferFeeConfig{
		MinFeeUSDCents:    config.MinFeeUSDCents,
		MaxFeeUSDCents:    config.MaxFeeUSDCents,
		DeciBps:           config.DeciBps,
		DestGasOverhead:   config.DestGasOverhead,
		DestBytesOverhead: config.DestBytesOverhead,
		IsEnabled:         config.IsEnabled,
	}, nil
}

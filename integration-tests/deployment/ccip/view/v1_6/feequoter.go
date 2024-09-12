package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
)

type FeeQuoter struct {
	view.Contract
	AuthorizedCallers      []string                               `json:"authorizedCallers,omitempty"`
	FeeTokens              []string                               `json:"feeTokens,omitempty"`
	StaticConfig           FeeQuoterStaticConfig                  `json:"staticConfig,omitempty"`
	DestinationChainConfig map[uint64]DestinationConfigWithTokens `json:"DestinationConfigWithTokens,omitempty"` // Map of DestinationChainSelectors
}

type DestinationConfigWithTokens struct {
	DestChainConfig      FeeQuoterDestChainConfig                 `json:"destChainConfig,omitempty"`
	TokenPriceFeedConfig map[string]FeeQuoterTokenPriceFeedConfig `json:"tokenPriceFeedConfig,omitempty"` // Map of Token addresses to FeeQuoterTokenPriceFeedConfig
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg  string `json:"maxFeeJuelsPerMsg,omitempty"`
	LinkToken          string `json:"linkToken,omitempty"`
	StalenessThreshold uint32 `json:"stalenessThreshold,omitempty"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled,omitempty"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg,omitempty"`
	MaxDataBytes                      uint32 `json:"maxDataBytes,omitempty"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit,omitempty"`
	DestGasOverhead                   uint32 `json:"destGasOverhead,omitempty"`
	DestGasPerPayloadByte             uint16 `json:"destGasPerPayloadByte,omitempty"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas,omitempty"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte,omitempty"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps,omitempty"`
	DefaultTokenFeeUSDCents           uint16 `json:"defaultTokenFeeUSDCents,omitempty"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead,omitempty"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit,omitempty"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth,omitempty"`
	NetworkFeeUSDCents                uint32 `json:"networkFeeUSDCents,omitempty"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder,omitempty"`
	ChainFamilySelector               string `json:"chainFamilySelector,omitempty"`
}

type FeeQuoterTokenPriceFeedConfig struct {
	DataFeedAddress string `json:"dataFeedAddress,omitempty"`
	TokenDecimals   uint8  `json:"tokenDecimals,omitempty"`
}

func FeeQuoterSnapshot(fqContract *fee_quoter.FeeQuoter, tokensByDestination map[uint64][]string) (FeeQuoter, error) {
	var fq FeeQuoter
	tv, err := fqContract.TypeAndVersion(nil)
	if err != nil {
		return fq, err
	}
	fq.TypeAndVersion = tv
	fq.Address = fqContract.Address().Hex()
	authorizedCallers, err := fqContract.GetAllAuthorizedCallers(nil)
	if err != nil {
		return fq, err
	}
	fq.AuthorizedCallers = make([]string, 0, len(authorizedCallers))
	for _, ac := range authorizedCallers {
		fq.AuthorizedCallers = append(fq.AuthorizedCallers, ac.Hex())
	}
	feeTokens, err := fqContract.GetFeeTokens(nil)
	if err != nil {
		return fq, err
	}
	fq.FeeTokens = make([]string, 0, len(feeTokens))
	for _, ft := range feeTokens {
		fq.FeeTokens = append(fq.FeeTokens, ft.Hex())
	}
	staticConfig, err := fqContract.GetStaticConfig(nil)
	if err != nil {
		return fq, err
	}
	fq.StaticConfig = FeeQuoterStaticConfig{
		MaxFeeJuelsPerMsg:  staticConfig.MaxFeeJuelsPerMsg.String(),
		LinkToken:          staticConfig.LinkToken.Hex(),
		StalenessThreshold: staticConfig.StalenessThreshold,
	}
	fq.DestinationChainConfig = make(map[uint64]DestinationConfigWithTokens)
	for destChainSelector, tokens := range tokensByDestination {
		destChainConfig, err := fqContract.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return fq, err
		}
		d := FeeQuoterDestChainConfig{
			IsEnabled:                         destChainConfig.IsEnabled,
			MaxNumberOfTokensPerMsg:           destChainConfig.MaxNumberOfTokensPerMsg,
			MaxDataBytes:                      destChainConfig.MaxDataBytes,
			MaxPerMsgGasLimit:                 destChainConfig.MaxPerMsgGasLimit,
			DestGasOverhead:                   destChainConfig.DestGasOverhead,
			DestGasPerPayloadByte:             destChainConfig.DestGasPerPayloadByte,
			DestDataAvailabilityOverheadGas:   destChainConfig.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    destChainConfig.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: destChainConfig.DestDataAvailabilityMultiplierBps,
			DefaultTokenFeeUSDCents:           destChainConfig.DefaultTokenFeeUSDCents,
			DefaultTokenDestGasOverhead:       destChainConfig.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:                 destChainConfig.DefaultTxGasLimit,
			GasMultiplierWeiPerEth:            destChainConfig.GasMultiplierWeiPerEth,
			NetworkFeeUSDCents:                destChainConfig.NetworkFeeUSDCents,
			EnforceOutOfOrder:                 destChainConfig.EnforceOutOfOrder,
			ChainFamilySelector:               fmt.Sprintf("%x", destChainConfig.ChainFamilySelector),
		}
		tokenPriceFeedConfig := make(map[string]FeeQuoterTokenPriceFeedConfig)
		for _, token := range tokens {
			t, err := fqContract.GetTokenPriceFeedConfig(nil, common.HexToAddress(token))
			if err != nil {
				return fq, err
			}
			tokenPriceFeedConfig[token] = FeeQuoterTokenPriceFeedConfig{
				DataFeedAddress: t.DataFeedAddress.Hex(),
				TokenDecimals:   t.TokenDecimals,
			}
		}
		fq.DestinationChainConfig[destChainSelector] = DestinationConfigWithTokens{
			DestChainConfig:      d,
			TokenPriceFeedConfig: tokenPriceFeedConfig,
		}
	}
	return fq, nil
}

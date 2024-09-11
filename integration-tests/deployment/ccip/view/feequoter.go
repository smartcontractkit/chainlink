package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type FeeQuoter struct {
	Contract
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

type FeeQuoterTokenTransferFeeConfig struct {
	MinFeeUSDCents    uint32 `json:"minFeeUSDCents,omitempty"`
	MaxFeeUSDCents    uint32 `json:"maxFeeUSDCents,omitempty"`
	DeciBps           uint16 `json:"deciBps,omitempty"`
	DestGasOverhead   uint32 `json:"destGasOverhead,omitempty"`
	DestBytesOverhead uint32 `json:"destBytesOverhead,omitempty"`
	IsEnabled         bool   `json:"isEnabled,omitempty"`
}

func FeeQuoterSnapshot(reader FeeQuoterReader, tokensByDestination map[uint64][]string) (FeeQuoter, error) {
	var fq FeeQuoter
	tv, err := reader.TypeAndVersion(nil)
	if err != nil {
		return fq, err
	}
	fq.TypeAndVersion = tv
	fq.Address = reader.Address().Hex()
	authorizedCallers, err := reader.GetAllAuthorizedCallers(nil)
	if err != nil {
		return fq, err
	}
	fq.AuthorizedCallers = make([]string, 0, len(authorizedCallers))
	for _, ac := range authorizedCallers {
		fq.AuthorizedCallers = append(fq.AuthorizedCallers, ac.Hex())
	}
	feeTokens, err := reader.GetFeeTokens(nil)
	if err != nil {
		return fq, err
	}
	fq.FeeTokens = make([]string, 0, len(feeTokens))
	for _, ft := range feeTokens {
		fq.FeeTokens = append(fq.FeeTokens, ft.Hex())
	}
	staticConfig, err := reader.GetStaticConfig(nil)
	if err != nil {
		return fq, err
	}
	fq.StaticConfig = staticConfig
	fq.DestinationChainConfig = make(map[uint64]DestinationConfigWithTokens)
	for destChainSelector, tokens := range tokensByDestination {
		destChainConfig, err := reader.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return fq, err
		}
		tokenPriceFeedConfig := make(map[string]FeeQuoterTokenPriceFeedConfig)
		for _, token := range tokens {
			t, err := reader.GetTokenPriceFeedConfig(nil, common.HexToAddress(token))
			if err != nil {
				return fq, err
			}
			tokenPriceFeedConfig[token] = t
		}
		fq.DestinationChainConfig[destChainSelector] = DestinationConfigWithTokens{
			DestChainConfig:      destChainConfig,
			TokenPriceFeedConfig: tokenPriceFeedConfig,
		}
	}
	return fq, nil
}

type FeeQuoterReader interface {
	ContractState
	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)
	GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error)
	GetStaticConfig(opts *bind.CallOpts) (FeeQuoterStaticConfig, error)
	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (FeeQuoterDestChainConfig, error)
	GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (FeeQuoterTokenPriceFeedConfig, error)
	GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error)
}

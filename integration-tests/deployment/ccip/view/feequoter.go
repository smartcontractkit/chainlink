package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type FeeQuoter struct {
	Contract
	AuthorizedCallers      []string                            `json:"authorizedCallers,omitempty"`
	FeeTokens              []string                            `json:"feeTokens,omitempty"`
	StaticConfig           FeeQuoterStaticConfig               `json:"staticConfig,omitempty"`
	DestinationChainConfig map[uint64]FeeQuoterDestChainConfig `json:"destinationChainConfig,omitempty"` // Map of DestinationChainSelectors to FeeQuoterDestChainConfig
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg  string         `json:"maxFeeJuelsPerMsg,omitempty"`
	LinkToken          common.Address `json:"linkToken,omitempty"`
	StalenessThreshold uint32         `json:"stalenessThreshold,omitempty"`
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

type TokenTransferFeeConfig struct {
	MinFeeUSDCents    uint32 `json:"minFeeUSDCents,omitempty"`
	MaxFeeUSDCents    uint32 `json:"maxFeeUSDCents,omitempty"`
	DeciBps           uint16 `json:"deciBps,omitempty"`
	DestGasOverhead   uint32 `json:"destGasOverhead,omitempty"`
	DestBytesOverhead uint32 `json:"destBytesOverhead,omitempty"`
	IsEnabled         bool   `json:"isEnabled,omitempty"`
}

type FeeQuoterReader interface {
	ContractState
	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)
	GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error)
	GetStaticConfig(opts *bind.CallOpts) (FeeQuoterStaticConfig, error)
	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (FeeQuoterDestChainConfig, error)
	GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (FeeQuoterTokenPriceFeedConfig, error)
	GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (TokenTransferFeeConfig, error)
}

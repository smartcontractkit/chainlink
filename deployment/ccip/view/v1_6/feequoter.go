package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	router1_2 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type FeeQuoterView struct {
	types.ContractMetaData
	AuthorizedCallers                       []string                                 `json:"authorizedCallers,omitempty"`
	FeeTokens                               []string                                 `json:"feeTokens,omitempty"`
	StaticConfig                            FeeQuoterStaticConfig                    `json:"staticConfig,omitempty"`
	DestinationChainConfigBasedOnTestRouter map[uint64]FeeQuoterDestChainConfig      `json:"destinationChainConfigBasedOnTestRouter,omitempty"`
	DestinationChainConfig                  map[uint64]FeeQuoterDestChainConfig      `json:"destinationChainConfig,omitempty"`
	TokenPriceFeedConfig                    map[string]FeeQuoterTokenPriceFeedConfig `json:"tokenPriceFeedConfig,omitempty"`
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            string `json:"maxFeeJuelsPerMsg,omitempty"`
	LinkToken                    string `json:"linkToken,omitempty"`
	TokenPriceStalenessThreshold uint32 `json:"tokenPriceStalenessThreshold,omitempty"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg,omitempty"`
	MaxDataBytes                      uint32 `json:"maxDataBytes,omitempty"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit,omitempty"`
	DestGasOverhead                   uint32 `json:"destGasOverhead,omitempty"`
	DestGasPerPayloadByteBase         uint8  `json:"destGasPerPayloadByteBase,omitempty"`
	DestGasPerPayloadByteHigh         uint8  `json:"destGasPerPayloadByteHigh,omitempty"`
	DestGasPerPayloadByteThreshold    uint16 `json:"destGasPerPayloadByteThreshold,omitempty"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas,omitempty"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte,omitempty"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps,omitempty"`
	DefaultTokenFeeUSDCents           uint16 `json:"defaultTokenFeeUSDCents,omitempty"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead,omitempty"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit,omitempty"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth,omitempty"`
	NetworkFeeUSDCents                uint32 `json:"networkFeeUSDCents,omitempty"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder"`
	ChainFamilySelector               string `json:"chainFamilySelector,omitempty"`
}

type FeeQuoterTokenPriceFeedConfig struct {
	DataFeedAddress string `json:"dataFeedAddress,omitempty"`
	TokenDecimals   uint8  `json:"tokenDecimals,omitempty"`
}

func GenerateFeeQuoterView(fqContract *fee_quoter.FeeQuoter, router, testRouter *router1_2.Router, tokens []common.Address) (FeeQuoterView, error) {
	fq := FeeQuoterView{}
	authorizedCallers, err := fqContract.GetAllAuthorizedCallers(nil)
	if err != nil {
		return FeeQuoterView{}, err
	}
	fq.AuthorizedCallers = make([]string, 0, len(authorizedCallers))
	for _, ac := range authorizedCallers {
		fq.AuthorizedCallers = append(fq.AuthorizedCallers, ac.Hex())
	}
	fq.ContractMetaData, err = types.NewContractMetaData(fqContract, fqContract.Address())
	if err != nil {
		return FeeQuoterView{}, fmt.Errorf("metadata error for FeeQuoter: %w", err)
	}
	feeTokens, err := fqContract.GetFeeTokens(nil)
	if err != nil {
		return FeeQuoterView{}, err
	}
	fq.FeeTokens = make([]string, 0, len(feeTokens))
	for _, ft := range feeTokens {
		fq.FeeTokens = append(fq.FeeTokens, ft.Hex())
	}
	staticConfig, err := fqContract.GetStaticConfig(nil)
	if err != nil {
		return FeeQuoterView{}, err
	}
	fq.StaticConfig = FeeQuoterStaticConfig{
		MaxFeeJuelsPerMsg:            staticConfig.MaxFeeJuelsPerMsg.String(),
		LinkToken:                    staticConfig.LinkToken.Hex(),
		TokenPriceStalenessThreshold: staticConfig.TokenPriceStalenessThreshold,
	}
	// find router contract in dependencies
	fq.DestinationChainConfig, err = generateDestChainConfig(fqContract, router)
	if err != nil {
		return FeeQuoterView{}, err
	}
	if testRouter != nil {
		fq.DestinationChainConfigBasedOnTestRouter, err = generateDestChainConfig(fqContract, testRouter)
		if err != nil {
			return FeeQuoterView{}, err
		}
	}
	fq.TokenPriceFeedConfig = make(map[string]FeeQuoterTokenPriceFeedConfig)
	for _, token := range tokens {
		t, err := fqContract.GetTokenPriceFeedConfig(nil, token)
		if err != nil {
			return FeeQuoterView{}, err
		}
		fq.TokenPriceFeedConfig[token.String()] = FeeQuoterTokenPriceFeedConfig{
			DataFeedAddress: t.DataFeedAddress.Hex(),
			TokenDecimals:   t.TokenDecimals,
		}
	}
	return fq, nil
}

func generateDestChainConfig(fqContract *fee_quoter.FeeQuoter, router *router1_2.Router) (map[uint64]FeeQuoterDestChainConfig, error) {
	destinationChainConfig := make(map[uint64]FeeQuoterDestChainConfig)
	destSelectors, err := v1_2.GetRemoteChainSelectors(router)
	if err != nil {
		return nil, fmt.Errorf("view error for FeeQuoter while getting remote chain selectors from router: %w", err)
	}
	for _, destChainSelector := range destSelectors {
		destChainConfig, err := fqContract.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return nil, err
		}
		destinationChainConfig[destChainSelector] = FeeQuoterDestChainConfig{
			IsEnabled:                         destChainConfig.IsEnabled,
			MaxNumberOfTokensPerMsg:           destChainConfig.MaxNumberOfTokensPerMsg,
			MaxDataBytes:                      destChainConfig.MaxDataBytes,
			MaxPerMsgGasLimit:                 destChainConfig.MaxPerMsgGasLimit,
			DestGasOverhead:                   destChainConfig.DestGasOverhead,
			DestGasPerPayloadByteBase:         destChainConfig.DestGasPerPayloadByteBase,
			DestGasPerPayloadByteHigh:         destChainConfig.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    destChainConfig.DestGasPerPayloadByteThreshold,
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
	}
	return destinationChainConfig, nil
}

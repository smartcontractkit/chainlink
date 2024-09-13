package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	router1_2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
)

type FeeQuoter struct {
	types.ContractMetaData
	AuthorizedCallers      []string                                 `json:"authorizedCallers,omitempty"`
	FeeTokens              []string                                 `json:"feeTokens,omitempty"`
	StaticConfig           FeeQuoterStaticConfig                    `json:"staticConfig,omitempty"`
	DestinationChainConfig map[uint64]FeeQuoterDestChainConfig      `json:"destinationChainConfig,omitempty"`
	TokenPriceFeedConfig   map[string]FeeQuoterTokenPriceFeedConfig `json:"tokenPriceFeedConfig,omitempty"`
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

func (fq *FeeQuoter) Snapshot(fqContractMeta types.ContractMetaData, dependenciesMeta []types.ContractMetaData, client bind.ContractBackend) error {
	if err := fqContractMeta.Validate(); err != nil {
		return fmt.Errorf("snapshot error for FeeQuoter: %w", err)
	}
	fq.ContractMetaData = fqContractMeta
	fqContract, err := fee_quoter.NewFeeQuoter(fqContractMeta.Address, client)
	if err != nil {
		return fmt.Errorf("failed to get fee quoter contract: %w", err)
	}
	authorizedCallers, err := fqContract.GetAllAuthorizedCallers(nil)
	if err != nil {
		return err
	}
	fq.AuthorizedCallers = make([]string, 0, len(authorizedCallers))
	for _, ac := range authorizedCallers {
		fq.AuthorizedCallers = append(fq.AuthorizedCallers, ac.Hex())
	}
	feeTokens, err := fqContract.GetFeeTokens(nil)
	if err != nil {
		return err
	}
	fq.FeeTokens = make([]string, 0, len(feeTokens))
	for _, ft := range feeTokens {
		fq.FeeTokens = append(fq.FeeTokens, ft.Hex())
	}
	staticConfig, err := fqContract.GetStaticConfig(nil)
	if err != nil {
		return err
	}
	fq.StaticConfig = FeeQuoterStaticConfig{
		MaxFeeJuelsPerMsg:  staticConfig.MaxFeeJuelsPerMsg.String(),
		LinkToken:          staticConfig.LinkToken.Hex(),
		StalenessThreshold: staticConfig.StalenessThreshold,
	}
	// find router contract in dependencies
	fq.DestinationChainConfig = make(map[uint64]FeeQuoterDestChainConfig)
	destSelectors, err := GetDestinationSelectors(dependenciesMeta, client)
	if err != nil {
		return fmt.Errorf("snapshot error for FeeQuoter: %w", err)
	}
	for _, destChainSelector := range destSelectors {
		destChainConfig, err := fqContract.GetDestChainConfig(nil, destChainSelector)
		if err != nil {
			return err
		}
		fq.DestinationChainConfig[destChainSelector] = FeeQuoterDestChainConfig{
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
	}
	fq.TokenPriceFeedConfig = make(map[string]FeeQuoterTokenPriceFeedConfig)
	tokens, err := GetSupportedTokens(dependenciesMeta, client)
	if err != nil {
		return fmt.Errorf("snapshot error for FeeQuoter: %w", err)
	}
	for _, token := range tokens {
		t, err := fqContract.GetTokenPriceFeedConfig(nil, token)
		if err != nil {
			return err
		}
		fq.TokenPriceFeedConfig[token.String()] = FeeQuoterTokenPriceFeedConfig{
			DataFeedAddress: t.DataFeedAddress.Hex(),
			TokenDecimals:   t.TokenDecimals,
		}
	}
	return nil
}

func GetSupportedTokens(dependenciesMeta []types.ContractMetaData, client bind.ContractBackend) ([]common.Address, error) {
	for _, dep := range dependenciesMeta {
		if dep.TypeAndVersion == types.TokenAdminRegistryTypeAndVersionV1_5 {
			taContract, err := token_admin_registry.NewTokenAdminRegistry(dep.Address, client)
			if err != nil {
				return nil, fmt.Errorf("failed to get router contract: %w", err)
			}
			// TODO : include pagination CCIP-3416
			tokens, err := taContract.GetAllConfiguredTokens(nil, 0, 10)
			if err != nil {
				return nil, fmt.Errorf("failed to get tokens from router: %w", err)
			}
			return tokens, nil
		}
	}
	return nil, fmt.Errorf("token admin registry not found in dependencies")
}

func GetDestinationSelectors(dependenciesMeta []types.ContractMetaData, client bind.ContractBackend) ([]uint64, error) {
	destSelectors := make([]uint64, 0)
	foundRouter := false
	for _, dep := range dependenciesMeta {
		if dep.TypeAndVersion == types.RouterTypeAndVersionV1_2 {
			foundRouter = true
			routerContract, err := router1_2.NewRouter(dep.Address, client)
			if err != nil {
				return nil, fmt.Errorf("failed to get router contract: %w", err)
			}
			offRamps, err := routerContract.GetOffRamps(nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get offRamps from router: %w", err)
			}
			for _, offRamp := range offRamps {
				destSelectors = append(destSelectors, offRamp.SourceChainSelector)
			}
		}
	}
	if !foundRouter {
		return nil, fmt.Errorf("router not found in dependencies")
	}
	return destSelectors, nil
}

package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink/deployment"
)

type FeeQuoterView struct {
	Version                uint8                                              `json:"version,omitempty"`
	Owner                  string                                             `json:"owner,omitempty"`
	ProposedOwner          string                                             `json:"proposedOwner,omitempty"`
	MaxFeeJuelsPerMsg      string                                             `json:"maxFeeJuelsPerMsg,omitempty"`
	LinkTokenMint          string                                             `json:"linkTokenMint,omitempty"`
	LinkTokenLocalDecimals uint8                                              `json:"linkTokenLocalDecimals,omitempty"`
	Onramp                 string                                             `json:"onRamp,omitempty"`
	DefaultCodeVersion     string                                             `json:"defaultCodeVersion,omitempty"`
	DestinationChainConfig map[uint64]FeeQuoterDestChainConfig                `json:"destinationChainConfig,omitempty"`
	BillingTokenConfig     map[string]FeeQuoterBillingTokenConfig             `json:"billingTokenConfig,omitempty"`
	TokenTransferConfig    map[uint64]map[string]FeeQuoterTokenTransferConfig `json:"tokenTransferConfig,omitempty"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled,omitempty"`
	LaneCodeVersion                   string `json:"laneCodeVersion,omitempty"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg,omitempty"`
	MaxDataBytes                      uint32 `json:"maxDataBytes,omitempty"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit,omitempty"`
	DestGasOverhead                   uint32 `json:"destGasOverhead,omitempty"`
	DestGasPerPayloadByteBase         uint32 `json:"destGasPerPayloadByteBase,omitempty"`
	DestGasPerPayloadByteHigh         uint32 `json:"destGasPerPayloadByteHigh,omitempty"`
	DestGasPerPayloadByteThreshold    uint32 `json:"destGasPerPayloadByteThreshold,omitempty"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas,omitempty"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte,omitempty"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps,omitempty"`
	DefaultTokenFeeUSDCents           uint16 `json:"defaultTokenFeeUSDCents,omitempty"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead,omitempty"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit,omitempty"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth,omitempty"`
	NetworkFeeUSDCents                uint32 `json:"networkFeeUSDCents,omitempty"`
	GasPriceStalenessThreshold        uint32 `json:"gasPriceStalenessThreshold,omitempty"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder,omitempty"`
	ChainFamilySelector               string `json:"chainFamilySelector,omitempty"`
}

type FeeQuoterBillingTokenConfig struct {
	Enabled                    bool   `json:"enabled,omitempty"`
	PremiumMultiplierWeiPerEth uint64 `json:"premiumMultiplierWeiPerEth,omitempty"`
}

type FeeQuoterTokenTransferConfig struct {
	MinFeeUsdcents    uint32 `json:"minFeeUsdcents,omitempty"`
	MaxFeeUsdcents    uint32 `json:"maxFeeUsdcents,omitempty"`
	DeciBps           uint16 `json:"deciBps,omitempty"`
	DestGasOverhead   uint32 `json:"destGasOverhead,omitempty"`
	DestBytesOverhead uint32 `json:"destBytesOverhead,omitempty"`
	IsEnabled         bool   `json:"isEnabled,omitempty"`
}

func GenerateFeeQuoterView(chain deployment.SolChain, program solana.PublicKey, remoteChains []uint64, tokens []solana.PublicKey) (FeeQuoterView, error) {
	fq := FeeQuoterView{}
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(program)
	err := chain.GetAccountDataBorshInto(context.Background(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		return fq, fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	fq.Version = fqConfig.Version
	fq.Owner = fqConfig.Owner.String()
	fq.ProposedOwner = fqConfig.ProposedOwner.String()
	fq.MaxFeeJuelsPerMsg = fqConfig.MaxFeeJuelsPerMsg.String()
	fq.LinkTokenMint = fqConfig.LinkTokenMint.String()
	fq.LinkTokenLocalDecimals = fqConfig.LinkTokenLocalDecimals
	fq.Onramp = fqConfig.Onramp.String()
	fq.DefaultCodeVersion = fqConfig.DefaultCodeVersion.String()
	fq.BillingTokenConfig = make(map[string]FeeQuoterBillingTokenConfig)
	fq.TokenTransferConfig = make(map[uint64]map[string]FeeQuoterTokenTransferConfig)
	fq.DestinationChainConfig = make(map[uint64]FeeQuoterDestChainConfig)
	for _, remote := range remoteChains {
		fqRemoteChainPDA, _, err := solState.FindFqDestChainPDA(remote, program)
		if err != nil {
			return fq, fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solFeeQuoter.DestChain
		if err = chain.GetAccountDataBorshInto(context.Background(), fqRemoteChainPDA, &destChainStateAccount); err != nil {
			return fq, fmt.Errorf("remote %d is not configured on solana chain feequoter %d", remote, chain.Selector)
		}
		fq.TokenTransferConfig[remote] = make(map[string]FeeQuoterTokenTransferConfig)
		fq.DestinationChainConfig[remote] = FeeQuoterDestChainConfig{
			IsEnabled:                         destChainStateAccount.Config.IsEnabled,
			LaneCodeVersion:                   destChainStateAccount.Config.LaneCodeVersion.String(),
			MaxNumberOfTokensPerMsg:           destChainStateAccount.Config.MaxNumberOfTokensPerMsg,
			MaxDataBytes:                      destChainStateAccount.Config.MaxDataBytes,
			MaxPerMsgGasLimit:                 destChainStateAccount.Config.MaxPerMsgGasLimit,
			DestGasOverhead:                   destChainStateAccount.Config.DestGasOverhead,
			DestGasPerPayloadByteBase:         destChainStateAccount.Config.DestGasPerPayloadByteBase,
			DestGasPerPayloadByteHigh:         destChainStateAccount.Config.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    destChainStateAccount.Config.DestGasPerPayloadByteThreshold,
			DestDataAvailabilityOverheadGas:   destChainStateAccount.Config.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    destChainStateAccount.Config.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: destChainStateAccount.Config.DestDataAvailabilityMultiplierBps,
			DefaultTokenFeeUSDCents:           destChainStateAccount.Config.DefaultTokenFeeUsdcents,
			DefaultTokenDestGasOverhead:       destChainStateAccount.Config.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:                 destChainStateAccount.Config.DefaultTxGasLimit,
			GasMultiplierWeiPerEth:            destChainStateAccount.Config.GasMultiplierWeiPerEth,
			NetworkFeeUSDCents:                destChainStateAccount.Config.NetworkFeeUsdcents,
			GasPriceStalenessThreshold:        destChainStateAccount.Config.GasPriceStalenessThreshold,
			EnforceOutOfOrder:                 destChainStateAccount.Config.EnforceOutOfOrder,
			ChainFamilySelector:               fmt.Sprintf("%x", destChainStateAccount.Config.ChainFamilySelector),
		}
		// TODO: save the configured chains/tokens to the AB so we can reconstruct state without the loop
		for _, token := range tokens {
			remoteBillingPDA, _, err := solState.FindFqPerChainPerTokenConfigPDA(remote, token, program)
			if err != nil {
				return fq, fmt.Errorf("failed to find remote billing token config pda for (remoteSelector: %d, mint: %s, feeQuoter: %s): %w", remote, token, program, err)
			}
			var remoteBillingAccount solFeeQuoter.PerChainPerTokenConfig
			if err := chain.GetAccountDataBorshInto(context.Background(), remoteBillingPDA, &remoteBillingAccount); err == nil {
				fq.TokenTransferConfig[remote][token.String()] = FeeQuoterTokenTransferConfig{
					MinFeeUsdcents:    remoteBillingAccount.TokenTransferConfig.MinFeeUsdcents,
					MaxFeeUsdcents:    remoteBillingAccount.TokenTransferConfig.MaxFeeUsdcents,
					DeciBps:           remoteBillingAccount.TokenTransferConfig.DeciBps,
					DestGasOverhead:   remoteBillingAccount.TokenTransferConfig.DestGasOverhead,
					DestBytesOverhead: remoteBillingAccount.TokenTransferConfig.DestBytesOverhead,
					IsEnabled:         remoteBillingAccount.TokenTransferConfig.IsEnabled,
				}
			}
		}
	}
	// TODO: save the configured chains/tokens to the AB so we can reconstruct state without the loop
	for _, token := range tokens {
		billingConfigPDA, _, err := solState.FindFqBillingTokenConfigPDA(token, program)
		if err != nil {
			return fq, fmt.Errorf("failed to find billing token config pda (mint: %s, feeQuoter: %s): %w", token.String(), program.String(), err)
		}
		var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
		if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount); err == nil {
			fq.BillingTokenConfig[token.String()] = FeeQuoterBillingTokenConfig{
				Enabled:                    token0ConfigAccount.Config.Enabled,
				PremiumMultiplierWeiPerEth: token0ConfigAccount.Config.PremiumMultiplierWeiPerEth,
			}
		}
	}

	return fq, nil
}

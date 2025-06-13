package aptos

import (
	"encoding/hex"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type CCIPView struct {
	aptosCommon.ContractMetaData

	FeeQuoter FeeQuoterView `json:"feeQuoter,omitempty"`
	RMNRemote RMNRemoteView `json:"rmnRemote,omitempty"`
}

type FeeQuoterView struct {
	FeeTokens               []string
	StaticConfig            FeeQuoterStaticConfig
	DestinationChainConfigs map[uint64]FeeQuoterDestChainConfig
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            string
	LinkToken                    string
	TokenPriceStalenessThreshold uint64
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool
	MaxNumberOfTokensPerMsg           uint16
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
	DestGasOverhead                   uint32
	DestGasPerPayloadByteBase         uint8
	DestGasPerPayloadByteHigh         uint8
	DestGasPerPayloadByteThreshold    uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	ChainFamilySelector               string
	EnforceOutOfOrder                 bool
	DefaultTokenFeeUsdCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	GasPriceStalenessThreshold        uint32
	NetworkFeeUsdCents                uint32
}

type RMNRemoteView struct {
	aptosCommon.ContractMetaData
	IsCursed             bool
	Config               RMNRemoteVersionedConfig
	CursedSubjectEntries []RMNRemoteCurseEntry
}

type RMNRemoteVersionedConfig struct {
	Version uint32
	Signers []RMNRemoteSigner
	Fsign   uint64 `json:"fSign"`
}

type RMNRemoteSigner struct {
	OnchainPublicKey string `json:"onchain_public_key"` // Follow EVM snake_case
	NodeIndex        uint64 `json:"node_index"`
}

type RMNRemoteCurseEntry struct {
	Subject  string
	Selector uint64
}

func GenerateCCIPView(chain cldf_aptos.Chain, ccipAddress aptos.AccountAddress, routerAddress aptos.AccountAddress) (CCIPView, error) {
	boundCCIP := ccip.Bind(ccipAddress, chain.Client)
	boundRouter := ccip_router.Bind(ccipAddress, chain.Client)

	ccipOwner, err := boundCCIP.Auth().Owner(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get ccip owner: %w", err)
	}

	// Router
	destChainSelectors, err := boundRouter.Router().GetDestChains(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get dest chain selectors: %w", err)
	}

	// FeeQuoter
	feeQuoterStaticConfig, err := boundCCIP.FeeQuoter().GetStaticConfig(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get feeQuoter static config: %w", err)
	}
	feeQuoterFeeTkns, err := boundCCIP.FeeQuoter().GetFeeTokens(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get feeQuoter feeTokens: %w", err)
	}
	feeQuoterFeeTokens := make([]string, len(feeQuoterFeeTkns))
	for i, feeToken := range feeQuoterFeeTkns {
		feeQuoterFeeTokens[i] = feeToken.StringLong()
	}
	destinationChainConfigs := make(map[uint64]FeeQuoterDestChainConfig, len(destChainSelectors))
	for _, selector := range destChainSelectors {
		destChainConfig, err := boundCCIP.FeeQuoter().GetDestChainConfig(nil, selector)
		if err != nil {
			return CCIPView{}, fmt.Errorf("failed to get dest chain config for chain %v: %w", selector, err)
		}
		destinationChainConfigs[selector] = FeeQuoterDestChainConfig{
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
			ChainFamilySelector:               fmt.Sprintf("%x", destChainConfig.ChainFamilySelector),
			EnforceOutOfOrder:                 destChainConfig.EnforceOutOfOrder,
			DefaultTokenFeeUsdCents:           destChainConfig.DefaultTokenFeeUsdCents,
			DefaultTokenDestGasOverhead:       destChainConfig.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:                 destChainConfig.DefaultTxGasLimit,
			GasMultiplierWeiPerEth:            destChainConfig.GasMultiplierWeiPerEth,
			GasPriceStalenessThreshold:        destChainConfig.GasPriceStalenessThreshold,
			NetworkFeeUsdCents:                destChainConfig.NetworkFeeUsdCents,
		}
	}

	// RMNRemote
	cursedSubjects, err := boundCCIP.RMNRemote().GetCursedSubjects(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get RMNRemote cursed subjects: %w", err)
	}
	cursedSubjectEntries := make([]RMNRemoteCurseEntry, len(cursedSubjects))
	for i, subject := range cursedSubjects {
		cursedSubjectEntries[i] = RMNRemoteCurseEntry{
			Subject:  hex.EncodeToString(subject),
			Selector: globals.FamilyAwareSubjectToSelector(globals.Subject(subject), chainselectors.FamilyAptos),
		}
	}
	version, rmnRemoteConfig, err := boundCCIP.RMNRemote().GetVersionedConfig(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get RMNRemote versioned config: %w", err)
	}
	rmnRemoteVersionedConfig := RMNRemoteVersionedConfig{
		Version: version,
		Signers: make([]RMNRemoteSigner, len(rmnRemoteConfig.Signers)),
		Fsign:   rmnRemoteConfig.FSign,
	}
	for i, signer := range rmnRemoteConfig.Signers {
		rmnRemoteVersionedConfig.Signers[i] = RMNRemoteSigner{
			OnchainPublicKey: string(signer.OnchainPublicKey),
			NodeIndex:        signer.NodeIndex,
		}
	}

	return CCIPView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address: ccipAddress.StringLong(),
			Owner:   ccipOwner.StringLong(),
		},
		FeeQuoter: FeeQuoterView{
			FeeTokens: feeQuoterFeeTokens,
			StaticConfig: FeeQuoterStaticConfig{
				MaxFeeJuelsPerMsg:            feeQuoterStaticConfig.MaxFeeJuelsPerMsg.String(),
				LinkToken:                    feeQuoterStaticConfig.LinkToken.StringLong(),
				TokenPriceStalenessThreshold: feeQuoterStaticConfig.TokenPriceStalenessThreshold,
			},
			DestinationChainConfigs: destinationChainConfigs,
		},
		RMNRemote: RMNRemoteView{
			ContractMetaData: aptosCommon.ContractMetaData{
				Address: ccipAddress.StringLong(),
				Owner:   ccipOwner.StringLong(),
			},
			IsCursed:             len(cursedSubjectEntries) != 0,
			Config:               rmnRemoteVersionedConfig,
			CursedSubjectEntries: cursedSubjectEntries,
		},
	}, nil
}

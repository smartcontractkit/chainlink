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

	FeeQuoter          FeeQuoterView          `json:"feeQuoter,omitempty"`
	RMNRemote          RMNRemoteView          `json:"rmnRemote,omitempty"`
	TokenAdminRegistry TokenAdminRegistryView `json:"tokenAdminRegistry,omitempty"`
	NonceManager       NonceManagerView       `json:"nonceManager,omitempty"`
	ReceiverRegistry   ReceiverRegistryView   `json:"receiverRegistry,omitempty"`
}

type FeeQuoterView struct {
	aptosCommon.ContractMetaData

	FeeTokens               []string                            `json:"feeTokens"`
	StaticConfig            FeeQuoterStaticConfig               `json:"staticConfig"`
	DestinationChainConfigs map[uint64]FeeQuoterDestChainConfig `json:"destinationChainConfigs"`
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            string `json:"maxFeeJuelsPerMsg"`
	LinkToken                    string `json:"linkToken"`
	TokenPriceStalenessThreshold uint64 `json:"tokenPriceStalenessThreshold"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg"`
	MaxDataBytes                      uint32 `json:"maxDataBytes"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit"`
	DestGasOverhead                   uint32 `json:"destGasOverhead"`
	DestGasPerPayloadByteBase         uint8  `json:"destGasPerPayloadByteBase"`
	DestGasPerPayloadByteHigh         uint8  `json:"destGasPerPayloadByteHigh"`
	DestGasPerPayloadByteThreshold    uint16 `json:"destGasPerPayloadByteThreshold"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps"`
	ChainFamilySelector               string `json:"chainFamilySelector"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder"`
	DefaultTokenFeeUsdCents           uint16 `json:"defaultTokenFeeUsdCents"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth"`
	GasPriceStalenessThreshold        uint32 `json:"gasPriceStalenessThreshold"`
	NetworkFeeUsdCents                uint32 `json:"networkFeeUsdCents"`
}

type RMNRemoteView struct {
	aptosCommon.ContractMetaData
	IsCursed             bool                     `json:"isCursed"`
	Config               RMNRemoteVersionedConfig `json:"config"`
	CursedSubjectEntries []RMNRemoteCurseEntry    `json:"cursedSubjectEntries"`
}

type RMNRemoteVersionedConfig struct {
	Version uint32            `json:"version"`
	Signers []RMNRemoteSigner `json:"signers"`
	Fsign   uint64            `json:"fSign"`
}

type RMNRemoteSigner struct {
	OnchainPublicKey string `json:"onchain_public_key"` // Follow EVM snake_case
	NodeIndex        uint64 `json:"node_index"`
}

type RMNRemoteCurseEntry struct {
	Subject  string `json:"subject"`
	Selector uint64 `json:"selector"`
}

type TokenAdminRegistryView struct {
	aptosCommon.ContractMetaData
}

type NonceManagerView struct {
	aptosCommon.ContractMetaData
}

type ReceiverRegistryView struct {
	aptosCommon.ContractMetaData
}

func GenerateCCIPView(chain cldf_aptos.Chain, ccipAddress aptos.AccountAddress, routerAddress aptos.AccountAddress) (CCIPView, error) {
	boundCCIP := ccip.Bind(ccipAddress, chain.Client)
	boundRouter := ccip_router.Bind(ccipAddress, chain.Client)

	ccipOwner, err := boundCCIP.Auth().Owner(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get owner of CCIP %s: %w", ccipAddress.StringLong(), err)
	}

	// Router
	destChainSelectors, err := boundRouter.Router().GetDestChains(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get destChains of router %s: %w", routerAddress.StringLong(), err)
	}

	// FeeQuoter
	feeQuoterTypeAndVersion, err := boundCCIP.FeeQuoter().TypeAndVersion(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get typeAndVersion of feeQuoter %s: %w", ccipAddress.StringLong(), err)
	}
	feeQuoterStaticConfig, err := boundCCIP.FeeQuoter().GetStaticConfig(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get staticConfig of feeQuoter %s: %w", ccipAddress.StringLong(), err)
	}
	feeQuoterFeeTkns, err := boundCCIP.FeeQuoter().GetFeeTokens(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get feeTokens of feeQuoter %s: %w", ccipAddress.StringLong(), err)
	}
	feeQuoterFeeTokens := make([]string, len(feeQuoterFeeTkns))
	for i, feeToken := range feeQuoterFeeTkns {
		feeQuoterFeeTokens[i] = feeToken.StringLong()
	}
	destinationChainConfigs := make(map[uint64]FeeQuoterDestChainConfig, len(destChainSelectors))
	for _, selector := range destChainSelectors {
		destChainConfig, err := boundCCIP.FeeQuoter().GetDestChainConfig(nil, selector)
		if err != nil {
			return CCIPView{}, fmt.Errorf("failed to get destChainConfig for chain %d of feeQuoter %s: %w", selector, ccipAddress.StringLong(), err)
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
			ChainFamilySelector:               hex.EncodeToString(destChainConfig.ChainFamilySelector),
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
	rmnRemoteTypeAndVersion, err := boundCCIP.RMNRemote().TypeAndVersion(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get typeAndVersion of RMNRemote %s: %w", ccipAddress.StringLong(), err)
	}
	cursedSubjects, err := boundCCIP.RMNRemote().GetCursedSubjects(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get cursedSubjects of RMNRemote %s: %w", ccipAddress.StringLong(), err)
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
		return CCIPView{}, fmt.Errorf("failed to get versionedConfig of RMNRemote %s: %w", ccipAddress.StringLong(), err)
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

	// Token Admin Registry
	tokenAdminRegistryTypeAndVersion, err := boundCCIP.TokenAdminRegistry().TypeAndVersion(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get typeAndVersion of TokenAdminRegistry %s: %w", ccipAddress.StringLong(), err)
	}

	// Nonce Manager
	nonceManagerTypeAndVersion, err := boundCCIP.NonceManager().TypeAndVersion(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get typeAndVersion of NonceManager %s: %w", ccipAddress.StringLong(), err)
	}

	// Receiver Registry
	receiverRegistryTypeAndVersion, err := boundCCIP.ReceiverRegistry().TypeAndVersion(nil)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get typeAndVersion of ReceiverRegistry %s: %w", ccipAddress.StringLong(), err)
	}

	return CCIPView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address: ccipAddress.StringLong(),
			Owner:   ccipOwner.StringLong(),
		},
		FeeQuoter: FeeQuoterView{
			ContractMetaData: aptosCommon.ContractMetaData{
				Address:        ccipAddress.StringLong(),
				TypeAndVersion: feeQuoterTypeAndVersion,
			},
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
				Address:        ccipAddress.StringLong(),
				TypeAndVersion: rmnRemoteTypeAndVersion,
			},
			IsCursed:             len(cursedSubjectEntries) != 0,
			Config:               rmnRemoteVersionedConfig,
			CursedSubjectEntries: cursedSubjectEntries,
		},
		TokenAdminRegistry: TokenAdminRegistryView{
			ContractMetaData: aptosCommon.ContractMetaData{
				Address:        ccipAddress.StringLong(),
				TypeAndVersion: tokenAdminRegistryTypeAndVersion,
			},
		},
		NonceManager: NonceManagerView{
			ContractMetaData: aptosCommon.ContractMetaData{
				Address:        ccipAddress.StringLong(),
				TypeAndVersion: nonceManagerTypeAndVersion,
			},
		},
		ReceiverRegistry: ReceiverRegistryView{
			ContractMetaData: aptosCommon.ContractMetaData{
				Address:        ccipAddress.StringLong(),
				TypeAndVersion: receiverRegistryTypeAndVersion,
			},
		},
	}, nil
}

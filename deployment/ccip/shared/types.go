package shared

import (
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	// Legacy
	CommitStore   deployment.ContractType = "CommitStore"
	PriceRegistry deployment.ContractType = "PriceRegistry"
	RMN           deployment.ContractType = "RMN"
	// Legacy v1.5 datastore aliases (same bindings as OnRamp/OffRamp @ v1.5.0; matches on-chain typeAndVersion).
	EVM2EVMOnRamp  deployment.ContractType = "EVM2EVMOnRamp"
	EVM2EVMOffRamp deployment.ContractType = "EVM2EVMOffRamp"

	// Not legacy
	MockRMN              deployment.ContractType = "MockRMN"
	RMNRemote            deployment.ContractType = "RMNRemote"
	ARMProxy             deployment.ContractType = "ARMProxy"
	WETH9                deployment.ContractType = "WETH9"
	Router               deployment.ContractType = "Router"
	TokenAdminRegistry   deployment.ContractType = "TokenAdminRegistry"
	TokenPoolFactory     deployment.ContractType = "TokenPoolFactory"
	RegistryModule       deployment.ContractType = "RegistryModuleOwnerCustom"
	NonceManager         deployment.ContractType = "NonceManager"
	FeeQuoter            deployment.ContractType = "FeeQuoter"
	CCIPHome             deployment.ContractType = "CCIPHome"
	RMNHome              deployment.ContractType = "RMNHome"
	OnRamp               deployment.ContractType = "OnRamp"
	OffRamp              deployment.ContractType = "OffRamp"
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	DonIDClaimer         deployment.ContractType = "DonIDClaimer"
	PriceFeed            deployment.ContractType = "PriceFeed"
	TokenGovernor        deployment.ContractType = "TokenGovernor"

	// Test contracts. Note test router maps to a regular router contract.
	TestRouter             deployment.ContractType = "TestRouter"
	Multicall3             deployment.ContractType = "Multicall3"
	CCIPReceiver           deployment.ContractType = "CCIPReceiver"
	LogMessageDataReceiver deployment.ContractType = "LogMessageDataReceiver"
	USDCMockTransmitter    deployment.ContractType = "USDCMockTransmitter"

	// Pools
	BurnMintToken                                   deployment.ContractType = "BurnMintToken"
	BurnMintERC20Token                              deployment.ContractType = "BurnMintERC20Token"
	FactoryBurnMintERC20Token                       deployment.ContractType = "FactoryBurnMintERC20Token"
	ERC20Token                                      deployment.ContractType = "ERC20Token"
	ERC677Token                                     deployment.ContractType = "ERC677Token"
	ERC677TokenHelper                               deployment.ContractType = "ERC677TokenHelper"
	BurnMintTokenPool                               deployment.ContractType = "BurnMintTokenPool"
	BurnWithFromMintTokenPool                       deployment.ContractType = "BurnWithFromMintTokenPool"
	BurnMintFastTransferTokenPool                   deployment.ContractType = "BurnMintFastTransferTokenPool"
	BurnMintWithExternalMinterFastTransferTokenPool deployment.ContractType = "BurnMintWithExternalMinterFastTransferTokenPool"
	BurnFromMintTokenPool                           deployment.ContractType = "BurnFromMintTokenPool"
	LockReleaseTokenPool                            deployment.ContractType = "LockReleaseTokenPool"
	USDCToken                                       deployment.ContractType = "USDCToken"
	USDCTokenMessenger                              deployment.ContractType = "USDCTokenMessenger"
	USDCTokenPool                                   deployment.ContractType = "USDCTokenPool"
	CCTPMessageTransmitterProxy                     deployment.ContractType = "CCTPMessageTransmitterProxy"
	HybridLockReleaseUSDCTokenPool                  deployment.ContractType = "HybridLockReleaseUSDCTokenPool"
	HybridWithExternalMinterFastTransferTokenPool   deployment.ContractType = "HybridWithExternalMinterFastTransferTokenPool"
	BurnMintWithExternalMinterTokenPool             deployment.ContractType = "BurnMintWithExternalMinterTokenPool"
	HybridWithExternalMinterTokenPool               deployment.ContractType = "HybridWithExternalMinterTokenPool"
	BurnMintERC20TransparentToken                   deployment.ContractType = "BurnMintERC20TransparentToken"
	BurnMintERC20PausableFreezableTransparentToken  deployment.ContractType = "BurnMintERC20PausableFreezableTransparentToken"
	USDCTokenPoolProxy                              deployment.ContractType = "USDCTokenPoolProxy"

	// Firedrill
	FiredrillEntrypointType deployment.ContractType = "FiredrillEntrypoint"

	// Treasury
	FeeAggregator deployment.ContractType = "FeeAggregator"

	// Solana
	Receiver             deployment.ContractType = "Receiver"
	SPL2022Tokens        deployment.ContractType = "SPL2022Tokens"
	SPLTokens            deployment.ContractType = "SPLTokens"
	WSOL                 deployment.ContractType = "WSOL"
	CCIPCommon           deployment.ContractType = "CCIPCommon"
	RemoteSource         deployment.ContractType = "RemoteSource"
	RemoteDest           deployment.ContractType = "RemoteDest"
	TokenPoolLookupTable deployment.ContractType = "TokenPoolLookupTable"
	CCTPTokenPool        deployment.ContractType = "CCTPTokenPool"
	BPFUpgradeable       deployment.ContractType = "BPFUpgradeable"
	SVMSignerRegistry    deployment.ContractType = "SVMSignerRegistry"
	// CLL Identifier
	CLLMetadata = "CLL"

	// Aptos
	AptosMCMSType               deployment.ContractType = "AptosManyChainMultisig"
	AptosCurseMCMSType          deployment.ContractType = "AptosCurseMCMS"
	AptosCCIPType               deployment.ContractType = "AptosCCIP"
	AptosReceiverType           deployment.ContractType = "AptosReceiver"
	AptosManagedTokenPoolType   deployment.ContractType = "AptosManagedTokenPool"
	AptosRegulatedTokenPoolType deployment.ContractType = "AptosRegulatedTokenPool"
	AptosManagedTokenType       deployment.ContractType = "AptosManagedTokenType"
	AptosRegulatedTokenType     deployment.ContractType = "AptosRegulatedTokenType"
	AptosTestTokenType          deployment.ContractType = "AptosTestToken"

	// TON, [NONEVM-1938] currently added necessary contract for unblocking e2e env setup
	TonCCIP     deployment.ContractType = "TonCCIP"
	TonReceiver deployment.ContractType = "TonReceiver"

	// Attestation Service
	EVMSignerRegistry deployment.ContractType = "SignerRegistry"

	// OpenZeppelin
	ProxyAdmin                  deployment.ContractType = "ProxyAdmin"
	TransparentUpgradeableProxy deployment.ContractType = "TransparentUpgradeableProxy"
)

type OCRParameters struct {
	DeltaProgress                           time.Duration `json:"deltaProgress"`
	DeltaResend                             time.Duration `json:"deltaResend"`
	DeltaInitial                            time.Duration `json:"deltaInitial"`
	DeltaRound                              time.Duration `json:"deltaRound"`
	DeltaGrace                              time.Duration `json:"deltaGrace"`
	DeltaCertifiedCommitRequest             time.Duration `json:"deltaCertifiedCommitRequest"`
	DeltaStage                              time.Duration `json:"deltaStage"`
	Rmax                                    uint64        `json:"rmax"`
	MaxDurationQuery                        time.Duration `json:"maxDurationQuery"`
	MaxDurationObservation                  time.Duration `json:"maxDurationObservation"`
	MaxDurationShouldAcceptAttestedReport   time.Duration `json:"maxDurationShouldAcceptAttestedReport"`
	MaxDurationShouldTransmitAcceptedReport time.Duration `json:"maxDurationShouldTransmitAcceptedReport"`
}

func (params OCRParameters) Validate() error {
	if params.DeltaProgress <= 0 {
		return errors.New("deltaProgress must be positive")
	}
	if params.DeltaResend <= 0 {
		return errors.New("deltaResend must be positive")
	}
	if params.DeltaInitial <= 0 {
		return errors.New("deltaInitial must be positive")
	}
	if params.DeltaRound <= 0 {
		return errors.New("deltaRound must be positive")
	}
	if params.DeltaGrace <= 0 {
		return errors.New("deltaGrace must be positive")
	}
	if params.DeltaCertifiedCommitRequest <= 0 {
		return errors.New("deltaCertifiedCommitRequest must be positive")
	}
	if params.DeltaStage < 0 {
		return errors.New("deltaStage must be positive or 0 for disabled")
	}
	if params.Rmax <= 0 {
		return errors.New("rmax must be positive")
	}
	if params.MaxDurationQuery <= 0 {
		return errors.New("maxDurationQuery must be positive")
	}
	if params.MaxDurationObservation <= 0 {
		return errors.New("maxDurationObservation must be positive")
	}
	if params.MaxDurationShouldAcceptAttestedReport <= 0 {
		return errors.New("maxDurationShouldAcceptAttestedReport must be positive")
	}
	if params.MaxDurationShouldTransmitAcceptedReport <= 0 {
		return errors.New("maxDurationShouldTransmitAcceptedReport must be positive")
	}
	return nil
}

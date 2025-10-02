package shared

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	// Legacy
	CommitStore   deployment.ContractType = "CommitStore"
	PriceRegistry deployment.ContractType = "PriceRegistry"
	RMN           deployment.ContractType = "RMN"

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
)

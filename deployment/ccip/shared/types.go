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

	// Test contracts. Note test router maps to a regular router contract.
	TestRouter             deployment.ContractType = "TestRouter"
	Multicall3             deployment.ContractType = "Multicall3"
	CCIPReceiver           deployment.ContractType = "CCIPReceiver"
	LogMessageDataReceiver deployment.ContractType = "LogMessageDataReceiver"
	USDCMockTransmitter    deployment.ContractType = "USDCMockTransmitter"

	// Pools
	BurnMintToken                  deployment.ContractType = "BurnMintToken"
	FactoryBurnMintERC20Token      deployment.ContractType = "FactoryBurnMintERC20Token"
	ERC20Token                     deployment.ContractType = "ERC20Token"
	ERC677Token                    deployment.ContractType = "ERC677Token"
	ERC677TokenHelper              deployment.ContractType = "ERC677TokenHelper"
	BurnMintTokenPool              deployment.ContractType = "BurnMintTokenPool"
	BurnWithFromMintTokenPool      deployment.ContractType = "BurnWithFromMintTokenPool"
	BurnFromMintTokenPool          deployment.ContractType = "BurnFromMintTokenPool"
	LockReleaseTokenPool           deployment.ContractType = "LockReleaseTokenPool"
	USDCToken                      deployment.ContractType = "USDCToken"
	USDCTokenMessenger             deployment.ContractType = "USDCTokenMessenger"
	USDCTokenPool                  deployment.ContractType = "USDCTokenPool"
	HybridLockReleaseUSDCTokenPool deployment.ContractType = "HybridLockReleaseUSDCTokenPool"

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
	// CLL Identifier
	CLLMetadata = "CLL"

	// Aptos
	AptosMCMSType             deployment.ContractType = "AptosManyChainMultisig"
	AptosCCIPType             deployment.ContractType = "AptosCCIP"
	AptosReceiverType         deployment.ContractType = "AptosReceiver"
	AptosManagedTokenPoolType deployment.ContractType = "AptosManagedTokenPool"
	AptosManagedTokenType     deployment.ContractType = "AptosManagedTokenType"
)

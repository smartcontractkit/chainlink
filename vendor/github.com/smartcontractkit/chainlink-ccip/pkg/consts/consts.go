package consts

// This package contains ChainReader and ChainWriter related constants.

// Contract Names
const (
	ContractNameOffRamp              = "OffRamp"
	ContractNameOnRamp               = "OnRamp"
	ContractNamePriceRegistry        = "PriceRegistry"
	ContractNameCapabilitiesRegistry = "CapabilitiesRegistry"
	ContractNameCCIPConfig           = "CCIPConfig"
)

// Method Names
// TODO: these should be better organized, maybe separate packages.
const (
	// Offramp methods
	MethodNameGetSourceChainConfig         = "GetSourceChainConfig"
	MethodNameOfframpGetDynamicConfig      = "OfframpGetDynamicConfig"
	MethodNameOfframpGetStaticConfig       = "OfframpGetStaticConfig"
	MethodNameGetLatestPriceSequenceNumber = "GetLatestPriceSequenceNumber"
	MethodNameIsBlessed                    = "IsBlessed"
	MethodNameGetMerkleRoot                = "GetMerkleRoot"
	MethodNameGetExecutionState            = "GetExecutionState"

	// Onramp methods
	MethodNameOnrampGetDynamicConfig        = "OnrampGetDynamicConfig"
	MethodNameOnrampGetStaticConfig         = "OnrampGetStaticConfig"
	MethodNameGetExpectedNextSequenceNumber = "GetExpectedNextSequenceNumber"

	// Price registry view/pure methods
	MethodNamePriceRegistryGetStaticConfig  = "GetStaticConfig"
	MethodNameGetDestChainConfig            = "GetDestChainConfig"
	MethodNameGetPremiumMultiplierWeiPerEth = "GetPremiumMultiplierWeiPerEth"
	MethodNameGetTokenTransferFeeConfig     = "GetTokenTransferFeeConfig"
	MethodNameProcessMessageArgs            = "ProcessMessageArgs"
	MethodNameValidatePoolReturnData        = "ValidatePoolReturnData"
	MethodNameGetValidatedTokenPrice        = "GetValidatedTokenPrice"
	MethodNameGetFeeTokens                  = "GetFeeTokens"

	/*
		// On EVM:
		function commit(
			bytes32[3] calldata reportContext,
			    bytes calldata report,
			    bytes32[] calldata rs,
			    bytes32[] calldata ss,
			    bytes32 rawVs // signatures
			  ) external
	*/
	MethodCommit = "Commit"

	// On EVM:
	// function execute(bytes32[3] calldata reportContext, bytes calldata report) external
	MethodExecute = "Execute"

	// Capability registry methods.
	// Used by the home chain reader.
	MethodNameGetCapability = "GetCapability"

	// CCIPConfig.sol methods.
	// Used by the home chain reader.
	MethodNameGetAllChainConfigs = "GetAllChainConfigs"
	MethodNameGetOCRConfig       = "GetOCRConfig"
)

// Event Names
const (
	EventNameCCIPSendRequested     = "CCIPSendRequested"
	EventNameExecutionStateChanged = "ExecutionStateChanged"
	EventNameCommitReportAccepted  = "CommitReportAccepted"
)

// Event Attributes
const (
	EventAttributeSequenceNumber = "SequenceNumber"
	EventAttributeSourceChain    = "SourceChain"
	EventAttributeDestChain      = "DestChain"
)

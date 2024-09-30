package evm

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	onrampABI               = evmtypes.MustGetABI(onramp.OnRampABI)
	capabilitiesRegsitryABI = evmtypes.MustGetABI(kcr.CapabilitiesRegistryABI)
	ccipConfigABI           = evmtypes.MustGetABI(ccip_config.CCIPConfigABI)
	feeQuoterABI            = evmtypes.MustGetABI(fee_quoter.FeeQuoterABI)
	nonceManagerABI         = evmtypes.MustGetABI(nonce_manager.NonceManagerABI)
	priceFeedABI            = evmtypes.MustGetABI(aggregator_v3_interface.AggregatorV3InterfaceABI)
	rmnRemoteABI            = evmtypes.MustGetABI(rmn_remote.RMNRemoteABI)
	rmnHomeABI              = evmtypes.MustGetABI(rmnHomeString)
)

// TODO: replace with generated ABI when the contract will be defined
var rmnHomeString = "[{\"inputs\":[],\"name\":\"getAllConfigs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// MustSourceReaderConfig returns a ChainReaderConfig that can be used to read from the onramp.
// The configuration is marshaled into JSON so that it can be passed to the relayer NewContractReader() method.
func MustSourceReaderConfig() []byte {
	rawConfig := SourceReaderConfig
	encoded, err := json.Marshal(rawConfig)
	if err != nil {
		panic(fmt.Errorf("failed to marshal ChainReaderConfig into JSON: %w", err))
	}

	return encoded
}

// MustDestReaderConfig returns a ChainReaderConfig that can be used to read from the offramp.
// The configuration is marshaled into JSON so that it can be passed to the relayer NewContractReader() method.
func MustDestReaderConfig() []byte {
	rawConfig := DestReaderConfig
	encoded, err := json.Marshal(rawConfig)
	if err != nil {
		panic(fmt.Errorf("failed to marshal ChainReaderConfig into JSON: %w", err))
	}

	return encoded
}

func MergeReaderConfigs(configs ...evmrelaytypes.ChainReaderConfig) evmrelaytypes.ChainReaderConfig {
	allContracts := make(map[string]evmrelaytypes.ChainContractReader)
	for _, c := range configs {
		for contractName, contractReader := range c.Contracts {
			allContracts[contractName] = contractReader
		}
	}

	return evmrelaytypes.ChainReaderConfig{Contracts: allContracts}
}

// DestReaderConfig returns a ChainReaderConfig that can be used to read from the offramp.
var DestReaderConfig = evmrelaytypes.ChainReaderConfig{
	Contracts: map[string]evmrelaytypes.ChainContractReader{
		consts.ContractNameOffRamp: {
			ContractABI: offramp.OffRampABI,
			ContractPollingFilter: evmrelaytypes.ContractPollingFilter{
				GenericEventNames: []string{
					mustGetEventName(consts.EventNameExecutionStateChanged, offrampABI),
					mustGetEventName(consts.EventNameCommitReportAccepted, offrampABI),
				},
			},
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetExecutionState: {
					ChainSpecificName: mustGetMethodName("getExecutionState", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetMerkleRoot: {
					ChainSpecificName: mustGetMethodName("getMerkleRoot", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetLatestPriceSequenceNumber: {
					ChainSpecificName: mustGetMethodName("getLatestPriceSequenceNumber", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameOffRampGetDestChainConfig: {
					ChainSpecificName: mustGetMethodName("getDestChainConfig", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameOffRampGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameOffRampGetDynamicConfig: {
					ChainSpecificName: mustGetMethodName("getDynamicConfig", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetSourceChainConfig: {
					ChainSpecificName: mustGetMethodName("getSourceChainConfig", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.EventNameCommitReportAccepted: {
					ChainSpecificName: mustGetEventName(consts.EventNameCommitReportAccepted, offrampABI),
					ReadType:          evmrelaytypes.Event,
				},
				consts.EventNameExecutionStateChanged: {
					ChainSpecificName: mustGetEventName(consts.EventNameExecutionStateChanged, offrampABI),
					ReadType:          evmrelaytypes.Event,
				},
				// TODO: remove deprecated config.
				consts.MethodNameOfframpGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				// TODO: remove deprecated config.
				consts.MethodNameOfframpGetDynamicConfig: {
					ChainSpecificName: mustGetMethodName("getDynamicConfig", offrampABI),
					ReadType:          evmrelaytypes.Method,
				},
			},
		},
		consts.ContractNameNonceManager: {
			ContractABI: nonce_manager.NonceManagerABI,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetInboundNonce: {
					ChainSpecificName: mustGetMethodName("getInboundNonce", nonceManagerABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetOutboundNonce: {
					ChainSpecificName: mustGetMethodName("getOutboundNonce", nonceManagerABI),
					ReadType:          evmrelaytypes.Method,
				},
			},
		},
		consts.ContractNameFeeQuoter: {
			ContractABI: fee_quoter.FeeQuoterABI,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameFeeQuoterGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameFeeQuoterGetTokenPrices: {
					ChainSpecificName: mustGetMethodName("getTokenPrices", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetDestChainConfig: {
					ChainSpecificName: mustGetMethodName("getDestChainConfig", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetPremiumMultiplierWeiPerEth: {
					ChainSpecificName: mustGetMethodName("getPremiumMultiplierWeiPerEth", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetTokenTransferFeeConfig: {
					ChainSpecificName: mustGetMethodName("getTokenTransferFeeConfig", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameProcessMessageArgs: {
					ChainSpecificName: mustGetMethodName("processMessageArgs", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameProcessPoolReturnData: {
					ChainSpecificName: mustGetMethodName("processPoolReturnData", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetValidatedTokenPrice: {
					ChainSpecificName: mustGetMethodName("getValidatedTokenPrice", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameGetFeeTokens: {
					ChainSpecificName: mustGetMethodName("getFeeTokens", feeQuoterABI),
					ReadType:          evmrelaytypes.Method,
				},
			},
		},
		consts.ContractNameRMNRemote: {
			ContractABI: rmn_remote.RMNRemoteABI,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetVersionedConfig: {
					ChainSpecificName: mustGetMethodName("getVersionedConfig", rmnRemoteABI),
					ReadType:          evmrelaytypes.Method,
				},
			},
		},
	},
}

// SourceReaderConfig returns a ChainReaderConfig that can be used to read from the onramp.
var SourceReaderConfig = evmrelaytypes.ChainReaderConfig{
	Contracts: map[string]evmrelaytypes.ChainContractReader{
		consts.ContractNameOnRamp: {
			ContractABI: onramp.OnRampABI,
			ContractPollingFilter: evmrelaytypes.ContractPollingFilter{
				GenericEventNames: []string{
					consts.EventNameCCIPMessageSent,
				},
			},
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				// all "{external|public} view" functions in the onramp except for getFee and getPoolBySourceToken are here.
				// getFee is not expected to get called offchain and is only called by end-user contracts.
				consts.MethodNameGetExpectedNextSequenceNumber: {
					ChainSpecificName: mustGetMethodName("getExpectedNextSequenceNumber", onrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.EventNameCCIPMessageSent: {
					ChainSpecificName: mustGetEventName("CCIPMessageSent", onrampABI),
					ReadType:          evmrelaytypes.Event,
				},
				consts.MethodNameOnRampGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", onrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				consts.MethodNameOnRampGetDynamicConfig: {
					ChainSpecificName: mustGetMethodName("getDynamicConfig", onrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				// TODO: Remove deprecated config.
				consts.MethodNameOnrampGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", onrampABI),
					ReadType:          evmrelaytypes.Method,
				},
				// TODO: Remove deprecated config.
				consts.MethodNameOnrampGetDynamicConfig: {
					ChainSpecificName: mustGetMethodName("getDynamicConfig", onrampABI),
					ReadType:          evmrelaytypes.Method,
				},
			},
		},
	},
}

var FeedReaderConfig = evmrelaytypes.ChainReaderConfig{
	Contracts: map[string]evmrelaytypes.ChainContractReader{
		consts.ContractNamePriceAggregator: {
			ContractABI: aggregator_v3_interface.AggregatorV3InterfaceABI,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetLatestRoundData: {
					ChainSpecificName: mustGetMethodName(consts.MethodNameGetLatestRoundData, priceFeedABI),
				},
				consts.MethodNameGetDecimals: {
					ChainSpecificName: mustGetMethodName(consts.MethodNameGetDecimals, priceFeedABI),
				},
			},
		},
	},
}

// HomeChainReaderConfigRaw returns a ChainReaderConfig that can be used to read from the home chain.
var HomeChainReaderConfigRaw = evmrelaytypes.ChainReaderConfig{
	Contracts: map[string]evmrelaytypes.ChainContractReader{
		consts.ContractNameCapabilitiesRegistry: {
			ContractABI: kcr.CapabilitiesRegistryABI,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetCapability: {
					ChainSpecificName: mustGetMethodName("getCapability", capabilitiesRegsitryABI),
				},
			},
		},
		consts.ContractNameCCIPConfig: {
			ContractABI: ccip_config.CCIPConfigABI,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetAllChainConfigs: {
					ChainSpecificName: mustGetMethodName("getAllChainConfigs", ccipConfigABI),
				},
				consts.MethodNameGetOCRConfig: {
					ChainSpecificName: mustGetMethodName("getOCRConfig", ccipConfigABI),
				},
			},
		},
		consts.ContractNameRMNHome: {
			ContractABI: rmnHomeString,
			Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
				consts.MethodNameGetAllConfigs: {
					ChainSpecificName: mustGetMethodName("getAllConfigs", rmnHomeABI),
				},
			},
		},
	},
}

func mustGetEventName(event string, tabi abi.ABI) string {
	e, ok := tabi.Events[event]
	if !ok {
		panic(fmt.Sprintf("missing event %s in onrampABI", event))
	}
	return e.Name
}

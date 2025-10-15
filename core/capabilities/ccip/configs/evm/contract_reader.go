package evm

import (
	"encoding/json"
	"fmt"
	"maps"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/evm"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/aggregator_v3_interface"
)

const (
	// DefaultCCIPLogsRetention defines the duration for which logs critical for Commit/Exec plugins processing are retained.
	// Although Exec relies on permissionlessExecThreshold which is lower than 24hours for picking eligible CommitRoots,
	// Commit still can reach to older logs because it filters them by sequence numbers. For instance, in case of RMN curse on chain,
	// we might have logs waiting in OnRamp to be committed first. When outage takes days we still would
	// be able to bring back processing without replaying any logs from chain. You can read that param as
	// "how long CCIP can be down and still be able to process all the messages after getting back to life".
	// Breaching this threshold would require replaying chain using LogPoller from the beginning of the outage.
	// Using same default retention as v1.5 https://github.com/smartcontractkit/ccip/pull/530/files
	DefaultCCIPLogsRetention = 30 * 24 * time.Hour // 30 days
)

var (
	onrampABI               = evmtypes.MustGetABI(onramp.OnRampABI)
	capabilitiesRegistryABI = evmtypes.MustGetABI(kcr.CapabilitiesRegistryABI)
	ccipHomeABI             = evmtypes.MustGetABI(ccip_home.CCIPHomeABI)
	feeQuoterABI            = evmtypes.MustGetABI(fee_quoter.FeeQuoterABI)
	nonceManagerABI         = evmtypes.MustGetABI(nonce_manager.NonceManagerABI)
	priceFeedABI            = evmtypes.MustGetABI(aggregator_v3_interface.AggregatorV3InterfaceABI)
	rmnRemoteABI            = evmtypes.MustGetABI(rmn_remote.RMNRemoteABI)
	rmnProxyABI             = evmtypes.MustGetABI(rmn_proxy_contract.RMNProxyABI)
	rmnHomeABI              = evmtypes.MustGetABI(rmn_home.RMNHomeABI)
	routerABI               = evmtypes.MustGetABI(router.RouterABI)
)

func MergeReaderConfigs(configs ...config.ChainReaderConfig) config.ChainReaderConfig {
	allContracts := make(map[string]config.ChainContractReader)
	for _, c := range configs {
		maps.Copy(allContracts, c.Contracts)
	}

	return config.ChainReaderConfig{Contracts: allContracts}
}

// DestReaderConfig returns a ChainReaderConfig that can be used to read from the offramp.
var DestReaderConfig = config.ChainReaderConfig{
	Contracts: map[string]config.ChainContractReader{
		consts.ContractNameOffRamp: {
			ContractABI: offramp.OffRampABI,
			ContractPollingFilter: config.ContractPollingFilter{
				GenericEventNames: []string{
					mustGetEventName(consts.EventNameExecutionStateChanged, offrampABI),
					mustGetEventName(consts.EventNameCommitReportAccepted, offrampABI),
				},
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetExecutionState: {
					ChainSpecificName: mustGetMethodName("getExecutionState", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetMerkleRoot: {
					ChainSpecificName: mustGetMethodName("getMerkleRoot", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetLatestPriceSequenceNumber: {
					ChainSpecificName: mustGetMethodName("getLatestPriceSequenceNumber", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameOffRampGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameOffRampGetDynamicConfig: {
					ChainSpecificName: mustGetMethodName("getDynamicConfig", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetSourceChainConfig: {
					ChainSpecificName: mustGetMethodName("getSourceChainConfig", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameOffRampGetAllSourceChainConfigs: {
					ChainSpecificName: mustGetMethodName("getAllSourceChainConfigs", offrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameOffRampLatestConfigDetails: {
					ChainSpecificName: mustGetMethodName("latestConfigDetails", offrampABI),
					ReadType:          config.Method,
				},
				consts.EventNameCommitReportAccepted: {
					ChainSpecificName: mustGetEventName(consts.EventNameCommitReportAccepted, offrampABI),
					ReadType:          config.Event,
				},
				consts.EventNameExecutionStateChanged: {
					ChainSpecificName: mustGetEventName(consts.EventNameExecutionStateChanged, offrampABI),
					ReadType:          config.Event,
					EventDefinitions: &config.EventDefinitions{
						GenericTopicNames: map[string]string{
							"sourceChainSelector": consts.EventAttributeSourceChain,
							"sequenceNumber":      consts.EventAttributeSequenceNumber,
						},
						GenericDataWordDetails: map[string]evm.DataWordDetail{
							consts.EventAttributeState: {
								Name: "state",
							},
						},
					},
				},
			},
		},
		consts.ContractNameNonceManager: {
			ContractABI: nonce_manager.NonceManagerABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetInboundNonce: {
					ChainSpecificName: mustGetMethodName("getInboundNonce", nonceManagerABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetOutboundNonce: {
					ChainSpecificName: mustGetMethodName("getOutboundNonce", nonceManagerABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameFeeQuoter: {
			ContractABI: fee_quoter.FeeQuoterABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameFeeQuoterGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameFeeQuoterGetTokenPrices: {
					ChainSpecificName: mustGetMethodName("getTokenPrices", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameFeeQuoterGetTokenPrice: {
					ChainSpecificName: mustGetMethodName("getTokenPrice", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetFeePriceUpdate: {
					ChainSpecificName: mustGetMethodName("getDestinationChainGasPrice", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetDestChainConfig: {
					ChainSpecificName: mustGetMethodName("getDestChainConfig", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetPremiumMultiplierWeiPerEth: {
					ChainSpecificName: mustGetMethodName("getPremiumMultiplierWeiPerEth", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetTokenTransferFeeConfig: {
					ChainSpecificName: mustGetMethodName("getTokenTransferFeeConfig", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameProcessMessageArgs: {
					ChainSpecificName: mustGetMethodName("processMessageArgs", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetValidatedTokenPrice: {
					ChainSpecificName: mustGetMethodName("getValidatedTokenPrice", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetFeeTokens: {
					ChainSpecificName: mustGetMethodName("getFeeTokens", feeQuoterABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameRMNRemote: {
			ContractABI: rmn_remote.RMNRemoteABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetVersionedConfig: {
					ChainSpecificName: mustGetMethodName("getVersionedConfig", rmnRemoteABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetReportDigestHeader: {
					ChainSpecificName: mustGetMethodName("getReportDigestHeader", rmnRemoteABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetCursedSubjects: {
					ChainSpecificName: mustGetMethodName("getCursedSubjects", rmnRemoteABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameRMNProxy: {
			ContractABI: rmn_proxy_contract.RMNProxyABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetARM: {
					ChainSpecificName: mustGetMethodName("getARM", rmnProxyABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameRouter: {
			ContractABI: router.RouterABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameRouterGetWrappedNative: {
					ChainSpecificName: mustGetMethodName("getWrappedNative", routerABI),
					ReadType:          config.Method,
				},
			},
		},
	},
}

// SourceReaderConfig returns a ChainReaderConfig that can be used to read from the onramp.
var SourceReaderConfig = config.ChainReaderConfig{
	Contracts: map[string]config.ChainContractReader{
		consts.ContractNameOnRamp: {
			ContractABI: onramp.OnRampABI,
			ContractPollingFilter: config.ContractPollingFilter{
				GenericEventNames: []string{
					consts.EventNameCCIPMessageSent,
				},
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				// all "{external|public} view" functions in the onramp except for getFee and getPoolBySourceToken are here.
				// getFee is not expected to get called offchain and is only called by end-user contracts.
				consts.MethodNameGetExpectedNextSequenceNumber: {
					ChainSpecificName: mustGetMethodName("getExpectedNextSequenceNumber", onrampABI),
					ReadType:          config.Method,
				},
				consts.EventNameCCIPMessageSent: {
					ChainSpecificName: mustGetEventName("CCIPMessageSent", onrampABI),
					ReadType:          config.Event,
					EventDefinitions: &config.EventDefinitions{
						GenericTopicNames: map[string]string{
							"destChainSelector": consts.EventAttributeDestChain,
							"sequenceNumber":    consts.EventAttributeSequenceNumber,
						},
					},
					OutputModifications: codec.ModifiersConfig{
						&codec.WrapperModifierConfig{Fields: map[string]string{
							"Message.FeeTokenAmount":      "Int",
							"Message.FeeValueJuels":       "Int",
							"Message.TokenAmounts.Amount": "Int",
						}},
					},
				},
				consts.MethodNameOnRampGetStaticConfig: {
					ChainSpecificName: mustGetMethodName("getStaticConfig", onrampABI),
					ReadType:          config.Method,
				},
				consts.MethodNameOnRampGetDynamicConfig: {
					ChainSpecificName: mustGetMethodName("getDynamicConfig", onrampABI),
					ReadType:          config.Method,
				},
				// TODO: swap with const.
				consts.MethodNameOnRampGetDestChainConfig: {
					ChainSpecificName: mustGetMethodName("getDestChainConfig", onrampABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameRouter: {
			ContractABI: router.RouterABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameRouterGetWrappedNative: {
					ChainSpecificName: mustGetMethodName("getWrappedNative", routerABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameFeeQuoter: {
			ContractABI: fee_quoter.FeeQuoterABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameFeeQuoterGetTokenPrices: {
					ChainSpecificName: mustGetMethodName("getTokenPrices", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameFeeQuoterGetTokenPrice: {
					ChainSpecificName: mustGetMethodName("getTokenPrice", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetFeePriceUpdate: {
					ChainSpecificName: mustGetMethodName("getDestinationChainGasPrice", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetDestChainConfig: {
					ChainSpecificName: mustGetMethodName("getDestChainConfig", feeQuoterABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetFeeTokens: {
					ChainSpecificName: mustGetMethodName("getFeeTokens", feeQuoterABI),
					ReadType:          config.Method,
				},
			},
		},
		consts.ContractNameRMNRemote: {
			ContractABI: rmn_remote.RMNRemoteABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetVersionedConfig: {
					ChainSpecificName: mustGetMethodName("getVersionedConfig", rmnRemoteABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetReportDigestHeader: {
					ChainSpecificName: mustGetMethodName("getReportDigestHeader", rmnRemoteABI),
					ReadType:          config.Method,
				},
				consts.MethodNameGetCursedSubjects: {
					ChainSpecificName: mustGetMethodName("getCursedSubjects", rmnRemoteABI),
					ReadType:          config.Method,
				},
			},
		},
	},
}

// FeedReaderConfig provides a ChainReaderConfig that can be used to read from a price feed
// that is deployed on-chain.
var FeedReaderConfig = config.ChainReaderConfig{
	Contracts: map[string]config.ChainContractReader{
		consts.ContractNamePriceAggregator: {
			ContractABI: aggregator_v3_interface.AggregatorV3InterfaceABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
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

var USDCReaderConfig = config.ChainReaderConfig{
	Contracts: map[string]config.ChainContractReader{
		consts.ContractNameCCTPMessageTransmitter: {
			ContractABI: MessageTransmitterABI,
			ContractPollingFilter: config.ContractPollingFilter{
				GenericEventNames: []string{consts.EventNameCCTPMessageSent},
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.EventNameCCTPMessageSent: {
					ChainSpecificName: consts.EventNameCCTPMessageSent,
					ReadType:          config.Event,
					EventDefinitions: &config.EventDefinitions{
						GenericDataWordDetails: map[string]evm.DataWordDetail{
							consts.CCTPMessageSentValue: {
								Name: consts.CCTPMessageSentValue,
								// Filtering by the 3rd word (indexing starts from 0) so it's ptr(2)
								Index: ptr(2),
								Type:  "bytes32",
							},
						},
					},
				},
			},
		},
	},
}

// HomeChainReaderConfigRaw returns a ChainReaderConfig that can be used to read from the home chain.
var HomeChainReaderConfigRaw = config.ChainReaderConfig{
	Contracts: map[string]config.ChainContractReader{
		consts.ContractNameCapabilitiesRegistry: {
			ContractABI: kcr.CapabilitiesRegistryABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetCapability: {
					ChainSpecificName: mustGetMethodName("getCapability", capabilitiesRegistryABI),
				},
			},
		},
		consts.ContractNameCCIPConfig: {
			ContractABI: ccip_home.CCIPHomeABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetAllChainConfigs: {
					ChainSpecificName: mustGetMethodName("getAllChainConfigs", ccipHomeABI),
				},
				consts.MethodNameGetOCRConfig: {
					ChainSpecificName: mustGetMethodName("getAllConfigs", ccipHomeABI),
				},
			},
		},
		consts.ContractNameRMNHome: {
			ContractABI: rmn_home.RMNHomeABI,
			ContractPollingFilter: config.ContractPollingFilter{
				PollingFilter: config.PollingFilter{
					Retention: sqlutil.Interval(DefaultCCIPLogsRetention),
				},
			},
			Configs: map[string]*config.ChainReaderDefinition{
				consts.MethodNameGetAllConfigs: {
					ChainSpecificName: mustGetMethodName("getAllConfigs", rmnHomeABI),
				},
			},
		},
	},
}

var HomeChainReaderConfig = mustMarshal(HomeChainReaderConfigRaw)

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func mustGetEventName(event string, tabi abi.ABI) string {
	e, ok := tabi.Events[event]
	if !ok {
		panic(fmt.Sprintf("missing event %s in onrampABI", event))
	}
	return e.Name
}

func ptr[T any](v T) *T {
	return &v
}

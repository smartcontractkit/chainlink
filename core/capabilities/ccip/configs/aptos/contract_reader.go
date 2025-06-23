package aptosconfig

import (
	"time"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
)

func GetChainReaderConfig() (config.ChainReaderConfig, error) {
	return config.ChainReaderConfig{
		IsLoopPlugin: true,
		Modules: map[string]*config.ChainReaderModule{
			consts.ContractNameRMNRemote: {
				Name: "rmn_remote",
				Functions: map[string]*config.ChainReaderFunction{
					consts.MethodNameGetReportDigestHeader: {
						Name: "get_report_digest_header",
					},
					consts.MethodNameGetVersionedConfig: {
						Name: "get_versioned_config",
						// ref: https://github.com/smartcontractkit/chainlink-ccip/blob/bee7c32c71cf0aec594c051fef328b4a7281a1fc/pkg/reader/ccip.go#L1440
						ResultTupleToStruct: []string{"version", "config"},
					},
					consts.MethodNameGetCursedSubjects: {
						Name: "get_cursed_subjects",
					},
				},
			},
			consts.ContractNameRMNProxy: {
				Name: "rmn_remote",
				Functions: map[string]*config.ChainReaderFunction{
					consts.MethodNameGetARM: {
						Name: "get_arm",
					},
				},
			},
			consts.ContractNameFeeQuoter: {
				Name: "fee_quoter",
				Functions: map[string]*config.ChainReaderFunction{
					consts.MethodNameFeeQuoterGetTokenPrice: {
						Name: "get_token_price",
						Params: []config.AptosFunctionParam{
							{
								Name:     "token",
								Type:     "address",
								Required: true,
							},
						},
					},
					consts.MethodNameFeeQuoterGetTokenPrices: {
						Name: "get_token_prices",
						Params: []config.AptosFunctionParam{
							{
								Name:     "tokens",
								Type:     "vector<address>",
								Required: true,
							},
						},
					},
					consts.MethodNameFeeQuoterGetStaticConfig: {
						Name: "get_static_config",
					},
					consts.MethodNameGetFeePriceUpdate: {
						Name: "get_dest_chain_gas_price",
						Params: []config.AptosFunctionParam{
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
			},
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*config.ChainReaderFunction{
					consts.MethodNameGetExecutionState: {
						Name: "get_execution_state",
						Params: []config.AptosFunctionParam{
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
							{
								Name:     "sequenceNumber",
								Type:     "u64",
								Required: true,
							},
						},
					},
					consts.MethodNameGetMerkleRoot: {
						Name: "get_merkle_root",
						Params: []config.AptosFunctionParam{
							{
								Name:     "root",
								Type:     "vector<u8>",
								Required: true,
							},
						},
					},
					consts.MethodNameOffRampLatestConfigDetails: {
						Name: "latest_config_details",
						Params: []config.AptosFunctionParam{
							{
								Name:     "ocrPluginType",
								Type:     "u8",
								Required: true,
							},
						},
						// wrap the returned OCR config
						// https://github.com/smartcontractkit/chainlink-ccip/blob/bee7c32c71cf0aec594c051fef328b4a7281a1fc/pkg/reader/ccip.go#L141
						ResultTupleToStruct: []string{"ocr_config"},
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						Name: "get_latest_price_sequence_number",
					},
					consts.MethodNameOffRampGetStaticConfig: {
						Name: "get_static_config",
					},
					consts.MethodNameOffRampGetDynamicConfig: {
						Name: "get_dynamic_config",
					},
					consts.MethodNameGetSourceChainConfig: {
						Name: "get_source_chain_config",
						Params: []config.AptosFunctionParam{
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					consts.EventNameExecutionStateChanged: {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "execution_state_changed_events",
						EventAccountAddress:   "offramp::get_state_address",
					},
					consts.EventNameCommitReportAccepted: {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "commit_report_accepted_events",
						EventAccountAddress:   "offramp::get_state_address",
					},
					"OCRConfigSet": {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "ocr3_base_state.config_set_events",
						EventAccountAddress:   "offramp::get_state_address",
					},
					"SourceChainConfigSet": {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "source_chain_config_set_events",
						EventAccountAddress:   "offramp::get_state_address",
					},
				},
			},
			consts.ContractNameOnRamp: {
				Name: "onramp",
				Functions: map[string]*config.ChainReaderFunction{
					consts.MethodNameOnRampGetDynamicConfig: {
						Name: "get_dynamic_config",
					},
					consts.MethodNameOnRampGetStaticConfig: {
						Name: "get_static_config",
					},
					consts.MethodNameOnRampGetDestChainConfig: {
						Name: "get_dest_chain_config",
						Params: []config.AptosFunctionParam{
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"sequenceNumber", "allowListEnabled", "router"},
					},
					consts.MethodNameGetExpectedNextSequenceNumber: {
						Name: "get_expected_next_sequence_number",
						Params: []config.AptosFunctionParam{
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					consts.EventNameCCIPMessageSent: {
						EventHandleStructName: "OnRampState",
						EventHandleFieldName:  "ccip_message_sent_events",
						EventAccountAddress:   "onramp::get_state_address",
						EventFilterRenames: map[string]string{
							"DestChain":   "DestChainSelector",
							"SourceChain": "Message.Header.SourceChainSelector",
						},
					},
				},
			},
		},
		EventSyncInterval: 12 * time.Second,
		EventSyncTimeout:  10 * time.Second,
		TxSyncInterval:    12 * time.Second,
		TxSyncTimeout:     10 * time.Second,
	}, nil
}

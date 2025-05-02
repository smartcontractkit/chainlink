package aptosconfig

import (
	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
)

func GetChainReaderConfig() (chainreader.ChainReaderConfig, error) {
	return chainreader.ChainReaderConfig{
		IsLoopPlugin: true,
		Modules: map[string]*chainreader.ChainReaderModule{
			// TODO: more offramp config and other modules
			consts.ContractNameRMNRemote: {
				Name: "rmn_remote",
				Functions: map[string]*chainreader.ChainReaderFunction{
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
				Functions: map[string]*chainreader.ChainReaderFunction{
					consts.MethodNameGetARM: {
						Name: "get_arm",
					},
				},
			},
			consts.ContractNameFeeQuoter: {
				Name: "fee_quoter",
				Functions: map[string]*chainreader.ChainReaderFunction{
					consts.MethodNameFeeQuoterGetTokenPrice: {
						Name: "get_token_price",
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "token",
								Type:     "address",
								Required: true,
							},
						},
					},
					consts.MethodNameFeeQuoterGetTokenPrices: {
						Name: "get_token_prices",
						Params: []chainreader.AptosFunctionParam{
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
						Params: []chainreader.AptosFunctionParam{
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
				Functions: map[string]*chainreader.ChainReaderFunction{
					consts.MethodNameGetExecutionState: {
						Name: "get_execution_state",
						Params: []chainreader.AptosFunctionParam{
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
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "root",
								Type:     "vector<u8>",
								Required: true,
							},
						},
					},
					consts.MethodNameOffRampLatestConfigDetails: {
						Name: "latest_config_details",
						Params: []chainreader.AptosFunctionParam{
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
						// TODO: field renames
					},
					consts.MethodNameOffRampGetDynamicConfig: {
						Name: "get_dynamic_config",
						// TODO: field renames
					},
					consts.MethodNameGetSourceChainConfig: {
						Name: "get_source_chain_config",
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*chainreader.ChainReaderEvent{
					consts.EventNameExecutionStateChanged: {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "execution_state_changed_events",
						EventAccountAddress:   "offramp::get_state_address",
						EventFieldRenames: map[string]chainreader.RenamedField{
							"source_chain_selector": {
								NewName: "SourceChainSelector",
							},
							"sequence_number": {
								NewName: "SequenceNumber",
							},
							"message_id": {
								NewName: "MessageId",
							},
							"message_hash": {
								NewName: "MessageHash",
							},
							"state": {
								NewName: "State",
							},
						},
					},
					consts.EventNameCommitReportAccepted: {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "commit_report_accepted_events",
						EventAccountAddress:   "offramp::get_state_address",
						EventFieldRenames: map[string]chainreader.RenamedField{
							"blessed_merkle_roots": {
								NewName: "BlessedMerkleRoots",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"source_chain_selector": {
										NewName: "SourceChainSelector",
									},
									"on_ramp_address": {
										NewName: "OnRampAddress",
									},
									"min_seq_nr": {
										NewName: "MinSeqNr",
									},
									"max_seq_nr": {
										NewName: "MaxSeqNr",
									},
									"merkle_root": {
										NewName: "MerkleRoot",
									},
								},
							},
							"unblessed_merkle_roots": {
								NewName: "UnblessedMerkleRoots",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"source_chain_selector": {
										NewName: "SourceChainSelector",
									},
									"on_ramp_address": {
										NewName: "OnRampAddress",
									},
									"min_seq_nr": {
										NewName: "MinSeqNr",
									},
									"max_seq_nr": {
										NewName: "MaxSeqNr",
									},
									"merkle_root": {
										NewName: "MerkleRoot",
									},
								},
							},
							"price_updates": {
								NewName: "PriceUpdates",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"token_price_updates": {
										NewName: "TokenPriceUpdates",
										SubFieldRenames: map[string]chainreader.RenamedField{
											"source_token": {
												NewName: "SourceToken",
											},
											"usd_per_token": {
												NewName: "UsdPerToken",
											},
										},
									},
									"gas_price_updates": {
										NewName: "GasPriceUpdates",
										SubFieldRenames: map[string]chainreader.RenamedField{
											"dest_chain_selector": {
												NewName: "DestChainSelector",
											},
											"usd_per_unit_gas": {
												NewName: "UsdPerUnitGas",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			consts.ContractNameOnRamp: {
				Name: "onramp",
				Functions: map[string]*chainreader.ChainReaderFunction{
					consts.MethodNameOnRampGetDynamicConfig: {
						Name: "get_dynamic_config",
					},
					consts.MethodNameOnRampGetStaticConfig: {
						Name: "get_static_config",
					},
					consts.MethodNameOnRampGetDestChainConfig: {
						Name: "get_dest_chain_config",
						Params: []chainreader.AptosFunctionParam{
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
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*chainreader.ChainReaderEvent{
					consts.EventNameCCIPMessageSent: {
						EventHandleStructName: "OnRampState",
						EventHandleFieldName:  "ccip_message_sent_events",
						EventAccountAddress:   "onramp::get_state_address",
						EventFieldRenames: map[string]chainreader.RenamedField{
							"dest_chain_selector": {
								NewName:         "DestChainSelector",
								SubFieldRenames: nil,
							},
							"sequence_number": {
								NewName:         "SequenceNumber",
								SubFieldRenames: nil,
							},
							"message": {
								NewName:         "Message",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"header": {
										NewName:         "Header",
									},
									"sender": {
										NewName:         "Sender",
									},
									"data": {
										NewName:         "Data",
									},
									"receiver": {
										NewName:         "Receiver",
									},
									"extra_args": {
										NewName:         "ExtraArgs",
									},
									"fee_token": {
										NewName:         "FeeToken",
									},
									"fee_token_amount": {
										NewName:         "FeeTokenAmount",
									},
									"fee_value_juels": {
										NewName:         "FeeValueJuels",
									},
									"token_amounts": {
										NewName:         "TokenAmounts",
										SubFieldRenames: nil, // TODO
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

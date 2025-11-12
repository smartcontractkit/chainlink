package suiconfig

import (
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"

	chainreaderConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
)

func PublicKeyToAddress(pubKeyHex string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}

	flagged := append([]byte{suikey.Ed25519Scheme}, pubKeyBytes...)

	hash := blake2b.Sum256(flagged)
	return hex.EncodeToString(hash[:]), nil
}

func GetChainReaderConfig(pubKeyStr string) (map[string]any, error) {
	fromAddress, err := PublicKeyToAddress(pubKeyStr)
	if err != nil {
		return map[string]any{}, fmt.Errorf("unable to derive Sui address from public key %s: %w", pubKeyStr, err)
	}
	fromAddress = "0x" + fromAddress

	return map[string]any{
		"IsLoopPlugin": true,
		"EventsIndexer": map[string]any{
			"PollingInterval": 10 * time.Second,
			"SyncTimeout":     10 * time.Second,
		},
		"TransactionsIndexer": map[string]any{
			"PollingInterval": 10 * time.Second,
			"SyncTimeout":     10 * time.Second,
		},
		"Modules": map[string]any{
			// TODO: more offramp config and other modules
			consts.ContractNameRMNRemote: map[string]any{
				"Name": "rmn_remote",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					"GetReportDigestHeader": {
						SignerAddress: fromAddress,
						Name:          "get_report_digest_header",
					},
					"GetVersionedConfig": {
						Name:          "get_versioned_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
						},
						// ref: https://github.com/smartcontractkit/chainlink-ccip/blob/bee7c32c71cf0aec594c051fef328b4a7281a1fc/pkg/reader/ccip.go#L1440
						ResultTupleToStruct: []string{"version", "config"},
					},
					"GetCursedSubjects": {
						Name:          "get_cursed_subjects",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
						},
					},
				},
			},
			consts.ContractNameRMNProxy: map[string]any{
				"Name": "rmn_remote",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					consts.MethodNameGetARM: {
						Name:          "get_arm",
						SignerAddress: fromAddress,
					},
				},
			},
			consts.ContractNameFeeQuoter: map[string]any{
				"Name": "fee_quoter",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					"GetTokenPrice": {
						Name:          "get_token_price",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
							{
								Name:     "token",
								Type:     "address",
								Required: true,
							},
						},
					},
					"GetTokenPrices": {
						Name:          "get_token_prices",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
							{
								Name:     "tokens",
								Type:     "vector<address>",
								Required: true,
							},
						},
					},
					"GetStaticConfig": {
						Name:          "get_static_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"max_fee_juels_per_msg, link_token, token_price_staleness_threshold"},
					},
					"GetDestinationChainGasPrice": {
						Name:          "get_dest_chain_gas_price",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
			},
			"OffRamp": map[string]any{
				"Name": "offramp",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					consts.MethodNameOffRampLatestConfigDetails: {
						Name:          "latest_config_details",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "off_ramp_state_id",
								PointerTag: &codec.PointerTag{
									Module:        "offramp",
									PointerName:   "OffRampStatePointer",
									FieldName:     "off_ramp_state_id",
									DerivationKey: "OffRampState",
								},
								Type:     "object_id",
								Required: true,
							},
							{
								Name:     "ocrPluginType",
								Type:     "u8",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"ocr_config"},
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						Name:          "get_latest_price_sequence_number",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "off_ramp_state_id",
								PointerTag: &codec.PointerTag{
									Module:        "offramp",
									PointerName:   "OffRampStatePointer",
									FieldName:     "off_ramp_state_id",
									DerivationKey: "OffRampState",
								},
								Type:     "object_id",
								Required: true,
							},
						},
					},

					consts.MethodNameOffRampGetStaticConfig: {
						Name:          "get_static_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
							{
								Name: "off_ramp_state_id",
								PointerTag: &codec.PointerTag{
									Module:        "offramp",
									PointerName:   "OffRampStatePointer",
									FieldName:     "off_ramp_state_id",
									DerivationKey: "OffRampState",
								},
								Type:     "object_id",
								Required: true,
							},
						},
					},
					consts.MethodNameOffRampGetDynamicConfig: {
						Name:          "get_dynamic_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
							{
								Name: "off_ramp_state_id",
								PointerTag: &codec.PointerTag{
									Module:        "offramp",
									PointerName:   "OffRampStatePointer",
									FieldName:     "off_ramp_state_id",
									DerivationKey: "OffRampState",
								},
								Type:     "object_id",
								Required: true,
							},
						},
					},
					consts.MethodNameGetSourceChainConfig: {
						Name:          "get_source_chain_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "object_ref_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "state_object",
									PointerName:   "CCIPObjectRefPointer",
									FieldName:     "object_ref_id",
									DerivationKey: "CCIPObjectRef",
								},
								Required: true,
							},
							{
								Name: "off_ramp_state_id",
								PointerTag: &codec.PointerTag{
									Module:        "offramp",
									PointerName:   "OffRampStatePointer",
									FieldName:     "off_ramp_state_id",
									DerivationKey: "OffRampState",
								},
								Type:     "object_id",
								Required: true,
							},
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				"Events": map[string]*chainreaderConfig.ChainReaderEvent{
					"ExecutionStateChanged": {
						Name:      "offramp",
						EventType: "ExecutionStateChanged",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "offramp",
							Event:  "ExecutionStateChanged",
						},
						EventFilterRenames: map[string]string{
							"SourceChain":    "sourceChainSelector",
							"SequenceNumber": "sequenceNumber",
							"State":          "state",
						},
					},
					"CommitReportAccepted": {
						Name:      "offramp",
						EventType: "CommitReportAccepted",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "offramp",
							Event:  "CommitReportAccepted",
						},
					},
					"ConfigSet": {
						Name:      "offramp",
						EventType: "ConfigSet",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "ocr3_base",
							Event:  "ConfigSet",
						},
					},
					"SourceChainConfigSet": {
						Name:      "offramp",
						EventType: "SourceChainConfigSet",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "offramp",
							Event:  "SourceChainConfigSet",
						},
					},
				},
			},
			"OnRamp": map[string]any{
				"Name": "onramp",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					"OnRampGetDynamicConfig": {
						Name:          "get_dynamic_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "on_ramp_state_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "onramp",
									PointerName:   "OnRampStatePointer",
									FieldName:     "on_ramp_state_id",
									DerivationKey: "OnRampState",
								},
								Required: true,
							},
						},
					},
					"OnRampGetStaticConfig": {
						Name:          "get_static_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "on_ramp_state_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "onramp",
									PointerName:   "OnRampStatePointer",
									FieldName:     "on_ramp_state_id",
									DerivationKey: "OnRampState",
								},
								Required: true,
							},
						},
					},
					"OnRampGetDestChainConfig": {
						Name:          "get_dest_chain_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "on_ramp_state_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "onramp",
									PointerName:   "OnRampStatePointer",
									FieldName:     "on_ramp_state_id",
									DerivationKey: "OnRampState",
								},
								Required: true,
							},
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"sequenceNumber", "allowListEnabled", "router"},
					},
					"GetExpectedNextSequenceNumber": {
						Name:          "get_expected_next_sequence_number",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name: "on_ramp_state_id",
								Type: "object_id",
								PointerTag: &codec.PointerTag{
									Module:        "onramp",
									PointerName:   "OnRampStatePointer",
									FieldName:     "on_ramp_state_id",
									DerivationKey: "OnRampState",
								},
								Required: true,
							},
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				"Events": map[string]*chainreaderConfig.ChainReaderEvent{
					"CCIPMessageSent": {
						Name:      "CCIPMessageSent",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "onramp",
							Event:  "CCIPMessageSent",
						},
						EventFilterRenames: map[string]string{
							"SequenceNumber": "sequenceNumber",
							"DestChain":      "destChainSelector",
							"SourceChain":    "sourceChainSelector",
						},
					},
				},
			},
		},
		"EventSyncInterval": 12 * time.Second,
		"EventSyncTimeout":  10 * time.Second,
		"TxSyncInterval":    12 * time.Second,
		"TxSyncTimeout":     10 * time.Second,
	}, nil
}

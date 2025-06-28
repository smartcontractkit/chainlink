package solana

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	types "github.com/smartcontractkit/chainlink-common/pkg/types/solana"
)

var (
	ccipOffRampIDL       = idl.FetchCCIPOfframpIDL()
	ccipFeeQuoterIDL     = idl.FetchFeeQuoterIDL()
	ccipRmnRemoteIDL     = idl.FetchRMNRemoteIDL()
	ccipCCTPTokenPoolIDL = idl.FetchCctpTokenPoolIDL()
	ccipRouterIDL        = idl.FetchCCIPRouterIDL()

	// defaultCCIPLogsRetention defines the duration for which logs critical for Commit/Exec plugins processing are retained.
	// Although Exec relies on permissionlessExecThreshold which is lower than 24hours for picking eligible CommitRoots,
	// Commit still can reach to older logs because it filters them by sequence numbers. For instance, in case of RMN curse on chain,
	// we might have logs waiting in OnRamp to be committed first. When outage takes days we still would
	// be able to bring back processing without replaying any logs from chain. You can read that param as
	// "how long CCIP can be down and still be able to process all the messages after getting back to life".
	// Breaching this threshold would require replaying chain using LogPoller from the beginning of the outage.
	// Using same default retention as v1.5 https://github.com/smartcontractkit/ccip/pull/530/files
	defaultCCIPLogsRetention = 30 * 24 * time.Hour // 30 days
)

func DestContractReaderConfig() (types.ContractReader, error) {
	var offRampIDL types.IDL
	if err := json.Unmarshal([]byte(ccipOffRampIDL), &offRampIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP OffRamp IDL, error: %w", err)
	}

	var feeQuoterIDL types.IDL
	if err := json.Unmarshal([]byte(ccipFeeQuoterIDL), &feeQuoterIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Fee Quoter IDL, error: %w", err)
	}

	var rmnRemoteIDL types.IDL
	if err := json.Unmarshal([]byte(ccipRmnRemoteIDL), &rmnRemoteIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP RMN Remote IDL, error: %w", err)
	}

	var cctpTokenPoolIDL types.IDL
	if err := json.Unmarshal([]byte(ccipCCTPTokenPoolIDL), &cctpTokenPoolIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP CCTP Token Pool IDL, error: %w", err)
	}

	feeQuoterIDL.Accounts = append(feeQuoterIDL.Accounts, types.IdlTypeDef{
		Name: "USDPerToken",
		Type: types.IdlTypeDefTy{
			Kind: types.IdlTypeDefTyKindStruct,
			Fields: &types.IdlTypeDefStruct{
				{
					Name: "tokenPrices",
					Type: types.IdlType{
						AsIdlTypeVec: &types.IdlTypeVec{Vec: types.IdlType{AsIdlTypeDefined: &types.IdlTypeDefined{Defined: "TimestampedPackedU224"}}},
					},
				},
			},
		},
	})

	// Prepend custom type so it takes priority over the IDL
	offRampIDL.Types = append([]types.IdlTypeDef{{
		Name: "OnRampAddress",
		Type: types.IdlTypeDefTy{
			Kind:  types.IdlTypeDefTyKindCustom,
			Codec: "onramp_address",
		},
	}}, offRampIDL.Types...)

	var routerIDL types.IDL
	if err := json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	trueVal := true

	locationFirst := codec.ElementExtractorLocationFirst
	return types.ContractReader{
		AddressShareGroups: [][]string{{consts.ContractNameRouter, consts.ContractNameNonceManager}, {consts.ContractNameRMNRemote, consts.ContractNameRMNProxy}},
		Namespaces: map[string]types.ChainContractReader{
			consts.ContractNameOffRamp: {
				IDL: offRampIDL,
				Reads: map[string]types.ReadDefinition{
					consts.EventNameExecutionStateChanged: {
						ChainSpecificName: consts.EventNameExecutionStateChanged,
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							PollingFilter: &types.PollingFilter{
								Retention:       &defaultCCIPLogsRetention,
								IncludeReverted: &trueVal,
							},
							IndexedField0: &types.IndexedField{
								OffChainPath: consts.EventAttributeSourceChain,
								OnChainPath:  "SourceChainSelector",
							},
							IndexedField1: &types.IndexedField{
								OffChainPath: consts.EventAttributeSequenceNumber,
								OnChainPath:  consts.EventAttributeSequenceNumber,
							},
							IndexedField2: &types.IndexedField{
								OffChainPath: consts.EventAttributeState,
								OnChainPath:  consts.EventAttributeState,
							},
						},
					},
					consts.EventNameCommitReportAccepted: {
						ChainSpecificName: "CommitReportAccepted",
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							PollingFilter: &types.PollingFilter{
								Retention: &defaultCCIPLogsRetention,
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"MerkleRoot": "UnblessedMerkleRoots"}},
							&codec.ElementExtractorModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"UnblessedMerkleRoots": &locationFirst}},
						},
					},
					consts.MethodNameOffRampLatestConfigDetails: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition:     types.PDATypeDef{Prefix: []byte("config")},
						// TODO: OutputModifications are currently disabled and a special workaround is built into chainlink-solana for now
						// OutputModifications: codec.ModifiersConfig{
						// 	&codec.WrapperModifierConfig{
						// 		Fields: map[string]string{"Ocr3": "OcrConfig"},
						// 	},
						// 	&codec.PropertyExtractorConfig{FieldName: "Ocr3"},
						// 	&codec.ElementExtractorFromOnchainModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"OcrConfig": &locationFirst}},
						// 	&codec.ByteToBooleanModifierConfig{Fields: []string{"Ocrtypes.ConfigInfo.IsSignatureVerificationEnabled"}},
						// },
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						ChainSpecificName: "GlobalState",
						ReadType:          types.Account,
						PDADefinition:     types.PDATypeDef{Prefix: []byte("state")},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{
								FieldName: "LatestPriceSequenceNumber",
							},
						},
					},
					consts.MethodNameOffRampGetStaticConfig: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"SvmChainSelector": "ChainSelector",
								},
							},
						},
						MultiReader: &types.MultiReader{
							Reads: []types.ReadDefinition{
								// CCIP expects a NonceManager address, in our case that's the Router
								{
									ChainSpecificName: "ReferenceAddresses",
									ReadType:          types.Account,
									PDADefinition: types.PDATypeDef{
										Prefix: []byte("reference_addresses"),
									},
									OutputModifications: codec.ModifiersConfig{
										&codec.RenameModifierConfig{Fields: map[string]string{"Router": "NonceManager"}},
									},
								},
							},
						},
					},
					consts.MethodNameOffRampGetDynamicConfig: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{"EnableManualExecutionAfter": "PermissionLessExecutionThresholdSeconds"},
							},
							// TODO: figure out how this will be properly configured, if it has to be added to SVM state
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"IsRMNVerificationDisabled": true}},
						},
						MultiReader: &types.MultiReader{
							Reads: []types.ReadDefinition{
								{
									ChainSpecificName: "ReferenceAddresses",
									ReadType:          types.Account,
									PDADefinition: types.PDATypeDef{
										Prefix: []byte("reference_addresses"),
									},
								},
							},
						},
					},
					consts.MethodNameGetSourceChainConfig: {
						ChainSpecificName: "SourceChain",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("source_chain_state"),
							Seeds:  []types.PDASeed{{Name: "NewChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "SourceChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "Config"},
							// TODO: figure out how this will be properly configured, if it has to be added to SVM state
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"IsRMNVerificationDisabled": true}},
						},
						MultiReader: &types.MultiReader{
							ReuseParams: true,
							Reads: []types.ReadDefinition{
								{
									ChainSpecificName: "ReferenceAddresses",
									ReadType:          types.Account,
									PDADefinition: types.PDATypeDef{
										Prefix: []byte("reference_addresses"),
									},
								},
								{
									// this seems like a hack to extract both State and Config fields?
									ChainSpecificName: "SourceChain",
									ReadType:          types.Account,
									PDADefinition: types.PDATypeDef{
										Prefix: []byte("source_chain_state"),
										Seeds:  []types.PDASeed{{Name: "NewChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
									},
									InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "SourceChainSelector"}}},
									OutputModifications: codec.ModifiersConfig{
										&codec.PropertyExtractorConfig{FieldName: "State"},
									},
								},
							},
						},
					},
				},
			},
			consts.ContractNameFeeQuoter: {
				IDL: feeQuoterIDL,
				Reads: map[string]types.ReadDefinition{
					consts.MethodNameFeeQuoterGetStaticConfig: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"MaxFeeJuelsPerMsg": "MaxFeeJuelsPerMsg",
									"LinkTokenMint":     "LinkToken",
								},
							},
						},
					},
					// This one is hacky, but works - [NONEVM-1320]
					consts.MethodNameFeeQuoterGetTokenPrices: {
						ChainSpecificName: "USDPerToken",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []types.PDASeed{
								{
									Name: "Tokens",
									Type: types.IdlType{
										AsIdlTypeVec: &types.IdlTypeVec{
											Vec: types.IdlType{AsString: types.IdlTypePublicKey},
										},
									},
								},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "TokenPrices"},
						},
					},
					consts.MethodNameFeeQuoterGetTokenPrice: {
						ChainSpecificName: "BillingTokenConfigWrapper",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []types.PDASeed{{
								Name: "Token",
								Type: types.IdlType{AsString: types.IdlTypePublicKey},
							}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "types.UsdPerToken"},
						},
					},
					consts.MethodNameGetFeePriceUpdate: {
						ChainSpecificName: "DestChain",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []types.PDASeed{{Name: "DestinationChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "State.UsdPerUnitGas"}},
					},
					consts.MethodNameGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []types.PDASeed{{Name: "DestinationChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "Config"},
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"DefaultTokenFeeUsdcents": "DefaultTokenFeeUSDCents",
									"NetworkFeeUsdcents":      "NetworkFeeUSDCents",
								},
							},
						},
						MultiReader: &types.MultiReader{
							ReuseParams: true,
							Reads: []types.ReadDefinition{
								{
									// this seems like a hack to extract both State and Config fields?
									ChainSpecificName: "DestChain",
									PDADefinition: types.PDATypeDef{
										Prefix: []byte("dest_chain"),
										Seeds:  []types.PDASeed{{Name: "DestinationChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
									},
									InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
									OutputModifications: codec.ModifiersConfig{
										&codec.PropertyExtractorConfig{FieldName: "State"},
									},
								},
							},
						},
					},
				},
			},
			consts.ContractNameRouter: {
				IDL: routerIDL,
				Reads: map[string]types.ReadDefinition{
					// TODO: PDA fetching is unnecessary here
					consts.MethodNameRouterGetWrappedNative: {
						ChainSpecificName: "Config",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"WrappedNative": solana.WrappedSol.String()}},
							&codec.PropertyExtractorConfig{FieldName: "WrappedNative"},
							// TODO: error: process Router results: get router wrapped native result: invalid type: '': source data must be an array or slice, got string"
						},
					},
				},
			},
			consts.ContractNameNonceManager: {
				IDL: routerIDL,
				Reads: map[string]types.ReadDefinition{
					consts.MethodNameGetInboundNonce: {
						ChainSpecificName: "Nonce",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("nonce"),
							Seeds: []types.PDASeed{
								{Name: "DestinationChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}},
								{Name: "Authority", Type: types.IdlType{AsString: types.IdlTypePublicKey}},
							},
						},
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{
								"DestinationChainSelector": "SourceChainSelector",
								"Authority":                "Sender",
							}}},
					},
				},
			},
			consts.ContractNameRMNProxy: {
				IDL: rmnRemoteIDL,
				Reads: map[string]types.ReadDefinition{
					consts.MethodNameGetARM: {
						// TODO: need to have definition or it'll complain
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							// create a field to extract it
							&codec.HardCodeModifierConfig{
								OffChainValues: map[string]any{"RmnRemoteAddress": ""},
							},
							&codec.PropertyExtractorConfig{
								FieldName: "RmnRemoteAddress",
							},
						},
						ResponseAddressHardCoder: &codec.HardCodeModifierConfig{
							// type doesn't matter it will be overridden with address internally, key is "" because it's a primitive value and not a field
							OffChainValues: map[string]any{"": ""},
						},
					},
				},
			},
			consts.ContractNameRMNRemote: {
				IDL: rmnRemoteIDL,
				Reads: map[string]types.ReadDefinition{
					consts.MethodNameGetVersionedConfig: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							// Disable fields so config isn't used, we only use global verification
							&codec.DropModifierConfig{
								Fields: []string{"Version"},
							},
						},
					},
					consts.MethodNameGetReportDigestHeader: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
					},
					consts.MethodNameGetCursedSubjects: {
						ChainSpecificName: "Curses",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("curses"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{
								FieldName: "CursedSubjects.Value",
							},
							&codec.WrapperModifierConfig{
								Fields: map[string]string{"": "CursedSubjects"},
							},
						},
					},
				},
			},
			consts.ContractNameUSDCTokenPool: {
				IDL: cctpTokenPoolIDL,
				Reads: map[string]types.ReadDefinition{
					consts.EventNameCCTPMessageSent: {
						ChainSpecificName: "CcipCctpMessageSentEvent",
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							PollingFilter: &types.PollingFilter{
								Retention: &defaultCCIPLogsRetention,
							},
							IndexedField0: &types.IndexedField{
								OffChainPath: consts.EventAttributeCCTPNonce,
								OnChainPath:  "CctpNonce",
							},
							IndexedField1: &types.IndexedField{
								OffChainPath: consts.EventAttributeSourceDomain,
								OnChainPath:  "SourceDomain",
							},
						},
					},
				},
			},
		},
	}, nil
}

func SourceContractReaderConfig() (types.ContractReader, error) {
	var routerIDL types.IDL
	if err := json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	var feeQuoterIDL types.IDL
	if err := json.Unmarshal([]byte(ccipFeeQuoterIDL), &feeQuoterIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Fee Quoter IDL, error: %w", err)
	}

	var cctpTokenPoolIDL types.IDL
	if err := json.Unmarshal([]byte(ccipCCTPTokenPoolIDL), &cctpTokenPoolIDL); err != nil {
		return types.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP CCTP Token Pool IDL, error: %w", err)
	}

	feeQuoterIDL.Accounts = append(feeQuoterIDL.Accounts, types.IdlTypeDef{
		Name: "USDPerToken",
		Type: types.IdlTypeDefTy{
			Kind: types.IdlTypeDefTyKindStruct,
			Fields: &types.IdlTypeDefStruct{
				{
					Name: "tokenPrices",
					Type: types.IdlType{
						AsIdlTypeVec: &types.IdlTypeVec{Vec: types.IdlType{AsIdlTypeDefined: &types.IdlTypeDefined{Defined: "TimestampedPackedU224"}}},
					},
				},
			},
		},
	})

	// Prepend custom type so it takes priority over the IDL
	routerIDL.Types = append([]types.IdlTypeDef{{
		Name: "CrossChainAmount",
		Type: types.IdlTypeDefTy{
			Kind:  types.IdlTypeDefTyKindCustom,
			Codec: "cross_chain_amount",
		},
	}}, routerIDL.Types...)

	return types.ContractReader{
		AddressShareGroups: [][]string{{consts.ContractNameRouter, consts.ContractNameOnRamp}},
		Namespaces: map[string]types.ChainContractReader{
			consts.ContractNameOnRamp: {
				IDL: routerIDL,
				Reads: map[string]types.ReadDefinition{
					consts.MethodNameGetExpectedNextSequenceNumber: {
						ChainSpecificName: "DestChain",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("dest_chain_state"),
							Seeds:  []types.PDASeed{{Name: "NewChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "State"},
							&codec.RenameModifierConfig{
								Fields: map[string]string{"SequenceNumber": "ExpectedNextSequenceNumber"},
							}},
					},
					consts.EventNameCCIPMessageSent: {
						ChainSpecificName: "CCIPMessageSent",
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							PollingFilter: &types.PollingFilter{
								Retention: &defaultCCIPLogsRetention,
							},
							IndexedField0: &types.IndexedField{
								OffChainPath: consts.EventAttributeSourceChain,
								OnChainPath:  "Message.Header.SourceChainSelector",
							},
							IndexedField1: &types.IndexedField{
								OffChainPath: consts.EventAttributeDestChain,
								OnChainPath:  "Message.Header.DestChainSelector",
							},
							IndexedField2: &types.IndexedField{
								OffChainPath: consts.EventAttributeSequenceNumber,
								OnChainPath:  "Message.Header.SequenceNumber",
							},
						},
					},
					consts.MethodNameOnRampGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("dest_chain_state"),
							Seeds:  []types.PDASeed{{Name: "NewChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						// response Router field will be populated with the bound address of the onramp
						ResponseAddressHardCoder: &codec.HardCodeModifierConfig{
							// type doesn't matter it will be overridden with address internally
							OffChainValues: map[string]any{"Router": ""},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "State"},
							&codec.RenameModifierConfig{
								Fields: map[string]string{"SequenceNumber": "ExpectedNextSequenceNumber"},
							},
						},
						MultiReader: &types.MultiReader{
							ReuseParams: true,
							Reads: []types.ReadDefinition{
								// this seems like a hack to extract both State and Config fields?
								{
									ChainSpecificName: "DestChain",
									ReadType:          types.Account,
									PDADefinition: types.PDATypeDef{
										Prefix: []byte("dest_chain_state"),
										Seeds:  []types.PDASeed{{Name: "NewChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
									},
									InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "DestChainSelector"}}},
									OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "Config"}},
								},
							},
						},
					},
					consts.MethodNameOnRampGetDynamicConfig: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition:     types.PDATypeDef{Prefix: []byte("config")},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{"Owner": "AllowListAdmin"},
							},
							// for some reason CCIP reader expects the data to be wrapped under DynamicConfig, but not on offramp...
							&codec.WrapperModifierConfig{
								Fields: map[string]string{"": "DynamicConfig"},
							},
						},
					},
				},
			},
			consts.ContractNameFeeQuoter: {
				IDL: feeQuoterIDL,
				Reads: map[string]types.ReadDefinition{
					consts.MethodNameFeeQuoterGetStaticConfig: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"MaxFeeJuelsPerMsg": "MaxFeeJuelsPerMsg",
									"LinkTokenMint":     "LinkToken",
								},
							},
						},
					},
					// this one is hacky, but should work NONEVM-1320
					consts.MethodNameFeeQuoterGetTokenPrices: {
						ChainSpecificName: "USDPerToken",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []types.PDASeed{
								{
									Name: "Tokens",
									Type: types.IdlType{
										AsIdlTypeVec: &types.IdlTypeVec{
											Vec: types.IdlType{AsString: types.IdlTypePublicKey},
										},
									},
								},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "TokenPrices"},
						},
					},
					consts.MethodNameFeeQuoterGetTokenPrice: {
						ChainSpecificName: "BillingTokenConfigWrapper",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []types.PDASeed{{
								Name: "Token",
								Type: types.IdlType{AsString: types.IdlTypePublicKey},
							}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "types.UsdPerToken"},
						},
					},
					consts.MethodNameGetFeePriceUpdate: {
						ChainSpecificName: "DestChain",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []types.PDASeed{{Name: "DestinationChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "State.UsdPerUnitGas"}},
					},
					consts.MethodNameGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []types.PDASeed{{Name: "DestinationChainSelector", Type: types.IdlType{AsString: types.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "Config"},
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"DefaultTokenFeeUsdcents": "DefaultTokenFeeUSDCents",
									"NetworkFeeUsdcents":      "NetworkFeeUSDCents",
								},
							},
						},
					},
				},
			},
			consts.ContractNameRouter: {
				IDL: routerIDL,
				Reads: map[string]types.ReadDefinition{
					// TODO: PDA fetching is unnecessary here
					consts.MethodNameRouterGetWrappedNative: {
						ChainSpecificName: "Config",
						ReadType:          types.Account,
						PDADefinition: types.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"WrappedNative": solana.WrappedSol.String()}},
							&codec.PropertyExtractorConfig{FieldName: "WrappedNative"},
							// TODO: error: process Router results: get router wrapped native result: invalid type: '': source data must be an array or slice, got string"
						},
					},
				},
			},
			consts.ContractNameUSDCTokenPool: {
				IDL: cctpTokenPoolIDL,
				Reads: map[string]types.ReadDefinition{
					consts.EventNameCCTPMessageSent: {
						ChainSpecificName: "CcipCctpMessageSentEvent",
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							PollingFilter: &types.PollingFilter{
								Retention: &defaultCCIPLogsRetention,
							},
							IndexedField0: &types.IndexedField{
								OffChainPath: consts.EventAttributeCCTPNonce,
								OnChainPath:  "CctpNonce",
							},
							IndexedField1: &types.IndexedField{
								OffChainPath: consts.EventAttributeSourceDomain,
								OnChainPath:  "SourceDomain",
							},
						},
					},
				},
			},
		},
	}, nil
}

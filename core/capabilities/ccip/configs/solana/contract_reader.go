package solana

import (
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	solanacodec "github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var ccipOffRampIDL = idl.FetchCCIPOfframpIDL()
var ccipFeeQuoterIDL = idl.FetchFeeQuoterIDL()
var ccipRmnRemoteIDL = idl.FetchRMNRemoteIDL()

func DestContractReaderConfig() (config.ContractReader, error) {
	var offRampIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipOffRampIDL), &offRampIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP OffRamp IDL, error: %w", err)
	}

	var feeQuoterIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipFeeQuoterIDL), &feeQuoterIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Fee Quoter IDL, error: %w", err)
	}

	var rmnRemoteIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipRmnRemoteIDL), &rmnRemoteIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP RMN Remote IDL, error: %w", err)
	}

	feeQuoterIDL.Accounts = append(feeQuoterIDL.Accounts, solanacodec.IdlTypeDef{
		Name: "USDPerToken",
		Type: solanacodec.IdlTypeDefTy{
			Kind: solanacodec.IdlTypeDefTyKindStruct,
			Fields: &solanacodec.IdlTypeDefStruct{
				{
					Name: "tokenPrices",
					Type: solanacodec.IdlType{
						AsIdlTypeVec: &solanacodec.IdlTypeVec{Vec: solanacodec.IdlType{AsIdlTypeDefined: &solanacodec.IdlTypeDefined{Defined: "TimestampedPackedU224"}}},
					},
				},
			},
		},
	})

	// Prepend custom type so it takes priority over the IDL
	offRampIDL.Types = append([]solanacodec.IdlTypeDef{{
		Name: "OnRampAddress",
		Type: solanacodec.IdlTypeDefTy{
			Kind:  solanacodec.IdlTypeDefTyKindCustom,
			Codec: "onramp_address",
		},
	}}, offRampIDL.Types...)

	var routerIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	trueVal := true

	locationFirst := codec.ElementExtractorLocationFirst
	return config.ContractReader{
		AddressShareGroups: [][]string{{consts.ContractNameRouter, consts.ContractNameNonceManager}, {consts.ContractNameRMNRemote, consts.ContractNameRMNProxy}},
		Namespaces: map[string]config.ChainContractReader{
			consts.ContractNameOffRamp: {
				IDL: offRampIDL,
				Reads: map[string]config.ReadDefinition{
					consts.EventNameExecutionStateChanged: {
						ChainSpecificName: consts.EventNameExecutionStateChanged,
						ReadType:          config.Event,
						EventDefinitions: &config.EventDefinitions{
							PollingFilter: &config.PollingFilter{
								IncludeReverted: &trueVal,
							},
							IndexedField0: &config.IndexedField{
								OffChainPath: consts.EventAttributeSourceChain,
								OnChainPath:  "SourceChainSelector",
							},
							IndexedField1: &config.IndexedField{
								OffChainPath: consts.EventAttributeSequenceNumber,
								OnChainPath:  consts.EventAttributeSequenceNumber,
							},
							IndexedField2: &config.IndexedField{
								OffChainPath: consts.EventAttributeState,
								OnChainPath:  consts.EventAttributeState,
							},
						},
					},
					consts.EventNameCommitReportAccepted: {
						ChainSpecificName: "CommitReportAccepted",
						ReadType:          config.Event,
						EventDefinitions: &config.EventDefinitions{
							PollingFilter: &config.PollingFilter{},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"MerkleRoot": "UnblessedMerkleRoots"}},
							&codec.ElementExtractorModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"UnblessedMerkleRoots": &locationFirst}},
						},
					},
					consts.MethodNameOffRampLatestConfigDetails: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition:     solanacodec.PDATypeDef{Prefix: []byte("config")},
						// TODO: OutputModifications are currently disabled and a special workaround is built into chainlink-solana for now
						// OutputModifications: codec.ModifiersConfig{
						// 	&codec.WrapperModifierConfig{
						// 		Fields: map[string]string{"Ocr3": "OcrConfig"},
						// 	},
						// 	&codec.PropertyExtractorConfig{FieldName: "Ocr3"},
						// 	&codec.ElementExtractorFromOnchainModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"OcrConfig": &locationFirst}},
						// 	&codec.ByteToBooleanModifierConfig{Fields: []string{"OcrConfig.ConfigInfo.IsSignatureVerificationEnabled"}},
						// },
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						ChainSpecificName: "GlobalState",
						ReadType:          config.Account,
						PDADefinition:     solanacodec.PDATypeDef{Prefix: []byte("state")},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{
								FieldName: "LatestPriceSequenceNumber",
							},
						},
					},
					consts.MethodNameOffRampGetStaticConfig: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"SvmChainSelector": "ChainSelector",
								},
							},
						},
						MultiReader: &config.MultiReader{
							Reads: []config.ReadDefinition{
								// CCIP expects a NonceManager address, in our case that's the Router
								{
									ChainSpecificName: "ReferenceAddresses",
									ReadType:          config.Account,
									PDADefinition: solanacodec.PDATypeDef{
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
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{"EnableManualExecutionAfter": "PermissionLessExecutionThresholdSeconds"},
							},
							// TODO: figure out how this will be properly configured, if it has to be added to SVM state
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"IsRMNVerificationDisabled": true}},
						},
						MultiReader: &config.MultiReader{
							Reads: []config.ReadDefinition{
								{
									ChainSpecificName: "ReferenceAddresses",
									ReadType:          config.Account,
									PDADefinition: solanacodec.PDATypeDef{
										Prefix: []byte("reference_addresses"),
									},
								},
							},
						},
					},
					consts.MethodNameGetSourceChainConfig: {
						ChainSpecificName: "SourceChain",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("source_chain_state"),
							Seeds:  []solanacodec.PDASeed{{Name: "NewChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "SourceChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "Config"},
							// TODO: figure out how this will be properly configured, if it has to be added to SVM state
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"IsRMNVerificationDisabled": true}},
						},
						MultiReader: &config.MultiReader{
							ReuseParams: true,
							Reads: []config.ReadDefinition{
								{
									ChainSpecificName: "ReferenceAddresses",
									ReadType:          config.Account,
									PDADefinition: solanacodec.PDATypeDef{
										Prefix: []byte("reference_addresses"),
									},
								},
								{
									// this seems like a hack to extract both State and Config fields?
									ChainSpecificName: "SourceChain",
									ReadType:          config.Account,
									PDADefinition: solanacodec.PDATypeDef{
										Prefix: []byte("source_chain_state"),
										Seeds:  []solanacodec.PDASeed{{Name: "NewChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
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
				Reads: map[string]config.ReadDefinition{
					consts.MethodNameFeeQuoterGetStaticConfig: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
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
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []solanacodec.PDASeed{
								{
									Name: "Tokens",
									Type: solanacodec.IdlType{
										AsIdlTypeVec: &solanacodec.IdlTypeVec{
											Vec: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey},
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
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []solanacodec.PDASeed{{
								Name: "Token",
								Type: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey},
							}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "Config.UsdPerToken"},
						},
					},
					consts.MethodNameGetFeePriceUpdate: {
						ChainSpecificName: "DestChain",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
						},
						InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "State.UsdPerUnitGas"}},
					},
					consts.MethodNameGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
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
						MultiReader: &config.MultiReader{
							ReuseParams: true,
							Reads: []config.ReadDefinition{
								{
									// this seems like a hack to extract both State and Config fields?
									ChainSpecificName: "DestChain",
									PDADefinition: solanacodec.PDATypeDef{
										Prefix: []byte("dest_chain"),
										Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
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
				Reads: map[string]config.ReadDefinition{
					// TODO: PDA fetching is unnecessary here
					consts.MethodNameRouterGetWrappedNative: {
						ChainSpecificName: "Config",
						PDADefinition: solanacodec.PDATypeDef{
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
				Reads: map[string]config.ReadDefinition{
					consts.MethodNameGetInboundNonce: {
						ChainSpecificName: "Nonce",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("nonce"),
							Seeds: []solanacodec.PDASeed{
								{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}},
								{Name: "Authority", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey}},
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
				Reads: map[string]config.ReadDefinition{
					consts.MethodNameGetARM: {
						// TODO: need to have definition or it'll complain
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
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
				Reads: map[string]config.ReadDefinition{
					consts.MethodNameGetVersionedConfig: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
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
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("config"),
						},
					},
					consts.MethodNameGetCursedSubjects: {
						ChainSpecificName: "Curses",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
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
		},
	}, nil
}

func SourceContractReaderConfig() (config.ContractReader, error) {
	var routerIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	var feeQuoterIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipFeeQuoterIDL), &feeQuoterIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Fee Quoter IDL, error: %w", err)
	}

	feeQuoterIDL.Accounts = append(feeQuoterIDL.Accounts, solanacodec.IdlTypeDef{
		Name: "USDPerToken",
		Type: solanacodec.IdlTypeDefTy{
			Kind: solanacodec.IdlTypeDefTyKindStruct,
			Fields: &solanacodec.IdlTypeDefStruct{
				{
					Name: "tokenPrices",
					Type: solanacodec.IdlType{
						AsIdlTypeVec: &solanacodec.IdlTypeVec{Vec: solanacodec.IdlType{AsIdlTypeDefined: &solanacodec.IdlTypeDefined{Defined: "TimestampedPackedU224"}}},
					},
				},
			},
		},
	})

	// Prepend custom type so it takes priority over the IDL
	routerIDL.Types = append([]solanacodec.IdlTypeDef{{
		Name: "CrossChainAmount",
		Type: solanacodec.IdlTypeDefTy{
			Kind:  solanacodec.IdlTypeDefTyKindCustom,
			Codec: "cross_chain_amount",
		},
	}}, routerIDL.Types...)

	return config.ContractReader{
		AddressShareGroups: [][]string{{consts.ContractNameRouter, consts.ContractNameOnRamp}},
		Namespaces: map[string]config.ChainContractReader{
			consts.ContractNameOnRamp: {
				IDL: routerIDL,
				Reads: map[string]config.ReadDefinition{
					consts.MethodNameGetExpectedNextSequenceNumber: {
						ChainSpecificName: "DestChain",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain_state"),
							Seeds:  []solanacodec.PDASeed{{Name: "NewChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
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
						ReadType:          config.Event,
						EventDefinitions: &config.EventDefinitions{
							PollingFilter: &config.PollingFilter{},
							IndexedField0: &config.IndexedField{
								OffChainPath: consts.EventAttributeSourceChain,
								OnChainPath:  "Message.Header.SourceChainSelector",
							},
							IndexedField1: &config.IndexedField{
								OffChainPath: consts.EventAttributeDestChain,
								OnChainPath:  "Message.Header.DestChainSelector",
							},
							IndexedField2: &config.IndexedField{
								OffChainPath: consts.EventAttributeSequenceNumber,
								OnChainPath:  "Message.Header.SequenceNumber",
							},
						},
					},
					consts.MethodNameOnRampGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain_state"),
							Seeds:  []solanacodec.PDASeed{{Name: "NewChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
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
						MultiReader: &config.MultiReader{
							ReuseParams: true,
							Reads: []config.ReadDefinition{
								// this seems like a hack to extract both State and Config fields?
								{
									ChainSpecificName: "DestChain",
									ReadType:          config.Account,
									PDADefinition: solanacodec.PDATypeDef{
										Prefix: []byte("dest_chain_state"),
										Seeds:  []solanacodec.PDASeed{{Name: "NewChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
									},
									InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"NewChainSelector": "DestChainSelector"}}},
									OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "Config"}},
								},
							},
						},
					},
					consts.MethodNameOnRampGetDynamicConfig: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition:     solanacodec.PDATypeDef{Prefix: []byte("config")},
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
				Reads: map[string]config.ReadDefinition{
					consts.MethodNameFeeQuoterGetStaticConfig: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
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
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []solanacodec.PDASeed{
								{
									Name: "Tokens",
									Type: solanacodec.IdlType{
										AsIdlTypeVec: &solanacodec.IdlTypeVec{
											Vec: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey},
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
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("fee_billing_token_config"),
							Seeds: []solanacodec.PDASeed{{
								Name: "Token",
								Type: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey},
							}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "Config.UsdPerToken"},
						},
					},
					consts.MethodNameGetFeePriceUpdate: {
						ChainSpecificName: "DestChain",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
						},
						InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestinationChainSelector": "DestChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "State.UsdPerUnitGas"}},
					},
					consts.MethodNameGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
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
				Reads: map[string]config.ReadDefinition{
					// TODO: PDA fetching is unnecessary here
					consts.MethodNameRouterGetWrappedNative: {
						ChainSpecificName: "Config",
						ReadType:          config.Account,
						PDADefinition: solanacodec.PDATypeDef{
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
		},
	}, nil
}

func MergeReaderConfigs(configs ...config.ContractReader) config.ContractReader {
	allNamespaces := make(map[string]config.ChainContractReader)
	for _, c := range configs {
		for namespace, method := range c.Namespaces {
			allNamespaces[namespace] = method
		}
	}

	return config.ContractReader{Namespaces: allNamespaces}
}

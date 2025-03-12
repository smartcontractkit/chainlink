package solana

import (
	"encoding/json"
	"fmt"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	solanacodec "github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var ccipOffRampIDL = idl.FetchCCIPOfframpIDL()
var ccipFeeQuoterIDL = idl.FetchFeeQuoterIDL()

// TODO add events when Querying is finished
func DestContractReaderConfig() (config.ContractReader, error) {
	var offRampIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipOffRampIDL), &offRampIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP OffRamp IDL, error: %w", err)
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

	var routerIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	locationFirst := codec.ElementExtractorLocationFirst
	return config.ContractReader{
		AddressShareGroups: [][]string{{consts.ContractNameRouter, consts.ContractNameNonceManager}},
		Namespaces: map[string]config.ChainContractReader{
			consts.ContractNameOffRamp: {
				IDL: offRampIDL,
				Reads: map[string]config.ReadDefinition{
					consts.EventNameExecutionStateChanged: {
						ChainSpecificName: consts.EventNameExecutionStateChanged,
						ReadType:          config.Event,
						EventDefinitions: &config.EventDefinitions{
							PollingFilter: &config.PollingFilter{},
							IndexedField0: &config.IndexedField{
								OffChainPath: consts.EventAttributeSourceChain,
								OnChainPath:  consts.EventAttributeSourceChain,
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
						OutputModifications: codec.ModifiersConfig{
							// TODO why does Solana have two of these in an array, but EVM has one
							&codec.WrapperModifierConfig{
								Fields: map[string]string{"Ocr3": "OcrConfig"},
							},
							&codec.PropertyExtractorConfig{FieldName: "Ocr3"},
							&codec.ElementExtractorFromOnchainModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"OcrConfig": &locationFirst}},
							&codec.ByteToBooleanModifierConfig{Fields: []string{"OcrConfig.ConfigInfo.IsSignatureVerificationEnabled"}},
						},
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						ChainSpecificName: "GlobalState",
						ReadType:          config.Account,
						PDADefinition:     solanacodec.PDATypeDef{Prefix: []byte("state")},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{"LatestPriceSequenceNumber": "LatestSeqNr"},
							}},
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
							// TODO, onchain comment cays that both can be populated, but EVM contracts only have 1, so we take first here
							//	  // OnRamp addresses supported from the source chain, each of them has a 64 byte address. So this can hold 2 addresses.
							//    // If only one address is configured, then the space for the second address must be zeroed.
							//    // Each address must be right padded with zeros if it is less than 64 bytes.
							&codec.ElementExtractorModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"OnRamp": &locationFirst}},
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
								Name: "Tokens",
								Type: solanacodec.IdlType{
									AsIdlTypeVec: &solanacodec.IdlTypeVec{
										Vec: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey},
									},
								},
							}}},
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
					consts.MethodNameRouterGetWrappedNative: {
						ChainSpecificName: "Config",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"LinkTokenMint": "LinkToken",
								},
							},
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
		},
	}, nil
}

// TODO add events when Querying is finished
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
							}},
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
								Name: "Tokens",
								Type: solanacodec.IdlType{
									AsIdlTypeVec: &solanacodec.IdlTypeVec{
										Vec: solanacodec.IdlType{AsString: solanacodec.IdlTypePublicKey},
									},
								},
							}}},
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
					consts.MethodNameRouterGetWrappedNative: {
						ChainSpecificName: "Config",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("config"),
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{
									"LinkTokenMint": "LinkToken",
								},
							},
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

package solana

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"

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
	type TimestampedUnixBig struct {
		Value     *big.Int `json:"value"`
		Timestamp uint32   `json:"timestamp"`
	}

	var offRampIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipOffRampIDL), &offRampIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP OffRamp IDL, error: %w", err)
	}

	var feeQuoterIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipFeeQuoterIDL), &feeQuoterIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Fee Quoter IDL, error: %w", err)
	}

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
							&codec.HardCodeModifierConfig{
								OnChainValues: map[string]any{
									// Also doesn't exist on Solana?
									"NonceManager": solana.PublicKey{},
									// I don't see this being used in the ccip plugin, but it's still worth figuring out what the equivalent is
									"TokenAdminRegistry": solana.PublicKey{},
									// these two don't exist on Solana
									"GasForCallExactCheck": 0,
									"RmnRemote":            solana.PublicKey{},
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
							&codec.HardCodeModifierConfig{
								OnChainValues: map[string]any{
									"IsRMNVerificationDisabled": false,
									"MessageInterceptor":        solana.PublicKey{},
								},
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
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"SourceChainSelector": "NewChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{
							&codec.PropertyExtractorConfig{FieldName: "SourceChainConfig"},
							// TODO, onchain comment cays that both can be populated, but EVM contracts only have 1, so we take first here
							//	  // OnRamp addresses supported from the source chain, each of them has a 64 byte address. So this can hold 2 addresses.
							//    // If only one address is configured, then the space for the second address must be zeroed.
							//    // Each address must be right padded with zeros if it is less than 64 bytes.
							&codec.ElementExtractorModifierConfig{Extractions: map[string]*codec.ElementExtractorLocation{"OnRamp": &locationFirst}},
						},
						// TODO enable multi reader for batch (should be a really small simple feature) https://github.com/smartcontractkit/chainlink-solana/pull/1070
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
				IDL: solanacodec.IDL{},
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
							&codec.HardCodeModifierConfig{
								// This comment is from onchain program
								// The following field is unused until the day we integrate with feeds to fetch fresh values
								OnChainValues: map[string]any{"StalenessThreshold": 0},
							},
						},
					},
					// This one is hacky, but works - [NONEVM-1320]
					consts.MethodNameFeeQuoterGetTokenPrices: {
						ChainSpecificName: "BillingTokenConfigWrapper",
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
							&codec.DropModifierConfig{
								Fields: []string{"Config"},
							},
							&codec.HardCodeModifierConfig{
								OffChainValues: map[string]any{
									"Response": make([]TimestampedUnixBig, 1000),
								},
							},
							&codec.PropertyExtractorConfig{FieldName: "Response"},
						},
						ReadType: config.Account,
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
						InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "DestinationChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "State.UsdPerUnitGas"}},
					},
					consts.MethodNameGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "DestinationChainSelector"}}},
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
								"SourceChainSelector": "DestinationChainSelector",
								"Sender":              "Authority",
							}}},
					},
				},
			},
		},
	}, nil
}

// TODO add events when Querying is finished
func SourceContractReaderConfig() (config.ContractReader, error) {
	type TimestampedUnixBig struct {
		Value     *big.Int `json:"value"`
		Timestamp uint32   `json:"timestamp"`
	}

	var routerIDL solanacodec.IDL
	if err := json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return config.ContractReader{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

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
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "NewChainSelector"}}},
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
							OffChainValues: map[string]any{"Router": solana.PublicKey{}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "NewChainSelector"}}},
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
									InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "NewChainSelector"}}},
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
							&codec.HardCodeModifierConfig{
								OnChainValues: map[string]any{
									// don't exist on Solana, just leave empty
									"ReentrancyGuardEntered": solana.PublicKey{},
									"MessageInterceptor":     solana.PublicKey{},
								},
							}},
					},
				},
			},
			consts.ContractNameFeeQuoter: {
				IDL: solanacodec.IDL{},
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
							&codec.HardCodeModifierConfig{
								//  comment from onchain
								// The following field is unused until the day we integrate with feeds to fetch fresh values
								OnChainValues: map[string]any{"StalenessThreshold": 0},
							},
						},
					},
					// TODO this one is hacky, but should work NONEVM-1320
					consts.MethodNameFeeQuoterGetTokenPrices: {
						ChainSpecificName: "BillingTokenConfigWrapper",
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
							&codec.DropModifierConfig{
								Fields: []string{"Config"},
							},
							&codec.HardCodeModifierConfig{
								OffChainValues: map[string]any{
									"Response": make([]TimestampedUnixBig, 1000),
								},
							},
							&codec.PropertyExtractorConfig{FieldName: "Response"},
						},
						ReadType: config.Account,
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
						InputModifications:  codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "DestinationChainSelector"}}},
						OutputModifications: codec.ModifiersConfig{&codec.PropertyExtractorConfig{FieldName: "State.UsdPerUnitGas"}},
					},
					consts.MethodNameGetDestChainConfig: {
						ChainSpecificName: "DestChain",
						PDADefinition: solanacodec.PDATypeDef{
							Prefix: []byte("dest_chain"),
							Seeds:  []solanacodec.PDASeed{{Name: "DestinationChainSelector", Type: solanacodec.IdlType{AsString: solanacodec.IdlTypeU64}}},
						},
						InputModifications: codec.ModifiersConfig{&codec.RenameModifierConfig{Fields: map[string]string{"DestChainSelector": "DestinationChainSelector"}}},
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

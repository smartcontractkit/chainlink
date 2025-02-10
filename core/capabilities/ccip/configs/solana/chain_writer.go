package solana

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	solanacodec "github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
)

// TODO: Replace with call for offramp once on-chain changes are merged
var ccipOfframpIDL string
var ccipRouterIDL = idl.FetchCCIPRouterIDL()

const (
	sourceChainSelectorPath = "Info.AbstractReports.Messages.Header.SourceChainSelector"
	destChainSelectorPath   = "Info.AbstractReports.Messages.Header.DestChainSelector"
	destTokenAddress        = "Info.AbstractReports.Messages.TokenAmounts.DestTokenAddress"
	merkleRootChainSelector = "Info.MerkleRoots.ChainSel"
	merkleRoot              = "Info.MerkleRoots.MerkleRoot"
)

func getCommitMethodConfig(fromAddress string, offrampProgramAddress string) chainwriter.MethodConfig {
	return chainwriter.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawReport": "Report"},
			},
		},
		ChainSpecificName: "commit",
		LookupTables: chainwriter.LookupTables{
			DerivedLookupTables: []chainwriter.DerivedLookupTable{
				getCommonAddressLookupTableConfig(offrampProgramAddress),
			},
		},
		Accounts: []chainwriter.Lookup{
			getOfframpAccountConfig(offrampProgramAddress),
			getReferenceAddressesConfig(offrampProgramAddress),
			chainwriter.PDALookups{
				Name:      "SourceChainState",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			chainwriter.PDALookups{
				Name:      "CommitReport",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("commit_report")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRoot}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			getSysVarInstructionConstant(),
			getFeeBillingSignerConfig(offrampProgramAddress),
			getFeeQuoterConfig(offrampProgramAddress),
			chainwriter.PDALookups{
				Name: "FeeQuoterAllowedPriceUpdater",
				// Fetch fee quoter public key to use as program ID for PDA
				PublicKey: getFeeQuoterConfig(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("allowed_price_updater")},
					{Dynamic: getFeeBillingSignerConfig(offrampProgramAddress)},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "FeeQuoterConfig",
				// Fetch fee quoter public key to use as program ID for PDA
				PublicKey: getFeeQuoterConfig(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("config")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "GlobalState",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("state")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "BillingTokenConfig",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("fee_billing_token_config")},
					{Dynamic: chainwriter.AccountLookup{Location: "Info.TokenPrices.TokenID"}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "ChainConfigGasPrice",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("dest_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
		},
		DebugIDLocation: "",
	}
}

func getExecuteMethodConfig(fromAddress string, offrampProgramAddress string) chainwriter.MethodConfig {
	return chainwriter.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawExecutionReport": "Report"},
			},
		},
		ChainSpecificName: "execute",
		ArgsTransform:     "CCIP",
		LookupTables: chainwriter.LookupTables{
			DerivedLookupTables: []chainwriter.DerivedLookupTable{
				{
					Name: "PoolLookupTable",
					Accounts: chainwriter.PDALookups{
						Name:      "TokenAdminRegistry",
						PublicKey: offrampAddressConstant(offrampProgramAddress),
						Seeds: []chainwriter.Seed{
							{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
						},
						IsSigner:   false,
						IsWritable: false,
						InternalField: chainwriter.InternalField{
							TypeName: "TokenAdminRegistry",
							Location: "LookupTable",
						},
					},
				},
				getCommonAddressLookupTableConfig(offrampProgramAddress),
			},
		},
		Accounts: []chainwriter.Lookup{
			getOfframpAccountConfig(offrampProgramAddress),
			getReferenceAddressesConfig(offrampProgramAddress),
			chainwriter.PDALookups{
				Name:      "SourceChainState",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: sourceChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "CommitReport",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("commit_report")},
					{Dynamic: chainwriter.AccountLookup{Location: sourceChainSelectorPath}},
					{Dynamic: chainwriter.AccountLookup{
						// The seed is the merkle root of the report, as passed into the input params.
						Location: merkleRoot,
					}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			chainwriter.AccountConstant{
				Name:       "Offramp",
				Address:    offrampProgramAddress,
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "AllowedOfframp",
				PublicKey: getRouterConfig(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("allowed_offramp")},
					{Dynamic: chainwriter.AccountLookup{Location: sourceChainSelectorPath}},
					{Dynamic: offrampAddressConstant(offrampProgramAddress)},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "ExternalExecutionConfig",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_execution_config")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			getSysVarInstructionConstant(),
			chainwriter.PDALookups{
				Name:      "ExternalTokenPoolsSigner",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_token_pools_signer")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountLookup{
				Name:       "UserAccounts",
				Location:   "Info.AbstractReports.Message.ExtraArgsDecoded.Accounts",
				IsWritable: chainwriter.MetaBool{BitmapLocation: "Info.AbstractReports.Message.ExtraArgsDecoded.IsWritableBitmap"},
				IsSigner:   chainwriter.MetaBool{Value: false},
			},
			chainwriter.PDALookups{
				Name: "ReceiverAssociatedTokenAccount",
				PublicKey: chainwriter.AccountConstant{
					Address: solana.SPLAssociatedTokenAccountProgramID.String(),
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte(fromAddress)},
					{Dynamic: chainwriter.AccountLookup{Location: "Info.AbstractReports.Messages.Receiver"}},
					{Dynamic: chainwriter.AccountsFromLookupTable{
						LookupTableName: "PoolLookupTable",
						IncludeIndexes:  []int{6},
					}},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name:      "PerChainTokenConfig",
				PublicKey: offrampAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("ccip_tokenpool_billing")},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "PoolChainConfig",
				PublicKey: chainwriter.AccountsFromLookupTable{
					LookupTableName: "PoolLookupTable",
					IncludeIndexes:  []int{2},
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("ccip_tokenpool_billing")},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountsFromLookupTable{
				LookupTableName: "PoolLookupTable",
				IncludeIndexes:  []int{},
			},
		},
		DebugIDLocation: "Info.AbstractReports.Messages.Header.MessageID",
	}
}

func GetSolanaChainWriterConfig(offrampProgramAddress string, fromAddress string) (chainwriter.ChainWriterConfig, error) {
	// check fromAddress
	pk, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	if pk.IsZero() {
		return chainwriter.ChainWriterConfig{}, errors.New("from address cannot be empty")
	}

	// validate CCIP Offramp IDL, errors not expected
	var offrampIDL solanacodec.IDL
	if err = json.Unmarshal([]byte(ccipOfframpIDL), &offrampIDL); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Offramp IDL, error: %w", err)
	}
	// validate CCIP Router IDL, errors not expected
	var routerIDL solanacodec.IDL
	if err = json.Unmarshal([]byte(ccipOfframpIDL), &routerIDL); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			"ccip-offramp": {
				Methods: map[string]chainwriter.MethodConfig{
					"execute": getExecuteMethodConfig(fromAddress, offrampProgramAddress),
					"commit":  getCommitMethodConfig(fromAddress, offrampProgramAddress),
				},
				IDL: ccipOfframpIDL,
			},
			// Required for the CCIP args transform configured for the execute method which relies on the TokenAdminRegistry stored in the router
			"ccip-router": {
				IDL: ccipRouterIDL,
			},
		},
	}

	return solConfig, nil
}

func getOfframpAccountConfig(offrampProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name: "OfframpAccountConfig",
		PublicKey: chainwriter.AccountConstant{
			Address: offrampProgramAddress,
		},
		Seeds: []chainwriter.Seed{
			{Static: []byte("config")},
		},
		IsSigner:   false,
		IsWritable: false,
	}
}

func offrampAddressConstant(offrampProgramAddress string) chainwriter.AccountConstant {
	return chainwriter.AccountConstant{
		Address: offrampProgramAddress,
	}
}

func getFeeQuoterConfig(offrampProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name:      "FeeQuoter",
		PublicKey: offrampAddressConstant(offrampProgramAddress),
		Seeds: []chainwriter.Seed{
			{Static: []byte("reference_addresses")},
		},
		IsSigner:   false,
		IsWritable: false,
		// Reads the address from the reference addresses account
		InternalField: chainwriter.InternalField{
			TypeName: "ReferenceAddresses",
			Location: "FeeQuoter",
		},
	}
}

func getRouterConfig(offrampProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name:      "Router",
		PublicKey: offrampAddressConstant(offrampProgramAddress),
		Seeds: []chainwriter.Seed{
			{Static: []byte("reference_addresses")},
		},
		IsSigner:   false,
		IsWritable: false,
		// Reads the address from the reference addresses account
		InternalField: chainwriter.InternalField{
			TypeName: "ReferenceAddresses",
			Location: "Router",
		},
	}
}

func getReferenceAddressesConfig(offrampProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name:      "ReferenceAddresses",
		PublicKey: offrampAddressConstant(offrampProgramAddress),
		Seeds: []chainwriter.Seed{
			{Static: []byte("reference_addresses")},
		},
		IsSigner:   false,
		IsWritable: false,
	}
}

func getFeeBillingSignerConfig(offrampProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name:      "FeeBillingSigner",
		PublicKey: offrampAddressConstant(offrampProgramAddress),
		Seeds: []chainwriter.Seed{
			{Static: []byte("fee_billing_signer")},
		},
		IsSigner:   false,
		IsWritable: false,
	}
}

func getCommonAddressLookupTableConfig(offrampProgramAddress string) chainwriter.DerivedLookupTable {
	return chainwriter.DerivedLookupTable{
		Name: "CommonAddressLookupTable",
		Accounts: chainwriter.PDALookups{
			Name:      "OfframpLookupTable",
			PublicKey: offrampAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("reference_addresses")},
			},
			InternalField: chainwriter.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "OfframpLookupTable",
			},
		},
	}
}

func getAuthorityAccountConstant(fromAddress string) chainwriter.AccountConstant {
	return chainwriter.AccountConstant{
		Name:       "Authority",
		Address:    fromAddress,
		IsSigner:   true,
		IsWritable: true,
	}
}

func getSystemProgramConstant() chainwriter.AccountConstant {
	return chainwriter.AccountConstant{
		Name:       "SystemProgram",
		Address:    solana.SystemProgramID.String(),
		IsSigner:   false,
		IsWritable: false,
	}
}

func getSysVarInstructionConstant() chainwriter.AccountConstant {
	return chainwriter.AccountConstant{
		Name:       "SysvarInstructions",
		Address:    solana.SysVarInstructionsPubkey.String(),
		IsSigner:   false,
		IsWritable: false,
	}
}

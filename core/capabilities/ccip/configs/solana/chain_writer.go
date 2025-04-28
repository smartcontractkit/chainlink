package solana

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	ccipconsts "github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	solanacodec "github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
)

var ccipOfframpIDL = idl.FetchCCIPOfframpIDL()
var ccipRouterIDL = idl.FetchCCIPRouterIDL()
var ccipCommonIDL = idl.FetchCommonIDL()

const (
	sourceChainSelectorPath       = "Info.AbstractReports.Messages.Header.SourceChainSelector"
	destTokenAddress              = "Info.AbstractReports.Messages.TokenAmounts.DestTokenAddress"
	tokenReceiverAddress          = "ExtraData.ExtraArgsDecoded.tokenReceiver"
	merkleRootSourceChainSelector = "Info.MerkleRoots.ChainSel"
	merkleRoot                    = "Info.MerkleRoots.MerkleRoot"
)

func getCommitMethodConfig(fromAddress string, offrampProgramAddress string, priceOnly bool) chainwriter.MethodConfig {
	chainSpecificName := "commit"
	if priceOnly {
		chainSpecificName = "commitPriceOnly"
	}
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
		ChainSpecificName: chainSpecificName,
		ArgsTransform:     "CCIPCommit",
		LookupTables: chainwriter.LookupTables{
			DerivedLookupTables: []chainwriter.DerivedLookupTable{
				getCommonAddressLookupTableConfig(offrampProgramAddress),
			},
		},
		Accounts:        buildCommitAccountsList(fromAddress, offrampProgramAddress, priceOnly),
		DebugIDLocation: "",
	}
}

func buildCommitAccountsList(fromAddress, offrampProgramAddress string, priceOnly bool) []chainwriter.Lookup {
	accounts := []chainwriter.Lookup{}
	accounts = append(accounts,
		getOfframpAccountConfig(offrampProgramAddress),
		getReferenceAddressesConfig(offrampProgramAddress),
	)
	if !priceOnly {
		accounts = append(accounts,
			chainwriter.Lookup{
				PDALookups: &chainwriter.PDALookups{
					Name:      "SourceChainState",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []chainwriter.Seed{
						{Static: []byte("source_chain_state")},
						{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: merkleRootSourceChainSelector}}},
					},
					IsSigner:   false,
					IsWritable: true,
				},
			},
			chainwriter.Lookup{
				PDALookups: &chainwriter.PDALookups{
					Name:      "CommitReport",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []chainwriter.Seed{
						{Static: []byte("commit_report")},
						{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: merkleRootSourceChainSelector}}},
						{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: merkleRoot}}},
					},
					IsSigner:   false,
					IsWritable: true,
				},
			},
		)
	}
	accounts = append(accounts,
		getAuthorityAccountConstant(fromAddress),
		getSystemProgramConstant(),
		getSysVarInstructionConstant(),
		getFeeBillingSignerConfig(offrampProgramAddress),
		getFeeQuoterProgramAccount(offrampProgramAddress),
		getFeeQuoterAllowedPriceUpdater(offrampProgramAddress),
		getFeeQuoterConfigLookup(offrampProgramAddress),
		getRMNRemoteProgramAccount(offrampProgramAddress),
		getRMNRemoteCursesLookup(offrampProgramAddress),
		getRMNRemoteConfigLookup(offrampProgramAddress),
		getGlobalStateConfig(offrampProgramAddress),
		getBillingTokenConfig(offrampProgramAddress),
		getChainConfigGasPriceConfig(offrampProgramAddress),
	)
	return accounts
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
		ChainSpecificName:        "execute",
		ArgsTransform:            "CCIPExecute",
		ComputeUnitLimitOverhead: 150_000,
		LookupTables: chainwriter.LookupTables{
			DerivedLookupTables: []chainwriter.DerivedLookupTable{
				{
					Name: "PoolLookupTable",
					Accounts: chainwriter.Lookup{
						PDALookups: &chainwriter.PDALookups{
							Name:      "TokenAdminRegistry",
							PublicKey: getRouterProgramAccount(offrampProgramAddress),
							Seeds: []chainwriter.Seed{
								{Static: []byte("token_admin_registry")},
								{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: destTokenAddress}}},
							},
							IsSigner:   false,
							IsWritable: false,
							InternalField: chainwriter.InternalField{
								TypeName: "TokenAdminRegistry",
								Location: "LookupTable",
								// TokenAdminRegistry is in the common program so need to provide the IDL
								IDL: ccipCommonIDL,
							},
						},
					},
					Optional: true, // Lookup table is optional if DestTokenAddress is not present in report
				},
				getCommonAddressLookupTableConfig(offrampProgramAddress),
			},
		},
		ATAs: []chainwriter.ATALookup{
			{
				Location:      destTokenAddress,
				WalletAddress: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: tokenReceiverAddress}},
				TokenProgram: chainwriter.Lookup{
					AccountsFromLookupTable: &chainwriter.AccountsFromLookupTable{
						LookupTableName: "PoolLookupTable",
						IncludeIndexes:  []int{6},
					},
				},
				MintAddress: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: destTokenAddress}},
				Optional:    true, // ATA lookup is optional if DestTokenAddress is not present in report
			},
		},
		Accounts: []chainwriter.Lookup{
			getOfframpAccountConfig(offrampProgramAddress),
			getReferenceAddressesConfig(offrampProgramAddress),
			{
				PDALookups: &chainwriter.PDALookups{
					Name:      "SourceChainState",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []chainwriter.Seed{
						{Static: []byte("source_chain_state")},
						{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: sourceChainSelectorPath}}},
					},
					IsSigner:   false,
					IsWritable: false,
				},
			},
			{
				PDALookups: &chainwriter.PDALookups{
					Name:      "CommitReport",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []chainwriter.Seed{
						{Static: []byte("commit_report")},
						{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: sourceChainSelectorPath}}},
						{Dynamic: chainwriter.Lookup{
							AccountLookup: &chainwriter.AccountLookup{
								// The seed is the merkle root of the report, as passed into the input params.
								Location: merkleRoot,
							}},
						},
					},
					IsSigner:   false,
					IsWritable: true,
				},
			},
			getAddressConstant(offrampProgramAddress),
			{
				PDALookups: &chainwriter.PDALookups{
					Name:      "AllowedOfframp",
					PublicKey: getRouterProgramAccount(offrampProgramAddress),
					Seeds: []chainwriter.Seed{
						{Static: []byte("allowed_offramp")},
						{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: sourceChainSelectorPath}}},
						{Dynamic: getAddressConstant(offrampProgramAddress)},
					},
					IsSigner:   false,
					IsWritable: false,
				},
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			getSysVarInstructionConstant(),
			getRMNRemoteProgramAccount(offrampProgramAddress),
			getRMNRemoteCursesLookup(offrampProgramAddress),
			getRMNRemoteConfigLookup(offrampProgramAddress),
			// logic receiver and user defined messaging accounts are appended in the CCIPExecute args transform
			// user token account, token billing config, pool chain config, and pool lookup table accounts
			// are appended to the accounts list in the CCIPExecute args transform for each token transfer
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
	if err = json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}
	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			ccipconsts.ContractNameOffRamp: {
				Methods: map[string]chainwriter.MethodConfig{
					ccipconsts.MethodExecute:         getExecuteMethodConfig(fromAddress, offrampProgramAddress),
					ccipconsts.MethodCommit:          getCommitMethodConfig(fromAddress, offrampProgramAddress, false),
					ccipconsts.MethodCommitPriceOnly: getCommitMethodConfig(fromAddress, offrampProgramAddress, true),
				},
				IDL: ccipOfframpIDL,
			},
		},
	}

	return solConfig, nil
}

func getOfframpAccountConfig(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name: "OfframpAccountConfig",
			PublicKey: chainwriter.Lookup{
				AccountConstant: &chainwriter.AccountConstant{
					Address: offrampProgramAddress,
				},
			},
			Seeds: []chainwriter.Seed{
				{Static: []byte("config")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getAddressConstant(address string) chainwriter.Lookup {
	return chainwriter.Lookup{
		AccountConstant: &chainwriter.AccountConstant{
			Address:    address,
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeQuoterProgramAccount(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      ccipconsts.ContractNameFeeQuoter,
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
			// Reads the address from the reference addresses account
			InternalField: chainwriter.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "FeeQuoter",
				IDL:      ccipOfframpIDL,
			},
		},
	}
}

func getRouterProgramAccount(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      ccipconsts.ContractNameRouter,
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
			// Reads the address from the reference addresses account
			InternalField: chainwriter.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "Router",
				IDL:      ccipOfframpIDL,
			},
		},
	}
}

func getReferenceAddressesConfig(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "ReferenceAddresses",
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeBillingSignerConfig(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "FeeBillingSigner",
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("fee_billing_signer")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeQuoterAllowedPriceUpdater(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name: "FeeQuoterAllowedPriceUpdater",
			// Fetch fee quoter public key to use as program ID for PDA
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("allowed_price_updater")},
				{Dynamic: getFeeBillingSignerConfig(offrampProgramAddress)},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeQuoterConfigLookup(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name: "FeeQuoterConfig",
			// Fetch fee quoter public key to use as program ID for PDA
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("config")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getRMNRemoteProgramAccount(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      ccipconsts.ContractNameRMNRemote,
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
			// Reads the address from the reference addresses account
			InternalField: chainwriter.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "RmnRemote",
				IDL:      ccipOfframpIDL,
			},
		},
	}
}

func getRMNRemoteCursesLookup(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "RMNRemoteCurses",
			PublicKey: getRMNRemoteProgramAccount(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("curses")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getRMNRemoteConfigLookup(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "RMNRemoteConfig",
			PublicKey: getRMNRemoteProgramAccount(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("config")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getGlobalStateConfig(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "GlobalState",
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("state")},
			},
			IsSigner:   false,
			IsWritable: true,
		},
		Optional: true,
	}
}

func getBillingTokenConfig(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "BillingTokenConfig",
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("fee_billing_token_config")},
				{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: "Info.TokenPriceUpdates.TokenID"}}},
			},
			IsSigner:   false,
			IsWritable: true,
		},
		Optional: true,
	}
}

func getChainConfigGasPriceConfig(offrampProgramAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		PDALookups: &chainwriter.PDALookups{
			Name:      "ChainConfigGasPrice",
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []chainwriter.Seed{
				{Static: []byte("dest_chain")},
				{Dynamic: chainwriter.Lookup{AccountLookup: &chainwriter.AccountLookup{Location: "Info.GasPriceUpdates.ChainSel"}}},
			},
			IsSigner:   false,
			IsWritable: true,
		},
		Optional: true,
	}
}

// getCommonAddressLookupTableConfig returns the lookup table config that fetches the lookup table address from a PDA on-chain
// The offramp contract contains a PDA with a ReferenceAddresses struct that stores the lookup table address in the OfframpLookupTable field
func getCommonAddressLookupTableConfig(offrampProgramAddress string) chainwriter.DerivedLookupTable {
	return chainwriter.DerivedLookupTable{
		Name: "CommonAddressLookupTable",
		Accounts: chainwriter.Lookup{
			PDALookups: &chainwriter.PDALookups{
				Name:      "OfframpLookupTable",
				PublicKey: getAddressConstant(offrampProgramAddress),
				Seeds: []chainwriter.Seed{
					{Static: []byte("reference_addresses")},
				},
				InternalField: chainwriter.InternalField{
					TypeName: "ReferenceAddresses",
					Location: "OfframpLookupTable",
					IDL:      ccipOfframpIDL,
				},
			},
		},
	}
}

func getAuthorityAccountConstant(fromAddress string) chainwriter.Lookup {
	return chainwriter.Lookup{
		AccountConstant: &chainwriter.AccountConstant{
			Name:       "Authority",
			Address:    fromAddress,
			IsSigner:   true,
			IsWritable: true,
		},
	}
}

func getSystemProgramConstant() chainwriter.Lookup {
	return chainwriter.Lookup{
		AccountConstant: &chainwriter.AccountConstant{
			Name:       "SystemProgram",
			Address:    solana.SystemProgramID.String(),
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getSysVarInstructionConstant() chainwriter.Lookup {
	return chainwriter.Lookup{
		AccountConstant: &chainwriter.AccountConstant{
			Name:       "SysvarInstructions",
			Address:    solana.SysVarInstructionsPubkey.String(),
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

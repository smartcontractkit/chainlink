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

var ccipRouterIDL = idl.FetchCCIPRouterIDL()

const (
	destChainSelectorPath   = "Info.AbstractReports.Messages.Header.DestChainSelector"
	destTokenAddress        = "Info.AbstractReports.Messages.TokenAmounts.DestTokenAddress"
	merkleRootChainSelector = "Info.MerkleRoots.ChainSel"
)

func getCommitMethodConfig(fromAddress string, routerProgramAddress string, commonAddressesLookupTable solana.PublicKey) chainwriter.MethodConfig {
	sysvarInstructionsAddress := solana.SysVarInstructionsPubkey.String()
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
			StaticLookupTables: []solana.PublicKey{
				commonAddressesLookupTable,
			},
		},
		Accounts: []chainwriter.Lookup{
			getRouterAccountConfig(routerProgramAddress),
			chainwriter.PDALookups{
				Name: "SourceChainState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			chainwriter.PDALookups{
				Name: "RouterReportAccount",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("commit_report")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
					{Dynamic: chainwriter.AccountLookup{
						Location: "Info.MerkleRoots.MerkleRoot",
					}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			chainwriter.AccountConstant{
				Name:       "SysvarInstructions",
				Address:    sysvarInstructionsAddress,
				IsSigner:   true,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "GlobalState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("state")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "BillingTokenConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("fee_billing_token_config")},
					{Dynamic: chainwriter.AccountLookup{Location: "Info.TokenPrices.TokenID"}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "ChainConfigGasPrice",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
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

func getExecuteMethodConfig(fromAddress string, routerProgramAddress string, commonAddressesLookupTable solana.PublicKey) chainwriter.MethodConfig {
	sysvarInstructionsAddress := solana.SysVarInstructionsPubkey.String()
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
						Name: "TokenAdminRegistry",
						PublicKey: chainwriter.AccountConstant{
							Address: routerProgramAddress,
						},
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
			},
			StaticLookupTables: []solana.PublicKey{
				commonAddressesLookupTable,
			},
		},
		Accounts: []chainwriter.Lookup{
			getRouterAccountConfig(routerProgramAddress),
			chainwriter.PDALookups{
				Name: "SourceChainState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "CommitReport",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_execution_config")},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
					{Dynamic: chainwriter.AccountLookup{
						// The seed is the merkle root of the report, as passed into the input params.
						Location: "Info.MerkleRoots.MerkleRoot",
					}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			chainwriter.PDALookups{
				Name: "ExternalExecutionConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_execution_config")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			chainwriter.AccountConstant{
				Name:       "SysvarInstructions",
				Address:    sysvarInstructionsAddress,
				IsSigner:   true,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "ExternalTokenPoolsSigner",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
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
				Name: "PerChainTokenConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
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
		DebugIDLocation: "AbstractReport.Message.MessageID",
	}
}

func GetSolanaChainWriterConfig(routerProgramAddress string, commonAddressesLookupTable solana.PublicKey, fromAddress string) (chainwriter.ChainWriterConfig, error) {
	// check fromAddress
	pk, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	if pk.IsZero() {
		return chainwriter.ChainWriterConfig{}, errors.New("from address cannot be empty")
	}

	// validate CCIP Router IDL, errors not expected
	var idl solanacodec.IDL
	if err = json.Unmarshal([]byte(ccipRouterIDL), &idl); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	// solConfig references the ccip_example_config.go from github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter, which is currently subject to change
	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			"ccip-router": {
				Methods: map[string]chainwriter.MethodConfig{
					"execute": getExecuteMethodConfig(fromAddress, routerProgramAddress, commonAddressesLookupTable),
					"commit":  getCommitMethodConfig(fromAddress, routerProgramAddress, commonAddressesLookupTable),
				},
				IDL: ccipRouterIDL},
		},
	}

	return solConfig, nil
}

func getRouterAccountConfig(routerProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name: "RouterAccountConfig",
		PublicKey: chainwriter.AccountConstant{
			Address: routerProgramAddress,
		},
		Seeds: []chainwriter.Seed{
			{Static: []byte("config")},
		},
		IsSigner:   false,
		IsWritable: false,
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

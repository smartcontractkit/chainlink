package solana

import (
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
)

var ccipRouterIDL = idl.FetchCCIPRouterIDL()

const (
	destChainSelectorPath = "Message.Header.DestChainSelector"
	destTokenAddress      = "Message.TokenAmounts.DestTokenAddress"
)

func getCommitMethodConfig(fromAddress string, routerProgramAddress string, sysvarInstructionsAddress string, computeBudgetProgramAddress string, commonAddressesLookupTable solana.PublicKey, routerAccountConfig chainwriter.PDALookups) chainwriter.MethodConfig {
	return chainwriter.MethodConfig{
		FromAddress:        fromAddress,
		InputModifications: nil,
		ChainSpecificName:  "commit",
		LookupTables: chainwriter.LookupTables{
			StaticLookupTables: []solana.PublicKey{
				commonAddressesLookupTable,
			},
		},
		Accounts: []chainwriter.Lookup{
			routerAccountConfig,
			chainwriter.PDALookups{
				Name: "SourceChainState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: "MerkleRoot.DestChainSelector"}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountConstant{
				Name:       "RouterProgram",
				Address:    routerProgramAddress,
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "RouterAccountConfig",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("config")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "RouterAccountState",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("state")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "RouterReportAccount",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Dynamic: chainwriter.AccountLookup{
						Location: "args.MerkleRoots",
					}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountConstant{
				Name:       "ComputeBudgetProgram",
				Address:    computeBudgetProgramAddress,
				IsSigner:   true,
				IsWritable: false,
			},
			chainwriter.AccountConstant{
				Name:       "SysvarInstructions",
				Address:    sysvarInstructionsAddress,
				IsSigner:   true,
				IsWritable: false,
			},
		},
		DebugIDLocation: "",
	}
}

func getExecuteProgramConfig(fromAddress string, routerProgramAddress string, sysvarInstructionsAddress string, computeBudgetProgramAddress string, commonAddressesLookupTable solana.PublicKey, routerAccountConfig chainwriter.PDALookups) chainwriter.MethodConfig {
	return chainwriter.MethodConfig{
		FromAddress:        fromAddress,
		InputModifications: nil,
		ChainSpecificName:  "execute",
		LookupTables: chainwriter.LookupTables{
			DerivedLookupTables: []chainwriter.DerivedLookupTable{
				{
					Name: "RegistryTokenState",
					Accounts: chainwriter.PDALookups{
						Name: "RegistryTokenState",
						PublicKey: chainwriter.AccountConstant{
							Address:    routerProgramAddress,
							IsSigner:   false,
							IsWritable: false,
						},
						Seeds: []chainwriter.Seed{
							{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
						},
						IsSigner:   false,
						IsWritable: false,
					},
				},
			},
			StaticLookupTables: []solana.PublicKey{
				commonAddressesLookupTable,
			},
		},
		Accounts: []chainwriter.Lookup{
			routerAccountConfig,
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
					{Static: []byte("commit_report")},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
					{Dynamic: chainwriter.AccountLookup{
						Location: "Info.MerkleRoot",
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
			chainwriter.AccountConstant{
				Name:       "Authority",
				Address:    fromAddress,
				IsSigner:   true,
				IsWritable: true,
			},
			chainwriter.AccountConstant{
				Name:       "SystemProgram",
				Address:    solana.SystemProgramID.String(),
				IsSigner:   false,
				IsWritable: false,
			},
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
				Name:     "UserAccounts",
				Location: "Message.ExtraArgs.Accounts",
			},
			chainwriter.PDALookups{
				Name: "ReceiverAssociatedTokenAccount",
				PublicKey: chainwriter.AccountConstant{
					Address: solana.SPLAssociatedTokenAccountProgramID.String(),
				},
				Seeds: []chainwriter.Seed{
					{Dynamic: chainwriter.AccountLookup{Location: "Message.Receiver"}},
					{Dynamic: chainwriter.AccountsFromLookupTable{
						LookupTableName: "RegistryTokenState",
						IncludeIndexes:  []int{5},
					}},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
				},
			},
			chainwriter.PDALookups{
				Name: "PerChainTokenConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountsFromLookupTable{
				LookupTableName: "RegistryTokenState",
				IncludeIndexes:  []int{},
			},
			chainwriter.PDALookups{
				Name: "RegistryTokenConfig",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "UserNoncePerChain",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Dynamic: chainwriter.AccountLookup{Location: "Message.Receiver"}},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
			},
			chainwriter.PDALookups{
				Name: "CPISigner",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_token_pools_signer")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountConstant{
				Name:       "ComputeBudgetProgram",
				Address:    computeBudgetProgramAddress,
				IsSigner:   true,
				IsWritable: false,
			},
		},
		DebugIDLocation: "Message.MessageID",
	}
}

func GetSolanaChainWriterConfig(routerProgramAddress string, commonAddressesLookupTable solana.PublicKey, fromAddress string) (chainwriter.ChainWriterConfig, error) {
	computeBudgetProgramAddress := solana.ComputeBudget.String()
	sysvarInstructionsAddress := solana.SysVarInstructionsPubkey.String()

	// check fromAddress
	_, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	// validate CCIP Router IDL, errors not expected
	var idl codec.IDL
	if err = json.Unmarshal([]byte(ccipRouterIDL), &idl); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	routeAccountConfig := chainwriter.PDALookups{
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

	// solConfig references the ccip_example_config.go from github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter, which is currently subject to change
	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			"ccip-router": {
				Methods: map[string]chainwriter.MethodConfig{
					"execute": getExecuteProgramConfig(fromAddress, routerProgramAddress, sysvarInstructionsAddress, computeBudgetProgramAddress, commonAddressesLookupTable, routeAccountConfig),
					"commit":  getCommitMethodConfig(fromAddress, routerProgramAddress, sysvarInstructionsAddress, computeBudgetProgramAddress, commonAddressesLookupTable, routeAccountConfig),
				},
				IDL: ccipRouterIDL},
		},
	}

	return solConfig, nil
}

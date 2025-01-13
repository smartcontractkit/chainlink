package solana

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
)

//go:embed ccip_router.json
var ccipRouter string

func GetSolanaChainWriterConfig(fromAddress string) (chainwriter.ChainWriterConfig, error) {
	// TODO once on-chain account lookup address are available, the routerProgramAddress and commonAddressesLookupTable should be updated
	routerProgramAddress := ""
	var commonAddressesLookupTable []byte

	computeBudgetProgramAddress := solana.ComputeBudget.String()
	sysvarInstructionsAddress := solana.SysVarInstructionsPubkey.String()

	// check fromAddress
	_, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	// validate CCIP Router IDL, errors not expected
	var idl codec.IDL
	if err = json.Unmarshal([]byte(ccipRouter), &idl); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	// solConfig references the ccip_example_config.go from github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter, which is currently subject to change
	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			"ccip-router": {
				Methods: map[string]chainwriter.MethodConfig{
					"execute": {
						FromAddress:        fromAddress,
						InputModifications: nil,
						ChainSpecificName:  "execute",
						LookupTables: chainwriter.LookupTables{
							// DerivedLookupTables are useful in both the ways described above.
							// 	a. The user can configure any type of look up to get a list of lookupTables to read from.
							// 	b. The ChainWriter reads from this lookup table and store the internal addresses in memory
							//	c. Later, in the []Accounts the user can specify which accounts to include in the TX with an chainwriter.AccountsFromLookupTable lookup.
							// 	d. Lastly, the lookup table is used to compress the size of the transaction.
							DerivedLookupTables: []chainwriter.DerivedLookupTable{
								{
									Name: "RegistryTokenState",
									// In this case, the user configured the lookup table accounts to use a PDALookup, which
									// generates a list of one of more PDA accounts based on the input parameters. Specifically,
									// there will be multiple PDA accounts if there are multiple addresses in the message, otherwise,
									// there will only be one PDA account to read from. The PDA account corresponds to the lookup table.
									Accounts: chainwriter.PDALookups{
										Name: "RegistryTokenState",
										PublicKey: chainwriter.AccountConstant{
											Address:    routerProgramAddress,
											IsSigner:   false,
											IsWritable: false,
										},
										// Seeds would be used if the user needed to look up addresses to use as seeds, which isn't the case here.
										Seeds: []chainwriter.Seed{
											{Dynamic: chainwriter.AccountLookup{Location: "Message.TokenAmounts.DestTokenAddress"}},
										},
										IsSigner:   false,
										IsWritable: false,
									},
								},
							},
							// Static lookup tables are the traditional use case (point 2 above) of Lookup tables. These are lookup
							// tables which contain commonly used addresses in all CCIP execute transactions. The ChainWriter reads
							// these lookup tables and appends them to the transaction to reduce the size of the transaction.
							StaticLookupTables: []solana.PublicKey{
								solana.PublicKey(commonAddressesLookupTable),
							},
						},
						// The Accounts field is where the user specifies which accounts to include in the transaction. Each Lookup
						// resolves to one or more on-chain addresses.
						Accounts: []chainwriter.Lookup{
							// The accounts can be of any of the following types:
							// 1. Account constant
							// 2. Account Lookup - Based on data from input parameters
							// 3. Lookup Table content - Get all the accounts from a lookup table
							// 4. PDA Account Lookup - Based on another account and a seed/s
							//	Nested PDA Account with seeds from:
							// 		-> input parameters
							// 		-> constant
							// 	PDALookups can resolve to multiple addresses if:
							// 		A) The PublicKey lookup resolves to multiple addresses (i.e. multiple token addresses)
							// 		B) The Seeds or ValueSeeds resolve to multiple values
							// PDA lokoup with constant seed
							chainwriter.PDALookups{
								Name: "RouterAccountConfig",
								PublicKey: chainwriter.AccountConstant{
									Address: routerProgramAddress,
								},
								Seeds: []chainwriter.Seed{
									{Static: []byte("config")},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							chainwriter.PDALookups{
								Name: "SourceChainState",
								// PublicKey is a constant account in this case, not a lookup.
								PublicKey: chainwriter.AccountConstant{
									Address: routerProgramAddress,
								},
								// Similar to the RegistryTokenState above, the user is looking up PDA accounts based on the dest tokens.
								Seeds: []chainwriter.Seed{
									{Static: []byte("source_chain_state")},
									{Dynamic: chainwriter.AccountLookup{Location: "Message.Header.DestChainSelector"}},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							// PDA lookup to get the Router Report Accounts.
							chainwriter.PDALookups{
								Name: "CommitReport",
								// The public key is a constant Router address.
								PublicKey: chainwriter.AccountConstant{
									Address: routerProgramAddress,
								},
								Seeds: []chainwriter.Seed{
									{Static: []byte("commit_report")},
									{Dynamic: chainwriter.AccountLookup{Location: "Message.Header.DestChainSelector"}},
									{Dynamic: chainwriter.AccountLookup{
										// The seed is the merkle root of the report, as passed into the input params.
										Location: "Info.MerkleRoot",
									}},
								},
								IsSigner:   false,
								IsWritable: true,
							},
							// Static PDA lookup
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
							// feePayer/authority address
							chainwriter.AccountConstant{
								Name:       "Authority",
								Address:    fromAddress,
								IsSigner:   true,
								IsWritable: true,
							},
							// Account constant
							chainwriter.AccountConstant{
								Name:       "SystemProgram",
								Address:    solana.SystemProgramID.String(),
								IsSigner:   false,
								IsWritable: false,
							},
							// Account constant
							chainwriter.AccountConstant{
								Name:       "SysvarInstructions",
								Address:    sysvarInstructionsAddress,
								IsSigner:   true,
								IsWritable: false,
							},
							// Static PDA lookup
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
							// User specified accounts - formatted as AccountMeta
							chainwriter.AccountLookup{
								Name:     "UserAccounts",
								Location: "Message.ExtraArgs.Accounts",
							},
							// PDA Account Lookup - Based on an account lookup and an address lookup
							chainwriter.PDALookups{
								Name: "ReceiverAssociatedTokenAccount",
								PublicKey: chainwriter.AccountConstant{
									Address: solana.SPLAssociatedTokenAccountProgramID.String(),
								},
								Seeds: []chainwriter.Seed{
									// receiver address
									{Dynamic: chainwriter.AccountLookup{Location: "Message.Receiver"}},
									// token program
									{Dynamic: chainwriter.AccountsFromLookupTable{
										LookupTableName: "RegistryTokenState",
										IncludeIndexes:  []int{5},
									}},
									// mint
									{Dynamic: chainwriter.AccountLookup{Location: "Message.TokenAmounts.DestTokenAddress"}},
								},
							},
							chainwriter.PDALookups{
								Name: "PerChainTokenConfig",
								// PublicKey is a constant account in this case, not a lookup.
								PublicKey: chainwriter.AccountConstant{
									Address: routerProgramAddress,
								},
								// Similar to the RegistryTokenState above, the user is looking up PDA accounts based on the dest tokens.
								Seeds: []chainwriter.Seed{
									{Dynamic: chainwriter.AccountLookup{Location: "Message.TokenAmounts.DestTokenAddress"}},
									{Dynamic: chainwriter.AccountLookup{Location: "Message.Header.DestChainSelector"}},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							// Lookup Table content - Get the accounts from the derived lookup table above
							chainwriter.AccountsFromLookupTable{
								LookupTableName: "RegistryTokenState",
								IncludeIndexes:  []int{}, // If left empty, all addresses will be included. Otherwise, only the specified indexes will be included.
							},
							// PDA Lookup for the RegistryTokenConfig.
							chainwriter.PDALookups{
								Name: "RegistryTokenConfig",
								// constant public key
								PublicKey: chainwriter.AccountConstant{
									Address:    routerProgramAddress,
									IsSigner:   false,
									IsWritable: false,
								},
								// The seed, once again, is the destination token address.
								Seeds: []chainwriter.Seed{
									{Dynamic: chainwriter.AccountLookup{Location: "Message.TokenAmounts.DestTokenAddress"}},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							// PDA lookup to get UserNoncePerChain
							chainwriter.PDALookups{
								Name: "UserNoncePerChain",
								// The public key is a constant Router address.
								PublicKey: chainwriter.AccountConstant{
									Address:    routerProgramAddress,
									IsSigner:   false,
									IsWritable: false,
								},
								// In this case, the user configured multiple seeds. These will be used in conjunction
								// with the public key to generate one or multiple PDA accounts.
								Seeds: []chainwriter.Seed{
									{Dynamic: chainwriter.AccountLookup{Location: "Message.Receiver"}},
									{Dynamic: chainwriter.AccountLookup{Location: "Message.Header.DestChainSelector"}},
								},
							},
							// PDA lokoup with constant seed
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
							// Account constant
							chainwriter.AccountConstant{
								Name:       "ComputeBudgetProgram",
								Address:    computeBudgetProgramAddress,
								IsSigner:   true,
								IsWritable: false,
							},
						},
						// TBD where this will be in the report
						// This will be appended to every error message
						DebugIDLocation: "Message.MessageID",
					},
					"commit": {
						FromAddress:        fromAddress,
						InputModifications: nil,
						ChainSpecificName:  "commit",
						LookupTables: chainwriter.LookupTables{
							StaticLookupTables: []solana.PublicKey{
								solana.PublicKey(commonAddressesLookupTable),
							},
						},
						Accounts: []chainwriter.Lookup{
							// Static PDA lookup
							chainwriter.PDALookups{
								Name: "RouterAccountConfig",
								PublicKey: chainwriter.AccountConstant{
									Address: routerProgramAddress,
								},
								Seeds: []chainwriter.Seed{
									{Static: []byte("config")},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							chainwriter.PDALookups{
								Name: "SourceChainState",
								// PublicKey is a constant account in this case, not a lookup.
								PublicKey: chainwriter.AccountConstant{
									Address: routerProgramAddress,
								},
								// Similar to the RegistryTokenState above, the user is looking up PDA accounts based on the dest tokens.
								Seeds: []chainwriter.Seed{
									{Static: []byte("source_chain_state")},
									{Dynamic: chainwriter.AccountLookup{Location: "MerkleRoot.DestChainSelector"}},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							// Account constant
							chainwriter.AccountConstant{
								Name:       "RouterProgram",
								Address:    routerProgramAddress,
								IsSigner:   false,
								IsWritable: false,
							},
							// PDA lokoup with constant seed
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
							// PDA lokoup with constant seed
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
							// PDA lookup to get the Router Report Accounts.
							chainwriter.PDALookups{
								Name: "RouterReportAccount",
								// The public key is a constant Router address.
								PublicKey: chainwriter.AccountConstant{
									Address:    routerProgramAddress,
									IsSigner:   false,
									IsWritable: false,
								},
								Seeds: []chainwriter.Seed{
									{Dynamic: chainwriter.AccountLookup{
										// The seed is the merkle root of the report, as passed into the input params.
										Location: "args.MerkleRoots",
									}},
								},
								IsSigner:   false,
								IsWritable: false,
							},
							// Account constant
							chainwriter.AccountConstant{
								Name:       "ComputeBudgetProgram",
								Address:    computeBudgetProgramAddress,
								IsSigner:   true,
								IsWritable: false,
							},
							// Account constant
							chainwriter.AccountConstant{
								Name:       "SysvarInstructions",
								Address:    sysvarInstructionsAddress,
								IsSigner:   true,
								IsWritable: false,
							},
						},
						DebugIDLocation: "",
					},
				},
				IDL: ccipRouter},
		},
	}

	return solConfig, nil
}

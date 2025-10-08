package environment

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/gagliardetto/solana-go"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
)

type PrepareFundCLNodesOpDeps struct {
	TestLogger  zerolog.Logger
	Env         *cldf.Environment
	Blockchains []*cre.Blockchain
	DonTopology *cre.DonTopology
}

type PrepareFundCLNodesOpInput struct {
	FundingPerChainFamilyForEachNode map[string]uint64
}

type PrepareFundCLNodesOpOutput struct {
	PrivateKeysPerChainFamily        map[string]map[uint64][]byte
	FundingPerChainFamilyForEachNode map[string]uint64
}

// PrepareCLNodesFundingOp prepares funding accounts for Chainlink nodes by generating new key pairs and funding them from the root account.
// It cannot be run in parallel with other operations that modify blockchain state, as it relies on the root account's balance.
var PrepareCLNodesFundingOp = operations.NewOperation[PrepareFundCLNodesOpInput, *PrepareFundCLNodesOpOutput, PrepareFundCLNodesOpDeps](
	"prepare-cl-nodes-funding-op",
	semver.MustParse("1.0.0"),
	"Prepare inputs for funding of Chainlink Nodes",
	func(b operations.Bundle, deps PrepareFundCLNodesOpDeps, input PrepareFundCLNodesOpInput) (*PrepareFundCLNodesOpOutput, error) {
		ctx := b.GetContext()

		output := &PrepareFundCLNodesOpOutput{
			PrivateKeysPerChainFamily:        make(map[string]map[uint64][]byte),
			FundingPerChainFamilyForEachNode: input.FundingPerChainFamilyForEachNode,
		}
		requiredFundingPerChain := make(map[uint64]uint64)
		for _, don := range deps.DonTopology.Dons.List() {
			for _, bc := range deps.Blockchains {
				if !flags.RequiresForwarderContract(don.Flags, bc.ChainID) && bc.SolChain == nil {
					continue
				}

				if bc.SolChain != nil {
					requiredFundingPerChain[bc.SolChain.ChainSelector] += input.FundingPerChainFamilyForEachNode[chainselectors.FamilySolana] * uint64(len(don.Nodes))
					continue
				}

				if bc.CtfOutput.Family == blockchain.FamilyTron {
					requiredFundingPerChain[bc.ChainSelector] += input.FundingPerChainFamilyForEachNode[chainselectors.FamilyTron] * uint64(len(don.Nodes))
					continue
				}

				requiredFundingPerChain[bc.ChainSelector] += input.FundingPerChainFamilyForEachNode[chainselectors.FamilyEVM] * uint64(len(don.Nodes))
			}
		}

		for _, bc := range deps.Blockchains {
			if requiredFundingPerChain[bc.ChainSelector] == 0 {
				continue
			}

			chainFamily := chainselectors.FamilyEVM
			if bc.SolChain != nil {
				chainFamily = chainselectors.FamilySolana
			} else if bc.CtfOutput.Family == blockchain.FamilyTron {
				chainFamily = chainselectors.FamilyTron
			}

			switch chainFamily {
			case chainselectors.FamilyEVM:
				if _, exists := output.PrivateKeysPerChainFamily[chainFamily]; !exists {
					output.PrivateKeysPerChainFamily[chainFamily] = make(map[uint64][]byte)
				}

				if _, exists := output.PrivateKeysPerChainFamily[chainFamily][bc.ChainSelector]; !exists {
					publicAddress, privateKeyBytes, pkErr := generatePubPrivKeyPairForEth()
					if pkErr != nil {
						return nil, pkgerrors.Wrap(pkErr, "failed to generate pub/priv key pair for EVM funding account")
					}

					fundingAmount := libc.MustSafeInt64(requiredFundingPerChain[bc.ChainSelector])
					fundingAmount += (fundingAmount / 5) // add 20% to cover gas fees

					_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bc.SethClient, libfunding.FundsToSend{
						ToAddress:  *publicAddress,
						Amount:     big.NewInt(fundingAmount),
						PrivateKey: bc.SethClient.MustGetRootPrivateKey(),
					})

					if fundingErr != nil {
						return nil, pkgerrors.Wrapf(fundingErr, "failed to fund funding account %s on chain %d", publicAddress.String(), bc.ChainID)
					}

					output.PrivateKeysPerChainFamily[chainFamily][bc.ChainSelector] = privateKeyBytes
				}
			case chainselectors.FamilySolana:
				if _, exists := output.PrivateKeysPerChainFamily[chainFamily]; !exists {
					output.PrivateKeysPerChainFamily[chainFamily] = make(map[uint64][]byte)
				}
				if _, exists := output.PrivateKeysPerChainFamily[chainFamily][bc.SolChain.ChainSelector]; !exists {
					private, pkErr := solana.NewRandomPrivateKey()
					if pkErr != nil {
						return nil, pkgerrors.Wrap(pkErr, "failed to generate private key for solana")
					}
					public := private.PublicKey()
					fundingErr := libfunding.SendFundsSol(ctx, zerolog.Logger{}, bc.SolClient, libfunding.FundsToSendSol{
						Recipent:   public,
						PrivateKey: bc.SolChain.PrivateKey,
						Amount:     requiredFundingPerChain[bc.SolChain.ChainSelector],
					})
					if fundingErr != nil {
						return nil, pkgerrors.Wrapf(fundingErr, " failed to fund funding accounts on chain %v", bc.SolChain.ChainID)
					}

					output.PrivateKeysPerChainFamily[chainFamily][bc.SolChain.ChainSelector] = private
				}
			case chainselectors.FamilyTron:
				// TRON uses its own built-in funding account, no preparation needed
				continue
			default:
				return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
			}
		}

		return output, nil
	},
)

func generatePubPrivKeyPairForEth() (*common.Address, []byte, error) {
	privateKey, pkErr := crypto.GenerateKey()
	if pkErr != nil {
		return nil, nil, pkgerrors.Wrap(pkErr, "failed to generate private key for funding accounts")
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("error casting public key to ECDSA")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &publicAddress, privateKeyBytes, nil
}

type FundCLNodesOpDeps struct {
	TestLogger        zerolog.Logger
	Env               *cldf.Environment
	BlockchainOutputs []*cre.Blockchain
	DonTopology       *cre.DonTopology
}

type FundCLNodesOpInput struct {
	PrivateKeyPerChainFamily    map[string]map[uint64][]byte
	FundingAmountPerChainFamily map[string]uint64
}

type FundCLNodesOpOutput struct {
}

var FundCLNodesOp = operations.NewOperation(
	"fund-cl-nodes-op",
	semver.MustParse("1.0.0"),
	"Fund Chainlink Nodes",
	func(b operations.Bundle, deps FundCLNodesOpDeps, input FundCLNodesOpInput) (*FundCLNodesOpOutput, error) {
		ctx := b.GetContext()
		for donIndex, don := range deps.DonTopology.Dons.List() {
			deps.TestLogger.Info().Msgf("Funding nodes for DON %s", don.Name)
			for _, bc := range deps.BlockchainOutputs {
				if !flags.RequiresForwarderContract(don.Flags, bc.ChainID) &&
					bc.SolChain == nil { // for now, we can only write to solana, so we consider forwarder is always present
					continue
				}
				for _, node := range deps.DonTopology.Dons.List()[donIndex].Nodes {
					chainFamily := chainselectors.FamilyEVM
					if bc.SolChain != nil {
						chainFamily = chainselectors.FamilySolana
					} else if bc.CtfOutput.Family == blockchain.FamilyTron {
						chainFamily = chainselectors.FamilyTron
					}

					fundingAmount, ok := input.FundingAmountPerChainFamily[chainFamily]
					if !ok {
						return nil, fmt.Errorf("missing funding amount for chain family %s", chainFamily)
					}

					switch chainFamily {
					case chainselectors.FamilyEVM:
						if err := fundEthAddress(ctx, deps.TestLogger, node, fundingAmount, bc, input.PrivateKeyPerChainFamily); err != nil {
							return nil, err
						}
					case chainselectors.FamilySolana:
						if err := fundSolanaAddress(ctx, deps.TestLogger, node, fundingAmount, bc, input.PrivateKeyPerChainFamily); err != nil {
							return nil, err
						}
					case chainselectors.FamilyTron:
						nodeAddress := getTronNodeAddress(node, bc)
						if nodeAddress == nil {
							deps.TestLogger.Info().Msgf("No EVM key for chainID %d (Tron) found for node %s. Skipping funding", bc.ChainID, node.Name)
							continue // Skip nodes without EVM keys for this chain
						}
						if err := FundTronAddress(ctx, deps.TestLogger, *nodeAddress, fundingAmount, bc, deps.Env); err != nil {
							return nil, err
						}
					default:
						return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
					}
				}
			}

			deps.TestLogger.Info().Msgf("Funded nodes for DON %s", don.Name)
		}

		return &FundCLNodesOpOutput{}, nil
	},
)

func fundEthAddress(ctx context.Context, testLogger zerolog.Logger, node *cre.Node, fundingAmount uint64, bc *cre.Blockchain, privateKeyPerChainFamily map[string]map[uint64][]byte) error {
	evmKey, ok := node.Keys.EVM[bc.ChainID]
	if !ok {
		testLogger.Info().Msgf("No EVM key for chainID %d found for node %s. Skipping funding", bc.ChainID, node.Name)
		return nil // Skip nodes without EVM keys for this chain
	}

	nodeAddress := evmKey.PublicAddress.String()
	testLogger.Info().Msgf("Attempting to fund EVM account %s", nodeAddress)

	fundingPrivateKey, ok := privateKeyPerChainFamily["evm"][bc.ChainSelector]
	if !ok {
		return fmt.Errorf("missing funding private key for chain familyevm, chain %d", bc.ChainID)
	}

	fundingKey, fkErr := crypto.ToECDSA(fundingPrivateKey)
	if fkErr != nil {
		return pkgerrors.Wrap(fkErr, "failed to convert funding private key to ECDSA")
	}

	_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bc.SethClient, libfunding.FundsToSend{
		ToAddress:  common.HexToAddress(nodeAddress),
		Amount:     big.NewInt(libc.MustSafeInt64(fundingAmount)),
		PrivateKey: fundingKey,
	})

	if fundingErr != nil {
		return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", nodeAddress)
	}
	testLogger.Info().Msgf("Successfully funded EVM account %s", nodeAddress)

	return nil
}

func fundSolanaAddress(ctx context.Context, testLogger zerolog.Logger, node *cre.Node, fundingAmount uint64, bc *cre.Blockchain, _ map[string]map[uint64][]byte) error {
	funder := bc.SolChain.PrivateKey
	solKey, ok := node.Keys.Solana[bc.SolChain.ChainID]
	if !ok {
		return fmt.Errorf("missing solana key for node %s on chain %s", node.Name, bc.SolChain.ChainID)
	}
	recipient := solana.MustPublicKeyFromBase58(solKey.PublicAddress.String())
	testLogger.Info().Msgf("Attempting to fund Solana account %s", recipient.String())

	err := libfunding.SendFundsSol(ctx, zerolog.Logger{}, bc.SolClient, libfunding.FundsToSendSol{
		Recipent:   recipient,
		PrivateKey: funder,
		Amount:     fundingAmount,
	})
	if err != nil {
		return fmt.Errorf("failed to fund Solana account for a node: %w", err)
	}
	testLogger.Info().Msgf("Successfully funded Solana account %s", recipient.String())

	return nil
}

func getTronNodeAddress(node *cre.Node, bc *cre.Blockchain) *common.Address {
	evmKey, ok := node.Keys.EVM[bc.ChainID]
	if !ok {
		return nil // Skip nodes without EVM keys for this chain
	}

	return &evmKey.PublicAddress
}

func FundTronAddress(ctx context.Context, testLogger zerolog.Logger, nodeAddress common.Address, fundingAmount uint64, bc *cre.Blockchain, env *cldf.Environment) error {
	receiverAddress := address.EVMAddressToAddress(nodeAddress)
	testLogger.Info().Msgf("Attempting to fund TRON account %s", nodeAddress)

	tronChains := env.BlockChains.TronChains()
	tronChain, exists := tronChains[bc.ChainSelector]
	if !exists {
		return fmt.Errorf("TRON chain not found for selector %d", bc.ChainSelector)
	}

	tx, err := tronChain.Client.Transfer(tronChain.Address, receiverAddress, libc.MustSafeInt64(fundingAmount))
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to create transfer transaction for TRON node %s", nodeAddress)
	}

	txInfo, err := tronChain.SendAndConfirm(ctx, tx, nil)
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to send and confirm transfer to TRON node %s", nodeAddress)
	}

	testLogger.Info().Msgf("Successfully funded TRON account %s with %d SUN, txHash: %s", receiverAddress.String(), fundingAmount, txInfo.ID)

	return nil
}

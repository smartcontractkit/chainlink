package environment

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"

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

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
)

type PrepareFundCLNodesOpDeps struct {
	TestLogger        zerolog.Logger
	Env               *cldf.Environment
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	DonTopology       *cre.DonTopology
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
		for _, metaDon := range deps.DonTopology.ToDonMetadata() {
			for _, bcOut := range deps.BlockchainOutputs {
				if !flags.RequiresForwarderContract(metaDon.Flags, bcOut.ChainID) && bcOut.SolChain == nil {
					continue
				}

				if bcOut.SolChain != nil {
					requiredFundingPerChain[bcOut.SolChain.ChainSelector] += input.FundingPerChainFamilyForEachNode[chainselectors.FamilySolana] * uint64(len(metaDon.NodesMetadata))
					continue
				}

				if bcOut.BlockchainOutput.Family == blockchain.FamilyTron {
					requiredFundingPerChain[bcOut.ChainSelector] += input.FundingPerChainFamilyForEachNode[chainselectors.FamilyTron] * uint64(len(metaDon.NodesMetadata))
					continue
				}

				requiredFundingPerChain[bcOut.ChainSelector] += input.FundingPerChainFamilyForEachNode[chainselectors.FamilyEVM] * uint64(len(metaDon.NodesMetadata))
			}
		}

		for _, bcOut := range deps.BlockchainOutputs {
			if requiredFundingPerChain[bcOut.ChainSelector] == 0 {
				continue
			}

			chainFamily := chainselectors.FamilyEVM
			if bcOut.SolChain != nil {
				chainFamily = chainselectors.FamilySolana
			} else if bcOut.BlockchainOutput.Family == blockchain.FamilyTron {
				chainFamily = chainselectors.FamilyTron
			}

			switch chainFamily {
			case chainselectors.FamilyEVM:
				if _, exists := output.PrivateKeysPerChainFamily[chainFamily]; !exists {
					output.PrivateKeysPerChainFamily[chainFamily] = make(map[uint64][]byte)
				}

				if _, exists := output.PrivateKeysPerChainFamily[chainFamily][bcOut.ChainSelector]; !exists {
					publicAddress, privateKeyBytes, pkErr := generatePubPrivKeyPairForEth()
					if pkErr != nil {
						return nil, pkgerrors.Wrap(pkErr, "failed to generate pub/priv key pair for EVM funding account")
					}

					fundingAmount := libc.MustSafeInt64(requiredFundingPerChain[bcOut.ChainSelector])
					fundingAmount += (fundingAmount / 5) // add 20% to cover gas fees

					_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bcOut.SethClient, libfunding.FundsToSend{
						ToAddress:  *publicAddress,
						Amount:     big.NewInt(fundingAmount),
						PrivateKey: bcOut.SethClient.MustGetRootPrivateKey(),
					})

					if fundingErr != nil {
						return nil, pkgerrors.Wrapf(fundingErr, "failed to fund funding account %s on chain %d", publicAddress.String(), bcOut.ChainID)
					}

					output.PrivateKeysPerChainFamily[chainFamily][bcOut.ChainSelector] = privateKeyBytes
				}
			case chainselectors.FamilySolana:
				if _, exists := output.PrivateKeysPerChainFamily[chainFamily]; !exists {
					output.PrivateKeysPerChainFamily[chainFamily] = make(map[uint64][]byte)
				}
				if _, exists := output.PrivateKeysPerChainFamily[chainFamily][bcOut.SolChain.ChainSelector]; !exists {
					private, pkErr := solana.NewRandomPrivateKey()
					if pkErr != nil {
						return nil, pkgerrors.Wrap(pkErr, "failed to generate private key for solana")
					}
					public := private.PublicKey()
					fundingErr := libfunding.SendFundsSol(ctx, zerolog.Logger{}, bcOut.SolClient, libfunding.FundsToSendSol{
						Recipent:   public,
						PrivateKey: bcOut.SolChain.PrivateKey,
						Amount:     requiredFundingPerChain[bcOut.SolChain.ChainSelector],
					})
					if fundingErr != nil {
						return nil, pkgerrors.Wrapf(fundingErr, " failed to fund funding accounts on chain %v", bcOut.SolChain.ChainID)
					}

					output.PrivateKeysPerChainFamily[chainFamily][bcOut.SolChain.ChainSelector] = private
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
	BlockchainOutputs []*cre.WrappedBlockchainOutput
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
		for donIndex, donMetadata := range deps.DonTopology.ToDonMetadata() {
			deps.TestLogger.Info().Msgf("Funding nodes for DON %s", donMetadata.Name)
			for _, bcOut := range deps.BlockchainOutputs {
				if !flags.RequiresForwarderContract(donMetadata.Flags, bcOut.ChainID) &&
					bcOut.SolChain == nil { // for now, we can only write to solana, so we consider forwarder is always present
					continue
				}
				for _, node := range deps.DonTopology.Dons.List()[donIndex].Nodes {
					chainFamily := chainselectors.FamilyEVM
					if bcOut.SolChain != nil {
						chainFamily = chainselectors.FamilySolana
					} else if bcOut.BlockchainOutput.Family == blockchain.FamilyTron {
						chainFamily = chainselectors.FamilyTron
					}

					fundingAmount, ok := input.FundingAmountPerChainFamily[chainFamily]
					if !ok {
						return nil, fmt.Errorf("missing funding amount for chain family %s", chainFamily)
					}

					switch chainFamily {
					case chainselectors.FamilyEVM:
						if err := fundEthAddress(ctx, deps.TestLogger, node, fundingAmount, bcOut, input.PrivateKeyPerChainFamily); err != nil {
							return nil, err
						}
					case chainselectors.FamilySolana:
						if err := fundSolanaAddress(ctx, deps.TestLogger, node, fundingAmount, bcOut, input.PrivateKeyPerChainFamily); err != nil {
							return nil, err
						}
					case chainselectors.FamilyTron:
						nodeAddress := getTronNodeAddress(node, bcOut)
						if err := FundTronAddress(ctx, deps.TestLogger, nodeAddress, fundingAmount, bcOut, deps.Env); err != nil {
							return nil, err
						}
					default:
						return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
					}
				}
			}

			deps.TestLogger.Info().Msgf("Funded nodes for DON %s", donMetadata.Name)
		}

		return &FundCLNodesOpOutput{}, nil
	},
)

func fundEthAddress(ctx context.Context, testLogger zerolog.Logger, node devenv.Node, fundingAmount uint64, bcOut *cre.WrappedBlockchainOutput, privateKeyPerChainFamily map[string]map[uint64][]byte) error {
	nodeAddress := node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]
	if nodeAddress == "" {
		return nil // Skip nodes without addresses for this chain
	}

	testLogger.Info().Msgf("Attempting to fund EVM account %s", nodeAddress)

	fundingPrivateKey, ok := privateKeyPerChainFamily["evm"][bcOut.ChainSelector]
	if !ok {
		return fmt.Errorf("missing funding private key for chain familyevm, chain %d", bcOut.ChainID)
	}

	fundingKey, fkErr := crypto.ToECDSA(fundingPrivateKey)
	if fkErr != nil {
		return pkgerrors.Wrap(fkErr, "failed to convert funding private key to ECDSA")
	}

	_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bcOut.SethClient, libfunding.FundsToSend{
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

func fundSolanaAddress(ctx context.Context, testLogger zerolog.Logger, node devenv.Node, fundingAmount uint64, bcOut *cre.WrappedBlockchainOutput, _ map[string]map[uint64][]byte) error {
	funder := bcOut.SolChain.PrivateKey
	recipient := solana.MustPublicKeyFromBase58(node.AccountAddr[bcOut.SolChain.ChainID])
	testLogger.Info().Msgf("Attempting to fund Solana account %s", recipient.String())

	err := libfunding.SendFundsSol(ctx, zerolog.Logger{}, bcOut.SolClient, libfunding.FundsToSendSol{
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

func getTronNodeAddress(node devenv.Node, bcOut *cre.WrappedBlockchainOutput) common.Address {
	nodeAddress := node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]
	if nodeAddress == "" {
		return common.Address{} // Skip nodes without addresses for this chain
	}

	return common.HexToAddress(nodeAddress)
}

func FundTronAddress(ctx context.Context, testLogger zerolog.Logger, nodeAddress common.Address, fundingAmount uint64, bcOut *cre.WrappedBlockchainOutput, env *cldf.Environment) error {
	receiverAddress := address.EVMAddressToAddress(nodeAddress)

	testLogger.Info().Msgf("Attempting to fund TRON account %s", nodeAddress)

	tronChains := env.BlockChains.TronChains()
	tronChain, exists := tronChains[bcOut.ChainSelector]
	if !exists {
		return fmt.Errorf("TRON chain not found for selector %d", bcOut.ChainSelector)
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

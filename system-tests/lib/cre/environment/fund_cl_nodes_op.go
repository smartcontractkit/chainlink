package environment

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
)

type FundCLNodesOpDeps struct {
	Env               *cldf.Environment
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	DonTopology       *cre.DonTopology
}

type FundCLNodesOpInput struct {
	FundAmount int64
}

type FundCLNodesOpOutput struct {
}

var FundCLNodesOp = operations.NewOperation[FundCLNodesOpInput, FundCLNodesOpOutput, FundCLNodesOpDeps](
	"fund-cl-nodes-op",
	semver.MustParse("1.0.0"),
	"Fund Chainlink Nodes",
	func(b operations.Bundle, deps FundCLNodesOpDeps, input FundCLNodesOpInput) (FundCLNodesOpOutput, error) {
		ctx := b.GetContext()
		// Fund the nodes
		concurrentNonceMap, concurrentNonceMapErr := NewConcurrentNonceMap(ctx, deps.BlockchainOutputs)
		if concurrentNonceMapErr != nil {
			return FundCLNodesOpOutput{}, pkgerrors.Wrap(concurrentNonceMapErr, "failed to create concurrent nonce map")
		}

		// Decrement the nonce for each chain, because we will increment it in the next loop
		for _, bcOut := range deps.BlockchainOutputs {
			concurrentNonceMap.Decrement(bcOut.ChainID)
		}

		errGroup := &errgroup.Group{}
		for _, metaDon := range deps.DonTopology.DonsWithMetadata {
			for _, bcOut := range deps.BlockchainOutputs {
				if !flags.RequiresForwarderContract(metaDon.Flags, bcOut.ChainID) {
					continue
				}
				for _, node := range metaDon.DON.Nodes {
					errGroup.Go(func() error {
						if bcOut.SolChain != nil {
							// handle funding here
							// TODO move it to lib funding and add confirmation
							fmt.Println("fund solana node")
							funder := bcOut.SolChain.PrivateKey
							recipient := solana.MustPublicKeyFromBase58(node.AccountAddr[bcOut.SolChain.ChainID])
							if recipient.IsZero() {
								fmt.Println("solana addr not found")
								return nil
							}
							bal, _ := bcOut.SolClient.GetBalance(ctx, funder.PublicKey(), rpc.CommitmentConfirmed)
							fmt.Println("sender balance:", bal.Value, "recipient:", recipient)

							recent, err := bcOut.SolClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
							if err != nil {
								return fmt.Errorf("failed to fund node: %w", err)
							}

							tx, err := solana.NewTransaction([]solana.Instruction{
								system.NewTransferInstruction(
									50_000_000,
									funder.PublicKey(),
									recipient,
								).Build(),
							},
								recent.Value.Blockhash,
								solana.TransactionPayer(funder.PublicKey()))
							if err != nil {
								return fmt.Errorf("failed to build fund transaction: %w", err)
							}

							_, err = tx.Sign(
								func(key solana.PublicKey) *solana.PrivateKey {
									if funder.PublicKey().Equals(key) {
										return &funder
									}
									return nil
								},
							)
							if err != nil {
								return fmt.Errorf("failed to sign fund transaction: %w", err)
							}

							_, err = bcOut.SolClient.SendTransaction(ctx, tx)
							if err != nil {
								return fmt.Errorf("failed to send fund transaction: %w", err)
							}

							return nil
						}

						nodeAddress := node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]
						if nodeAddress == "" {
							return nil
						}

						nonce := concurrentNonceMap.Increment(bcOut.ChainID)

						_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bcOut.SethClient, libfunding.FundsToSend{
							ToAddress:  common.HexToAddress(nodeAddress),
							Amount:     big.NewInt(input.FundAmount),
							PrivateKey: bcOut.SethClient.MustGetRootPrivateKey(),
							Nonce:      ptr.Ptr(nonce),
						})
						if fundingErr != nil {
							return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", nodeAddress)
						}
						return nil
					})
				}
			}
		}

		if err := errGroup.Wait(); err != nil {
			return FundCLNodesOpOutput{}, pkgerrors.Wrap(err, "failed to fund nodes")
		}

		return FundCLNodesOpOutput{}, nil
	},
)

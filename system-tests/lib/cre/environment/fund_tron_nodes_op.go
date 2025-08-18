package environment

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	pkgerrors "github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

type FundTronNodesOpDeps struct {
	Env               *cldf.Environment
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	DonTopology       *cre.DonTopology
}

type FundTronNodesOpInput struct {
	FundAmount int64 // Amount in SUN (1 TRX = 1_000_000 SUN)
}

type FundTronNodesOpOutput struct {
}

var FundTronNodesOp = operations.NewOperation[FundTronNodesOpInput, FundTronNodesOpOutput, FundTronNodesOpDeps](
	"fund-tron-nodes-op",
	semver.MustParse("1.0.0"),
	"Fund Tron Chainlink Nodes",
	func(b operations.Bundle, deps FundTronNodesOpDeps, input FundTronNodesOpInput) (FundTronNodesOpOutput, error) {
		ctx := b.GetContext()
		logger := b.Logger

		logger.Infow("Starting Tron node funding",
			"fundAmount", input.FundAmount,
			"fundAmountTRX", float64(input.FundAmount)/1_000_000)

		errGroup := &errgroup.Group{}
		for _, metaDon := range deps.DonTopology.DonsWithMetadata {
			for _, bcOut := range deps.BlockchainOutputs {
				// All blockchain outputs passed to this operation are Tron chains
				if !flags.RequiresForwarderContract(metaDon.Flags, bcOut.ChainID) {
					continue
				}

				// Get chain selector for this Tron chain
				chainSelector, err := chainselectors.SelectorFromChainId(bcOut.ChainID)
				if err != nil {
					return FundTronNodesOpOutput{}, pkgerrors.Wrapf(err, "failed to get chain selector for chain ID %d", bcOut.ChainID)
				}

				expectedAddressKey := node.AddressKeyFromSelector(chainSelector)

				// Use the simple approach: get node addresses from topology metadata
				for _, nodeMetadata := range metaDon.NodesMetadata {
					errGroup.Go(func() error {
						// Use the simple FindLabelValue function to get the address
						nodeAddress, err := node.FindLabelValue(nodeMetadata, expectedAddressKey)
						if err != nil {
							// Skip nodes that don't have an address for this chain
							return nil
						}

						// Convert hex address to Tron address
						receiverAddress := address.EVMAddressToAddress(common.HexToAddress(nodeAddress))

						logger.Infow("Funding Tron node",
							"chainID", bcOut.ChainID,
							"nodeAddress", nodeAddress,
							"receiverAddress", receiverAddress.String(),
							"amount", input.FundAmount,
							"amountTRX", float64(input.FundAmount)/1_000_000)

						// Create transfer transaction
						tx, err := bcOut.TronChain.Client.Transfer(bcOut.TronChain.Address, receiverAddress, input.FundAmount)
						if err != nil {
							return pkgerrors.Wrapf(err, "failed to create transfer transaction for node %s", nodeAddress)
						}

						// Send and confirm the transaction
						txInfo, err := bcOut.TronChain.SendAndConfirm(ctx, tx, nil)
						if err != nil {
							return pkgerrors.Wrapf(err, "failed to send and confirm transfer to node %s", nodeAddress)
						}

						logger.Infow("Successfully funded Tron node",
							"chainID", bcOut.ChainID,
							"nodeAddress", nodeAddress,
							"receiverAddress", receiverAddress.String(),
							"amount", input.FundAmount,
							"amountTRX", float64(input.FundAmount)/1_000_000,
							"txHash", txInfo.ID)

						return nil
					})
				}
			}
		}

		if err := errGroup.Wait(); err != nil {
			return FundTronNodesOpOutput{}, pkgerrors.Wrap(err, "failed to fund Tron nodes")
		}

		logger.Infow("Tron node funding completed successfully")
		return FundTronNodesOpOutput{}, nil
	},
)

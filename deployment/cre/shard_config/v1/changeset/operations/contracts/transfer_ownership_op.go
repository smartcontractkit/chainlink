package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	shard_config "github.com/smartcontractkit/chainlink-evm/contracts/cre/gobindings/shardconfig/generated/v1_0_0/shard_config"
)

// TransferOwnershipOpInput contains the input parameters for transferring ownership.
type TransferOwnershipOpInput struct {
	ChainSelector uint64         `json:"chainSelector"`
	ContractAddr  string         `json:"contractAddr"`
	NewOwner      common.Address `json:"newOwner"`
}

// TransferOwnershipOpOutput contains the output of transferring ownership.
type TransferOwnershipOpOutput struct {
	TxHash string `json:"txHash"`
}

// TransferOwnershipOpDeps contains the dependencies for the transfer ownership operation.
type TransferOwnershipOpDeps struct {
	Env *cldf.Environment
}

// TransferOwnershipOp is an operation that transfers ownership of a ShardConfig contract.
var TransferOwnershipOp = operations.NewOperation(
	"transfer-ownership-v1-op",
	semver.MustParse("1.0.0"),
	"Transfer ShardConfig ownership",
	func(b operations.Bundle, deps TransferOwnershipOpDeps, input TransferOwnershipOpInput) (TransferOwnershipOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return TransferOwnershipOpOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		contract, err := shard_config.NewShardConfig(
			common.HexToAddress(input.ContractAddr),
			chain.Client,
		)
		if err != nil {
			return TransferOwnershipOpOutput{}, fmt.Errorf("failed to bind ShardConfig: %w", err)
		}

		tx, err := contract.TransferOwnership(chain.DeployerKey, input.NewOwner)
		if err != nil {
			return TransferOwnershipOpOutput{}, fmt.Errorf("failed to transfer ownership: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return TransferOwnershipOpOutput{}, fmt.Errorf("failed to confirm tx: %w", err)
		}

		deps.Env.Logger.Infof("Transferred ShardConfig ownership to %s at %s on chain %d", input.NewOwner.Hex(), input.ContractAddr, chain.Selector)

		return TransferOwnershipOpOutput{TxHash: tx.Hash().String()}, nil
	},
)

package contracts

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	shard_config "github.com/smartcontractkit/chainlink-evm/contracts/cre/gobindings/shardconfig/generated/v1_0_0/shard_config"
)

// UpdateShardCountOpInput contains the input parameters for updating the shard count.
type UpdateShardCountOpInput struct {
	ChainSelector uint64 `json:"chainSelector"`
	ContractAddr  string `json:"contractAddr"`
	NewShardCount uint64 `json:"newShardCount"`
}

// UpdateShardCountOpOutput contains the output of updating the shard count.
type UpdateShardCountOpOutput struct {
	TxHash string `json:"txHash"`
}

// UpdateShardCountOpDeps contains the dependencies for the update operation.
type UpdateShardCountOpDeps struct {
	Env *cldf.Environment
}

// UpdateShardCountOp is an operation that updates the desired shard count on a ShardConfig contract.
var UpdateShardCountOp = operations.NewOperation(
	"update-shard-count-v1-op",
	semver.MustParse("1.0.0"),
	"Update ShardConfig shard count",
	func(b operations.Bundle, deps UpdateShardCountOpDeps, input UpdateShardCountOpInput) (UpdateShardCountOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return UpdateShardCountOpOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		contract, err := shard_config.NewShardConfig(
			common.HexToAddress(input.ContractAddr),
			chain.Client,
		)
		if err != nil {
			return UpdateShardCountOpOutput{}, fmt.Errorf("failed to bind ShardConfig: %w", err)
		}

		tx, err := contract.SetDesiredShardCount(chain.DeployerKey, big.NewInt(int64(input.NewShardCount)))
		if err != nil {
			return UpdateShardCountOpOutput{}, fmt.Errorf("failed to set shard count: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return UpdateShardCountOpOutput{}, fmt.Errorf("failed to confirm tx: %w", err)
		}

		deps.Env.Logger.Infof("Updated ShardConfig shard count to %d at %s on chain %d", input.NewShardCount, input.ContractAddr, chain.Selector)

		return UpdateShardCountOpOutput{TxHash: tx.Hash().String()}, nil
	},
)

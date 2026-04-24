package contracts

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	shard_config "github.com/smartcontractkit/chainlink-evm/contracts/cre/gobindings/shardconfig/generated/v1_0_0/shard_config"
)

// DeployShardConfigOpInput contains the input parameters for deploying a ShardConfig contract.
type DeployShardConfigOpInput struct {
	ChainSelector     uint64 `json:"chainSelector"`
	InitialShardCount uint64 `json:"initialShardCount"`
	Qualifier         string `json:"qualifier,omitempty"`
}

// DeployShardConfigOpOutput contains the output of deploying a ShardConfig contract.
type DeployShardConfigOpOutput struct {
	Address       string   `json:"address"`
	ChainSelector uint64   `json:"chainSelector"`
	Labels        []string `json:"labels"`
	Qualifier     string   `json:"qualifier"`
	Type          string   `json:"type"`
	Version       string   `json:"version"`
}

// DeployShardConfigOpDeps contains the dependencies for the deploy operation.
type DeployShardConfigOpDeps struct {
	Env *cldf.Environment
}

// DeployShardConfigOp is an operation that deploys the ShardConfig contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployShardConfigOp = operations.NewOperation(
	"deploy-shard-config-v1-op",
	semver.MustParse("1.0.0"),
	"Deploy ShardConfig Contract",
	func(b operations.Bundle, deps DeployShardConfigOpDeps, input DeployShardConfigOpInput) (DeployShardConfigOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployShardConfigOpOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		addr, tx, _, err := shard_config.DeployShardConfig(
			chain.DeployerKey,
			chain.Client,
			big.NewInt(int64(input.InitialShardCount)),
		)
		if err != nil {
			return DeployShardConfigOpOutput{}, fmt.Errorf("failed to deploy ShardConfig: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return DeployShardConfigOpOutput{}, fmt.Errorf("failed to confirm deployment: %w", err)
		}

		// Bind to the deployed contract to get type and version
		shardConfig, err := shard_config.NewShardConfig(addr, chain.Client)
		if err != nil {
			return DeployShardConfigOpOutput{}, fmt.Errorf("failed to bind to ShardConfig: %w", err)
		}

		tvStr, err := shardConfig.TypeAndVersion(&bind.CallOpts{})
		if err != nil {
			return DeployShardConfigOpOutput{}, fmt.Errorf("failed to get type and version: %w", err)
		}

		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return DeployShardConfigOpOutput{}, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}

		deps.Env.Logger.Infof("Deployed %s on chain %d at %s", tv.String(), chain.Selector, addr.String())

		return DeployShardConfigOpOutput{
			Address:       addr.String(),
			ChainSelector: input.ChainSelector,
			Labels:        tv.Labels.List(),
			Qualifier:     input.Qualifier,
			Type:          string(tv.Type),
			Version:       tv.Version.String(),
		}, nil
	},
)

package sharding

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"

	deployment_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	shard_config_changeset "github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

type SetupShardingInput struct {
	Logger   zerolog.Logger
	CreEnv   *cre.Environment
	Topology *cre.Topology
}

func SetupSharding(input SetupShardingInput) error {
	deployInput := shard_config_changeset.DeployShardConfigInput{
		ChainSelector:     input.CreEnv.RegistryChainSelector,
		InitialShardCount: uint64(input.Topology.DonsMetadata.ShardCount()),
	}

	vErr := shard_config_changeset.DeployShardConfig{}.VerifyPreconditions(*input.CreEnv.CldfEnvironment, deployInput)
	if vErr != nil {
		return fmt.Errorf("preconditions verification for Shard Config contract failed: %w", vErr)
	}

	out, dErr := shard_config_changeset.DeployShardConfig{}.Apply(*input.CreEnv.CldfEnvironment, deployInput)
	if dErr != nil {
		return fmt.Errorf("failed to deploy Shard Config contract: %w", dErr)
	}

	crecontracts.MergeAllDataStores(input.CreEnv, out)

	shardConfigAddr := crecontracts.MustGetAddressFromDataStore(input.CreEnv.CldfEnvironment.DataStore, input.CreEnv.RegistryChainSelector, deployment_contracts.ShardConfig.String(), semver.MustParse("1"), "")
	input.Logger.Info().Msgf("Deployed Shard Config v1 contract on chain %d at %s", input.CreEnv.RegistryChainSelector, shardConfigAddr)

	// once we have a changeset for job creation:
	// get shard leader DON
	// create job(s) on shard leader DON

	return nil
}

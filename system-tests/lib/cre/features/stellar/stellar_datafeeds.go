package stellar

import (
	"context"
	"fmt"

	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	datafeeds "github.com/smartcontractkit/chainlink-stellar/deployment/data-feeds"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
)

func DeployAndConfigureStellarDataFeedsCache(
	ctx context.Context,
	chain *stellchain.Blockchain,
	creEnv *cre.Environment,
	dataID [32]byte,
	description string,
	workflowOwner [20]byte,
	workflowName [10]byte,
) (string, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return "", err
	}
	owner := stellarChain.Signer.Address()
	if fundErr := chain.Fund(ctx, owner, 0); fundErr != nil {
		return "", fmt.Errorf("failed to fund stellar deployer %s via friendbot: %w", owner, fundErr)
	}
	forwarder, ok := forwarderAddress(creEnv.CldfEnvironment.DataStore, chain.ChainSelector())
	if !ok {
		return "", fmt.Errorf("missing Stellar forwarder in datastore for chain selector %d", chain.ChainSelector())
	}
	deployer, err := stellardeployment.NewDeployerFromChain(stellarChain)
	if err != nil {
		return "", fmt.Errorf("failed to build stellar deployer: %w", err)
	}
	wasm, err := datafeeds.Artifact(datafeeds.DataFeedsCacheWasm)
	if err != nil {
		return "", fmt.Errorf("failed to source data feeds cache WASM: %w", err)
	}

	salt := [32]byte{0x44} // 'D', distinct from the receiver and read-fixture salts
	cacheID, err := deployer.DeployContractBytesWithArgs(ctx, wasm, salt, []xdr.ScVal{scval.AddressToScVal(owner)})
	if err != nil {
		return "", fmt.Errorf("failed to deploy stellar data feeds cache: %w", err)
	}

	client := data_feeds_cache.NewDataFeedsCacheClient(deployer, cacheID)
	if err := client.AddFeedAdmin(ctx, owner); err != nil {
		return "", fmt.Errorf("failed to add stellar data feeds admin: %w", err)
	}
	entry := data_feeds_cache.FeedConfigEntry{
		DataId: dataID,
		Config: data_feeds_cache.FeedConfig{
			Description: description,
			WorkflowPermissions: []data_feeds_cache.WorkflowPermission{{
				AllowedSender:        forwarder,
				AllowedWorkflowOwner: workflowOwner,
				AllowedWorkflowName:  workflowName,
			}},
		},
	}
	if err := client.SetFeedConfigs(ctx, owner, []data_feeds_cache.FeedConfigEntry{entry}); err != nil {
		return "", fmt.Errorf("failed to set stellar data feeds config: %w", err)
	}
	return cacheID, nil
}

func DataFeedsCacheLatestRound(ctx context.Context, chain *stellchain.Blockchain, cacheID string, dataID [32]byte) (*data_feeds_cache.RoundData, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return nil, err
	}
	deployer, err := stellardeployment.NewDeployerFromChain(stellarChain)
	if err != nil {
		return nil, fmt.Errorf("failed to build stellar deployer: %w", err)
	}
	rounds, err := data_feeds_cache.NewDataFeedsCacheClient(deployer, cacheID).LatestRound(ctx, [][32]byte{dataID})
	if err != nil {
		return nil, fmt.Errorf("failed to read stellar data feeds latest round: %w", err)
	}
	return rounds[0], nil
}

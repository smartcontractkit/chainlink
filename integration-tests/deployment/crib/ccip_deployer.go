package crib

import (
	"context"
	"errors"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can be started with the CR address in Capabilities.ExternalRegistry
func DeployHomeChainContracts(lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64) (deployment.CapabilityRegistryConfig, deployment.AddressBook, error) {
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	ab := deployment.NewMemoryAddressBook()
	capReg, err := ccipdeployment.DeployCapReg(lggr, ab, chains[homeChainSel])
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	evmChainID, err := chainsel.ChainIdFromSelector(homeChainSel)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	return deployment.CapabilityRegistryConfig{
		NetworkType: relay.NetworkEVM,
		EVMChainID:  evmChainID,
		Contract:    capReg.Address,
	}, ab, nil
}

func DeployCCIPAndAddLanes(lggr logger.Logger, envCfg devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) error {
	e, _, err := devenv.NewEnvironment(context.Background(), lggr, envCfg)
	if err != nil {
		return err
	}
	if e == nil {
		return errors.New("environment is nil")
	}

	_, err = ccipdeployment.DeployFeeds(lggr, ab, e.Chains[feedChainSel])
	if err != nil {
		return err
	}
	feeTokenContracts, err := ccipdeployment.DeployFeeTokensToChains(lggr, ab, e.Chains)
	if err != nil {
		return err
	}
	tenv := ccipdeployment.DeployedEnv{
		Env:               *e,
		HomeChainSel:      homeChainSel,
		FeedChainSel:      feedChainSel,
		FeeTokenContracts: feeTokenContracts,
		Ab:                ab,
	}

	state, err := ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	if err != nil {
		return err
	}
	if state.Chains[tenv.HomeChainSel].LinkToken == nil {
		return errors.New("link token not deployed")
	}

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccipdeployment.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccipdeployment.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccipdeployment.LinkSymbol].Address().String()),
			Decimals:          ccipdeployment.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	tokenConfig.UpsertTokenInfo(ccipdeployment.WethSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccipdeployment.WethSymbol].Address().String()),
			Decimals:          ccipdeployment.WethDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	mcmsCfg, err := ccipdeployment.NewTestMCMSConfig(tenv.Env)
	if err != nil {
		return err
	}
	output, err := changeset.InitialDeployChangeSet(tenv.Ab, tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:       tenv.HomeChainSel,
		FeedChainSel:       tenv.FeedChainSel,
		ChainsToDeploy:     tenv.Env.AllChainSelectors(),
		TokenConfig:        tokenConfig,
		MCMSConfig:         mcmsCfg,
		CapabilityRegistry: state.Chains[tenv.HomeChainSel].CapabilityRegistry.Address(),
		FeeTokenContracts:  tenv.FeeTokenContracts,
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	if err != nil {
		return err
	}
	// Get new state after migration.
	state, err = ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	if err != nil {
		return err
	}

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := tenv.Env.Offchain.ProposeJob(context.Background(),
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			if err != nil {
				return err
			}
		}
	}

	// Add all lanes
	return ccipdeployment.AddLanesForAll(tenv.Env, state)
}

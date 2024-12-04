package crib

import (
	"context"
	"errors"
	chainsel "github.com/smartcontractkit/chain-selectors"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"golang.org/x/exp/maps"
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
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
	capReg, err := changeset.DeployCapReg(lggr,
		// deploying cap reg for the first time on a blank chain state
		changeset.CCIPOnChainState{
			Chains: make(map[uint64]changeset.CCIPChainState),
		}, ab, chains[homeChainSel])
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

func DeployCCIPAndAddLanes(lggr logger.Logger, envCfg devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	ctx := context.Background()
	e, _, err := devenv.NewEnvironment(func() context.Context { return context.Background() }, lggr, envCfg)
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	if e == nil {
		return DeployCCIPOutput{}, errors.New("environment is nil")
	}

	_, err = changeset.DeployFeeds(lggr, ab, e.Chains[feedChainSel], big.NewInt(9000000), big.NewInt(9000000))
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	e.ExistingAddresses = ab
	tenv := changeset.DeployedEnv{
		Env:          *e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedChainSel,
	}
	chains := tenv.Env.AllChainSelectors()
	out, err := changeset.DeployPrerequisites(tenv.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: chains,
	})
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	err = tenv.Env.ExistingAddresses.Merge(out.AddressBook)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	state, err := changeset.LoadOnchainState(tenv.Env)
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	if state.Chains[tenv.HomeChainSel].LinkToken == nil {
		return DeployCCIPOutput{}, errors.New("link token not deployed")
	}

	//  Deploy contracts to new chain
	cfg, err := commonchangeset.CreateMCMSConfig(tenv.Env.AllDeployerKeys())
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	var mcmsConfigs = make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range chains {
		mcmsConfigs[chain] = cfg
	}
	out, err = commonchangeset.DeployMCMSWithTimelock(tenv.Env, mcmsConfigs)
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	err = tenv.Env.ExistingAddresses.Merge(out.AddressBook)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	chainConfig := make(map[uint64]changeset.CCIPOCRParams)
	for _, chain := range chains {
		chainConfig[chain] = changeset.DefaultOCRParams(tenv.FeedChainSel, nil, nil)
	}
	err = changeset.DeployCCIPContracts(tenv.Env, tenv.Env.ExistingAddresses, changeset.NewChainsConfig{
		HomeChainSel:       tenv.HomeChainSel,
		FeedChainSel:       tenv.FeedChainSel,
		ChainConfigByChain: chainConfig,
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	// Get new state after migration.

	state, err = changeset.LoadOnchainState(tenv.Env)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	out, err = changeset.CCIPCapabilityJobspec(tenv.Env, struct{}{})
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	for nodeID, jobs := range out.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := tenv.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			if err != nil {
				return DeployCCIPOutput{}, err
			}
		}
	}

	// Add all lanes
	err = changeset.AddLanesForAll(tenv.Env, state)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := ab.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     maps.Keys(out.JobSpecs),
	}, err
}

package crib

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can be started with the CR address in Capabilities.ExternalRegistry
// DeployHomeChainContracts is to 1. Set up crib with chains and chainlink nodes ( cap reg is not known yet so not setting the config with capreg address)
// Call DeployHomeChain changeset with nodeinfo ( the peer id and all)
func DeployHomeChainContracts(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64, feedChainSel uint64) (deployment.CapabilityRegistryConfig, deployment.AddressBook, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if e == nil {
		return deployment.CapabilityRegistryConfig{}, nil, errors.New("environment is nil")
	}

	tenv := changeset.DeployedEnv{
		Env:          *e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedChainSel,
	}

	// Call DeployHomeChain changeset with nodeinfo ( the peer id and all)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("Failed to get node info from env", err)
	}
	p2pIds := nodes.NonBootstraps().PeerIDs()
	out, err := changeset.DeployHomeChain(tenv.Env, changeset.DeployHomeChainConfig{
		HomeChainSel:     homeChainSel,
		RMNStaticConfig:  changeset.NewTestRMNStaticConfig(),
		RMNDynamicConfig: changeset.NewTestRMNDynamicConfig(),
		NodeOperators:    changeset.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
		NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
			"NodeOperator": p2pIds,
		},
	})
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("Failed to deploy home chain", err)
	}
	err = tenv.Env.ExistingAddresses.Merge(out.AddressBook)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("Failed to merge addresses after deploying home chain", err)
	}

	fmt.Printf("Deployed home chain contracts\n")

	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("Failed to load on chain state", err)
	}
	capRegAddr := state.Chains[homeChainSel].CapabilityRegistry.Address()
	if capRegAddr == common.HexToAddress("0x") {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("Cap Reg address not found")
	}
	capRegConfig := deployment.CapabilityRegistryConfig{
		EVMChainID:  homeChainSel,
		Contract:    state.Chains[homeChainSel].CapabilityRegistry.Address(),
		NetworkType: relay.NetworkEVM,
	}
	return capRegConfig, e.ExistingAddresses, nil
}

func DeployCCIPAndAddLanes(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to initiate new environment", err)
	}
	e.ExistingAddresses = ab
	allChainIds := e.AllChainSelectors()

	out, err := commonchangeset.DeployLinkToken(*e, allChainIds)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to deploy link token", err)
	}
	err = e.ExistingAddresses.Merge(out.AddressBook)

	// deploy pre requisites
	out, err = changeset.DeployPrerequisites(*e, changeset.DeployPrerequisiteConfig{
		ChainSelectors: allChainIds,
	})
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to deploy prerequisites", err)
	}
	fmt.Printf("PreReq AddressBook %+v\n", out.AddressBook)
	err = e.ExistingAddresses.Merge(out.AddressBook)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to merge addresses after deploying prereqs", err)
	}

	// deploy MCMS With Timelock
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain], err = commonchangeset.CreateMCMSConfig(e.AllDeployerKeys())
		if err != nil {
			return DeployCCIPOutput{}, fmt.Errorf("Failed to create MCMS config for %d", chain, err)
		}
	}
	out, err = commonchangeset.DeployMCMSWithTimelock(*e, cfg)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to deploy MCMS with timelock", err)
	}
	err = e.ExistingAddresses.Merge(out.AddressBook)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to merge addresses after deploying MCMS with timelock", err)
	}

	// deploy ccip chain contracts
	out, err = changeset.DeployChainContracts(*e, changeset.DeployChainContractsConfig{
		ChainSelectors:    allChainIds,
		HomeChainSelector: homeChainSel,
	})
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to deploy chain contracts", err)
	}
	err = e.ExistingAddresses.Merge(out.AddressBook)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to merge addresses after deploying chain contracts", err)
	}

	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to load onchain state", err)
	}

	out, err = changeset.CCIPCapabilityJobspec(*e, struct{}{})
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to get CCIP capability jobspec", err)
	}
	for nodeID, jobs := range out.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			if err != nil {
				return DeployCCIPOutput{}, fmt.Errorf("failed to propose job: %w", err)
			}
		}
	}

	// Add all lanes
	err = changeset.AddLanesForAll(*e, state)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to add lanes", err)
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("Failed to get convert address book to address book map", err)
	}
	fmt.Printf("New Addresses %+v\n", addresses)
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     maps.Keys(out.JobSpecs),
	}, err
}

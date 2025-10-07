package environment

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	deployment_devenv "github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	docker_blockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/docker"
	k8s_blockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/kubernetes"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

// BuildFromSavedState rebuilds the CLDF environment and per‑chain clients from
// artifacts produced by a previous local CRE run.
// Inputs:
//   - cachedInput: outputs from starting the environment via CTFv2 configs
//     (node sets, Job Distributor, blockchain nodes).
//   - envArtifact: CLDF deployment output including JD config and DON
//     topology/metadata.
//
// Artifact paths are recorded in `artifact_paths.json` in the environment
// directory (typically `core/scripts/cre/environment`).
// Returns the reconstructed CLDF environment, wrapped blockchain outputs, and an error.
func BuildFromSavedState(ctx context.Context, cldLogger logger.Logger, cachedInput *envconfig.Config, envArtifact *EnvArtifact) (*cre.Environment, []*cre.Blockchain, error) {
	if cachedInput == nil {
		return nil, nil, errors.New("cached input cannot be nil")
	}

	if envArtifact == nil {
		return nil, nil, errors.New("environment artifact cannot be nil")
	}

	var blockchainDeployers map[blockchain.ChainFamily]blockchains.Deployer
	if cachedInput.Infra.IsDocker() {
		blockchainDeployers = docker_blockchains.NewDeployerSet()
	} else {
		blockchainDeployers = k8s_blockchains.NewDeployerSet(framework.L, cachedInput.Infra.CRIB.Namespace, CribConfigsDir)
	}

	deployedBlockchains, startErr := blockchains.Start(
		cldLogger,
		cachedInput.Blockchains,
		blockchainDeployers,
	)
	if startErr != nil {
		return nil, nil, errors.Wrap(startErr, "failed to start blockchains")
	}

	addressBook := cldf.NewMemoryAddressBookFromMap(envArtifact.AddressBook)
	datastore := datastore.NewMemoryDataStore()
	for _, addrRef := range envArtifact.AddressRefs {
		addErr := datastore.AddressRefStore.Add(addrRef)
		if addErr != nil {
			return nil, nil, errors.Wrapf(addErr, "failed to add address ref to datastore %v", addrRef)
		}
	}

	allNodeInfo := make([]deployment_devenv.NodeInfo, 0)
	allNodeIDs := make([]string, 0)
	devenvDons := make([]*deployment_devenv.DON, 0, len(envArtifact.DONs))

	for idx, don := range envArtifact.DONs {
		_, ok := envArtifact.Nodes[don.DonName]
		if !ok {
			return nil, nil, errors.Errorf("no nodes found for don %s", don.DonName)
		}

		for id := range envArtifact.Nodes[don.DonName].Nodes {
			allNodeIDs = append(allNodeIDs, id)
		}

		// a maximum of 1 bootstrap is supported due to environment constraints
		bootstrapNodesCount := 0
		if envArtifact.Topology.ToDonMetadata()[idx].ContainsBootstrapNode() {
			bootstrapNodesCount = 1
		}

		nodeInfo, err := crenode.GetNodeInfo(cachedInput.NodeSets[idx].Out, cachedInput.NodeSets[idx].Name, don.DonID, bootstrapNodesCount)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get node info for don %s", don.DonName)
		}
		offChain, offChainErr := deployment_devenv.NewJDClient(ctx, deployment_devenv.JDConfig{
			WSRPC:    envArtifact.JdConfig.ExternalGRPCUrl,
			GRPC:     envArtifact.JdConfig.ExternalGRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		})
		if offChainErr != nil {
			return nil, nil, errors.Wrapf(offChainErr, "failed to create offchain client for don %s", don.DonName)
		}

		jd, ok := offChain.(*deployment_devenv.JobDistributor)
		if !ok {
			return nil, nil, errors.Errorf("offchain client is not a JobDistributor for don %s", don.DonName)
		}
		registeredDon, donErr := deployment_devenv.NewRegisteredDON(ctx, nodeInfo, *jd)
		if donErr != nil {
			return nil, nil, errors.Wrapf(donErr, "failed to create DON for don %s", don.DonName)
		}

		devenvDons = append(devenvDons, registeredDon)
		allNodeInfo = append(allNodeInfo, nodeInfo...)
	}

	donsMetadata, metaErr := cre.NewDonsMetadata(envArtifact.Topology.ToDonMetadata(), *cachedInput.Infra)
	if metaErr != nil {
		return nil, nil, errors.Wrapf(metaErr, "failed to recreate dons metadata from artifact")
	}

	dons, donsErr := cre.NewDons(donsMetadata, devenvDons)
	if donsErr != nil {
		return nil, nil, errors.Wrapf(donsErr, "failed to create Dons from metadata")
	}

	offChain, offChainErr := deployment_devenv.NewJDClient(ctx, deployment_devenv.JDConfig{
		WSRPC:    envArtifact.JdConfig.ExternalGRPCUrl,
		GRPC:     envArtifact.JdConfig.ExternalGRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: allNodeInfo,
	})
	if offChainErr != nil {
		return nil, nil, errors.Wrapf(offChainErr, "failed to create offchain client")
	}

	chainConfigs := make([]deployment_devenv.ChainConfig, 0, len(deployedBlockchains.Outputs))
	for _, output := range deployedBlockchains.Outputs {
		cfg, cfgErr := cre.ChainConfigFromWrapped(output)
		if cfgErr != nil {
			return nil, nil, errors.Wrapf(cfgErr, "failed to build chain config from write for blockchain %s", output.CtfOutput.Family)
		}
		chainConfigs = append(chainConfigs, cfg)
	}

	blockChains, chainErr := deployment_devenv.NewChains(cldLogger, chainConfigs)
	if chainErr != nil {
		return nil, nil, errors.Wrapf(chainErr, "failed to create block chains")
	}

	cldEnv := cldf.NewEnvironment(
		"cre",
		cldLogger,
		addressBook,
		datastore.Seal(),
		allNodeIDs,
		offChain,
		func() context.Context {
			return ctx
		},
		focr.XXXGenerateTestOCRSecrets(),
		blockChains,
	)

	topology, tErr := cre.NewTopology(cachedInput.NodeSets, *cachedInput.Infra)
	if tErr != nil {
		return nil, nil, errors.Wrap(tErr, "failed to recreate topology from artifact")
	}

	return &cre.Environment{
		CldfEnvironment: cldEnv,
		DonTopology:     cre.NewDonTopology(envArtifact.Topology.HomeChainSelector, topology, dons),
	}, deployedBlockchains.Outputs, nil
}

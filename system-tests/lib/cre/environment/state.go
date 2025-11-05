package environment

import (
	"context"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	blockchain_sets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/sets"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
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
func BuildFromSavedState(ctx context.Context, cldLogger logger.Logger, cachedInput *envconfig.Config, envArtifact *EnvArtifact) (*cre.Environment, *cre.Dons, error) {
	if cachedInput == nil {
		return nil, nil, errors.New("cached input cannot be nil")
	}

	if envArtifact == nil {
		return nil, nil, errors.New("environment artifact cannot be nil")
	}

	blockchainDeployers := blockchain_sets.NewDeployerSet(framework.L, cachedInput.Infra, infra.CribConfigsDir)
	deployedBlockchains, startErr := blockchains.Start(
		ctx,
		framework.L,
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

	allNodeIDs := make([]string, 0)
	donsSlice := make([]*cre.Don, 0, len(envArtifact.DONs))

	jdConfig := jd.JDConfig{
		GRPC:  envArtifact.JdConfig.ExternalGRPCUrl,
		WSRPC: envArtifact.JdConfig.ExternalGRPCUrl,
		Creds: insecure.NewCredentials(),
	}

	offChain, offChainErr := jd.NewJDClient(jdConfig)
	if offChainErr != nil {
		return nil, nil, errors.Wrap(offChainErr, "failed to create offchain client")
	}

	topology, tErr := cre.NewTopology(cachedInput.NodeSets, *cachedInput.Infra)
	if tErr != nil {
		return nil, nil, errors.Wrap(tErr, "failed to recreate topology from artifact")
	}

	for idx, don := range envArtifact.DONs {
		_, ok := envArtifact.Nodes[don.DonName]
		if !ok {
			return nil, nil, errors.Errorf("no nodes found for don %s", don.DonName)
		}

		for id := range envArtifact.Nodes[don.DonName].Nodes {
			allNodeIDs = append(allNodeIDs, id)
		}

		startedDON, donErr := cre.NewDON(ctx, topology.DonsMetadata.List()[idx], cachedInput.NodeSets[idx].Out.CLNodes)
		if donErr != nil {
			return nil, nil, errors.Wrapf(donErr, "failed to create DON for don %s", don.DonName)
		}
		donsSlice = append(donsSlice, startedDON)
	}

	cldfBlockchains := make([]cldf_chain.BlockChain, 0, len(deployedBlockchains.Outputs))
	for _, db := range deployedBlockchains.Outputs {
		chain, chainErr := db.ToCldfChain()
		if chainErr != nil {
			return nil, nil, errors.Wrap(chainErr, "failed to create cldf chain from blockchain")
		}
		cldfBlockchains = append(cldfBlockchains, chain)
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
		cldf_chain.NewBlockChainsFromSlice(cldfBlockchains),
	)

	dons := cre.NewDons(donsSlice, envArtifact.GatewayConnectors)
	linkDonsToJDInput := &cre.LinkDonsToJDInput{
		Blockchains:     deployedBlockchains.Outputs,
		CldfEnvironment: cldEnv,
		Topology:        topology,
		Dons:            dons,
	}
	linkErr := cre.LinkToJobDistributor(ctx, linkDonsToJDInput)
	if linkErr != nil {
		return nil, nil, errors.Wrap(linkErr, "failed to link dons to JD")
	}

	return &cre.Environment{
		CldfEnvironment:       cldEnv,
		Blockchains:           deployedBlockchains.Outputs,
		RegistryChainSelector: deployedBlockchains.Outputs[0].ChainSelector(),
		Provider:              *cachedInput.Infra,
		CapabilityConfigs:     envArtifact.CapabilityConfigs,
		ContractVersions:      envArtifact.ContractVersions,
	}, dons, nil
}

func SetDefaultPrivateKeyIfEmpty(defaultPrivateKey string) error {
	if os.Getenv("PRIVATE_KEY") == "" {
		setErr := os.Setenv("PRIVATE_KEY", defaultPrivateKey)
		if setErr != nil {
			return fmt.Errorf("failed to set PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}

func SetDefaultSolanaPrivateKeyIfEmpty(key solana.PrivateKey) error {
	if os.Getenv("SOLANA_PRIVATE_KEY") == "" {
		setErr := os.Setenv("SOLANA_PRIVATE_KEY", key.String())
		if setErr != nil {
			return fmt.Errorf("failed to set SOLANA_PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set SOLANA_PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}

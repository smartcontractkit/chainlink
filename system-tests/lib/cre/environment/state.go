package environment

import (
	"context"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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
func BuildFromSavedState(ctx context.Context, cldLogger logger.Logger, cachedInput *envconfig.Config, envArtifact *EnvArtifact) (*cre.Environment, []*cre.WrappedBlockchainOutput, error) {
	if cachedInput == nil {
		return nil, nil, errors.New("cached input cannot be nil")
	}

	if envArtifact == nil {
		return nil, nil, errors.New("environment artifact cannot be nil")
	}

	if pkErr := SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return nil, nil, pkErr
	}
	// just in case
	if pkErr := SetDefaultSolanaPrivateKeyIfEmpty(defaultSolanaPrivateKey); pkErr != nil {
		return nil, nil, pkErr
	}

	wrappedBlockchainOutputs := make([]*cre.WrappedBlockchainOutput, 0)

	for _, bc := range cachedInput.Blockchains {
		if bc.Type == blockchain.FamilySolana {
			initErr := initSolanaInput(&bc)
			if initErr != nil {
				return nil, nil, errors.Wrap(initErr, "failed to init solana")
			}
			w, err := wrapSolana(&bc, bc.Out)
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to wrap solana")
			}
			wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, w)
			continue
		}

		if bc.Type == blockchain.FamilyTron {
			w, err := wrapTron(&bc, bc.Out)
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to wrap tron")
			}
			wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, w)
			continue
		}

		w, err := wrapEVM(bc.Out)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to wrap evm")
		}
		wrappedBlockchainOutputs = append(wrappedBlockchainOutputs, w)
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
	dons := make([]*cre.DON, 0, len(envArtifact.DONs))

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
		dons = append(dons, startedDON)
	}

	chainConfigs := make([]cre.ChainConfig, 0, len(wrappedBlockchainOutputs))
	for _, output := range wrappedBlockchainOutputs {
		cfg, cfgErr := cre.ChainConfigFromWrapped(output)
		if cfgErr != nil {
			return nil, nil, errors.Wrapf(cfgErr, "failed to build chain config from write for blockchain %s", output.BlockchainOutput.Family)
		}
		chainConfigs = append(chainConfigs, cfg)
	}

	blockChains, chainErr := cre.NewChains(cldLogger, chainConfigs)
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

	linkDonsToJDInput := &cre.LinkDonsToJDInput{
		JDClient:          offChain,
		BlockchainOutputs: wrappedBlockchainOutputs,
		CldfEnvironment:   cldEnv,
		Topology:          topology,
		DONs:              dons,
	}
	var linkErr error
	cldEnv, linkErr = cre.LinkToJobDistributor(ctx, linkDonsToJDInput)
	if linkErr != nil {
		return nil, nil, errors.Wrap(linkErr, "failed to link dons to JD")
	}

	return &cre.Environment{
		CldfEnvironment: cldEnv,
		DonTopology:     cre.NewDonTopology(envArtifact.RegistryChainSelector, topology, cre.NewDons(dons)),
	}, wrappedBlockchainOutputs, nil
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

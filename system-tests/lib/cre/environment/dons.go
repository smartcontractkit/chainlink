package environment

import (
	"context"
	"fmt"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type StartedDON struct {
	NodeOutput *cre.WrappedNodeOutput
	DON        *cre.Don
}

type StartedDONs []*StartedDON

func (s *StartedDONs) NodeOutputs() []*cre.WrappedNodeOutput {
	outputs := make([]*cre.WrappedNodeOutput, len(*s))
	for idx, don := range *s {
		outputs[idx] = don.NodeOutput
	}
	return outputs
}

func (s *StartedDONs) DONs() []*cre.Don {
	dons := make([]*cre.Don, len(*s))
	for idx, don := range *s {
		dons[idx] = don.DON
	}
	return dons
}

func StartDONs(
	ctx context.Context,
	lggr zerolog.Logger,
	topology *cre.Topology,
	infraInput infra.Provider,
	registryChainBlockchainOutput *blockchain.Output,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
	nodeSets []*cre.NodeSet,
) (*StartedDONs, error) {
	if infraInput.Type == infra.CRIB {
		deployCribDonsInput := &crib.DeployCribDonsInput{
			Topology:       topology,
			NodeSet:        nodeSets,
			CribConfigsDir: infra.CribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		var devspaceErr error
		nodeSets, devspaceErr = crib.DeployDons(ctx, deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with crib-sdk")
		}
	}

	for donIdx, donMetadata := range topology.DonsMetadata.List() {
		if !copyCapabilityBinaries {
			continue
		}

		customBinariesPaths := make(map[cre.CapabilityFlag]string)
		for flag, config := range capabilityConfigs {
			if flags.HasFlagForAnyChain(donMetadata.Flags, flag) && config.BinaryPath != "" {
				customBinariesPaths[flag] = config.BinaryPath
			}
		}

		executableErr := crecapabilities.MakeBinariesExecutable(customBinariesPaths)
		if executableErr != nil {
			return nil, pkgerrors.Wrap(executableErr, "failed to make binaries executable")
		}

		var err error
		ns, err := crecapabilities.AppendBinariesPathsNodeSpec(nodeSets[donIdx], donMetadata, customBinariesPaths)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
		}
		nodeSets[donIdx] = ns
	}

	// Add env vars, which were provided programmatically, to the node specs
	// or fail, if node specs already had some env vars set in the TOML config
	for donIdx, donMetadata := range topology.DonsMetadata.List() {
		hasEnvVarsInTomlConfig := false
		for nodeIdx, nodeSpec := range nodeSets[donIdx].NodeSpecs {
			if len(nodeSpec.Node.EnvVars) > 0 {
				hasEnvVarsInTomlConfig = true
				break
			}

			nodeSets[donIdx].NodeSpecs[nodeIdx].Node.EnvVars = nodeSets[donIdx].EnvVars
		}

		if hasEnvVarsInTomlConfig && len(nodeSets[donIdx].EnvVars) > 0 {
			return nil, fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
		}
	}

	errGroup, _ := errgroup.WithContext(ctx)
	var resultMap sync.Map

	for idx, nodeSet := range nodeSets {
		errGroup.Go(func() error {
			startTime := time.Now()
			lggr.Info().Msgf("Starting DON named %s", nodeSet.Name)
			nodeSet.Input.NodeSpecs = nodeSet.ExtractCTFInputs()
			nodeset, nodesetErr := ns.NewSharedDBNodeSetWithContext(ctx, nodeSet.Input, registryChainBlockchainOutput)
			if nodesetErr != nil {
				return pkgerrors.Wrapf(nodesetErr, "failed to start nodeSet named %s", nodeSet.Name)
			}

			don, donErr := cre.NewDON(ctx, topology.DonsMetadata.List()[idx], nodeset.CLNodes)
			if donErr != nil {
				return pkgerrors.Wrapf(donErr, "failed to create DON from node set named %s", nodeSet.Name)
			}

			resultMap.Store(idx, &StartedDON{
				NodeOutput: &cre.WrappedNodeOutput{
					Output:       nodeset,
					NodeSetName:  nodeSet.Name,
					Capabilities: nodeSet.ComputedCapabilities,
				},
				DON: don,
			})

			lggr.Info().Msgf("DON %s started in %.2f seconds", nodeSet.Name, time.Since(startTime).Seconds())

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		infra.PrintFailedContainerLogs(lggr, 30)
		return nil, err
	}

	startedDONs := make(StartedDONs, len(nodeSets))
	resultMap.Range(func(key, value any) bool {
		// key is index in the original slice
		startedDONs[key.(int)] = value.(*StartedDON)
		return true
	})

	return &startedDONs, nil
}

func FundNodes(ctx context.Context, testLogger zerolog.Logger, dons *cre.Dons, blockchains []blockchains.Blockchain, fundingAmountPerChainFamily map[string]uint64) error {
	for _, don := range dons.List() {
		testLogger.Info().Msgf("Funding nodes for DON %s", don.Name)
		for _, bc := range blockchains {
			if !flags.RequiresForwarderContract(don.Flags, bc.ChainID()) && !bc.IsFamily(chainselectors.FamilySolana) { // for now, we can only write to solana, so we consider forwarder is always present
				continue
			}

			chainFamily := bc.CtfOutput().Family
			fundingAmount, ok := fundingAmountPerChainFamily[chainFamily]
			if !ok {
				return fmt.Errorf("missing funding amount for chain family %s", chainFamily)
			}

			for _, node := range don.Nodes {
				address, addrErr := nodeAddress(node, chainFamily, bc)
				if addrErr != nil {
					return pkgerrors.Wrapf(addrErr, "failed to get address for node %s on chain family %s and chain %d", node.Name, chainFamily, bc.ChainID())
				}

				if address == "" {
					testLogger.Info().Msgf("No key for chainID %d found for node %s. Skipping funding", bc.ChainID(), node.Name)
					continue // Skip nodes without keys for this chain
				}

				err := bc.Fund(ctx, address, fundingAmount)
				if err != nil {
					return err
				}
			}
		}

		testLogger.Info().Msgf("Funded nodes for DON %s", don.Name)
	}

	return nil
}

func nodeAddress(node *cre.Node, chainFamily string, bc blockchains.Blockchain) (string, error) {
	switch chainFamily {
	case chainselectors.FamilyEVM, chainselectors.FamilyTron:
		evmKey, ok := node.Keys.EVM[bc.ChainID()]
		if !ok {
			return "", nil // Skip nodes without EVM keys for this chain
		}

		return evmKey.PublicAddress.String(), nil
	case chainselectors.FamilySolana:
		solBc := bc.(*solana.Blockchain)
		solKey, ok := node.Keys.Solana[solBc.SolanaChainID]
		if !ok {
			return "", nil // Skip nodes without Solana keys for this chain
		}
		return solKey.PublicAddress.String(), nil
	default:
		return "", fmt.Errorf("unsupported chain family %s", chainFamily)
	}
}

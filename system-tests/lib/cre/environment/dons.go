package environment

import (
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
)

func StartDONs(lggr zerolog.Logger, nixShell *nix.Shell, topology *cre.Topology, infraInput infra.Input, registryChainBlockchainOutput *blockchain.Output, capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet) ([]*cre.WrappedNodeOutput, error) {
	startTime := time.Now()
	lggr.Info().Msgf("Starting %d DONs", len(capabilitiesAwareNodeSets))

	if infraInput.Type == infra.CRIB {
		lggr.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &cre.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  capabilitiesAwareNodeSets,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		var devspaceErr error
		capabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with crib-sdk")
		}
	}

	nodeSetOutput := make([]*cre.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))

	// TODO we could parallelize this as well in the future, but for single DON env this doesn't matter
	for _, nodeSetInput := range capabilitiesAwareNodeSets {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, registryChainBlockchainOutput)
		if nodesetErr != nil {
			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}

		nodeSetOutput = append(nodeSetOutput, &cre.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.ComputedCapabilities,
		})
	}

	lggr.Info().Msgf("DONs started in %.2f seconds", time.Since(startTime).Seconds())

	return nodeSetOutput, nil
}

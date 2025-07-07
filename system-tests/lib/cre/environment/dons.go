package environment

import (
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func StartDONs(
	lggr zerolog.Logger,
	nixShell *nix.Shell,
	topology *cretypes.Topology,
	infraType libtypes.InfraType,
	registryChainBlockchainOutput *blockchain.Output,
	capabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet,
) ([]*cretypes.WrappedNodeOutput, error) {
	startTime := time.Now()
	lggr.Info().Msgf("Starting %d DONs", len(capabilitiesAwareNodeSets))

	if infraType == libtypes.CRIB {
		lggr.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &cretypes.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  capabilitiesAwareNodeSets,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var devspaceErr error
		capabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with devspace")
		}
	}

	nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))

	// TODO we could parallelize this as well in the future, but for single DON env this doesn't matter
	for _, nodeSetInput := range capabilitiesAwareNodeSets {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, registryChainBlockchainOutput)
		if nodesetErr != nil {
			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}

		nodeSetOutput = append(nodeSetOutput, &cretypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	lggr.Info().Msgf("DONs started in %.2f seconds", time.Since(startTime).Seconds())

	return nodeSetOutput, nil
}

package environment

import (
	"context"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type StartedDON struct {
	NodeOutput *cre.WrappedNodeOutput
	DON        *cre.DON
}

type StartedDONs []*StartedDON

func (s *StartedDONs) NodeOutputs() []*cre.WrappedNodeOutput {
	outputs := make([]*cre.WrappedNodeOutput, 0, len(*s))
	for idx, don := range *s {
		outputs[idx] = don.NodeOutput
	}
	return outputs
}

func (s *StartedDONs) DONs() []*cre.DON {
	dons := make([]*cre.DON, 0, len(*s))
	for idx, don := range *s {
		dons[idx] = don.DON
	}
	return dons
}

func StartDONs(ctx context.Context, lggr zerolog.Logger, topology *cre.Topology, infraInput infra.Provider, registryChainBlockchainOutput *blockchain.Output, capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet) (*StartedDONs, error) {
	startTime := time.Now()
	lggr.Info().Msgf("Starting %d DONs", len(capabilitiesAwareNodeSets))

	if infraInput.Type == infra.CRIB {
		lggr.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &cre.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  capabilitiesAwareNodeSets,
			CribConfigsDir: cribConfigsDir,
			Namespace:      infraInput.CRIB.Namespace,
		}

		var devspaceErr error
		capabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with crib-sdk")
		}
	}

	startedDONs := StartedDONs{}
	// TODO we could parallelize this as well in the future, but for single DON env this doesn't matter
	for idx, nodeSetInput := range capabilitiesAwareNodeSets {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, registryChainBlockchainOutput)
		if nodesetErr != nil {
			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}

		don, donErr := cre.NewDON(ctx, topology.DonsMetadata.List()[idx], nodeset.CLNodes)
		if donErr != nil {
			return nil, pkgerrors.Wrapf(donErr, "failed to create DON from node set named %s", nodeSetInput.Name)
		}

		startedDONs = append(startedDONs, &StartedDON{
			NodeOutput: &cre.WrappedNodeOutput{
				Output:       nodeset,
				NodeSetName:  nodeSetInput.Name,
				Capabilities: nodeSetInput.ComputedCapabilities,
			},
			DON: don,
		})
	}

	lggr.Info().Msgf("DONs started in %.2f seconds", time.Since(startTime).Seconds())

	return &startedDONs, nil
}

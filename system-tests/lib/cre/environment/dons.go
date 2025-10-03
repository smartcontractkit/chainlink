package environment

import (
	"context"
	"slices"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func StartDONs(ctx context.Context, lggr zerolog.Logger, topology *cre.Topology, infraInput infra.Provider, registryChainBlockchainOutput *blockchain.Output, capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet) ([]*cre.WrappedNodeOutput, error) {
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

	errGroup, _ := errgroup.WithContext(ctx)

	// to preserve order of node sets, we will collect results in a channel and then sort them by index
	// because later on we rely on the order of node sets matching the order of DONs in topology
	type orderedNodeOutput struct {
		Index  int
		Output *cre.WrappedNodeOutput
	}

	resultCh := make(chan *orderedNodeOutput, len(capabilitiesAwareNodeSets))
	for idx, nodeSetInput := range capabilitiesAwareNodeSets {
		startTime := time.Now()
		lggr.Info().Msgf("Starting DON named %s", nodeSetInput.Name)
		errGroup.Go(func() error {
			nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, registryChainBlockchainOutput)
			if nodesetErr != nil {
				return pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
			}

			resultCh <- &orderedNodeOutput{
				Output: &cre.WrappedNodeOutput{
					Output:       nodeset,
					NodeSetName:  nodeSetInput.Name,
					Capabilities: nodeSetInput.ComputedCapabilities,
				},
				Index: idx,
			}

			lggr.Info().Msgf("DON %s started in %.2f seconds", nodeSetInput.Name, time.Since(startTime).Seconds())

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}
	close(resultCh)

	orderedOutput := make([]*orderedNodeOutput, len(capabilitiesAwareNodeSets))
	for res := range resultCh {
		orderedOutput[res.Index] = res
	}

	slices.SortFunc(orderedOutput, func(a, b *orderedNodeOutput) int {
		return a.Index - b.Index
	})

	nodeSetOutput := make([]*cre.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))
	for _, res := range orderedOutput {
		nodeSetOutput = append(nodeSetOutput, res.Output)
	}

	return nodeSetOutput, nil
}

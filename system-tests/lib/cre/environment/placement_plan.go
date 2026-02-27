package environment

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

type placementPlan struct {
	NodeSetPlacement    *nodeSetPlacementSummary
	HasRemoteComponents bool
}

type nodeSetPlacementSummary struct {
	HasLocalTargets  bool
	HasRemoteTargets bool
}

func buildPlacementPlan(
	configuredBlockchains []*config.Blockchain,
	jdInput *config.JobDistributor,
	nodeSets []*cre.NodeSet,
) (*placementPlan, error) {
	nodeSetPlacement, err := summarizeNodeSetPlacement(nodeSets)
	if err != nil {
		return nil, err
	}
	if err := validateUnsupportedPlacements(configuredBlockchains, nodeSetPlacement); err != nil {
		return nil, err
	}

	return &placementPlan{
		NodeSetPlacement:    nodeSetPlacement,
		HasRemoteComponents: hasRemoteComponents(configuredBlockchains, jdInput, nodeSets),
	}, nil
}

func hasRemoteComponents(blockchains []*config.Blockchain, jdInput *config.JobDistributor, nodeSets []*cre.NodeSet) bool {
	for _, configuredBlockchain := range blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Placement == config.PlacementRemote {
			return true
		}
	}
	if jdInput != nil && jdInput.Placement == config.PlacementRemote {
		return true
	}
	for _, nodeSet := range nodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Placement) == string(config.PlacementRemote) {
			return true
		}
	}
	return false
}

func summarizeNodeSetPlacement(nodeSets []*cre.NodeSet) (*nodeSetPlacementSummary, error) {
	summary := &nodeSetPlacementSummary{}
	for _, nodeSet := range nodeSets {
		if nodeSet == nil {
			continue
		}
		configPlacement := strings.TrimSpace(nodeSet.Placement)
		if configPlacement == "" || configPlacement == string(config.PlacementLocal) {
			summary.HasLocalTargets = true
			continue
		}
		if configPlacement == string(config.PlacementRemote) {
			summary.HasRemoteTargets = true
			continue
		}
		return nil, fmt.Errorf("invalid nodeset placement: %s", nodeSet.Placement)
	}

	return summary, nil
}

func validateUnsupportedPlacements(
	configuredBlockchains []*config.Blockchain,
	nodeSetPlacement *nodeSetPlacementSummary,
) error {
	if nodeSetPlacement == nil || !nodeSetPlacement.HasRemoteTargets {
		return nil
	}
	for _, bc := range configuredBlockchains {
		if bc == nil {
			continue
		}
		if bc.Placement == config.PlacementLocal {
			return fmt.Errorf(
				"remote nodesets with local blockchains are not supported in this PoC. " +
					"Set all blockchains to placement=remote, or run nodesets with placement=local so nodes stay colocated with local blockchains",
			)
		}
	}
	return nil
}

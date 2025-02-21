package don

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func globalBootstraperNodeData(donTopologies []*types.DonWithMetadata) (string, string, error) {
	var findHost = func(n devenv.Node) string {
		for _, label := range n.Labels() {
			if label.Key == node.HostLabelKey {
				return *label.Value
			}
		}
		return ""
	}

	if len(donTopologies) == 1 {
		bootstrapNode, err := node.FindOneWithLabel(donTopologies[0].DON, &ptypes.Label{Key: node.RoleLabelKey, Value: ptr.Ptr(types.BootstrapNode)})
		if err != nil {
			return "", "", errors.Wrap(err, "failed to find bootstrap node")
		}

		// if there is only one DON, then the global bootstrapper is the bootstrap node of the DON
		peerID, err := node.ToP2PID(*bootstrapNode, node.KeyExtractingTransformFn)
		if err != nil {
			return "", "", errors.Wrapf(err, "failed to get peer ID for node %s", donTopologies[0].DON.Nodes[0].Name)
		}

		bootstrapNodeHost := findHost(*bootstrapNode)
		if bootstrapNodeHost == "" {
			return "", "", errors.New("failed to get bootstrap node host from labels")
		}

		return peerID, bootstrapNodeHost, nil
	} else if len(donTopologies) > 1 {
		// if there's more than one DON, then peering capabilitity needs to point to the same bootstrap node
		// for all the DONs, and so we need to find it first. For us, it will always be the bootstrap node of the workflow DON.
		for _, donTopology := range donTopologies {
			if flags.HasFlag(donTopology.Flags, types.WorkflowDON) {
				bootstrapNode, err := node.FindOneWithLabel(donTopology.DON, &ptypes.Label{Key: node.RoleLabelKey, Value: ptr.Ptr(types.BootstrapNode)})
				if err != nil {
					return "", "", errors.Wrap(err, "failed to find bootstrap node")
				}

				peerID, err := node.ToP2PID(*bootstrapNode, node.KeyExtractingTransformFn)
				if err != nil {
					return "", "", errors.Wrapf(err, "failed to get peer ID for node %s", bootstrapNode.Name)
				}

				bootstrapNodeHost := findHost(*bootstrapNode)
				if bootstrapNodeHost == "" {
					return "", "", errors.New("failed to get bootstrap node host from labels")
				}

				return peerID, bootstrapNodeHost, nil
			}
		}

		return "", "", errors.New("expected at least one workflow DON")
	}

	return "", "", errors.New("expected at least one DON topology")
}

func FindPeeringData(donTopologies []*types.DonWithMetadata) (types.PeeringData, error) {
	globalBootstraperPeerID, globalBootstraperHost, err := globalBootstraperNodeData(donTopologies)
	if err != nil {
		return types.PeeringData{}, err
	}

	return types.PeeringData{
		GlobalBootstraperPeerID: globalBootstraperPeerID,
		GlobalBootstraperHost:   globalBootstraperHost,
	}, nil
}

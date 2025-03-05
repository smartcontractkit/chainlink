package don

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func globalBootstraperNodeData(topology *types.Topology) (string, string, error) {
	if len(topology.DonsMetadata) == 1 {
		bootstrapNode, err := node.FindOneWithLabel(topology.DonsMetadata[0].NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.BootstrapNode}, node.EqualLabels)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to find bootstrap node")
		}

		// if there is only one DON, then the global bootstrapper is the bootstrap node of the DON
		peerID, err := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to get peer ID for the bootstrap node")
		}

		bootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
		if hostErr != nil {
			return "", "", errors.Wrap(hostErr, "failed to get bootstrap node host from labels")
		}

		return peerID, bootstrapNodeHost, nil
	} else if len(topology.DonsMetadata) > 1 {
		// if there's more than one DON, then peering capabilitity needs to point to the same bootstrap node
		// for all the DONs, and so we need to find it first. For us, it will always be the bootstrap node of the workflow DON.
		for _, donTopology := range topology.DonsMetadata {
			if flags.HasFlag(donTopology.Flags, types.WorkflowDON) {
				bootstrapNode, err := node.FindOneWithLabel(donTopology.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.BootstrapNode}, node.EqualLabels)
				if err != nil {
					return "", "", errors.Wrap(err, "failed to find bootstrap node")
				}

				peerID, err := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
				if err != nil {
					return "", "", errors.Wrapf(err, "failed to get peer ID for workernode %s", "CHANGE ME")
				}

				bootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
				if hostErr != nil {
					return "", "", errors.Wrap(hostErr, "failed to get bootstrap node host from labels")
				}

				return peerID, bootstrapNodeHost, nil
			}
		}

		return "", "", errors.New("expected at least one workflow DON")
	}

	return "", "", errors.New("expected at least one DON topology")
}

func FindPeeringData(donTopologies *types.Topology) (types.CapabilitiesPeeringData, error) {
	globalBootstraperPeerID, globalBootstraperHost, err := globalBootstraperNodeData(donTopologies)
	if err != nil {
		return types.CapabilitiesPeeringData{}, err
	}

	return types.CapabilitiesPeeringData{
		GlobalBootstraperPeerID: globalBootstraperPeerID,
		GlobalBootstraperHost:   globalBootstraperHost,
	}, nil
}

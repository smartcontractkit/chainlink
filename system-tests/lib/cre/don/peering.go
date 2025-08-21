package don

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func globalBootstraperNodeData(topology *cre.Topology) (string, string, error) {
	if len(topology.DonsMetadata) == 0 {
		return "", "", errors.New("expected at least one DON topology")
	}

	// if there's more than one DON, then peering capabilitity needs to point to the same bootstrap node
	// for all the DONs, and so we need to find it first. For us, it will always be the bootstrap node of the workflow DON.
	for _, donTopology := range topology.DonsMetadata {
		if flags.HasFlag(donTopology.Flags, cre.WorkflowDON) {
			bootstrapNode, err := node.FindOneWithLabel(donTopology.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
			if err != nil {
				return "", "", errors.Wrap(err, "failed to find bootstrap node")
			}

			peerID, err := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
			if err != nil {
				return "", "", errors.Wrap(err, "failed to get bootstrap node's peerID from labels")
			}

			bootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
			if hostErr != nil {
				return "", "", errors.Wrap(hostErr, "failed to get bootstrap node's host from labels")
			}

			return peerID, bootstrapNodeHost, nil
		}
	}

	return "", "", errors.New("expected at least one workflow DON")
}

func FindPeeringData(donTopologies *cre.Topology) (cre.CapabilitiesPeeringData, cre.OCRPeeringData, error) {
	globalBootstraperPeerID, globalBootstraperHost, err := globalBootstraperNodeData(donTopologies)
	if err != nil {
		return cre.CapabilitiesPeeringData{}, cre.OCRPeeringData{}, err
	}

	return cre.CapabilitiesPeeringData{
			GlobalBootstraperPeerID: globalBootstraperPeerID,
			GlobalBootstraperHost:   globalBootstraperHost,
			Port:                    CapabilitiesPeeringPort,
		}, cre.OCRPeeringData{
			OCRBootstraperPeerID: globalBootstraperPeerID,
			OCRBootstraperHost:   globalBootstraperHost,
			Port:                 OCRPeeringPort,
		}, nil
}

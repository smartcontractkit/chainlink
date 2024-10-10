package keystone

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type updateNodesRequest struct {
	p2pToCapabilities map[string][]registeredCapability
	nodes             []*ocr2Node
}

func UpdateNodes(lggr logger.Logger, req *registerNodesRequest) (*registerNodesResponse, error) {
	req.registry.UpdateNodes(req.chain.DeployerKey, nil)
	return nil, nil
}

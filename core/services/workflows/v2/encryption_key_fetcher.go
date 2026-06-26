package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"sort"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

type encryptionKeyFetcher struct {
	lggr        logger.Logger
	capRegistry core.CapabilitiesRegistry
}

var _ host.EncryptionKeyFetcher = (*encryptionKeyFetcher)(nil)

func (e *encryptionKeyFetcher) GetEncryptionKeys(ctx context.Context) ([]string, error) {
	e.lggr.Debug("Fetching encryption keys...")
	myNode, err := e.capRegistry.LocalNode(ctx)
	if err != nil {
		return nil, errors.New("failed to get local node from registry" + err.Error())
	}

	encryptionKeys := make([]string, 0, len(myNode.WorkflowDON.Members))
	for _, peerID := range myNode.WorkflowDON.Members {
		peerNode, err := e.capRegistry.NodeByPeerID(ctx, peerID)
		if err != nil {
			return nil, errors.New("failed to get node info for peerID: " + peerID.String() + " - " + err.Error())
		}
		encryptionKeys = append(encryptionKeys, hex.EncodeToString(peerNode.EncryptionPublicKey[:]))
	}
	// Sort the encryption keys to ensure consistent ordering across all nodes.
	sort.Strings(encryptionKeys)
	return encryptionKeys, nil
}

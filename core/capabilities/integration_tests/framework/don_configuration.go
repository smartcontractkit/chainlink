package framework

import (
	"fmt"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type DonConfiguration struct {
	commoncap.DON
	name       string
	keys       []ethkey.KeyV2
	KeyBundles []ocr2key.KeyBundle
	peerIDs    []peer
}

// NewDonConfigurationParams exists purely to make it obvious in the test code what DON configuration is being used
type NewDonConfigurationParams struct {
	Name             string
	NumNodes         int
	F                uint8
	AcceptsWorkflows bool
}

func NewDonConfiguration(don NewDonConfigurationParams) (DonConfiguration, error) {
	if !(don.NumNodes >= int(3*don.F+1)) {
		return DonConfiguration{}, fmt.Errorf("invalid configuration, number of nodes must be at least 3*F+1")
	}

	keyBundles, peerIDs, err := getKeyBundlesAndPeerIDs(don.Name, don.NumNodes)
	if err != nil {
		return DonConfiguration{}, fmt.Errorf("failed to get key bundles and peer IDs: %w", err)
	}

	donPeers := make([]p2ptypes.PeerID, len(peerIDs))
	var donKeys []ethkey.KeyV2
	for i := 0; i < len(peerIDs); i++ {
		peerID := p2ptypes.PeerID{}
		err = peerID.UnmarshalText([]byte(peerIDs[i].PeerID))
		if err != nil {
			return DonConfiguration{}, fmt.Errorf("failed to unmarshal peer ID: %w", err)
		}
		donPeers[i] = peerID
		newKey, err := ethkey.NewV2()
		if err != nil {
			return DonConfiguration{}, fmt.Errorf("failed to create key: %w", err)
		}
		donKeys = append(donKeys, newKey)
	}

	donConfiguration := DonConfiguration{
		DON: commoncap.DON{
			Members:          donPeers,
			F:                don.F,
			ConfigVersion:    1,
			AcceptsWorkflows: don.AcceptsWorkflows,
		},
		name:       don.Name,
		peerIDs:    peerIDs,
		keys:       donKeys,
		KeyBundles: keyBundles,
	}
	return donConfiguration, nil
}

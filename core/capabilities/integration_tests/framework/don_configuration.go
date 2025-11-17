package framework

import (
	"errors"
	"fmt"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type DonConfiguration struct {
	commoncap.DON
	name       string
	keys       []ethkey.KeyV2
	KeyBundles []ocr2key.KeyBundle
	p2pKeys    []p2pkey.KeyV2
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
		return DonConfiguration{}, errors.New("invalid configuration, number of nodes must be at least 3*F+1")
	}

	keyBundles, p2pKeys, err := getKeyBundlesAndP2PKeys(don.Name, don.NumNodes)
	if err != nil {
		return DonConfiguration{}, fmt.Errorf("failed to get key bundles and p2p keys: %w", err)
	}

	donPeers := make([]p2ptypes.PeerID, len(p2pKeys))
	var donKeys []ethkey.KeyV2
	for i := range p2pKeys {
		donPeers[i] = p2ptypes.PeerID(p2pKeys[i].PeerID())
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
		p2pKeys:    p2pKeys,
		keys:       donKeys,
		KeyBundles: keyBundles,
	}
	return donConfiguration, nil
}

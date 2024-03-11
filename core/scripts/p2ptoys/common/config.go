package common

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/smartcontractkit/libocr/commontypes"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type Config struct {
	Nodes         []string `json:"nodes"`
	Bootstrappers []string `json:"bootstrappers"`

	// parsed values below
	NodeKeys       []ed25519.PrivateKey
	NodePeerIDs    []ragep2ptypes.PeerID
	NodePeerIDsStr []string

	BootstrapperKeys      []ed25519.PrivateKey
	BootstrapperPeerInfos []ragep2ptypes.PeerInfo
	BootstrapperLocators  []commontypes.BootstrapperLocator
}

const (
	// bootstrappers will listen on 127.0.0.1 ports 9000, 9001, 9002, etc.
	BootstrapStartPort = 9000

	// nodes will listen on 127.0.0.1 ports 8000, 8001, 8002, etc.
	NodeStartPort = 8000
)

func ParseConfigFromFile(fileName string) (*Config, error) {
	rawConfig, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, err
	}

	for _, hexKey := range config.Nodes {
		key, peerID, err := parseKey(hexKey)
		if err != nil {
			return nil, err
		}
		config.NodeKeys = append(config.NodeKeys, key)
		config.NodePeerIDs = append(config.NodePeerIDs, peerID)
		config.NodePeerIDsStr = append(config.NodePeerIDsStr, peerID.String())
	}

	for _, hexKey := range config.Bootstrappers {
		key, peerID, err := parseKey(hexKey)
		if err != nil {
			return nil, err
		}
		config.BootstrapperKeys = append(config.BootstrapperKeys, key)
		config.BootstrapperPeerInfos = append(config.BootstrapperPeerInfos, ragep2ptypes.PeerInfo{ID: peerID})
	}

	locators := []commontypes.BootstrapperLocator{}
	for id, peerID := range config.BootstrapperPeerInfos {
		addr := fmt.Sprintf("127.0.0.1:%d", BootstrapStartPort+id)
		locators = append(locators, commontypes.BootstrapperLocator{
			PeerID: peerID.ID.String(),
			Addrs:  []string{addr},
		})
		config.BootstrapperPeerInfos[id].Addrs = []ragep2ptypes.Address{ragep2ptypes.Address(addr)}
	}
	config.BootstrapperLocators = locators

	return &config, nil
}

func parseKey(hexKey string) (ed25519.PrivateKey, ragep2ptypes.PeerID, error) {
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, ragep2ptypes.PeerID{}, err
	}
	key := ed25519.PrivateKey(b)
	peerID, err := ragep2ptypes.PeerIDFromPrivateKey(key)
	if err != nil {
		return nil, ragep2ptypes.PeerID{}, err
	}
	return key, peerID, nil
}

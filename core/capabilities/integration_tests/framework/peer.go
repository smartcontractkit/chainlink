package framework

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type peerIDAndOCRSigner struct {
	PeerID p2ptypes.PeerID
	Signer string
}

func peerToNode(nopID uint32, p peerIDAndOCRSigner) (kcr.CapabilitiesRegistryNodeParams, error) {
	sig := strings.TrimPrefix(p.Signer, "0x")
	signerB, err := hex.DecodeString(sig)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return kcr.CapabilitiesRegistryNodeParams{
		NodeOperatorId:      nopID,
		P2pId:               p.PeerID,
		Signer:              sigb,
		EncryptionPublicKey: testutils.Random32Byte(),
	}, nil
}

func getKeyBundlesAndP2PKeys(donName string, numNodes int) ([]ocr2key.KeyBundle, []p2pkey.KeyV2, error) {
	var keyBundles []ocr2key.KeyBundle
	var donPeerKeys []p2pkey.KeyV2
	for range numNodes {
		key, err := p2pkey.NewV2()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create p2p key: %w", err)
		}

		keyBundle, err := ocr2key.New(chaintype.EVM)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create key bundle: %w", err)
		}

		keyBundles = append(keyBundles, keyBundle)

		donPeerKeys = append(donPeerKeys, key)
	}
	return keyBundles, donPeerKeys, nil
}

type peerWrapper struct {
	peer p2pPeer
}

func (t peerWrapper) Start(ctx context.Context) error {
	return nil
}

func (t peerWrapper) Close() error {
	return nil
}

func (t peerWrapper) Ready() error {
	return nil
}

func (t peerWrapper) HealthReport() map[string]error {
	return nil
}

func (t peerWrapper) Name() string {
	return "peerWrapper"
}

func (t peerWrapper) GetPeer() p2ptypes.Peer {
	return t.peer
}

type p2pPeer struct {
	id p2ptypes.PeerID
}

func (t p2pPeer) Start(ctx context.Context) error {
	return nil
}

func (t p2pPeer) Close() error {
	return nil
}

func (t p2pPeer) Ready() error {
	return nil
}

func (t p2pPeer) HealthReport() map[string]error {
	return nil
}

func (t p2pPeer) Name() string {
	return "p2pPeer"
}

func (t p2pPeer) ID() p2ptypes.PeerID {
	return t.id
}

func (t p2pPeer) UpdateConnections(peers map[p2ptypes.PeerID]p2ptypes.StreamConfig) error {
	return nil
}

func (t p2pPeer) Send(peerID p2ptypes.PeerID, msg []byte) error {
	return nil
}

func (t p2pPeer) Receive() <-chan p2ptypes.Message {
	return nil
}

func (t p2pPeer) IsBootstrap() bool {
	return false
}

func getSignerStringFromOCRKeyBundle(keyBundle ocr2key.KeyBundle) (string, error) {
	if keyBundle == nil {
		return "", errors.New("keyBundle is nil")
	}

	return fmt.Sprintf("0x%x", keyBundle.PublicKey()), nil
}

func ptr[T any](t T) *T { return &t }

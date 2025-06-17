package framework

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type peer struct {
	PeerID string
	Signer string
}

func peerIDToBytes(peerID string) ([32]byte, error) {
	var peerIDB ragetypes.PeerID
	err := peerIDB.UnmarshalText([]byte(peerID))
	if err != nil {
		return [32]byte{}, err
	}

	return peerIDB, nil
}

func peers(ps []peer) ([][32]byte, error) {
	out := [][32]byte{}
	for _, p := range ps {
		b, err := peerIDToBytes(p.PeerID)
		if err != nil {
			return nil, err
		}

		out = append(out, b)
	}

	return out, nil
}

func peerToNode(nopID uint32, p peer) (kcr.CapabilitiesRegistryNodeParams, error) {
	peerIDB, err := peerIDToBytes(p.PeerID)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert peerID: %w", err)
	}

	sig := strings.TrimPrefix(p.Signer, "0x")
	signerB, err := hex.DecodeString(sig)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return kcr.CapabilitiesRegistryNodeParams{
		NodeOperatorId:      nopID,
		P2pId:               peerIDB,
		Signer:              sigb,
		EncryptionPublicKey: testutils.Random32Byte(),
	}, nil
}

func getKeyBundlesAndPeerIDs(donName string, numNodes int) ([]ocr2key.KeyBundle, []peer, error) {
	var keyBundles []ocr2key.KeyBundle
	var donPeerIDs []peer
	for i := 0; i < numNodes; i++ {
		peerID := NewPeerID(donName, i)

		keyBundle, err := ocr2key.New(chaintype.EVM)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create key bundle: %w", err)
		}

		keyBundles = append(keyBundles, keyBundle)

		pk := keyBundle.PublicKey()

		p := peer{
			PeerID: peerID,
			Signer: fmt.Sprintf("0x%x", pk),
		}

		donPeerIDs = append(donPeerIDs, p)
	}
	return keyBundles, donPeerIDs, nil
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

func NewPeerID(donName string, nodeOrdinal int) string {
	privKeyString := fmt.Sprintf("privatekey:%s:%d", donName, nodeOrdinal)
	privKey := sha256.Sum256([]byte(privKeyString))
	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID)
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}

func ptr[T any](t T) *T { return &t }

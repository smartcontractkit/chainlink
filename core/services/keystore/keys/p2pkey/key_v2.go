package p2pkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

type Raw []byte

func (raw Raw) Key() KeyV2 {
	privKey, err := cryptop2p.UnmarshalPrivateKey(raw)
	if err != nil {
		panic(err)
	}
	key, err := fromPrivkey(privKey)
	if err != nil {
		panic(err)
	}
	return key
}

func (raw Raw) String() string {
	return "<P2P Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &KeyV2{}

type KeyV2 struct {
	cryptop2p.PrivKey
	peerID PeerID
}

func NewV2() (KeyV2, error) {
	privKey, _, err := cryptop2p.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	return fromPrivkey(privKey)
}

func MustNewV2XXXTestingOnly(k *big.Int) KeyV2 {
	seed := make([]byte, ed25519.SeedSize)
	copy(seed, k.Bytes())
	pk := ed25519.NewKeyFromSeed(seed[:])
	p2pPrivKey, err := cryptop2p.UnmarshalEd25519PrivateKey(pk[:])
	if err != nil {
		panic(err)
	}
	key, err := fromPrivkey(p2pPrivKey)
	if err != nil {
		panic(err)
	}
	return key
}

func (key KeyV2) ID() string {
	return peer.ID(key.peerID).String()
}

func (key KeyV2) Raw() Raw {
	marshalledPrivK, err := cryptop2p.MarshalPrivateKey(key.PrivKey)
	if err != nil {
		panic(err)
	}
	return marshalledPrivK
}

func (key KeyV2) PeerID() PeerID {
	return key.peerID
}

func (key KeyV2) PublicKeyHex() string {
	pubKeyBytes, err := key.GetPublic().Raw()
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(pubKeyBytes)
}

func (key KeyV2) String() string {
	return fmt.Sprintf("P2PKeyV2{PrivateKey: <redacted>, PeerID: %s}", key.peerID.Raw())
}

func (key KeyV2) GoString() string {
	return key.String()
}

func fromPrivkey(privKey cryptop2p.PrivKey) (KeyV2, error) {
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		return KeyV2{}, err
	}
	return KeyV2{
		PrivKey: privKey,
		peerID:  PeerID(peerID),
	}, nil
}

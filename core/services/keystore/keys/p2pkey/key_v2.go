package p2pkey

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

var libp2pPBPrefix = []byte{0x08, 0x01, 0x12, 0x40}

func KeyFor(raw internal.Raw) KeyV2 {
	privKey, err := unmarshalPrivateKey(raw)
	if err != nil {
		panic(err)
	}
	key, err := fromPrivkey(privKey)
	if err != nil {
		panic(err)
	}
	return key
}

func unmarshalPrivateKey(raw internal.Raw) (ed25519.PrivateKey, error) {
	if !bytes.HasPrefix(raw.Bytes(), libp2pPBPrefix) {
		return nil, errors.New("invalid key: missing libp2p protobuf prefix")
	}
	return raw.Bytes()[len(libp2pPBPrefix):], nil
}

func marshalPrivateKey(key ed25519.PrivateKey) ([]byte, error) {
	return bytes.Join([][]byte{libp2pPBPrefix, key}, nil), nil
}

var _ fmt.GoStringer = &KeyV2{}

type KeyV2 struct {
	PrivKey ed25519.PrivateKey
	peerID  PeerID
}

func NewV2() (KeyV2, error) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	return fromPrivkey(privKey)
}

func MustNewV2XXXTestingOnly(k *big.Int) KeyV2 {
	seed := make([]byte, ed25519.SeedSize)
	copy(seed, k.Bytes())
	pk := ed25519.NewKeyFromSeed(seed)
	key, err := fromPrivkey(pk)
	if err != nil {
		panic(err)
	}
	return key
}

func (key KeyV2) ID() string {
	return types.PeerID(key.peerID).String()
}

func (key KeyV2) Raw() internal.Raw {
	marshalledPrivK, err := marshalPrivateKey(key.PrivKey)
	if err != nil {
		panic(err)
	}
	return internal.NewRaw(marshalledPrivK)
}

func (key KeyV2) PeerID() PeerID {
	return key.peerID
}

func (key KeyV2) PublicKeyHex() string {
	pubKeyBytes := key.PrivKey.Public().(ed25519.PublicKey)
	return hex.EncodeToString(pubKeyBytes)
}

func (key KeyV2) String() string {
	return fmt.Sprintf("P2PKeyV2{PrivateKey: <redacted>, PeerID: %s}", key.peerID.Raw())
}

func (key KeyV2) GoString() string {
	return key.String()
}

func fromPrivkey(privKey ed25519.PrivateKey) (KeyV2, error) {
	peerID, err := types.PeerIDFromPrivateKey(privKey)
	if err != nil {
		return KeyV2{}, err
	}
	return KeyV2{
		PrivKey: privKey,
		peerID:  PeerID(peerID),
	}, nil
}

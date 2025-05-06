package p2pkey

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
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
	key.raw = raw
	return key
}

func unmarshalPrivateKey(raw internal.Raw) (ed25519.PrivateKey, error) {
	b := internal.Bytes((raw))
	if !bytes.HasPrefix(b, libp2pPBPrefix) {
		return nil, errors.New("invalid key: missing libp2p protobuf prefix")
	}
	return b[len(libp2pPBPrefix):], nil
}

func marshalPrivateKey(key ed25519.PrivateKey) ([]byte, error) {
	return bytes.Join([][]byte{libp2pPBPrefix, key}, nil), nil
}

type KeyV2 struct {
	raw    internal.Raw
	signer crypto.Signer

	peerID PeerID
}

func fromPrivkey(privKey ed25519.PrivateKey) (KeyV2, error) {
	peerID, err := types.PeerIDFromPrivateKey(privKey)
	if err != nil {
		return KeyV2{}, err
	}
	return KeyV2{
		signer: privKey,
		peerID: PeerID(peerID),
	}, nil
}

func NewV2() (KeyV2, error) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	key, err := fromPrivkey(privKey)
	if err != nil {
		return KeyV2{}, err
	}
	marshalledPrivK, err := marshalPrivateKey(privKey)
	if err != nil {
		return KeyV2{}, err
	}
	key.raw = internal.NewRaw(marshalledPrivK)
	return key, nil
}

func MustNewV2XXXTestingOnly(k *big.Int) KeyV2 {
	seed := make([]byte, ed25519.SeedSize)
	copy(seed, k.Bytes())
	pk := ed25519.NewKeyFromSeed(seed)
	key, err := fromPrivkey(pk)
	if err != nil {
		panic(err)
	}
	marshalledPrivK, err := marshalPrivateKey(pk)
	if err != nil {
		panic(err)
	}
	key.raw = internal.NewRaw(marshalledPrivK)
	return key
}

func (key KeyV2) ID() string {
	return types.PeerID(key.peerID).String()
}

func (key KeyV2) Raw() internal.Raw {
	return key.raw
}

func (key KeyV2) PeerID() PeerID {
	return key.peerID
}

func (key KeyV2) PublicKeyHex() string {
	pubKeyBytes := key.signer.Public().(ed25519.PublicKey)
	return hex.EncodeToString(pubKeyBytes)
}

func (key KeyV2) Public() crypto.PublicKey { return key.signer.Public() }

func (key KeyV2) Sign(rand io.Reader, message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return key.signer.Sign(rand, message, opts)
}

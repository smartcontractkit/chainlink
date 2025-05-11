package tonkey

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

// Key represents a TON ed25519 key
type Key struct {
	raw    internal.Raw
	signFn func(io.Reader, []byte, crypto.SignerOpts) ([]byte, error)
	pubKey ed25519.PublicKey
}

// KeyFor loads a Key from a raw seed
func KeyFor(raw internal.Raw) Key {
	seed := internal.Bytes(raw)
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)

	return Key{
		raw:    raw,
		signFn: privKey.Sign,
		pubKey: pubKey,
	}
}

// New creates a new Key using secure randomness
func New() (Key, error) {
	reader := crypto_rand.Reader
	return newFrom(reader)
}

// newFrom creates a new Key from a reader
func newFrom(reader io.Reader) (Key, error) {
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return Key{}, err
	}
	rawSeed := priv.Seed()
	return Key{
		raw:    internal.NewRaw(rawSeed),
		signFn: priv.Sign,
		pubKey: pub,
	}, nil
}

// MustNewInsecure creates a Key or panics if an error occurs
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

// ID returns the hex-encoded public key (same as PublicKeyStr)
func (key Key) ID() string {
	return key.PublicKeyStr()
}

// PublicKeyStr returns the hex-encoded public key
func (key Key) PublicKeyStr() string {
	return hex.EncodeToString(key.pubKey)
}

// GetPublic returns the raw ed25519 public key bytes
func (key Key) GetPublic() ed25519.PublicKey {
	return key.pubKey
}

// Raw returns the private key seed
func (key Key) Raw() internal.Raw {
	return key.raw
}

// Sign signs a message using ed25519
func (key Key) Sign(msg []byte) ([]byte, error) {
	rng := crypto_rand.Reader
	hash := crypto.Hash(0)
	return key.signFn(rng, msg, hash)
}

// PubkeyToAddressWithVersionAndWorkchain returns the TON wallet address for the given wallet version and workchain
func (key Key) PubkeyToAddressWithVersionAndWorkchain(version wallet.Version, workchain int8) (*address.Address, error) {
	privKey := ed25519.NewKeyFromSeed(internal.Bytes(key.raw))
	w, err := wallet.FromPrivateKeyWithOptions(nil, privKey, version, wallet.WithWorkchain(workchain))
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}
	return w.WalletAddress(), nil
}

// PubkeyToAddress returns the wallet v3 address for workchain 0
func (key Key) PubkeyToAddress() *address.Address {
	ver := wallet.V3
	addr, err := key.PubkeyToAddressWithVersionAndWorkchain(ver, 0)
	if err != nil {
		panic(fmt.Errorf("failed to get address: %w", err))
	}
	return addr
}

// UserFriendlyAddress returns the user-friendly version of the Ton address
// https://docs.ton.org/v3/concepts/dive-into-ton/ton-blockchain/smart-contract-addresses#user-friendly-address
func (key Key) UserFriendlyAddress() string {
	address := key.PubkeyToAddress()
	return address.String()
}

// RawAddress returns the raw version of the Ton address
// https://docs.ton.org/v3/concepts/dive-into-ton/ton-blockchain/smart-contract-addresses#raw-address
func (key Key) RawAddress() string {
	address := key.PubkeyToAddress()
	return address.StringRaw()
}

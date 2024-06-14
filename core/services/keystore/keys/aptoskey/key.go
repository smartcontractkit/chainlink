package aptoskey

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"fmt"
	"io"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/sha3"
)

// AccountAddress is a 32 byte address on the Aptos blockchain
// It can represent an Object, an Account, and much more.
// Extracting this out from the aptos sdk as there are still breaking changes
// https://github.com/aptos-labs/aptos-go-sdk
type AccountAddress [32]byte

// Raw represents the ETH private key
type Raw []byte

// Key gets the Key
func (raw Raw) Key() Key {
	privKey := ed25519.NewKeyFromSeed(raw)
	pubKey := privKey.Public().(ed25519.PublicKey)
	accountAddress := PubkeyToAddress(pubKey)
	return Key{
		Address: accountAddress,
		privkey: privKey,
		pubKey:  pubKey,
	}
}

// String returns description
func (raw Raw) String() string {
	return "<Aptos Raw Private Key>"
}

// GoString wraps String()
func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

// Key represents Aptos key
type Key struct {
	Address AccountAddress
	privkey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

// New creates new Key
func New() (Key, error) {
	return newFrom(crypto_rand.Reader)
}

func PubkeyToAddress(pubkey ed25519.PublicKey) AccountAddress {
	authKey := sha3.Sum256(append([]byte(pubkey), 0x00))
	accountAddress := AccountAddress(authKey)
	return accountAddress
}

// MustNewInsecure return Key if no error
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

// newFrom creates new Key from a provided random reader
func newFrom(reader io.Reader) (Key, error) {
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return Key{}, err
	}
	accountAddress := PubkeyToAddress(pub)
	return Key{
		Address: accountAddress,
		privkey: priv,
		pubKey:  pub,
	}, nil
}

// ID gets Key ID
func (key Key) ID() string {
	return key.PublicKeyStr()
}

// GetPublic get Key's public key
func (key Key) GetPublic() ed25519.PublicKey {
	return key.pubKey
}

// PublicKeyStr return base58 encoded public key
func (key Key) PublicKeyStr() string {
	return base58.Encode(key.pubKey)
}

// Raw returns the seed from private key
func (key Key) Raw() Raw {
	return key.privkey.Seed()
}

// String is the print-friendly format of the Key
func (key Key) String() string {
	return fmt.Sprintf("AptosKey{PrivateKey: <redacted>, Address: %s}", key.Address)
}

// GoString wraps String()
func (key Key) GoString() string {
	return key.String()
}

// Sign is used to sign a message
func (key Key) Sign(msg []byte) ([]byte, error) {
	return key.privkey.Sign(crypto_rand.Reader, msg, crypto.Hash(0)) // no specific hash function used
}

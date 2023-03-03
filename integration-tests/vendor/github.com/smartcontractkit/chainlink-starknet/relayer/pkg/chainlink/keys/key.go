package keys

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/dontpanicdao/caigo"
)

// Raw represents the Stark private key
type Raw []byte

// Key gets the Key
func (raw Raw) Key() Key {
	k := Key{}
	var err error

	k.priv = new(big.Int).SetBytes(raw)
	k.pub.X, k.pub.Y, err = caigo.Curve.PrivateToPoint(k.priv)
	if err != nil {
		panic(err) // key not generated
	}
	return k
}

// String returns description
func (raw Raw) String() string {
	return "<StarkNet Raw Private Key>"
}

// GoString wraps String()
func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

type PublicKey struct {
	X, Y *big.Int
}

// Key represents StarkNet key
type Key struct {
	priv *big.Int
	pub  PublicKey
}

// New creates new Key
func New() (Key, error) {
	return newFrom(crypto_rand.Reader)
}

// MustNewInsecure return Key if no error
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

func newFrom(reader io.Reader) (Key, error) {
	return GenerateKey(reader)
}

// ID gets Key ID
func (key Key) ID() string {
	return key.AccountAddressStr()
}

// this is the derived contract address, the contract is deployed using the StarkKeyStr
// This is the primary identifier for onchain interactions
// the private key is identified by this
func (key Key) AccountAddressStr() string {
	return "0x" + hex.EncodeToString(PubKeyToAccount(key.pub, defaultContractHash, defaultSalt))
}

// StarkKeyStr is the starknet public key associated to the private key
// it is the X component of the ECDSA pubkey and used in the deployment of the account contract
// this func is used in exporting it via CLI and API
func (key Key) StarkKeyStr() string {
	return "0x" + hex.EncodeToString(PubKeyToStarkKey(key.pub))
}

// Raw from private key
func (key Key) Raw() Raw {
	return key.priv.Bytes()
}

// String is the print-friendly format of the Key
func (key Key) String() string {
	return fmt.Sprintf("StarkNetKey{PrivateKey: <redacted>, Contract Address: %s}", key.AccountAddressStr())
}

// GoString wraps String()
func (key Key) GoString() string {
	return key.String()
}

// ToPrivKey returns the key usable for signing.
func (key Key) ToPrivKey() *big.Int {
	return key.priv
}

// PublicKey copies public key object
func (key Key) PublicKey() PublicKey {
	return key.pub
}

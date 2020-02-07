// package vrfkey tracks the secret keys associated with VRF proofs. It
// is a separate package from ../store to increase encapsulation of the keys,
// and reduce the risk of them leaking.
//
// The three types, PrivateKey, PublicKey and EncryptedSecretKey are all aspects
// of the one keypair.
package vrfkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"go.dedis.ch/kyber/v3"
)

// PrivateKey represents the secret used to construct a VRF proof.
//
// Don't serialize directly, use Encrypt method, with user-supplied passphrase.
// The unencrypted PrivateKey struct should only live in-memory.
//
// Only use it if you absolutely need it (i.e., for a novel crypto protocol.)
// Implement whatever cryptography you need on this struct, so your callers
// don't need to know the secret key explicitly. (See, e.g., MarshaledProof.)
type PrivateKey struct {
	k         kyber.Scalar
	PublicKey PublicKey
}

// PublicKey is a secp256k1 point in compressed format
type PublicKey [CompressedPublicKeyLength]byte

// EncryptedSecretKey contains encrypted private key to be serialized to DB
//
// We could re-use geth's key handling, here, but this makes it much harder to
// misuse VRF proving keys as ethereum keys or vice versa.
type EncryptedSecretKey struct {
	PublicKey PublicKey     `gorm:"primary_key;type:varchar(68)"`
	VRFKey    gethKeyStruct `json:"vrf_key" gorm:"type:text"`
}

// Copied from go-ethereum/accounts/keystore/key.go's encryptedKeyJSONV3
type gethKeyStruct struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

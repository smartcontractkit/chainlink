package signature

import (
	"bytes"
	"crypto/ecdsa"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

// Curve is the elliptic Curve on which on-chain messages are to be signed
var Curve = secp256k1.S256()

// OnChainPublicKey is the public key used to cryptographically identify an
// oracle to the on-chain smart contract.
type OnChainPublicKey ecdsa.PublicKey

// Equal returns true iff k and k2 represent the same public key
func (k OnChainPublicKey) Equal(k2 OnChainPublicKey) bool {
	return bytes.Equal(
		common.Address(k.Address()).Bytes(),
		common.Address(k2.Address()).Bytes(),
	)
}

type EthAddresses = map[types.OnChainSigningAddress]commontypes.OracleID

// VerifyOnChain returns an error unless signature is a valid signature by one
// of the signers, in which case it returns the ID of the signer
func VerifyOnChain(msg []byte, signature []byte, signers EthAddresses,
) (commontypes.OracleID, error) {
	author, err := crypto.SigToPub(onChainHash(msg), signature)
	if err != nil {
		return commontypes.OracleID(math.MaxUint8), errors.Wrapf(err, "while trying to recover "+
			"sender from sig %x on msg %+v", signature, msg)
	}
	oid, ok := signers[(*OnChainPublicKey)(author).Address()]
	if ok {
		return oid, nil
	} else {
		return commontypes.OracleID(math.MaxUint8), errors.Errorf("signer is not on whitelist")
	}
}

// OnchainPrivateKey is the secret key oracles use to sign messages destined for
// verification by the on-chain smart contract.
type OnchainPrivateKey ecdsa.PrivateKey

// Sign signs message with k
func (k *OnchainPrivateKey) Sign(msg []byte) (signature []byte, err error) {
	sig, err := crypto.Sign(onChainHash(msg), (*ecdsa.PrivateKey)(k))
	return sig, err
}

func onChainHash(msg []byte) []byte {
	return crypto.Keccak256(msg)
}

func (k OnChainPublicKey) Address() types.OnChainSigningAddress {
	return types.OnChainSigningAddress(crypto.PubkeyToAddress(ecdsa.PublicKey(k)))
}

func (k OnchainPrivateKey) Address() types.OnChainSigningAddress {
	return types.OnChainSigningAddress(crypto.PubkeyToAddress(k.PublicKey))
}

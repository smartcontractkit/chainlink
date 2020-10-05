package signature

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

var Curve = secp256k1.S256()

type OnChainPublicKey ecdsa.PublicKey

var zero = big.NewInt(0)

func (k OnChainPublicKey) fieldEqual(a, b *big.Int) bool {
	difference := big.NewInt(0).Sub(a, b)
	modulus := ecdsa.PublicKey(k).Params().P
	return big.NewInt(1).Mod(difference, modulus).Cmp(zero) == 0
}

func (k OnChainPublicKey) Equal(k2 OnChainPublicKey) bool {
	return bytes.Equal(
		common.Address(k.Address()).Bytes(),
		common.Address(k2.Address()).Bytes(),
	)
}

var halfGroupOrder = big.NewInt(1).Rsh(Curve.Params().N, 1)

type EthAddresses = map[types.OnChainSigningAddress]types.OracleID

func VerifyOnChain(msg []byte, signature []byte, signers EthAddresses,
) (types.OracleID, error) {
	author, err := crypto.SigToPub(onChainHash(msg), signature)
	if err != nil {
		return types.OracleID(-1), errors.Wrapf(err, "while trying to recover "+
			"sender from sig %x on msg %+v", signature, msg)
	}
	oid, ok := signers[(*OnChainPublicKey)(author).Address()]
	if ok {
		return oid, nil
	} else {
		return types.OracleID(-1), errors.Errorf("signer is not on whitelist")
	}
}

type OnchainPrivateKey ecdsa.PrivateKey

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

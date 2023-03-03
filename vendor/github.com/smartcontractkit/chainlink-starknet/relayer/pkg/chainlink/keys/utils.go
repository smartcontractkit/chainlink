package keys

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/NethermindEth/juno/pkg/crypto/pedersen"
	"github.com/dontpanicdao/caigo"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

// constants
var (
	byteLen = 32

	// note: the contract hash must match the corresponding OZ gauntlet command hash - otherwise addresses will not correspond
	defaultContractHash, _ = new(big.Int).SetString("0x0750cd490a7cd1572411169eaa8be292325990d33c5d4733655fe6b926985062", 0)
	defaultSalt            = big.NewInt(100)
)

// PubKeyToContract implements the pubkey to deployed account given contract hash + salt
func PubKeyToAccount(pubkey PublicKey, classHash, salt *big.Int) []byte {
	hash := pedersen.ArrayDigest(
		new(big.Int).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS")),
		big.NewInt(0),
		salt,      // salt
		classHash, // classHash
		pedersen.ArrayDigest(pubkey.X),
	)

	// pad big.Int to 32 bytes if needed
	return starknet.PadBytes(hash.Bytes(), byteLen)
}

// PubToStarkKey implements the pubkey to starkkey functionality: https://github.com/0xs34n/starknet.js/blob/cd61356974d355aa42f07a3d63f7ccefecbd913c/src/utils/ellipticCurve.ts#L49
func PubKeyToStarkKey(pubkey PublicKey) []byte {
	return starknet.PadBytes(pubkey.X.Bytes(), byteLen)
}

// reimplements parts of https://github.com/dontpanicdao/caigo/blob/main/utils.go#L85
// generate the PK as a pseudo-random number in the interval [1, CurveOrder - 1]
// using io.Reader, and Key struct
func GenerateKey(material io.Reader) (k Key, err error) {
	max := new(big.Int).Sub(caigo.Curve.N, big.NewInt(1))

	k.priv, err = rand.Int(material, max)
	if err != nil {
		return k, err
	}

	k.pub.X, k.pub.Y, err = caigo.Curve.PrivateToPoint(k.priv)
	if err != nil {
		return k, err
	}

	if !caigo.Curve.IsOnCurve(k.pub.X, k.pub.Y) {
		return k, fmt.Errorf("key gen is not on stark curve")
	}

	return k, nil
}

package ethkey

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

var curve = crypto.S256()

func KeyFor(raw internal.Raw) KeyV2 {
	var privateKey ecdsa.PrivateKey
	d := big.NewInt(0).SetBytes(internal.Bytes(raw))
	privateKey.PublicKey.Curve = curve
	privateKey.D = d
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	k := newKeyV2(&privateKey)
	k.raw = raw
	return k
}

type KeyV2 struct {
	raw          internal.Raw
	getPK        func() *ecdsa.PrivateKey
	Address      common.Address
	EIP55Address types.EIP55Address
}

func newKeyV2(privKey *ecdsa.PrivateKey) KeyV2 {
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	eip55 := types.EIP55AddressFromAddress(address)
	return KeyV2{
		getPK:        func() *ecdsa.PrivateKey { return privKey },
		Address:      address,
		EIP55Address: eip55,
	}
}

func NewV2() (KeyV2, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	return FromPrivateKey(privateKeyECDSA), nil
}

func FromPrivateKey(privKey *ecdsa.PrivateKey) (key KeyV2) {
	key = newKeyV2(privKey)
	key.raw = internal.NewRaw(privKey.D.Bytes())
	return
}

func (key KeyV2) ID() string {
	return key.Address.Hex()
}

func (key KeyV2) Raw() internal.Raw { return key.raw }

func (key KeyV2) Sign(data []byte) ([]byte, error) { return crypto.Sign(data, key.getPK()) }

// Cmp uses byte-order address comparison to give a stable comparison between two keys
func (key KeyV2) Cmp(key2 KeyV2) int {
	return bytes.Compare(key.Address.Bytes(), key2.Address.Bytes())
}

func (key KeyV2) SignerFn(chainID *big.Int) bind.SignerFn {
	return func(from common.Address, tx *gethtypes.Transaction) (*gethtypes.Transaction, error) {
		signer := gethtypes.LatestSignerForChainID(chainID)
		h := signer.Hash(tx)
		sig, err := key.Sign(h[:])
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}
		return tx.WithSignature(signer, sig)
	}
}

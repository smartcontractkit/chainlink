package metatx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	forwarder_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	BankERC20TokenName    = "BankERC20"
	BankERC20TokenSymbol  = "BANK"
	BankERC20TokenVersion = "v1"
)

var TypeHash = crypto.Keccak256([]byte("ForwardRequest(address from,address target,uint256 nonce,bytes data,uint256 validUntilTime)"))

func SignMetaTransfer(
	forwarder forwarder_wrapper.Forwarder,
	ownerPrivateKey *ecdsa.PrivateKey,
	owner, sourceTokenAddress common.Address,
	calldataHash [32]byte,
	deadline *big.Int,
	tokenName, tokenVersion string,
) (signature []byte, domainSeparatorHash [32]byte, typeHash [32]byte, nonce *big.Int, err error) {
	genericParams, err := forwarder.GENERICPARAMS(nil)
	if err != nil {
		return nil, [32]byte{}, [32]byte{}, nil, errors.Wrapf(err, "failed to get domainType of forwarder: %x", forwarder.Address())
	}

	typeHashRaw := crypto.Keccak256([]byte(fmt.Sprintf("ForwardRequest(%s)", genericParams)))

	// Sanity check: make sure TypeHash is equal to keccak256 of the constants in Forwarder.sol
	if !bytes.Equal(TypeHash, typeHashRaw) {
		return nil, [32]byte{}, [32]byte{}, nil, errors.Errorf("unexpected domainType hash. Expected: %x Actual: %x", TypeHash, typeHashRaw)
	}

	copy(typeHash[:], typeHashRaw[:])

	nonce, err = rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil))
	if err != nil {
		panic(err)
	}

	domainSeparator, err := forwarder.GetDomainSeparator(nil, tokenName, tokenVersion)
	if err != nil {
		return nil, [32]byte{}, [32]byte{}, nil, errors.Wrap(err, "failed to get domain separator from contract")
	}
	domainSeparatorHashRaw := crypto.Keccak256(domainSeparator)
	copy(domainSeparatorHash[:], domainSeparatorHashRaw[:])

	copy(typeHash[:], TypeHash[:])
	message := []byte{0x19, 0x01} // \x19\x01
	message = append(message, domainSeparatorHashRaw[:]...)

	encodedCall, err := utils.ABIEncode(
		`
	[
			{"name": "typeHash","type":"bytes32"},
			{"name": "from","type":"address"},
			{"name": "target", "type": "address"},
			{"name": "nonce", "type": "uint256"},
			{"name": "data", "type": "bytes32"},
			{"name": "validUntilTime", "type": "uint256"}
	]
	`, typeHash, owner, sourceTokenAddress, nonce, calldataHash, deadline,
	)

	if err != nil {
		return nil, [32]byte{}, [32]byte{}, nil, errors.Wrap(err, "failed to abi encode")
	}

	encodedHash := crypto.Keccak256(encodedCall)

	message = append(message, encodedHash...)
	messageDigest := crypto.Keccak256(message)
	rawSignature, err := crypto.Sign(messageDigest, ownerPrivateKey)
	if err != nil {
		return nil, [32]byte{}, [32]byte{}, nil, errors.Wrap(err, "failed to sign message digest")
	}

	// decompose signature into v, r and s
	// the returned byte array is [R || S || V]
	if len(rawSignature) != 65 {
		panic("rawSignature should be 65 bytes long")
	}
	var (
		v uint8
		r [32]byte
		s [32]byte
	)
	rSlice := rawSignature[:32] // first 32 bytes is R
	copy(r[:], rSlice[:])
	sSlice := rawSignature[32:64] // second 32 bytes is S
	copy(s[:], sSlice[:])
	v = uint8(rawSignature[64])
	if v == 1 {
		v = 28
	} else {
		v = 27
	}

	signature = append(signature, r[:]...)
	signature = append(signature, s[:]...)
	signature = append(signature, v)

	return
}

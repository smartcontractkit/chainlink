package metatx

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/meta_erc20"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func SignMetaTransfer(
	metaToken meta_erc20.MetaERC20Interface,
	ownerPrivateKey *ecdsa.PrivateKey,
	owner, to common.Address,
	amount *big.Int,
	deadline *big.Int,
) (v uint8, r [32]byte, s [32]byte, err error) {
	nonce, err := metaToken.Nonces(nil, owner)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, errors.Wrapf(err, "failed to get nonce of %s", owner.Hex())
	}

	domainSeparator, err := metaToken.DOMAINSEPARATOR(nil)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, errors.Wrap(err, "failed to get domain separator from contract")
	}

	typeHash, err := metaToken.METATRANSFERTYPEHASH(nil)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, errors.Wrap(err, "failed to get type hash from contract")
	}

	message := []byte{0x19, 0x01} // \x19\x01
	message = append(message, domainSeparator[:]...)
	encodedCall, err := utils.ABIEncode(
		`
[
		{"name": "metaTransferTypeHash","type":"bytes32"},
		{"name": "owner","type":"address"},
		{"name": "to", "type": "address"},
		{"name": "amount", "type": "uint256"},
		{"name": "nonce", "type": "uint256"},
		{"name": "deadline", "type": "uint256"}
]
`, typeHash, owner, to, amount, nonce, deadline,
	)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, errors.Wrap(err, "failed to abi encode")
	}

	encodedHash := crypto.Keccak256(encodedCall)

	message = append(message, encodedHash...)
	messageDigest := crypto.Keccak256(message)
	signature, err := crypto.Sign(messageDigest, ownerPrivateKey)
	if err != nil {
		return 0, [32]byte{}, [32]byte{}, errors.Wrap(err, "failed to sign message digest")
	}

	// decompose signature into v, r and s
	// the returned byte array is [R || S || V]
	if len(signature) != 65 {
		panic("signature should be 65 bytes long")
	}
	rSlice := signature[:32] // first 32 bytes is R
	copy(r[:], rSlice[:])
	sSlice := signature[32:64] // second 32 bytes is S
	copy(s[:], sSlice[:])
	v = uint8(signature[64])
	if v == 1 {
		v = 28
	} else {
		v = 27
	}

	return
}

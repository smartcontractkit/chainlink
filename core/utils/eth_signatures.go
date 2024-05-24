package utils

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

const EthSignedMessagePrefix = "\x19Ethereum Signed Message:\n"

func GetSignersEthAddress(msg []byte, sig []byte) (recoveredAddr common.Address, err error) {
	if len(sig) != 65 {
		return recoveredAddr, errors.New("invalid signature: signature length must be 65 bytes")
	}

	// Adjust the V component of the signature in case it uses 27 or 28 instead of 0 or 1
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}
	if sig[64] != 0 && sig[64] != 1 {
		return recoveredAddr, errors.New("invalid signature: invalid V component")
	}

	prefixedMsg := fmt.Sprintf("%s%d%s", EthSignedMessagePrefix, len(msg), msg)
	hash := crypto.Keccak256Hash([]byte(prefixedMsg))

	sigPublicKey, err := crypto.SigToPub(hash[:], sig)
	if err != nil {
		return recoveredAddr, err
	}

	recoveredAddr = crypto.PubkeyToAddress(*sigPublicKey)
	return recoveredAddr, nil
}

func GenerateEthPrefixedMsgHash(msg []byte) (hash common.Hash) {
	prefixedMsg := fmt.Sprintf("%s%d%s", EthSignedMessagePrefix, len(msg), msg)
	return crypto.Keccak256Hash([]byte(prefixedMsg))
}

func GenerateEthSignature(privateKey *ecdsa.PrivateKey, msg []byte) (signature []byte, err error) {
	hash := GenerateEthPrefixedMsgHash(msg)
	return crypto.Sign(hash[:], privateKey)
}

package common

import (
	"crypto/ecdsa"
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/slices"
)

func Uint32ToBytes(val uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, val)
	return result
}

func BytesToUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

// input string can't have any 0x0 characters
func StringToAlignedBytes(input string, size int) []byte {
	aligned := make([]byte, size)
	copy(aligned, input)
	return aligned
}

func AlignedBytesToString(data []byte) string {
	idx := slices.IndexFunc(data, func(b byte) bool { return b == 0 })
	return string(data[:idx])
}

func SignData(privateKey *ecdsa.PrivateKey, data ...[]byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data...)
	return crypto.Sign(hash.Bytes(), privateKey)
}

func ValidateSignature(signature []byte, data ...[]byte) (signerAddress []byte, err error) {
	hash := crypto.Keccak256Hash(data...)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return
	}
	ecdsaPubKey, _ := crypto.UnmarshalPubkey(sigPublicKey)
	signerAddress = crypto.PubkeyToAddress(*ecdsaPubKey).Bytes()

	signatureNoRecoverID := signature[:len(signature)-1]
	if !crypto.VerifySignature(sigPublicKey, hash.Bytes(), signatureNoRecoverID) {
		return nil, errors.New("invalid signature")
	}
	return
}

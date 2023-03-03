package contract

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/pkg/errors"

	dkg_hash "github.com/smartcontractkit/ocr2vrf/types/hash"

	"go.dedis.ch/kyber/v3"
)

type KeyData struct {
	PublicKey kyber.Point
	Hashes    []dkg_hash.Hash
}

func (kd *KeyData) MarshalBinary(keyID [32]byte) (rv []byte, err error) {
	km, err := kd.PublicKey.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not serialize key for onchain report")
	}
	return keyDataEncoding.Pack(keyID, km, kd.Hashes)
}

func (kd *KeyData) String() string {
	km, err := kd.PublicKey.MarshalBinary()
	if err != nil {
		panic(errors.Errorf("could not extract binary data for key %v", kd.PublicKey))
	}
	hashStrings := make([]string, len(kd.Hashes))
	for i, h := range kd.Hashes {
		hashStrings[i] = h.String()
	}
	return fmt.Sprintf("{Key: 0x%x, Hashes: [%s]}", km, strings.Join(hashStrings, ", "))
}

var keyDataEncoding = getKeyDataEncoding()

func getKeyDataEncoding() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return []abi.Argument{
		{Name: "keyID", Type: mustNewType("bytes32")},
		{Name: "publicKey", Type: mustNewType("bytes")},
		{Name: "hashes", Type: mustNewType("bytes32[]")},
	}
}

func MakeKeyDataFromOnchainKeyData(
	kd OnchainKeyData, kg kyber.Group,
) (KeyData, error) {
	pk := kg.Point()
	if err := pk.UnmarshalBinary(kd.PublicKey); err != nil {
		return KeyData{}, errors.Wrap(err, "could not unmarshal onchain key")
	}
	var hashes []dkg_hash.Hash
	for _, h := range kd.Hashes {
		hashes = append(hashes, h)
	}
	return KeyData{pk, hashes}, nil
}

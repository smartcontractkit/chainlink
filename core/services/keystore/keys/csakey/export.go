package csakey

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const keyTypeIdentifier = "CSA"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	data, err := corekeys.FromEncryptedCSAKey(keyJSON, password)
	switch {
	case errors.Is(err, corekeys.ErrInvalidExportFormat):
		return internal.FromEncryptedJSON(
			keyTypeIdentifier,
			keyJSON,
			password,
			adulteratedPassword,
			func(_ internal.EncryptedKeyExport, rawPrivKey internal.Raw) (KeyV2, error) {
				return KeyFor(rawPrivKey), nil
			},
		)
	case err != nil:
		return KeyV2{}, err
	}

	return KeyFor(internal.NewRaw(data)), nil
}

func (k KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return internal.ToEncryptedJSON(
		keyTypeIdentifier,
		k,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key KeyV2, cryptoJSON keystore.CryptoJSON) internal.EncryptedKeyExport {
			return internal.EncryptedKeyExport{
				KeyType:   id,
				PublicKey: key.PublicKeyString(),
				Crypto:    cryptoJSON,
			}
		},
	)
}

func adulteratedPassword(password string) string {
	return "csakey" + password
}

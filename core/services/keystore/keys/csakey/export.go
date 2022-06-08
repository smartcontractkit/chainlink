package csakey

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "CSA"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	return keys.FromEncryptedJSON(keyTypeIdentifier, keyJSON, password, adulteratedPassword, func(raw []byte) KeyV2 {
		return Raw(raw).Key()
	})
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(keyTypeIdentifier, key.Raw(), key.PublicKeyString(), password, scryptParams, adulteratedPassword)
}

func adulteratedPassword(password string) string {
	return "csakey" + password
}

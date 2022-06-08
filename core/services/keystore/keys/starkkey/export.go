package starkkey

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "StarkNet"

// FromEncryptedJSON gets key from json and password
func FromEncryptedJSON(keyJSON []byte, password string) (Key, error) {
	return keys.FromEncryptedJSON(keyTypeIdentifier, keyJSON, password, adulteratedPassword, func(raw []byte) Key {
		return Raw(raw).Key()
	})
}

// ToEncryptedJSON returns encrypted JSON representing key
func (key Key) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return keys.ToEncryptedJSON(keyTypeIdentifier, key.Raw(), key.PublicKeyStr(), password, scryptParams, adulteratedPassword)
}

func adulteratedPassword(password string) string {
	return "starkkey" + password
}

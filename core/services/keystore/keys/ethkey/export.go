package ethkey

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	var export EncryptedEthKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return KeyV2{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return KeyV2{}, errors.Wrap(err, "failed to decrypt Eth key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

type EncryptedEthKeyExport struct {
	KeyType string              `json:"keyType"`
	Address EIP55Address        `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	// DEV: uuid is derived directly from the address, since it is not stored internally
	id, err := uuid.FromBytes(key.Address.Bytes()[:16])
	if err != nil {
		return nil, errors.Wrapf(err, "could not generate ethkey UUID")
	}
	dKey := &keystore.Key{
		Id:         id,
		Address:    key.Address.Address(),
		PrivateKey: key.privateKey,
	}
	return keystore.EncryptKey(dKey, password, scryptParams.N, scryptParams.P)
}

func adulteratedPassword(password string) string {
	return "ethkey" + password
}

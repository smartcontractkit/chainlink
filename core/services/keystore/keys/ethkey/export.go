package ethkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

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

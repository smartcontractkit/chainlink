package keystore

import (
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

func BuildBeholderAuth(keyStore Master) (authHeaders map[string]string, pubKeyHex string, err error) {
	csaKeys, err := keyStore.CSA().GetAll()
	if err != nil {
		return nil, "", err
	}
	csaKey := csaKeys[0]
	csaPrivKey := csaKey.Raw().Bytes()
	authHeaders = beholder.BuildAuthHeaders(csaPrivKey)
	pubKeyHex = hex.EncodeToString(csaKey.PublicKey)
	return
}

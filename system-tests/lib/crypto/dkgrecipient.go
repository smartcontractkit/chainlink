package crypto

import (
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type DKGRecipientKeys struct {
	EncryptedJSONs [][]byte
	PubKeys        []dkgocrtypes.P256ParticipantPublicKey
	Password       string
}

func GenerateDKGRecipientKeys(password string, n int) (*DKGRecipientKeys, error) {
	result := &DKGRecipientKeys{
		Password: password,
	}
	for i := 0; i < n; i++ {
		key, err := dkgrecipientkey.New()
		if err != nil {
			return nil, err
		}
		d, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
		if err != nil {
			return nil, err
		}

		result.EncryptedJSONs = append(result.EncryptedJSONs, d)
		result.PubKeys = append(result.PubKeys, key.PublicKey())
	}
	return result, nil
}

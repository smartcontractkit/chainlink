package crypto

import (
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type DKGRecipientKey struct {
	EncryptedJSON []byte
	PubKey        dkgocrtypes.P256ParticipantPublicKey
	Password      string
}

func NewDKGRecipientKey(password string) (*DKGRecipientKey, error) {
	result := &DKGRecipientKey{
		Password: password,
	}
	key, err := dkgrecipientkey.New()
	if err != nil {
		return nil, err
	}
	d, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
	if err != nil {
		return nil, err
	}

	result.EncryptedJSON = d
	result.PubKey = key.PublicKey()

	return result, nil
}

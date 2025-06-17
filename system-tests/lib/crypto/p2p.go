package crypto

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func GenerateP2PKeys(password string, n int) (*types.P2PKeys, error) {
	result := &types.P2PKeys{
		Password: password,
	}
	for i := 0; i < n; i++ {
		key, err := p2pkey.NewV2()
		if err != nil {
			return nil, err
		}
		d, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
		if err != nil {
			return nil, err
		}

		result.EncryptedJSONs = append(result.EncryptedJSONs, d)
		result.PeerIDs = append(result.PeerIDs, key.PeerID().String())
	}
	return result, nil
}

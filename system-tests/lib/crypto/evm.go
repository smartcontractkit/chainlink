package crypto

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func GenerateEVMKeys(password string, n int) (*types.EVMKeys, error) {
	result := &types.EVMKeys{
		Password: password,
	}
	for range n {
		key, addr, err := clclient.NewETHKey(password)
		if err != nil {
			return result, nil
		}
		result.EncryptedJSONs = append(result.EncryptedJSONs, key)
		result.PublicAddresses = append(result.PublicAddresses, addr)
	}
	return result, nil
}

package crypto

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

type EVMKeys struct {
	EncryptedJSONs  [][]byte
	PublicAddresses []common.Address
	Password        string
	ChainID         int
}

func GenerateEVMKeys(password string, n int) (*EVMKeys, error) {
	result := &EVMKeys{
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

package crypto

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/pkg/errors"

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

/*
Generates new private and public key pair

Returns a new public address and a private key
*/
func GenerateNewKeyPair() (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, pkErr := crypto.GenerateKey()
	if pkErr != nil {
		return common.Address{}, nil, errors.Wrap(pkErr, "failed to generate a new private key (EOA)")
	}

	publicKeyAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	return publicKeyAddr, privateKey, nil
}

package crypto

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type SolKeys struct {
	EncryptedJSONs  [][]byte
	PublicAddresses []solana.PublicKey
	Password        string
	ChainID         string
}

func GenerateSolKeys(password string, n int) (*SolKeys, error) {
	result := &SolKeys{
		Password: password,
	}
	for range n {
		key, err := solkey.New()
		if err != nil {
			return nil, fmt.Errorf("err create solkey: %w", err)
		}

		enc, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt key: %w", err)
		}

		result.EncryptedJSONs = append(result.EncryptedJSONs, enc)
		result.PublicAddresses = append(result.PublicAddresses, key.PublicKey())
	}

	return result, nil
}

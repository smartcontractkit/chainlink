package crypto

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/solkey"
)

type SolKey struct {
	EncryptedJSON []byte
	PublicAddress solana.PublicKey
	Password      string
	ChainID       string
}

func NewSolKey(password, chainID string) (*SolKey, error) {
	key, err := solkey.New()
	if err != nil {
		return nil, fmt.Errorf("err create solkey: %w", err)
	}

	enc, err := key.ToEncryptedJSON(password, keystore.DefaultScryptParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	return &SolKey{
		EncryptedJSON: enc,
		PublicAddress: key.PublicKey(),
		Password:      password,
		ChainID:       chainID,
	}, nil
}

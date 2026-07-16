package crypto

import (
	"fmt"

	aptossdk "github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/aptoskey"
)

type AptosKey struct {
	EncryptedJSON []byte
	PublicKey     string
	Account       string
	Password      string
}

func NewAptosKey(password string) (*AptosKey, error) {
	key, err := aptoskey.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create aptos key: %w", err)
	}

	enc, err := key.ToEncryptedJSON(password, keystore.DefaultScryptParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt aptos key: %w", err)
	}

	account, err := NormalizeAptosAccount(key.Account())
	if err != nil {
		return nil, fmt.Errorf("failed to normalize aptos account: %w", err)
	}

	return &AptosKey{
		EncryptedJSON: enc,
		PublicKey:     key.PublicKeyStr(),
		Account:       account,
		Password:      password,
	}, nil
}

func NormalizeAptosAccount(raw string) (string, error) {
	var addr aptossdk.AccountAddress
	if err := addr.ParseStringRelaxed(raw); err != nil {
		return "", err
	}
	return addr.StringLong(), nil
}

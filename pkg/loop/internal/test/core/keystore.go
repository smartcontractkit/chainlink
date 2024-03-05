package core

import (
	"bytes"
	"context"
	"fmt"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var Keystore = staticKeystore{
	staticKeystoreConfig: staticKeystoreConfig{
		Account: libocr.Account("testaccount"),
		encoded: []byte{5: 11},
		signed:  []byte{13: 37},
	},
}

var _ types.Keystore = (*staticKeystore)(nil)
var _ testtypes.Evaluator[types.Keystore] = (*staticKeystore)(nil)

type staticKeystoreConfig struct {
	Account libocr.Account
	encoded []byte
	signed  []byte
}

type staticKeystore struct {
	staticKeystoreConfig
}

func (s staticKeystore) Accounts(ctx context.Context) (accounts []string, err error) {
	return []string{string(s.Account)}, nil
}

func (s staticKeystore) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	if string(s.Account) != id {
		return nil, fmt.Errorf("expected id %q but got %q", s.Account, id)
	}
	if !bytes.Equal(s.encoded, data) {
		return nil, fmt.Errorf("expected encoded data %x but got %x", s.encoded, data)
	}
	return s.signed, nil
}

func (s staticKeystore) Evaluate(ctx context.Context, other types.Keystore) error {
	accounts, err := s.Accounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}
	if len(accounts) != 1 {
		return fmt.Errorf("expected 1 account but got %d", len(accounts))
	}
	if accounts[0] != string(s.Account) {
		return fmt.Errorf("expected account %q but got %q", s.Account, accounts[0])
	}

	signed, err := other.Sign(ctx, string(s.Account), s.encoded)
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}
	if !bytes.Equal(s.signed, signed) {
		return fmt.Errorf("expected signed data %x but got %x", s.signed, signed)
	}
	return nil
}

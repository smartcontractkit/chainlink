package types

import (
	"context"
)

type Keystore interface {
	Accounts(ctx context.Context) (accounts []string, err error)
	// Sign returns data signed by account.
	// nil data can be used as a no-op to check for account existence.
	Sign(ctx context.Context, account string, data []byte) (signed []byte, err error)
}

type ErrorLog interface {
	SaveError(ctx context.Context, msg string) error
}

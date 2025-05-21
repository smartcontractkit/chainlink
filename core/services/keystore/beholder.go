package keystore

import (
	"context"
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

func BuildBeholderAuth(ctx context.Context, keyStore CSA) (authHeaders map[string]string, pubKeyHex string, err error) {
	csaKey, err := GetDefault(ctx, keyStore)
	if err != nil {
		return nil, "", err
	}

	authHeaders, err = beholder.NewAuthHeaders(csaKey)
	if err != nil {
		return
	}
	pubKeyHex = hex.EncodeToString(csaKey.PublicKey)
	return
}

package starkkey_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/caigo"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

func TestKeyStoreAdapter(t *testing.T) {
	var (
		starknetSenderAddr = "legit"
	)

	starkKey, err := starkkey.New()
	require.NoError(t, err)

	lk := starkkey.NewLooppKeystore(func(id string) (starkkey.Key, error) {
		if id != starknetSenderAddr {
			return starkkey.Key{}, fmt.Errorf("error getting key for sender %s: %w", id, ErrSenderNoExist)
		}
		return starkKey, nil
	})
	adapter := starkkey.NewKeystoreAdapter(lk)
	// test that adapter implements the loopp spec. signing nil data should not error
	// on existing sender id
	signed, err := adapter.Loopp().Sign(context.Background(), starknetSenderAddr, nil)
	require.Nil(t, signed)
	require.NoError(t, err)

	signed, err = adapter.Loopp().Sign(context.Background(), "not an address", nil)
	require.Nil(t, signed)
	require.Error(t, err)

	hash, err := caigo.Curve.PedersenHash([]*big.Int{big.NewInt(42)})
	require.NoError(t, err)
	r, s, err := adapter.Sign(context.Background(), starknetSenderAddr, hash)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.NotNil(t, s)

	pubx, puby, err := caigo.Curve.PrivateToPoint(starkKey.ToPrivKey())
	require.NoError(t, err)
	require.True(t, caigo.Curve.Verify(hash, r, s, pubx, puby))
}

var ErrSenderNoExist = errors.New("sender does not exist")

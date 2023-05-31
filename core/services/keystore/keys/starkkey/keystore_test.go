package starkkey_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/caigo"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

func TestKeyStoreAdapter(t *testing.T) {
	var (
		starknetPK         = generateTestKey(t)
		starknetSenderAddr = "legit"
	)

	lk := starkkey.NewLooppKeystore(func(id string) (*big.Int, error) {
		if id != starknetSenderAddr {
			return nil, fmt.Errorf("error getting key for sender %s: %w", id, ErrSenderNoExist)
		}
		return starknetPK, nil
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

	pubx, puby, err := caigo.Curve.PrivateToPoint(starknetPK)
	require.NoError(t, err)
	require.True(t, caigo.Curve.Verify(hash, r, s, pubx, puby))
}

func generateTestKey(t *testing.T) *big.Int {
	// sadly generating a key can fail, but it should happen infrequently
	// best effort here to  avoid flaky tests
	var generatorDuration = 1 * time.Second
	d, exists := t.Deadline()
	if exists {
		generatorDuration = time.Until(d) / 2
	}
	timer := time.NewTicker(generatorDuration)
	defer timer.Stop()
	var key *big.Int
	var generationErr error

	generated := func() bool {
		select {
		case <-timer.C:
			key = nil
			generationErr = fmt.Errorf("failed to generate test key in allotted time")
			return true
		default:
			key, generationErr = caigo.Curve.GetRandomPrivateKey()
			if generationErr == nil {
				return true
			}
		}
		return false
	}

	//nolint:all
	for !generated() {
	}

	require.NoError(t, generationErr)
	return key
}

var ErrSenderNoExist = errors.New("sender does not exist")

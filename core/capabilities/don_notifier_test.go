package capabilities_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

func TestDonNotifier_WaitForDon(t *testing.T) {
	notifier := capabilities.NewDonNotifier()
	don := commoncap.DON{
		ID: 1,
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		notifier.NotifyDonSet(don)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err := notifier.WaitForDon(ctx)
	require.NoError(t, err)
	assert.Equal(t, don, result)

	result, err = notifier.WaitForDon(ctx)
	require.NoError(t, err)
	assert.Equal(t, don, result)
}

func TestDonNotifier_WaitForDon_ContextTimeout(t *testing.T) {
	notifier := capabilities.NewDonNotifier()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := notifier.WaitForDon(ctx)
	require.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

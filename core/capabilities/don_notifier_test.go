package capabilities_test

import (
	"context"
	"testing"

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
	notifyCh := make(chan struct{})
	setCh := make(chan struct{})

	go func() {
		<-notifyCh
		notifier.NotifyDonSet(don)
		close(setCh)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	close(notifyCh)
	result, err := notifier.WaitForDon(ctx)
	require.NoError(t, err)
	assert.Equal(t, don, result)

	// Verify that a second read returns same DON without ever waiting
	<-setCh
	cancel()
	result, err = notifier.WaitForDon(ctx)
	require.NoError(t, err) // should not error because the value is cached
	assert.Equal(t, don, result)
}

func TestDonNotifier_WaitForDon_ContextTimeout(t *testing.T) {
	notifier := capabilities.NewDonNotifier()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := notifier.WaitForDon(ctx)
	require.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestDonNotifier_DonUpdate(t *testing.T) {
	notifier := capabilities.NewDonNotifier()
	notifyChs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
	}
	setCh := make(chan struct{})

	// Set the first DON
	don1 := commoncap.DON{
		ID: 1,
	}
	go func() {
		<-notifyChs[0]
		notifier.NotifyDonSet(don1)
		close(setCh)
	}()

	// Update to second DON
	don2 := commoncap.DON{
		ID: 2,
	}
	go func() {
		<-notifyChs[1]
		notifier.NotifyDonSet(don2)
		close(setCh)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	close(notifyChs[0])
	result, err := notifier.WaitForDon(ctx)
	require.NoError(t, err)
	assert.Equal(t, don1, result)

	close(notifyChs[1])
	<-setCh
	result, err = notifier.WaitForDon(ctx)
	require.NoError(t, err)
	assert.Equal(t, don2, result)
}

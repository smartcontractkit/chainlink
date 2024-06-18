package bridges_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/bridges/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestBridgeCache_Type(t *testing.T) {
	t.Parallel()

	t.Run("loads bridge from data source - single", func(t *testing.T) {
		t.Parallel()

		mORM := new(mocks.ORM)
		lggr, _ := logger.NewLogger()
		cache := bridges.NewCache(mORM, lggr, bridges.DefaultUpsertInterval)

		bridge := bridges.BridgeName("test")
		expected := bridges.BridgeType{
			Name: bridge,
		}

		// first call to find should fallthrough to data source
		mORM.On("FindBridge", mock.Anything, bridge).Return(expected, nil)
		result, err := cache.FindBridge(context.Background(), bridge)

		require.NoError(t, err)
		assert.Equal(t, expected, result)

		// calling find again should return from cache
		result, err = cache.FindBridge(context.Background(), bridge)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("loads bridge from data source - multiple", func(t *testing.T) {
		t.Parallel()

		mORM := new(mocks.ORM)
		lggr, _ := logger.NewLogger()
		cache := bridges.NewCache(mORM, lggr, bridges.DefaultUpsertInterval)

		ctx := context.Background()
		nameA := bridges.BridgeName("A")
		nameB := bridges.BridgeName("B")
		nameC := bridges.BridgeName("C")

		initialExpected := []bridges.BridgeType{
			{Name: nameA},
			{Name: nameB},
		}
		finalExpected := append(initialExpected, bridges.BridgeType{Name: nameC})

		mORM.On("FindBridges", mock.Anything, []bridges.BridgeName{nameA, nameB}).Return(initialExpected, nil)
		types, err := cache.FindBridges(ctx, []bridges.BridgeName{nameA, nameB})

		require.NoError(t, err)
		assert.Equal(t, initialExpected, types)

		// second call to FindBridges only includes the value not in the cache
		mORM.On("FindBridges", mock.Anything, []bridges.BridgeName{nameC}).
			Return([]bridges.BridgeType{{Name: nameC}}, nil)

		types, err = cache.FindBridges(ctx, []bridges.BridgeName{nameA, nameB, nameC})

		require.NoError(t, err)
		assert.Equal(t, finalExpected, types)
	})

	t.Run("creates, updates, and deletes bridge", func(t *testing.T) {
		t.Parallel()

		mORM := new(mocks.ORM)
		lggr, _ := logger.NewLogger()
		cache := bridges.NewCache(mORM, lggr, bridges.DefaultUpsertInterval)

		ctx := context.Background()
		bridge := bridges.BridgeName("test")
		expected := &bridges.BridgeType{
			Name:          bridge,
			Confirmations: 42,
		}

		mORM.On("CreateBridgeType", mock.Anything, expected).Return(nil)
		assert.NoError(t, cache.CreateBridgeType(ctx, expected))

		result, err := cache.FindBridge(ctx, bridge)

		require.NoError(t, err)
		assert.Equal(t, *expected, result)

		btr := &bridges.BridgeTypeRequest{
			Confirmations: 21,
		}

		mORM.On("UpdateBridgeType", mock.Anything, expected, btr).Return(nil).Run(func(args mock.Arguments) {
			btp := args.Get(1).(*bridges.BridgeType)
			req := args.Get(2).(*bridges.BridgeTypeRequest)

			btp.Confirmations = req.Confirmations
		})
		require.NoError(t, cache.UpdateBridgeType(ctx, expected, btr))
		assert.Equal(t, btr.Confirmations, expected.Confirmations)

		result, err = cache.FindBridge(ctx, bridge)

		require.NoError(t, err)
		assert.Equal(t, *expected, result)

		mORM.On("DeleteBridgeType", mock.Anything, expected).Return(nil)
		require.NoError(t, cache.DeleteBridgeType(ctx, expected))

		// bridge type is removed from cache so call to find fallsback to the data store
		mORM.On("FindBridge", mock.Anything, bridge).Return(bridges.BridgeType{}, errors.New("not found"))
		_, err = cache.FindBridge(ctx, bridge)
		require.NotNil(t, err)
	})
}

func TestBridgeCache_Response(t *testing.T) {
	t.Parallel()

	t.Run("loads response from data source", func(t *testing.T) {
		t.Parallel()

		mORM := new(mocks.ORM)
		lggr, _ := logger.NewLogger()
		cache := bridges.NewCache(mORM, lggr, bridges.DefaultUpsertInterval)

		ctx := context.Background()
		dotId := "test"
		specId := int32(42)
		responseData := []byte("test")

		mORM.On("GetCachedResponseWithFinished", mock.Anything, dotId, specId, time.Second).
			Return(responseData, time.Now(), nil)

		response, err := cache.GetCachedResponse(ctx, dotId, specId, time.Second)

		require.NoError(t, err)
		assert.Equal(t, responseData, response)

		response, err = cache.GetCachedResponse(ctx, dotId, specId, time.Second)

		require.NoError(t, err)
		assert.Equal(t, responseData, response)
	})

	t.Run("async upserts bridge response", func(t *testing.T) {
		t.Parallel()

		mORM := new(mocks.ORM)
		lggr, _ := logger.NewLogger()
		cache := bridges.NewCache(mORM, lggr, bridges.DefaultUpsertInterval)

		t.Cleanup(func() {
			require.NoError(t, cache.Close())
		})

		require.NoError(t, cache.Start(context.Background()))

		ctx := context.Background()
		dotId := "test"
		specId := int32(42)
		expected := []byte("test")
		chBulkUpsertCalled := make(chan struct{}, 1)

		mORM.On("BulkUpsertBridgeResponse", mock.Anything, mock.Anything).Return(nil).Run(func(_ mock.Arguments) {
			chBulkUpsertCalled <- struct{}{}
		})

		require.NoError(t, cache.UpsertBridgeResponse(ctx, dotId, specId, expected))

		response, err := cache.GetCachedResponse(ctx, dotId, specId, time.Second)

		require.NoError(t, err)
		assert.Equal(t, expected, response)

		timer := time.NewTimer(2 * bridges.DefaultUpsertInterval)

		select {
		case <-chBulkUpsertCalled:
			timer.Stop()
		case <-timer.C:
			timer.Stop()

			t.Log("upsert interval exceeded without expected upsert")
			t.FailNow()
		}
	})
}

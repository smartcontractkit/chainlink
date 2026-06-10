package shared

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type codeAtStub struct {
	calls   int
	readyOn int
	code    []byte
	err     error
}

func (s *codeAtStub) CodeAt(context.Context, common.Address, *big.Int) ([]byte, error) {
	s.calls++
	if s.err != nil && s.calls == 1 {
		return nil, s.err
	}
	if s.readyOn > 0 && s.calls < s.readyOn {
		return nil, nil
	}
	return s.code, nil
}

func TestWaitForContractCode(t *testing.T) {
	t.Parallel()

	addr := common.HexToAddress("0x1234")
	bytecode := []byte{0x60, 0x80}

	t.Run("returns immediately when code is present", func(t *testing.T) {
		t.Parallel()
		stub := &codeAtStub{readyOn: 1, code: bytecode}

		err := retryUntilContractCode(context.Background(), stub, addr)
		require.NoError(t, err)
		assert.Equal(t, 1, stub.calls)
	})

	t.Run("retries until code appears", func(t *testing.T) {
		t.Parallel()
		stub := &codeAtStub{readyOn: 2, code: bytecode}

		err := retryUntilContractCode(context.Background(), stub, addr)
		require.NoError(t, err)
		assert.Equal(t, 2, stub.calls)
	})

	t.Run("returns context error when code never appears", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), 750*time.Millisecond)
		t.Cleanup(cancel)

		err := retryUntilContractCode(ctx, &codeAtStub{readyOn: 100}, addr)
		require.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("retries RPC errors", func(t *testing.T) {
		t.Parallel()
		stub := &codeAtStub{
			readyOn: 2,
			code:    bytecode,
			err:     errors.New("rpc unavailable"),
		}

		err := retryUntilContractCode(context.Background(), stub, addr)
		require.NoError(t, err)
		assert.Equal(t, 2, stub.calls)
	})
}

package legacygasstation

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/jsonrpc"
)

type testRequestHandler struct {
	chainID *utils.Big
}

func (h testRequestHandler) SendTransaction(*gin.Context, types.SendTransactionRequest) (*types.SendTransactionResponse, *jsonrpc.Error) {
	return nil, nil
}

func (h testRequestHandler) CCIPChainSelector() *utils.Big {
	return h.chainID
}

func TestRegisterDeregister(t *testing.T) {
	lggr := logger.TestLogger(t)
	rr := NewRequestRouter(lggr)
	chainID := utils.NewBigI(5)
	h := testRequestHandler{
		chainID: chainID,
	}
	err := rr.RegisterHandler(h)
	require.NoError(t, err)

	ctx := gin.Context{}
	req := types.SendTransactionRequest{
		SourceChainID: chainID.ToInt().Uint64(),
	}
	_, jsonrpcErr := rr.SendTransaction(&ctx, req)
	require.Nil(t, jsonrpcErr)

	rr.DeregisterHandler(chainID)

	_, jsonrpcErr = rr.SendTransaction(&ctx, req)
	require.NotNil(t, jsonrpcErr)
	require.Equal(t, jsonrpc.InvalidRequestError, jsonrpcErr.Code)
}

func TestRegisterInvalidChainID(t *testing.T) {
	lggr := logger.TestLogger(t)
	rr := NewRequestRouter(lggr)
	h := testRequestHandler{
		chainID: nil,
	}
	err := rr.RegisterHandler(h)
	require.Error(t, err)
}

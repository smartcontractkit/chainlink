package vault

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	vault_api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type mockGatewaySender struct {
	resp *jsonrpc.Response[json.RawMessage]
}

func (m *mockGatewaySender) SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error {
	m.resp = resp
	return nil
}

func TestVault_Handler_ListSecretIdentifiers_OwnerEmpty(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*Request]()
	hdlr := requests.NewHandler[*Request, *Response](lggr, store, clock, expiry)
	service := NewService(lggr, clock, expiry, hdlr)
	servicetest.Run(t, service)

	gs := &mockGatewaySender{}
	handler, err := NewHandler(service, gs, lggr)
	require.NoError(t, err)

	req := vault.ListSecretIdentifiersRequest{
		RequestId: "request-id",
	}
	reqb, err := json.Marshal(&req)
	require.NoError(t, err)

	rmsg := json.RawMessage(reqb)
	msg := &jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "test_id",
		Method:  vault_api.MethodSecretsList,
		Params:  &rmsg,
	}
	err = handler.HandleGatewayMessage(context.Background(), "test_gateway", msg)
	require.NoError(t, err)

	assert.Nil(t, gs.resp.Result)
	assert.Contains(t, gs.resp.Error.Message, "owner must not be empty")
}

func TestVault_Handler_ListSecretIdentifiers_RequestIdEmpty(t *testing.T) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*Request]()
	hdlr := requests.NewHandler[*Request, *Response](lggr, store, clock, expiry)
	service := NewService(lggr, clock, expiry, hdlr)
	servicetest.Run(t, service)

	gs := &mockGatewaySender{}
	handler, err := NewHandler(service, gs, lggr)
	require.NoError(t, err)

	req := vault.ListSecretIdentifiersRequest{
		RequestId: "",
		Owner:     "owner-id",
	}
	reqb, err := json.Marshal(&req)
	require.NoError(t, err)

	rmsg := json.RawMessage(reqb)
	msg := &jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "test_id",
		Method:  vault_api.MethodSecretsList,
		Params:  &rmsg,
	}
	err = handler.HandleGatewayMessage(context.Background(), "test_gateway", msg)
	require.NoError(t, err)

	assert.Nil(t, gs.resp.Result)
	assert.Contains(t, gs.resp.Error.Message, "request ID must not be empty")
}

func TestVault_Handler_ListSecretIdentifiers(t *testing.T) {
	lggr := logger.TestLogger(t)
	listResp := &vault.ListSecretIdentifiersResponse{
		Identifiers: []*vault.SecretIdentifier{},
	}
	requestID := "request-id"
	respb, err := json.Marshal(listResp)
	require.NoError(t, err)
	resp := &Response{
		ID:      requestID,
		Payload: respb,
		Format:  "json",
	}

	service := &mockService{resp: resp}

	gs := &mockGatewaySender{}
	handler, err := NewHandler(service, gs, lggr)
	require.NoError(t, err)

	req := vault.ListSecretIdentifiersRequest{
		RequestId: requestID,
		Owner:     "owner-id",
	}
	reqb, err := json.Marshal(&req)
	require.NoError(t, err)

	rmsg := json.RawMessage(reqb)
	msg := &jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "test_id",
		Method:  vault_api.MethodSecretsList,
		Params:  &rmsg,
	}
	err = handler.HandleGatewayMessage(context.Background(), "test_gateway", msg)
	require.NoError(t, err)

	assert.Nil(t, gs.resp.Error)

	gotListResp := &vault.ListSecretIdentifiersResponse{}
	err = json.Unmarshal(*gs.resp.Result, gotListResp)
	require.NoError(t, err)
	assert.True(t, proto.Equal(listResp, gotListResp))
}

type mockService struct {
	*Service
	resp *Response
}

func (m *mockService) ListSecretIdentifiers(ctx context.Context, req *vault.ListSecretIdentifiersRequest) (*Response, error) {
	return m.resp, nil
}

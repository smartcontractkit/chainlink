package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

type methodBinding struct {
	contractName string
	method       string
	client       evmclient.Client
	codec        commontypes.Codec
}

var _ readBinding = &methodBinding{}

func (m *methodBinding) SetCodec(codec commontypes.RemoteCodec) {
	m.codec = codec
}

func (m *methodBinding) Register(_ context.Context, _ common.Address) error {
	return nil
}

func (m *methodBinding) Unregister(_ context.Context, _ common.Address) error {
	return nil
}

func (m *methodBinding) UnregisterAll(_ context.Context) error {
	return nil
}

func (m *methodBinding) QueryOne(_ context.Context, _ common.Address, _ query.Filter, _ query.LimitAndSort, _ any) ([]commontypes.Sequence, error) {
	return nil, nil
}

func (m *methodBinding) GetLatestValue(ctx context.Context, address common.Address, params, returnValue any) error {
	data, err := m.codec.Encode(ctx, params, wrapItemType(m.contractName, m.method, true))
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		To:   &address,
		From: address,
		Data: data,
	}

	bytes, err := m.client.CallContract(ctx, callMsg, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return m.codec.Decode(ctx, bytes, returnValue, wrapItemType(m.contractName, m.method, false))
}

func (m *methodBinding) Bind(_ context.Context, _ common.Address) error {
	return nil
}

func (m *methodBinding) UnBind(_ context.Context, _ common.Address) error {
	return nil
}

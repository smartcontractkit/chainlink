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
	address      common.Address
	contractName string
	method       string
	client       evmclient.Client
	codec        commontypes.Codec
	bound        bool
}

var _ readBinding = &methodBinding{}

func (m *methodBinding) SetCodec(codec commontypes.RemoteCodec) {
	m.codec = codec
}

func (m *methodBinding) Register(_ context.Context) error {
	return nil
}

func (m *methodBinding) Unregister(_ context.Context) error {
	return nil
}

func (m *methodBinding) UnregisterAll(_ context.Context) error {
	return nil
}

func (m *methodBinding) GetLatestValue(ctx context.Context, params, returnValue any) error {
	if !m.bound {
		return fmt.Errorf("%w: method not bound", commontypes.ErrInvalidType)
	}

	methodKey := formatKey(m.contractName, m.method)
	data, err := m.codec.Encode(ctx, params, wrapItemType(methodKey, true))
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		To:   &m.address,
		From: m.address,
		Data: data,
	}

	bytes, err := m.client.CallContract(ctx, callMsg, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return m.codec.Decode(ctx, bytes, returnValue, wrapItemType(methodKey, false))
}

func (m *methodBinding) QueryKey(_ context.Context, _ query.KeyFilter, _ query.LimitAndSort, _ any) ([]commontypes.Sequence, error) {
	return nil, nil
}

func (m *methodBinding) Bind(_ context.Context, binding commontypes.BoundContract) error {
	m.address = common.HexToAddress(binding.Address)
	m.bound = true
	return nil
}

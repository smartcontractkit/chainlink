package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

var ErrNoContractExists = errors.New("contract does not exist at address")

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

func (m *methodBinding) Bind(ctx context.Context, binding commontypes.BoundContract) error {
	addr := common.HexToAddress(binding.Address)

	// check for contract byte code at the latest block and provided address
	byteCode, err := m.client.CodeAt(ctx, addr, nil)
	if err != nil {
		return err
	}

	if len(byteCode) < 0 {
		return fmt.Errorf("%w: %s", ErrNoContractExists, addr)
	}

	m.address = addr
	m.bound = true

	return nil
}

func (m *methodBinding) Register(_ context.Context) error {
	return nil
}

func (m *methodBinding) Unregister(_ context.Context) error {
	return nil
}

func (m *methodBinding) GetLatestValue(ctx context.Context, params, returnValue any) error {
	if !m.bound {
		return fmt.Errorf("%w: method not bound", commontypes.ErrInvalidType)
	}

	data, err := m.codec.Encode(ctx, params, wrapItemType(m.contractName, m.method, true))
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

	return m.codec.Decode(ctx, bytes, returnValue, wrapItemType(m.contractName, m.method, false))
}

func (m *methodBinding) QueryKey(_ context.Context, _ query.KeyFilter, _ query.LimitAndSort, _ any) ([]commontypes.Sequence, error) {
	return nil, nil
}

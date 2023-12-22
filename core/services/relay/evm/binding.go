package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type readBinding interface {
	GetLatestValue(ctx context.Context, params any) ([]byte, error)
	Bind(binding commontypes.BoundContract) error
	SetCodec(codec commontypes.Codec)
	Register() error
	Unregister() error
}

type methodBinding struct {
	address      common.Address
	contractName string
	method       string
	client       evmclient.Client
	codec        commontypes.Codec
	bound        bool
}

var _ readBinding = &methodBinding{}

func (m *methodBinding) SetCodec(codec commontypes.Codec) {
	m.codec = codec
}

func (m *methodBinding) Register() error {
	return nil
}

func (m *methodBinding) Unregister() error {
	return nil
}

func (m *methodBinding) GetLatestValue(ctx context.Context, params any) ([]byte, error) {
	if !m.bound {
		return nil, fmt.Errorf("%w: method not bound", commontypes.ErrInvalidType)
	}

	data, err := m.codec.Encode(ctx, params, wrapItemType(m.contractName, m.method, true))
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		To:   &m.address,
		From: m.address,
		Data: data,
	}

	return m.client.CallContract(ctx, callMsg, nil)
}

func (m *methodBinding) Bind(binding commontypes.BoundContract) error {
	m.address = common.HexToAddress(binding.Address)
	m.bound = true
	return nil
}

type eventBinding struct {
	address      common.Address
	contractName string
	eventName    string
	lp           logpoller.LogPoller
	hash         common.Hash
	codec        commontypes.Codec
	pending      bool
	subscribed   bool
	bound        bool
}

func (e *eventBinding) SetCodec(codec commontypes.Codec) {
	e.codec = codec
}

func (e *eventBinding) Register() error {
	e.subscribed = true
	if !e.bound {
		return nil
	}

	if err := e.lp.RegisterFilter(logpoller.Filter{
		Name:      wrapItemType(e.contractName, e.eventName, false),
		EventSigs: evmtypes.HashArray{e.hash},
		Addresses: evmtypes.AddressArray{e.address},
	}); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return nil
}

func (e *eventBinding) Unregister() error {
	e.subscribed = false
	if !e.bound {
		return nil
	}

	if err := e.lp.UnregisterFilter(wrapItemType(e.contractName, e.eventName, false)); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return nil
}

var _ readBinding = &eventBinding{}

func (e *eventBinding) GetLatestValue(_ context.Context, _ any) ([]byte, error) {
	if !e.bound {
		return nil, fmt.Errorf("%w: event not bound", commontypes.ErrInvalidType)
	}

	confs := logpoller.Finalized
	if e.pending {
		confs = logpoller.Unconfirmed
	}
	log, err := e.lp.LatestLogByEventSigWithConfs(e.hash, e.address, confs)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "not found") || strings.Contains(errStr, "no rows") {
			return nil, fmt.Errorf("%w: %w", commontypes.ErrNotFound, err)
		}
		return nil, err
	}

	return log.Data, nil
}

func (e *eventBinding) Bind(binding commontypes.BoundContract) error {
	wasSubscribed := e.subscribed
	if wasSubscribed {
		if err := e.Unregister(); err != nil {
			return err
		}
	}
	e.address = common.HexToAddress(binding.Address)
	e.pending = binding.Pending
	e.bound = true

	if wasSubscribed {
		return e.Register()
	}
	return nil
}

package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

type NoContractExistsError struct {
	Address common.Address
}

func (e NoContractExistsError) Error() string {
	return fmt.Sprintf("contract does not exist at address: %s", e.Address)
}

type methodBinding struct {
	lggr                 logger.Logger
	ht                   logpoller.HeadTracker
	address              common.Address
	contractName         string
	method               string
	client               evmclient.Client
	codec                commontypes.Codec
	bound                bool
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations
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

	if len(byteCode) == 0 {
		return NoContractExistsError{Address: addr}
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

func (m *methodBinding) GetLatestValue(ctx context.Context, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	if !m.bound {
		return fmt.Errorf("%w: method not bound", commontypes.ErrInvalidType)
	}

	data, err := m.codec.Encode(ctx, params, WrapItemType(m.contractName, m.method, true))
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		To:   &m.address,
		From: m.address,
		Data: data,
	}

	block, err := m.blockNumberFromConfidence(ctx, confidenceLevel)
	if err != nil {
		return err
	}

	bytes, err := m.client.CallContract(ctx, callMsg, block)
	if err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return m.codec.Decode(ctx, bytes, returnVal, WrapItemType(m.contractName, m.method, false))
}

func (m *methodBinding) QueryKey(_ context.Context, _ query.KeyFilter, _ query.LimitAndSort, _ any) ([]commontypes.Sequence, error) {
	return nil, nil
}

func (m *methodBinding) blockNumberFromConfidence(ctx context.Context, confidenceLevel primitives.ConfidenceLevel) (*big.Int, error) {
	confirmations, err := confidenceToConfirmations(m.confirmationsMapping, confidenceLevel)
	if err != nil {
		err = fmt.Errorf("%w for contract: %s, method: %s", err, m.contractName, m.method)
		if confidenceLevel == primitives.Unconfirmed {
			m.lggr.Errorf("%v, now falling back to default contract call behaviour that calls latest state", err)
			return nil, nil
		}
		return nil, err
	}

	_, finalized, err := m.ht.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return nil, err
	}

	if confirmations == evmtypes.Finalized {
		return big.NewInt(finalized.Number), nil
	} else if confirmations == evmtypes.Unconfirmed {
		return nil, nil
	}

	return nil, fmt.Errorf("unknown evm confirmations: %v for contract: %s, method: %s", confirmations, m.contractName, m.method)
}
